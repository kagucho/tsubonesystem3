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
	"github.com/kagucho/tsubonesystem3/backend/db"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/common"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/context"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/token/authorizer"
	"net/http"
	"strconv"
	"time"
)

func parseTime(encoded string) (time.Time, error) {
	// RFC 7519 NumericDate
	parsed, parseError := strconv.ParseInt(encoded, 10, 64)
	if parseError != nil {
		return time.Time{}, parseError
	}

	return time.Unix(parsed, 0).In(time.Local), nil
}

func CreateServeHTTP(writer http.ResponseWriter, request *http.Request,
	context context.Context, claim authorizer.Claim) {
	start, startError := parseTime(request.FormValue(`start`))
	if startError != nil {
		common.ServeError(writer,
			common.Error{Description: startError.Error()},
			http.StatusBadRequest)

		return
	}

	end, endError := parseTime(request.FormValue(`end`))
	if endError != nil {
		common.ServeError(writer,
			common.Error{Description: endError.Error()},
			http.StatusBadRequest)

		return
	}

	due, dueError := parseTime(request.FormValue(`due`))
	if dueError != nil {
		common.ServeError(writer,
			common.Error{Description: dueError.Error()},
			http.StatusBadRequest)

		return
	}

	name := request.FormValue(`name`)
	place := request.FormValue(`place`)
	inviteds := request.FormValue(`inviteds`)
	invitedsName := request.FormValue(`inviteds_name`)
	details := request.FormValue(`details`)

	mails, insertError := context.DB.InsertParty(name, start, end, place, due, inviteds, invitedsName, details)
	if insertError == db.IncorrectIdentity {
		common.ServeError(writer,
			common.Error{Description: `unknown inviteds`},
			http.StatusBadRequest)

		return
	} else if insertError != nil {
		panic(insertError)
	}

	inviteError := context.Mail.Invite(request.Host, mails, name, start, end, place, invitedsName, due, details)
	if inviteError != nil {
		common.ServeMailError(writer)

		return
	}

	common.ServeJSON(writer, struct{}{}, http.StatusOK)
}
