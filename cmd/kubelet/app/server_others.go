//go:build !windows
// +build !windows

/*
Copyright 2021 The Kubernetes Authors.

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

package app

import (
	"fmt"
	"os"
)

func checkPermissions() error {
	if uid := os.Getuid(); uid != 0 {
		return fmt.Errorf("kubelet needs to run as uid `0`. It is being run as %d", uid)
	}

	if err := checkInitialUserNamespace(); err != nil {
		return err
	}

	return nil
}

// checkInitialUserNamespace checks if kubelet is running in the initial user namespace.
// http://man7.org/linux/man-pages/man7/user_namespaces.7.html
func checkInitialUserNamespace() error {
	uidMap, err := ioutil.ReadFile("/proc/self/uid_map")
	if err != nil {
		return err
	}

	if strings.TrimSpace(string(uidMap)) != "0\t0\t4294967295" {
		return fmt.Errorf("kubelet is not running in the initial user namespace")
	}

	return nil
}
