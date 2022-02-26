/*
Copyright 2022 The Kubernetes Authors.

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

package resourceclaim

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/core/validation"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

// resourceclaimStrategy implements behavior for ResourceClaim objects
type resourceclaimStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating
// ResourceClaim objects via the REST API.
var Strategy = resourceclaimStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

func (resourceclaimStrategy) NamespaceScoped() bool {
	return true
}

// GetResetFields returns the set of fields that get reset by the strategy and
// should not be modified by the user. For a new ResourceClaim that is the
// status.
func (resourceclaimStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	fields := map[fieldpath.APIVersion]*fieldpath.Set{
		"v1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("status"),
		),
	}

	return fields
}

func (resourceclaimStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	claim := obj.(*core.ResourceClaim)
	// Status must not be set by user on create.
	claim.Status = core.ResourceClaimStatus{}
}

func (resourceclaimStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	claim := obj.(*core.ResourceClaim)
	return validation.ValidateResourceClaim(claim)
}

func (resourceclaimStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (resourceclaimStrategy) Canonicalize(obj runtime.Object) {
}

func (resourceclaimStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (resourceclaimStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newClaim := obj.(*core.ResourceClaim)
	oldClaim := old.(*core.ResourceClaim)
	newClaim.Status = oldClaim.Status
}

func (resourceclaimStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	newClaim := obj.(*core.ResourceClaim)
	oldClaim := old.(*core.ResourceClaim)
	errorList := validation.ValidateResourceClaim(newClaim)
	return append(errorList, validation.ValidateResourceClaimUpdate(newClaim, oldClaim)...)
}

func (resourceclaimStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

func (resourceclaimStrategy) AllowUnconditionalUpdate() bool {
	return true
}

type resourceclaimStatusStrategy struct {
	resourceclaimStrategy
}

var StatusStrategy = resourceclaimStatusStrategy{Strategy}

// GetResetFields returns the set of fields that get reset by the strategy and
// should not be modified by the user. For a status update that is the spec.
func (resourceclaimStatusStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	fields := map[fieldpath.APIVersion]*fieldpath.Set{
		"v1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
		),
	}

	return fields
}

func (resourceclaimStatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (resourceclaimStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	newClaim := obj.(*core.ResourceClaim)
	oldClaim := old.(*core.ResourceClaim)
	return validation.ValidateResourceClaimStatusUpdate(newClaim, oldClaim)
}

// WarningsOnUpdate returns warnings for the given update.
func (resourceclaimStatusStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

// Match returns a generic matcher for a given label and field selector.
func Match(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	claim, ok := obj.(*core.ResourceClaim)
	if !ok {
		return nil, nil, errors.New("not a resourceclaim")
	}
	return labels.Set(claim.Labels), toSelectableFields(claim), nil
}

// toSelectableFields returns a field set that represents the object
func toSelectableFields(claim *core.ResourceClaim) fields.Set {
	fields := generic.ObjectMetaFieldsSet(&claim.ObjectMeta, true)
	return fields
}
