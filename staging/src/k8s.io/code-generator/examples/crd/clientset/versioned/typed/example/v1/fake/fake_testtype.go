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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1 "k8s.io/code-generator/examples/crd/apis/example/v1"
	examplev1 "k8s.io/code-generator/examples/crd/applyconfiguration/example/v1"
)

// FakeTestTypes implements TestTypeInterface
type FakeTestTypes struct {
	Fake *FakeExampleV1
	ns   string
}

var testtypesResource = v1.SchemeGroupVersion.WithResource("testtypes")

var testtypesKind = v1.SchemeGroupVersion.WithKind("TestType")

// Get takes name of the testType, and returns the corresponding testType object, and an error if there is any.
func (c *FakeTestTypes) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.TestType, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(testtypesResource, c.ns, name), &v1.TestType{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.TestType), err
}

// List takes label and field selectors, and returns the list of TestTypes that match those selectors.
func (c *FakeTestTypes) List(ctx context.Context, opts metav1.ListOptions) (result *v1.TestTypeList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(testtypesResource, testtypesKind, c.ns, opts), &v1.TestTypeList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.TestTypeList{ListMeta: obj.(*v1.TestTypeList).ListMeta}
	for _, item := range obj.(*v1.TestTypeList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested testTypes.
func (c *FakeTestTypes) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(testtypesResource, c.ns, opts))

}

// Create takes the representation of a testType and creates it.  Returns the server's representation of the testType, and an error, if there is any.
func (c *FakeTestTypes) Create(ctx context.Context, testType *v1.TestType, opts metav1.CreateOptions) (result *v1.TestType, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(testtypesResource, c.ns, testType), &v1.TestType{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.TestType), err
}

// Update takes the representation of a testType and updates it. Returns the server's representation of the testType, and an error, if there is any.
func (c *FakeTestTypes) Update(ctx context.Context, testType *v1.TestType, opts metav1.UpdateOptions) (result *v1.TestType, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(testtypesResource, c.ns, testType), &v1.TestType{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.TestType), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeTestTypes) UpdateStatus(ctx context.Context, testType *v1.TestType, opts metav1.UpdateOptions) (*v1.TestType, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(testtypesResource, "status", c.ns, testType), &v1.TestType{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.TestType), err
}

// Delete takes name of the testType and deletes it. Returns an error if one occurs.
func (c *FakeTestTypes) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(testtypesResource, c.ns, name, opts), &v1.TestType{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeTestTypes) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(testtypesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1.TestTypeList{})
	return err
}

// Patch applies the patch and returns the patched testType.
func (c *FakeTestTypes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.TestType, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(testtypesResource, c.ns, name, pt, data, subresources...), &v1.TestType{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.TestType), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied testType.
func (c *FakeTestTypes) Apply(ctx context.Context, testType *examplev1.TestTypeApplyConfiguration, opts metav1.ApplyOptions) (result *v1.TestType, err error) {
	if testType == nil {
		return nil, fmt.Errorf("testType provided to Apply must not be nil")
	}
	data, err := json.Marshal(testType)
	if err != nil {
		return nil, err
	}

	manager := "default-test-manager"
	if m := opts.FieldManager; m != "" {
		manager = m
	}

	name := testType.Name
	if name == nil {
		return nil, fmt.Errorf("testType.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewApplySubresourceAction(testtypesResource, c.ns, *name, data, manager, opts.Force), &v1.TestType{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.TestType), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeTestTypes) ApplyStatus(ctx context.Context, testType *examplev1.TestTypeApplyConfiguration, opts metav1.ApplyOptions) (result *v1.TestType, err error) {
	if testType == nil {
		return nil, fmt.Errorf("testType provided to Apply must not be nil")
	}
	data, err := json.Marshal(testType)
	if err != nil {
		return nil, err
	}

	manager := "default-test-manager"
	if m := opts.FieldManager; m != "" {
		manager = m
	}

	name := testType.Name
	if name == nil {
		return nil, fmt.Errorf("testType.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewApplySubresourceAction(testtypesResource, c.ns, *name, data, manager, opts.Force, "status"), &v1.TestType{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.TestType), err
}
