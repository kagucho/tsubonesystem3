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
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/kagucho/tsubonesystem3/backend/encoding"
	"log"
)

type ClubCommon struct {
	Name    string `json:"name"`
	Chief   string `json:"chief"`
}

// Club is a structure holding the information about a club.
type Club struct {
	ClubCommon
	Members ClubMemberChan `json:"members"`
}

type ClubEntryCommon struct {
	ClubCommon
	ID      string `json:"id"`
}

type ClubEntry struct {
	ClubEntryCommon
	Members []string `json:"members"`
}

// ClubEntryChan is a reciever of db.ClubEntry.
type ClubEntryChan <-chan ClubEntry

type clubEntry struct {
	ClubEntryCommon
	members *[]string `json:"members"`
}

/*
ClubMemberResult is a structure holding a result of querying the members
belonging to a club.
*/
type ClubMemberResult struct {
	ID    string
	Error error
}

// ClubMemberChan is a reciever of db.ClubMemberResult.
type ClubMemberChan <-chan ClubMemberResult

/*
MarshalJSON returns the JSON encoding of the remaining entries and closes the
channel.

This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (entryChan ClubEntryChan) MarshalJSON() ([]byte, error) {
	return encoding.MarshalJSONArray(func() (interface{}, error, bool) {
		result, present := <-entryChan
		return result, nil, present
	})
}

/*
MarshalJSON returns the JSON encoding of the members and closes the channel.
This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (memberChan ClubMemberChan) MarshalJSON() ([]byte, error) {
	return encoding.MarshalJSONArray(func() (interface{}, error, bool) {
		result, present := <-memberChan
		return result.ID, result.Error, present
	})
}

/*
DeleteClub deletes a club identified by the given ID.

It returns db.ErrIncorrectIdentity if the given ID is incorrect. Other errors
tell db.DB is bad.
*/
func (db DB) DeleteClub(id string) error {
	result, execErr := db.stmts[stmtDeleteClub].Exec(id)
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
}

/*
InsertClub inserts a club with the given properties.

It may return one of the following errors:
db.ErrBadOmission tells the ID or name is omitted.
db.ErrDupEntry tells a club with the given ID already exists.
db.ErrIncorrectIdentity tells the ID of the chief is incorrect.
db.ErrInvalid tells some of the given properties is invalid.

Other errors tell db.DB is bad.
*/
func (db DB) InsertClub(id, name, chief string) error {
	if id == `` || name == `` {
		return ErrBadOmission
	}

	if !validateID(id) {
		return ErrInvalid
	}

	_, err := db.stmts[stmtInsertClub].Exec(id, name, chief)

	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		switch mysqlErr.Number {
		case erDupEntry:
			return ErrDupEntry

		case erDataTooLong:
			fallthrough
		case erTruncatedWrongValueForField:
			return ErrInvalid

		case erNoReferencedRow:
			fallthrough
		case erNoReferencedRow2:
			return ErrIncorrectIdentity
		}
	}

	return err
}

/*
QueryClub returns db.Club corresponding with the given ID.

It returns db.ErrIncorrectIdentity if the given ID is incorrect. Other errors
tell db.DB is bad.

Resources will be holded until Members gets closed.
*/
func (db DB) QueryClub(id string) (Club, error) {
	var clubID uint8
	var club Club

	tx, err := db.sql.BeginTx(context.Background(),
		&sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  true,
		})
	if err != nil {
		return club, err
	}

	if err := tx.Stmt(db.stmts[stmtSelectClubByID]).QueryRow(id).Scan(
		&clubID, &club.Name, &club.Chief); err != nil {
		if err := tx.Commit(); err != nil {
			log.Print(err)
		}

		if err == sql.ErrNoRows {
			err = ErrIncorrectIdentity
		}

		return club, err
	}

	members := make(chan ClubMemberResult)

	go func() {
		defer func() {
			close(members)

			if err := tx.Commit(); err != nil {
				log.Print(err)
			}
		}()

		rows, err := tx.Stmt(db.stmts[stmtSelectMemberIDsByInternalClub]).Query(clubID)
		if err != nil {
			members <- ClubMemberResult{Error: err}
			return
		}

		defer rows.Close()

		for rows.Next() {
			var result ClubMemberResult
			result.Error = rows.Scan(&result.ID)

			members <- result
			if result.Error != nil {
				return
			}
		}
	}()

	club.Members = members

	return club, nil
}

