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

package main

import (
	goflag "flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"

	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	_ "k8s.io/component-base/metrics/prometheus/restclient" // for client metric registration
	_ "k8s.io/component-base/metrics/prometheus/version"    // for version metric registration
	"k8s.io/component-base/version"
	"k8s.io/component-base/version/verflag"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/cluster/ports"
	cadvisortest "k8s.io/kubernetes/pkg/kubelet/cadvisor/testing"
	"k8s.io/kubernetes/pkg/kubelet/cm"
	"k8s.io/kubernetes/pkg/kubelet/cri/remote"
	fakeremote "k8s.io/kubernetes/pkg/kubelet/cri/remote/fake"
	"k8s.io/kubernetes/pkg/kubemark"
	utiltaints "k8s.io/kubernetes/pkg/util/taints"
)

type hollowKubeletConfig struct {
	KubeconfigPath      string
	KubeletPort         int
	KubeletReadOnlyPort int
	NodeName            string
	ContentType         string
	NodeLabels          map[string]string
	RegisterWithTaints  []core.Taint
}

const (
	maxPods     = 110
	podsPerCore = 0
)

func (c *hollowKubeletConfig) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.KubeconfigPath, "kubeconfig", "/kubeconfig/kubeconfig", "Path to kubeconfig file.")
	fs.IntVar(&c.KubeletPort, "kubelet-port", ports.KubeletPort, "Port on which HollowKubelet should be listening.")
	fs.IntVar(&c.KubeletReadOnlyPort, "kubelet-read-only-port", ports.KubeletReadOnlyPort, "Read-only port on which Kubelet is listening.")
	fs.StringVar(&c.NodeName, "name", "fake-node", "Name of this Hollow Node.")
	fs.StringVar(&c.ContentType, "kube-api-content-type", "application/vnd.kubernetes.protobuf", "ContentType of requests sent to apiserver.")
	bindableNodeLabels := cliflag.ConfigurationMap(c.NodeLabels)
	fs.Var(&bindableNodeLabels, "node-labels", "Additional node labels")
	fs.Var(utiltaints.NewTaintsVar(&c.RegisterWithTaints), "register-with-taints", "Register the node with the given list of taints (comma separated \"<key>=<value>:<effect>\"). No-op if register-node is false.")
}

func (c *hollowKubeletConfig) createClientConfigFromFile() (*restclient.Config, error) {
	clientConfig, err := clientcmd.LoadFromFile(c.KubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("error while loading kubeconfig from file %v: %v", c.KubeconfigPath, err)
	}
	config, err := clientcmd.NewDefaultClientConfig(*clientConfig, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error while creating kubeconfig: %v", err)
	}
	config.ContentType = c.ContentType
	config.QPS = 10
	config.Burst = 20
	return config, nil
}

func (c *hollowKubeletConfig) createHollowKubeletOptions() *kubemark.HollowKubletOptions {
	return &kubemark.HollowKubletOptions{
		NodeName:            c.NodeName,
		KubeletPort:         c.KubeletPort,
		KubeletReadOnlyPort: c.KubeletReadOnlyPort,
		MaxPods:             maxPods,
		PodsPerCore:         podsPerCore,
		NodeLabels:          c.NodeLabels,
		RegisterWithTaints:  c.RegisterWithTaints,
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	command := newHollowKubeletCommand()

	// TODO: once we switch everything over to Cobra commands, we can go back to calling
	// cliflag.InitFlags() (by removing its pflag.Parse() call). For now, we have to set the
	// normalize func and add the go flag set by hand.
	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	// cliflag.InitFlags()
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

// newHollowKubeletCommand creates a *cobra.Command object with default parameters
func newHollowKubeletCommand() *cobra.Command {
	s := &hollowKubeletConfig{
		NodeLabels: make(map[string]string),
	}

	cmd := &cobra.Command{
		Use:  "hollow-kubelet",
		Long: "hollow-kubelet pretends to be an ordinary kubelet but doesn't start any containers or mount any volumes",
		Run: func(cmd *cobra.Command, args []string) {
			verflag.PrintAndExitIfRequested()
			run(s)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}
	s.addFlags(cmd.Flags())

	return cmd
}

func run(config *hollowKubeletConfig) {
	// To help debugging, immediately log version
	klog.Infof("Version: %+v", version.Get())

	// create a client to communicate with API server.
	clientConfig, err := config.createClientConfigFromFile()
	if err != nil {
		klog.Fatalf("Failed to create a ClientConfig: %v. Exiting.", err)
	}

	client, err := clientset.NewForConfig(clientConfig)
	if err != nil {
		klog.Fatalf("Failed to create a ClientSet: %v. Exiting.", err)
	}

	f, c := kubemark.GetHollowKubeletConfig(config.createHollowKubeletOptions())

	heartbeatClientConfig := *clientConfig
	heartbeatClientConfig.Timeout = c.NodeStatusUpdateFrequency.Duration
	// The timeout is the minimum of the lease duration and status update frequency
	leaseTimeout := time.Duration(c.NodeLeaseDurationSeconds) * time.Second
	if heartbeatClientConfig.Timeout > leaseTimeout {
		heartbeatClientConfig.Timeout = leaseTimeout
	}

	heartbeatClientConfig.QPS = float32(-1)
	heartbeatClient, err := clientset.NewForConfig(&heartbeatClientConfig)
	if err != nil {
		klog.Fatalf("Failed to create a ClientSet: %v. Exiting.", err)
	}

	cadvisorInterface := &cadvisortest.Fake{
		NodeName: config.NodeName,
	}
	containerManager := cm.NewStubContainerManager()

	endpoint, err := fakeremote.GenerateEndpoint()
	if err != nil {
		klog.Fatalf("Failed to generate fake endpoint %v.", err)
	}
	fakeRemoteRuntime := fakeremote.NewFakeRemoteRuntime()
	if err = fakeRemoteRuntime.Start(endpoint); err != nil {
		klog.Fatalf("Failed to start fake runtime %v.", err)
	}
	defer fakeRemoteRuntime.Stop()
	runtimeService, err := remote.NewRemoteRuntimeService(endpoint, 15*time.Second)
	if err != nil {
		klog.Fatalf("Failed to init runtime service %v.", err)
	}

	remoteImageService, err := remote.NewRemoteImageService(f.RemoteImageEndpoint, 15*time.Second)
	if err != nil {
		klog.Fatalf("Failed to init image service %v.", err)
	}

	hollowKubelet := kubemark.NewHollowKubelet(
		f, c,
		client,
		heartbeatClient,
		cadvisorInterface,
		remoteImageService,
		runtimeService,
		containerManager,
	)
	hollowKubelet.Run()
}
