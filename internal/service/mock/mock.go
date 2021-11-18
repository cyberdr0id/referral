// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	repository "github.com/cyberdr0id/referral/internal/repository"
	gomock "github.com/golang/mock/gomock"
)

// MockAuth is a mock of Auth interface.
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth.
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance.
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// LogIn mocks base method.
func (m *MockAuth) LogIn(name, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogIn", name, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LogIn indicates an expected call of LogIn.
func (mr *MockAuthMockRecorder) LogIn(name, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogIn", reflect.TypeOf((*MockAuth)(nil).LogIn), name, password)
}

// SignUp mocks base method.
func (m *MockAuth) SignUp(name, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", name, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockAuthMockRecorder) SignUp(name, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockAuth)(nil).SignUp), name, password)
}

// MockReferral is a mock of Referral interface.
type MockReferral struct {
	ctrl     *gomock.Controller
	recorder *MockReferralMockRecorder
}

// MockReferralMockRecorder is the mock recorder for MockReferral.
type MockReferralMockRecorder struct {
	mock *MockReferral
}

// NewMockReferral creates a new mock instance.
func NewMockReferral(ctrl *gomock.Controller) *MockReferral {
	mock := &MockReferral{ctrl: ctrl}
	mock.recorder = &MockReferralMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReferral) EXPECT() *MockReferralMockRecorder {
	return m.recorder
}

// AddCandidate mocks base method.
func (m *MockReferral) AddCandidate(name, surname string, fileID int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCandidate", name, surname, fileID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCandidate indicates an expected call of AddCandidate.
func (mr *MockReferralMockRecorder) AddCandidate(name, surname, fileID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCandidate", reflect.TypeOf((*MockReferral)(nil).AddCandidate), name, surname, fileID)
}

// GetCVID mocks base method.
func (m *MockReferral) GetCVID(id int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCVID", id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCVID indicates an expected call of GetCVID.
func (mr *MockReferralMockRecorder) GetCVID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCVID", reflect.TypeOf((*MockReferral)(nil).GetCVID), id)
}

// GetRequests mocks base method.
func (m *MockReferral) GetRequests(id int) ([]repository.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRequests", id)
	ret0, _ := ret[0].([]repository.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRequests indicates an expected call of GetRequests.
func (mr *MockReferralMockRecorder) GetRequests(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRequests", reflect.TypeOf((*MockReferral)(nil).GetRequests), id)
}

// UpdateRequest mocks base method.
func (m *MockReferral) UpdateRequest(id, status string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRequest", id, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRequest indicates an expected call of UpdateRequest.
func (mr *MockReferralMockRecorder) UpdateRequest(id, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRequest", reflect.TypeOf((*MockReferral)(nil).UpdateRequest), id, status)
}