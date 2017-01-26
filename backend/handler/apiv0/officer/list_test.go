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

package officer

import (
	"github.com/kagucho/tsubonesystem3/backend/db"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/context"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/token/authorizer"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListServeHTTP(t *testing.T) {
	t.Parallel()

	db, dbError := db.Prepare()
	if dbError != nil {
		t.Fatal(dbError)
	}

	defer db.Close()

	request := httptest.NewRequest(`GET`,
		`http://kagucho.net/member/officers`, nil)

	recorder := httptest.NewRecorder()

	ListServeHTTP(recorder, request, context.Context{DB: db}, authorizer.Claim{})

	if recorder.Code != http.StatusOK {
		t.Errorf(`invalid code; expected %v, got %v`,
			http.StatusOK, recorder.Code)
	}

	const expected = `[{"id":"president","member":{"id":"1stDisplayID","mail":"1st@kagucho.net","nickname":"1 !\\%_1\"#","realname":"$&\\%_2'(","tel":"000-000-001"},"name":"局長"},{"id":"vice","member":{"id":"1stDisplayID","mail":"1st@kagucho.net","nickname":"1 !\\%_1\"#","realname":"$&\\%_2'(","tel":"000-000-001"},"name":"副局長"}]
`
	if result := recorder.Body.String(); result != expected {
		t.Error(`expected `, expected, `, got `, result)
	}
}
