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

package db

import (
	"database/sql"
	"github.com/kagucho/tsubonesystem3/scope"
	"testing"
)

func (db DB) testGetScope(t *testing.T) {
	for _, test := range [...]struct {
		description string
		user        string
		password    string
		scope       scope.Scope
		scopeError  error
	}{
		{
			`president`, `1stDisplayID`, `1stPassword`,
			scope.Scope{}.Set(
				scope.Management).Set(
				scope.Privacy).Set(
				scope.Basic),
			nil,
		}, {
			`invalidUser`, ``, `1stPassword`,
			scope.Scope{}, sql.ErrNoRows,
		}, {
			`invalidPassword`, `1stDisplayID`, ``,
			scope.Scope{}, nil,
		},
	} {
		test := test

		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			result, queryError :=
				db.GetScope(test.user, test.password)
			if queryError != test.scopeError {
				t.Errorf(`expected %v, got %v`,
					test.scopeError, queryError)
			}

			if result != test.scope {
				t.Errorf(`expected %v, got %v`,
					test.scope, result)
			}
		})
	}
}
