//go:build !providerless
// +build !providerless

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

package gce

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetLoadBalancer(t *testing.T) {
	t.Parallel()

	vals := DefaultTestClusterValues()
	gce, err := fakeGCECloud(vals)
	if err != nil {
		t.Fatal(err)
	}

	apiService := fakeLoadbalancerService("")

	// When a load balancer has not been created
	status, found, err := gce.GetLoadBalancer(context.Background(), vals.ClusterName, apiService)
	if err != nil {
		t.Fatal(err)
	}
	if status != nil {
		t.Errorf("expected=%v, got=%v", nil, status)
	}
	if found != false {
		t.Errorf("expected=%t, got=%t", false, found)
	}

	nodeNames := []string{"test-node-1"}
	nodes, err := createAndInsertNodes(gce, nodeNames, vals.ZoneName)
	if err != nil {
		t.Fatal(err)
	}
	expectedStatus, err := gce.EnsureLoadBalancer(context.Background(), vals.ClusterName, apiService, nodes)
	if err != nil {
		t.Fatal(err)
	}

	status, found, err = gce.GetLoadBalancer(context.Background(), vals.ClusterName, apiService)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(expectedStatus, status) {
		t.Errorf("expected=%v, got=%v", expectedStatus, status)
	}
	if found != true {
		t.Errorf("expected=%t, got=%t", true, found)
	}
}

func TestEnsureLoadBalancerCreatesExternalLb(t *testing.T) {
	t.Parallel()

	vals := DefaultTestClusterValues()
	gce, err := fakeGCECloud(vals)
	if err != nil {
		t.Fatal(err)
	}

	nodeNames := []string{"test-node-1"}
	nodes, err := createAndInsertNodes(gce, nodeNames, vals.ZoneName)
	if err != nil {
		t.Fatal(err)
	}

	apiService := fakeLoadbalancerService("")
	status, err := gce.EnsureLoadBalancer(context.Background(), vals.ClusterName, apiService, nodes)
	if err != nil {
		t.Fatal(err)
	}
	if status.Ingress == nil {
		t.Error("status.Ingress is nil")
	}
	assertExternalLbResources(t, gce, apiService, vals, nodeNames)
}

