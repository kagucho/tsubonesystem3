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

// Chief is a structure to hold information of a chief.
type Chief struct {
  ID string `json:"id"`
  Mail string `json:"mail"`
  Nickname string `json:"nickname"`
  Realname string `json:"realname"`
  Tel string `json:"tel"`
}

type ClubDetail struct {
  Chief Chief `json:"chief"`
  Members <-chan ClubMember `json:"members"`
  MembersErrors <-chan error `json:"-"`
  Name string `json:"name"`
}

type ClubMember struct {
  Entrance uint16 `json:"entrance"`
  ID string `json:"id"`
  Nickname string `json:"nickname"`
  Realname string `json:"realname"`
}

type ClubNameEntry struct {
  ID string `json:"id"`
  Name string `json:"name"`
}

type ClubEntry struct {
  ClubNameEntry
  Chief Chief `json:"chief"`
}

func (db DB) QueryClub(id string) (ClubDetail, error) {
  var chiefID uint16
  var clubID uint8
  var club ClubDetail

  if scanError := db.sql.QueryRow(
       `SELECT chief,id,name FROM clubs WHERE display_id=?`,
       id).Scan(&chiefID, &clubID, &club.Name)
     scanError != nil {
    return ClubDetail{}, scanError
  }

  members := make(chan ClubMember)
  membersErrors := make(chan error)

  go func() {
    defer close(members)
    defer close(membersErrors)

    rows, queryError :=
      db.sql.Query(`SELECT member FROM club_member WHERE club=?`, clubID)
    if queryError != nil {
      membersErrors <- queryError
      return
    }

    defer rows.Close()

    for rows.Next() {
      var memberID uint16

      if scanError := rows.Scan(&memberID); scanError != nil {
        membersErrors <- scanError
        continue
      }

      var member ClubMember
      if scanError := db.sql.QueryRow(
           `SELECT entrance,display_id,nickname,realname FROM members WHERE id=?`,
           memberID).Scan(
             &member.Entrance, &member.ID, &member.Nickname, &member.Realname)
         scanError != nil {
        membersErrors <- scanError
        continue
      }

      members <- member
    }
  }()

  club.Members = members
  club.MembersErrors = membersErrors

  if scanError := db.sql.QueryRow(
       `SELECT display_id,mail,nickname,realname,tel FROM members WHERE id=?`,
       chiefID).Scan(
         &club.Chief.ID, &club.Chief.Mail,
         &club.Chief.Nickname, &club.Chief.Realname, &club.Chief.Tel)
     scanError != nil {
    return ClubDetail{}, scanError
  }

  return club, nil
}

func (db DB) QueryClubName(id string) (string, error) {
  var name string

  scanError := db.sql.QueryRow(
    `SELECT name FROM clubs WHERE display_id=?`, id).Scan(&name)

  return name, scanError
}

func (db DB) QueryClubNames() (<-chan ClubNameEntry, <-chan error) {
  entryChan := make(chan ClubNameEntry)
  errorChan := make(chan error)

  go func() {
    defer close(entryChan)
    defer close(errorChan)

    rows, queryError := db.sql.Query(`SELECT display_id,name FROM clubs`)
    if queryError != nil {
      errorChan <- queryError
      return
    }

    for rows.Next() {
      var entry ClubNameEntry
      scanError := rows.Scan(&entry.ID, &entry.Name)
      if scanError == nil {
        entryChan <- entry
      } else {
        errorChan <- scanError
      }
    }
  }()

  return entryChan, errorChan
}

func (db DB) QueryClubs() (<-chan ClubEntry, <-chan error) {
  clubChan := make(chan ClubEntry)
  errorChan := make(chan error)

  go func() {
    defer close(clubChan)
    defer close(errorChan)

    clubs, queryError := db.sql.Query(`SELECT chief,display_id,name FROM clubs`)
    if queryError != nil {
      errorChan <- queryError
      return
    }

    for clubs.Next() {
      var club ClubEntry
      var chiefID uint16

      if scanError := clubs.Scan(&chiefID, &club.ID, &club.Name)
         scanError != nil {
        errorChan <- scanError
        continue
      }

      if scanError := db.sql.QueryRow(
           `SELECT display_id,mail,nickname,realname,tel FROM members WHERE id=?`,
           chiefID).Scan(
             &club.Chief.ID, &club.Chief.Mail,
             &club.Chief.Nickname, &club.Chief.Realname, &club.Chief.Tel)
         scanError != nil {
        errorChan <- scanError
        continue
      }

      clubChan <- club
    }
  }()

  return clubChan, errorChan
}
