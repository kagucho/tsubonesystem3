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

package apiv0

import (
	"github.com/kagucho/tsubonesystem3/backend/db"
	"github.com/kagucho/tsubonesystem3/backend/encoding"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/util"
	"github.com/kagucho/tsubonesystem3/backend/scope"
	"net/http"
	"strconv"
)

func mailDeleteServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path == `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	if (authorize(writer, request, shared, scope.Management) == claim{}) {
		return
	}

	switch err := shared.DB.DeleteMail(request.URL.Path[1:]); err {
	case db.ErrIncorrectIdentity:
		util.ServeErrorDefault(writer, http.StatusNotFound)

	case nil:
		util.ServeJSON(writer, struct{}{}, http.StatusOK)

	default:
		panic(err)
	}
}

func mailGetServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path == `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	if (authorize(writer, request, shared, scope.Member) == claim{}) {
		return
	}

	switch mail, err := shared.DB.QueryMail(request.URL.Path[1:]); err {
	case db.ErrIncorrectIdentity:
		util.ServeErrorDefault(writer, http.StatusNotFound)

	case nil:
		util.ServeJSON(writer, mail, http.StatusOK)

	default:
		panic(err)
	}
}

func mailPatchServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path == `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	if (authorize(writer, request, shared, scope.Management) == claim{}) {
		return
	}

	var date encoding.Time

	if dateString := request.FormValue(`date`); dateString != `` {
		var err error

		date, err = encoding.ParseQueryTime(dateString)
		if err != nil {
			description := `syntax error in date`
			status := http.StatusBadRequest

			if err == strconv.ErrRange {
				description = `date out of range`
				status = http.StatusUnprocessableEntity
			}

			util.ServeError(writer,
				util.Error{Description: description},
				status)

			return
		}
	}

	switch err := shared.DB.UpdateMail(request.URL.Path[1:],
		request.PostFormValue(`recipients`), date,
		request.PostFormValue(`from`), request.PostFormValue(`to`),
		request.PostFormValue(`body`)); err {
	case db.ErrIncorrectIdentity:
		// FIXME: should be StatusUnprocessableEntity if mail is found
		// but recipients or from is not.
		util.ServeErrorDefault(writer, http.StatusNotFound)

	case db.ErrInvalid:
		util.ServeErrorDefault(writer, http.StatusUnprocessableEntity)

	case nil:
		util.ServeJSON(writer, struct{}{}, http.StatusOK)

	default:
		panic(err)
	}
}

func mailPutServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if len(request.URL.Path) < 2 {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	authorized := authorize(writer, request, shared, scope.Member)
	if (authorized == claim{}) {
		return
	}

	body := request.PostFormValue(`body`)
	to := request.PostFormValue(`to`)
	subject := request.URL.Path[1:]

	switch nickname, recipients, err := shared.DB.InsertMail(
		request.PostFormValue(`recipients`), authorized.sub, to, subject, body); err {
	case db.ErrBadOmission:
		util.ServeError(writer,
			util.Error{Description: `recipients, to, subject, and body are required`},
			http.StatusUnprocessableEntity)

	case db.ErrDupEntry:
		// FIXME: should be StatusMethodNotAllowed
		util.ServeError(writer,
			util.Error{Description: `duplicate subject`},
			http.StatusUnprocessableEntity)

	case db.ErrIncorrectIdentity:
		// FIXME: should be StatusUnprocessableEntity if mail is found
		// but user or from is not.
		util.ServeErrorDefault(writer, http.StatusNotFound)

	case db.ErrInvalid:
		util.ServeErrorDefault(writer, http.StatusUnprocessableEntity)

	case nil:
		if err := shared.Mail.Message(request.Host, recipients, authorized.sub, nickname, to, subject, body); err != nil {
			util.ServeMailError(writer)

			return
		}

		util.ServeJSON(writer, struct{}{}, http.StatusCreated)

	default:
		panic(err)
	}
}

func mailsGetServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path != `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	if (authorize(writer, request, shared, scope.Member) == claim{}) {
		return
	}

	util.ServeJSON(writer, shared.DB.QueryMails(), http.StatusOK)
}
