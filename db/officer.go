/*
  Copyright (C) 2016  Kagucho <kagucho.net@gmail.com>

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU Affero General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU Affero General Public License for more details.

  You should have received a copy of the GNU Affero General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package db

import `strings`

type OfficerDetail struct {
  Member OfficerMember `json:"member"`
  Name string `json:"name"`
  Scope []string `json:"scope"`
}

// Officer is a structure to hold information of an officer.
type OfficerEntry struct {
  ID string `json:"id"`
  Member OfficerMember `json:"member"`
  Name string `json:"name"`
}

type OfficerMember struct {
  ID string `json:"id"`
  Mail string `json:"mail"`
  Nickname string `json:"nickname"`
  Realname string `json:"realname"`
  Tel string `json:"tel"`
}

func (db DB) QueryOfficer(id string) (OfficerDetail, error) {
  var detail OfficerDetail
  var member uint16
  var scope string

  if scanError := db.sql.QueryRow(
       `SELECT member,name,scope FROM officers WHERE display_id=?`,
       id).Scan(&member, &detail.Name, &scope)
     scanError != nil {
    return OfficerDetail{}, scanError
  }

  if scanError := db.sql.QueryRow(
       `SELECT display_id,mail,nickname,realname,tel FROM members WHERE id=?`,
       member).Scan(
         &detail.Member.ID, &detail.Member.Mail,
         &detail.Member.Nickname, &detail.Member.Realname,
         &detail.Member.Tel)
     scanError != nil {
    return OfficerDetail{}, scanError
  }

  detail.Scope = strings.Split(scope, `,`)

  return detail, nil
}

func (db DB) QueryOfficerName(id string) (string, error) {
  var name string

  scanError := db.sql.QueryRow(
    `SELECT name FROM officers WHERE display_id=?`, id).Scan(&name)

  return name, scanError
}

// QueryOfficers returns a channel which provides information about officers.
func (db DB) QueryOfficers() (<-chan OfficerEntry, <-chan error) {
  officerChan := make(chan OfficerEntry)
  errorChan := make(chan error)

  go func() {
    defer close(officerChan)
    defer close(errorChan)

    rows, queryError :=
      db.sql.Query(`SELECT display_id,member,name FROM officers`)
    if queryError != nil {
      errorChan <- queryError
      return
    }

    defer rows.Close()

    for rows.Next() {
      var officer OfficerEntry
      var member uint16

      if scanError := rows.Scan(&officer.ID, &member, &officer.Name)
         scanError != nil {
        errorChan <- scanError
        continue
      }

      if scanError := db.sql.QueryRow(
           `SELECT display_id,mail,nickname,realname,tel FROM members WHERE id=?`,
           member).Scan(
             &officer.Member.ID, &officer.Member.Mail,
             &officer.Member.Nickname, &officer.Member.Realname,
             &officer.Member.Tel)
         scanError != nil {
        errorChan <- scanError
        continue
      }

      officerChan <- officer
    }
  }()

  return officerChan, errorChan
}
