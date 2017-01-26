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

package member

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

	request := httptest.NewRequest(`GET`, `http://kagucho.net/member/list`, nil)

	recorder := httptest.NewRecorder()

	ListServeHTTP(recorder, request, context.Context{DB: db}, authorizer.Claim{})

	if recorder.Code != http.StatusOK {
		t.Error(`invalid code; expected `, http.StatusOK, `, got `, recorder.Code)
	}

	const expected = `[{"affiliation":"理学部第一部 数理情報科学科","entrance":1901,"id":"1stDisplayID","nickname":"1 !\\%_1\"#","ob":false,"realname":"$&\\%_2'("},{"entrance":1901,"id":"2ndDisplayID","nickname":"2 !%_1\"#","ob":false,"realname":"$&\\%_2'("},{"entrance":1901,"id":"3rdDisplayID","nickname":"3 !\\%*1\"#","ob":false,"realname":"$&\\%_2'("},{"entrance":1901,"id":"4thDisplayID","nickname":"4 !)_1\"#","ob":false,"realname":"$&\\%_2'("},{"entrance":1901,"id":"5thDisplayID","nickname":"5 !\\%_1\"#","ob":false,"realname":"$&%+2'("},{"entrance":2155,"id":"6thDisplayID","nickname":"6 !\\%_1\"#","ob":false,"realname":"$&\\%+2'("},{"entrance":1901,"id":"7thDisplayID","nickname":"7 !\\%_1\"#","ob":true,"realname":"$&,_2'("}]
`
	if result := recorder.Body.String(); result != expected {
		t.Error(`expected `, expected, `, got `, result)
	}
}
