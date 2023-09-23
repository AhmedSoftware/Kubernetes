//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by defaulter-gen. DO NOT EDIT.

package v1beta1

import (
	v1beta1 "k8s.io/api/admissionregistration/v1beta1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	v1 "k8s.io/kubernetes/pkg/apis/admissionregistration/v1"
)

// RegisterDefaults adds defaulters functions to the given scheme.
// Public to allow building arbitrary schemes.
// All generated defaulters are covering - they call all nested defaulters.
func RegisterDefaults(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&v1beta1.MutatingWebhookConfiguration{}, func(obj interface{}) {
		SetObjectDefaults_MutatingWebhookConfiguration(obj.(*v1beta1.MutatingWebhookConfiguration))
	})
	scheme.AddTypeDefaultingFunc(&v1beta1.MutatingWebhookConfigurationList{}, func(obj interface{}) {
		SetObjectDefaults_MutatingWebhookConfigurationList(obj.(*v1beta1.MutatingWebhookConfigurationList))
	})
	scheme.AddTypeDefaultingFunc(&v1beta1.ValidatingAdmissionPolicy{}, func(obj interface{}) {
		SetObjectDefaults_ValidatingAdmissionPolicy(obj.(*v1beta1.ValidatingAdmissionPolicy))
	})
	scheme.AddTypeDefaultingFunc(&v1beta1.ValidatingAdmissionPolicyBinding{}, func(obj interface{}) {
		SetObjectDefaults_ValidatingAdmissionPolicyBinding(obj.(*v1beta1.ValidatingAdmissionPolicyBinding))
	})
	scheme.AddTypeDefaultingFunc(&v1beta1.ValidatingAdmissionPolicyBindingList{}, func(obj interface{}) {
		SetObjectDefaults_ValidatingAdmissionPolicyBindingList(obj.(*v1beta1.ValidatingAdmissionPolicyBindingList))
	})
	scheme.AddTypeDefaultingFunc(&v1beta1.ValidatingAdmissionPolicyList{}, func(obj interface{}) {
		SetObjectDefaults_ValidatingAdmissionPolicyList(obj.(*v1beta1.ValidatingAdmissionPolicyList))
	})
	scheme.AddTypeDefaultingFunc(&v1beta1.ValidatingWebhookConfiguration{}, func(obj interface{}) {
		SetObjectDefaults_ValidatingWebhookConfiguration(obj.(*v1beta1.ValidatingWebhookConfiguration))
	})
	scheme.AddTypeDefaultingFunc(&v1beta1.ValidatingWebhookConfigurationList{}, func(obj interface{}) {
		SetObjectDefaults_ValidatingWebhookConfigurationList(obj.(*v1beta1.ValidatingWebhookConfigurationList))
	})
	return nil
}

func SetObjectDefaults_MutatingWebhookConfiguration(in *v1beta1.MutatingWebhookConfiguration) {
	for i := range in.Webhooks {
		a := &in.Webhooks[i]
		SetDefaults_MutatingWebhook(a)
		if a.ClientConfig.Service != nil {
			SetDefaults_ServiceReference(a.ClientConfig.Service)
		}
		for j := range a.Rules {
			b := &a.Rules[j]
			v1.SetDefaults_Rule(&b.Rule)
		}
	}
}

func SetObjectDefaults_MutatingWebhookConfigurationList(in *v1beta1.MutatingWebhookConfigurationList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetObjectDefaults_MutatingWebhookConfiguration(a)
	}
}

func SetObjectDefaults_ValidatingAdmissionPolicy(in *v1beta1.ValidatingAdmissionPolicy) {
	if in.Spec.MatchConstraints != nil {
		for i := range in.Spec.MatchConstraints.ResourceRules {
			a := &in.Spec.MatchConstraints.ResourceRules[i]
			v1.SetDefaults_Rule(&a.RuleWithOperations.Rule)
		}
		for i := range in.Spec.MatchConstraints.ExcludeResourceRules {
			a := &in.Spec.MatchConstraints.ExcludeResourceRules[i]
			v1.SetDefaults_Rule(&a.RuleWithOperations.Rule)
		}
		SetDefaults_MatchResources(in.Spec.MatchConstraints)
	}
	SetDefaults_ValidatingAdmissionPolicySpec(&in.Spec)
}

func SetObjectDefaults_ValidatingAdmissionPolicyBinding(in *v1beta1.ValidatingAdmissionPolicyBinding) {
	if in.Spec.MatchResources != nil {
		for i := range in.Spec.MatchResources.ResourceRules {
			a := &in.Spec.MatchResources.ResourceRules[i]
			v1.SetDefaults_Rule(&a.RuleWithOperations.Rule)
		}
		for i := range in.Spec.MatchResources.ExcludeResourceRules {
			a := &in.Spec.MatchResources.ExcludeResourceRules[i]
			v1.SetDefaults_Rule(&a.RuleWithOperations.Rule)
		}
		SetDefaults_MatchResources(in.Spec.MatchResources)
	}
}

func SetObjectDefaults_ValidatingAdmissionPolicyBindingList(in *v1beta1.ValidatingAdmissionPolicyBindingList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetObjectDefaults_ValidatingAdmissionPolicyBinding(a)
	}
}

func SetObjectDefaults_ValidatingAdmissionPolicyList(in *v1beta1.ValidatingAdmissionPolicyList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetObjectDefaults_ValidatingAdmissionPolicy(a)
	}
}

func SetObjectDefaults_ValidatingWebhookConfiguration(in *v1beta1.ValidatingWebhookConfiguration) {
	for i := range in.Webhooks {
		a := &in.Webhooks[i]
		SetDefaults_ValidatingWebhook(a)
		if a.ClientConfig.Service != nil {
			SetDefaults_ServiceReference(a.ClientConfig.Service)
		}
		for j := range a.Rules {
			b := &a.Rules[j]
			v1.SetDefaults_Rule(&b.Rule)
		}
	}
}

func SetObjectDefaults_ValidatingWebhookConfigurationList(in *v1beta1.ValidatingWebhookConfigurationList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetObjectDefaults_ValidatingWebhookConfiguration(a)
	}
}
