/*
	Copyright (C) 2017  Kagucho <kagucho.net@gmail.com>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published
	by the Free Software Foundation, either version 3 of the License, or (at
	your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package configuration

import (
	"crypto/sha256"
	"testing"
)

func TestDBDSN(t *testing.T) {
	if DBDSN == "" {
		t.Fail()
	}
}

func TestDBPasswordKey(t *testing.T) {
	if result := len(DBPasswordKey); result != sha256.Size224 {
		t.Errorf(`invalid length; expected %v, got %v`,
			sha256.Size224, result)
	}

	if DBPasswordKey == `` {
		t.Fail()
	}
}

func TestListenNet(t *testing.T) {
	if ListenNet == "" {
		t.Fail()
	}
}

func TestListenAddress(t *testing.T) {
	if ListenAddress == "" {
		t.Fail()
	}
}
