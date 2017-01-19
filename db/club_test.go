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
	"testing"
)

func (db DB) testQueryClub(t *testing.T) {
	t.Run(`valid`, func(t *testing.T) {
		t.Parallel()

		detail, queryError := db.QueryClub(`prog`)
		if queryError != nil {
			t.Fatal(queryError)
		}

		if expected := (Chief{`2ndDisplayID`, ``, `2 !%_1"#`, `$&\%_2'(`, ``}); detail.Chief != expected {
			t.Errorf(`invalid chief; expected %v, got %v`,
				expected, detail.Chief)
		}

		if expected := `Prog部`; detail.Name != expected {
			t.Errorf(`invalid name; expected %v, got %v`,
				expected, detail.Name)
		}

		const membersExpected = `[{"entrance":1901,"id":"2ndDisplayID","nickname":"2 !%_1\"#","realname":"$&\\%_2'("},{"entrance":1901,"id":"1stDisplayID","nickname":"1 !\\%_1\"#","realname":"$&\\%_2'("}]`
		membersResult, membersError := detail.Members.MarshalJSON()
		if membersError != nil {
			t.Error(`invalid members; `, membersError)
		}

		if resultString := string(membersResult); resultString != membersExpected {
			t.Error(`invalid members; expected `, membersExpected,
				` got `, resultString)
		}
	})

	t.Run(`invalid`, func(t *testing.T) {
		t.Parallel()

		detail, queryError := db.QueryClub(``)

		if (detail != Club{}) {
			t.Error(`invalid club; expected zero value, got `,
				detail)
		}

		if queryError != sql.ErrNoRows {
			t.Errorf(`invalid error; expected %v, got %v`,
				sql.ErrNoRows, queryError)
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

func (db DB) testQueryClubNames(t *testing.T) {
	const expected = `[{"id":"prog","name":"Prog部"},{"id":"web","name":"Web部"}]`
	result, resultError := db.QueryClubNames().MarshalJSON()
	if resultError != nil {
		t.Error(resultError)
	}

	if resultString := string(result); resultString != expected {
		t.Error(`expected `, expected, `, got `, resultString)
	}
}

func (db DB) testQueryClubs(t *testing.T) {
	const expected = `[{"id":"prog","name":"Prog部","chief":{"id":"2ndDisplayID","mail":"","nickname":"2 !%_1\"#","realname":"$&\\%_2'(","tel":""}},{"id":"web","name":"Web部","chief":{"id":"1stDisplayID","mail":"1st@kagucho.net","nickname":"1 !\\%_1\"#","realname":"$&\\%_2'(","tel":"012-345-567"}}]`

	result, resultError := db.QueryClubs().MarshalJSON()
	if resultError != nil {
		t.Error(resultError)
	}

	if resultString := string(result); resultString != expected {
		t.Error(`expected `, expected, `, got `, resultString)
	}
}
