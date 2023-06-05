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

	v1alpha2 "k8s.io/api/resource/v1alpha2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	resourcev1alpha2 "k8s.io/client-go/applyconfigurations/resource/v1alpha2"
	testing "k8s.io/client-go/testing"
)

// FakeResourceClasses implements ResourceClassInterface
type FakeResourceClasses struct {
	Fake *FakeResourceV1alpha2
}

var resourceclassesResource = v1alpha2.SchemeGroupVersion.WithResource("resourceclasses")

var resourceclassesKind = v1alpha2.SchemeGroupVersion.WithKind("ResourceClass")

// Get takes name of the resourceClass, and returns the corresponding resourceClass object, and an error if there is any.
func (c *FakeResourceClasses) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha2.ResourceClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(resourceclassesResource, name), &v1alpha2.ResourceClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.ResourceClass), err
}

// List takes label and field selectors, and returns the list of ResourceClasses that match those selectors.
func (c *FakeResourceClasses) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha2.ResourceClassList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(resourceclassesResource, resourceclassesKind, opts), &v1alpha2.ResourceClassList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha2.ResourceClassList{ListMeta: obj.(*v1alpha2.ResourceClassList).ListMeta}
	for _, item := range obj.(*v1alpha2.ResourceClassList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested resourceClasses.
func (c *FakeResourceClasses) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(resourceclassesResource, opts))
}

// Create takes the representation of a resourceClass and creates it.  Returns the server's representation of the resourceClass, and an error, if there is any.
func (c *FakeResourceClasses) Create(ctx context.Context, resourceClass *v1alpha2.ResourceClass, opts v1.CreateOptions) (result *v1alpha2.ResourceClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(resourceclassesResource, resourceClass), &v1alpha2.ResourceClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.ResourceClass), err
}

// Update takes the representation of a resourceClass and updates it. Returns the server's representation of the resourceClass, and an error, if there is any.
func (c *FakeResourceClasses) Update(ctx context.Context, resourceClass *v1alpha2.ResourceClass, opts v1.UpdateOptions) (result *v1alpha2.ResourceClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(resourceclassesResource, resourceClass), &v1alpha2.ResourceClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.ResourceClass), err
}

// Delete takes name of the resourceClass and deletes it. Returns an error if one occurs.
func (c *FakeResourceClasses) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(resourceclassesResource, name, opts), &v1alpha2.ResourceClass{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeResourceClasses) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(resourceclassesResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha2.ResourceClassList{})
	return err
}

// Patch applies the patch and returns the patched resourceClass.
func (c *FakeResourceClasses) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha2.ResourceClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(resourceclassesResource, name, pt, data, subresources...), &v1alpha2.ResourceClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.ResourceClass), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied resourceClass.
func (c *FakeResourceClasses) Apply(ctx context.Context, resourceClass *resourcev1alpha2.ResourceClassApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha2.ResourceClass, err error) {
	if resourceClass == nil {
		return nil, fmt.Errorf("resourceClass provided to Apply must not be nil")
	}
	data, err := json.Marshal(resourceClass)
	if err != nil {
		return nil, err
	}

	manager := "default-test-manager"
	if m := opts.FieldManager; m != "" {
		manager = m
	}

	name := resourceClass.Name
	if name == nil {
		return nil, fmt.Errorf("resourceClass.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewRootApplySubresourceAction(resourceclassesResource, *name, data, manager, opts.Force), &v1alpha2.ResourceClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.ResourceClass), err
}
