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
	"github.com/kagucho/tsubonesystem3/backend/scope"
	"log"
	"strings"
)

// GetScope returns the scope of the member identified with the given
// credentials.
func (db DB) GetScope(user string, password string) (scope.Scope, error) {
	var result scope.Scope
	var id uint16

	if passwordError := func() error {
		rows, queryError := db.stmts[stmtSelectMemberIDPassword].Query(user)
		if queryError != nil {
			return queryError
		}

		defer func() {
			if closeError := rows.Close(); closeError != nil {
				log.Print(closeError)
			}
		}()

		rows.Next()

		var dbPassword sql.RawBytes
		if scanError := rows.Scan(&id, &dbPassword); scanError != nil {
			return IncorrectIdentity
		}

		return verifyPassword(password, dbPassword)
	}(); passwordError != nil {
		return result, passwordError
	}

	result = result.Set(scope.User).Set(scope.Member)

	rows, queryError := db.stmts[stmtSelectOfficerScopeByInternalMember].Query(id)
	if queryError != nil {
		return result, queryError
	}

	for rows.Next() {
		var dbScope string

		if scanError := rows.Scan(&dbScope); scanError != nil {
			return result, scanError
		}

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
