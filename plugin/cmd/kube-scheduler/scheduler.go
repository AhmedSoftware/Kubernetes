/*
Copyright 2014 The Kubernetes Authors All rights reserved.

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

package main

import (
	"runtime"

	"github.com/golang/glog"

	"k8s.io/kubernetes/pkg/healthz"
	"k8s.io/kubernetes/pkg/util"
	"k8s.io/kubernetes/plugin/cmd/kube-scheduler/app"
)

func init() {
	healthz.DefaultHealthz()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cmd := app.NewSchedulerServerCommand()

	util.InitLogs()
	defer util.FlushLogs()

	if err := cmd.Execute(); err != nil {
		glog.Fatal(err)
	}
}