/*
QueryClubName returns the name of the club identified with the given ID.

It returns db.ErrIncorrectIdentity if the given ID is incorrect. Other errors
tell db.DB is bad.
*/
func (db DB) QueryClubName(id string) (string, error) {
	var name string
	err := db.stmts[stmtSelectClubNameByID].QueryRow(id).Scan(&name)
	if err == sql.ErrNoRows {
		err = ErrIncorrectIdentity
	}

	return name, err
}

/*
QueryClubs returns db.ClubChan which represents all the clubs.

A result recieved from the returned channel tells an error if db.DB is bad.

Resources will be holded until the channel gets closed.
*/
func (db DB) QueryClubs() (ClubEntryChan, error) {
	tx, txErr := db.sql.BeginTx(context.Background(),
		&sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  true,
		})
	if txErr != nil {
		return nil, txErr
	}

	defer func() {
		if commitErr := tx.Commit(); commitErr != nil {
			log.Print(commitErr)
		}
	}()

	entries := make(map[uint8]clubEntry)

	clubsErr := func() error {
		rows, queryErr := db.stmts[stmtSelectClubs].Query()
		if queryErr != nil {
			return queryErr
		}

		defer func() {
			if closeErr := rows.Close(); closeErr != nil {
				log.Print(closeErr)
			}
		}()

		for rows.Next() {
			var dbID uint8
			var entry ClubEntryCommon
			var members []string

			scanErr := rows.Scan(&dbID,
				&entry.ID, &entry.Name, &entry.Chief)
			if scanErr != nil {
				return scanErr
			}

			entries[dbID] = clubEntry{entry, &members}
		}

		return nil
	}()
	if clubsErr != nil {
		return nil, clubsErr
	}

	membersErr := func() error {
		rows, queryErr := db.stmts[stmtSelectClubInternalIDMemberID].Query()
		if queryErr != nil {
			return queryErr
		}

		defer func() {
			if closeErr := rows.Close(); closeErr != nil {
				log.Print(closeErr)
			}
		}()

		for rows.Next() {
			var dbID uint8
			var memberID string

			scanErr := rows.Scan(&dbID, &memberID)
			if scanErr != nil {
				return scanErr
			}

			members := entries[dbID].members
			*members = append(*members, memberID)
		}

		return nil
	}()
	if membersErr != nil {
		return nil, membersErr
	}

	entryChan := make(chan ClubEntry)

	go func() {
		defer close(entryChan)

		for _, entry := range entries {
			entryChan <- ClubEntry{
				entry.ClubEntryCommon,
				*entry.members,
			}
		}
	}()

	return entryChan, nil
}

/*
UpdateClub updates the club identified with the given ID, with the given
properties.

It may return one of the following errors:
db.ErrIncorrectIdentity tells the given ID of the club or the one of the chief
is incorrect.
db.ErrInvalid tells some of the given properties is invalid.

Other errors tell db.DB is bad.
*/
func (db DB) UpdateClub(id, name, chief string) error {
	arguments := append(make([]interface{}, 0, 3), sql.Named(`id`, id))

	if name != `` {
		arguments = append(arguments, sql.Named(`name`, name))
	}

	if chief != `` {
		arguments = append(arguments, sql.Named(`chief`, chief))
	}

	result, execErr := db.stmts[stmtUpdateClub].Exec(arguments...)
	if execErr != nil {
		if mysqlErr, ok := execErr.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case erDataTooLong:
				fallthrough
			case erTruncatedWrongValueForField:
				return ErrInvalid

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
