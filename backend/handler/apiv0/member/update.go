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
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/common"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/context"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/token/authorizer"
	"github.com/kagucho/tsubonesystem3/backend/mail"
	"net/http"
	netMail "net/mail"
)

func CreateServeHTTP(writer http.ResponseWriter, request *http.Request, context context.Context, claim authorizer.Claim) {
	address := request.FormValue(`mail`)
	id := request.FormValue(`id`)
	nickname := request.FormValue(`nickname`)

	asciiAddress, addressError := mail.AddressToASCII(address)
	if addressError != nil {
		common.ServeError(writer, common.Error{Description: `invalid mail`},
			http.StatusBadRequest)

		return
	}

	createError := context.DB.InsertMember(id, asciiAddress, nickname)
	if createError != nil {
		common.ServeErrorDefault(writer, http.StatusBadRequest)

		return
	}

	token, tokenError := context.Token.IssueTmpUserAccess(id)
	if tokenError != nil {
		panic(tokenError)
	}

	mailError := context.Mail.SendCreation(request.Host, netMail.Address{nickname, asciiAddress}, id, token)
	if mailError != nil {
		common.ServeMailError(writer)

		return
	}

	common.ServeJSON(writer, struct{}{}, http.StatusOK)
}

/*
func ManageMailServeHTTP(writer http.ResponseWriter, request *http.Request, context context.Context, claim authorizer.Claim) {
	id := request.FormValue(`id`)
	mail := request.FormValue(`mail`)

	updateError := context.DB.UpdateMemberMail(id, mail)
	if updateError != nil {
		common.ServeErrorDefault(writer, http.StatusBadRequest)

		return
	}

	mailError := context.Mail.MailConfirmation(id, mail, nickname, request.Host)
	if mailError != nil {
		common.ServeError(writer,
			common.Error{ID: `mail_failure`, Description: mailError.Error()},
			http.StatusOK)

		return
	}

	common.ServeJSON(writer, struct{}{}, http.StatusOK)
}
*/
