// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ConsenSys/fc-retrieval/common/pkg/request (interfaces: HttpCommunications)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	fcrmessages "github.com/ConsenSys/fc-retrieval/common/pkg/fcrmessages"
	gomock "github.com/golang/mock/gomock"
)

// MockHttpCommunications is a mock of HttpCommunications interface.
type MockHttpCommunications struct {
	ctrl     *gomock.Controller
	recorder *MockHttpCommunicationsMockRecorder
}

// MockHttpCommunicationsMockRecorder is the mock recorder for MockHttpCommunications.
type MockHttpCommunicationsMockRecorder struct {
	mock *MockHttpCommunications
}

// NewMockHttpCommunications creates a new mock instance.
func NewMockHttpCommunications(ctrl *gomock.Controller) *MockHttpCommunications {
	mock := &MockHttpCommunications{ctrl: ctrl}
	mock.recorder = &MockHttpCommunicationsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHttpCommunications) EXPECT() *MockHttpCommunicationsMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockHttpCommunications) Delete(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockHttpCommunicationsMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockHttpCommunications)(nil).Delete), arg0)
}

// GetJSON mocks base method.
func (m *MockHttpCommunications) GetJSON(arg0 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJSON", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJSON indicates an expected call of GetJSON.
func (mr *MockHttpCommunicationsMockRecorder) GetJSON(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJSON", reflect.TypeOf((*MockHttpCommunications)(nil).GetJSON), arg0)
}

// SendJSON mocks base method.
func (m *MockHttpCommunications) SendJSON(arg0 string, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendJSON", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendJSON indicates an expected call of SendJSON.
func (mr *MockHttpCommunicationsMockRecorder) SendJSON(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendJSON", reflect.TypeOf((*MockHttpCommunications)(nil).SendJSON), arg0, arg1)
}

// SendMessage mocks base method.
func (m *MockHttpCommunications) SendMessage(arg0 string, arg1 *fcrmessages.FCRMessage) (*fcrmessages.FCRMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", arg0, arg1)
	ret0, _ := ret[0].(*fcrmessages.FCRMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockHttpCommunicationsMockRecorder) SendMessage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockHttpCommunications)(nil).SendMessage), arg0, arg1)
}
