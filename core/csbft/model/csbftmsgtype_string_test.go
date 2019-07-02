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

// Code generated by "stringer -type=CsBftMsgType"; DO NOT EDIT.

package model

import "testing"

func TestCsBftMsgType_String(t *testing.T) {
	tests := []struct {
		name string
		i    CsBftMsgType
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.String(); got != tt.want {
				t.Errorf("CsBftMsgType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCsBftMsgType_String2(t *testing.T) {
	n := CsBftMsgType(1)
	n.String()
	n = CsBftMsgType(11)
	n.String()
	n = CsBftMsgType(55)
	n.String()
}
