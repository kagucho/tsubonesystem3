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

// Attendance is a type representing an attendance.
type Attendance uint16

// The values Attendance may have.
const (
	AttendanceUninvited Attendance = iota
	AttendanceInvited
	AttendanceDeclined
	AttendanceAccepted
)

/*
PartyCommon is a structure holding information about a party common for queries.
*/
type PartyCommon struct {
	Creator  encoding.ZeroString `json:"creator"`
	Start    encoding.Time       `json:"start"`
	End      encoding.Time       `json:"end"`
	Place    string              `json:"place"`
	Inviteds string              `json:"inviteds"`
	Due      encoding.Time       `json:"due"`
}

// PartyEntry is a strucutre holding basic information about a party.
type PartyEntry struct {
	PartyCommon
	Name string `json:"name"`
}

/*
PartyAttendanceResult is a structure representing a result of querying
db.PartyAttendance.
*/
type PartyAttendanceResult struct {
	Member     string
	Attendance Attendance
	Error      error
}

// PartyAttendanceChan is a reciever of db.PartyAttendanceResult.
type PartyAttendanceChan <-chan PartyAttendanceResult

// PartyDetail is a structure holding details of a party.
type PartyDetail struct {
	PartyCommon
	Details     string              `json:"details"`
	Attendances PartyAttendanceChan `json:"attendances"`
}

/*
PartyUser is a structure holding information of a party and the attendance of
the user.
*/
type PartyUser struct {
	PartyEntry
	User Attendance `json:"user"`
}

/*
PartyUserResult is a structure representing a result of querying db.PartyUser.
*/
type PartyUserResult struct {
	PartyUser
	Error error
}

// PartyUserChan is a reciever of db.PartyUserResult.
type PartyUserChan <-chan PartyUserResult

/*
PartyNameResult is a structure representing a result of querying names of
parties.
*/
type PartyNameResult struct {
	Error error
	Name  string
}

// PartyNameChan is a reciever of db.PartyNameResult.
type PartyNameChan <-chan PartyNameResult

var attendanceJSON = [][]byte{
	AttendanceUninvited: []byte(`"uninvited"`),
	AttendanceInvited:   []byte(`"invited"`),
	AttendanceDeclined:  []byte(`"declined"`),
	AttendanceAccepted:  []byte(`"accepted"`),
}

/*
MarshalJSON returns the JSON encoding of the attendance.

This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (attendance Attendance) MarshalJSON() ([]byte, error) {
	return attendanceJSON[attendance], nil
}

/*
MarshalJSON returns the JSON encoding of the remaining attendances and closes
the channel.

This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (attendanceChan PartyAttendanceChan) MarshalJSON() ([]byte, error) {
	return encoding.MarshalJSONObject(func() (string, interface{}, error, bool) {
		result, present := <-attendanceChan
		return result.Member, result.Attendance, result.Error, present
	})
}

/*
MarshalJSON returns the JSON encoding of the remaining parties and closes the
channel.

This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (userChan PartyUserChan) MarshalJSON() ([]byte, error) {
	return encoding.MarshalJSONArray(func() (interface{}, error, bool) {
		result, present := <-userChan
		return result.PartyUser, result.Error, present
	})
}

/*
MarshalJSON returns the JSON encoding of the remaining names and closes the
channel.

This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (nameChan PartyNameChan) MarshalJSON() ([]byte, error) {
	return encoding.MarshalJSONArray(func() (interface{}, error, bool) {
		result, present := <-nameChan
		return result.Name, result.Error, present
	})
}

/*
DeleteParty deletes a party.

It returns db.ErrIncorrectIdentity if any party named with the given name and
created by the member identified by the given ID does not exist. Other errors
tell db.DB is bad.
*/
func (db DB) DeleteParty(name, creator string) error {
	result, execErr := db.stmts[stmtDeleteParty].Exec(name, creator)
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
InsertParty inserts a party with the given properties and returns the email
addresses of the invited members.

It may return one of the following errors:
db.ErrBadOmission tells some of the parameters is omitted.
db.ErrDupEntry tells the name is duplicate.
db.ErrIncorrectIdentity tells the ID of the creator or one of the invited
members is incorrect.
db.ErrInvalid tells some of the properties is invalid.
*/
func (db DB) InsertParty(name, creator string, start, end encoding.Time, place string, due encoding.Time, invitedIDs, inviteds, details string) ([]string, error) {
	if (name == `` || start == encoding.Time{} || end == encoding.Time{} || place == `` || due == encoding.Time{} || invitedIDs == `` || inviteds == `` || details == ``) {
		return nil, ErrBadOmission
	}

	invitedsNumber, invitedsDB := stringListToDBList(invitedIDs)
	invitedMails := make([]string, 0, invitedsNumber)

	rows, err := db.stmts[stmtCallInsertParty].Query(name, creator,
		start.Generic(), end.Generic(), place, due.Generic(),
		inviteds, invitedsDB, invitedsNumber, details)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case erDataTooLong:
				fallthrough
			case erTruncatedWrongValueForField:
				fallthrough
			case erWrongValue:
				return nil, ErrInvalid

			case erDupEntry:
				return nil, ErrDupEntry

			case erNoReferencedRow:
				fallthrough
			case erNoReferencedRow2:
				fallthrough
			case erSignalException:
				return nil, ErrIncorrectIdentity
			}
		}

		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Print(err)
		}
	}()

	for rows.Next() {
		var mail string

		if err := rows.Scan(&mail); err != nil {
			return nil, err
		}

		if mail != `` {
			invitedMails = append(invitedMails, mail)
		}
	}

	return invitedMails, nil
}

