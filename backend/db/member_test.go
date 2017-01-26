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
	"strconv"
	"testing"
)

func (db DB) testQueryMember(t *testing.T) {
	t.Run(`valid`, func(t *testing.T) {
		t.Parallel()

		detail, queryError := db.QueryMember(`1stDisplayID`)
		if queryError != nil {
			t.Fatal(queryError)
		}

		if expected := `理学部第一部 数理情報科学科`; detail.Affiliation != expected {
			t.Errorf(`invalid affiliation; expected %q, got %q`,
				expected, detail.Affiliation)
		}

		if expected := uint16(1901); detail.Entrance != expected {
			t.Error(`invalid entrance; expected %v, got %v`,
				expected, detail.Entrance)
		}

		if expected := `男`; detail.Gender != expected {
			t.Errorf(`invalid gender; expected %q, got %q`,
				expected, detail.Gender)
		}

		if expected := `1st@kagucho.net`; detail.Mail != expected {
			t.Errorf(`invalid mail; expected %q, got %q`,
				expected, detail.Mail)
		}

		if expected := `1 !\%_1"#`; detail.Nickname != expected {
			t.Errorf(`invalid nickname; expected %q, got %q`,
				expected, detail.Nickname)
		}

		if expected := false; detail.OB != expected {
			t.Errorf(`invalid ob; expected %v, got %v`,
				expected, detail.OB)
		}

		if expected := `$&\%_2'(`; detail.Realname != expected {
			t.Errorf(`invalid realname; expected %q, got %q`,
				expected, detail.Realname)
		}

		if expected := `000-000-001`; detail.Tel != expected {
			t.Errorf(`invalid tel; expected %q, got %q`,
				expected, detail.Tel)
		}

		expectedPositions := [...]Position{
			{`president`, `局長`}, {`vice`, `副局長`},
		}

		positionCount := 0
		var positionsResult [len(expectedPositions)]Position
		for result := range detail.Positions {
			if result.Error == nil {
				if positionCount < len(expectedPositions) {
					positionsResult[positionCount] = result.Value
					positionCount++
				}
			} else {
				t.Error(result.Error)
			}
		}

		if positionCount != len(expectedPositions) {
			t.Error(`invalid positions; expected`, len(expectedPositions), `positions`)
		} else if positionsResult != expectedPositions {
			t.Errorf(`invalid positions; expected %v, got %v`,
				expectedPositions, positionsResult)
		}

		expectedClubs := [...]MemberClub{
			{true, `web`, `Web部`}, {false, `prog`, `Prog部`},
		}

		clubCount := 0
		var clubsResult [len(expectedClubs)]MemberClub
		for result := range detail.Clubs {
			if result.Error == nil {
				if clubCount < len(expectedClubs) {
					clubsResult[clubCount] = result.Value
					clubCount++
				}
			} else {
				t.Error(result.Error)
			}
		}

		if clubCount != len(expectedClubs) {
			t.Error(`invalid clubs; expected`, len(expectedClubs), `clubs`)
		} else if clubsResult != expectedClubs {
			t.Errorf(`invalid clubs; expected %v, got %v`,
				expectedClubs, clubsResult)
		}
	})

	t.Run(`invalid`, func(t *testing.T) {
		t.Parallel()

		detail, queryError := db.QueryMember(``)

		if (detail != Member{}) {
			t.Errorf(`invalid member; expected zero value, got %v`,
				detail)
		}

		if queryError != sql.ErrNoRows {
			t.Errorf(`invalid error; expected %v, got %v`,
				sql.ErrNoRows, queryError)
		}
	})
}

func (db DB) testQueryMemberGraph(t *testing.T) {
	graph, queryError := db.QueryMemberGraph(`1stDisplayID`)

	if queryError != nil {
		t.Fatal(queryError)
	}

	expected := MemberGraph{`男`, `1 !\%_1"#`}
	if graph != expected {
		t.Errorf(`expected %v, got %v`, expected, graph)
	}
}

