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

/*
Package db implements an abstraction for the database.

String List

Some functions have a space-delimited (%x20) list as an argument. It is
expected to be compatible with similar lists used in RFC 6749.

RFC 6749 - The OAuth 2.0 Authorization Framework
https://tools.ietf.org/html/rfc6749
*/
package db

import (
	"database/sql"
	"errors"
	"github.com/kagucho/tsubonesystem3/configuration"
	"log"
)

/*
DB is a structure holding the connection to the database. It should be
initialized with db.New.
*/
type DB struct {
	sql   *sql.DB
	stmts [stmtNumber]*sql.Stmt
}

// ErrDupEntry is an error telling the entry is duplicate.
var ErrDupEntry = errors.New(`duplicate entry`)

// ErrIncorrectIdentity is an error telling the identity is incorrect.
var ErrIncorrectIdentity = errors.New(`incorrect identity`)

// ErrInvalid is an error telling some of the given parameter is invalid.
var ErrInvalid = errors.New(`invalid`)

// ErrBadOmision is an error telling omitting parameters is forbidden.
var ErrBadOmission = errors.New(`bad omission`)

/*
	MariaDB Error Codes - MariaDB Knowledge Base
	https://mariadb.com/kb/en/mariadb/mariadb-error-codes/
*/
const (
	erDupEntry                    = 1062
	erNoReferencedRow             = 1216
	erRowIsReferenced             = 1217
	erTruncatedWrongValueForField = 1366
	erDataTooLong                 = 1406
	erRowIsReferenced2            = 1451
	erNoReferencedRow2            = 1452
	erWrongValue                  = 1525
	erSignalException             = 1644
)

// New returns a new db.DB. Resources will be holded until Close gets called.
func New() (DB, error) {
	var db DB
	var err error

	db.sql, err = sql.Open(`mysql`, configuration.DBDSN)
	if err != nil {
		return db, err
	}

	db.sql.SetMaxOpenConns(128)

	_, err = db.sql.Exec(`SET sql_mode='ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ZERO_DATE,NO_ZERO_IN_DATE,STRICT_ALL_TABLES'`)
	if err != nil {
		if err := db.Close(); err != nil {
			log.Print(err)
		}

		return db, err
	}

	err = db.prepareStmts()
	if err != nil {
		return db, err
	}

	return db, nil
}

/*
Close closes the connection to the database. It must be called before disposing
db.DB returned by db.New.
*/
func (db DB) Close() error {
	return db.sql.Close()
}

func stringListToDBList(list string) (uint, []byte) {
	count := uint(0)
	bytes := []byte(list)
	for index, character := range bytes {
		if character == ' ' {
			bytes[index] = ','
			count++
		}
	}

	return count, bytes
}

func validateID(id string) bool {
	for index := 0; index < len(id); index++ {
		/*
			URL Standard
			5.2. application/x-www-form-urlencoded serializing
			https://url.spec.whatwg.org/#urlencoded-serializing

			> 0x2A
			> 0x2D
			> 0x2E
			> 0x30 to 0x39
			> 0x41 to 0x5A
			> 0x5F
			> 0x61 to 0x7A
			>
			> Append a code point whose value is byte to output.

			Accept only those characters.
		*/
		if !(id[index] == 0x2A || id[index] == 0x2D || id[index] == 0x2E ||
			(id[index] >= 0x30 && id[index] <= 0x39) ||
			(id[index] >= 0x41 && id[index] <= 0x5A) ||
			id[index] == 0x5F ||
			(id[index] >= 0x61 && id[index] <= 0x7A)) {
			return false
		}
	}

	return true
}
