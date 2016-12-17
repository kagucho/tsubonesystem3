/*
  Copyright (C) 2016  Kagucho <kagucho.net@gmail.com>

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU Affero General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU Affero General Public License for more details.

  You should have received a copy of the GNU Affero General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

// Package scope implements the common expression of scope.
package scope

// The index of the flags for the scopes.
const (
  Basic uint = iota
  Privacy
  Management
)

// Scope is the common expression of scope
type Scope struct {
  expression uint
}

// IsSet returns whether the flag identified with the given index is set.
func (scope Scope) IsSet(index uint) bool {
  return scope.expression & (1 << index) != 0
}

// IsSetAny returns whether any flag is set.
func (scope Scope) IsSetAny() bool {
  return scope.expression != 0
}

/*
  Set returns scope.Scope which has the existing flags and the flag identified
  with the given index.
 */
func (scope Scope) Set(index uint) Scope {
  return Scope{scope.expression | (1 << index)}
}
