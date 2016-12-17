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
  `github.com/kagucho/tsubonesystem3/handler/apiv0/common`
  `github.com/kagucho/tsubonesystem3/handler/apiv0/token/authorizer`
  `net/http`
)

// ListServeHTTP serves the functionality of API v0 to list clubs via HTTP.
func ListServeHTTP(writer http.ResponseWriter, request *http.Request,
                   dbInstance db.DB, claim authorizer.Claim) {
  clubs, errors := dbInstance.QueryClubs()
  common.ServeJSONChan(writer, clubs, clubs, errors, http.StatusOK)
}

func ListNameServeHTTP(writer http.ResponseWriter, request *http.Request,
                       dbInstance db.DB, claim authorizer.Claim) {
  clubs, errors := dbInstance.QueryClubNames()
  common.ServeJSONChan(writer, clubs, clubs, errors, http.StatusOK)
}
