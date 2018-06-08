/*
Copyright 2016 The Kubernetes Authors.

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

package cronjob

import (
	"context"

	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/api/pod"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/batch/validation"
	corevalidation "k8s.io/kubernetes/pkg/apis/core/validation"
)

// cronJobStrategy implements verification logic for Replication Controllers.
type cronJobStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating CronJob objects.
var Strategy = cronJobStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// DefaultGarbageCollectionPolicy returns Orphan because that was the default
// behavior before the server-side garbage collection was implemented.
func (cronJobStrategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.OrphanDependents
}

// NamespaceScoped returns true because all scheduled jobs need to be within a namespace.
func (cronJobStrategy) NamespaceScoped() bool {
	return true
}

// PrepareForCreate clears the status of a scheduled job before creation.
func (cronJobStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	cronJob := obj.(*batch.CronJob)
	cronJob.Status = batch.CronJobStatus{}

	pod.DropDisabledAlphaFields(&cronJob.Spec.JobTemplate.Spec.Template.Spec)
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (cronJobStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newCronJob := obj.(*batch.CronJob)
	oldCronJob := old.(*batch.CronJob)
	newCronJob.Status = oldCronJob.Status

	pod.DropDisabledAlphaFields(&newCronJob.Spec.JobTemplate.Spec.Template.Spec)
	pod.DropDisabledAlphaFields(&oldCronJob.Spec.JobTemplate.Spec.Template.Spec)
}

// Validate validates a new scheduled job.
func (cronJobStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	cronJob := obj.(*batch.CronJob)
	allErrs := validation.ValidateCronJob(cronJob)
	allErrs = append(allErrs, apimachineryvalidation.ValidateRatchetingCreate(
		&cronJob.Spec.JobTemplate.Spec.Template.Spec,
		field.NewPath("spec", "template", "spec", "template", "spec"),
		corevalidation.RatchetingPodSpecValidations)...,
	)
	return allErrs
}

// Canonicalize normalizes the object after validation.
func (cronJobStrategy) Canonicalize(obj runtime.Object) {
}

func (cronJobStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// AllowCreateOnUpdate is false for scheduled jobs; this means a POST is needed to create one.
func (cronJobStrategy) AllowCreateOnUpdate() bool {
	return false
}

// ValidateUpdate is the default update validation for an end user.
func (cronJobStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	cronJob := obj.(*batch.CronJob)
	oldCronJob := old.(*batch.CronJob)
	allErrs := validation.ValidateCronJobUpdate(cronJob, oldCronJob)
	allErrs = append(allErrs, apimachineryvalidation.ValidateRatchetingUpdate(
		&cronJob.Spec.JobTemplate.Spec.Template.Spec,
		&oldCronJob.Spec.JobTemplate.Spec.Template.Spec,
		field.NewPath("spec", "template", "spec", "template", "spec"),
		corevalidation.RatchetingPodSpecValidations)...,
	)
	return allErrs
}

type cronJobStatusStrategy struct {
	cronJobStrategy
}

var StatusStrategy = cronJobStatusStrategy{Strategy}

func (cronJobStatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newJob := obj.(*batch.CronJob)
	oldJob := old.(*batch.CronJob)
	newJob.Spec = oldJob.Spec
}

func (cronJobStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}
