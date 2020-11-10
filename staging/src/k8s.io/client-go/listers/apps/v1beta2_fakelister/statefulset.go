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

package v1beta2_fakelister

import (
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	v1beta2 "k8s.io/client-go/listers/apps/v1beta2"
	cache "k8s.io/client-go/tools/cache"
)

// StatefulSetLister helps list StatefulSets.
// All objects returned here must be treated as read-only.
type FakeStatefulSetLister interface {
	v1beta2.StatefulSetLister
	// Add adds the given object to the lister
	Add(obj ...*appsv1beta2.StatefulSet) error
	// Update updates the given object in the lister
	Update(obj *appsv1beta2.StatefulSet) error
	// Delete deletes the given object from lister
	Delete(obj *appsv1beta2.StatefulSet) error
}

// statefulSetLister implements the StatefulSetLister interface.
type statefulSetLister struct {
	index cache.Indexer
	v1beta2.StatefulSetLister
}

// NewStatefulSetLister returns a new StatefulSetLister.
func NewFakeStatefulSetLister() FakeStatefulSetLister {
	indexers := v1beta2.NewStatefulSetDefaultIndexer()
	index := cache.NewIndexer(cache.DeletionHandlingMetaNamespaceKeyFunc, indexers)
	lister := v1beta2.NewStatefulSetLister(index)
	return &statefulSetLister{
		index:             index,
		StatefulSetLister: lister,
	}
}

// Add adds the given object to the lister
func (s *statefulSetLister) Add(obj ...*appsv1beta2.StatefulSet) error {
	for _, curr := range obj {
		if err := s.index.Add(curr); err != nil {
			return err
		}
	}
	return nil
}

// Update updates the given object in the lister
func (s *statefulSetLister) Update(obj *appsv1beta2.StatefulSet) error {
	return s.index.Update(obj)
}

// Delete deletes the given object from lister
func (s *statefulSetLister) Delete(obj *appsv1beta2.StatefulSet) error {
	return s.index.Delete(obj)
}
