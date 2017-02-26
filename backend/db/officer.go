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
	"github.com/kagucho/tsubonesystem3/json"
	"strings"
)

type OfficerCommon struct {
	Member OfficerMember `json:"member"`
	Name   string        `json:"name"`
}

// OfficerDetail is a structure to hold the information about an officer
type OfficerDetail struct {
	OfficerCommon
	Scope []string `json:"scope"`
}

type OfficerEntry struct {
	OfficerCommon
	ID string `json:"id"`
}

type OfficerEntryChan <-chan OfficerEntryResult

type OfficerEntryResult struct {
	Error error
	Value OfficerEntry
}

// OfficerMember is a structure to hold the information about a member who is
// an officer.
type OfficerMember struct {
	ID       string `json:"id"`
	Mail     string `json:"mail"`
	Nickname string `json:"nickname"`
	Realname string `json:"realname,omitempty"`
	Tel      string `json:"tel,omitempty"`
}

func (entryChan OfficerEntryChan) MarshalJSON() ([]byte, error) {
	return json.MarshalChan(entryChan)
}

// QueryOfficerDetail returns db.OfficerDetail of the officer identified with
// the given ID.
func (db DB) QueryOfficerDetail(id string) (OfficerDetail, error) {
	var detail OfficerDetail
	var scope string

	scanError := db.stmts[stmtSelectOfficer].QueryRow(id).Scan(
		&detail.Name, &scope,
		&detail.Member.ID, &detail.Member.Mail,
		&detail.Member.Nickname, &detail.Member.Realname,
		&detail.Member.Tel)
	if scanError == sql.ErrNoRows {
		return detail, IncorrectIdentity
	} else if scanError != nil {
		return detail, scanError
	}

	detail.Scope = strings.Split(scope, `,`)

	return detail, nil
}

// QueryOfficerName returns the name of the officer identified with the given
// ID.
func (db DB) QueryOfficerName(id string) (string, error) {
	var name string
	scanError := db.stmts[stmtSelectOfficerName].QueryRow(id).Scan(&name)
	if scanError == sql.ErrNoRows {
		scanError = IncorrectIdentity
	}

	return name, scanError
}

// QueryOfficers returns db.OfficerEntryChan which represents all the officers.
func (db DB) QueryOfficers() OfficerEntryChan {
	resultChan := make(chan OfficerEntryResult)

	go func() {
		defer close(resultChan)

		rows, queryError := db.stmts[stmtSelectOfficers].Query()
		if queryError != nil {
			resultChan <- OfficerEntryResult{Error: queryError}
			return
		}

		defer rows.Close()

		for rows.Next() {
			var result OfficerEntryResult

			result.Error = rows.Scan(&result.Value.ID, &result.Value.Name,
				&result.Value.Member.ID,
				&result.Value.Member.Mail,
				&result.Value.Member.Nickname,
				&result.Value.Member.Realname,
				&result.Value.Member.Tel)

			resultChan <- result

			if result.Error != nil {
				return
			}
		}
	}()

	return resultChan
}
