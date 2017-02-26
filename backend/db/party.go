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
	"github.com/go-sql-driver/mysql"
	"github.com/kagucho/tsubonesystem3/json"
	"log"
	"time"
)

type Attendance uint

const (
	AttendanceUninvited Attendance = iota
	AttendanceInvited
	AttendanceAccepted
	AttendanceDeclined
)

type Party struct {
	Name          string   `json:"name"`
	Start         NullTime `json:"start"`
	End           NullTime `json:"end"`
	Place         string   `json:"place"`
	Inviteds      string   `json:"inviteds"`
	Due           NullTime `json:"due"`
}

type PartyUser struct {
	Party
	User Attendance `json:"user"`
}

type PartyUserChan <-chan PartyUserResult

type PartyUserResult struct {
	Error error
	Value PartyUser
}

type PartyNameChan <-chan PartyNameResult

type PartyNameResult struct {
	Error error
	Value string
}

var attendanceJSON = [][]byte{
	AttendanceUninvited: []byte(`"uninvited"`),
	AttendanceInvited:   []byte(`"invited"`),
	AttendanceAccepted:  []byte(`"accepted"`),
	AttendanceDeclined:  []byte(`"declined"`),
}

func (attendance Attendance) MarshalJSON() ([]byte, error) {
	return attendanceJSON[attendance], nil
}

func (partyUserChan PartyUserChan) MarshalJSON() ([]byte, error) {
	return json.MarshalChan(partyUserChan)
}

func (nameChan PartyNameChan) MarshalJSON() ([]byte, error) {
	return json.MarshalChan(nameChan)
}

func (db DB) InsertParty(name string, start, end time.Time, place string, due time.Time, inviteds, invitedsName, details string) (returningMails []string, returningError error) {
	invitedIDs, invitedMails, invitedError := db.queryMemberInternalIDMails(inviteds)
	if invitedError != nil {
		return nil, invitedError
	}

	tx, txError := db.sql.Begin()
	if txError != nil {
		return nil, txError
	}

	defer func() {
		if recovered := recover(); recovered == nil {
			returningError = tx.Commit()
		} else {
			if rollbackError := tx.Rollback(); rollbackError != nil {
				log.Print(rollbackError)
			}

			var ok bool
			returningError, ok = recovered.(error)
			if !ok {
				panic(recovered)
			}
		}
	}()

	result, execError := tx.Stmt(db.stmts[stmtInsertParty]).Exec(name, start, end, place, due, invitedsName, details)
	if execError != nil {
		panic(execError)
	}

	partyID, idError := result.LastInsertId()
	if idError != nil {
		panic(idError)
	}

	for _, invitedID := range invitedIDs {
		if _, execError := tx.Stmt(db.stmts[stmtInsertInternalAttendance]).Exec(partyID, invitedID); execError != nil {
			panic(execError)
		}
	}

	return invitedMails, nil
}

func (db DB) QueryParties(user string) PartyUserChan {
	resultChan := make(chan PartyUserResult)

	go func() {
		defer close(resultChan)

		var result PartyUserResult

		/*
			TODO: This memory allocation is so innocent, but it can
			be improved with API changes in database/sql. See
			mysql_store_result in MySQL Connect/C.
		*/
		parties := make(map[uint16]Party)

		result.Error = func() error {
			rows, queryError := db.stmts[stmtSelectParties].Query()
			if queryError != nil {
				return queryError
			}

			defer func() {
				if closeError := rows.Close(); closeError != nil {
					log.Print(closeError)
				}
			}()

			for rows.Next() {
				var id uint16
				var party Party

				if scanError := rows.Scan(&id, &party.Name,
					(*mysql.NullTime)(&party.Start),
					(*mysql.NullTime)(&party.End),
					&party.Place, &party.Inviteds,
					(*mysql.NullTime)(&party.Due)); scanError != nil {
					return scanError
				}

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
			rows, queryError := db.stmts[stmtSelectAttendancesByMember].Query(user)
			if queryError != nil {
				return queryError
			}

			defer func() {
				if closeError := rows.Close(); closeError != nil {
					log.Print(closeError)
				}
			}()

			for rows.Next() {
				var attending sql.NullBool
				var id uint16
				if scanError := rows.Scan(&id, &attending); scanError != nil {
					return scanError
				}

				if !attending.Valid {
					attendances[id] = AttendanceInvited
				} else if attending.Bool {
					attendances[id] = AttendanceAccepted
				} else {
					attendances[id] = AttendanceDeclined
				}
			}

			return nil
		}()
		if result.Error != nil {
			resultChan <- result

			return
		}

		for id, party := range parties {
			partyUser := PartyUser{Party: party}

			if attendance, ok := attendances[id]; ok {
				partyUser.User = attendance
			}

			resultChan <- PartyUserResult{Value: partyUser}
		}
	}()

	return resultChan
}

func (db DB) QueryPartyNames() PartyNameChan {
	resultChan := make(chan PartyNameResult)

	go func() {
		defer close(resultChan)

		rows, queryError := db.stmts[stmtSelectPartyNames].Query()
		if queryError != nil {
			resultChan <- PartyNameResult{Error: queryError}

			return
		}

		defer func() {
			if closeError := rows.Close(); closeError != nil {
				log.Print(closeError)
			}
		}()

		for rows.Next() {
			var result PartyNameResult
			result.Error = rows.Scan(&result.Value)

			resultChan <- result

			if result.Error != nil {
				return
			}
		}
	}()

	return resultChan
}

func (db DB) UpdateAttendance(attending bool, party, member string) error {
	result, execError := db.stmts[stmtUpdateAttendance].Exec(attending, party, member)
	if execError != nil {
		return execError
	}

	affected, affectedError := result.RowsAffected()
	if affectedError != nil {
		return affectedError
	}

	if affected != 1 {
		return IncorrectIdentity
	}

	return nil
}
