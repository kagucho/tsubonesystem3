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

// Package db implements an abstraction for the database.
package db

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/kagucho/tsubonesystem3/configuration"
	"github.com/kagucho/tsubonesystem3/json"
	"log"
)

type NullTime mysql.NullTime

// DB is a structure to keep the connection to the database.
type DB struct {
	sql   *sql.DB
	stmts [stmtNumber]*sql.Stmt
}

var IncorrectIdentity = errors.New(`incorrect identity`)

func (nullTime NullTime) MarshalJSON() ([]byte, error) {
	if !nullTime.Valid {
		return []byte(`null`), nil
	}

	return json.MarshalTime(nullTime.Time)
}

// Prepare prepares the database.
func Prepare() (DB, error) {
	var db DB
	var prepareError error

	db.sql, prepareError = sql.Open(`mysql`, configuration.DBDSN)
	if prepareError != nil {
		return db, prepareError
	}

	_, prepareError = db.sql.Exec(`SET sql_mode='TRADITIONAL'`)
	if prepareError != nil {
		if closeError := db.Close(); closeError != nil {
			log.Print(closeError)
		}

		return db, prepareError
	}

	prepareError = db.prepareStmts()
	if prepareError != nil {
		return db, prepareError
	}

	return db, nil
}

// Close closes the connection to the database.
func (db DB) Close() error {
	return db.sql.Close()
}
