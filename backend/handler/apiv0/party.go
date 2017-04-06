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
	"fmt"
	"github.com/kagucho/tsubonesystem3/backend/db"
	"github.com/kagucho/tsubonesystem3/backend/encoding"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/util"
	"github.com/kagucho/tsubonesystem3/backend/scope"
	"net/http"
	"strconv"
)

func partyDeleteServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path == `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	authorized := authorize(writer, request, shared, scope.Member)
	if (authorized == claim{}) {
		return
	}

	switch err := shared.DB.DeleteParty(request.URL.Path[1:], authorized.sub); err {
	case db.ErrIncorrectIdentity:
		util.ServeErrorDefault(writer, http.StatusNotFound)

	case nil:
		util.ServeJSON(writer, struct{}{}, http.StatusOK)

	default:
		panic(err)
	}
}

func partyGetServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path == `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
	}

	if (authorize(writer, request, shared, scope.Member) == claim{}) {
		return
	}

	switch party, err := shared.DB.QueryParty(request.URL.Path[1:]); err {
	case db.ErrIncorrectIdentity:
		util.ServeErrorDefault(writer, http.StatusNotFound)

	case nil:
		util.ServeJSON(writer, party, http.StatusOK)

	default:
		panic(err)
	}
}

func partyPatchServeHTTP(writer http.ResponseWriter, request *http.Request,
	shared shared) {
	if request.URL.Path == `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	authorized := authorize(writer, request, shared, scope.Member)
	if (authorized == claim{}) {
		return
	}

	var start encoding.Time
	var end encoding.Time
	var due encoding.Time
	party := request.URL.Path[1:]
	place := request.FormValue(`place`)
	inviteds := request.FormValue(`inviteds`)
	invitedIDs := request.FormValue(`invited_ids`)
	details := request.FormValue(`details`)

	if startString := request.FormValue(`start`); startString != `` {
		var err error

		start, err = encoding.ParseQueryTime(startString)
		if err != nil {
			description := `syntax error in start`
			status := http.StatusBadRequest

			if err == strconv.ErrRange {
				description = `start out of range`
				status = http.StatusUnprocessableEntity
			}

			util.ServeError(writer,
				util.Error{Description: description},
				status)

			return
		}
	}

	if endString := request.FormValue(`end`); endString != `` {
		var err error

		end, err = encoding.ParseQueryTime(endString)
		if err != nil {
			description := `syntax error in end`
			status := http.StatusBadRequest

			if err == strconv.ErrRange {
				description = `end out of range`
				status = http.StatusUnprocessableEntity
			}

			util.ServeError(writer,
				util.Error{Description: description},
				status)

			return
		}
	}

	if dueString := request.FormValue(`due`); dueString != `` {
		var err error

		due, err = encoding.ParseQueryTime(dueString)
		if err != nil {
			description := `syntax error in due`
			status := http.StatusBadRequest

			if err == strconv.ErrRange {
				description = `due out of range`
				status = http.StatusUnprocessableEntity
			}

			util.ServeError(writer,
				util.Error{Description: description},
				status)

			return
		}
	}

	if (start != encoding.Time{} || end != encoding.Time{} || place != `` || due != encoding.Time{} || inviteds != `` || invitedIDs != `` || details != ``) {
		switch err := shared.DB.UpdateParty(
			party, authorized.sub, start, end,
			place, due, inviteds, invitedIDs,
			details); err {
		case db.ErrIncorrectIdentity:
			/*
				FIXME: should be StatusUnprocessableEntity if
				party is found but members are not.
			*/
			util.ServeErrorDefault(writer, http.StatusNotFound)
			return

		case db.ErrInvalid:
			util.ServeErrorDefault(writer, http.StatusUnprocessableEntity)
			return

		case nil:

		default:
			panic(err)
		}
	}

	if attendingString := request.FormValue(`attending`); attendingString != `` {
		var attending bool
		switch ; attendingString {
		case `0`:
			attending = false

		case `1`:
			attending = true

		default:
			util.ServeError(writer,
				util.Error{Description: fmt.Sprintf(`invalid attending: '%v'`, attending)},
				http.StatusBadRequest)

			return
		}

		switch err := shared.DB.UpdateAttendance(attending,
			party, authorized.sub); err {
		case db.ErrIncorrectIdentity:
			/*
				FIXME: should be StatusUnprocessableEntity if
				party is found but creator is not.
			*/
			util.ServeErrorDefault(writer, http.StatusNotFound)
			return

		case nil:

		default:
			panic(err)
		}
	}

	util.ServeJSON(writer, struct{}{}, http.StatusOK)
}

func partyPutServeHTTP(writer http.ResponseWriter, request *http.Request,
	shared shared) {
	if len(request.URL.Path) < 2 {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	authorized := authorize(writer, request, shared, scope.Member)
	if (authorized == claim{}) {
		return
	}

	start, startErr := encoding.ParseQueryTime(request.PostFormValue(`start`))
	if startErr != nil {
		description := `syntax error in start`
		status := http.StatusBadRequest

		if startErr == strconv.ErrRange {
			description = `start out of range`
			status = http.StatusUnprocessableEntity
		}

		util.ServeError(writer,
			util.Error{Description: description},
			status)

		return
	}

	end, endErr := encoding.ParseQueryTime(request.PostFormValue(`end`))
	if endErr != nil {
		description := `syntax error in end`
		status := http.StatusBadRequest

		if endErr == strconv.ErrRange {
			description = `end out of range`
			status = http.StatusUnprocessableEntity
		}

		util.ServeError(writer,
			util.Error{Description: description},
			status)

		return
	}

	due, dueErr := encoding.ParseQueryTime(request.PostFormValue(`due`))
	if dueErr != nil {
		description := `syntax error in due`
		status := http.StatusBadRequest

		if endErr == strconv.ErrRange {
			description = `due out of range`
			status = http.StatusUnprocessableEntity
		}

		util.ServeError(writer,
			util.Error{Description: description},
			status)

		return
	}

	name := request.URL.Path[1:]
	place := request.PostFormValue(`place`)
	invitedIDs := request.PostFormValue(`invited_ids`)
	inviteds := request.PostFormValue(`inviteds`)
	details := request.PostFormValue(`details`)

	switch mails, err := shared.DB.InsertParty(name,
		authorized.sub, start, end, place,
		due, invitedIDs, inviteds, details); err {
	case db.ErrBadOmission:
		util.ServeError(writer,
			util.Error{Description: `place, invited_ids, inviteds, and details are required`},
			http.StatusUnprocessableEntity)

	case db.ErrDupEntry:
		// FIXME: should be StatusMethodNotAllowed if id is duplicate.
		util.ServeError(writer,
			util.Error{Description: `duplicate name`},
			http.StatusUnprocessableEntity)

	case db.ErrIncorrectIdentity:
		util.ServeError(writer,
			util.Error{Description: `unknown inviteds`},
			http.StatusUnprocessableEntity)

	case db.ErrInvalid:
		util.ServeError(writer,
			util.Error{Description: `invalid request`},
			http.StatusUnprocessableEntity)

	case nil:
		err := shared.Mail.SendInvitations(request.Host, mails, name, start, end, place, inviteds, due, details)
		if err != nil {
			util.ServeMailError(writer)

			return
		}

		util.ServeJSON(writer, struct{}{}, http.StatusCreated)

	default:
		panic(err)
	}
}

func partiesGetServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path != `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	authorized := authorize(writer, request, shared, scope.Member)
	if (authorized == claim{}) {
		return
	}

	util.ServeJSON(writer,
		shared.DB.QueryParties(authorized.sub), http.StatusOK)
}
