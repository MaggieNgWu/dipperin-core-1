// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/caiqingfeng/dipperin-core/core/cs-chain (interfaces: CacheDB)

// Package cs_chain is a generated GoMock package.
package cs_chain

import (
	common "github.com/dipperin/dipperin-core/common"
	model "github.com/dipperin/dipperin-core/core/model"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCacheDB is a mock of CacheDB interface
type MockCacheDB struct {
	ctrl     *gomock.Controller
	recorder *MockCacheDBMockRecorder
}

// MockCacheDBMockRecorder is the mock recorder for MockCacheDB
type MockCacheDBMockRecorder struct {
	mock *MockCacheDB
}

// NewMockCacheDB creates a new mock instance
func NewMockCacheDB(ctrl *gomock.Controller) *MockCacheDB {
	mock := &MockCacheDB{ctrl: ctrl}
	mock.recorder = &MockCacheDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCacheDB) EXPECT() *MockCacheDBMockRecorder {
	return m.recorder
}

// GetSeenCommits mocks base method
func (m *MockCacheDB) GetSeenCommits(arg0 uint64, arg1 common.Hash) ([]model.AbstractVerification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSeenCommits", arg0, arg1)
	ret0, _ := ret[0].([]model.AbstractVerification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSeenCommits indicates an expected call of GetSeenCommits
func (mr *MockCacheDBMockRecorder) GetSeenCommits(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSeenCommits", reflect.TypeOf((*MockCacheDB)(nil).GetSeenCommits), arg0, arg1)
}

// SaveSeenCommits mocks base method
func (m *MockCacheDB) SaveSeenCommits(arg0 uint64, arg1 common.Hash, arg2 []model.AbstractVerification) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveSeenCommits", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveSeenCommits indicates an expected call of SaveSeenCommits
func (mr *MockCacheDBMockRecorder) SaveSeenCommits(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSeenCommits", reflect.TypeOf((*MockCacheDB)(nil).SaveSeenCommits), arg0, arg1, arg2)
}
