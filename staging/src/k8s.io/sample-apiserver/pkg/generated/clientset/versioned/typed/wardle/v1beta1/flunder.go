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

package v1beta1

import (
	context "context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
	wardlev1beta1 "k8s.io/sample-apiserver/pkg/apis/wardle/v1beta1"
	applyconfigurationwardlev1beta1 "k8s.io/sample-apiserver/pkg/generated/applyconfiguration/wardle/v1beta1"
	scheme "k8s.io/sample-apiserver/pkg/generated/clientset/versioned/scheme"
)

// FlundersGetter has a method to return a FlunderInterface.
// A group's client should implement this interface.
type FlundersGetter interface {
	Flunders(namespace string) FlunderInterface
}

// FlunderInterface has methods to work with Flunder resources.
type FlunderInterface interface {
	Create(ctx context.Context, flunder *wardlev1beta1.Flunder, opts v1.CreateOptions) (*wardlev1beta1.Flunder, error)
	Update(ctx context.Context, flunder *wardlev1beta1.Flunder, opts v1.UpdateOptions) (*wardlev1beta1.Flunder, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, flunder *wardlev1beta1.Flunder, opts v1.UpdateOptions) (*wardlev1beta1.Flunder, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*wardlev1beta1.Flunder, error)
	List(ctx context.Context, opts v1.ListOptions) (*wardlev1beta1.FlunderList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *wardlev1beta1.Flunder, err error)
	Apply(ctx context.Context, flunder *applyconfigurationwardlev1beta1.FlunderApplyConfiguration, opts v1.ApplyOptions) (result *wardlev1beta1.Flunder, err error)
	// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
	ApplyStatus(ctx context.Context, flunder *applyconfigurationwardlev1beta1.FlunderApplyConfiguration, opts v1.ApplyOptions) (result *wardlev1beta1.Flunder, err error)
	FlunderExpansion
}

// flunders implements FlunderInterface
type flunders struct {
	*gentype.ClientWithListAndApply[*wardlev1beta1.Flunder, *wardlev1beta1.FlunderList, *applyconfigurationwardlev1beta1.FlunderApplyConfiguration]
}

// newFlunders returns a Flunders
func newFlunders(c *WardleV1beta1Client, namespace string) *flunders {
	return &flunders{
		gentype.NewClientWithListAndApply[*wardlev1beta1.Flunder, *wardlev1beta1.FlunderList, *applyconfigurationwardlev1beta1.FlunderApplyConfiguration](
			"flunders",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *wardlev1beta1.Flunder { return &wardlev1beta1.Flunder{} },
			func() *wardlev1beta1.FlunderList { return &wardlev1beta1.FlunderList{} },
		),
	}
}
