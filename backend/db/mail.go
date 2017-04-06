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

// RecipientResult is a structure representing a result of querying a recipient.
type RecipientResult struct {
	Error     error
	Recipient string
}

// RecipientChan is a reciever of db.RecipientResult.
type RecipientChan <-chan RecipientResult

/*
MailCommon is a structure holding information about a email common for some
queries.
*/
type MailCommon struct {
	Date encoding.Time       `json:"date"`
	From encoding.ZeroString `json:"from"`
	To   string              `json:"to"`
}

// MailDetail is a structure holding details about a email.
type MailDetail struct {
	MailCommon
	Recipients RecipientChan `json:"recipients"`
	Body       string        `json:"body"`
}

// MailEntry is a structure holding basic information about a email.
type MailEntry struct {
	MailCommon
	Subject string `json:"subject"`
}

/*
MailEntryResult is a structure representing a result of querying db.MailEntry.
*/
type MailEntryResult struct {
	MailEntry
	Error error
}

// MailEntryChan is a reciever of db.MailEntryResult.
type MailEntryChan <-chan MailEntryResult

/*
MarshalJSON returns the JSON encoding of the remaining recipients and closes
the channel.
This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (recipientChan RecipientChan) MarshalJSON() ([]byte, error) {
	return encoding.MarshalJSONArray(func() (interface{}, error, bool) {
		result, present := <-recipientChan
		return result.Recipient, result.Error, present
	})
}

/*
MarshalJSON returns the JSON encoding of the remaining entries and closes
the channel.
This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (entryChan MailEntryChan) MarshalJSON() ([]byte, error) {
	return encoding.MarshalJSONArray(func() (interface{}, error, bool) {
		result, present := <-entryChan
		return result.MailEntry, result.Error, present
	})
}

/*
InsertMail inserts an email with the given properties and returns the nickname
of From and the IDs of the recipients.

It may return one of the following errors:
db.ErrBadOmission tells recipients, to, subject, or body is omitted.
db.ErrDupEntry tells the subject is duplicate.
db.ErrIncorrectIdentity tells one of the IDs of the recipients or the one of
From is incorrect.
db.ErrInvalid tells one of the given properties is invalid.

Other errors tell db.DB is bad.
*/
func (db DB) InsertMail(recipients, from, to, subject, body string) (returnedFrom string, returnedMails []string, returnedErr error) {
	var fromDBID uint16
	var fromNickname string

	if recipients == `` || to == `` || subject == `` || body == `` {
		return ``, nil, ErrBadOmission
	}

	recipientsCount, recipientsDB := stringListToDBList(recipients)
	recipientMails := make([]string, recipientsCount)

	if nicknameErr := db.stmts[stmtSelectMemberInternalIDNicknameByID].QueryRow(from).Scan(&fromDBID, &fromNickname); nicknameErr != nil {
		if nicknameErr == sql.ErrNoRows {
			nicknameErr = ErrIncorrectIdentity
		}

		return ``, nil, nicknameErr
	}

	rows, mailErr := db.stmts[stmtCallInsertMail].Query(recipientsDB, recipientsCount, fromDBID, to, subject, body)
	if mailErr != nil {
		if mysqlErr, ok := mailErr.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case erDataTooLong:
				fallthrough
			case erTruncatedWrongValueForField:
				return ``, nil, ErrInvalid

			case erDupEntry:
				return ``, nil, ErrDupEntry

			case erNoReferencedRow:
				fallthrough
			case erNoReferencedRow2:
				fallthrough
			case erSignalException:
				return ``, nil, ErrIncorrectIdentity
			}
		}

		return ``, nil, mailErr
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Print(closeErr)
		}
	}()

	for rows.Next() {
		var mail string

		if scanErr := rows.Scan(&mail); scanErr != nil {
			return ``, nil, scanErr
		}

		recipientMails = append(recipientMails, mail)
	}

	return fromNickname, recipientMails, nil
}

