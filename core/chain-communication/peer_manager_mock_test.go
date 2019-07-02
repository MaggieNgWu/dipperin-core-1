// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dipperin/dipperin-core/core/chain-communication (interfaces: PeerManager)

// Package chain_communication is a generated GoMock package.
package chain_communication

import (
	//chain_communication "github.com/dipperin/dipperin-core/core/chain-communication"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockPeerManager is a mock of PeerManager interface
type MockPeerManager struct {
	ctrl     *gomock.Controller
	recorder *MockPeerManagerMockRecorder
}

// MockPeerManagerMockRecorder is the mock recorder for MockPeerManager
type MockPeerManagerMockRecorder struct {
	mock *MockPeerManager
}

// NewMockPeerManager creates a new mock instance
func NewMockPeerManager(ctrl *gomock.Controller) *MockPeerManager {
	mock := &MockPeerManager{ctrl: ctrl}
	mock.recorder = &MockPeerManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPeerManager) EXPECT() *MockPeerManagerMockRecorder {
	return m.recorder
}

// BestPeer mocks base method
func (m *MockPeerManager) BestPeer() PmAbstractPeer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BestPeer")
	ret0, _ := ret[0].(PmAbstractPeer)
	return ret0
}

// BestPeer indicates an expected call of BestPeer
func (mr *MockPeerManagerMockRecorder) BestPeer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BestPeer", reflect.TypeOf((*MockPeerManager)(nil).BestPeer))
}

// GetPeer mocks base method
func (m *MockPeerManager) GetPeer(arg0 string) PmAbstractPeer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeer", arg0)
	ret0, _ := ret[0].(PmAbstractPeer)
	return ret0
}

// GetPeer indicates an expected call of GetPeer
func (mr *MockPeerManagerMockRecorder) GetPeer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeer", reflect.TypeOf((*MockPeerManager)(nil).GetPeer), arg0)
}

// GetPeers mocks base method
func (m *MockPeerManager) GetPeers() map[string]PmAbstractPeer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeers")
	ret0, _ := ret[0].(map[string]PmAbstractPeer)
	return ret0
}

// GetPeers indicates an expected call of GetPeers
func (mr *MockPeerManagerMockRecorder) GetPeers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeers", reflect.TypeOf((*MockPeerManager)(nil).GetPeers))
}

// IsSync mocks base method
func (m *MockPeerManager) IsSync() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSync")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsSync indicates an expected call of IsSync
func (mr *MockPeerManagerMockRecorder) IsSync() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSync", reflect.TypeOf((*MockPeerManager)(nil).IsSync))
}

// RemovePeer mocks base method
func (m *MockPeerManager) RemovePeer(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemovePeer", arg0)
}

// RemovePeer indicates an expected call of RemovePeer
func (mr *MockPeerManagerMockRecorder) RemovePeer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemovePeer", reflect.TypeOf((*MockPeerManager)(nil).RemovePeer), arg0)
}
