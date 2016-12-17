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
  `database/sql`
  `strconv`
  `testing`
)

func (db DB) testQueryMember(t *testing.T) {
  t.Run(`valid`, func(t *testing.T) {
    t.Parallel()

    detail, queryError := db.QueryMember(`1stDisplayID`)
    if queryError != nil {
      t.Fatal(queryError)
    }

    if expected := `理学部第一部 数理情報科学科`
       detail.Affiliation != expected {
      t.Errorf(`invalid affiliation; expected %q, got %q`,
               expected, detail.Affiliation)
    }

    if expected := uint16(1901); detail.Entrance != expected {
      t.Error(`invalid entrance; expected `, expected,
              `, got `, detail.Entrance)
    }

    if expected := `男`; detail.Gender != expected {
      t.Errorf(`invalid gender; expected %q, got %q`, expected, detail.Gender)
    }

    if expected := `1st@kagucho.net`; detail.Mail != expected {
      t.Errorf(`invalid mail; expected %q, got %q`, expected, detail.Mail)
    }

    if expected := ` !\%_1"#`; detail.Nickname != expected {
      t.Errorf(`invalid nickname; expected %q, got %q`,
               expected, detail.Nickname)
    }

    if expected := false; detail.OB != expected {
      t.Error(`invalid ob; expected `, expected, `, got `, detail.OB)
    }

    if expected := `$&\%_2'(`; detail.Realname != expected {
      t.Errorf(`invalid realname; expected %q, got %q`,
               expected, detail.Realname)
    }

    if expected := `012-345-567`; detail.Tel != expected {
      t.Errorf(`invalid tel; expected %q, got %q`, expected, detail.Tel)
    }

    positionCount := 0
    var positionResult MemberPosition
    for detail.Positions != nil || detail.PositionsErrors != nil {
      select {
      case result, present := <-detail.Positions:
        if present {
          positionResult = result
          positionCount++
        } else {
          detail.Positions = nil
        }

      case result, present := <-detail.PositionsErrors:
        if present {
          t.Error(result)
        } else {
          detail.PositionsErrors = nil
        }
      }
    }

    if expected := 1; positionCount != expected {
      t.Error(`invalid positions number; expected `, expected,
              `, got `, positionCount)
    }

    if expected := (MemberPosition{`president`, `局長`})
       positionResult != expected {
      t.Error(`invalid position; expected `, expected ,
              `, got `, positionResult)
    }

    clubCount := 0
    var clubResult MemberClub
    for detail.Clubs != nil || detail.ClubsErrors != nil {
      select {
      case result, present := <-detail.Clubs:
        if present {
          clubResult = result
          clubCount++
        } else {
          detail.Clubs = nil
        }

      case result, present := <-detail.ClubsErrors:
        if present {
          t.Error(result)
        } else {
          detail.ClubsErrors = nil
        }
      }
    }

    if expected := 1; clubCount != expected {
      t.Error(`invalid clubs number; expected `, expected, `, got `, clubCount)
    } else {
      if expected := (MemberClub{false, `prog`, `Prog部`})
         clubResult != expected {
        t.Error(`invalid club; expected `, expected, `, got `, clubResult)
      }
    }
  })

  t.Run(`invalid`, func(t *testing.T) {
    t.Parallel()

    detail, queryError := db.QueryMember(``)

    if (detail != MemberDetail{}) {
      t.Error(`invalid detail; expected zero value, got `, detail)
    }

    if queryError != sql.ErrNoRows {
      t.Error(`invalid error; expected `, sql.ErrNoRows, `, got `, queryError)
    }
  })
}

func (db DB) testQueryMemberGraph(t *testing.T) {
  graph, queryError := db.QueryMemberGraph(`1stDisplayID`)

  if queryError != nil {
    t.Fatal(queryError)
  }

  expected := MemberGraph{`男`, ` !\%_1"#`}
  if graph != expected {
    t.Error(`expected `, expected, `, got `, graph)
  }
}

func (db DB) testQueryMembers(t *testing.T) {
  expected := [...]MemberEntry{
    {
      Entrance: 1901, ID: `1stDisplayID`, Nickname: ` !\%_1"#`, OB: false,
      Realname: `$&\%_2'(`,
    }, {
      Entrance: 1901, ID: `2ndDisplayID`, Nickname: ` !%_1"#`, OB: false,
      Realname: `$&\%_2'(`,
    }, {
      Entrance: 1901, ID: `3rdDisplayID`, Nickname: ` !\%*1"#`, OB: false,
      Realname: `$&\%_2'(`,
    }, {
      Entrance: 1901, ID: `4thDisplayID`, Nickname: ` !)_1"#`, OB: false,
      Realname: `$&\%_2'(`,
    }, {
      Entrance: 1901, ID: `5thDisplayID`, Nickname: ` !\%_1"#`, OB: false,
      Realname: `$&%+2'(`,
    }, {
      Entrance: 2155, ID: `6thDisplayID`, Nickname: ` !\%_1"#`, OB: false,
      Realname: `$&\%+2'(`,
    }, {
      Entrance: 1901, ID: `7thDisplayID`, Nickname: ` !\%_1"#`, OB: true,
      Realname: `$&,_2'(`,
    },
  }

  count := 0
  members, errors := db.QueryMembers()
  for members != nil || errors != nil {
    select {
    case result, present := <-members:
      if present {
        if count < len(expected) {
          if result != expected[count] {
            t.Error(`invalid member; expected `, expected[count],
                    `, got `, result)
          }

          count++
        } else if count == len(expected) {
          t.Error(`invalid member; expected `, len(expected),
                  ` member, got more`)
          count++
        }
      } else {
        members = nil
      }

    case result, present := <-errors:
      if present {
        t.Error(result)
      } else {
        errors = nil
      }
    }
  }
}

