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

import "github.com/kagucho/tsubonesystem3/chanjson"

// Chief is a structure to hold the information about a member.
type Chief struct {
	ID       string `json:"id"`
	Mail     string `json:"mail"`
	Nickname string `json:"nickname"`
	Realname string `json:"realname"`
	Tel      string `json:"tel"`
}

// Club is a structure to hold the information about a club.
type Club struct {
	Chief         Chief             `json:"chief"`
	Members       chanjson.ChanJSON `json:"members"`
	Name          string            `json:"name"`
}

type clubMemberResult struct {
	Error error
	Value struct {
		Entrance uint16 `json:"entrance"`
		ID       string `json:"id"`
		Nickname string `json:"nickname"`
		Realname string `json:"realname"`
	}
}

type clubName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type clubNameResult struct {
	Error error
	Value clubName
}

type club struct {
	clubName
	Chief Chief `json:"chief"`
}

type clubResult struct {
	Error error
	Value club
}

// QueryClub returns db.Club corresponding with the given ID.
func (db DB) QueryClub(id string) (Club, error) {
	var chiefID uint16
	var clubID uint8
	var club Club

	if scanError := db.stmts[stmtSelectClub].QueryRow(id).Scan(&chiefID, &clubID, &club.Name); scanError != nil {
		return Club{}, scanError
	}

	if scanError := db.stmts[stmtSelectOfficerMemberInternal].QueryRow(chiefID).Scan(
		&club.Chief.ID, &club.Chief.Mail,
		&club.Chief.Nickname, &club.Chief.Realname,
		&club.Chief.Tel); scanError != nil {
		return Club{}, scanError
	}

	members := make(chan clubMemberResult)

	go func() {
		defer close(members)

		rows, queryError := db.stmts[stmtSelectClubMembers].Query(clubID)
		if queryError != nil {
			members <- clubMemberResult{Error: queryError}
			return
		}

		defer rows.Close()

		for rows.Next() {
			var memberID uint16

			if scanError := rows.Scan(&memberID); scanError != nil {
				members <- clubMemberResult{Error: scanError}
				return
			}

			var result clubMemberResult
			result.Error = db.stmts[stmtSelectBriefMemberInternal].QueryRow(memberID).Scan(
				&result.Value.Entrance, &result.Value.ID,
				&result.Value.Nickname, &result.Value.Realname)

			members <- result

			if result.Error != nil {
				return
			}
		}
	}()

	club.Members = chanjson.New(members)

	return club, nil
}

// QueryClubName returns the name of the club identified with the given ID.
func (db DB) QueryClubName(id string) (string, error) {
	var name string
	scanError := db.stmts[stmtSelectClubName].QueryRow(id).Scan(&name)

	return name, scanError
}

// QueryClubNames returns chanjson.ChanJSON which represents the names of
// all the clubs.
func (db DB) QueryClubNames() chanjson.ChanJSON {
	resultChan := make(chan clubNameResult)

	go func() {
		defer close(resultChan)

		rows, queryError := db.stmts[stmtSelectClubNames].Query()
		if queryError != nil {
			resultChan <- clubNameResult{Error: queryError}
			return
		}

		for rows.Next() {
			var result clubNameResult

			result.Error = rows.Scan(
				&result.Value.ID, &result.Value.Name)

			resultChan <- result

			if result.Error != nil {
				return
			}
		}
	}()

	return chanjson.New(resultChan)
}

// QueryClubs returns chanjson.ChanJSON which represents all the clubs.
func (db DB) QueryClubs() chanjson.ChanJSON {
	resultChan := make(chan clubResult)

	go func() {
		defer close(resultChan)

		clubs, queryError := db.stmts[stmtSelectClubs].Query()
		if queryError != nil {
			resultChan <- clubResult{Error: queryError}
			return
		}

		for clubs.Next() {
			var result clubResult
			var chiefID uint16

			if scanError := clubs.Scan(&chiefID, &result.Value.ID, &result.Value.Name); scanError != nil {
				resultChan <- clubResult{Error: scanError}
				return
			}

			result.Error = db.stmts[stmtSelectOfficerMemberInternal].QueryRow(chiefID).Scan(
				&result.Value.Chief.ID, &result.Value.Chief.Mail,
				&result.Value.Chief.Nickname, &result.Value.Chief.Realname,
				&result.Value.Chief.Tel)

			resultChan <- result
		}
	}()

	return chanjson.New(resultChan)
}
