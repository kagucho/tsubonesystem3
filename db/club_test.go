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
  `testing`
)

func (db DB) testQueryClub(t *testing.T) {
  t.Run(`valid`, func(t *testing.T) {
    t.Parallel()

    detail, queryError := db.QueryClub(`prog`)
    if queryError != nil {
      t.Fatal(queryError)
    }

    if expected := (Chief{`2ndDisplayID`, ``, ` !%_1"#`, `$&\%_2'(`, ``})
       detail.Chief != expected {
      t.Errorf(`invalid chief; expected %v, got %v`, expected, detail.Chief)
    }

    if expected := `Prog部`; detail.Name != expected {
      t.Errorf(`invalid name; expected %v, got %v`, expected, detail.Name)
    }

    memberCount := 0
    memberExpected := [...]ClubMember{
      ClubMember{1901, `2ndDisplayID`, ` !%_1"#`, `$&\%_2'(`},
      ClubMember{1901, `1stDisplayID`, ` !\%_1"#`, `$&\%_2'(`},
    }
    var memberResult [len(memberExpected)]ClubMember
    for detail.Members != nil || detail.MembersErrors != nil {
      select {
      case result, present := <-detail.Members:
        if present {
          if memberCount < len(memberResult) {
            memberResult[memberCount] = result
          }

          memberCount++
        } else {
          detail.Members = nil
        }

      case result, present := <-detail.MembersErrors:
        if present {
          t.Error(result)
        } else {
          detail.MembersErrors = nil
        }
      }
    }

    if memberCount != len(memberExpected) {
      t.Errorf(`invalid members number; expected %v, got %v`,
               len(memberExpected), memberCount)
    } else if memberResult != memberExpected {
      t.Errorf(`invalid members result; expected %v, got %v`,
               memberExpected, memberResult)
    }
  })

  t.Run(`invalid`, func(t *testing.T) {
    t.Parallel()

    detail, queryError := db.QueryClub(``)

    if (detail != ClubDetail{}) {
      t.Error(`invalid detail; expected zero value, got `, detail)
    }

    if queryError != sql.ErrNoRows {
      t.Errorf(`invalid error; expected %v, got %v`, sql.ErrNoRows, queryError)
    }
  })
}

func (db DB) testQueryClubName(t *testing.T) {
  if name, queryError := db.QueryClubName(`prog`); queryError != nil {
    t.Error(queryError)
  } else if name != `Prog部` {
    t.Errorf(`expected "Prog部", got %q`, name)
  }
}

func (db DB) testQueryClubs(t *testing.T) {
  clubs, errors := db.QueryClubs()

  var club ClubEntry
  count := 0
  for clubs != nil || errors != nil {
    select {
    case result, present := <-clubs:
      if present {
        club = result
        count++
      } else {
        clubs = nil
      }

    case result, present := <-errors:
      if present {
        t.Error(result)
      } else {
        errors = nil
      }
    }
  }

  if expected := 1; count != expected {
    t.Errorf(`invalid clubs number; expected %v, got %v`, expected, count)
  } else {
    if expected := (ClubEntry{
         Chief{`2ndDisplayID`, ``, ` !%_1"#`, `$&\%_2'(`, ``},
         `prog`, `Prog部`,
       }); club != expected {
      t.Errorf(`invalid club; expected %v, got %v`, expected, club)
    }
  }
}
