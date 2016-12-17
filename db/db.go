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

// Package db implements an abstraction for the database.
package db

import (
  `database/sql`

  // Any Driver import should be blank.
  _ `github.com/go-sql-driver/mysql`

  `github.com/kagucho/tsubonesystem3/configuration`
)

// DB is a structure to keep the connection to the database.
type DB struct {
  sql *sql.DB
}

// Open returns a new connection to the database.
func Open() (DB, error) {
  sql, openError := sql.Open(`mysql`, configuration.DBDSN)
  return DB{sql}, openError
}

// Close closes the connection to the database.
func (db DB) Close() error {
  return db.sql.Close()
}