// DB columns TODO: write up SQL and add the below as the comment.
// { 0, 1, ` !\%_1"#`, `$&\%_2'(` } valid names which need to be escaped
// { 0, 2, ` !%_1"#`, `$&\%_2'(` } nickname lacking `\`
// { 0, 3, ` !\%*1"#`, `$&\%_2'(` } one invalid character in nickname
// { 0, 4, ` !\)_1"#`, `$&\%_2'(` } another invalid character in nickname
// { 0, 5, ` !\%_1"#`, `$&%+2'(` } realname lacking `\`
// { 0, 6, ` !\%_1"#`, `$&\%+2'(` } one invalid character in realname
// { 0, 7, ` !\%_1"#`, `$&\,_2'(` } another invalid character in realname
// { 65535, 8, ` !\%_1"#`, `$&\%_2'(` } different entrance
func (db DB) testQueryMembersCount(t *testing.T) {
  for _, test := range [...]struct{
        description string
        entrance int
        nickname string
        realname string
        status MemberStatus
        expectedCount uint16
        expectedError string
      }{
        {
          `full`, 1901, `\%_1`, `\%_2`,
          MemberStatusActive | MemberStatusOB, 1, ``,
        },

        //  !\%_1"# | $&\%_2'(
        // Test whether it checks entrance and ignores if the given value is 0.
        {
          `1901NoEntrance`, 0, `\%_1`, `\%_2`,
          MemberStatusActive | MemberStatusOB, 1, ``,
        },
        {
          `1901BadEntrance`, 2155, `\%_1`, `\%_2`,
          MemberStatusActive | MemberStatusOB, 0, ``,
        },
        {
          `21552155`, 2155, `!\%_1"#`, `$&\%+2'(`,
          MemberStatusActive | MemberStatusOB, 1, ``,
        },
        {
          `2155NoEntrance`, 0, `!\%_1"#`, `$&\%+2'(`,
          MemberStatusActive | MemberStatusOB, 1, ``,
        },
        {
          `2155BadEntrance`, 1901, `!\%_1"#`, `$&\%+2'(`,
          MemberStatusActive | MemberStatusOB, 0, ``,
        },

        // Test whether it checks entrance

        // Test whether it checks nickname. It should escape `\`, `%` and `_` in
        // patterns.
        {
          `badNickname`, 1901, `\%1`, `\%_2`,
          MemberStatusActive | MemberStatusOB, 0, ``,
        },

        // Test whether it checks status.
        {`activeNone`, 1901, `\%_1`, `\%_2`, 0, 0, ``},
        {`activeActive`, 1901, `\%_1`, `\%_2`, MemberStatusActive, 1, ``},
        {`activeOB`, 1901, `\%_1`, `\%_2`, MemberStatusOB, 0, ``},
        {`obNone`, 1901, `!\%_1"#`, `$&,_2'(`, 0, 0, ``},
        {`obActive`, 1901, `!\%_1"#`, `$&,_2'(`, MemberStatusActive, 0, ``},
        {`obOB`, 1901, `!\%_1"#`, `$&,_2'(`, MemberStatusOB, 1, ``},
        {
          `obAll`, 1901, `!\%_1"#`, `$&,_2'(`,
          MemberStatusActive | MemberStatusOB, 1, ``,
        },
        {`invalidStatus`, 1901, `\%_1`, `\%_2`, 4, 0, `invalid status 4`},

        // Test whether it checks realname. It should escape `\`, `%` and `_` in
        // patterns.
        {`wrongNickname`, 1901, `\%_1`, `\%2`, MemberStatusActive, 0, ``},
      } {
    test := test

    t.Run(test.description, func(t *testing.T) {
      t.Parallel()
      count, queryError := db.QueryMembersCount(
        test.entrance, test.nickname, test.realname, test.status)

      if count != test.expectedCount {
        t.Error(`invalid count; expected `, test.expectedCount, `, got `, count)
      }

      if (queryError == nil && test.expectedError != ``) ||
         (queryError != nil && queryError.Error() != test.expectedError) {
        t.Errorf(`invalid error; expected %q, got %q`,
                 test.expectedError, queryError)
      }
    })
  }
}

func TestValidateMemberEntrance(t *testing.T) {
  t.Parallel()

  for _, test := range [...]struct{
        entrance int
        expected bool
      }{{1900, false}, {1901, true}, {2155, true}, {2156, false}} {
    test := test

    t.Run(strconv.Itoa(test.entrance), func(t *testing.T) {
      t.Parallel()

      if result := ValidateMemberEntrance(test.entrance)
         result != test.expected {
        t.Error(`expected `, test.expected, `, got `, result)
      }
    })
  }
}
