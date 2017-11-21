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

package garbagecollector

import (
	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

// apiResource consults the REST mapper to translate an <apiVersion, kind,
// namespace> tuple to a unversioned.APIResource struct.
func (gc *GarbageCollector) apiResource(apiVersion, kind string, namespaced bool) (*metav1.APIResource, error) {
	fqKind := schema.FromAPIVersionAndKind(apiVersion, kind)
	mapping, err := gc.restMapper.RESTMapping(fqKind.GroupKind(), apiVersion)
	if err != nil {
		return nil, newRESTMappingError(kind, apiVersion)
	}
	glog.V(5).Infof("map kind %s, version %s to resource %s", kind, apiVersion, mapping.Resource)
	resource := metav1.APIResource{
		Name:       mapping.Resource,
		Namespaced: namespaced,
		Kind:       kind,
	}
	return &resource, nil
}

func (gc *GarbageCollector) deleteObject(item objectReference, policy *metav1.DeletionPropagation) error {
	fqKind := schema.FromAPIVersionAndKind(item.APIVersion, item.Kind)
	client, err := gc.clientPool.ClientForGroupVersionKind(fqKind)
	if err != nil {
		return err
	}
	resource, err := gc.apiResource(item.APIVersion, item.Kind, len(item.Namespace) != 0)
	if err != nil {
		return err
	}
	uid := item.UID
	preconditions := metav1.Preconditions{UID: &uid}
	deleteOptions := metav1.DeleteOptions{Preconditions: &preconditions, PropagationPolicy: policy}
	return client.Resource(resource, item.Namespace).Delete(item.Name, &deleteOptions)
}

func (gc *GarbageCollector) getObject(item objectReference) (*unstructured.Unstructured, error) {
	fqKind := schema.FromAPIVersionAndKind(item.APIVersion, item.Kind)
	client, err := gc.clientPool.ClientForGroupVersionKind(fqKind)
	if err != nil {
		return nil, err
	}
	resource, err := gc.apiResource(item.APIVersion, item.Kind, len(item.Namespace) != 0)
	if err != nil {
		return nil, err
	}
	return client.Resource(resource, item.Namespace).Get(item.Name, metav1.GetOptions{})
}

func (gc *GarbageCollector) updateObject(item objectReference, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	fqKind := schema.FromAPIVersionAndKind(item.APIVersion, item.Kind)
	client, err := gc.clientPool.ClientForGroupVersionKind(fqKind)
	if err != nil {
		return nil, err
	}
	resource, err := gc.apiResource(item.APIVersion, item.Kind, len(item.Namespace) != 0)
	if err != nil {
		return nil, err
	}
	return client.Resource(resource, item.Namespace).Update(obj)
}

func (gc *GarbageCollector) patchObject(item objectReference, patch []byte) (*unstructured.Unstructured, error) {
	fqKind := schema.FromAPIVersionAndKind(item.APIVersion, item.Kind)
	client, err := gc.clientPool.ClientForGroupVersionKind(fqKind)
	if err != nil {
		return nil, err
	}
	resource, err := gc.apiResource(item.APIVersion, item.Kind, len(item.Namespace) != 0)
	if err != nil {
		return nil, err
	}
	return client.Resource(resource, item.Namespace).Patch(item.Name, types.StrategicMergePatchType, patch)
}
