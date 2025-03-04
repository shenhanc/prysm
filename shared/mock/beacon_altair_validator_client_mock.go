// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/prysmaticlabs/prysm/proto/prysm/v2 (interfaces: BeaconNodeValidatorAltair_StreamBlocksClient)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v2 "github.com/prysmaticlabs/prysm/proto/prysm/v2"
	metadata "google.golang.org/grpc/metadata"
)

// MockBeaconNodeValidatorAltair_StreamBlocksClient is a mock of BeaconNodeValidatorAltair_StreamBlocksClient interface
type MockBeaconNodeValidatorAltair_StreamBlocksClient struct {
	ctrl     *gomock.Controller
	recorder *MockBeaconNodeValidatorAltair_StreamBlocksClientMockRecorder
}

// MockBeaconNodeValidatorAltair_StreamBlocksClientMockRecorder is the mock recorder for MockBeaconNodeValidatorAltair_StreamBlocksClient
type MockBeaconNodeValidatorAltair_StreamBlocksClientMockRecorder struct {
	mock *MockBeaconNodeValidatorAltair_StreamBlocksClient
}

// NewMockBeaconNodeValidatorAltair_StreamBlocksClient creates a new mock instance
func NewMockBeaconNodeValidatorAltair_StreamBlocksClient(ctrl *gomock.Controller) *MockBeaconNodeValidatorAltair_StreamBlocksClient {
	mock := &MockBeaconNodeValidatorAltair_StreamBlocksClient{ctrl: ctrl}
	mock.recorder = &MockBeaconNodeValidatorAltair_StreamBlocksClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBeaconNodeValidatorAltair_StreamBlocksClient) EXPECT() *MockBeaconNodeValidatorAltair_StreamBlocksClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method
func (m *MockBeaconNodeValidatorAltair_StreamBlocksClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend
func (mr *MockBeaconNodeValidatorAltair_StreamBlocksClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockBeaconNodeValidatorAltair_StreamBlocksClient)(nil).CloseSend))
}

// Context mocks base method
func (m *MockBeaconNodeValidatorAltair_StreamBlocksClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockBeaconNodeValidatorAltair_StreamBlocksClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockBeaconNodeValidatorAltair_StreamBlocksClient)(nil).Context))
}

// Header mocks base method
func (m *MockBeaconNodeValidatorAltair_StreamBlocksClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header
func (mr *MockBeaconNodeValidatorAltair_StreamBlocksClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockBeaconNodeValidatorAltair_StreamBlocksClient)(nil).Header))
}

// Recv mocks base method
func (m *MockBeaconNodeValidatorAltair_StreamBlocksClient) Recv() (*v2.StreamBlocksResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*v2.StreamBlocksResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv
func (mr *MockBeaconNodeValidatorAltair_StreamBlocksClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockBeaconNodeValidatorAltair_StreamBlocksClient)(nil).Recv))
}

// RecvMsg mocks base method
func (m *MockBeaconNodeValidatorAltair_StreamBlocksClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg
func (mr *MockBeaconNodeValidatorAltair_StreamBlocksClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockBeaconNodeValidatorAltair_StreamBlocksClient)(nil).RecvMsg), arg0)
}

// SendMsg mocks base method
func (m *MockBeaconNodeValidatorAltair_StreamBlocksClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockBeaconNodeValidatorAltair_StreamBlocksClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockBeaconNodeValidatorAltair_StreamBlocksClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method
func (m *MockBeaconNodeValidatorAltair_StreamBlocksClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer
func (mr *MockBeaconNodeValidatorAltair_StreamBlocksClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockBeaconNodeValidatorAltair_StreamBlocksClient)(nil).Trailer))
}
