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
	"crypto/hmac"
	"crypto/sha256"
	"github.com/kagucho/tsubonesystem3/configuration"
	"github.com/kagucho/tsubonesystem3/scope"
	"strings"
)

// GetScope returns the scope of the member identified with the given
// credentials.
func (db DB) GetScope(user string, password string) (scope.Scope, error) {
	var result scope.Scope
	var dbPassword string
	var id uint16

	if scanError := db.stmts[stmtSelectMemberIDPassword].QueryRow(user).Scan(&id, &dbPassword); scanError != nil {
		return scope.Scope{}, scanError
	}

	hash := hmac.New(sha256.New224, []byte(configuration.DBPasswordKey))
	hash.Write([]byte(password))
	if hmac.Equal([]byte(dbPassword), hash.Sum(nil)) {
		var dbScope string
		if scanError := db.stmts[stmtSelectOfficerScopeInternal].QueryRow(id).Scan(&dbScope); scanError != nil {
			return scope.Scope{}, scanError
		}

		result = result.Set(scope.Basic)
		for _, flag := range strings.Split(dbScope, `,`) {
			switch flag {
			case `management`:
				result = result.Set(scope.Management)

			case `privacy`:
				result = result.Set(scope.Privacy)
			}
		}
	}

	return result, nil
}

func (db DB) QueryTemporary(id string) (bool, error) {
	var dbPassword string

	if scanError := db.stmts[stmtSelectMemberPassword].QueryRow(id).Scan(&dbPassword); scanError != nil {
		return false, scanError
	}

	for index := 0; index < len(dbPassword); index++ {
		if dbPassword[index] != 0 {
			return false, nil
		}
	}

	return true, nil
}
