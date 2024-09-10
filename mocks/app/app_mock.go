// Code generated by MockGen. DO NOT EDIT.
// Source: ./app/app.go
// +build testing
//
// Generated by this command:
//
//	mockgen -package=app -source=./app/app.go -destination=mocks/app/app_mock.go
//

// Package app is a generated GoMock package.
package app

import (
	context "context"
	reflect "reflect"

	app "github.com/Asmodai/gohacks/app"
	config "github.com/Asmodai/gohacks/config"
	logger "github.com/Asmodai/gohacks/logger"
	process "github.com/Asmodai/gohacks/process"
	semver "github.com/Asmodai/gohacks/semver"
	gomock "go.uber.org/mock/gomock"
)

// MockApplication is a mock of Application interface.
type MockApplication struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationMockRecorder
}

// MockApplicationMockRecorder is the mock recorder for MockApplication.
type MockApplicationMockRecorder struct {
	mock *MockApplication
}

// NewMockApplication creates a new mock instance.
func NewMockApplication(ctrl *gomock.Controller) *MockApplication {
	mock := &MockApplication{ctrl: ctrl}
	mock.recorder = &MockApplicationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplication) EXPECT() *MockApplicationMockRecorder {
	return m.recorder
}

// Commit mocks base method.
func (m *MockApplication) Commit() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit")
	ret0, _ := ret[0].(string)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockApplicationMockRecorder) Commit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockApplication)(nil).Commit))
}

// Configuration mocks base method.
func (m *MockApplication) Configuration() config.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Configuration")
	ret0, _ := ret[0].(config.Config)
	return ret0
}

// Configuration indicates an expected call of Configuration.
func (mr *MockApplicationMockRecorder) Configuration() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Configuration", reflect.TypeOf((*MockApplication)(nil).Configuration))
}

// Context mocks base method.
func (m *MockApplication) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockApplicationMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockApplication)(nil).Context))
}

// Init mocks base method.
func (m *MockApplication) Init() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Init")
}

// Init indicates an expected call of Init.
func (mr *MockApplicationMockRecorder) Init() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockApplication)(nil).Init))
}

// IsDebug mocks base method.
func (m *MockApplication) IsDebug() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDebug")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsDebug indicates an expected call of IsDebug.
func (mr *MockApplicationMockRecorder) IsDebug() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDebug", reflect.TypeOf((*MockApplication)(nil).IsDebug))
}

// IsRunning mocks base method.
func (m *MockApplication) IsRunning() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRunning")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsRunning indicates an expected call of IsRunning.
func (mr *MockApplicationMockRecorder) IsRunning() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRunning", reflect.TypeOf((*MockApplication)(nil).IsRunning))
}

// Logger mocks base method.
func (m *MockApplication) Logger() logger.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logger")
	ret0, _ := ret[0].(logger.Logger)
	return ret0
}

// Logger indicates an expected call of Logger.
func (mr *MockApplicationMockRecorder) Logger() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logger", reflect.TypeOf((*MockApplication)(nil).Logger))
}

// Name mocks base method.
func (m *MockApplication) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockApplicationMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockApplication)(nil).Name))
}

// ProcessManager mocks base method.
func (m *MockApplication) ProcessManager() process.Manager {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessManager")
	ret0, _ := ret[0].(process.Manager)
	return ret0
}

// ProcessManager indicates an expected call of ProcessManager.
func (mr *MockApplicationMockRecorder) ProcessManager() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessManager", reflect.TypeOf((*MockApplication)(nil).ProcessManager))
}

// Run mocks base method.
func (m *MockApplication) Run() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run")
}

// Run indicates an expected call of Run.
func (mr *MockApplicationMockRecorder) Run() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockApplication)(nil).Run))
}

// SetMainLoop mocks base method.
func (m *MockApplication) SetMainLoop(arg0 app.MainLoopFn) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMainLoop", arg0)
}

// SetMainLoop indicates an expected call of SetMainLoop.
func (mr *MockApplicationMockRecorder) SetMainLoop(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMainLoop", reflect.TypeOf((*MockApplication)(nil).SetMainLoop), arg0)
}

// SetOnCHLD mocks base method.
func (m *MockApplication) SetOnCHLD(arg0 app.OnSignalFn) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetOnCHLD", arg0)
}

// SetOnCHLD indicates an expected call of SetOnCHLD.
func (mr *MockApplicationMockRecorder) SetOnCHLD(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOnCHLD", reflect.TypeOf((*MockApplication)(nil).SetOnCHLD), arg0)
}

// SetOnExit mocks base method.
func (m *MockApplication) SetOnExit(arg0 app.OnSignalFn) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetOnExit", arg0)
}

// SetOnExit indicates an expected call of SetOnExit.
func (mr *MockApplicationMockRecorder) SetOnExit(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOnExit", reflect.TypeOf((*MockApplication)(nil).SetOnExit), arg0)
}

// SetOnHUP mocks base method.
func (m *MockApplication) SetOnHUP(arg0 app.OnSignalFn) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetOnHUP", arg0)
}

// SetOnHUP indicates an expected call of SetOnHUP.
func (mr *MockApplicationMockRecorder) SetOnHUP(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOnHUP", reflect.TypeOf((*MockApplication)(nil).SetOnHUP), arg0)
}

// SetOnStart mocks base method.
func (m *MockApplication) SetOnStart(arg0 app.OnSignalFn) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetOnStart", arg0)
}

// SetOnStart indicates an expected call of SetOnStart.
func (mr *MockApplicationMockRecorder) SetOnStart(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOnStart", reflect.TypeOf((*MockApplication)(nil).SetOnStart), arg0)
}

// SetOnUSR1 mocks base method.
func (m *MockApplication) SetOnUSR1(arg0 app.OnSignalFn) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetOnUSR1", arg0)
}

// SetOnUSR1 indicates an expected call of SetOnUSR1.
func (mr *MockApplicationMockRecorder) SetOnUSR1(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOnUSR1", reflect.TypeOf((*MockApplication)(nil).SetOnUSR1), arg0)
}

// SetOnUSR2 mocks base method.
func (m *MockApplication) SetOnUSR2(arg0 app.OnSignalFn) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetOnUSR2", arg0)
}

// SetOnUSR2 indicates an expected call of SetOnUSR2.
func (mr *MockApplicationMockRecorder) SetOnUSR2(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOnUSR2", reflect.TypeOf((*MockApplication)(nil).SetOnUSR2), arg0)
}

// SetOnWINCH mocks base method.
func (m *MockApplication) SetOnWINCH(arg0 app.OnSignalFn) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetOnWINCH", arg0)
}

// SetOnWINCH indicates an expected call of SetOnWINCH.
func (mr *MockApplicationMockRecorder) SetOnWINCH(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOnWINCH", reflect.TypeOf((*MockApplication)(nil).SetOnWINCH), arg0)
}

// Terminate mocks base method.
func (m *MockApplication) Terminate() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Terminate")
}

// Terminate indicates an expected call of Terminate.
func (mr *MockApplicationMockRecorder) Terminate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Terminate", reflect.TypeOf((*MockApplication)(nil).Terminate))
}

// Version mocks base method.
func (m *MockApplication) Version() *semver.SemVer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Version")
	ret0, _ := ret[0].(*semver.SemVer)
	return ret0
}

// Version indicates an expected call of Version.
func (mr *MockApplicationMockRecorder) Version() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Version", reflect.TypeOf((*MockApplication)(nil).Version))
}