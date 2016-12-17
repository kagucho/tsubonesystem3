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

package member

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

  request := httptest.NewRequest(`GET`, `http://kagucho.net/member/list`, nil)

  recorder := httptest.NewRecorder()

  ListServeHTTP(recorder, request, db, authorizer.Claim{})

  if recorder.Code != http.StatusOK {
    t.Error(`invalid code; expected `, http.StatusOK, `, got `, recorder.Code)
  }

  const expected = `[{"entrance":1901,"id":"1stDisplayID","nickname":" !\\%_1\"#","ob":false,"realname":"$\u0026\\%_2'("},{"entrance":1901,"id":"2ndDisplayID","nickname":" !%_1\"#","ob":false,"realname":"$\u0026\\%_2'("},{"entrance":1901,"id":"3rdDisplayID","nickname":" !\\%*1\"#","ob":false,"realname":"$\u0026\\%_2'("},{"entrance":1901,"id":"4thDisplayID","nickname":" !)_1\"#","ob":false,"realname":"$\u0026\\%_2'("},{"entrance":1901,"id":"5thDisplayID","nickname":" !\\%_1\"#","ob":false,"realname":"$\u0026%+2'("},{"entrance":2155,"id":"6thDisplayID","nickname":" !\\%_1\"#","ob":false,"realname":"$\u0026\\%+2'("},{"entrance":1901,"id":"7thDisplayID","nickname":" !\\%_1\"#","ob":true,"realname":"$\u0026,_2'("}]`
  if result := recorder.Body.String(); result != expected {
    t.Error(`expected `, expected, `, got `, result)
  }
}