/*
QueryParties returns db.PartyUserChan representing all parties.

Resources will be holded until the channel gets closed.
*/
func (db DB) QueryParties(user string) PartyUserChan {
	resultChan := make(chan PartyUserResult)

	go func() {
		defer close(resultChan)

		var result PartyUserResult
		var tx *sql.Tx

		/*
			TODO: This memory allocation is so innocent, but that
			can be improved with API changes in database/sql. See
			mysql_store_result in MySQL Connect/C.
		*/
		parties := make(map[uint16]PartyEntry)

		tx, result.Error = db.sql.BeginTx(context.Background(),
			&sql.TxOptions{
				Isolation: sql.LevelSerializable,
				ReadOnly:  true,
			})
		if result.Error != nil {
			resultChan <- result

			return
		}

		result.Error = func() error {
			rows, err := tx.Stmt(db.stmts[stmtSelectParties]).Query()
			if err != nil {
				return err
			}

			defer func() {
				if err := rows.Close(); err != nil {
					log.Print(err)
				}
			}()

			for rows.Next() {
				var id uint16
				var creator sql.NullString
				var party PartyEntry
				var start mysql.NullTime
				var end mysql.NullTime
				var due mysql.NullTime

				if err := rows.Scan(&id,
					&party.Name, &creator,
					&start, &end,
					&party.Place, &party.Inviteds,
					&due); err != nil {
					return err
				}

				party.Creator = encoding.ZeroString(creator.String)
				party.Start = encoding.NewTime(start.Time)
				party.End = encoding.NewTime(end.Time)
				party.Due = encoding.NewTime(due.Time)

				parties[id] = party
			}

			return nil
		}()
		if result.Error != nil {
			resultChan <- result

			return
		}

		/*
			TODO: This is also bad of Golang... Go 1.8 cannot
			efficiently read-and-modify items in map, so prepare
			another map.
		*/
		attendances := make(map[uint16]Attendance, len(parties))

		result.Error = func() error {
			rows, err := tx.Stmt(db.stmts[stmtSelectAttendancesByMember]).Query(user)
			if err != nil {
				return err
			}

			defer func() {
				if err := rows.Close(); err != nil {
					log.Print(err)
				}
			}()

			for rows.Next() {
				var attendance uint16
				var id uint16

				if err := rows.Scan(&id, &attendance); err != nil {
					return err
				}

				attendances[id] = Attendance(attendance)
			}

			return nil
		}()
		if result.Error != nil {
			resultChan <- result

			return
		}

		for id, party := range parties {
			partyUser := PartyUser{PartyEntry: party}

			if attendance, ok := attendances[id]; ok {
				partyUser.User = attendance
			}

			resultChan <- PartyUserResult{PartyUser: partyUser}
		}
	}()

	return resultChan
}

