/*
Copyright The Kubernetes Authors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1_fakelister

import (
	networkingv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/client-go/listers/networking/v1"
	cache "k8s.io/client-go/tools/cache"
)

// NetworkPolicyLister helps list NetworkPolicies.
// All objects returned here must be treated as read-only.
type FakeNetworkPolicyLister interface {
	v1.NetworkPolicyLister
	// Add adds the given object to the lister
	Add(obj ...*networkingv1.NetworkPolicy) error
	// Update updates the given object in the lister
	Update(obj *networkingv1.NetworkPolicy) error
	// Delete deletes the given object from lister
	Delete(obj *networkingv1.NetworkPolicy) error
}

// networkPolicyLister implements the NetworkPolicyLister interface.
type networkPolicyLister struct {
	index cache.Indexer
	v1.NetworkPolicyLister
}

// NewNetworkPolicyLister returns a new NetworkPolicyLister.
func NewFakeNetworkPolicyLister() FakeNetworkPolicyLister {
	indexers := v1.NewNetworkPolicyDefaultIndexer()
	index := cache.NewIndexer(cache.DeletionHandlingMetaNamespaceKeyFunc, indexers)
	lister := v1.NewNetworkPolicyLister(index)
	return &networkPolicyLister{
		index:               index,
		NetworkPolicyLister: lister,
	}
}

// Add adds the given object to the lister
func (s *networkPolicyLister) Add(obj ...*networkingv1.NetworkPolicy) error {
	for _, curr := range obj {
		if err := s.index.Add(curr); err != nil {
			return err
		}
	}
	return nil
}

// Update updates the given object in the lister
func (s *networkPolicyLister) Update(obj *networkingv1.NetworkPolicy) error {
	return s.index.Update(obj)
}

// Delete deletes the given object from lister
func (s *networkPolicyLister) Delete(obj *networkingv1.NetworkPolicy) error {
	return s.index.Delete(obj)
}
