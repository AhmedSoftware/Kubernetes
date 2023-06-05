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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	v1beta1 "k8s.io/api/flowcontrol/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	flowcontrolv1beta1 "k8s.io/client-go/applyconfigurations/flowcontrol/v1beta1"
	testing "k8s.io/client-go/testing"
)

// FakeFlowSchemas implements FlowSchemaInterface
type FakeFlowSchemas struct {
	Fake *FakeFlowcontrolV1beta1
}

var flowschemasResource = v1beta1.SchemeGroupVersion.WithResource("flowschemas")

var flowschemasKind = v1beta1.SchemeGroupVersion.WithKind("FlowSchema")

// Get takes name of the flowSchema, and returns the corresponding flowSchema object, and an error if there is any.
func (c *FakeFlowSchemas) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.FlowSchema, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(flowschemasResource, name), &v1beta1.FlowSchema{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.FlowSchema), err
}

// List takes label and field selectors, and returns the list of FlowSchemas that match those selectors.
func (c *FakeFlowSchemas) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.FlowSchemaList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(flowschemasResource, flowschemasKind, opts), &v1beta1.FlowSchemaList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.FlowSchemaList{ListMeta: obj.(*v1beta1.FlowSchemaList).ListMeta}
	for _, item := range obj.(*v1beta1.FlowSchemaList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested flowSchemas.
func (c *FakeFlowSchemas) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(flowschemasResource, opts))
}

// Create takes the representation of a flowSchema and creates it.  Returns the server's representation of the flowSchema, and an error, if there is any.
func (c *FakeFlowSchemas) Create(ctx context.Context, flowSchema *v1beta1.FlowSchema, opts v1.CreateOptions) (result *v1beta1.FlowSchema, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(flowschemasResource, flowSchema), &v1beta1.FlowSchema{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.FlowSchema), err
}

// Update takes the representation of a flowSchema and updates it. Returns the server's representation of the flowSchema, and an error, if there is any.
func (c *FakeFlowSchemas) Update(ctx context.Context, flowSchema *v1beta1.FlowSchema, opts v1.UpdateOptions) (result *v1beta1.FlowSchema, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(flowschemasResource, flowSchema), &v1beta1.FlowSchema{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.FlowSchema), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeFlowSchemas) UpdateStatus(ctx context.Context, flowSchema *v1beta1.FlowSchema, opts v1.UpdateOptions) (*v1beta1.FlowSchema, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(flowschemasResource, "status", flowSchema), &v1beta1.FlowSchema{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.FlowSchema), err
}

// Delete takes name of the flowSchema and deletes it. Returns an error if one occurs.
func (c *FakeFlowSchemas) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(flowschemasResource, name, opts), &v1beta1.FlowSchema{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFlowSchemas) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(flowschemasResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.FlowSchemaList{})
	return err
}

// Patch applies the patch and returns the patched flowSchema.
func (c *FakeFlowSchemas) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.FlowSchema, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(flowschemasResource, name, pt, data, subresources...), &v1beta1.FlowSchema{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.FlowSchema), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied flowSchema.
func (c *FakeFlowSchemas) Apply(ctx context.Context, flowSchema *flowcontrolv1beta1.FlowSchemaApplyConfiguration, opts v1.ApplyOptions) (result *v1beta1.FlowSchema, err error) {
	if flowSchema == nil {
		return nil, fmt.Errorf("flowSchema provided to Apply must not be nil")
	}
	data, err := json.Marshal(flowSchema)
	if err != nil {
		return nil, err
	}

	manager := "default-test-manager"
	if m := opts.FieldManager; m != "" {
		manager = m
	}

	name := flowSchema.Name
	if name == nil {
		return nil, fmt.Errorf("flowSchema.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewRootApplySubresourceAction(flowschemasResource, *name, data, manager, opts.Force), &v1beta1.FlowSchema{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.FlowSchema), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeFlowSchemas) ApplyStatus(ctx context.Context, flowSchema *flowcontrolv1beta1.FlowSchemaApplyConfiguration, opts v1.ApplyOptions) (result *v1beta1.FlowSchema, err error) {
	if flowSchema == nil {
		return nil, fmt.Errorf("flowSchema provided to Apply must not be nil")
	}
	data, err := json.Marshal(flowSchema)
	if err != nil {
		return nil, err
	}

	manager := "default-test-manager"
	if m := opts.FieldManager; m != "" {
		manager = m
	}

	name := flowSchema.Name
	if name == nil {
		return nil, fmt.Errorf("flowSchema.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewRootApplySubresourceAction(flowschemasResource, *name, data, manager, opts.Force, "status"), &v1beta1.FlowSchema{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.FlowSchema), err
}
