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
	"github.com/kagucho/tsubonesystem3/backend/scope"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/util"
	"net/http"
)

func clubDeleteServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path == `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	if (authorize(writer, request, shared, scope.Management) == claim{}) {
		return
	}

	switch err := shared.DB.DeleteClub(request.URL.Path[1:]); err {
	case db.ErrIncorrectIdentity:
		util.ServeErrorDefault(writer, http.StatusNotFound)

	case nil:
		util.ServeJSON(writer, struct{}{}, http.StatusOK)

	default:
		panic(err)
	}
}

func clubGetServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path == `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	if (authorize(writer, request, shared, scope.Member) == claim{}) {
		return
	}

	switch detail, err := shared.DB.QueryClub(request.URL.Path[1:]); err {
	case db.ErrIncorrectIdentity:
		util.ServeErrorDefault(writer, http.StatusNotFound)

	case nil:
		util.ServeJSON(writer, detail, http.StatusOK)

	default:
		panic(err)
	}
}

func clubPatchServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path == `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	if (authorize(writer, request, shared, scope.Management) == claim{}) {
		return
	}

	switch err := shared.DB.UpdateClub(
		request.URL.Path[1:],
		request.PostFormValue(`name`),
		request.PostFormValue(`chief`)); err {
	case db.ErrIncorrectIdentity:
		util.ServeErrorDefault(writer, http.StatusNotFound)

	case db.ErrInvalid:
		util.ServeErrorDefault(writer, http.StatusUnprocessableEntity)

	case nil:
		util.ServeJSON(writer, struct{}{}, http.StatusOK)

	default:
		panic(err)
	}
}

func clubPutServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if len(request.URL.Path) < 2 {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	if (authorize(writer, request, shared, scope.Management) == claim{}) {
		return
	}

	switch err := shared.DB.InsertClub(request.URL.Path[1:],
		request.PostFormValue(`name`),
		request.PostFormValue(`chief`)); err {
	case db.ErrBadOmission:
		util.ServeError(writer,
			util.Error{Description: `id and name are required`},
			http.StatusUnprocessableEntity)

	case db.ErrDupEntry:
		// FIXME: should be StatusMethodNotAllowed if id is duplicate
		util.ServeError(writer,
			util.Error{Description: `duplicate id or name`},
			http.StatusUnprocessableEntity)

	case db.ErrIncorrectIdentity:
		util.ServeErrorDefault(writer, http.StatusNotFound)

	case db.ErrInvalid:
		util.ServeErrorDefault(writer, http.StatusUnprocessableEntity)

	case nil:
		util.ServeJSON(writer, struct{}{}, http.StatusCreated)

	default:
		panic(err)
	}
}

func clubsGetServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path != `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	clubChan, err := shared.DB.QueryClubs()
	if err != nil {
		panic(err)
	}

	util.ServeJSON(writer, clubChan, http.StatusOK)
}
