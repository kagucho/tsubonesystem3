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
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"
	"github.com/kagucho/tsubonesystem3/backend/mail"
	"github.com/kagucho/tsubonesystem3/configuration"
	"github.com/kagucho/tsubonesystem3/json"
	"golang.org/x/crypto/pbkdf2"
	"log"
	"strings"
	"unicode/utf8"
)

type MemberAddress struct {
	Mail     string
	Nickname string
}

// MemberClub is a structure to hold the information about a club which a member
// belongs to.
type MemberClub struct {
	Chief bool   `json:"chief"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

type MemberClubChan <-chan MemberClubResult

// MemberClubResult is a structure to hold the result of an action to query
// the information about a club which a member belongs to.
type MemberClubResult struct {
	Error error
	Value MemberClub
}

type MemberClubIDChan <-chan MemberClubIDResult

type MemberClubIDResult struct {
	Error error
	Value string
}

type MemberCommon struct {
	Affiliation string `json:"affiliation,omitempty"`
	Entrance    uint16 `json:"entrance,omitempty"`
	Nickname    string `json:"nickname"`
	OB          bool   `json:"ob"`
	Realname    string `json:"realname,omitempty"`
}

// MemberDetail is a structure to hold the details of a member.
type MemberDetail struct {
	MemberCommon
	Clubs     MemberClubChan `json:"clubs"`
	Confirmed bool           `json:"confirmed"`
	Gender    string         `json:"gender,omitempty"`
	Mail      string         `json:"mail"`
	Positions PositionChan   `json:"positions"`
	Tel       string         `json:"tel,omitempty"`
}

type MemberEntry struct {
	MemberCommon
	ID string `json:"id"`
}

type MemberEntryChan <-chan MemberEntryResult

type MemberEntryResult struct {
	Error error
	Value MemberEntry
}

// MemberGraph is a structure to hold the information about a member to render
// a graph.
type MemberGraph struct {
	Gender   string
	Nickname string
}

type MemberRole struct {
	Clubs     MemberClubIDChan `json:"clubs"`
	ID        string           `json:"id"`
	Nickname  string           `json:"nickname"`
	OB        bool             `json:"ob"`
}

type MemberRoleChan <-chan MemberRoleResult

type MemberRoleResult struct {
	Error error
	Value MemberRole
}

// MemberStatus is an unsigned integer which describes the acceptable status
// of members for querying.
type MemberStatus uint

// These are flags for MemberStatus.
const (
	MemberStatusOB     MemberStatus = 1 << iota
	MemberStatusActive MemberStatus = 1 << iota
)

type memberAttendance struct {
	party     uint16
	attending bool
}

// Position is a structure to hold the information about a position.
type Position struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PositionChan <-chan PositionResult

// PositionResult is a structure to hold the result of an action to query the
// information about a position where a member is.
type PositionResult struct {
	Error error
	Value Position
}

var MemberInvalidID = errors.New(`invalid id`)
var MemberInvalidMail = errors.New(`invalid mail`)
var MemberInvalidNickname = errors.New(`invalid nickname`)
var MemberInvalidPassword = errors.New(`invalid password`)
var MemberInvalidTel = errors.New(`invalid tel`)

func (entryChan MemberEntryChan) MarshalJSON() ([]byte, error) {
	return json.MarshalChan(entryChan)
}

func (clubChan MemberClubChan) MarshalJSON() ([]byte, error) {
	return json.MarshalChan(clubChan)
}

func (idChan MemberClubIDChan) MarshalJSON() ([]byte, error) {
	return json.MarshalChan(idChan)
}

func (positionChan PositionChan) MarshalJSON() ([]byte, error) {
	return json.MarshalChan(positionChan)
}

func (roleChan MemberRoleChan) MarshalJSON() ([]byte, error) {
	return json.MarshalChan(roleChan)
}

func (db DB) InsertMember(id, mail, nickname string) error {
	if !validateID(id) {
		return MemberInvalidID
	}

	if !validateMail(mail) {
		return MemberInvalidMail
	}

	if !validateNickname(nickname) {
		return MemberInvalidNickname
	}

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

// QueryMemberDetail returns db.MemberDetail of the member identified with the
// given ID.
func (db DB) QueryMemberDetail(id string) (MemberDetail, error) {
	var dbFlags string
	var dbMember uint16
	var output MemberDetail

	if scanError := db.stmts[stmtSelectMember].QueryRow(id).Scan(
		&dbMember, &output.Affiliation, &output.Entrance, &dbFlags,
		&output.Gender, &output.Mail, &output.Nickname, &output.Realname,
		&output.Tel); scanError == sql.ErrNoRows {
		return MemberDetail{}, IncorrectIdentity
	} else if scanError != nil {
		return MemberDetail{}, scanError
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

		rows, queryError := db.stmts[stmtSelectClubsByInternalMember].Query(dbMember)
		if queryError != nil {
			clubs <- MemberClubResult{Error: queryError}
			return
		}

		defer rows.Close()

		for rows.Next() {
			var result MemberClubResult
			var clubChief uint16
			result.Error = rows.Scan(&clubChief, &result.Value.ID, &result.Value.Name)
			if result.Error != nil {
				clubs <- result
				return
			}

			result.Value.Chief = dbMember == clubChief

			clubs <- result
		}
	}()

	positions := make(chan PositionResult)
	output.Positions = positions

	go func() {
		defer close(positions)

		rows, queryError := db.stmts[stmtSelectMemberOfficer].Query(dbMember)
		if queryError != nil {
			positions <- PositionResult{Error: queryError}
			return
		}

		defer rows.Close()

		for rows.Next() {
			var result PositionResult

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
	if scanError == sql.ErrNoRows {
		scanError = IncorrectIdentity
	}

	return graph, scanError
}

func (db DB) QueryMemberMails(ids string) ([]string, error) {
	count := 0
	idBytes := []byte(ids)
	for index, character := range idBytes {
		if character == ' ' {
			idBytes[index] = ','
			count++
		}
	}

	rows, queryError := db.stmts[stmtSelectMemberMails].Query(idBytes)
	if queryError != nil {
		return nil, queryError
	}

	mails := make([]string, 0, count)
	for rows.Next() {
		var mail string
		if scanError := rows.Scan(&mail); scanError != nil {
			return nil, scanError
		}

		if mail != `` {
			mails = append(mails, mail)
		}

		count--
	}

	if count > 0 {
		return nil, IncorrectIdentity
	}

	return mails, nil
}

func (db DB) QueryMemberNickname(id string) (string, error) {
	var nickname string
	scanError := db.stmts[stmtSelectMemberNickname].QueryRow(id).Scan(&nickname)
	if scanError == sql.ErrNoRows {
		scanError = IncorrectIdentity
	}

	return nickname, scanError
}

func (db DB) QueryMemberRoles() MemberRoleChan {
	resultChan := make(chan MemberRoleResult)

	go func() {
		defer close(resultChan)

		rows, queryError := db.stmts[stmtSelectMemberRoles].Query()
		if queryError != nil {
			resultChan <- MemberRoleResult{Error: queryError}
			return
		}

		defer rows.Close()

		for rows.Next() {
			var flags string
			var dbID  uint16
			var result MemberRoleResult

			result.Error = rows.Scan(
				&result.Value.ID, &flags, &dbID,
				&result.Value.Nickname)
			if result.Error != nil {
				resultChan <- result

				return
			}

			clubs := make(chan MemberClubIDResult)

			go func() {
				defer close(clubs)

				rows, queryError := db.stmts[stmtSelectClubIDsByInternalMember].Query(dbID)
				if queryError != nil {
					clubs <- MemberClubIDResult{Error: queryError}

					return
				}

				defer rows.Close()

				for rows.Next() {
					var result MemberClubIDResult

					result.Error = rows.Scan(&result.Value)
					clubs <- result

					if result.Error != nil {
						return
					}
				}
			}()

			result.Value.Clubs = clubs
			result.Value.OB = flagsHasOB(flags)
			resultChan <- result
		}
	}()

	return resultChan
}

func (db DB) QueryMemberTmp(id string) (bool, error) {
	rows, queryError := db.stmts[stmtSelectMemberPassword].Query(id)
	if queryError != nil {
		return false, queryError
	}

	defer func() {
		if closeError := rows.Close(); closeError != nil {
			log.Print(closeError)
		}
	}()

	rows.Next()

	var dbPassword sql.RawBytes
	if scanError := rows.Scan(&dbPassword); scanError != nil {
		return false, IncorrectIdentity
	}

	var result byte

	for _, value := range dbPassword {
		// Compare in a constant time to prevent side channel attaks.
		// Achieve this by eliminating any conditional branches.
		result |= value
	}

	return result == 0, nil
}

// QueryMembers returns db.MemberEntryChan which represents all the members.
func (db DB) QueryMembers() MemberEntryChan {
	resultChan := make(chan MemberEntryResult)

	go func() {
		defer close(resultChan)

		rows, queryError := db.stmts[stmtSelectMembers].Query()
		if queryError != nil {
			resultChan <- MemberEntryResult{Error: queryError}
			return
		}

		defer rows.Close()

		for rows.Next() {
			var flags string
			var result MemberEntryResult

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

	return resultChan
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
		if !validatePassword(password) {
			return MemberInvalidPassword
		}

		dbPassword, hashError := makeDBPassword(password)
		if hashError != nil {
			return hashError
		}

		expressions = append(expressions, `password=?`)
		arguments = append(arguments, dbPassword)
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
		if !validateMail(mail) {
			return MemberInvalidMail
		}

		expressions = append(expressions, `flags=flags&~1, mail=?`)
		arguments = append(arguments, mail)
	}

	if nickname != `` {
		if !validateNickname(nickname) {
			return MemberInvalidNickname
		}

		expressions = append(expressions, `nickname=?`)
		arguments = append(arguments, nickname)
	}

	if realname != `` {
		expressions = append(expressions, `realname=?`)
		arguments = append(arguments, realname)
	}

	if tel != `` {
		if !validateTel(tel) {
			return MemberInvalidTel
		}

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
				panic(recovered)
			}
		}
	}()

	if len(expressions) > 0 {
		arguments = append(arguments, id)

		_, execError := tx.Exec(
			strings.Join([]string{`UPDATE members SET`, strings.Join(expressions, `,`), `WHERE display_id=?`}, ` `),
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

		pendings, pendingsError := db.txQueryInternalClubs(tx, clubs)
		if pendingsError != nil {
			panic(pendingsError)
		}

		toDeletes, diffError := memberDiffClubs(db, tx, dbID, pendings)
		if diffError != nil {
			panic(diffError)
		}

		if len(toDeletes) > 0 {
			bytes := make([]byte, 0, len(toDeletes) * 6)
			index := 0
			for {
				for toDelete := toDeletes[index]; toDelete > 0; toDelete /= 10 {
					bytes = append(bytes, '0' + byte(toDelete) % 10)
				}

				index++
				if index >= len(toDeletes) {
					break
				}

				bytes = append(bytes, ',')
			}

			if _, execError := tx.Stmt(db.stmts[stmtDeleteClubMemberByInternal]).Exec(bytes); execError != nil {
				panic(execError)
			}
		}

		for pending := range pendings {
			if _, execError := tx.Stmt(db.stmts[stmtInsertInternalClubMember]).Exec(pending, dbID); execError != nil {
				panic(execError)
			}
		}
	}

	return nil
}

func (db DB) UpdatePassword(id, currentPassword, newPassword string) error {
	if !validatePassword(newPassword) {
		return MemberInvalidPassword
	}

	if scanError := func() error {
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
			return IncorrectIdentity
		}

		return verifyPassword(currentPassword, dbPassword)
	}(); scanError != nil {
		return scanError
	}

	newDBPassword, hashError := makeDBPassword(newPassword)
	if hashError != nil {
		return hashError
	}

	_, execError := db.stmts[stmtUpdateMemberPassword].Exec(newDBPassword, id)

	return execError
}

func flagsHasOB(flags string) bool {
	for _, flag := range strings.Split(flags, `,`) {
		if flag == `ob` {
			return true
		}
	}

	return false
}

/*
	DRAFT NIST Special Publication 800-63B
	Digital Identity Guidelines
	Authentication and Lifecycle Management
	5. Authenticator and Verifier Requirements
	https://pages.nist.gov/800-63-3/sp800-63b.html#sec5

*/
func hashPassword(raw string, salt []byte) ([]byte, error) {
	combinedSalt := bytes.NewBuffer(make([]byte, 0, len(configuration.DBPasswordSalt) + len(salt)))

	/* 
		> A keyed hash function (e.g., HMAC [FIPS198-1]), with the key
		> stored separately from the hashed authenticators (e.g., in a
		> hardware security module) SHOULD be used to further resist
		> dictionary attacks against the stored hashed authenticators.

		FIXME: Catenation doesn't make sense. We may prepare.hash
		context initialized with configuration.DBPasswordSalt.
		However, Golang doesn't provide any feature to clone hash
		contexts.
	*/
	if _, writeError := combinedSalt.WriteString(configuration.DBPasswordSalt); writeError != nil {
		return nil, writeError
	}

	/*
		> The salt value SHALL be a 32-bit or longer random value
		> generated by an approved random bit generator and stored along
		> with the hash result.

		salt may be stored in the database, but it must be unique for
		each users, which prevents attackers from constructing a
		rainbow table.
	*/
	if _, writeError := combinedSalt.Write(salt); writeError != nil {
		return nil, writeError
	}

	/*
		> Secrets SHALL be hashed with a salt value using an approved
		> hash function such as PBKDF2 as described in [SP 800-132].
		> At least 10,000 iterations of the hash function SHOULD be
		> performed.

		Choose SHA-512 because it could be relatively fast even for
		generic computers with Intel CPU thanks to SHA extensions.

		Intel® SHA Extensions | Intel® Software
		https://software.intel.com/en-us/articles/intel-sha-extensions
	*/
	return pbkdf2.Key([]byte(raw), combinedSalt.Bytes(), 16384, sha512.Size, sha512.New), nil
}

func makeDBPassword(raw string) ([]byte, error) {
	salt := make([]byte, sha512.BlockSize, sha512.BlockSize + sha512.Size)

	if _, randError := rand.Read(salt); randError != nil {
		return nil, randError
	}

	hashed, hashError := hashPassword(raw, salt)
	if hashError != nil {
		return nil, hashError
	}

	buffer := bytes.NewBuffer(salt)

	if _, writeError := buffer.Write(hashed); writeError != nil {
		return nil, writeError
	}

	return buffer.Bytes(), nil
}

func verifyPassword(raw string, db []byte) error {
	hashed, hashError := hashPassword(raw, db[:sha512.BlockSize])
	if hashError != nil {
		return hashError
	}

	if subtle.ConstantTimeCompare(hashed, db[sha512.BlockSize:]) != 1 {
		return IncorrectIdentity
	}

	return nil
}

func memberDiffClubs(db DB, tx *sql.Tx, member uint16, clubs map[uint8]struct{}) ([]uint16, error) {
	deleted := make([]uint16, 0)
	registered, queryError := tx.Stmt(db.stmts[stmtSelectClubInternalByInternalMember]).Query(member)
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

func (db DB) queryMemberInternalIDMails(ids string) ([]uint16, []string, error) {
	count := 0
	idBytes := []byte(ids)
	for index, character := range idBytes {
		if character == ' ' {
			idBytes[index] = ','
			count++
		}
	}

	rows, queryError := db.stmts[stmtSelectMemberInternalIDMails].Query(idBytes)
	if queryError != nil {
		return nil, nil, queryError
	}

	defer func() {
		if closeError := rows.Close(); closeError != nil {
			log.Print(closeError)
		}
	}()

	internalIDs := make([]uint16, 0, count)
	mails := make([]string, 0, count)

	for rows.Next() {
		var id uint16
		var mail string
		if scanError := rows.Scan(&id, &mail); scanError != nil {
			return nil, nil, scanError
		}

		internalIDs = append(internalIDs, id)
		if mail != `` {
			mails = append(mails, mail)
		}
	}

	if len(internalIDs) < count {
		return nil, nil, IncorrectIdentity
	}

	return internalIDs, mails, nil
}

func validateID(id string) bool {
	for index := 0; index < len(id); index++ {
		/*
			URL Standard
			5.2. application/x-www-form-urlencoded serializing
			https://url.spec.whatwg.org/#urlencoded-serializing

			> 0x2A
			> 0x2D
			> 0x2E
			> 0x30 to 0x39
			> 0x41 to 0x5A
			> 0x5F
			> 0x61 to 0x7A
			>
			> Append a code point whose value is byte to output.

			Accept only those characters.
		*/
		if !(id[index] == 0x2A || id[index] == 0x2D || id[index] == 0x2E ||
			(id[index] >= 0x30 && id[index] <= 0x39) ||
			(id[index] >= 0x41 && id[index] <= 0x5A) ||
			id[index] == 0x5F ||
			(id[index] >= 0x61 && id[index] <= 0x7A)) {
			return false
		}
	}

	return true
}

func validateMail(dbMail string) bool {
	return mail.ValidateAddressHTML(dbMail)
}

func validateNickname(nickname string) bool {
	index := 0
	for index < len(nickname) {
		if (nickname[index] & 0x80 == 0) {
			if (nickname[index] < 0x20 || nickname[index] == '"') {
				return false
			}

			index++
		} else {
			decoded, encodedLen := utf8.DecodeRuneInString(nickname[index:])
			if decoded == utf8.RuneError {
				return false
			}

			index += encodedLen
		}
	}

	return true
}

func validatePassword(password string) bool {
	if len(password) > sha512.BlockSize {
		return false
	}

	for index := 0; index < len(password); index++ {
		// Accept only the ASCII printable characters.
		if password[index] < 0x20 || password[index] >= 0x80 {
			return false
		}
	}

	return true
}

func validateTel(tel string) bool {
	for index := 0; index < len(tel); index++ {
		/*
			RFC 3986 - Uniform Resource Identifier (URI): Generic Syntax
			https://tools.ietf.org/html/rfc3986#section-2
			2.2.  Reserved Characters

			Allow characters valid in hier-part.
		*/
		if !(tel[index] == 0x21 || tel[index] == 0x24 ||
			(tel[index] >= 0x26 && tel[index] <= 0x39) ||
			tel[index] == 0x3B || tel[index] == 0x3D ||
			(tel[index] >= 0x41 && tel[index] <= 0x5A) ||
			tel[index] == 0x5F ||
			(tel[index] >= 0x61 && tel[index] <= 0x7A) ||
			tel[index] == 0x7E) {
			return false
		}
	}

	return true
}
