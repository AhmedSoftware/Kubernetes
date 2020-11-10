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
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/client-go/listers/admissionregistration/v1"
	cache "k8s.io/client-go/tools/cache"
)

// MutatingWebhookConfigurationLister helps list MutatingWebhookConfigurations.
// All objects returned here must be treated as read-only.
type FakeMutatingWebhookConfigurationLister interface {
	v1.MutatingWebhookConfigurationLister
	// Add adds the given object to the lister
	Add(obj ...*admissionregistrationv1.MutatingWebhookConfiguration) error
	// Update updates the given object in the lister
	Update(obj *admissionregistrationv1.MutatingWebhookConfiguration) error
	// Delete deletes the given object from lister
	Delete(obj *admissionregistrationv1.MutatingWebhookConfiguration) error
}

// mutatingWebhookConfigurationLister implements the MutatingWebhookConfigurationLister interface.
type mutatingWebhookConfigurationLister struct {
	index cache.Indexer
	v1.MutatingWebhookConfigurationLister
}

// NewMutatingWebhookConfigurationLister returns a new MutatingWebhookConfigurationLister.
func NewFakeMutatingWebhookConfigurationLister() FakeMutatingWebhookConfigurationLister {
	indexers := v1.NewMutatingWebhookConfigurationDefaultIndexer()
	index := cache.NewIndexer(cache.DeletionHandlingMetaNamespaceKeyFunc, indexers)
	lister := v1.NewMutatingWebhookConfigurationLister(index)
	return &mutatingWebhookConfigurationLister{
		index:                              index,
		MutatingWebhookConfigurationLister: lister,
	}
}

// Add adds the given object to the lister
func (s *mutatingWebhookConfigurationLister) Add(obj ...*admissionregistrationv1.MutatingWebhookConfiguration) error {
	for _, curr := range obj {
		if err := s.index.Add(curr); err != nil {
			return err
		}
	}
	return nil
}

// Update updates the given object in the lister
func (s *mutatingWebhookConfigurationLister) Update(obj *admissionregistrationv1.MutatingWebhookConfiguration) error {
	return s.index.Update(obj)
}

// Delete deletes the given object from lister
func (s *mutatingWebhookConfigurationLister) Delete(obj *admissionregistrationv1.MutatingWebhookConfiguration) error {
	return s.index.Delete(obj)
}
