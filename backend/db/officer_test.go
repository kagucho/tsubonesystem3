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
	"database/sql"
	"reflect"
	"testing"
)

func (db DB) testQueryOfficer(t *testing.T) {
	t.Run(`valid`, func(t *testing.T) {
		t.Parallel()

		expected := Officer{
			OfficerMember{
				`1stDisplayID`, `1st@kagucho.net`,
				`1 !\%_1"#`, `$&\%_2'(`, `000-000-001`,
			},
			`局長`,
			[]string{`management`, `privacy`},
		}
		if detail, err := db.QueryOfficer(`president`); err != nil {
			t.Error(err)
		} else if !reflect.DeepEqual(detail, expected) {
			t.Errorf(`expected %v, got %v`, expected, detail)
		}
	})

	t.Run(`invalid`, func(t *testing.T) {
		t.Parallel()

		if detail, err := db.QueryOfficer(``); err != sql.ErrNoRows {
			t.Errorf(`invalid error; expected %v, got %v`,
				sql.ErrNoRows, err)
		} else if !reflect.DeepEqual(detail, Officer{}) {
			t.Error(`expected zero value, got `, detail)
		}
	})
}

func (db DB) testQueryOfficerName(t *testing.T) {
	if name, err := db.QueryOfficerName(`president`); err != nil {
		t.Error(err)
	} else if name != `局長` {
		t.Errorf(`expected "局長", got %q`, name)
	}
}

func (db DB) testQueryOfficers(t *testing.T) {
	const expected = `[{"id":"president","member":{"id":"1stDisplayID","mail":"1st@kagucho.net","nickname":"1 !\\%_1\"#","realname":"$&\\%_2'(","tel":"000-000-001"},"name":"局長"},{"id":"vice","member":{"id":"1stDisplayID","mail":"1st@kagucho.net","nickname":"1 !\\%_1\"#","realname":"$&\\%_2'(","tel":"000-000-001"},"name":"副局長"}]`
	result, err := db.QueryOfficers().MarshalJSON()

	if err != nil {
		t.Error(err)
	}

	if resultString := string(result); resultString != expected {
		t.Error(`expected `, expected, `, got `, resultString)
	}
}
