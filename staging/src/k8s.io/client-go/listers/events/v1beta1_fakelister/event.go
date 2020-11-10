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

package v1beta1_fakelister

import (
	eventsv1beta1 "k8s.io/api/events/v1beta1"
	v1beta1 "k8s.io/client-go/listers/events/v1beta1"
	cache "k8s.io/client-go/tools/cache"
)

// EventLister helps list Events.
// All objects returned here must be treated as read-only.
type FakeEventLister interface {
	v1beta1.EventLister
	// Add adds the given object to the lister
	Add(obj ...*eventsv1beta1.Event) error
	// Update updates the given object in the lister
	Update(obj *eventsv1beta1.Event) error
	// Delete deletes the given object from lister
	Delete(obj *eventsv1beta1.Event) error
}

// eventLister implements the EventLister interface.
type eventLister struct {
	index cache.Indexer
	v1beta1.EventLister
}

// NewEventLister returns a new EventLister.
func NewFakeEventLister() FakeEventLister {
	indexers := v1beta1.NewEventDefaultIndexer()
	index := cache.NewIndexer(cache.DeletionHandlingMetaNamespaceKeyFunc, indexers)
	lister := v1beta1.NewEventLister(index)
	return &eventLister{
		index:       index,
		EventLister: lister,
	}
}

// Add adds the given object to the lister
func (s *eventLister) Add(obj ...*eventsv1beta1.Event) error {
	for _, curr := range obj {
		if err := s.index.Add(curr); err != nil {
			return err
		}
	}
	return nil
}

// Update updates the given object in the lister
func (s *eventLister) Update(obj *eventsv1beta1.Event) error {
	return s.index.Update(obj)
}

// Delete deletes the given object from lister
func (s *eventLister) Delete(obj *eventsv1beta1.Event) error {
	return s.index.Delete(obj)
}
