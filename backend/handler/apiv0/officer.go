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
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/util"
	"github.com/kagucho/tsubonesystem3/backend/scope"
	"net/http"
)

func officerDeleteServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path == `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	authorized := authorize(writer, request, shared, scope.Management)
	if (authorized == claim{}) {
		return
	}

	switch err := shared.DB.DeleteOfficer(authorized.sub, request.URL.Path[1:]); err {
	case db.ErrIncorrectIdentity:
		util.ServeErrorDefault(writer, http.StatusNotFound)

	case db.ErrOfficerSuicide:
		util.ServeError(writer,
			util.Error{Description: `removing user's own management permission`},
			http.StatusUnprocessableEntity)

	case nil:
		util.ServeJSON(writer, struct{}{}, http.StatusOK)

	default:
		panic(err)
	}
}

func officerGetServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path == `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	if (authorize(writer, request, shared, scope.Member) == claim{}) {
		return
	}

	officer, err := shared.DB.QueryOfficerDetail(request.URL.Path[1:])
	switch err {
	case nil:
		util.ServeJSON(writer, officer, http.StatusOK)

	case db.ErrIncorrectIdentity:
		util.ServeErrorDefault(writer, http.StatusNotFound)

	default:
		panic(err)
	}
}

func officerPatchServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path == `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	authorized := authorize(writer, request, shared, scope.Management)
	if (authorized == claim{}) {
		return
	}

	err := request.ParseForm()
	if err != nil {
		util.ServeJSON(writer,
			util.Error{Description: err.Error()},
			http.StatusBadRequest)

		return
	}

	scope := db.NoScopeUpdate
	if formScope := request.PostForm[`scope`]; formScope != nil {
		scope = formScope[0]
	}

	switch err := shared.DB.UpdateOfficer(authorized.sub,
		request.PostForm.Get(`id`), request.PostForm.Get(`name`),
		request.PostForm.Get(`member`), scope); err {
	case db.ErrDupEntry:
		util.ServeError(writer,
			util.Error{Description: `duplicate name`},
			http.StatusBadRequest)

	case db.ErrIncorrectIdentity:
		/*
			FIXME: should be StatusUnprocessableEntity if officer is
			found but member is not.
		*/
		util.ServeErrorDefault(writer, http.StatusNotFound)

	case db.ErrInvalid:
		util.ServeErrorDefault(writer, http.StatusUnprocessableEntity)

	case db.ErrOfficerSuicide:
		util.ServeJSON(writer,
			util.Error{Description: `removing user's own management permission`},
			http.StatusUnprocessableEntity)

	case nil:
		util.ServeJSON(writer, struct{}{}, http.StatusOK)

	default:
		panic(err)
	}
}

func officerPutServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if len(request.URL.Path) < 2 {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	if (authorize(writer, request, shared, scope.Management) == claim{}) {
		return
	}

	switch err := shared.DB.InsertOfficer(request.URL.Path[1:],
		request.PostFormValue(`name`), request.PostFormValue(`member`),
		request.PostFormValue(`scope`)); err {
	case db.ErrBadOmission:
		util.ServeError(writer,
			util.Error{Description: `id and name are required`},
			http.StatusUnprocessableEntity)

	case db.ErrDupEntry:
		// FIXME: should be StatusMethodNotAllowed if id is duplicate.
		util.ServeError(writer,
			util.Error{Description: `duplicate id or name`},
			http.StatusUnprocessableEntity)

	case db.ErrIncorrectIdentity:
		util.ServeError(writer,
			util.Error{Description: `incorrect member`},
			http.StatusUnprocessableEntity)

	case db.ErrInvalid:
		util.ServeErrorDefault(writer, http.StatusUnprocessableEntity)

	case nil:
		util.ServeJSON(writer, struct{}{}, http.StatusCreated)

	default:
		panic(err)
	}
}

func officersGetServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path != `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	if (authorize(writer, request, shared, scope.Member) == claim{}) {
		return
	}

	util.ServeJSON(writer, shared.DB.QueryOfficers(), http.StatusOK)
}

func officersNamesGetServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path != `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	if (authorize(writer, request, shared, scope.Member) == claim{}) {
		return
	}

	util.ServeJSON(writer, shared.DB.QueryOfficerNames(), http.StatusOK)
}
