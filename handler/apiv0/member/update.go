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
	"net/http"
	"github.com/kagucho/tsubonesystem3/db"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/common"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/context"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/token/authorizer"
	"strconv"
	"strings"
)

func parseClubs(clubsString string) []string {
	if clubsString != `` {
		return strings.Split(clubsString, ` `)
	}

	return nil
}

func parseEntrance(entranceString string) (int, error) {
	entrance, entranceError := strconv.Atoi(entranceString)
	if entranceError != nil {
		return 0, entranceError
	}

	if !db.ValidateMemberEntrance(entrance) {
		return 0, strconv.ErrRange
	}

	return entrance, nil
}

func validateID(id string) bool {
	for index := 0; index < len(id); index++ {
		/*
			URL Standard
			5.2. application/x-www-form-urlencoded serializing
			https://url.spec.whatwg.org/#urlencoded-serializing

			> 0x2A
			> 0x2D
			> 0x2E
			> 0x30 to 0x39
			> 0x41 to 0x5A
			> 0x5F
			> 0x61 to 0x7A
			>
			> Append a code point whose value is byte to output.

			Accept only those characters.
		*/
		if !(id[index] == 0x2A || id[index] == 0x2D || id[index] == 0x2E ||
			(id[index] >= 0x30 && id[index] <= 0x39) ||
			(id[index] >= 0x41 && id[index] <= 0x5A) ||
			id[index] == 0x5F ||
			(id[index] >= 0x61 && id[index] <= 0x7A)) {
			return false
		}
	}

	return true
}

func validateTel(tel string) bool {
	for index := 0; index < len(tel); index++ {
		/*
			RFC 3986 - Uniform Resource Identifier (URI): Generic Syntax
			https://tools.ietf.org/html/rfc3986#section-2
			2.2.  Reserved Characters

			Allow characters valid in hier-part.
		*/
		if !(tel[index] == 0x21 || tel[index] == 0x24 ||
			(tel[index] >= 0x26 && tel[index] <= 0x39) ||
			tel[index] == 0x3B || tel[index] == 0x3D ||
			(tel[index] >= 0x41 && tel[index] <= 0x5A) ||
			tel[index] == 0x5F ||
			(tel[index] >= 0x61 && tel[index] <= 0x7A) ||
			tel[index] == 0x7E) {
			return false
		}
	}

	return true
}

func CreateServeHTTP(writer http.ResponseWriter, request *http.Request, context context.Context, claim authorizer.Claim) {
	id := request.FormValue(`id`)
	if !validateID(id) {
		common.ServeError(writer,
			common.Error{Description: `invalid id`},
			http.StatusBadRequest)

		return
	}

	mail := request.FormValue(`mail`)
	nickname := request.FormValue(`nickname`)

	createError := context.DB.InsertMember(id, mail, nickname)
	if createError != nil {
		common.ServeErrorDefault(writer, http.StatusBadRequest)

		return
	}

	token, tokenError := context.Token.IssueTemporaryAccessUpdater(id)
	if tokenError != nil {
		common.ServeErrorDefault(writer, http.StatusInternalServerError)

		return
	}

	mailError := context.Mail.MailCreation(id, mail, nickname, request.URL.Host, token)
	if mailError != nil {
		common.ServeError(writer,
			common.Error{Description: mailError.Error()},
			http.StatusBadRequest)

		return
	}

	common.ServeJSON(writer, struct{}{}, http.StatusOK)
}

func DeclareOBServeHTTP(writer http.ResponseWriter, request *http.Request, context context.Context, claim authorizer.Claim) {
	if declareError := context.DB.DeclareOB(claim.Sub); declareError == nil {
		common.ServeJSON(writer, struct{}{}, http.StatusOK)
	} else {
		common.ServeErrorDefault(writer, http.StatusBadRequest)
	}
}

func UpdateServeHTTP(writer http.ResponseWriter, request *http.Request, context context.Context, claim authorizer.Claim) {
	var entrance int
	if entranceString := request.PostFormValue(`entrance`); entranceString == `` {
		entrance = 0
	} else {
		var entranceError error
		entrance, entranceError = parseEntrance(entranceString)
		if entranceError != nil {
			common.ServeError(writer,
				common.Error{Description: entranceError.Error()},
				http.StatusBadRequest)

			return
		}
	}

	tel := request.PostFormValue(`tel`)
	if !validateTel(tel) {
		common.ServeError(writer,
			common.Error{Description: `invalid tel`},
			http.StatusBadRequest)

		return
	}

	password := request.PostFormValue(`password`)
	if password != `` {
		// TODO this token should be valid only when the password is not set.
	}

	clubs := parseClubs(request.PostFormValue(`clubs`))

	updateError := context.DB.UpdateMember(claim.Sub,
		password,
		request.PostFormValue(`affiliation`),
		clubs, entrance,
		request.PostFormValue(`gender`),
		request.PostFormValue(`mail`),
		request.PostFormValue(`nickname`),
		request.PostFormValue(`realname`),
		tel)
	if updateError == nil {
		common.ServeJSON(writer, struct{}{}, http.StatusOK)
	} else {
		common.ServeErrorDefault(writer, http.StatusBadRequest)
	}
}
