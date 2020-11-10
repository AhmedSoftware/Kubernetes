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
	crv1 "k8s.io/apiextensions-apiserver/examples/client-go/pkg/apis/cr/v1"
	v1 "k8s.io/apiextensions-apiserver/examples/client-go/pkg/client/listers/cr/v1"
	cache "k8s.io/client-go/tools/cache"
)

// ExampleLister helps list Examples.
// All objects returned here must be treated as read-only.
type FakeExampleLister interface {
	v1.ExampleLister
	// Add adds the given object to the lister
	Add(obj ...*crv1.Example) error
	// Update updates the given object in the lister
	Update(obj *crv1.Example) error
	// Delete deletes the given object from lister
	Delete(obj *crv1.Example) error
}

// exampleLister implements the ExampleLister interface.
type exampleLister struct {
	index cache.Indexer
	v1.ExampleLister
}

// NewExampleLister returns a new ExampleLister.
func NewFakeExampleLister() FakeExampleLister {
	indexers := v1.NewExampleDefaultIndexer()
	index := cache.NewIndexer(cache.DeletionHandlingMetaNamespaceKeyFunc, indexers)
	lister := v1.NewExampleLister(index)
	return &exampleLister{
		index:         index,
		ExampleLister: lister,
	}
}

// Add adds the given object to the lister
func (s *exampleLister) Add(obj ...*crv1.Example) error {
	for _, curr := range obj {
		if err := s.index.Add(curr); err != nil {
			return err
		}
	}
	return nil
}

// Update updates the given object in the lister
func (s *exampleLister) Update(obj *crv1.Example) error {
	return s.index.Update(obj)
}

// Delete deletes the given object from lister
func (s *exampleLister) Delete(obj *crv1.Example) error {
	return s.index.Delete(obj)
}
