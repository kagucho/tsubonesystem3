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

package user

import (
	"github.com/kagucho/tsubonesystem3/backend/db"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/common"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/context"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/token/authorizer"
	"github.com/kagucho/tsubonesystem3/backend/mail"
	"net/http"
)

// DetailServeHTTP serves the detail of the member identified with the given ID
// via HTTP.
func DetailServeHTTP(writer http.ResponseWriter, request *http.Request,
	context context.Context, claim authorizer.Claim) {
	detail, queryError := context.DB.QueryMemberDetail(claim.Sub)
	switch queryError {
	case nil:
		var mailError error
		detail.Mail, mailError = mail.AddressToUnicode(detail.Mail)
		if mailError != nil {
			panic(mailError)
		}

		common.ServeJSON(writer, detail, http.StatusOK)

	case db.IncorrectIdentity:
		common.ServeErrorDefault(writer, http.StatusBadRequest)

	default:
		panic(queryError)
	}
}
