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

package secret

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// fakeManager implements Manager interface for testing purposes.
// simple operations to apiserver.
type fakeManager struct {
}

// NewFakeManager creates empty/fake secret manager
func NewFakeManager() Manager {
	return &fakeManager{}
}

// GetSecret returns a nil secret for testing
func (s *fakeManager) GetSecret(namespace, name string) (*v1.Secret, error) {
	return nil, nil
}

// RegisterPod implements the RegisterPod method for testing purposes.
func (s *fakeManager) RegisterPod(pod *v1.Pod) {
}

// UnregisterPod implements the UnregisterPod method for testing purposes.
func (s *fakeManager) UnregisterPod(pod *v1.Pod) {
}

// GetPodsNeedSyncObjects implements the GetPodsNeedSyncObjects method for testing purposes.
func (s *fakeManager) GetPodsNeedSyncObjects() []types.UID {
	return nil
}