/*
QueryParty queries details of a party identified by the given name.

It returns db.ErrIncorrectIdentity if the name is incorrect.

Resources will be holded until Attendances gets closed.
*/
func (db DB) QueryParty(name string) (PartyDetail, error) {
	var party PartyDetail
	var creator sql.NullString
	var start mysql.NullTime
	var end mysql.NullTime
	var due mysql.NullTime
	var id uint16

	tx, err := db.sql.BeginTx(context.Background(),
		&sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  true,
		})
	if err != nil {
		return party, err
	}

	if err := tx.Stmt(db.stmts[stmtSelectParty]).QueryRow(name).Scan(
		&id, &creator, &start, &end, &party.Place,
		&party.Inviteds, &due, &party.Details); err != nil {
		if err := tx.Commit(); err != nil {
			log.Print(err)
		}

		if err == sql.ErrNoRows {
			return party, ErrIncorrectIdentity
		}

		return party, err
	}

	attendances := make(chan PartyAttendanceResult)

	go func() {
		defer func() {
			close(attendances)

			if err := tx.Commit(); err != nil {
				log.Print(err)
			}
		}()

		rows, err := tx.Stmt(db.stmts[stmtSelectAttendancesByInternalParty]).Query(id)
		if err != nil {
			attendances <- PartyAttendanceResult{Error: err}
			return
		}

		defer func() {
			if err := rows.Close(); err != nil {
				log.Print(err)
			}
		}()

		for rows.Next() {
			var result PartyAttendanceResult
			result.Error = rows.Scan(&result.Member, (*uint16)(&result.Attendance))

			attendances <- result
			if result.Error != nil {
				return
			}
		}
	}()

	party.Creator = encoding.ZeroString(creator.String)
	party.Start = encoding.NewTime(start.Time)
	party.End = encoding.NewTime(end.Time)
	party.Due = encoding.NewTime(due.Time)
	party.Attendances = attendances

	return party, nil
}

/*
UpdateParty updates the party identified by the given name and created by the
member identified by the given ID, with the given properties.

It may return one of the following errors:
ErrIncorrectIdentity tells the name of the party or the ID of the creator or
one of the invited members is incorrect.
ErrInvalid tells some of the properties is invalid.

Other errors tell db.DB is bad.
*/
func (db DB) UpdateParty(name, creator string, start, end encoding.Time, place string, due encoding.Time, inviteds, invitedIDs, details string) error {
	arguments := make([]interface{}, 10)
	arguments[0] = name
	arguments[1] = creator

	if (start != encoding.Time{}) {
		arguments[2] = start.Generic()
	}

	if (end != encoding.Time{}) {
		arguments[3] = end.Generic()
	}

	if place != `` {
		arguments[4] = place
	}

	if (due != encoding.Time{}) {
		arguments[5] = due.Generic()
	}

	if inviteds != `` {
		arguments[8] = inviteds
	}

	if invitedIDs != `` {
		arguments[7], arguments[6] = stringListToDBList(invitedIDs)
	}

	if details != `` {
		arguments[9] = details
	}

	result, execErr := db.stmts[stmtCallUpdateParty].Exec(arguments...)
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
				fallthrough
			case erSignalException:
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

/*
UpdateAttendance updates attendance of a member identified by the given ID for
the party identified by the given name.

It returns ErrIncorrectIdentity if the given ID or name is incorrect. Other
errors tell db.DB is bad.
*/
func (db DB) UpdateAttendance(attending bool, party, member string) error {
	attendance := 2
	if attending {
		attendance = 3
	}

	result, execErr := db.stmts[stmtUpdateAttendance].Exec(attendance, party, member)
	if execErr != nil {
		return execErr
	}

	affected, affectedErr := result.RowsAffected()
	if affectedErr != nil {
		return affectedErr
	}

	if affected != 1 {
		return ErrIncorrectIdentity
	}

	return nil
}
