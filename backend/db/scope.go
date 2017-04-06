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
	"crypto/md5"
	"crypto/subtle"
	"database/sql"
	"github.com/kagucho/tsubonesystem3/backend/scope"
	"log"
	"strings"
)

func (db DB) Authenticate(id, password string) error {
	rows, queryErr := db.stmts[stmtSelectMemberPasswordByID].Query(id)
	if queryErr != nil {
		return queryErr
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Print(closeErr)
		}
	}()

	rows.Next()

	var dbPassword sql.RawBytes
	if scanErr := rows.Scan(&dbPassword); scanErr != nil {
		return ErrIncorrectIdentity
	}

	return verifyPassword(password, dbPassword)
}

/*
GetScope returns the scope of the member identified with the given
credentials.

It returns db.ErrIncorrectIdentity if the credentials are incorrect. Other
errors tell db.DB is bad.
*/
func (db DB) GetScope(id string, password string) (scope.Scope, error) {
	var result scope.Scope
	var dbID uint16

	if passwordErr := func() error {
		rows, queryErr := db.stmts[stmtSelectMemberInternalIDPasswordByID].Query(id)
		if queryErr != nil {
			return queryErr
		}

		defer func() {
			if closeErr := rows.Close(); closeErr != nil {
				log.Print(closeErr)
			}
		}()

		rows.Next()

		var dbPassword sql.RawBytes
		if scanErr := rows.Scan(&dbID, &dbPassword); scanErr != nil {
			return ErrIncorrectIdentity
		}

		if verifyPassword(password, dbPassword) == nil {
			return nil
		}

		// TODO: remove MD5 support
		hashed := md5.Sum([]byte(password))
		if subtle.ConstantTimeCompare(hashed[:], dbPassword[:md5.Size]) == 0 {
			return ErrIncorrectIdentity
		}

		newPassword, newPasswordErr := makeDBPassword(password)
		if newPasswordErr != nil {
			return newPasswordErr
		}

		result, execErr := db.stmts[stmtUpdateMemberPassword].Exec(newPassword, id)
		if execErr != nil {
			return execErr
		}

		affected, affectedErr := result.RowsAffected()
		if affectedErr != nil {
			return affectedErr
		}

		if affected <= 0 {
			return ErrIncorrectIdentity
		}

		return nil
	}(); passwordErr != nil {
		return result, passwordErr
	}

	result = result.Set(scope.User).Set(scope.Member)

	rows, queryErr := db.stmts[stmtSelectOfficerScopeByInternalMember].Query(dbID)
	if queryErr != nil {
		return result, queryErr
	}

	for rows.Next() {
		var dbScope string

		if scanErr := rows.Scan(&dbScope); scanErr != nil {
			return result, scanErr
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
