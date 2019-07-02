// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dipperin/dipperin-core/core/chain-communication (interfaces: TxPool)

// Package chain_communication is a generated GoMock package.
package chain_communication

import (
	common "github.com/dipperin/dipperin-core/common"
	bloom "github.com/dipperin/dipperin-core/core/bloom"
	model "github.com/dipperin/dipperin-core/core/model"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockTxPool is a mock of TxPool interface
type MockTxPool struct {
	ctrl     *gomock.Controller
	recorder *MockTxPoolMockRecorder
}

// MockTxPoolMockRecorder is the mock recorder for MockTxPool
type MockTxPoolMockRecorder struct {
	mock *MockTxPool
}

// NewMockTxPool creates a new mock instance
func NewMockTxPool(ctrl *gomock.Controller) *MockTxPool {
	mock := &MockTxPool{ctrl: ctrl}
	mock.recorder = &MockTxPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTxPool) EXPECT() *MockTxPoolMockRecorder {
	return m.recorder
}

// AddLocal mocks base method
func (m *MockTxPool) AddLocal(arg0 model.AbstractTransaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLocal", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddLocal indicates an expected call of AddLocal
func (mr *MockTxPoolMockRecorder) AddLocal(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLocal", reflect.TypeOf((*MockTxPool)(nil).AddLocal), arg0)
}

// AddLocals mocks base method
func (m *MockTxPool) AddLocals(arg0 []model.AbstractTransaction) []error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLocals", arg0)
	ret0, _ := ret[0].([]error)
	return ret0
}

// AddLocals indicates an expected call of AddLocals
func (mr *MockTxPoolMockRecorder) AddLocals(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLocals", reflect.TypeOf((*MockTxPool)(nil).AddLocals), arg0)
}

// AddRemote mocks base method
func (m *MockTxPool) AddRemote(arg0 model.AbstractTransaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRemote", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRemote indicates an expected call of AddRemote
func (mr *MockTxPoolMockRecorder) AddRemote(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRemote", reflect.TypeOf((*MockTxPool)(nil).AddRemote), arg0)
}

// AddRemotes mocks base method
func (m *MockTxPool) AddRemotes(arg0 []model.AbstractTransaction) []error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRemotes", arg0)
	ret0, _ := ret[0].([]error)
	return ret0
}

// AddRemotes indicates an expected call of AddRemotes
func (mr *MockTxPoolMockRecorder) AddRemotes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRemotes", reflect.TypeOf((*MockTxPool)(nil).AddRemotes), arg0)
}

// ConvertPoolToMap mocks base method
func (m *MockTxPool) ConvertPoolToMap() map[common.Hash]model.AbstractTransaction {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConvertPoolToMap")
	ret0, _ := ret[0].(map[common.Hash]model.AbstractTransaction)
	return ret0
}

// ConvertPoolToMap indicates an expected call of ConvertPoolToMap
func (mr *MockTxPoolMockRecorder) ConvertPoolToMap() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConvertPoolToMap", reflect.TypeOf((*MockTxPool)(nil).ConvertPoolToMap))
}

// GetTxsEstimator mocks base method
func (m *MockTxPool) GetTxsEstimator(arg0 *bloom.Bloom) *bloom.HybridEstimator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxsEstimator", arg0)
	ret0, _ := ret[0].(*bloom.HybridEstimator)
	return ret0
}

// GetTxsEstimator indicates an expected call of GetTxsEstimator
func (mr *MockTxPoolMockRecorder) GetTxsEstimator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxsEstimator", reflect.TypeOf((*MockTxPool)(nil).GetTxsEstimator), arg0)
}

// Pending mocks base method
func (m *MockTxPool) Pending() (map[common.Address][]model.AbstractTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pending")
	ret0, _ := ret[0].(map[common.Address][]model.AbstractTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Pending indicates an expected call of Pending
func (mr *MockTxPoolMockRecorder) Pending() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pending", reflect.TypeOf((*MockTxPool)(nil).Pending))
}

// Queueing mocks base method
func (m *MockTxPool) Queueing() (map[common.Address][]model.AbstractTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Queueing")
	ret0, _ := ret[0].(map[common.Address][]model.AbstractTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Queueing indicates an expected call of Queueing
func (mr *MockTxPoolMockRecorder) Queueing() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Queueing", reflect.TypeOf((*MockTxPool)(nil).Queueing))
}

// Stats mocks base method
func (m *MockTxPool) Stats() (int, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stats")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(int)
	return ret0, ret1
}

// Stats indicates an expected call of Stats
func (mr *MockTxPoolMockRecorder) Stats() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stats", reflect.TypeOf((*MockTxPool)(nil).Stats))
}
