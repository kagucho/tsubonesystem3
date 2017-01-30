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
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"github.com/kagucho/tsubonesystem3/chanjson"
	"github.com/kagucho/tsubonesystem3/configuration"
	"log"
	"runtime/debug"
	"strings"
)

// MemberGraph is a structure to hold the information about a member to render
// a graph.
type MemberGraph struct {
	Gender   string
	Nickname string
}

// Position is a structure to hold the information about a position.
type Position struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// MemberClub is a structure to hold the information about a club which a member
// belongs to.
type MemberClub struct {
	Chief bool   `json:"chief"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

// MemberClubResult is a structure to hold the result of an action to query
// the information about a club which a member belongs to.
type MemberClubResult struct {
	Error error
	Value MemberClub
}

// MemberPositionResult is a structure to hold the result of an action to query
// the information about a position where a member is.
type MemberPositionResult struct {
	Error error
	Value Position
}

// MemberStatus is an unsigned integer which describes the acceptable status
// of members for querying.
type MemberStatus uint

// These are flags for MemberStatus.
const (
	MemberStatusOB     MemberStatus = 1 << iota
	MemberStatusActive MemberStatus = 1 << iota
)

// Member is a structure to hold the details of a member.
type Member struct {
	Affiliation string
	Clubs       <-chan MemberClubResult
	Confirmed   bool
	Entrance    uint16
	Gender      string
	Mail        string
	Nickname    string
	OB          bool
	Positions   <-chan MemberPositionResult
	Realname    string
	Tel         string
}

type member struct {
	Affiliation string `json:"affiliation,omitempty"`
	Entrance    uint16 `json:"entrance,omitempty"`
	ID          string `json:"id"`
	Nickname    string `json:"nickname"`
	OB          bool   `json:"ob"`
	Realname    string `json:"realname,omitempty"`
}

type memberResult struct {
	Error error
	Value member
}

func flagsHasOB(flags string) bool {
	for _, flag := range strings.Split(flags, `,`) {
		if flag == `ob` {
			return true
		}
	}

	return false
}

func hashPassword(password string) ([]byte, error) {
	hash := hmac.New(sha256.New224, []byte(configuration.DBPasswordKey))

	if _, writeError := hash.Write([]byte(password)); writeError != nil {
		return nil, writeError
	}

	return hash.Sum(nil), nil
}

func memberQueryClubIDs(tx *sql.Tx, clubs []string) (map[uint8]struct{}, error) {
	clubInterfaces := make([]interface{}, len(clubs))
	clubIDs := make(map[uint8]struct{}, len(clubs))

	for index, value := range clubs {
		clubInterfaces[index] = value
	}

	rows, queryError := tx.Query(
		strings.Join([]string{
			`SELECT id FROM clubs WHERE display_id IN(?`,
			strings.Repeat(`,?`, len(clubs)-1), `)`,
		}, ``),
		clubInterfaces...)
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

func memberDiffClubs(db DB, tx *sql.Tx, member uint16, clubs map[uint8]struct{}) ([]interface{}, error) {
	deleted := make([]interface{}, 0)
	registered, queryError := tx.Stmt(db.stmts[stmtSelectMemberClubIDs]).Query(member)
	if queryError != nil {
		return nil, queryError
	}

	defer registered.Close()

	for registered.Next() {
		var id uint16
		var club uint8
		if scanError := registered.Scan(&id, &club); scanError != nil {
			return nil, scanError
		}

		if _, present := clubs[club]; !present {
			deleted = append(deleted, id)
		} else {
			delete(clubs, club)
		}
	}

	return deleted, nil
}

func (db DB) InsertMember(id, mail, nickname string) error {
	_, execError := db.stmts[stmtInsertMember].Exec(id, mail, nickname)
	return execError
}

func (db DB) DeclareMemberOB(id string) error {
	_, execError := db.stmts[stmtDeclareMemberOB].Exec(id)
	return execError
}

func (db DB) ConfirmMember(id string) error {
	_, execError := db.stmts[stmtConfirmMember].Exec(id)
	return execError
}

func (db DB) DeleteMember(id string) error {
	_, execError := db.stmts[stmtDeleteMember].Exec(id)
	return execError
}

// QueryMember returns db.MemberDetail of the member identified with the given
// ID.
func (db DB) QueryMember(id string) (Member, error) {
	var dbFlags string
	var dbMember uint16
	var output Member

	if scanError := db.stmts[stmtSelectMember].QueryRow(id).Scan(
		&dbMember, &output.Affiliation, &output.Entrance, &dbFlags,
		&output.Gender, &output.Mail, &output.Nickname, &output.Realname,
		&output.Tel); scanError != nil {
		return Member{}, scanError
	}

	for _, flag := range strings.Split(dbFlags, `,`) {
		switch flag {
		case `confirmed`:
			output.Confirmed = true

		case `ob`:
			output.OB = true
		}
	}

	output.OB = flagsHasOB(dbFlags)

	clubs := make(chan MemberClubResult)
	output.Clubs = clubs

	go func() {
		defer close(clubs)

		rows, queryError := db.stmts[stmtSelectMemberClubs].Query(dbMember)
		if queryError != nil {
			clubs <- MemberClubResult{Error: queryError}
			return
		}

		defer rows.Close()

		for rows.Next() {
			var dbClub uint8

			if scanError := rows.Scan(&dbClub); scanError != nil {
				clubs <- MemberClubResult{Error: scanError}
				return
			}

			var result MemberClubResult
			var clubChief uint16
			result.Error = db.stmts[stmtSelectClubInternal].QueryRow(dbClub).Scan(
				&clubChief, &result.Value.ID, &result.Value.Name)
			if result.Error != nil {
				clubs <- result
				return
			}

			result.Value.Chief = dbMember == clubChief

			clubs <- result
		}
	}()

	positions := make(chan MemberPositionResult)
	output.Positions = positions

	go func() {
		defer close(positions)

		rows, queryError := db.stmts[stmtSelectMemberOfficer].Query(dbMember)
		if queryError != nil {
			positions <- MemberPositionResult{Error: queryError}
			return
		}

		defer rows.Close()

		for rows.Next() {
			var result MemberPositionResult

			result.Error = rows.Scan(&result.Value.ID, &result.Value.Name)
			positions <- result

			if result.Error != nil {
				return
			}
		}
	}()

	return output, nil
}

// QueryMemberGraph returns db.MemberGraph of the member identified with the
// given ID.
func (db DB) QueryMemberGraph(id string) (MemberGraph, error) {
	var graph MemberGraph

	scanError := db.stmts[stmtSelectMemberGraph].QueryRow(id).Scan(
		&graph.Gender, &graph.Nickname)

	return graph, scanError
}

// QueryMembers returns chanjson.ChanJSON which represents all the members.
func (db DB) QueryMembers() chanjson.ChanJSON {
	resultChan := make(chan memberResult)

	go func() {
		defer close(resultChan)

		rows, queryError := db.stmts[stmtSelectMembers].Query()
		if queryError != nil {
			resultChan <- memberResult{Error: queryError}
			return
		}

		defer rows.Close()

		for rows.Next() {
			var flags string
			var result memberResult

			result.Error = rows.Scan(
				&result.Value.Affiliation,
				&result.Value.ID,
				&result.Value.Entrance,
				&flags,
				&result.Value.Nickname,
				&result.Value.Realname)

			if result.Error == nil {
				result.Value.OB = flagsHasOB(flags)
				resultChan <- result
			} else {
				resultChan <- result

				return
			}
		}
	}()

	return chanjson.New(resultChan)
}

// QueryMembersCount returns the number of the members who matches the given
// conditions.
func (db DB) QueryMembersCount(entrance int, nickname string, realname string,
	status MemberStatus) (uint16, error) {
	pattern := func(raw string) string {
		return strings.Join(
			[]string{
				`%`,
				strings.Replace(
					strings.Replace(
						strings.Replace(
							raw,
							`\`, `\\`, -1),
						`%`, `\%`, -1),
					`_`, `\_`, -1),
				`%`,
			}, ``)
	}

	arguments := make([]interface{}, 0, 4)
	arguments = append(arguments, pattern(nickname))
	arguments = append(arguments, pattern(realname))

	if entrance == 0 {
		arguments = append(arguments, nil)
	} else {
		arguments = append(arguments, entrance)
	}

	switch status {
	case 0:
		return 0, nil

	case MemberStatusOB:
		arguments = append(arguments, 1)

	case MemberStatusActive:
		arguments = append(arguments, 0)

	case MemberStatusOB | MemberStatusActive:
		arguments = append(arguments, nil)

	default:
		return 0, fmt.Errorf(`invalid status %v`, status)
	}

	var count uint16

	scanError := db.stmts[stmtCountMembers].QueryRow(arguments...).Scan(&count)

	return count, scanError
}

func (db DB) UpdateMember(id, password, affiliation string, clubs []string, entrance int, gender, mail, nickname, realname, tel string) (returning error) {
	expressions := make([]string, 0, 5)
	arguments := make([]interface{}, 0, 5)

	if password != `` {
		hashedPassword, hashError := hashPassword(password)
		if hashError != nil {
			return hashError
		}

		expressions = append(expressions, `password=?`)
		arguments = append(arguments, hashedPassword)
	}

	if affiliation != `` {
		expressions = append(expressions, `affiliation=?`)
		arguments = append(arguments, affiliation)
	}

	if entrance != 0 {
		expressions = append(expressions, `entrance=?`)
		arguments = append(arguments, entrance)
	}

	if gender != `` {
		expressions = append(expressions, `gender=?`)
		arguments = append(arguments, gender)
	}

	if mail != `` {
		expressions = append(expressions, `flags=EXPORT_SET(flags, '', 'confirmed'), mail=?`)
		arguments = append(arguments, mail)
	}

	if nickname != `` {
		expressions = append(expressions, `nickname=?`)
		arguments = append(arguments, nickname)
	}

	if realname != `` {
		expressions = append(expressions, `realname=?`)
		arguments = append(arguments, realname)
	}

	if tel != `` {
		expressions = append(expressions, `tel=?`)
		arguments = append(arguments, tel)
	}

	tx, txError := db.sql.Begin()
	if txError != nil {
		return txError
	}

	defer func() {
		if recovered := recover(); recovered == nil {
			returning = tx.Commit()
		} else {
			if rollbackError := tx.Rollback(); rollbackError != nil {
				log.Print(rollbackError)
			}

			var ok bool
			returning, ok = recovered.(error)
			if !ok {
				recoveredString := fmt.Sprint(recovered)

				log.Print(recoveredString)
				debug.PrintStack()

				returning = errors.New(recoveredString)
			}
		}
	}()

	if len(expressions) > 0 {
		arguments = append(arguments, id)

		_, execError := tx.Exec(
			strings.Join([]string{`UPDATE members SET`, strings.Join(expressions, `,`), `WHERE id=?`}, ` `),
			arguments...)
		if execError != nil {
			panic(execError)
		}
	}

	if len(clubs) > 0 {
		var dbID uint16
		if scanError := tx.Stmt(db.stmts[stmtSelectMemberID]).QueryRow(id).Scan(&dbID); scanError != nil {
			panic(scanError)
		}

		pendings, pendingsError := memberQueryClubIDs(tx, clubs)
		if pendingsError != nil {
			panic(pendingsError)
		}

		toDelete, diffError := memberDiffClubs(db, tx, dbID, pendings)
		if diffError != nil {
			panic(diffError)
		}

		if len(toDelete) > 0 {
			if _, execError := db.sql.Exec(
				strings.Join([]string{
					`DELETE FROM club_member WHERE id IN(?`,
					strings.Repeat(`,?`, len(toDelete)-1), `)`,
				}, ``),
				toDelete...); execError != nil {
				panic(execError)
			}
		}

		for pending := range pendings {
			if _, execError := tx.Stmt(db.stmts[stmtInsertClubMemberInternal]).Exec(pending, dbID); execError != nil {
				panic(execError)
			}
		}
	}

	return nil
}

func (db DB) UpdatePassword(id, oldPassword, newPassword string) error {
	if scanError := func() error {
		hashedOldPassword, hashError := hashPassword(oldPassword)
		if hashError != nil {
			return hashError
		}

		rows, queryError := db.stmts[stmtSelectMemberPassword].Query(id)
		if queryError != nil {
			return queryError
		}

		defer func() {
			if closeError := rows.Close(); closeError != nil {
				log.Print(closeError)
			}
		}()

		rows.Next()

		var dbPassword sql.RawBytes
		if scanError := rows.Scan(&dbPassword); scanError != nil {
			return sql.ErrNoRows
		}

		if !hmac.Equal(dbPassword, hashedOldPassword) {
			return errors.New(`incorrect password`)
		}

		return nil
	}(); scanError != nil {
		return scanError
	}

	hashedNewPassword, hashError := hashPassword(newPassword)
	if hashError != nil {
		return hashError
	}

	_, execError := db.stmts[stmtUpdatePassword].Exec(hashedNewPassword, id)

	return execError
}

// ValidateMemberEntrance returns whether the given entrance year is valid.
func ValidateMemberEntrance(entrance int) bool {
	return entrance >= 1901 && entrance <= 2155
}

func ValidatePassword(password string) bool {
	return len(password) <= sha256.Size224
}
