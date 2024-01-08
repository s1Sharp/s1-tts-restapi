// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/s1Sharp/s1-tts-restapi/service (interfaces: UserService)

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/s1Sharp/s1-tts-restapi/internal/models"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// FindUserByEmail mocks base method.
func (m *MockUserService) FindUserByEmail(arg0 string) (*models.DBUserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByEmail", arg0)
	ret0, _ := ret[0].(*models.DBUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByEmail indicates an expected call of FindUserByEmail.
func (mr *MockUserServiceMockRecorder) FindUserByEmail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByEmail", reflect.TypeOf((*MockUserService)(nil).FindUserByEmail), arg0)
}

// FindUserById mocks base method.
func (m *MockUserService) FindUserById(arg0 string) (*models.DBUserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserById", arg0)
	ret0, _ := ret[0].(*models.DBUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserById indicates an expected call of FindUserById.
func (mr *MockUserServiceMockRecorder) FindUserById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserById", reflect.TypeOf((*MockUserService)(nil).FindUserById), arg0)
}

// RemoveUserById mocks base method.
func (m *MockUserService) RemoveUserById(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUserById", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUserById indicates an expected call of RemoveUserById.
func (mr *MockUserServiceMockRecorder) RemoveUserById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUserById", reflect.TypeOf((*MockUserService)(nil).RemoveUserById), arg0)
}

// UpdateUserById mocks base method.
func (m *MockUserService) UpdateUserById(arg0 string, arg1 *models.UpdateInput) (*models.DBUserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserById", arg0, arg1)
	ret0, _ := ret[0].(*models.DBUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserById indicates an expected call of UpdateUserById.
func (mr *MockUserServiceMockRecorder) UpdateUserById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserById", reflect.TypeOf((*MockUserService)(nil).UpdateUserById), arg0, arg1)
}
