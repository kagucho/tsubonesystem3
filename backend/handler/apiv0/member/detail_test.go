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
	"github.com/kagucho/tsubonesystem3/backend/scope"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDetailServeHTTP(t *testing.T) {
	t.Parallel()

	db, dbError := db.Prepare()
	if dbError != nil {
		t.Fatal(dbError)
	}

	context := context.Context{DB: db}

	func() {
		defer db.Close()

		for _, test := range [...]struct {
			description  string
			requestQuery string
			requestClaim authorizer.Claim
			responseCode int
			responseBody string
		}{
			{
				`invalidID`, ``,
				authorizer.Claim{
					Sub:   `1stDisplayID`,
					Scope: scope.Scope{}.Set(scope.Basic),
				}, http.StatusBadRequest,
				`{"error":"invalid_id","error_description":"invalid ID","error_uri":"https://tools.ietf.org/html/rfc7231#section-6.5.1"}
`,
			}, {
				`president`, `id=1stDisplayID`,
				authorizer.Claim{
					Scope: scope.Scope{}.Set(scope.Basic),
				}, http.StatusOK,
				`{"affiliation":"理学部第一部 数理情報科学科","clubs":[{"chief":true,"id":"web","name":"Web部"},{"chief":false,"id":"prog","name":"Prog部"}],"entrance":1901,"gender":"男","mail":"1st@kagucho.net","nickname":"1 !\\%_1\"#","ob":false,"positions":[{"id":"president","name":"局長"},{"id":"vice","name":"副局長"}],"realname":"$&\\%_2'(","tel":"000-000-001"}
`,
			}, {
				`chief`, `id=2ndDisplayID`,
				authorizer.Claim{
					Scope: scope.Scope{}.Set(scope.Basic),
				}, http.StatusOK,
				`{"clubs":[{"chief":true,"id":"prog","name":"Prog部"}],"entrance":1901,"gender":"女","mail":"","nickname":"2 !%_1\"#","ob":false,"positions":[],"realname":"$&\\%_2'(","tel":"000-000-002"}
`,
			}, {
				`same`, `id=3rdDisplayID`,
				authorizer.Claim{
					Sub:   `3rdDisplayID`,
					Scope: scope.Scope{}.Set(scope.Basic),
				}, http.StatusOK,
				`{"clubs":[],"entrance":1901,"mail":"","nickname":"3 !\\%*1\"#","ob":false,"positions":[],"realname":"$&\\%_2'(","tel":"000-000-003"}
`,
			}, {
				`privacy`, `id=3rdDisplayID`,
				authorizer.Claim{
					Scope: scope.Scope{}.Set(scope.Basic).Set(scope.Privacy),
				}, http.StatusOK,
				`{"clubs":[],"entrance":1901,"mail":"","nickname":"3 !\\%*1\"#","ob":false,"positions":[],"realname":"$&\\%_2'(","tel":"000-000-003"}
`,
			}, {
				`normal`, `id=3rdDisplayID`,
				authorizer.Claim{
					Scope: scope.Scope{}.Set(scope.Basic),
				}, http.StatusOK,
				`{"clubs":[],"entrance":1901,"mail":"","nickname":"3 !\\%*1\"#","ob":false,"positions":[],"realname":"$&\\%_2'("}
`,
			},
		} {
			t.Run(test.description, func(t *testing.T) {
				request := httptest.NewRequest(`GET`,
					`http://kagucho.net/api/v0/member/detail?`+test.requestQuery, nil)

				recorder := httptest.NewRecorder()
				DetailServeHTTP(recorder, request, context, test.requestClaim)

				if recorder.Code != test.responseCode {
					t.Error(`expected `, test.responseCode,
						`, got `, recorder.Code)
				}

				if result := recorder.Body.String(); result != test.responseBody {
					t.Error(`expected `, test.responseBody,
						`, got `, result)
				}
			})
		}
	}()

	t.Run(`internalServerError`, func(t *testing.T) {
		t.Parallel()

		request := httptest.NewRequest(`GET`,
			`http://kagucho.net/api/v0/member/detail?id=1stDisplayID`, nil)

		recorder := httptest.NewRecorder()
		DetailServeHTTP(recorder, request, context,
			authorizer.Claim{Scope: scope.Scope{}.Set(scope.Basic)})

		if recorder.Code != http.StatusInternalServerError {
			t.Errorf(`expected %v, got %v`,
				http.StatusInternalServerError, recorder.Code)
		}
	})
}
