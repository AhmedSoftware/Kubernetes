/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

package schedulercache

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

// CreateNodeNameToInfoMap obtains a list of pods and pivots that list into a map where the keys are node names
// and the values are the aggregated information for that node.
func CreateNodeNameToInfoMap(pods []*api.Pod) map[string]*NodeInfo {
	nodeNameToInfo := make(map[string]*NodeInfo)
	for _, pod := range pods {
		nodeName := pod.Spec.NodeName
		nodeInfo, ok := nodeNameToInfo[nodeName]
		if !ok {
			nodeInfo = NewNodeInfo()
			nodeNameToInfo[nodeName] = nodeInfo
		}
		nodeInfo.addPod(pod)
	}
	return nodeNameToInfo
}

type ReferenceWithSelector struct {
	Reference *api.ObjectReference
	Selector  labels.Selector
}

// CreateNodeNameToInfoMapWithGroupingObjects objects calls CreateNodeNameToInfoMap and additionally
// updates references in NodeInfo structs based on the given pairs of references and selectors.
func CreateNodeNameToInfoMapWithGroupingObjects(pods []*api.Pod, refs []*ReferenceWithSelector) map[string]*NodeInfo {
	result := CreateNodeNameToInfoMap(pods)
	for _, ref := range refs {
		for _, pod := range pods {
			if pod.Namespace != ref.Reference.Namespace {
				continue
			}
			if ref.Selector.Matches(labels.Set(pod.Labels)) {
				result[pod.Spec.NodeName].AddReference(ref.Reference)
			}
		}
	}
	return result
}
