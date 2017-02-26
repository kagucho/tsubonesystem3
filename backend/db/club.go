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

// Chief is a structure to hold the information about a member.
type Chief struct {
	ID       string `json:"id"`
	Mail     string `json:"mail"`
	Nickname string `json:"nickname"`
	Realname string `json:"realname",omitempty`
	Tel      string `json:"tel",omitempty`
}

// ClubDetail is a structure to hold the information about a club.
type ClubDetail struct {
	Chief   Chief          `json:"chief"`
	Members ClubMemberChan `json:"members"`
	Name    string         `json:"name"`
}

type ClubMemberChan <-chan ClubMemberResult

type ClubMemberResult struct {
	Error error
	Value struct {
		Entrance uint16 `json:"entrance",omitempty`
		ID       string `json:"id"`
		Nickname string `json:"nickname"`
		Realname string `json:"realname",omitempty`
	}
}

type ClubName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ClubNameChan <-chan ClubNameResult

type ClubNameResult struct {
	Error error
	Value ClubName
}

type ClubEntry struct {
	ClubName
	Chief Chief `json:"chief"`
}

type ClubEntryChan <-chan ClubEntryResult

type ClubEntryResult struct {
	Error error
	Value ClubEntry
}

func (entryChan ClubEntryChan) MarshalJSON() ([]byte, error) {
	return json.MarshalChan(entryChan)
}

func (memberChan ClubMemberChan) MarshalJSON() ([]byte, error) {
	return json.MarshalChan(memberChan)
}

func (nameChan ClubNameChan) MarshalJSON() ([]byte, error) {
	return json.MarshalChan(nameChan)
}

// QueryClubDetail returns db.ClubDetail corresponding with the given ID.
func (db DB) QueryClubDetail(id string) (ClubDetail, error) {
	var clubID uint8
	var club ClubDetail

	if scanError := db.stmts[stmtSelectClub].QueryRow(id).Scan(
		&clubID, &club.Name,
		&club.Chief.ID, &club.Chief.Mail,
		&club.Chief.Nickname, &club.Chief.Realname,
		&club.Chief.Tel); scanError == sql.ErrNoRows {
		return ClubDetail{}, IncorrectIdentity
	} else if scanError != nil {
		return ClubDetail{}, scanError
	}

	members := make(chan ClubMemberResult)

	go func() {
		defer close(members)

		rows, queryError := db.stmts[stmtSelectMembersByClub].Query(clubID)
		if queryError != nil {
			members <- ClubMemberResult{Error: queryError}
			return
		}

		defer rows.Close()

		for rows.Next() {
			var result ClubMemberResult

			result.Error = rows.Scan(
				&result.Value.Entrance, &result.Value.ID,
				&result.Value.Nickname, &result.Value.Realname)

			members <- result

			if result.Error != nil {
				return
			}
		}
	}()

	club.Members = members

	return club, nil
}

// QueryClubName returns the name of the club identified with the given ID.
func (db DB) QueryClubName(id string) (string, error) {
	var name string
	scanError := db.stmts[stmtSelectClubName].QueryRow(id).Scan(&name)
	if scanError == sql.ErrNoRows {
		scanError = IncorrectIdentity
	}

	return name, scanError
}

// QueryClubNames returns db.ClubNameChan which represents the names of all the
// clubs.
func (db DB) QueryClubNames() ClubNameChan {
	resultChan := make(chan ClubNameResult)

	go func() {
		defer close(resultChan)

		rows, queryError := db.stmts[stmtSelectClubNames].Query()
		if queryError != nil {
			resultChan <- ClubNameResult{Error: queryError}
			return
		}

		for rows.Next() {
			var result ClubNameResult

			result.Error = rows.Scan(
				&result.Value.ID, &result.Value.Name)

			resultChan <- result

			if result.Error != nil {
				return
			}
		}
	}()

	return resultChan
}

// QueryClubs returns db.ClubEntryChan which represents all the clubs.
func (db DB) QueryClubs() ClubEntryChan {
	resultChan := make(chan ClubEntryResult)

	go func() {
		defer close(resultChan)

		rows, queryError := db.stmts[stmtSelectClubs].Query()
		if queryError != nil {
			resultChan <- ClubEntryResult{Error: queryError}
			return
		}

		for rows.Next() {
			var result ClubEntryResult

			result.Error = rows.Scan(
				&result.Value.ID, &result.Value.Name,
				&result.Value.Chief.ID, &result.Value.Chief.Mail,
				&result.Value.Chief.Nickname, &result.Value.Chief.Realname,
				&result.Value.Chief.Tel)

			resultChan <- result
		}
	}()

	return resultChan
}

func (db DB) txQueryInternalClubs(tx *sql.Tx, clubs []string) (map[uint8]struct{}, error) {
	clubIDs := make(map[uint8]struct{}, len(clubs))

	rows, queryError := tx.Stmt(db.stmts[stmtSelectInternalClubs]).Query(strings.Join(clubs, `,`))
	if queryError != nil {
		return nil, queryError
	}

	defer rows.Close()

	for rows.Next() {
		var id uint8
		if scanError := rows.Scan(&id); scanError != nil {
			return nil, scanError
		}

		clubIDs[id] = struct{}{}
	}

	return clubIDs, nil
}
