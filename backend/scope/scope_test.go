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

package scope

import "testing"

func TestIsSet(t *testing.T) {
	t.Parallel()

	t.Run(`true`, func(t *testing.T) {
		t.Parallel()

		if (!Scope{3}.IsSet(1)) {
			t.Fail()
		}
	})

	t.Run(`false`, func(t *testing.T) {
		t.Parallel()

		if (Scope{1}.IsSet(1)) {
			t.Fail()
		}
	})
}

func TestIsSetAny(t *testing.T) {
	t.Parallel()

	t.Run(`true`, func(t *testing.T) {
		t.Parallel()

		if (!Scope{1}.IsSetAny()) {
			t.Fail()
		}
	})

	t.Run(`false`, func(t *testing.T) {
		t.Parallel()

		if (Scope{0}.IsSetAny()) {
			t.Fail()
		}
	})
}

func TestSet(t *testing.T) {
	t.Parallel()

	expected := Scope{1}
	result := Scope{}.Set(0)
	if result != expected {
		t.Error(`expected `, expected, `, got `, result)
	}
}