func TestEnsureLoadBalancerCreatesInternalLb(t *testing.T) {
	t.Parallel()

	vals := DefaultTestClusterValues()
	gce, err := fakeGCECloud(vals)
	if err != nil {
		t.Fatal(err)
	}

	nodeNames := []string{"test-node-1"}
	nodes, err := createAndInsertNodes(gce, nodeNames, vals.ZoneName)
	if err != nil {
		t.Fatal(err)
	}

	apiService := fakeLoadbalancerService(string(LBTypeInternal))
	apiService, err = gce.client.CoreV1().Services(apiService.Namespace).Create(context.TODO(), apiService, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	status, err := gce.EnsureLoadBalancer(context.Background(), vals.ClusterName, apiService, nodes)
	if err != nil {
		t.Fatal(err)
	}
	if status.Ingress == nil {
		t.Error("status.Ingress is nil")
	}
	assertInternalLbResources(t, gce, apiService, vals, nodeNames)
}

func TestEnsureLoadBalancerDeletesExistingInternalLb(t *testing.T) {
	t.Parallel()

	vals := DefaultTestClusterValues()
	gce, err := fakeGCECloud(vals)
	if err != nil {
		t.Fatal(err)
	}

	nodeNames := []string{"test-node-1"}
	nodes, err := createAndInsertNodes(gce, nodeNames, vals.ZoneName)
	if err != nil {
		t.Fatal(err)
	}

	apiService := fakeLoadbalancerService("")
	createInternalLoadBalancer(gce, apiService, nil, nodeNames, vals.ClusterName, vals.ClusterID, vals.ZoneName)

	status, err := gce.EnsureLoadBalancer(context.Background(), vals.ClusterName, apiService, nodes)
	if err != nil {
		t.Fatal(err)
	}
	if status.Ingress == nil {
		t.Error("status.Ingress is nil")
	}

	assertExternalLbResources(t, gce, apiService, vals, nodeNames)
	assertInternalLbResourcesDeleted(t, gce, apiService, vals, false)
}

func TestEnsureLoadBalancerDeletesExistingExternalLb(t *testing.T) {
	t.Parallel()

	vals := DefaultTestClusterValues()
	gce, err := fakeGCECloud(vals)
	if err != nil {
		t.Fatal(err)
	}

	nodeNames := []string{"test-node-1"}
	nodes, err := createAndInsertNodes(gce, nodeNames, vals.ZoneName)
	if err != nil {
		t.Fatal(err)
	}

	apiService := fakeLoadbalancerService("")
	createExternalLoadBalancer(gce, apiService, nodeNames, vals.ClusterName, vals.ClusterID, vals.ZoneName)

	apiService = fakeLoadbalancerService(string(LBTypeInternal))
	apiService, err = gce.client.CoreV1().Services(apiService.Namespace).Create(context.TODO(), apiService, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	status, err := gce.EnsureLoadBalancer(context.Background(), vals.ClusterName, apiService, nodes)
	if err != nil {
		t.Fatal(err)
	}
	if status.Ingress == nil {
		t.Error("status.Ingress is nil")
	}

	assertInternalLbResources(t, gce, apiService, vals, nodeNames)
	assertExternalLbResourcesDeleted(t, gce, apiService, vals, false)
}

func TestEnsureLoadBalancerDeletedDeletesExternalLb(t *testing.T) {
	t.Parallel()

	vals := DefaultTestClusterValues()
	gce, err := fakeGCECloud(vals)
	if err != nil {
		t.Fatal(err)
	}

	nodeNames := []string{"test-node-1"}
	_, err = createAndInsertNodes(gce, nodeNames, vals.ZoneName)
	if err != nil {
		t.Fatal(err)
	}

	apiService := fakeLoadbalancerService("")
	createExternalLoadBalancer(gce, apiService, nodeNames, vals.ClusterName, vals.ClusterID, vals.ZoneName)

	err = gce.EnsureLoadBalancerDeleted(context.Background(), vals.ClusterName, apiService)
	if err != nil {
		t.Fatal(err)
	}
	assertExternalLbResourcesDeleted(t, gce, apiService, vals, true)
}

func TestEnsureLoadBalancerDeletedDeletesInternalLb(t *testing.T) {
	t.Parallel()

	vals := DefaultTestClusterValues()
	gce, err := fakeGCECloud(vals)
	if err != nil {
		t.Fatal(err)
	}

	nodeNames := []string{"test-node-1"}
	_, err = createAndInsertNodes(gce, nodeNames, vals.ZoneName)
	if err != nil {
		t.Fatal(err)
	}

	apiService := fakeLoadbalancerService(string(LBTypeInternal))
	apiService, err = gce.client.CoreV1().Services(apiService.Namespace).Create(context.TODO(), apiService, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	createInternalLoadBalancer(gce, apiService, nil, nodeNames, vals.ClusterName, vals.ClusterID, vals.ZoneName)

	err = gce.EnsureLoadBalancerDeleted(context.Background(), vals.ClusterName, apiService)
	if err != nil {
		t.Fatal(err)
	}
	assertInternalLbResourcesDeleted(t, gce, apiService, vals, true)
}

func TestProjectsBasePath(t *testing.T) {
	t.Parallel()
	vals := DefaultTestClusterValues()
	gce, err := fakeGCECloud(vals)
	// Loadbalancer controller code expects basepath to contain the projects string.
	expectProjectsBasePath := "https://compute.googleapis.com/compute/v1/projects/"
	// See https://github.com/kubernetes/kubernetes/issues/102757, the endpoint can have mtls in some cases.
	expectMtlsProjectsBasePath := "https://compute.mtls.googleapis.com/compute/v1/projects/"
	if err != nil {
		t.Fatal(err)
	}
	if gce.projectsBasePath != expectProjectsBasePath && gce.projectsBasePath != expectMtlsProjectsBasePath {
		t.Errorf("Compute projectsBasePath has changed. Got %q, want %q or %q", gce.projectsBasePath, expectProjectsBasePath, expectMtlsProjectsBasePath)
	}
}
