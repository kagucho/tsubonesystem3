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
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/kagucho/tsubonesystem3/backend/encoding"
	"github.com/kagucho/tsubonesystem3/backend/mail"
	"github.com/kagucho/tsubonesystem3/configuration"
	"golang.org/x/crypto/pbkdf2"
	"log"
	"strings"
	"sync"
	"unicode/utf8"
)

/*
MemberClub is a structure holding the information about a club which a member
belongs to.
*/
type MemberClub struct {
	Chief bool   `json:"chief"`
	ID    string `json:"id"`
}

/*
MemberClubResult is a structure representing a result of querying
db.MemberClub.
*/
type MemberClubResult struct {
	MemberClub
	Error error
}

// MemberClubChan is a reciever of MemberClubResult.
type MemberClubChan <-chan MemberClubResult

// MemberClubQuerier is a structure to provide a feature to query db.MemberClub.
type MemberClubQuerier struct {
	context *querierContext
	stmt    *sql.Stmt
}

/*
MemberCommon is a structure holding information about a member common for
queries.
*/
type MemberCommon struct {
	Affiliation encoding.ZeroString `json:"affiliation"`
	Entrance    encoding.ZeroUint16 `json:"entrance"`
	Nickname    string              `json:"nickname"`
	OB          bool                `json:"ob"`
	Realname    encoding.ZeroString `json:"realname"`
}

// MemberDetail is a structure holding details of a member.
type MemberDetail struct {
	MemberCommon
	Clubs     MemberClubQuerier   `json:"clubs"`
	Confirmed bool                `json:"confirmed"`
	Gender    encoding.ZeroString `json:"gender"`
	Mail      string              `json:"mail"`
	Positions PositionQuerier     `json:"positions"`
	Tel       encoding.ZeroString `json:"tel"`
}

// MemberEntry is a structure holding the basic information of a member.
type MemberEntry struct {
	MemberCommon
	ID string `json:"id"`
}

/*
MemberEntryResult is a structure representing a result of querying
db.MemberEntry.
*/
type MemberEntryResult struct {
	MemberEntry
	Error error
}

// MemberEntryChan is a reciever of db.MemberEntryResult.
type MemberEntryChan <-chan MemberEntryResult

/*
MemberGraph is a structure holding the information about a member to render a
graph.
*/
type MemberGraph struct {
	Gender   string
	Nickname string
}

/*
MemberMailResult is a structure representing a result of querying the ID
and email address of a member.
*/
type MemberMailResult struct {
	Error error
	ID    string
	Mail  string
}

// MemberMailChan is a reciever of db.MemberMailResult.
type MemberMailChan <-chan MemberMailResult

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

// PositionResult is a structure holding the result of an action to query the
// information about a position where a member is.
type PositionResult struct {
	ID string
	Error error
}

// PositionChan is a reciver of db.PositionResult.
type PositionChan <-chan PositionResult

// PositionQuerier is a structure to provide a feature to query db.Position.
type PositionQuerier struct {
	context *querierContext
	stmt    *sql.Stmt
}

type querierContext struct {
	member uint16
	mutex  sync.Mutex
	tx     *sql.Tx
}

/*
ErrMemberIsOfficer is an error telling the member is an officer, which is
unexpected.
*/
var ErrMemberIsOfficer = errors.New(`member is officer`)

