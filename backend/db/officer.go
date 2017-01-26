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
	"github.com/kagucho/tsubonesystem3/chanjson"
	"strings"
)

// Officer is a structure to hold the information about an officer
type Officer struct {
	Member OfficerMember `json:"member"`
	Name   string        `json:"name"`
	Scope  []string      `json:"scope"`
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

type officer struct {
	ID     string        `json:"id"`
	Member OfficerMember `json:"member"`
	Name   string        `json:"name"`
}

type officerResult struct {
	Error error
	Value officer
}

// QueryOfficer returns db.Officer of the officer identified with the given ID.
func (db DB) QueryOfficer(id string) (Officer, error) {
	var detail Officer
	var member uint16
	var scope string

	if scanError := db.stmts[stmtSelectOfficer].QueryRow(id).Scan(&member, &detail.Name, &scope); scanError != nil {
		return Officer{}, scanError
	}

	if scanError := db.stmts[stmtSelectOfficerMemberInternal].QueryRow(member).Scan(
		&detail.Member.ID, &detail.Member.Mail,
		&detail.Member.Nickname, &detail.Member.Realname,
		&detail.Member.Tel); scanError != nil {
		return Officer{}, scanError
	}

	detail.Scope = strings.Split(scope, `,`)

	return detail, nil
}

// QueryOfficerName returns the name of the officer identified with the given
// ID.
func (db DB) QueryOfficerName(id string) (string, error) {
	var name string
	scanError := db.stmts[stmtSelectOfficerName].QueryRow(id).Scan(&name)

	return name, scanError
}

// QueryOfficers returns chanjson.ChanJSON which represents all the officers.
func (db DB) QueryOfficers() chanjson.ChanJSON {
	resultChan := make(chan officerResult)

	go func() {
		defer close(resultChan)

		rows, queryError := db.stmts[stmtSelectOfficers].Query()
		if queryError != nil {
			resultChan <- officerResult{Error: queryError}
			return
		}

		defer rows.Close()

		for rows.Next() {
			var member uint16
			var result officerResult

			result.Error = rows.Scan(&result.Value.ID, &member, &result.Value.Name)
			if result.Error != nil {
				resultChan <- result
				return
			}

			result.Error = db.stmts[stmtSelectOfficerMemberInternal].QueryRow(member).Scan(
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

	return chanjson.New(resultChan)
}
