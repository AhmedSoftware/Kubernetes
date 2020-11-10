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
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/client-go/listers/core/v1"
	cache "k8s.io/client-go/tools/cache"
)

// ResourceQuotaLister helps list ResourceQuotas.
// All objects returned here must be treated as read-only.
type FakeResourceQuotaLister interface {
	v1.ResourceQuotaLister
	// Add adds the given object to the lister
	Add(obj ...*corev1.ResourceQuota) error
	// Update updates the given object in the lister
	Update(obj *corev1.ResourceQuota) error
	// Delete deletes the given object from lister
	Delete(obj *corev1.ResourceQuota) error
}

// resourceQuotaLister implements the ResourceQuotaLister interface.
type resourceQuotaLister struct {
	index cache.Indexer
	v1.ResourceQuotaLister
}

// NewResourceQuotaLister returns a new ResourceQuotaLister.
func NewFakeResourceQuotaLister() FakeResourceQuotaLister {
	indexers := v1.NewResourceQuotaDefaultIndexer()
	index := cache.NewIndexer(cache.DeletionHandlingMetaNamespaceKeyFunc, indexers)
	lister := v1.NewResourceQuotaLister(index)
	return &resourceQuotaLister{
		index:               index,
		ResourceQuotaLister: lister,
	}
}

// Add adds the given object to the lister
func (s *resourceQuotaLister) Add(obj ...*corev1.ResourceQuota) error {
	for _, curr := range obj {
		if err := s.index.Add(curr); err != nil {
			return err
		}
	}
	return nil
}

// Update updates the given object in the lister
func (s *resourceQuotaLister) Update(obj *corev1.ResourceQuota) error {
	return s.index.Update(obj)
}

// Delete deletes the given object from lister
func (s *resourceQuotaLister) Delete(obj *corev1.ResourceQuota) error {
	return s.index.Delete(obj)
}
