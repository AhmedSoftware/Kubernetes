/*
Copyright 2014 The Kubernetes Authors.

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

package kubernetesservice

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/intstr"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	v1informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes/fake"
	v1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	core "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
	netutils "k8s.io/utils/net"
)

func TestCreateOrUpdateMasterService(t *testing.T) {
	singleStack := corev1.IPFamilyPolicySingleStack
	ns := metav1.NamespaceDefault
	om := func(name string) metav1.ObjectMeta {
		return metav1.ObjectMeta{Namespace: ns, Name: name}
	}

	createTests := []struct {
		testName     string
		serviceName  string
		servicePorts []corev1.ServicePort
		serviceType  corev1.ServiceType
		expectCreate *corev1.Service // nil means none expected
	}{
		{
			testName:    "service does not exist",
			serviceName: "foo",
			servicePorts: []corev1.ServicePort{
				{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
			},
			serviceType: corev1.ServiceTypeClusterIP,
			expectCreate: &corev1.Service{
				ObjectMeta: om("foo"),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
					},
					Selector:        nil,
					ClusterIP:       "1.2.3.4",
					IPFamilyPolicy:  &singleStack,
					SessionAffinity: corev1.ServiceAffinityNone,
					Type:            corev1.ServiceTypeClusterIP,
				},
			},
		},
	}
	for _, test := range createTests {
		t.Run(test.testName, func(t *testing.T) {
			master := Controller{}
			fakeClient := fake.NewSimpleClientset()
			serviceStore := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
			master.serviceLister = v1listers.NewServiceLister(serviceStore)
			master.client = fakeClient
			master.CreateOrUpdateMasterServiceIfNeeded(test.serviceName, netutils.ParseIPSloppy("1.2.3.4"), test.servicePorts, test.serviceType, false)
			creates := []core.CreateAction{}
			for _, action := range fakeClient.Actions() {
				if action.GetVerb() == "create" {
					creates = append(creates, action.(core.CreateAction))
				}
			}
			if test.expectCreate != nil {
				if len(creates) != 1 {
					t.Errorf("case %q: unexpected creations: %v", test.testName, creates)
				} else {
					obj := creates[0].GetObject()
					if e, a := test.expectCreate.Spec, obj.(*corev1.Service).Spec; !reflect.DeepEqual(e, a) {
						t.Errorf("case %q: expected create:\n%#v\ngot:\n%#v\n", test.testName, e, a)
					}
				}
			}
			if test.expectCreate == nil && len(creates) > 1 {
				t.Errorf("case %q: no create expected, yet saw: %v", test.testName, creates)
			}
		})
	}

	reconcileTests := []struct {
		testName     string
		serviceName  string
		servicePorts []corev1.ServicePort
		serviceType  corev1.ServiceType
		service      *corev1.Service
		expectUpdate *corev1.Service // nil means none expected
	}{
		{
			testName:    "service definition wrong port",
			serviceName: "foo",
			servicePorts: []corev1.ServicePort{
				{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
			},
			serviceType: corev1.ServiceTypeClusterIP,
			service: &corev1.Service{
				ObjectMeta: om("foo"),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Name: "foo", Port: 8000, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
					},
					Selector:        nil,
					ClusterIP:       "1.2.3.4",
					SessionAffinity: corev1.ServiceAffinityNone,
					Type:            corev1.ServiceTypeClusterIP,
				},
			},
			expectUpdate: &corev1.Service{
				ObjectMeta: om("foo"),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
					},
					Selector:        nil,
					ClusterIP:       "1.2.3.4",
					SessionAffinity: corev1.ServiceAffinityNone,
					Type:            corev1.ServiceTypeClusterIP,
				},
			},
		},
		{
			testName:    "service definition missing port",
			serviceName: "foo",
			servicePorts: []corev1.ServicePort{
				{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
				{Name: "baz", Port: 1000, Protocol: "TCP", TargetPort: intstr.FromInt32(1000)},
			},
			serviceType: corev1.ServiceTypeClusterIP,
			service: &corev1.Service{
				ObjectMeta: om("foo"),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
					},
					Selector:        nil,
					ClusterIP:       "1.2.3.4",
					SessionAffinity: corev1.ServiceAffinityNone,
					Type:            corev1.ServiceTypeClusterIP,
				},
			},
			expectUpdate: &corev1.Service{
				ObjectMeta: om("foo"),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
						{Name: "baz", Port: 1000, Protocol: "TCP", TargetPort: intstr.FromInt32(1000)},
					},
					Selector:        nil,
					ClusterIP:       "1.2.3.4",
					SessionAffinity: corev1.ServiceAffinityNone,
					Type:            corev1.ServiceTypeClusterIP,
				},
			},
		},
		{
			testName:    "service definition incorrect port",
			serviceName: "foo",
			servicePorts: []corev1.ServicePort{
				{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
			},
			serviceType: corev1.ServiceTypeClusterIP,
			service: &corev1.Service{
				ObjectMeta: om("foo"),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Name: "bar", Port: 1000, Protocol: "UDP", TargetPort: intstr.FromInt32(1000)},
					},
					Selector:        nil,
					ClusterIP:       "1.2.3.4",
					SessionAffinity: corev1.ServiceAffinityNone,
					Type:            corev1.ServiceTypeClusterIP,
				},
			},
			expectUpdate: &corev1.Service{
				ObjectMeta: om("foo"),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
					},
					Selector:        nil,
					ClusterIP:       "1.2.3.4",
					SessionAffinity: corev1.ServiceAffinityNone,
					Type:            corev1.ServiceTypeClusterIP,
				},
			},
		},
		{
			testName:    "service definition incorrect port name",
			serviceName: "foo",
			servicePorts: []corev1.ServicePort{
				{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
			},
			serviceType: corev1.ServiceTypeClusterIP,
			service: &corev1.Service{
				ObjectMeta: om("foo"),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Name: "foo", Port: 1000, Protocol: "UDP", TargetPort: intstr.FromInt32(1000)},
					},
					Selector:        nil,
					ClusterIP:       "1.2.3.4",
					SessionAffinity: corev1.ServiceAffinityNone,
					Type:            corev1.ServiceTypeClusterIP,
				},
			},
			expectUpdate: &corev1.Service{
				ObjectMeta: om("foo"),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
					},
					Selector:        nil,
					ClusterIP:       "1.2.3.4",
					SessionAffinity: corev1.ServiceAffinityNone,
					Type:            corev1.ServiceTypeClusterIP,
				},
			},
		},
		{
			testName:    "service definition incorrect target port",
			serviceName: "foo",
			servicePorts: []corev1.ServicePort{
				{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
			},
			serviceType: corev1.ServiceTypeClusterIP,
			service: &corev1.Service{
				ObjectMeta: om("foo"),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(1000)},
					},
					Selector:        nil,
					ClusterIP:       "1.2.3.4",
					SessionAffinity: corev1.ServiceAffinityNone,
					Type:            corev1.ServiceTypeClusterIP,
				},
			},
			expectUpdate: &corev1.Service{
				ObjectMeta: om("foo"),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
					},
					Selector:        nil,
					ClusterIP:       "1.2.3.4",
					SessionAffinity: corev1.ServiceAffinityNone,
					Type:            corev1.ServiceTypeClusterIP,
				},
			},
		},
		{
			testName:    "service definition incorrect protocol",
			serviceName: "foo",
			servicePorts: []corev1.ServicePort{
				{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
			},
			serviceType: corev1.ServiceTypeClusterIP,
			service: &corev1.Service{
				ObjectMeta: om("foo"),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Name: "foo", Port: 8080, Protocol: "UDP", TargetPort: intstr.FromInt32(8080)},
					},
					Selector:        nil,
					ClusterIP:       "1.2.3.4",
					SessionAffinity: corev1.ServiceAffinityNone,
					Type:            corev1.ServiceTypeClusterIP,
				},
			},
			expectUpdate: &corev1.Service{
				ObjectMeta: om("foo"),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
					},
					Selector:        nil,
					ClusterIP:       "1.2.3.4",
					SessionAffinity: corev1.ServiceAffinityNone,
					Type:            corev1.ServiceTypeClusterIP,
				},
			},
		},
		{
			testName:    "service definition has incorrect type",
			serviceName: "foo",
			servicePorts: []corev1.ServicePort{
				{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
			},
			serviceType: corev1.ServiceTypeClusterIP,
			service: &corev1.Service{
				ObjectMeta: om("foo"),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
					},
					Selector:        nil,
					ClusterIP:       "1.2.3.4",
					SessionAffinity: corev1.ServiceAffinityNone,
					Type:            corev1.ServiceTypeNodePort,
				},
			},
			expectUpdate: &corev1.Service{
				ObjectMeta: om("foo"),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
					},
					Selector:        nil,
					ClusterIP:       "1.2.3.4",
					SessionAffinity: corev1.ServiceAffinityNone,
					Type:            corev1.ServiceTypeClusterIP,
				},
			},
		},
		{
			testName:    "service definition satisfies",
			serviceName: "foo",
			servicePorts: []corev1.ServicePort{
				{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
			},
			serviceType: corev1.ServiceTypeClusterIP,
			service: &corev1.Service{
				ObjectMeta: om("foo"),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
					},
					Selector:        nil,
					ClusterIP:       "1.2.3.4",
					SessionAffinity: corev1.ServiceAffinityNone,
					Type:            corev1.ServiceTypeClusterIP,
				},
			},
			expectUpdate: nil,
		},
	}
	for _, test := range reconcileTests {
		t.Run(test.testName, func(t *testing.T) {
			master := Controller{}
			fakeClient := fake.NewSimpleClientset(test.service)
			serviceInformer := v1informers.NewServiceInformer(fakeClient, metav1.NamespaceDefault, 12*time.Hour, cache.Indexers{})
			serviceStore := serviceInformer.GetIndexer()
			err := serviceStore.Add(test.service)
			if err != nil {
				t.Fatalf("unexpected error adding service %v to the store: %v", test.service, err)
			}
			master.serviceLister = v1listers.NewServiceLister(serviceStore)
			master.client = fakeClient
			err = master.CreateOrUpdateMasterServiceIfNeeded(test.serviceName, netutils.ParseIPSloppy("1.2.3.4"), test.servicePorts, test.serviceType, true)
			if err != nil {
				t.Errorf("case %q: unexpected error: %v", test.testName, err)
			}
			updates := []core.UpdateAction{}
			for _, action := range fakeClient.Actions() {
				if action.GetVerb() == "update" {
					updates = append(updates, action.(core.UpdateAction))
				}
			}
			if test.expectUpdate != nil {
				if len(updates) != 1 {
					t.Errorf("case %q: unexpected updates: %v", test.testName, updates)
				} else {
					obj := updates[0].GetObject()
					if e, a := test.expectUpdate.Spec, obj.(*corev1.Service).Spec; !reflect.DeepEqual(e, a) {
						t.Errorf("case %q: expected update:\n%#v\ngot:\n%#v\n", test.testName, e, a)
					}
				}
			}
			if test.expectUpdate == nil && len(updates) > 0 {
				t.Errorf("case %q: no update expected, yet saw: %v", test.testName, updates)
			}
		})
	}

	nonReconcileTests := []struct {
		testName     string
		serviceName  string
		servicePorts []corev1.ServicePort
		serviceType  corev1.ServiceType
		service      *corev1.Service
		expectUpdate *corev1.Service // nil means none expected
	}{
		{
			testName:    "service definition wrong port, no expected update",
			serviceName: "foo",
			servicePorts: []corev1.ServicePort{
				{Name: "foo", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromInt32(8080)},
			},
			serviceType: corev1.ServiceTypeClusterIP,
			service: &corev1.Service{
				ObjectMeta: om("foo"),
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Name: "foo", Port: 1000, Protocol: "TCP", TargetPort: intstr.FromInt32(1000)},
					},
					Selector:        nil,
					ClusterIP:       "1.2.3.4",
					SessionAffinity: corev1.ServiceAffinityNone,
					Type:            corev1.ServiceTypeClusterIP,
				},
			},
			expectUpdate: nil,
		},
	}
	for _, test := range nonReconcileTests {
		t.Run(test.testName, func(t *testing.T) {
			master := Controller{}
			fakeClient := fake.NewSimpleClientset(test.service)
			master.client = fakeClient
			serviceInformer := v1informers.NewServiceInformer(fakeClient, metav1.NamespaceDefault, 12*time.Hour, cache.Indexers{})
			serviceStore := serviceInformer.GetIndexer()
			err := serviceStore.Add(test.service)
			if err != nil {
				t.Fatalf("unexpected error adding service %v to the store: %v", test.service, err)
			}
			master.serviceLister = v1listers.NewServiceLister(serviceStore)

			err = master.CreateOrUpdateMasterServiceIfNeeded(test.serviceName, netutils.ParseIPSloppy("1.2.3.4"), test.servicePorts, test.serviceType, false)
			if err != nil {
				t.Errorf("case %q: unexpected error: %v", test.testName, err)
			}
			updates := []core.UpdateAction{}
			for _, action := range fakeClient.Actions() {
				if action.GetVerb() == "update" {
					updates = append(updates, action.(core.UpdateAction))
				}
			}
			if test.expectUpdate != nil {
				if len(updates) != 1 {
					t.Errorf("case %q: unexpected updates: %v", test.testName, updates)
				} else {
					obj := updates[0].GetObject()
					if e, a := test.expectUpdate.Spec, obj.(*corev1.Service).Spec; !reflect.DeepEqual(e, a) {
						t.Errorf("case %q: expected update:\n%#v\ngot:\n%#v\n", test.testName, e, a)
					}
				}
			}
			if test.expectUpdate == nil && len(updates) > 0 {
				t.Errorf("case %q: no update expected, yet saw: %v", test.testName, updates)
			}
		})
	}
}

func TestUpdateKubernetesService(t *testing.T) {
	var healthy atomic.Bool
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if healthy.Load() {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "OK")
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, "NOOK")
		}
	}))
	ts.EnableHTTP2 = true
	ts.StartTLS()
	defer ts.Close()
	transport, ok := ts.Client().Transport.(*http.Transport)
	if !ok {
		t.Fatalf("failed to assert *http.Transport")
	}
	restConfig := &rest.Config{
		Host:      ts.URL,
		Transport: utilnet.SetTransportDefaults(transport),
		// These fields are required to create a REST client.
		ContentConfig: rest.ContentConfig{
			GroupVersion:         &schema.GroupVersion{},
			NegotiatedSerializer: &serializer.CodecFactory{},
		},
	}
	restClient, err := rest.RESTClientFor(restConfig)
	if err != nil {
		t.Fatalf("failed to create REST client: %v", err)
	}

	fakeClient := fake.NewSimpleClientset()
	serviceStore := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	err = serviceStore.Add(&corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubernetesServiceName,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "foo", Port: 1000, Protocol: "TCP", TargetPort: intstr.FromInt32(1000)},
			},
			Selector:        nil,
			ClusterIP:       "1.2.3.4",
			SessionAffinity: corev1.ServiceAffinityNone,
			Type:            corev1.ServiceTypeClusterIP,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error adding service to the store: %v", err)
	}
	reconciler := &fakeReconciler{actions: make([]string, 0)}
	config := Config{
		EndpointReconciler: reconciler,
	}
	master := Controller{
		Config:        config,
		client:        fakeClient,
		restClient:    restClient,
		serviceLister: v1listers.NewServiceLister(serviceStore),
	}

	// healthy should ReconcileEndpoints
	healthy.Store(true)
	err = master.UpdateKubernetesService(false)
	if err != nil {
		t.Fatal(err)
	}
	lastReconcileAction := reconciler.actions[len(reconciler.actions)-1]
	if lastReconcileAction != "ReconcileEndpoints" {
		t.Fatalf("expected the reconciler to Reconcile endpoints, got %s ", lastReconcileAction)
	}
	// unhealthy should RemoveEndpoints
	healthy.Store(false)
	err = master.UpdateKubernetesService(false)
	if err != nil {
		t.Fatal(err)
	}
	lastReconcileAction = reconciler.actions[len(reconciler.actions)-1]
	if lastReconcileAction != "RemoveEndpoints" {
		t.Fatalf("expected the reconciler to Remove endpoints, got %s ", lastReconcileAction)
	}

	// healthy should ReconcileEndpoints
	healthy.Store(true)
	err = master.UpdateKubernetesService(false)
	if err != nil {
		t.Fatal(err)
	}
	lastReconcileAction = reconciler.actions[len(reconciler.actions)-1]
	if lastReconcileAction != "ReconcileEndpoints" {
		t.Fatalf("expected the reconciler to Reconcile endpoints, got %s ", lastReconcileAction)
	}

	// endpoint down should RemoveEndpoints
	ts.Close()
	err = master.UpdateKubernetesService(false)
	if err != nil {
		t.Fatal(err)
	}
	lastReconcileAction = reconciler.actions[len(reconciler.actions)-1]
	if lastReconcileAction != "RemoveEndpoints" {
		t.Fatalf("expected the reconciler to Remove endpoints, got %s ", lastReconcileAction)
	}

}

type fakeReconciler struct {
	actions []string
}

func (r *fakeReconciler) ReconcileEndpoints(serviceName string, ip net.IP, endpointPorts []corev1.EndpointPort, reconcilePorts bool) error {
	r.actions = append(r.actions, "ReconcileEndpoints")
	return nil
}

func (r *fakeReconciler) RemoveEndpoints(serviceName string, ip net.IP, endpointPorts []corev1.EndpointPort) error {
	r.actions = append(r.actions, "RemoveEndpoints")
	return nil
}

func (r *fakeReconciler) StopReconciling() {
	r.actions = append(r.actions, "StopReconciling")
}

func (r *fakeReconciler) Destroy() {
	r.actions = append(r.actions, "Destroy")
}