/*
MarshalJSON returns the JSON encoding of the remaining entries and closes the
channel.

This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (entryChan MemberEntryChan) MarshalJSON() ([]byte, error) {
	return encoding.MarshalJSONArray(func() (interface{}, error, bool) {
		result, present := <-entryChan
		return result.MemberEntry, result.Error, present
	})
}

/*
MarshalJSON returns the JSON encoding of the remaining clubs a member belonging
and closes the channel.

This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (clubChan MemberClubChan) MarshalJSON() ([]byte, error) {
	return encoding.MarshalJSONArray(func() (interface{}, error, bool) {
		result, present := <-clubChan
		return result.MemberClub, result.Error, present
	})
}

/*
MarshalJSON returns the JSON encoding of the remaining pairs of ID and email
address, and closes the channel.

This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (mailChan MemberMailChan) MarshalJSON() ([]byte, error) {
	return encoding.MarshalJSONObject(func() (string, interface{}, error, bool) {
		var value interface{}

		result, present := <-mailChan
		value = result.Mail
		if value == `` {
			value = nil
		}

		return string(result.ID), value, result.Error, present
	})
}

/*
MarshalJSON returns the JSON encoding of the remaining positions and closes the
channel.

This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (positionChan PositionChan) MarshalJSON() ([]byte, error) {
	return encoding.MarshalJSONArray(func() (interface{}, error, bool) {
		result, present := <-positionChan
		return result.ID, result.Error, present
	})
}

/*
InsertMember inserts a member with the given properties.

It may return one of the following errors:
db.ErrBadOmission tells the ID, email address, or nickname is omitted.
db.ErrDupEntry tells the ID is duplicate.
db.ErrInvalid tells some of the properties is invalid.

Other errors tell db.DB is bad.
*/
func (db DB) InsertMember(id, mail, nickname string) error {
	if id == `` || mail == `` || nickname == `` {
		return ErrBadOmission
	}

	if !validateID(id) || !validateMemberMail(mail) || !validateMemberNickname(nickname) {
		return ErrInvalid
	}

	_, err := db.stmts[stmtInsertMember].Exec(id, mail, nickname)

	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		switch mysqlErr.Number {
		case erDataTooLong:
			fallthrough
		case erTruncatedWrongValueForField:
			return ErrInvalid

		case erDupEntry:
			return ErrDupEntry
		}
	}

	return err
}

