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

package podschedulinghints

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/kubernetes/pkg/apis/resource"
)

var hints = &resource.PodSchedulingHints{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "valid-pod",
		Namespace: "default",
	},
	Spec: resource.PodSchedulingHintsSpec{
		SelectedNode: "worker",
	},
}

func TestPodSchedulingStrategy(t *testing.T) {
	if !Strategy.NamespaceScoped() {
		t.Errorf("PodSchedulingHints must be namespace scoped")
	}
	if Strategy.AllowCreateOnUpdate() {
		t.Errorf("PodSchedulingHints should not allow create on update")
	}
}

func TestPodSchedulingStrategyCreate(t *testing.T) {
	ctx := genericapirequest.NewDefaultContext()
	hints := hints.DeepCopy()

	Strategy.PrepareForCreate(ctx, hints)
	errs := Strategy.Validate(ctx, hints)
	if len(errs) != 0 {
		t.Errorf("unexpected error validating for create %v", errs)
	}
}

func TestPodSchedulingStrategyUpdate(t *testing.T) {
	t.Run("no-changes-okay", func(t *testing.T) {
		ctx := genericapirequest.NewDefaultContext()
		hints := hints.DeepCopy()
		newHints := hints.DeepCopy()
		newHints.ResourceVersion = "4"

		Strategy.PrepareForUpdate(ctx, newHints, hints)
		errs := Strategy.ValidateUpdate(ctx, newHints, hints)
		if len(errs) != 0 {
			t.Errorf("unexpected validation errors: %v", errs)
		}
	})

	t.Run("name-change-not-allowed", func(t *testing.T) {
		ctx := genericapirequest.NewDefaultContext()
		hints := hints.DeepCopy()
		newHints := hints.DeepCopy()
		newHints.Name = "valid-claim-2"
		newHints.ResourceVersion = "4"

		Strategy.PrepareForUpdate(ctx, newHints, hints)
		errs := Strategy.ValidateUpdate(ctx, newHints, hints)
		if len(errs) == 0 {
			t.Errorf("expected a validation error")
		}
	})
}
