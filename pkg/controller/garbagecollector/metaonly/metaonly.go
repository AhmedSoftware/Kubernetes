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

package metaonly

import (
	"fmt"
	"reflect"
	"strings"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/runtime/serializer"
)

func (obj *MetadataOnlyObject) GetObjectKind() unversioned.ObjectKind     { return obj }
func (obj *MetadataOnlyObjectList) GetObjectKind() unversioned.ObjectKind { return obj }

type metaOnlyJSONScheme struct{}

// This function can be extended to mapping different gvk to different MetadataOnlyObject,
// which embedded with different version of ObjectMeta. Currently the system
// only supports v1.ObjectMeta.
func gvkToMetadataOnlyObject(gvk unversioned.GroupVersionKind) runtime.Object {
	if strings.HasSuffix(gvk.Kind, "List") {
		return &MetadataOnlyObjectList{}
	} else {
		return &MetadataOnlyObject{}
	}
}

func NewMetadataCodecFactory() serializer.CodecFactory {
	// populating another scheme from api.Scheme, registering every kind with
	// MetadataOnlyObject (or MetadataOnlyObjectList).
	scheme := runtime.NewScheme()
	allTypes := api.Scheme.AllKnownTypes()
	for kind := range allTypes {
		if kind.Version == runtime.APIVersionInternal {
			continue
		}
		metaOnlyObject := gvkToMetadataOnlyObject(kind)
		scheme.AddKnownTypeWithName(kind, metaOnlyObject)
	}
	scheme.AddUnversionedTypes(api.Unversioned, &unversioned.Status{})
	return serializer.NewCodecFactory(scheme)
}

// String converts a MetadataOnlyObject to a human-readable string.
func (metaOnly MetadataOnlyObject) String() string {
	return fmt.Sprintf("%s/%s, name: %s, DeletionTimestamp:%v", metaOnly.TypeMeta.APIVersion, metaOnly.TypeMeta.Kind, metaOnly.ObjectMeta.Name, metaOnly.ObjectMeta.DeletionTimestamp)
}

// PrintAsMetadataOnlyObject is a helper function that converts an interface{} to
// *MetadataOnlyObject and then convert it to a human-readable string.
func PrintAsMetadataOnlyObject(obj interface{}) string {
	metaOnly, ok := obj.(*MetadataOnlyObject)
	if !ok {
		return fmt.Sprintf("expected MetadataOnlyObject, got %s", reflect.TypeOf(obj))
	}
	return metaOnly.String()
}
