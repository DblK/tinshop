// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DblK/tinshop/repository (interfaces: Stats)

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	repository "github.com/DblK/tinshop/repository"
	gomock "github.com/golang/mock/gomock"
)

// MockStats is a mock of Stats interface.
type MockStats struct {
	ctrl     *gomock.Controller
	recorder *MockStatsMockRecorder
}

// MockStatsMockRecorder is the mock recorder for MockStats.
type MockStatsMockRecorder struct {
	mock *MockStats
}

// NewMockStats creates a new mock instance.
func NewMockStats(ctrl *gomock.Controller) *MockStats {
	mock := &MockStats{ctrl: ctrl}
	mock.recorder = &MockStatsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStats) EXPECT() *MockStatsMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStats) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStatsMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStats)(nil).Close))
}

// ListVisit mocks base method.
func (m *MockStats) ListVisit(arg0 *repository.Switch) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ListVisit", arg0)
}

// ListVisit indicates an expected call of ListVisit.
func (mr *MockStatsMockRecorder) ListVisit(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListVisit", reflect.TypeOf((*MockStats)(nil).ListVisit), arg0)
}

// Load mocks base method.
func (m *MockStats) Load() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Load")
}

// Load indicates an expected call of Load.
func (mr *MockStatsMockRecorder) Load() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockStats)(nil).Load))
}
