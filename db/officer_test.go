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
  `reflect`
  `testing`
)

func (db DB) testQueryOfficer(t *testing.T) {
  t.Run(`valid`, func(t *testing.T) {
    t.Parallel()

    expected := OfficerDetail{
      OfficerMember{
        `1stDisplayID`, `1st@kagucho.net`, ` !\%_1"#`, `$&\%_2'(`,
        `012-345-567`,
      },
      `局長`,
      []string{ `management`, `privacy` },
    }
    if detail, queryError := db.QueryOfficer(`president`); queryError != nil {
      t.Error(queryError)
    } else if !reflect.DeepEqual(detail, expected) {
      t.Errorf(`expected %v, got %v`, expected, detail)
    }
  })

  t.Run(`invalid`, func(t *testing.T) {
    t.Parallel()

    if detail, queryError := db.QueryOfficer(``); queryError != sql.ErrNoRows {
      t.Errorf(`invalid error; expected %v, got %v`, sql.ErrNoRows, queryError)
    } else if !reflect.DeepEqual(detail, OfficerDetail{}) {
      t.Error(`expected zero value, got `, detail)
    }
  })
}

func (db DB) testQueryOfficerName(t *testing.T) {
  if name, queryError := db.QueryOfficerName(`president`); queryError != nil {
    t.Error(queryError)
  } else if name != `局長` {
    t.Errorf(`expected "局長", got %q`, name)
  }
}

func (db DB) testQueryOfficers(t *testing.T) {
  officerChan, errorChan := db.QueryOfficers()

  index := 0
  var officer OfficerEntry
  for officerChan != nil || errorChan != nil {
    select {
    case result, present := <-officerChan:
      if present {
        if (officer == OfficerEntry{}) {
          officer = result
        }
        index++
      } else {
        officerChan = nil
      }

    case result, present := <-errorChan:
      if present {
        t.Error(result)
      } else {
        errorChan = nil
      }
    }
  }

  if index != 1 {
    t.Errorf(`invalid officers number; expected 1, got %v`, index)
  } else {
    expected := OfficerEntry{
      `president`,
      OfficerMember{
        `1stDisplayID`, `1st@kagucho.net`, ` !\%_1"#`, `$&\%_2'(`,
        `012-345-567`,
      },
      `局長`,
    }

    if !reflect.DeepEqual(officer, expected) {
      t.Errorf(`expected %v, got %v`, expected, officer)
    }
  }
}