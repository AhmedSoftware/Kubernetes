#!/bin/bash

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

# This script downloads and installs the Kubernetes client and server binaries.
# It is intended to be called from an extracted Kubernetes release tarball.
#
# We automatically choose the correct client binaries to download.
#
# Options:
#  Set KUBERNETES_SERVER_ARCH to choose the server (Kubernetes cluster)
#  architecture to download:
#    * amd64 [default]
#    * arm
#    * arm64
#    * ppc64le
#
#  Set KUBERNETES_SKIP_CONFIRM to skip the installation confirmation prompt.
#  Set KUBERNETES_RELEASE_URL to choose where to download binaries from.
#    (Default https://storage.googleapis.com/kubernetes-release/release).

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(cd $(dirname "${BASH_SOURCE}")/.. && pwd)

KUBERNETES_RELEASE_URL="${KUBERNETES_RELEASE_URL:-https://storage.googleapis.com/kubernetes-release/release}"

function detect_kube_release() {
  if [[ ! -e "${KUBE_ROOT}/version" ]]; then
    echo "Can't determine Kubernetes release." >&2
    echo "This script should only be run from a prebuilt Kubernetes release." >&2
    echo "Did you mean to use get-kube.sh instead?" >&2
    exit 1
  fi

  KUBERNETES_RELEASE=$(cat "${KUBE_ROOT}/version")
  DOWNLOAD_URL_PREFIX="${KUBERNETES_RELEASE_URL}/${KUBERNETES_RELEASE}"
}

function detect_client_info() {
  local uname=$(uname)
  if [[ "${uname}" == "Darwin" ]]; then
    CLIENT_PLATFORM="darwin"
  elif [[ "${uname}" == "Linux" ]]; then
    CLIENT_PLATFORM="linux"
  else
    echo "Unknown, unsupported platform: (${uname})." >&2
    echo "Supported platforms: Linux, Darwin." >&2
    echo "Bailing out." >&2
    exit 2
  fi

  local machine=$(uname -m)
  if [[ "${machine}" == "x86_64" ]]; then
    CLIENT_ARCH="amd64"
  elif [[ "${machine}" == "i686" ]]; then
    CLIENT_ARCH="386"
  elif [[ "${machine}" == "arm" ]]; then
    CLIENT_ARCH="arm"
  elif [[ "${machine}" == "arm64" ]]; then
    CLIENT_ARCH="arm64"
  else
    echo "Unknown, unsupported architecture (${machine})." >&2
    echo "Supported architectures x86_64, i686, arm, arm64." >&2
    echo "Bailing out." >&2
    exit 3
  fi
}

function md5sum_file() {
  if which md5 >/dev/null 2>&1; then
    md5 -q "$1"
  else
    md5sum "$1" | awk '{ print $1 }'
  fi
}

function sha1sum_file() {
  if which shasum >/dev/null 2>&1; then
    shasum -a1 "$1" | awk '{ print $1 }'
  else
    sha1sum "$1" | awk '{ print $1 }'
  fi
}

function download_tarball() {
  local -r download_path="$1"
  local -r file="$2"
  url="${DOWNLOAD_URL_PREFIX}/${file}"
  mkdir -p "${download_path}"
  if [[ $(which curl) ]]; then
    curl -L "${url}" -o "${download_path}/${file}"
  elif [[ $(which wget) ]]; then
    wget "${url}" -O "${download_path}/${file}"
  else
    echo "Couldn't find curl or wget.  Bailing out." >&2
    exit 4
  fi
  echo
  local md5sum=$(md5sum_file "${download_path}/${file}")
  echo "md5sum(${file})=${md5sum}"
  local sha1sum=$(sha1sum_file "${download_path}/${file}")
  echo "sha1sum(${file})=${sha1sum}"
  echo
  # TODO: add actual verification
}

function extract_tarball() {
  local -r tarfile="$1"
  local -r platform="$2"
  local -r arch="$3"

  platforms_dir="${KUBE_ROOT}/platforms/${platform}/${arch}"
  echo "Extracting ${tarfile} into ${platforms_dir}"
  mkdir -p "${platforms_dir}"
  # Tarball looks like kubernetes/{client,server}/bin/BINARY"
  tar -xzf "${tarfile}" --strip-components 3 -C "${platforms_dir}"
  # Create convenience symlink
  ln -sf "${platforms_dir}" "$(dirname ${tarfile})/bin"
  echo "Add '$(dirname ${tarfile})/bin' to your PATH to use newly-installed binaries."
}

detect_kube_release

SERVER_PLATFORM="linux"
SERVER_ARCH="${KUBERNETES_SERVER_ARCH:-amd64}"
SERVER_TAR="kubernetes-server-${SERVER_PLATFORM}-${SERVER_ARCH}.tar.gz"

detect_client_info
CLIENT_TAR="kubernetes-client-${CLIENT_PLATFORM}-${CLIENT_ARCH}.tar.gz"

echo "Kubernetes release: ${KUBERNETES_RELEASE}"
echo "Server: ${SERVER_PLATFORM}/${SERVER_ARCH}"
echo "Client: ${CLIENT_PLATFORM}/${CLIENT_ARCH}"
echo

# TODO: remove this check and default to true when we stop shipping server
# tarballs in kubernetes.tar.gz
DOWNLOAD_SERVER_TAR=false
if [[ ! -e "${KUBE_ROOT}/server/${SERVER_TAR}" ]]; then
  DOWNLOAD_SERVER_TAR=true
  echo "Will download ${SERVER_TAR} from ${DOWNLOAD_URL_PREFIX}"
fi

# TODO: remove this check and default to true when we stop shipping kubectl
# in kubernetes.tar.gz
DOWNLOAD_CLIENT_TAR=false
if [[ ! -x "${KUBE_ROOT}/platforms/${CLIENT_PLATFORM}/${CLIENT_ARCH}/kubectl" ]]; then
  DOWNLOAD_CLIENT_TAR=true
  echo "Will download and extract ${CLIENT_TAR} from ${DOWNLOAD_URL_PREFIX}"
fi

if [[ "${DOWNLOAD_CLIENT_TAR}" == false && "${DOWNLOAD_SERVER_TAR}" == false ]]; then
  echo "Nothing additional to download."
  exit 0
fi

if [[ -z "${KUBERNETES_SKIP_CONFIRM-}" ]]; then
  echo "Is this ok? [Y]/n"
  read confirm
  if [[ "$confirm" == "n" ]]; then
    echo "Aborting."
    exit 0
  fi
fi

if "${DOWNLOAD_SERVER_TAR}"; then
download_tarball "${KUBE_ROOT}/server" "${SERVER_TAR}"
fi

if "${DOWNLOAD_CLIENT_TAR}"; then
  download_tarball "${KUBE_ROOT}/client" "${CLIENT_TAR}"
  extract_tarball "${KUBE_ROOT}/client/${CLIENT_TAR}" "${CLIENT_PLATFORM}" "${CLIENT_ARCH}"
fi
