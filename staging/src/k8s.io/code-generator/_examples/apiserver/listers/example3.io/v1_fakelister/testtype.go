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
	cache "k8s.io/client-go/tools/cache"
	example3iov1 "k8s.io/code-generator/_examples/apiserver/apis/example3.io/v1"
	v1 "k8s.io/code-generator/_examples/apiserver/listers/example3.io/v1"
)

// TestTypeLister helps list TestTypes.
// All objects returned here must be treated as read-only.
type FakeTestTypeLister interface {
	v1.TestTypeLister
	// Add adds the given object to the lister
	Add(obj ...*example3iov1.TestType) error
	// Update updates the given object in the lister
	Update(obj *example3iov1.TestType) error
	// Delete deletes the given object from lister
	Delete(obj *example3iov1.TestType) error
}

// testTypeLister implements the TestTypeLister interface.
type testTypeLister struct {
	index cache.Indexer
	v1.TestTypeLister
}

// NewTestTypeLister returns a new TestTypeLister.
func NewFakeTestTypeLister() FakeTestTypeLister {
	indexers := v1.NewTestTypeDefaultIndexer()
	index := cache.NewIndexer(cache.DeletionHandlingMetaNamespaceKeyFunc, indexers)
	lister := v1.NewTestTypeLister(index)
	return &testTypeLister{
		index:          index,
		TestTypeLister: lister,
	}
}

// Add adds the given object to the lister
func (s *testTypeLister) Add(obj ...*example3iov1.TestType) error {
	for _, curr := range obj {
		if err := s.index.Add(curr); err != nil {
			return err
		}
	}
	return nil
}

// Update updates the given object in the lister
func (s *testTypeLister) Update(obj *example3iov1.TestType) error {
	return s.index.Update(obj)
}

// Delete deletes the given object from lister
func (s *testTypeLister) Delete(obj *example3iov1.TestType) error {
	return s.index.Delete(obj)
}
