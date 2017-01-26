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

	passwordValid, passwordError := func() (bool, error) {
		hashedPassword, hashError := hashPassword(password)
		if hashError != nil {
			return false, hashError
		}

		rows, queryError := db.stmts[stmtSelectMemberIDPassword].Query(user)
		if queryError != nil {
			return false, queryError
		}

		defer func() {
			if closeError := rows.Close(); closeError != nil {
				log.Print(closeError)
			}
		}()

		rows.Next()

		var dbPassword sql.RawBytes
		if scanError := rows.Scan(&id, &dbPassword); scanError != nil {
			return false, sql.ErrNoRows
		}

		return hmac.Equal(dbPassword, hashedPassword), nil
	}()

	if passwordValid {
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

	return result, passwordError
}

func (db DB) QueryTemporary(id string) (bool, error) {
	rows, queryError := db.stmts[stmtSelectMemberPassword].Query(id)
	if queryError != nil {
		return false, queryError
	}

	defer func() {
		if closeError := rows.Close(); closeError != nil {
			log.Print(closeError)
		}
	}()

	var dbPassword sql.RawBytes
	if scanError := rows.Scan(&dbPassword); scanError != nil {
		return false, scanError
	}

	for _, value := range dbPassword {
		if value != 0 {
			return false, nil
		}
	}

	return true, nil
}
