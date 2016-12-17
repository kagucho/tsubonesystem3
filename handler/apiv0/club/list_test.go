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

package club

import (
  `github.com/kagucho/tsubonesystem3/db`
  `github.com/kagucho/tsubonesystem3/handler/apiv0/token/authorizer`
  `net/http`
  `net/http/httptest`
  `testing`
)

func TestListServeHTTP(t *testing.T) {
  t.Parallel()

  db, dbError := db.Open()
  if dbError != nil {
    t.Fatal(dbError)
  }

  defer db.Close()

  request := httptest.NewRequest(`GET`, `http://kagucho.net/club/list`, nil)
  recorder := httptest.NewRecorder()

  ListServeHTTP(recorder, request, db, authorizer.Claim{})

  if recorder.Code != http.StatusOK {
    t.Error(`invalid code; expected `, http.StatusOK, `, got `, recorder.Code)
  }

  const expected = `[{"chief":{"id":"2ndDisplayID","mail":"","nickname":" !%_1\"#","realname":"$\u0026\\%_2'(","tel":""},"id":"prog","name":"Prog部"}]`
  if result := recorder.Body.String(); result != expected {
    t.Error(`expected `, expected, `, got `, result)
  }
}
