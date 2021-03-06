// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mdanialr/sns_backend/internal/database/sql (interfaces: SNS)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
)

// MockSNS is a mock of SNS interface.
type MockSNS struct {
	ctrl     *gomock.Controller
	recorder *MockSNSMockRecorder
}

// MockSNSMockRecorder is the mock recorder for MockSNS.
type MockSNSMockRecorder struct {
	mock *MockSNS
}

// NewMockSNS creates a new mock instance.
func NewMockSNS(ctrl *gomock.Controller) *MockSNS {
	mock := &MockSNS{ctrl: ctrl}
	mock.recorder = &MockSNSMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSNS) EXPECT() *MockSNSMockRecorder {
	return m.recorder
}

// CreateSend mocks base method.
func (m *MockSNS) CreateSend(arg0 context.Context, arg1 database.CreateSendParams) (database.Send, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSend", arg0, arg1)
	ret0, _ := ret[0].(database.Send)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSend indicates an expected call of CreateSend.
func (mr *MockSNSMockRecorder) CreateSend(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSend", reflect.TypeOf((*MockSNS)(nil).CreateSend), arg0, arg1)
}

// CreateShorten mocks base method.
func (m *MockSNS) CreateShorten(arg0 context.Context, arg1 database.CreateShortenParams) (database.Shorten, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateShorten", arg0, arg1)
	ret0, _ := ret[0].(database.Shorten)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateShorten indicates an expected call of CreateShorten.
func (mr *MockSNSMockRecorder) CreateShorten(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateShorten", reflect.TypeOf((*MockSNS)(nil).CreateShorten), arg0, arg1)
}

// DeleteSend mocks base method.
func (m *MockSNS) DeleteSend(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSend", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSend indicates an expected call of DeleteSend.
func (mr *MockSNSMockRecorder) DeleteSend(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSend", reflect.TypeOf((*MockSNS)(nil).DeleteSend), arg0, arg1)
}

// DeleteShorten mocks base method.
func (m *MockSNS) DeleteShorten(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteShorten", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteShorten indicates an expected call of DeleteShorten.
func (mr *MockSNSMockRecorder) DeleteShorten(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteShorten", reflect.TypeOf((*MockSNS)(nil).DeleteShorten), arg0, arg1)
}

// GetSend mocks base method.
func (m *MockSNS) GetSend(arg0 context.Context, arg1 int64) (database.Send, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSend", arg0, arg1)
	ret0, _ := ret[0].(database.Send)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSend indicates an expected call of GetSend.
func (mr *MockSNSMockRecorder) GetSend(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSend", reflect.TypeOf((*MockSNS)(nil).GetSend), arg0, arg1)
}

// GetSendByUrl mocks base method.
func (m *MockSNS) GetSendByUrl(arg0 context.Context, arg1 string) (database.Send, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSendByUrl", arg0, arg1)
	ret0, _ := ret[0].(database.Send)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSendByUrl indicates an expected call of GetSendByUrl.
func (mr *MockSNSMockRecorder) GetSendByUrl(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSendByUrl", reflect.TypeOf((*MockSNS)(nil).GetSendByUrl), arg0, arg1)
}

// GetShorten mocks base method.
func (m *MockSNS) GetShorten(arg0 context.Context, arg1 int64) (database.Shorten, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShorten", arg0, arg1)
	ret0, _ := ret[0].(database.Shorten)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShorten indicates an expected call of GetShorten.
func (mr *MockSNSMockRecorder) GetShorten(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShorten", reflect.TypeOf((*MockSNS)(nil).GetShorten), arg0, arg1)
}

// GetShortenByUrl mocks base method.
func (m *MockSNS) GetShortenByUrl(arg0 context.Context, arg1 string) (database.Shorten, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShortenByUrl", arg0, arg1)
	ret0, _ := ret[0].(database.Shorten)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShortenByUrl indicates an expected call of GetShortenByUrl.
func (mr *MockSNSMockRecorder) GetShortenByUrl(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShortenByUrl", reflect.TypeOf((*MockSNS)(nil).GetShortenByUrl), arg0, arg1)
}

// ListSend mocks base method.
func (m *MockSNS) ListSend(arg0 context.Context) ([]database.Send, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSend", arg0)
	ret0, _ := ret[0].([]database.Send)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSend indicates an expected call of ListSend.
func (mr *MockSNSMockRecorder) ListSend(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSend", reflect.TypeOf((*MockSNS)(nil).ListSend), arg0)
}

// ListShorten mocks base method.
func (m *MockSNS) ListShorten(arg0 context.Context, arg1 string, arg2 database.DBOrder) ([]database.Shorten, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListShorten", arg0, arg1, arg2)
	ret0, _ := ret[0].([]database.Shorten)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListShorten indicates an expected call of ListShorten.
func (mr *MockSNSMockRecorder) ListShorten(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListShorten", reflect.TypeOf((*MockSNS)(nil).ListShorten), arg0, arg1, arg2)
}

// UpdateSend mocks base method.
func (m *MockSNS) UpdateSend(arg0 context.Context, arg1 database.UpdateSendParams) (database.Send, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSend", arg0, arg1)
	ret0, _ := ret[0].(database.Send)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSend indicates an expected call of UpdateSend.
func (mr *MockSNSMockRecorder) UpdateSend(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSend", reflect.TypeOf((*MockSNS)(nil).UpdateSend), arg0, arg1)
}

// UpdateShorten mocks base method.
func (m *MockSNS) UpdateShorten(arg0 context.Context, arg1 database.UpdateShortenParams) (database.Shorten, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateShorten", arg0, arg1)
	ret0, _ := ret[0].(database.Shorten)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateShorten indicates an expected call of UpdateShorten.
func (mr *MockSNSMockRecorder) UpdateShorten(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateShorten", reflect.TypeOf((*MockSNS)(nil).UpdateShorten), arg0, arg1)
}
