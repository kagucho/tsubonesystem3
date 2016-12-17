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

import (
  `fmt`
  `strings`
)

type MemberClub struct {
  Chief bool `json:"chief"`
  ID string `json:"id"`
  Name string `json:"name"`
}

// MemberEntry is a structure to hold information about a member retrieved by
// QueryMembers.
type MemberEntry struct {
  Entrance uint16 `json:"entrance"`
  ID string `json:"id"`
  Nickname string `json:"nickname"`
  OB bool `json:"ob"`
  Realname string `json:"realname"`
}

type MemberGraph struct {
  Gender string
  Nickname string
}

type MemberPosition struct {
  ID string `json:"id"`
  Name string `json:"name"`
}

type MemberStatus uint

const (
  MemberStatusOB MemberStatus = 1 << iota
  MemberStatusActive MemberStatus = 1 << iota
)

// MemberDetail is a structure to hold the details of a member.
type MemberDetail struct {
  Affiliation string
  Clubs <-chan MemberClub
  ClubsErrors <-chan error
  Entrance uint16
  Gender string
  Mail string
  Nickname string
  OB bool
  Positions <-chan MemberPosition
  PositionsErrors <-chan error
  Realname string
  Tel string
}

// QueryMember returns a MemberDetail holding the details of the member
// identified with the given ID.
func (db DB) QueryMember(display string) (MemberDetail, error) {
  var dbMember uint16
  var output MemberDetail
  if scanError := db.sql.QueryRow(
       `SELECT id,affiliation,entrance,gender,mail,nickname,ob,realname,tel FROM members WHERE display_id=?`,
       display).Scan(
         &dbMember, &output.Affiliation, &output.Entrance, &output.Gender,
         &output.Mail, &output.Nickname, &output.OB, &output.Realname,
         &output.Tel)
     scanError != nil {
    return MemberDetail{}, scanError
  }

  clubs := make(chan MemberClub)
  clubsErrors := make(chan error)

  go func() {
    defer close(clubs)
    defer close(clubsErrors)

    rows, queryError :=
      db.sql.Query(`SELECT club FROM club_member WHERE member=?`, dbMember)
    if queryError != nil {
      clubsErrors <- queryError
      return
    }

    defer rows.Close()

    for rows.Next() {
      var dbClub uint8

      if scanError := rows.Scan(&dbClub); scanError != nil {
        clubsErrors <- scanError
        continue
      }

      var club MemberClub
      var clubChief uint16
      if scanError := db.sql.QueryRow(
           `SELECT chief,display_id,name FROM clubs WHERE id=?`,
           dbClub).Scan(&clubChief, &club.ID, &club.Name);
         scanError != nil {
        clubsErrors <- scanError
        continue
      }

      club.Chief = dbMember == clubChief

      clubs <- club
    }
  }()

  positions := make(chan MemberPosition)
  positionsErrors := make(chan error)

  go func() {
    defer close(positions)
    defer close(positionsErrors)

    rows, queryError :=
      db.sql.Query(`SELECT display_id,name FROM officers WHERE member=?`,
                   dbMember)
    if queryError != nil {
      positionsErrors <- queryError
      return
    }

    defer rows.Close()

    for rows.Next() {
      var position MemberPosition

      if scanError := rows.Scan(&position.ID, &position.Name)
         scanError != nil {
        positionsErrors <- scanError
        continue
      }

      positions <- position
    }
  }()

  output.Clubs = clubs
  output.ClubsErrors = clubsErrors

  output.Positions = positions
  output.PositionsErrors = positionsErrors

  return output, nil
}

func (db DB) QueryMemberGraph(id string) (MemberGraph, error) {
  var graph MemberGraph

  scanError := db.sql.QueryRow(
    `SELECT gender,nickname FROM members WHERE display_id=?`, id).Scan(
      &graph.Gender, &graph.Nickname)

  return graph, scanError
}

// QueryMembers returns a channel which provides information about members.
func (db DB) QueryMembers() (<-chan MemberEntry, <-chan error) {
  memberChan := make(chan MemberEntry)
  errorChan := make(chan error)

  go func() {
    defer close(memberChan)
    defer close(errorChan)

    rows, queryError := db.sql.Query(`SELECT entrance,display_id,nickname,ob,realname FROM members`)
    if queryError != nil {
      errorChan <- queryError
      return
    }

    defer rows.Close()

    for rows.Next() {
      var member MemberEntry

      if scanError := rows.Scan(&member.Entrance, &member.ID, &member.Nickname,
                                &member.OB, &member.Realname)
         scanError != nil {
        errorChan <- scanError
        continue
      }

      memberChan <- member
    }
  }()

  return memberChan, errorChan
}

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

  queries := make([]string, 0, 3)
  queries = append(queries, `SELECT COUNT(*) FROM members WHERE nickname LIKE ? AND realname LIKE ?`)

  arguments := make([]interface{}, 0, 3)
  arguments = append(arguments, pattern(nickname))
  arguments = append(arguments, pattern(realname))

  if entrance != 0 {
    queries = append(queries, `entrance=?`)
    arguments = append(arguments, entrance)
  }

  switch (status) {
  case 0:
    return 0, nil

  case MemberStatusOB:
    queries = append(queries, `ob`)

  case MemberStatusActive:
    queries = append(queries, `NOT ob`)

  case MemberStatusOB | MemberStatusActive:

  default:
    return 0, fmt.Errorf(`invalid status %v`, status)
  }

  var count uint16

  scanError := db.sql.QueryRow(
    strings.Join(queries, ` AND `), arguments...).Scan(&count)

  return count, scanError
}

func ValidateMemberEntrance(entrance int) bool {
  return entrance >= 1901 && entrance <= 2155
}
