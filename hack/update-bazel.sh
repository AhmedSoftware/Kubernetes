#!/usr/bin/env bash
# Copyright 2016 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

export KUBE_ROOT=$(dirname "${BASH_SOURCE}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::util::ensure-gnu-sed

# Remove generated files prior to running kazel.
# TODO(spxtr): Remove this line once Bazel is the only way to build.
rm -f "${KUBE_ROOT}/pkg/generated/openapi/zz_generated.openapi.go"

# Ensure that we find the binaries we build before anything else.
export GOBIN="${KUBE_OUTPUT_BINPATH}"
PATH="${GOBIN}:${PATH}"

# Install tools we need, but only from vendor/...
go install ./vendor/github.com/bazelbuild/bazel-gazelle/cmd/gazelle

go install ./vendor/github.com/kubernetes/repo-infra/kazel

cd "${KUBE_ROOT}"
go get github.com/bazelbuild/buildtools/buildozer
touch "${KUBE_ROOT}/vendor/BUILD"

dozer() {
  # returns 0 on successful mod and 3 on no change
  buildozer "$@" && return 0 || [[ "$?" == 3 ]] && return 0
  return 1
}


add-generated-deps() {
  if which -s bazel; then
    bazel build //build:all-generated-sources
    for path in $(find -H bazel-genfiles -iname *.go.deps); do
      deps="${path#bazel-genfiles/build/}" # rel/path/to/foo.go.deps
      pkg="${deps%/*}" # rel/path/to
      if [[ ! -f "$pkg/BUILD" && ! -f "$pkg/BUILD.bazel" ]]; then
        echo "$pkg: no BUILD.bazel file"
        continue
      fi
      deps=($(cat "$path"))
      deps="$(IFS=' '  ; echo "${deps[*]}")"
      dozer "add deps $deps" "//${pkg}:go_default_library"
    done
  fi
}

gazelle-fix() {
  if ! which -s bazel && [[ "${CALLED_FROM_MAIN_MAKEFILE:-}" == "" ]]; then
    echo "Please use "make update" or install http://bazel.build"
  fi
  gazelle fix \
      -build_file_name=BUILD,BUILD.bazel \
      -external=vendored \
      -proto=legacy \
      -mode=fix
  # gazelle gets confused by our staging/ directory, prepending an extra
  # "k8s.io/kubernetes/staging/src" to the import path.
  # gazelle won't follow the symlinks in vendor/, so we can't just exclude
  # staging/. Instead we just fix the bad paths with sed.
  find staging -name BUILD -o -name BUILD.bazel | \
    xargs ${SED} -i 's|\(importpath = "\)k8s.io/kubernetes/staging/src/\(.*\)|\1\2|'

  # Add deps for any generated dependencies
  if which -s bazel; then
    add-generated-deps
  fi
}

gazelle-fix

kazel

update-k8s-gengo() {
  local name="$1" # something like deepcopy
  local pkg="$2"  # something like k8s.io/code-generator/cmd/deepcopy-gen
  local out="$3"  # something like zz_generated.deepcopy.go
  local match="$4"  # something like +k8s:deepcopy-gen=

  local all_rule="k8s_${name}_all"  # k8s_deepcopy_all, which generates out for matching packages
  local rule="k8s_${name}" # k8s_deepcopy, which copies out the file for a particular package

  # look for packages that contain match
  echo "Looking for packages with a $match comment in go code..."
  want=($(find . -name *.go | grep -v "$pkg" | (xargs grep -l "$match" || true) | (xargs -n 1 dirname || true) | sort -u | $SED -e 's|./staging/src/|./vendor/|' | (xargs go list || true) | $SED -e 's|k8s.io/kubernetes/||' | sort -u))
  echo "[$(IFS=$'\n ' ; echo "${want[*]}" | sed -e 's|\(^.*$\)|  "\1",|')]"

  # Ensure that k8s_deepcopy_all() rule exists
  if ! grep -q "${all_rule}(" build/BUILD; then
    echo Adding $all_rule to build/BULID
    touch build/BUILD
    echo "$all_rule(name=\"${name}-sources\")" >> build/BUILD
  else
    echo $all_rule found in build/BUILD
  fi
  dozer "new_load //build:deepcopy.bzl $all_rule" //build:__pkg__
  dozer "set packages ${want[*]}" "//build:$name-sources"
  have=$(find . -name BUILD -or -name BUILD.bazel | (xargs grep -l "$rule(" || true))
  if [[ -n "$have" ]]; then
    echo Deleting existing "$rule" commands... $have
    case "$(uname -s)" in
      Darwin*)
        $SED -i -e "/^$rule/d" $have
        ;;
      *)
        $SED -i -e "/^$rule/d" $have
        ;;
    esac
  fi
  echo "Adding $rule() rule"
  for w in "${want[@]}"; do
    if [[ $w == */$pkg ]]; then
      echo skipping $w...
      continue
    fi
    if [[ -f $w/BUILD.bazel ]]; then
      echo "$rule(outs=[\"${out}\"])" >> $w/BUILD
      dozer "new_load //build:deepcopy.bzl $rule" //$w:__pkg__
    elif  [[ -f $w/BUILD ]]; then
      echo "$rule(outs=[\"${out}\"])" >> $w/BUILD
      dozer "new_load //build:deepcopy.bzl $rule" //$w:__pkg__
    else
      echo cannot find build file for $w
      continue
    fi
    dozer "new_load //build:deepcopy.bzl $rule" //$w:__pkg__
  done
  echo "Deleting $out files"
  find . -iname "${out}" | xargs rm
}
update-k8s-gengo deepcopy k8s.io/code-generator/cmd/deepcopy-gen zz_generated.deepcopy.go '+k8s:deepcopy-gen='
update-k8s-gengo defaulter k8s.io/code-generator/cmd/defaulter-gen zz_generated.defaults.go '+k8s:defaulter-gen='
echo Running gazelle to cleanup any changes...
gazelle-fix
