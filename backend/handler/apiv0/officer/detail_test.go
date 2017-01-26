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

func TestDetailServeHTTP(t *testing.T) {
	db, dbError := db.Prepare()
	if dbError != nil {
		t.Fatal(dbError)
	}

	context := context.Context{DB: db}

	func() {
		defer db.Close()

		for _, test := range [...]struct {
			description  string
			request      string
			responseCode int
			responseBody string
		}{
			{
				`invalid`, ``, http.StatusBadRequest,
				`{"error":"invalid_id","error_description":"invalid ID","error_uri":"https://tools.ietf.org/html/rfc7231#section-6.5.1"}
`,
			}, {
				`valid`, `id=president`, http.StatusOK,
				`{"member":{"id":"1stDisplayID","mail":"1st@kagucho.net","nickname":"1 !\\%_1\"#","realname":"$&\\%_2'(","tel":"000-000-001"},"name":"局長","scope":["management","privacy"]}
`,
			},
		} {
			request := httptest.NewRequest(`GET`,
				`http://kagucho.net/api/v0/officer/detail?`+test.request, nil)

			recorder := httptest.NewRecorder()
			DetailServeHTTP(recorder, request, context, authorizer.Claim{})

			if recorder.Code != test.responseCode {
				t.Errorf(`expected %v, got %v`,
					test.responseCode, recorder.Code)
			}

			if result := recorder.Body.String(); result != test.responseBody {
				t.Error(`expected `, test.responseBody,
					`, got `, result)
			}
		}
	}()

	t.Run(`internalServerError`, func(t *testing.T) {
		t.Parallel()

		request := httptest.NewRequest(`GET`,
			`http://kagucho.net/api/v0/officer/detail?id=president`, nil)

		recorder := httptest.NewRecorder()
		DetailServeHTTP(recorder, request, context, authorizer.Claim{})

		if recorder.Code != http.StatusInternalServerError {
			t.Errorf(`expected %v, got %v`,
				http.StatusInternalServerError, recorder.Code)
		}
	})
}