func (db DB) testQueryMembers(t *testing.T) {
	const expected = `[{"affiliation":"理学部第一部 数理情報科学科","entrance":1901,"id":"1stDisplayID","nickname":"1 !\\%_1\"#","ob":false,"realname":"$&\\%_2'("},{"entrance":1901,"id":"2ndDisplayID","nickname":"2 !%_1\"#","ob":false,"realname":"$&\\%_2'("},{"entrance":1901,"id":"3rdDisplayID","nickname":"3 !\\%*1\"#","ob":false,"realname":"$&\\%_2'("},{"entrance":1901,"id":"4thDisplayID","nickname":"4 !)_1\"#","ob":false,"realname":"$&\\%_2'("},{"entrance":1901,"id":"5thDisplayID","nickname":"5 !\\%_1\"#","ob":false,"realname":"$&%+2'("},{"entrance":2155,"id":"6thDisplayID","nickname":"6 !\\%_1\"#","ob":false,"realname":"$&\\%+2'("},{"entrance":1901,"id":"7thDisplayID","nickname":"7 !\\%_1\"#","ob":true,"realname":"$&,_2'("}]`
	result, resultError := db.QueryMembers().MarshalJSON()

	if resultError != nil {
		t.Error(resultError)
	}

	if resultString := string(result); resultString != expected {
		t.Error(`expected `, expected, `, got `, resultString)
	}
}

func (db DB) testQueryMembersCount(t *testing.T) {
	for _, test := range [...]struct {
		description   string
		entrance      int
		nickname      string
		realname      string
		status        MemberStatus
		expectedCount uint16
		expectedError string
	}{
		{
			`full`, 1901, `\%_1`, `\%_2`,
			MemberStatusActive | MemberStatusOB, 1, ``,
		},

		//  !\%_1"# | $&\%_2'(
		/*
			Test whether it checks entrance and ignores if the
			given value is 0.
		*/
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

		/*
			Test whether it checks nickname. It should escape `\`,
			`%` and `_` in patterns.
		*/
		{
			`badNickname`, 1901, `\%1`, `\%_2`,
			MemberStatusActive | MemberStatusOB, 0, ``,
		},

		// Test whether it checks status.
		{
			`activeNone`, 1901, `\%_1`, `\%_2`,
			0, 0, ``,
		},
		{
			`activeActive`, 1901, `\%_1`, `\%_2`,
			MemberStatusActive, 1, ``,
		},
		{
			`activeOB`, 1901, `\%_1`, `\%_2`,
			MemberStatusOB, 0, ``,
		},
		{
			`obNone`, 1901, `!\%_1"#`, `$&,_2'(`,
			0, 0, ``,
		},
		{
			`obActive`, 1901, `!\%_1"#`, `$&,_2'(`,
			MemberStatusActive, 0, ``,
		},
		{
			`obOB`, 1901, `!\%_1"#`, `$&,_2'(`,
			MemberStatusOB, 1, ``,
		},
		{
			`obAll`, 1901, `!\%_1"#`, `$&,_2'(`,
			MemberStatusActive | MemberStatusOB, 1, ``,
		},
		{
			`invalidStatus`, 1901, `\%_1`, `\%_2`,
			4, 0, `invalid status 4`,
		},

		/*
			Test whether it checks realname. It should escape `\`,
			`%` and `_` in patterns.
		*/
		{
			`wrongNickname`, 1901, `\%_1`, `\%2`,
			MemberStatusActive, 0, ``,
		},
	} {
		test := test

		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			count, queryError := db.QueryMembersCount(
				test.entrance, test.nickname, test.realname,
				test.status)

			if count != test.expectedCount {
				t.Errorf(`invalid count; expected %v, got %v`,
					test.expectedCount, count)
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

	for _, test := range [...]struct {
		entrance int
		expected bool
	}{{1900, false}, {1901, true}, {2155, true}, {2156, false}} {
		test := test

		t.Run(strconv.Itoa(test.entrance), func(t *testing.T) {
			t.Parallel()

			if result := ValidateMemberEntrance(test.entrance); result != test.expected {
				t.Errorf(`expected %v, got %v`,
					test.expected, result)
			}
		})
	}
}
