// Code generated by MockGen. DO NOT EDIT.
// Source: ./timedcache/cache.go
//
// Generated by this command:
//
//	mockgen -package=timedcache -source=./timedcache/cache.go -destination=mocks/timedcache/cache_mock.go
//

// Package timedcache is a generated GoMock package.
package timedcache

import (
	reflect "reflect"
	time "time"

	timedcache "github.com/Asmodai/gohacks/timedcache"
	gomock "go.uber.org/mock/gomock"
)

// MockTimedCache is a mock of TimedCache interface.
type MockTimedCache struct {
	ctrl     *gomock.Controller
	recorder *MockTimedCacheMockRecorder
}

// MockTimedCacheMockRecorder is the mock recorder for MockTimedCache.
type MockTimedCacheMockRecorder struct {
	mock *MockTimedCache
}

// NewMockTimedCache creates a new mock instance.
func NewMockTimedCache(ctrl *gomock.Controller) *MockTimedCache {
	mock := &MockTimedCache{ctrl: ctrl}
	mock.recorder = &MockTimedCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTimedCache) EXPECT() *MockTimedCacheMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockTimedCache) Add(arg0, arg1 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockTimedCacheMockRecorder) Add(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockTimedCache)(nil).Add), arg0, arg1)
}

// Count mocks base method.
func (m *MockTimedCache) Count() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count")
	ret0, _ := ret[0].(int)
	return ret0
}

// Count indicates an expected call of Count.
func (mr *MockTimedCacheMockRecorder) Count() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockTimedCache)(nil).Count))
}

// Delete mocks base method.
func (m *MockTimedCache) Delete(arg0 any) (any, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockTimedCacheMockRecorder) Delete(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTimedCache)(nil).Delete), arg0)
}

// Expired mocks base method.
func (m *MockTimedCache) Expired() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Expired")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Expired indicates an expected call of Expired.
func (mr *MockTimedCacheMockRecorder) Expired() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Expired", reflect.TypeOf((*MockTimedCache)(nil).Expired))
}

// Flush mocks base method.
func (m *MockTimedCache) Flush() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Flush")
}

// Flush indicates an expected call of Flush.
func (mr *MockTimedCacheMockRecorder) Flush() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flush", reflect.TypeOf((*MockTimedCache)(nil).Flush))
}

// Get mocks base method.
func (m *MockTimedCache) Get(arg0 any) (any, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockTimedCacheMockRecorder) Get(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTimedCache)(nil).Get), arg0)
}

// LastUpdated mocks base method.
func (m *MockTimedCache) LastUpdated() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastUpdated")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// LastUpdated indicates an expected call of LastUpdated.
func (mr *MockTimedCacheMockRecorder) LastUpdated() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastUpdated", reflect.TypeOf((*MockTimedCache)(nil).LastUpdated))
}

// OnEvicted mocks base method.
func (m *MockTimedCache) OnEvicted(arg0 timedcache.OnEvictFn) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnEvicted", arg0)
}

// OnEvicted indicates an expected call of OnEvicted.
func (mr *MockTimedCacheMockRecorder) OnEvicted(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnEvicted", reflect.TypeOf((*MockTimedCache)(nil).OnEvicted), arg0)
}

// Replace mocks base method.
func (m *MockTimedCache) Replace(arg0, arg1 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Replace", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Replace indicates an expected call of Replace.
func (mr *MockTimedCacheMockRecorder) Replace(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Replace", reflect.TypeOf((*MockTimedCache)(nil).Replace), arg0, arg1)
}

// Set mocks base method.
func (m *MockTimedCache) Set(arg0, arg1 any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", arg0, arg1)
}

// Set indicates an expected call of Set.
func (mr *MockTimedCacheMockRecorder) Set(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockTimedCache)(nil).Set), arg0, arg1)
}
