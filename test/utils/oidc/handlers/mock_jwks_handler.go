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

// Code generated by mockery v2.40.3. DO NOT EDIT.

package handlers

import (
	mock "github.com/stretchr/testify/mock"
	jose "gopkg.in/square/go-jose.v2"
)

// MockJWKsHandler is an autogenerated mock type for the JWKsHandler type
type MockJWKsHandler struct {
	mock.Mock
}

type MockJWKsHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *MockJWKsHandler) EXPECT() *MockJWKsHandler_Expecter {
	return &MockJWKsHandler_Expecter{mock: &_m.Mock}
}

// KeySet provides a mock function with given fields:
func (_m *MockJWKsHandler) KeySet() jose.JSONWebKeySet {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for KeySet")
	}

	var r0 jose.JSONWebKeySet
	if rf, ok := ret.Get(0).(func() jose.JSONWebKeySet); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(jose.JSONWebKeySet)
	}

	return r0
}

// MockJWKsHandler_KeySet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'KeySet'
type MockJWKsHandler_KeySet_Call struct {
	*mock.Call
}

// KeySet is a helper method to define mock.On call
func (_e *MockJWKsHandler_Expecter) KeySet() *MockJWKsHandler_KeySet_Call {
	return &MockJWKsHandler_KeySet_Call{Call: _e.mock.On("KeySet")}
}

func (_c *MockJWKsHandler_KeySet_Call) Run(run func()) *MockJWKsHandler_KeySet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockJWKsHandler_KeySet_Call) Return(_a0 jose.JSONWebKeySet) *MockJWKsHandler_KeySet_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockJWKsHandler_KeySet_Call) RunAndReturn(run func() jose.JSONWebKeySet) *MockJWKsHandler_KeySet_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockJWKsHandler creates a new instance of MockJWKsHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockJWKsHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockJWKsHandler {
	mock := &MockJWKsHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
