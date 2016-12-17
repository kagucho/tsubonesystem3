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
  `database/sql`
  `fmt`
  `github.com/kagucho/tsubonesystem3/db`
  `github.com/kagucho/tsubonesystem3/handler/apiv0/common`
  `github.com/kagucho/tsubonesystem3/handler/apiv0/token/authorizer`
  `net/http`
)

// DetailServeHTTP serves the functionality of API v0 to detail a club via HTTP.
func DetailServeHTTP(writer http.ResponseWriter, request *http.Request,
                     dbInstance db.DB, claim authorizer.Claim) {
  serve := func() func() {
    defer common.Recover(writer)

    id := request.FormValue(`id`)
    detail, queryError := dbInstance.QueryClub(id)

    switch queryError {
    case nil:
      return func() {
        common.ServeJSONChan(
          writer, detail, detail.Members, detail.MembersErrors, http.StatusOK)
      }

    case sql.ErrNoRows:
      return func() {
        common.ServeError(writer, `invalid_id`,
          fmt.Sprintf(`unknown ID: %q`, id), ``, http.StatusBadRequest)
      }

    default:
      panic(queryError)
    }
  }()

  if serve != nil {
    serve()
  }
}
