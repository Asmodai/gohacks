// Code generated by MockGen. DO NOT EDIT.
// Source: ./database/databasemgr.go
//
// Generated by this command:
//
//	mockgen -package=database -source=./database/databasemgr.go -destination=mocks/database/databasemgr_mock.go
//

// Package database is a generated GoMock package.
package database

import (
	reflect "reflect"

	database "github.com/Asmodai/gohacks/database"
	gomock "go.uber.org/mock/gomock"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
	isgomock struct{}
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// CheckDB mocks base method.
func (m *MockManager) CheckDB(arg0 database.Database) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckDB", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckDB indicates an expected call of CheckDB.
func (mr *MockManagerMockRecorder) CheckDB(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckDB", reflect.TypeOf((*MockManager)(nil).CheckDB), arg0)
}

// Open mocks base method.
func (m *MockManager) Open(arg0, arg1 string) (database.Database, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open", arg0, arg1)
	ret0, _ := ret[0].(database.Database)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Open indicates an expected call of Open.
func (mr *MockManagerMockRecorder) Open(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockManager)(nil).Open), arg0, arg1)
}

// OpenConfig mocks base method.
func (m *MockManager) OpenConfig(arg0 *database.Config) (database.Database, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenConfig", arg0)
	ret0, _ := ret[0].(database.Database)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenConfig indicates an expected call of OpenConfig.
func (mr *MockManagerMockRecorder) OpenConfig(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenConfig", reflect.TypeOf((*MockManager)(nil).OpenConfig), arg0)
}
