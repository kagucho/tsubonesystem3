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

package officer

import (
  `github.com/kagucho/tsubonesystem3/db`
  `github.com/kagucho/tsubonesystem3/handler/apiv0/token/authorizer`
  `net/http`
  `net/http/httptest`
  `testing`
)

func TestDetailServeHTTP(t *testing.T) {
  db, dbError := db.Open()
  if dbError != nil {
    t.Fatal(dbError)
  }

  func() {
    defer db.Close()

    for _, test := range [...]struct{
          description string
          request string
          responseCode int
          responseBody string
        }{
          {
            `invalid`, ``, http.StatusBadRequest,
            `{"error":"invalid_id","error_description":"unknown ID: \"\"","error_uri":"https://tools.ietf.org/html/rfc7231#section-6.5.1"}`,
          }, {
            `valid`, `id=president`, http.StatusOK,
            `{"member":{"id":"1stDisplayID","mail":"1st@kagucho.net","nickname":" !\\%_1\"#","realname":"$\u0026\\%_2'(","tel":"012-345-567"},"name":"局長","scope":["management","privacy"]}`,
          },
        } {
      request := httptest.NewRequest(`GET`,
        `http://kagucho.net/api/v0/officer/detail?` + test.request, nil)

      recorder := httptest.NewRecorder()
      DetailServeHTTP(recorder, request, db, authorizer.Claim{})

      if recorder.Code != test.responseCode {
        t.Error(`expected `, test.responseCode, `, got `, recorder.Code)
      }

      if result := recorder.Body.String(); result != test.responseBody {
        t.Errorf(`expected `, test.responseBody, `, got `, result)
      }
    }
  }()

  t.Run(`internalServerError`, func(t *testing.T) {
    t.Parallel()

    request := httptest.NewRequest(`GET`,
      `http://kagucho.net/api/v0/officer/detail?id=president`, nil)

    recorder := httptest.NewRecorder()
    DetailServeHTTP(recorder, request, db, authorizer.Claim{})

    if recorder.Code != http.StatusInternalServerError {
      t.Errorf(`expected %v, got %v`,
               http.StatusInternalServerError, recorder.Code)
    }
  })
}
