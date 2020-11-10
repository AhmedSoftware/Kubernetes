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
	certificatesv1 "k8s.io/api/certificates/v1"
	v1 "k8s.io/client-go/listers/certificates/v1"
	cache "k8s.io/client-go/tools/cache"
)

// CertificateSigningRequestLister helps list CertificateSigningRequests.
// All objects returned here must be treated as read-only.
type FakeCertificateSigningRequestLister interface {
	v1.CertificateSigningRequestLister
	// Add adds the given object to the lister
	Add(obj ...*certificatesv1.CertificateSigningRequest) error
	// Update updates the given object in the lister
	Update(obj *certificatesv1.CertificateSigningRequest) error
	// Delete deletes the given object from lister
	Delete(obj *certificatesv1.CertificateSigningRequest) error
}

// certificateSigningRequestLister implements the CertificateSigningRequestLister interface.
type certificateSigningRequestLister struct {
	index cache.Indexer
	v1.CertificateSigningRequestLister
}

// NewCertificateSigningRequestLister returns a new CertificateSigningRequestLister.
func NewFakeCertificateSigningRequestLister() FakeCertificateSigningRequestLister {
	indexers := v1.NewCertificateSigningRequestDefaultIndexer()
	index := cache.NewIndexer(cache.DeletionHandlingMetaNamespaceKeyFunc, indexers)
	lister := v1.NewCertificateSigningRequestLister(index)
	return &certificateSigningRequestLister{
		index:                           index,
		CertificateSigningRequestLister: lister,
	}
}

// Add adds the given object to the lister
func (s *certificateSigningRequestLister) Add(obj ...*certificatesv1.CertificateSigningRequest) error {
	for _, curr := range obj {
		if err := s.index.Add(curr); err != nil {
			return err
		}
	}
	return nil
}

// Update updates the given object in the lister
func (s *certificateSigningRequestLister) Update(obj *certificatesv1.CertificateSigningRequest) error {
	return s.index.Update(obj)
}

// Delete deletes the given object from lister
func (s *certificateSigningRequestLister) Delete(obj *certificatesv1.CertificateSigningRequest) error {
	return s.index.Delete(obj)
}