/*
DeclareMemberOB declares the memeber identified by the given ID is now an OB.

It returns db.ErrIncorrectIdentity if the ID is incorrect. Other errors tell
db.DB is bad.
*/
func (db DB) DeclareMemberOB(id string) error {
	result, execErr := db.stmts[stmtDeclareMemberOB].Exec(id)
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
ConfirmMember confirms the email address of the member identified by the given
ID.

It returns db.ErrIncorrectIdentity if the ID is incorrect. Other errors tell
db.DB is bad.
*/
func (db DB) ConfirmMember(id string) error {
	result, execErr := db.stmts[stmtConfirmMember].Exec(id)
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
DeleteMember deletes the member identified by the given ID.

It returns db.ErrMemberIsOfficer if the member is an officer. It returns
db.ErrIncorrectIdentity if the ID is incorrect. Other errors tell db.DB is bad.
*/
func (db DB) DeleteMember(id string) error {
	result, execErr := db.stmts[stmtDeleteMember].Exec(id)
	if execErr != nil {
		if mysqlErr, ok := execErr.(*mysql.MySQLError); ok && mysqlErr.Number == erRowIsReferenced || mysqlErr.Number == erRowIsReferenced2 {
			return ErrMemberIsOfficer
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
QueryMemberDetail returns db.MemberDetail of the member identified with the
given ID.

It returns db.ErrIncorrectIdentity if the ID is incorrect. Other errors tell
db.DB is bad.

tl;dr:
If you need Query of Clubs and/or Positions, use them and call End as soon as
possible! For channels returned by them: use them and drain atomically!

Long explanation:
Any other transactions will be blocked until End function of db.MemberDetail
gets called. Clubs and Positions of db.MemberDetail will be invalid after that.
Channels returned by Query or MarshalJSON of Clubs and Positions will block each
other. For example, if you call Query of Clubs while a channel returned by
another Query is left open, it will result in dead blocking.
*/
func (db DB) QueryMemberDetail(id string) (MemberDetail, error) {
	var dbFlags string
	var dbMember uint16
	var output MemberDetail

	tx, err := db.sql.BeginTx(context.Background(),
		&sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  true,
		})
	if err != nil {
		return output, err
	}

	if err := tx.Stmt(db.stmts[stmtSelectMemberByID]).QueryRow(id).Scan(
		&dbMember, (*string)(&output.Affiliation),
		(*uint16)(&output.Entrance), &dbFlags,
		(*string)(&output.Gender), &output.Mail,
		&output.Nickname, (*string)(&output.Realname),
		(*string)(&output.Tel)); err != nil {
		if err := tx.Commit(); err != nil {
			log.Print(err)
		}

		if err == sql.ErrNoRows {
			err = ErrIncorrectIdentity
		}

		return output, err
	}

	for _, flag := range strings.Split(dbFlags, `,`) {
		switch flag {
		case `confirmed`:
			output.Confirmed = true

		case `ob`:
			output.OB = true
		}
	}

	querier := querierContext{tx: tx, member: dbMember}
	output.Clubs = MemberClubQuerier{&querier, db.stmts[stmtSelectClubsByInternalMember]}
	output.Positions = PositionQuerier{&querier, db.stmts[stmtSelectOfficerIDByMemberID]}

	return output, nil
}

/*
QueryMemberGraph returns db.MemberGraph of the member identified with the given
ID.

It returns db.ErrIncorrectIdentity if the given ID is incorrect. Other errors
tell db.DB is bad.
*/
func (db DB) QueryMemberGraph(id string) (MemberGraph, error) {
	var graph MemberGraph

	err := db.stmts[stmtSelectMemberGraphByID].QueryRow(id).Scan(
		(*string)(&graph.Gender), &graph.Nickname)
	if err == sql.ErrNoRows {
		err = ErrIncorrectIdentity
	}

	return graph, err
}

/*
QueryMemberMails returns db.MemberMailChan representing the email addresses of
all members.

Resources will be holded until the channel gets closed.
*/
func (db DB) QueryMemberMails() MemberMailChan {
	resultChan := make(chan MemberMailResult)

	go func() {
		defer close(resultChan)

		rows, err := db.stmts[stmtSelectMemberIDMails].Query()
		if err != nil {
			resultChan <- MemberMailResult{Error: err}
			return
		}

		defer func() {
			if err := rows.Close(); err != nil {
				log.Print(err)
			}
		}()

		for rows.Next() {
			var result MemberMailResult

			result.Error = rows.Scan(&result.ID, &result.Mail)
			resultChan <- result
			if result.Error != nil {
				return
			}
		}
	}()

	return resultChan
}

/*
QueryMemberNickname returns the nickname of the member identified by the given
ID.

It returns ErrIncorrectIdentity if the ID is incorrect. Other errors tell db.DB
is bad.
*/
func (db DB) QueryMemberNickname(id string) (string, error) {
	var nickname string

	err := db.stmts[stmtSelectMemberNicknameByID].QueryRow(id).Scan(&nickname)
	if err == sql.ErrNoRows {
		err = ErrIncorrectIdentity
	}

	return nickname, err
}

/*
QueryMemberTmp returns a Boolean telling whether the member identified by the
given ID has not completed his registration.

It returns db.ErrIncorrectIdentity if the ID is incorrect. Other errors tell
db.DB is bad.
*/
func (db DB) QueryMemberTmp(id string) (bool, error) {
	rows, err := db.stmts[stmtSelectMemberPasswordByID].Query(id)
	if err != nil {
		return false, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Print(err)
		}
	}()

	rows.Next()

	var dbPassword sql.RawBytes
	if err := rows.Scan(&dbPassword); err != nil {
		return false, ErrIncorrectIdentity
	}

	var result byte

	for _, value := range dbPassword {
		// Compare in a constant time to prevent side channel attaks.
		// Achieve this by eliminating any conditional branches.
		result |= value
	}

	return result == 0, nil
}

/*
QueryMembers returns db.MemberEntryChan which represents all members.

Resources will be holded until the channel gets closed.
*/
func (db DB) QueryMembers() MemberEntryChan {
	resultChan := make(chan MemberEntryResult)

	go func() {
		defer close(resultChan)

		rows, err := db.stmts[stmtSelectMembers].Query()
		if err != nil {
			resultChan <- MemberEntryResult{Error: err}
			return
		}

		defer rows.Close()

		for rows.Next() {
			var flags string
			var result MemberEntryResult

			result.Error = rows.Scan(
				(*string)(&result.Affiliation), &result.ID,
				(*uint16)(&result.Entrance), &flags,
				&result.Nickname, (*string)(&result.Realname))

			if result.Error == nil {
				result.OB = flagsHasOB(flags)
				resultChan <- result
			} else {
				resultChan <- result

				return
			}
		}
	}()

	return resultChan
}

/*
QueryMembersCount returns the number of the members who matches the given
conditions.

It returns db.ErrInvalid if the given status is invalid. Other errors tell db.DB
is bad.
*/
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

	arguments := append(make([]interface{}, 0, 4),
		pattern(nickname), pattern(realname))

	if entrance == 0 {
		arguments = append(arguments, nil)
	} else {
		arguments = append(arguments, entrance)
	}

	switch status {
	case 0:
		return 0, nil

	case MemberStatusOB:
		arguments = append(arguments, 2)

	case MemberStatusActive:
		arguments = append(arguments, 0)

	case MemberStatusOB | MemberStatusActive:
		arguments = append(arguments, nil)

	default:
		return 0, ErrInvalid
	}

	var count uint16

	err := db.stmts[stmtCountMembers].QueryRow(arguments...).Scan(&count)

	return count, err
}

/*
UpdateMember updates a member identified with the given ID, with the given
properties.

It returns db.ErrIncorrectIdentity if the ID is incorrect. It returns
db.ErrInvalid if some of the properties is invalid. Other errors tell db.DB is
bad.
*/
func (db DB) UpdateMember(id string, confirm, ob bool, password, affiliation, clubs string, entrance int, gender, mail, nickname, realname, tel string) error {
	andMask := 3
	orMask := 0
	arguments := make([]interface{}, 13)
	arguments[0] = id

	if confirm {
		orMask |= 2
	}

	if ob {
		orMask |= 1
	}

	if password != `` {
		if !validateMemberPassword(password) {
			return ErrInvalid
		}

		dbPassword, err := makeDBPassword(password)
		if err != nil {
			return err
		}

		arguments[3] = dbPassword
	}

	if affiliation != `` {
		arguments[4] = affiliation
	}

	if clubs != `` {
		arguments[6], arguments[5] = stringListToDBList(clubs)
	}

	if entrance != 0 {
		arguments[7] = entrance
	}

	if gender != `` {
		arguments[8] = gender
	}

	if mail != `` {
		if !validateMemberMail(mail) {
			return ErrInvalid
		}

		andMask = 1
		orMask &= andMask
		arguments[9] = mail
	}

	if nickname != `` {
		if !validateMemberNickname(nickname) {
			return ErrInvalid
		}

		arguments[10] = nickname
	}

	if realname != `` {
		arguments[11] = realname
	}

	if tel != `` {
		arguments[12] = tel
	}

	arguments[1] = andMask
	arguments[2] = orMask

	log.Print(arguments)
	_, err := db.stmts[stmtCallUpdateMember].Exec(arguments...)
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
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

	return err
}

/*
UpdatePassword updates.the password of the member identified by the given ID.

It returns db.ErrIncorrectIdentity if the ID or the given current password is
incorrect. It returns db.ErrInvalid if the given new password is invalid. Other
errors tell db.DB is bad.
*/
func (db DB) UpdatePassword(id, currentPassword, newPassword string) error {
	if !validateMemberPassword(newPassword) {
		return ErrInvalid
	}

	if err := func() error {
		rows, err := db.stmts[stmtSelectMemberPasswordByID].Query(id)
		if err != nil {
			return err
		}

		defer func() {
			if err := rows.Close(); err != nil {
				log.Print(err)
			}
		}()

		rows.Next()

		var dbPassword sql.RawBytes
		if err := rows.Scan(&dbPassword); err != nil {
			return ErrIncorrectIdentity
		}

		return verifyPassword(currentPassword, dbPassword)
	}(); err != nil {
		return err
	}

	newDBPassword, hashErr := makeDBPassword(newPassword)
	if hashErr != nil {
		return hashErr
	}

	result, execErr := db.stmts[stmtUpdateMemberPassword].Exec(newDBPassword, id)
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
Query returns db.MemberClubChan representing the clubs a member belonging to.

It will blocks any other operations sharing the context or the database until
the channel gets closed.
*/
func (querier MemberClubQuerier) Query() MemberClubChan {
	clubs := make(chan MemberClubResult)

	go func() {
		defer close(clubs)

		querier.context.mutex.Lock()
		defer querier.context.mutex.Unlock()

		rows, err := querier.context.tx.Stmt(querier.stmt).Query(querier.context.member)
		if err != nil {
			clubs <- MemberClubResult{Error: err}
			return
		}

		defer func() {
			if err := rows.Close(); err != nil {
				log.Print(err)
			}
		}()

		for rows.Next() {
			var result MemberClubResult
			var clubChief uint16
			result.Error = rows.Scan(&clubChief, &result.ID)
			if result.Error != nil {
				clubs <- result
				return
			}

			result.Chief = querier.context.member == clubChief

			clubs <- result
		}
	}()

	return clubs
}

/*
MarshalJSON returns the JSON encoding of the clubs of a member.

This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (querier MemberClubQuerier) MarshalJSON() ([]byte, error) {
	return querier.Query().MarshalJSON()
}

// End ends the transaction and releases all resources.
func (detail MemberDetail) End() error {
	return detail.Clubs.context.tx.Commit()
}

/*
Query returns db.PositionChan representing the positions of a member.

It will blocks any other operations sharing the context or the database until
the channel gets closed.
*/
func (querier PositionQuerier) Query() PositionChan {
	positions := make(chan PositionResult)

	go func() {
		defer close(positions)

		querier.context.mutex.Lock()
		defer querier.context.mutex.Unlock()

		rows, err := querier.context.tx.Stmt(querier.stmt).Query(querier.context.member)
		if err != nil {
			positions <- PositionResult{Error: err}
			return
		}

		defer func() {
			if err := rows.Close(); err != nil {
				log.Print(err)
			}
		}()

		for rows.Next() {
			var result PositionResult

			result.Error = rows.Scan(&result.ID)
			positions <- result

			if result.Error != nil {
				return
			}
		}
	}()

	return positions
}

/*
MarshalJSON returns the JSON encoding of the positions of a member.

This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (querier PositionQuerier) MarshalJSON() ([]byte, error) {
	return querier.Query().MarshalJSON()
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
	RFC 5802 - Salted Challenge Response Authentication Mechanism (SCRAM) SASL and GSS-API Mechanisms
	3.  SCRAM Algorithm Overview
	https://tools.ietf.org/html/rfc5802#section-3
	Designed to be compatible with ClientKey for future extensions.

*/
func hashPassword(raw string, salt, prefix []byte) ([]byte, error) {
	/*
		> SaltedPassword  := Hi(Normalize(password), salt, i)

		DRAFT NIST Special Publication 800-63B
		Digital Identity Guidelines
		Authentication and Lifecycle Management
		5. Authenticator and Verifier Requirements
		https://pages.nist.gov/800-63-3/sp800-63b.html#sec5
		> Secrets SHALL be hashed with a salt value using an approved
		> hash function such as PBKDF2 as described in [SP 800-132].
		> At least 10,000 iterations of the hash function SHOULD be
		> performed.

		Choose SHA-512 because it could be relatively fast even for
		generic computers with Intel CPU thanks to SHA extensions.

		Intel® SHA Extensions | Intel® Software
		https://software.intel.com/en-us/articles/intel-sha-extensions
	*/
	saltedPassword := pbkdf2.Key([]byte(raw), salt, 16384, sha512.Size, sha512.New)

	/*
		RFC 5802 - Salted Challenge Response Authentication Mechanism (SCRAM) SASL and GSS-API Mechanisms
		3.  SCRAM Algorithm Overview
		https://tools.ietf.org/html/rfc5802#section-3
		> ClientKey       := HMAC(SaltedPassword, "Client Key")
	*/
	hmacSHA512 := hmac.New(sha512.New, []byte(configuration.DBPasswordKey))

	if _, err := hmacSHA512.Write(saltedPassword); err != nil {
		return nil, err
	}

	return hmacSHA512.Sum(prefix), nil
}

func makeDBPassword(raw string) ([]byte, error) {
	salt := make([]byte, sha512.BlockSize, sha512.BlockSize+sha512.Size)

	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	return hashPassword(raw, salt, salt)
}

func verifyPassword(raw string, db []byte) error {
	hashed, err := hashPassword(raw, db[:sha512.BlockSize], nil)
	if err != nil {
		return err
	}

	if !hmac.Equal(hashed, db[sha512.BlockSize:]) {
		return ErrIncorrectIdentity
	}

	return nil
}

func validateMemberMail(dbMail string) bool {
	return mail.ValidateAddress(dbMail)
}

func validateMemberNickname(nickname string) bool {
	index := 0
	for index < len(nickname) {
		if nickname[index]&0x80 == 0 {
			if nickname[index] < 0x20 || nickname[index] == '"' {
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

func validateMemberPassword(password string) bool {
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