/*
DeleteMail delets the email with the given subject.

It returns db.ErrIncorrectIdentity if the given subject is incorrect. Other
errors tell db.DB is bad.
*/
func (db DB) DeleteMail(subject string) error {
	result, execErr := db.stmts[stmtDeleteMail].Exec(subject)
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
QueryMail returns db.MailDetail describing the email with the given subject.

It returns db.ErrIncorrectIdentity if the given subject is incorrect. Other
errors tell db.DB is bad.

Resources will be holded until Recipients gets closed.
*/
func (db DB) QueryMail(subject string) (MailDetail, error) {
	var date mysql.NullTime
	var dbID uint16
	var from sql.NullString
	var mail MailDetail

	tx, err := db.sql.BeginTx(context.Background(),
		&sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  true,
		})
	if err != nil {
		return mail, err
	}

	if err := tx.Stmt(db.stmts[stmtSelectMailBySubject]).QueryRow(subject).Scan(
		&dbID, &date, &from, &mail.To, &mail.Body); err != nil {
		if err := tx.Commit(); err != nil {
			log.Print(err)
		}

		if err == sql.ErrNoRows {
			err = ErrIncorrectIdentity
		}

		return mail, err
	}

	recipientChan := make(chan RecipientResult)

	go func() {
		defer func() {
			close(recipientChan)

			if err := tx.Commit(); err != nil {
				log.Print(err)
			}
		}()

		rows, err := tx.Stmt(db.stmts[stmtSelectRecipientsByInternalMail]).Query(dbID)
		if err != nil {
			recipientChan <- RecipientResult{Error: err}
			return
		}

		defer func() {
			if err := rows.Close(); err != nil {
				log.Print(err)
			}
		}()

		for rows.Next() {
			var result RecipientResult
			result.Error = rows.Scan(&result.Recipient)

			recipientChan <- result
			if result.Error != nil {
				return
			}
		}
	}()

	mail.Date = encoding.NewTime(date.Time)
	mail.From = encoding.ZeroString(from.String)
	mail.Recipients = recipientChan

	return mail, nil
}

/*
QueryMails returns a channel sending db.MailEntryResult of all mails.

db.MailEntryResult will have an error if db.DB is bad.

Resources will be holded until the channel gets closed.
*/
func (db DB) QueryMails() MailEntryChan {
	resultChan := make(chan MailEntryResult)

	go func() {
		defer close(resultChan)

		rows, err := db.stmts[stmtSelectMails].Query()
		if err != nil {
			resultChan <- MailEntryResult{Error: err}
			return
		}

		defer func() {
			if err := rows.Close(); err != nil {
				log.Print(err)
			}
		}()

		for rows.Next() {
			var date mysql.NullTime
			var from sql.NullString
			var result MailEntryResult

			result.Error = rows.Scan(&date, &from,
				&result.To, &result.Subject)
			result.Date = encoding.NewTime(date.Time)
			result.From = encoding.ZeroString(from.String)

			resultChan <- result
			if result.Error != nil {
				return
			}
		}
	}()

	return resultChan
}

/*
UpdateMail updates the email with the given subject, with the given properties.

It may return one of the following errors:
db.ErrIncorrectIdentity tells the IDs of the recipients or from is incorrect.
db.ErrInvalid tells some of the given properties is invalid.

Other errors tell db.DB is bad.
*/
func (db DB) UpdateMail(subject, recipients string, date encoding.Time, from, to, body string) error {
	arguments := make([]interface{}, 7)
	arguments[0] = subject

	if recipients != `` {
		arguments[2], arguments[1] = stringListToDBList(recipients)
	}

	if (date != encoding.Time{}) {
		arguments[3] = date
	}

	if from != `` {
		arguments[4] = from
	}

	if to != `` {
		arguments[5] = to
	}

	if body != `` {
		arguments[6] = body
	}

	result, execErr := db.stmts[stmtCallUpdateMail].Exec(arguments...)
	if execErr != nil {
		if mysqlErr, ok := execErr.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case erDataTooLong:
				fallthrough
			case erTruncatedWrongValueForField:
				fallthrough
			case erWrongValue:
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
