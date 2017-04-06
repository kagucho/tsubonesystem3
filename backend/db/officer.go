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
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/kagucho/tsubonesystem3/backend/encoding"
	"log"
	"strings"
)

// OfficerName is a structure associating the ID and name of an officer.
type OfficerName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

/*
OfficerNameResult is a strucutre representing a result of querying
db.OfficerName.
*/
type OfficerNameResult struct {
	OfficerName
	Error error
}

// OfficerNameChan is a reciever of db.OfficerNameResult.
type OfficerNameChan <-chan OfficerNameResult

// OfficerDetail is a structure holding details of an officer
type OfficerDetail struct {
	Member string   `json:"member"`
	Name   string   `json:"name"`
	Scope  []string `json:"scope"`
}

// OfficerEntry is a structure holding basic information of an officer.
type OfficerEntry struct {
	OfficerName
	Member string `json:"member"`
}

/*
OfficerEntryResult is a strcuture representing a result of querying
db.OfficerEntry.
*/
type OfficerEntryResult struct {
	OfficerEntry
	Error error
}

// OfficerEntryChan is a reciever of db.OfficerEntryResult.
type OfficerEntryChan <-chan OfficerEntryResult

// NoScopeUpdate is a string telling not to upate the scope.
const NoScopeUpdate = "\000"

/*
ErrOfficerSuicide is an error telling the operator is removing his own
management permission.
*/
var ErrOfficerSuicide = errors.New(`removing operator's own management permission`)

/*
MarshalJSON returns the JSON encoding of the remaining entries and closes the
channel.

This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (entryChan OfficerEntryChan) MarshalJSON() ([]byte, error) {
	return encoding.MarshalJSONArray(func() (interface{}, error, bool) {
		result, present := <-entryChan
		return result.OfficerEntry, result.Error, present
	})
}

/*
DeleteOfficer deletes an officer identified by the given ID.

It returns db.ErrIncorrectIdentity if the given ID of the operator or the
officer is incorrect. It returns db.ErrOfficerSuicide if the operation is
expected to remove the management permission of the operator. Other errors tell
db.DB is bad.
*/
func (db DB) DeleteOfficer(operator, id string) error {
	result, execErr := db.stmts[stmtCallDeleteOfficer].Exec(operator, id)
	if execErr != nil {
		if mysqlErr, ok := execErr.(*mysql.MySQLError); ok && mysqlErr.Number == erSignalException {
			return ErrOfficerSuicide
		}

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
}

/*
InsertOfficer inserts an operator with the given properties.

It may return one of the following errors:
db.ErrBadOmission tells the ID or name is omitted.
db.ErrDupEntry tells the ID or name is duplicate.
db.ErrIncorrectIdentity tells the ID of the member is incorrect.
db.ErrInvalid tells the name of the scope is invalid.

Other errors tell db.DB is bad.
*/
func (db DB) InsertOfficer(id, name, member, scope string) error {
	_, scopeBytes := stringListToDBList(scope)
	_, err := db.stmts[stmtInsertOfficer].Exec(id, name, scopeBytes, member)

	if id == `` || name == `` {
		return ErrBadOmission
	}

	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		switch mysqlErr.Number {
		case erDataTooLong:
			fallthrough
		case erTruncatedWrongValueForField:
			return ErrInvalid

		case erDupEntry:
			return ErrDupEntry

		case erNoReferencedRow:
			fallthrough
		case erNoReferencedRow2:
			return ErrIncorrectIdentity
		}
	}

	return err
}

/*
QueryOfficerDetail returns db.OfficerDetail of the officer identified with the
given ID.

It returns db.ErrIncorrectIdentity if the given ID is incorrect. Other errors
tell db.DB is bad.
*/
func (db DB) QueryOfficerDetail(id string) (OfficerDetail, error) {
	var detail OfficerDetail
	var scope string

	err := db.stmts[stmtSelectOfficerByID].QueryRow(id).Scan(
		&detail.Name, &scope, &detail.Member)
	if err == sql.ErrNoRows {
		return detail, ErrIncorrectIdentity
	} else if err != nil {
		return detail, err
	}

	detail.Scope = strings.Split(scope, `,`)

	return detail, nil
}

/*
QueryOfficerName returns the name of the officer identified with the given ID.

It returns db.ErrIncorrectIdentity if the given ID is incorrect. Other errors
tell db.DB is bad.
*/
func (db DB) QueryOfficerName(id string) (string, error) {
	var name string
	err := db.stmts[stmtSelectOfficerNameByID].QueryRow(id).Scan(&name)
	if err == sql.ErrNoRows {
		err = ErrIncorrectIdentity
	}

	return name, err
}

/*
QueryOfficerNames returns db.OfficerNameChan representing the names of all
officers.

Resources will be holded until the channel gets closed.
*/
func (db DB) QueryOfficerNames() OfficerNameChan {
	resultChan := make(chan OfficerNameResult)

	go func() {
		defer close(resultChan)

		rows, err := db.stmts[stmtSelectOfficerIDNames].Query()
		if err != nil {
			resultChan <- OfficerNameResult{Error: err}
			return
		}

		defer func() {
			if err := rows.Close(); err != nil {
				log.Print(err)
			}
		}()

		for rows.Next() {
			var result OfficerNameResult
			result.Error = rows.Scan(&result.ID, &result.Name)

			resultChan <- result
			if result.Error != nil {
				return
			}
		}
	}()

	return resultChan
}

/*
QueryOfficers returns db.OfficerEntryChan which represents all the officers.

Resources will be holded until the channel gets closed.
*/
func (db DB) QueryOfficers() OfficerEntryChan {
	resultChan := make(chan OfficerEntryResult)

	go func() {
		defer close(resultChan)

		rows, err := db.stmts[stmtSelectOfficers].Query()
		if err != nil {
			resultChan <- OfficerEntryResult{Error: err}
			return
		}

		defer func() {
			if err := rows.Close(); err != nil {
				log.Print(err)
			}
		}()

		for rows.Next() {
			var result OfficerEntryResult

			result.Error = rows.Scan(
				&result.ID, &result.Name, &result.Member)

			resultChan <- result
			if result.Error != nil {
				return
			}
		}
	}()

	return resultChan
}

/*
UpdateOfficer updates the officer identified by the given ID, with the given
properties.

It may return one of the following properties:
db.ErrDupEntry tells the name is duplicate.
db.ErrIncorrectIdentity tells the ID of the officer or member is incorrect.
db.ErrInvalid tells some of the properties is invalid.

Other errors tell db.DB is bad.
*/
func (db DB) UpdateOfficer(operator, id, name, member, scope string) error {
	arguments := make([]interface{}, 5)
	arguments[0] = operator
	arguments[1] = id

	if name != `` {
		arguments[2] = name
	}

	if member != `` {
		arguments[3] = member
	}

	if scope != NoScopeUpdate {
		_, arguments[4] = stringListToDBList(scope)
	}

	result, execErr := db.stmts[stmtCallUpdateOfficer].Exec(arguments...)
	if execErr != nil {
		if mysqlErr, ok := execErr.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case erDataTooLong:
				fallthrough
			case erTruncatedWrongValueForField:
				return ErrInvalid

			case erDupEntry:
				return ErrDupEntry

			case erSignalException:
				return ErrOfficerSuicide

			case erNoReferencedRow:
				fallthrough
			case erNoReferencedRow2:
				return ErrIncorrectIdentity
			}
		}

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
}
