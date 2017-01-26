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

package private

import "testing"

func TestGetRoute(t *testing.T) {
	t.Parallel()

	for key := range routes {
		if GetRoute(key).Handle == nil {
			t.Errorf(`invalid %v; got (nil) as Handle`, key)
		}
	}

	if result := GetRoute(`/invalid`).Handle; result != nil {
		t.Error(`invalid /invalid; expected (nil) as Handle, got `, result)
	}
}
