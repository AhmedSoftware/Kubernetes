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

package addons

import (
	"fmt"
	"net"

	"runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kuberuntime "k8s.io/apimachinery/pkg/runtime"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	"k8s.io/kubernetes/cmd/kubeadm/app/images"
	kubeadmutil "k8s.io/kubernetes/cmd/kubeadm/app/util"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/v1"
	extensions "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	"k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

// CreateEssentialAddons creates the kube-proxy and kube-dns addons
func CreateEssentialAddons(cfg *kubeadmapi.MasterConfiguration, client *clientset.Clientset) error {
	proxyConfigMapBytes, err := kubeadmutil.ParseTemplate(KubeProxyConfigMap, struct{ MasterEndpoint string }{
		// Fetch this value from the kubeconfig file
		MasterEndpoint: fmt.Sprintf("https://%s:%d", cfg.API.AdvertiseAddresses[0], cfg.API.Port),
	})
	if err != nil {
		return fmt.Errorf("error when parsing kube-proxy configmap template: %v", err)
	}

	proxyDaemonSetBytes, err := kubeadmutil.ParseTemplate(KubeProxyDaemonSet, struct{ Image, ClusterCIDR string }{
		Image:       images.GetCoreImage("proxy", cfg, kubeadmapi.GlobalEnvParams.HyperkubeImage),
		ClusterCIDR: getClusterCIDR(cfg.Networking.PodSubnet),
	})
	if err != nil {
		return fmt.Errorf("error when parsing kube-proxy daemonset template: %v", err)
	}

	dnsDeploymentBytes, err := kubeadmutil.ParseTemplate(KubeDNSDeployment, struct {
		ImageRepository, Arch, Version, DNSDomain string
		Replicas                                  int
	}{
		ImageRepository: kubeadmapi.GlobalEnvParams.RepositoryPrefix,
		Arch:            runtime.GOARCH,
		// TODO: Support larger amount of replicas?
		Replicas:  1,
		Version:   KubeDNSVersion,
		DNSDomain: cfg.Networking.DNSDomain,
	})
	if err != nil {
		return fmt.Errorf("error when parsing kube-dns deployment template: %v", err)
	}

	dnsip, err := getDNSIP(client)
	if err != nil {
		return err
	}

	dnsServiceBytes, err := kubeadmutil.ParseTemplate(KubeDNSService, struct{ DNSIP string }{
		DNSIP: dnsip,
	})
	if err != nil {
		return fmt.Errorf("error when parsing kube-proxy configmap template: %v", err)
	}

	err = CreateKubeProxyAddon(proxyConfigMapBytes, proxyDaemonSetBytes, client)
	if err != nil {
		return err
	}
	fmt.Println("[addons] Created essential addon: kube-proxy")

	err = CreateKubeDNSAddon(dnsDeploymentBytes, dnsServiceBytes, client)
	if err != nil {
		return err
	}
	fmt.Println("[addons] Created essential addon: kube-dns")
	return nil
}

func CreateKubeProxyAddon(configMapBytes, daemonSetbytes []byte, client *clientset.Clientset) error {
	kubeproxyConfigMap := &v1.ConfigMap{}
	if err := kuberuntime.DecodeInto(api.Codecs.UniversalDecoder(), configMapBytes, kubeproxyConfigMap); err != nil {
		return fmt.Errorf("unable to decode kube-proxy configmap %v", err)
	}

	if _, err := client.CoreV1().ConfigMaps(metav1.NamespaceSystem).Create(kubeproxyConfigMap); err != nil {
		return fmt.Errorf("unable to create a new kube-proxy configmap: %v", err)
	}

	kubeproxyDaemonSet := &extensions.DaemonSet{}
	if err := kuberuntime.DecodeInto(api.Codecs.UniversalDecoder(), daemonSetbytes, kubeproxyDaemonSet); err != nil {
		return fmt.Errorf("unable to decode kube-proxy daemonset %v", err)
	}

	if _, err := client.ExtensionsV1beta1().DaemonSets(metav1.NamespaceSystem).Create(kubeproxyDaemonSet); err != nil {
		return fmt.Errorf("unable to create a new kube-proxy daemonset: %v", err)
	}
	return nil
}

func CreateKubeDNSAddon(deploymentBytes, serviceBytes []byte, client *clientset.Clientset) error {
	kubednsDeployment := &extensions.Deployment{}
	if err := kuberuntime.DecodeInto(api.Codecs.UniversalDecoder(), deploymentBytes, kubednsDeployment); err != nil {
		return fmt.Errorf("unable to decode kube-dns deployment %v", err)
	}

	// TODO: All these .Create(foo) calls should instead be more like "kubectl apply -f" commands; they should not fail if there are existing objects with the same name
	if _, err := client.ExtensionsV1beta1().Deployments(metav1.NamespaceSystem).Create(kubednsDeployment); err != nil {
		return fmt.Errorf("unable to create a new kube-dns deployment: %v", err)
	}

	kubednsService := &v1.Service{}
	if err := kuberuntime.DecodeInto(api.Codecs.UniversalDecoder(), serviceBytes, kubednsService); err != nil {
		return fmt.Errorf("unable to decode kube-dns service %v", err)
	}

	if _, err := client.CoreV1().Services(metav1.NamespaceSystem).Create(kubednsService); err != nil {
		return fmt.Errorf("unable to create a new kube-dns service: %v", err)
	}
	return nil
}

// getDNSIP fetches the kubernetes service's ClusterIP and appends a "0" to it in order to get the DNS IP
func getDNSIP(client *clientset.Clientset) (string, error) {
	k8ssvc, err := client.CoreV1().Services(metav1.NamespaceDefault).Get("kubernetes", metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("couldn't fetch information about the kubernetes service: %v", err)
	}

	if len(k8ssvc.Spec.ClusterIP) == 0 {
		return "", fmt.Errorf("couldn't fetch a valid clusterIP from the kubernetes service")
	}

	dnsIP := fmt.Sprintf("%s0", k8ssvc.Spec.ClusterIP)

	// Check that it's a valid IP
	realIP := net.ParseIP(dnsIP)
	if realIP == nil {
		return "", fmt.Errorf("could not parse dns ip %q: %v", dnsIP, err)
	}
	return dnsIP, nil
}

func getClusterCIDR(podsubnet string) string {
	if len(podsubnet) == 0 {
		return ""
	}
	return "--cluster-cidr" + podsubnet
}
