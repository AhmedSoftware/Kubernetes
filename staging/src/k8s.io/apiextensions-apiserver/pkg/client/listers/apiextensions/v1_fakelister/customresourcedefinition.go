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
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	v1 "k8s.io/apiextensions-apiserver/pkg/client/listers/apiextensions/v1"
	cache "k8s.io/client-go/tools/cache"
)

// CustomResourceDefinitionLister helps list CustomResourceDefinitions.
// All objects returned here must be treated as read-only.
type FakeCustomResourceDefinitionLister interface {
	v1.CustomResourceDefinitionLister
	// Add adds the given object to the lister
	Add(obj ...*apiextensionsv1.CustomResourceDefinition) error
	// Update updates the given object in the lister
	Update(obj *apiextensionsv1.CustomResourceDefinition) error
	// Delete deletes the given object from lister
	Delete(obj *apiextensionsv1.CustomResourceDefinition) error
}

// customResourceDefinitionLister implements the CustomResourceDefinitionLister interface.
type customResourceDefinitionLister struct {
	index cache.Indexer
	v1.CustomResourceDefinitionLister
}

// NewCustomResourceDefinitionLister returns a new CustomResourceDefinitionLister.
func NewFakeCustomResourceDefinitionLister() FakeCustomResourceDefinitionLister {
	indexers := v1.NewCustomResourceDefinitionDefaultIndexer()
	index := cache.NewIndexer(cache.DeletionHandlingMetaNamespaceKeyFunc, indexers)
	lister := v1.NewCustomResourceDefinitionLister(index)
	return &customResourceDefinitionLister{
		index:                          index,
		CustomResourceDefinitionLister: lister,
	}
}

// Add adds the given object to the lister
func (s *customResourceDefinitionLister) Add(obj ...*apiextensionsv1.CustomResourceDefinition) error {
	for _, curr := range obj {
		if err := s.index.Add(curr); err != nil {
			return err
		}
	}
	return nil
}

// Update updates the given object in the lister
func (s *customResourceDefinitionLister) Update(obj *apiextensionsv1.CustomResourceDefinition) error {
	return s.index.Update(obj)
}

// Delete deletes the given object from lister
func (s *customResourceDefinitionLister) Delete(obj *apiextensionsv1.CustomResourceDefinition) error {
	return s.index.Delete(obj)
}
