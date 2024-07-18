// Code generated by MockGen. DO NOT EDIT.
// Source: ./events/interface.go
// +build testing
//
// Generated by this command:
//
//	mockgen -package=events -source=./events/interface.go -destination=events/interface_mock.go
//

// Package events is a generated GoMock package.
package events

import (
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockEvent is a mock of Event interface.
type MockEvent struct {
	ctrl     *gomock.Controller
	recorder *MockEventMockRecorder
}

// MockEventMockRecorder is the mock recorder for MockEvent.
type MockEventMockRecorder struct {
	mock *MockEvent
}

// NewMockEvent creates a new mock instance.
func NewMockEvent(ctrl *gomock.Controller) *MockEvent {
	mock := &MockEvent{ctrl: ctrl}
	mock.recorder = &MockEventMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEvent) EXPECT() *MockEventMockRecorder {
	return m.recorder
}

// String mocks base method.
func (m *MockEvent) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockEventMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockEvent)(nil).String))
}

// When mocks base method.
func (m *MockEvent) When() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "When")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// When indicates an expected call of When.
func (mr *MockEventMockRecorder) When() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "When", reflect.TypeOf((*MockEvent)(nil).When))
}
