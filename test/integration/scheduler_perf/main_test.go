/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package benchmark

import (
	"flag"
	"testing"

	"k8s.io/kubernetes/test/integration/framework"
)

func init() {
	// default logging flags to more reasonable values for this test
	flag.Set("v", "1")
	flag.Set("logtostderr", "false")
	flag.Set("log_dir", ".")
	flag.Parse()
}

func TestMain(m *testing.M) {
	framework.EtcdMain(m.Run)
}
