/*
Copyright 2020 The Kubernetes Authors.

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

package netpol

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
	e2epod "k8s.io/kubernetes/test/e2e/framework/pod"
	netutils "k8s.io/utils/net"
)

// TestPod represents an actual running pod. For each Pod defined by the model,
// there will be a corresponding TestPod. TestPod includes some runtime info
// (namespace name, service IP) which is not available in the model.
type TestPod struct {
	Namespace     string
	Name          string
	ContainerName string
	ServiceIP     string
}

func (pod TestPod) PodString() PodString {
	return NewPodString(pod.Namespace, pod.Name)
}

// kubeManager provides a convenience interface to kube functionality that we leverage for polling NetworkPolicy connections.
// Its responsibilities are:
//   - creating resources (pods, deployments, namespaces, services, network policies)
//   - modifying and cleaning up resources
type kubeManager struct {
	framework      *framework.Framework
	clientSet      clientset.Interface
	namespaceNames []string
	allPods        []TestPod
	allPodStrings  []PodString
	dnsDomain      string
}

// newKubeManager is a utility function that wraps creation of the kubeManager instance.
func newKubeManager(framework *framework.Framework, dnsDomain string) *kubeManager {
	return &kubeManager{
		framework: framework,
		clientSet: framework.ClientSet,
		dnsDomain: dnsDomain,
	}
}

// initializeCluster initialized the cluster, creating namespaces pods and services as needed.
func (k *kubeManager) initializeClusterFromModel(ctx context.Context, model *Model) error {
	var createdPods []*v1.Pod
	for _, ns := range model.Namespaces {
		// no labels needed, we just need the default kubernetes.io/metadata.name label
		namespace, err := k.framework.CreateNamespace(ctx, ns.BaseName, nil)
		if err != nil {
			return err
		}
		namespaceName := namespace.Name
		k.namespaceNames = append(k.namespaceNames, namespaceName)

		for _, pod := range ns.Pods {
			framework.Logf("creating pod %s/%s with matching service", namespaceName, pod.Name)

			// note that we defer the logic of pod (i.e. node selector) specifics to the model
			// which is aware of linux vs windows pods
			kubePod, err := k.createPod(ctx, pod.KubePod(namespaceName))
			if err != nil {
				return err
			}

			createdPods = append(createdPods, kubePod)
			svc, err := k.createService(ctx, pod.Service(namespaceName))
			if err != nil {
				return err
			}
			if netutils.ParseIPSloppy(svc.Spec.ClusterIP) == nil {
				return fmt.Errorf("empty IP address found for service %s/%s", svc.Namespace, svc.Name)
			}

			k.allPods = append(k.allPods, TestPod{
				Namespace:     kubePod.Namespace,
				Name:          kubePod.Name,
				ContainerName: pod.Containers[0].Name(),
				ServiceIP:     svc.Spec.ClusterIP,
			})
			k.allPodStrings = append(k.allPodStrings, NewPodString(kubePod.Namespace, kubePod.Name))
		}
	}

	for _, createdPod := range createdPods {
		err := e2epod.WaitTimeoutForPodReadyInNamespace(ctx, k.clientSet, createdPod.Name, createdPod.Namespace, framework.PodStartTimeout)
		if err != nil {
			return fmt.Errorf("unable to wait for pod %s/%s: %w", createdPod.Namespace, createdPod.Name, err)
		}
	}

	return nil
}

func (k *kubeManager) AllPods() []TestPod {
	return k.allPods
}

func (k *kubeManager) AllPodStrings() []PodString {
	return k.allPodStrings
}

func (k *kubeManager) DNSDomain() string {
	return k.dnsDomain
}

func (k *kubeManager) NamespaceNames() []string {
	return k.namespaceNames
}

// getPod gets a pod by namespace and name.
func (k *kubeManager) getPod(ctx context.Context, ns string, name string) (*v1.Pod, error) {
	kubePod, err := k.clientSet.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get pod %s/%s: %w", ns, name, err)
	}
	return kubePod, nil
}

// probeConnectivity execs into a pod and checks its connectivity to another pod returning the true if the connectivity state is the expected.
// Implements the Prober interface.
func (k *kubeManager) probeConnectivity(args *probeConnectivityArgs) (bool, error) {
	port := strconv.Itoa(args.toPort)
	if args.addrTo == "" {
		return false, fmt.Errorf("empty addrTo field")
	}
	negative := ""
	if !args.expectConnectivity {
		negative = "!"
	}

	now := time.Now()
	framework.Logf("Starting probe from pod %v to %v at %v", args.podFrom, args.addrTo, now)
	cmd := []string{"timeout", "30", "/bin/sh", "-c", fmt.Sprintf(`until %s /agnhost connect --timeout=3s --protocol=%s %s ; do echo "Date: $(date)" ; sleep 1; done`,
		negative, strings.ToLower(string(args.protocol)), net.JoinHostPort(args.addrTo, port))}

	// NOTE: The return value of this function[probeConnectivity] should be true if the probe is successful and false otherwise.
	_, stderr, probeError := k.executeRemoteCommand(args.nsFrom, args.podFrom, args.containerFrom, cmd)
	framework.Logf("probe %s/%s -> %s took %v expected %v :: stderr - %s probeError - %v", args.nsFrom, args.podFrom, args.addrTo, time.Since(now), args.expectConnectivity, stderr, probeError)
	return probeError == nil, probeError
}

// executeRemoteCommand executes a remote shell command on the given pod.
func (k *kubeManager) executeRemoteCommand(namespace string, pod string, containerName string, command []string) (string, string, error) {
	return e2epod.ExecWithOptions(k.framework, e2epod.ExecOptions{
		Command:            command,
		Namespace:          namespace,
		PodName:            pod,
		ContainerName:      containerName,
		Stdin:              nil,
		CaptureStdout:      true,
		CaptureStderr:      true,
		PreserveWhitespace: false,
	})
}

// createService is a convenience function for service setup.
func (k *kubeManager) createService(ctx context.Context, service *v1.Service) (*v1.Service, error) {
	ns := service.Namespace
	name := service.Name

	createdService, err := k.clientSet.CoreV1().Services(ns).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to create service %s/%s: %w", ns, name, err)
	}
	return createdService, nil
}

// createPod is a convenience function for pod setup.
func (k *kubeManager) createPod(ctx context.Context, pod *v1.Pod) (*v1.Pod, error) {
	ns := pod.Namespace
	framework.Logf("creating pod %s/%s", ns, pod.Name)

	createdPod, err := k.clientSet.CoreV1().Pods(ns).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to create pod %s/%s: %w", ns, pod.Name, err)
	}
	return createdPod, nil
}

// cleanNetworkPolicies is a convenience function for deleting network policies before startup of any new test.
func (k *kubeManager) cleanNetworkPolicies(ctx context.Context) error {
	for _, ns := range k.namespaceNames {
		framework.Logf("deleting policies in %s ..........", ns)
		l, err := k.clientSet.NetworkingV1().NetworkPolicies(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("unable to list network policies in ns %s: %w", ns, err)
		}
		for _, np := range l.Items {
			framework.Logf("deleting network policy %s/%s", ns, np.Name)
			err = k.clientSet.NetworkingV1().NetworkPolicies(ns).Delete(ctx, np.Name, metav1.DeleteOptions{})
			if err != nil {
				return fmt.Errorf("unable to delete network policy %s/%s: %w", ns, np.Name, err)
			}
		}
	}
	return nil
}

// createNetworkPolicy is a convenience function for creating network policies.
func (k *kubeManager) createNetworkPolicy(ctx context.Context, ns string, netpol *networkingv1.NetworkPolicy) (*networkingv1.NetworkPolicy, error) {
	framework.Logf("creating network policy %s/%s", ns, netpol.Name)
	netpol.ObjectMeta.Namespace = ns
	np, err := k.clientSet.NetworkingV1().NetworkPolicies(ns).Create(ctx, netpol, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to create network policy %s/%s: %w", ns, netpol.Name, err)
	}
	return np, nil
}

// updateNetworkPolicy is a convenience function for updating network policies.
func (k *kubeManager) updateNetworkPolicy(ctx context.Context, ns string, netpol *networkingv1.NetworkPolicy) (*networkingv1.NetworkPolicy, error) {
	framework.Logf("updating network policy %s/%s", ns, netpol.Name)
	netpol.ObjectMeta.Namespace = ns
	np, err := k.clientSet.NetworkingV1().NetworkPolicies(ns).Update(ctx, netpol, metav1.UpdateOptions{})
	if err != nil {
		return np, fmt.Errorf("unable to update network policy %s/%s: %w", ns, netpol.Name, err)
	}
	return np, nil
}

// getNamespace gets a namespace object from kubernetes.
func (k *kubeManager) getNamespace(ctx context.Context, ns string) (*v1.Namespace, error) {
	selectedNameSpace, err := k.clientSet.CoreV1().Namespaces().Get(ctx, ns, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get namespace %s: %w", ns, err)
	}
	return selectedNameSpace, nil
}
