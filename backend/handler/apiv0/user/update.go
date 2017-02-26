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
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/context"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/common"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/token/authorizer"
	"github.com/kagucho/tsubonesystem3/backend/mail"
	"log"
	"net/http"
	netMail "net/mail"
	"strconv"
	"strings"
)

func parseClubs(clubsString string) []string {
	if clubsString != `` {
		return strings.Split(clubsString, ` `)
	}

	return nil
}

func ConfirmServeHTTP(writer http.ResponseWriter, request *http.Request, context context.Context, claim authorizer.Claim) {
	tokenClaim, tokenError := context.Token.AuthenticateMail(request.FormValue(`token`))
	if tokenError.IsError() {
		common.ServeErrorDefault(writer, http.StatusBadRequest)

		return
	}

	if claim.Sub != tokenClaim.Sub {
		common.ServeErrorDefault(writer, http.StatusBadRequest)

		return
	}

	confirmError := context.DB.ConfirmMember(claim.Sub)
	if confirmError != nil {
		common.ServeErrorDefault(writer, http.StatusBadRequest)

		return
	}

	common.ServeJSON(writer, struct{}{}, http.StatusOK)
}

func DeclareOBServeHTTP(writer http.ResponseWriter, request *http.Request, context context.Context, claim authorizer.Claim) {
	if declareError := context.DB.DeclareMemberOB(claim.Sub); declareError == nil {
		common.ServeJSON(writer, struct{}{}, http.StatusOK)
	} else {
		common.ServeErrorDefault(writer, http.StatusBadRequest)
	}
}

func UpdatePasswordServeHTTP(writer http.ResponseWriter, request *http.Request, context context.Context, claim authorizer.Claim) {
	switch updateError := context.DB.UpdatePassword(claim.Sub, request.PostFormValue(`current`), request.PostFormValue(`new`)); updateError {
	case db.IncorrectIdentity:
		common.ServeError(writer,
			common.Error{Description: `invalid current`},
			http.StatusBadRequest)

	case db.MemberInvalidPassword:
		common.ServeError(writer,
			common.Error{Description: `invalid new`},
			http.StatusBadRequest)

	case nil:
		common.ServeJSON(writer, struct{}{}, http.StatusOK)

	default:
		panic(updateError)
	}
}

func UpdateServeHTTP(writer http.ResponseWriter, request *http.Request, context context.Context, claim authorizer.Claim) {
	var entrance int
	if entranceString := request.PostFormValue(`entrance`); entranceString == `` {
		entrance = 0
	} else {
		var entranceError error
		entrance, entranceError = strconv.Atoi(entranceString)
		if entranceError != nil {
			common.ServeError(writer,
				common.Error{Description: entranceError.Error()},
				http.StatusBadRequest)

			return
		}
	}

	password := request.PostFormValue(`password`)
	if password != `` && !claim.Tmp {
		common.ServeError(writer,
			common.Error{Description: `cannot use this endpoint to set password without temporary token`},
			http.StatusBadRequest)

		return
	}

	affiliation := request.PostFormValue(`affiliation`)
	clubs := parseClubs(request.PostFormValue(`clubs`))
	gender := request.PostFormValue(`gender`)
	address := request.PostFormValue(`mail`)
	nickname := request.PostFormValue(`nickname`)
	realname := request.PostFormValue(`realname`)
	tel := request.PostFormValue(`tel`)

	if address != `` {
		var addressError error
		address, addressError = mail.AddressToASCII(address)
		if addressError != nil {
			panic(addressError)
		}
	}

	updateError := context.DB.UpdateMember(claim.Sub,
		password, affiliation, clubs, entrance,
		gender, address, nickname, realname,
		tel)
	if updateError != nil {
		log.Print(updateError)
		common.ServeErrorDefault(writer, http.StatusBadRequest)

		return
	}

	if address != `` {
		token, tokenError := context.Token.IssueMail(claim.Sub)
		if tokenError != nil {
			panic(tokenError)
		}

		nickname, queryError := context.DB.QueryMemberNickname(claim.Sub)
		if queryError != nil {
			panic(queryError)
		}

		mailError := context.Mail.SendConfirmation(request.Host,
			netMail.Address{nickname, address}, claim.Sub, token)
		if mailError != nil {
			common.ServeError(writer,
				common.Error{
					ID:          `mail_failure`,
					Description: `failed to mail`,
				}, http.StatusOK)

			log.Print(mailError)

			return
		}
	}

	common.ServeJSON(writer, struct{}{}, http.StatusOK)
}
