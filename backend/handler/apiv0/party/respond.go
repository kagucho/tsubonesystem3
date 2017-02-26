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

package party

import (
	"fmt"
	"github.com/kagucho/tsubonesystem3/backend/db"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/common"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/context"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/token/authorizer"
	"net/http"
)

func RespondServeHTTP(writer http.ResponseWriter, request *http.Request,
	context context.Context, claim authorizer.Claim) {
	var attending bool

	switch attendingString := request.FormValue(`attending`); attendingString {
	case `0`:
		attending = false

	case `1`:
		attending = true

	default:
		common.ServeError(writer,
			common.Error{Description: fmt.Sprint(`invalid attending: '%v'`, attending)},
			http.StatusBadRequest)

		return
	}

	switch dbError := context.DB.UpdateAttendance(attending,
		request.FormValue(`party`), claim.Sub); dbError {
	case nil:
		common.ServeJSON(writer, struct{}{}, http.StatusOK)

	case db.IncorrectIdentity:
		common.ServeError(writer,
			common.Error{Description: `unknown combination of party and user`},
			http.StatusBadRequest)

	default:
		panic(dbError)
	}
}
