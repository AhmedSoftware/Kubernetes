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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	v1 "k8s.io/client-go/applyconfigurations/core/v1"
)

// IngressClassSpecApplyConfiguration represents an declarative configuration of the IngressClassSpec type for use
// with apply.
type IngressClassSpecApplyConfiguration struct {
	Controller *string                                         `json:"controller,omitempty"`
	Parameters *v1.TypedLocalObjectReferenceApplyConfiguration `json:"parameters,omitempty"`
}

// IngressClassSpecApplyConfiguration constructs an declarative configuration of the IngressClassSpec type for use with
// apply.
func IngressClassSpec() *IngressClassSpecApplyConfiguration {
	return &IngressClassSpecApplyConfiguration{}
}

// WithController sets the Controller field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Controller field is set to the value of the last call.
func (b *IngressClassSpecApplyConfiguration) WithController(value string) *IngressClassSpecApplyConfiguration {
	b.Controller = &value
	return b
}

// WithParameters sets the Parameters field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Parameters field is set to the value of the last call.
func (b *IngressClassSpecApplyConfiguration) WithParameters(value *v1.TypedLocalObjectReferenceApplyConfiguration) *IngressClassSpecApplyConfiguration {
	b.Parameters = value
	return b
}
