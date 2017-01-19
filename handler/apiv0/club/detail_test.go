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

package club

import (
	"github.com/kagucho/tsubonesystem3/db"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/token/authorizer"
	"github.com/kagucho/tsubonesystem3/scope"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDetailServeHTTP(t *testing.T) {
	db, dbError := db.Prepare()
	if dbError != nil {
		t.Fatal(dbError)
	}

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
				`valid`, `id=prog`, http.StatusOK,
				`{"chief":{"id":"2ndDisplayID","mail":"","nickname":"2 !%_1\"#","realname":"$&\\%_2'(","tel":""},"members":[{"entrance":1901,"id":"2ndDisplayID","nickname":"2 !%_1\"#","realname":"$&\\%_2'("},{"entrance":1901,"id":"1stDisplayID","nickname":"1 !\\%_1\"#","realname":"$&\\%_2'("}],"name":"Prog部"}
`,
			},
		} {
			test := test

			t.Run(test.description, func(t *testing.T) {
				request := httptest.NewRequest(`GET`,
					`http://kagucho.net/api/v0/club/detail?`+test.request, nil)

				recorder := httptest.NewRecorder()
				DetailServeHTTP(recorder, request, db,
					authorizer.Claim{Scope: scope.Scope{}.Set(scope.Basic)})

				if recorder.Code != test.responseCode {
					t.Errorf(`expected %v, got %v`,
						test.responseCode, recorder.Code)
				}

				if result := recorder.Body.String(); result != test.responseBody {
					t.Error(`expected `, test.responseBody, `, got `, result)
				}
			})
		}
	}()

	t.Run(`internalServerError`, func(t *testing.T) {
		t.Parallel()

		request := httptest.NewRequest(`GET`,
			`http://kagucho.net/api/v0/club/detail?id=prog`, nil)

		recorder := httptest.NewRecorder()
		DetailServeHTTP(recorder, request, db,
			authorizer.Claim{Scope: scope.Scope{}.Set(scope.Basic)})

		if recorder.Code != http.StatusInternalServerError {
			t.Errorf(`expected %v, got %v`,
				http.StatusInternalServerError, recorder.Code)
		}
	})
}
