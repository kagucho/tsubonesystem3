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

import "testing"

func TestDB(t *testing.T) {
	var db DB

	if !t.Run(`Prepare`, func(t *testing.T) {
		var err error
		db, err = Prepare()
		if err != nil {
			t.Fatal(err)
		}
	}) {
		t.FailNow()
	}

	t.Run(`QueryClub`, db.testQueryClub)
	t.Run(`QueryClubName`, db.testQueryClubName)
	t.Run(`QueryClubNames`, db.testQueryClubNames)
	t.Run(`QueryClubs`, db.testQueryClubs)
	t.Run(`QueryMember`, db.testQueryMember)
	t.Run(`QueryMemberGraph`, db.testQueryMemberGraph)
	t.Run(`QueryMembers`, db.testQueryMembers)
	t.Run(`QueryMembersCount`, db.testQueryMembersCount)
	t.Run(`QueryOfficer`, db.testQueryOfficer)
	t.Run(`QueryOfficerName`, db.testQueryOfficerName)
	t.Run(`QueryOfficers`, db.testQueryOfficers)
	t.Run(`GetScope`, db.testGetScope)

	t.Run(`Close`, func(t *testing.T) {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	})
}
