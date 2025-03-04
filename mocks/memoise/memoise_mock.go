// Code generated by MockGen. DO NOT EDIT.
// Source: ./memoise/memoise.go
//
// Generated by this command:
//
//	mockgen -package=memoise -source=./memoise/memoise.go -destination=mocks/memoise/memoise_mock.go
//

// Package memoise is a generated GoMock package.
package memoise

import (
	reflect "reflect"

	memoise "github.com/Asmodai/gohacks/memoise"
	gomock "go.uber.org/mock/gomock"
)

// MockMemoise is a mock of Memoise interface.
type MockMemoise struct {
	ctrl     *gomock.Controller
	recorder *MockMemoiseMockRecorder
	isgomock struct{}
}

// MockMemoiseMockRecorder is the mock recorder for MockMemoise.
type MockMemoiseMockRecorder struct {
	mock *MockMemoise
}

// NewMockMemoise creates a new mock instance.
func NewMockMemoise(ctrl *gomock.Controller) *MockMemoise {
	mock := &MockMemoise{ctrl: ctrl}
	mock.recorder = &MockMemoiseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMemoise) EXPECT() *MockMemoiseMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockMemoise) Check(arg0 string, arg1 memoise.CallbackFn) (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", arg0, arg1)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Check indicates an expected call of Check.
func (mr *MockMemoiseMockRecorder) Check(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockMemoise)(nil).Check), arg0, arg1)
}
