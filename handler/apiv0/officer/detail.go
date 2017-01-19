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
	"database/sql"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/common"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/context"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/token/authorizer"
	"net/http"
)

// DetailServeHTTP serves the functionality of API v0 to detail officers via
// HTTP.
func DetailServeHTTP(writer http.ResponseWriter, request *http.Request,
	context context.Context, claim authorizer.Claim) {
	serve := func() func() {
		defer common.Recover(writer)

		id := request.FormValue(`id`)
		officer, queryError := context.DB.QueryOfficer(id)
		switch queryError {
		case nil:
			return func() {
				common.ServeJSON(writer, officer, http.StatusOK)
			}

		case sql.ErrNoRows:
			return func() {
				common.ServeError(writer,
					common.Error{
						ID: `invalid_id`,
						Description: `invalid ID`,
					}, http.StatusBadRequest)
			}

		default:
			panic(queryError)
		}
	}()

	if serve != nil {
		serve()
	}
}
