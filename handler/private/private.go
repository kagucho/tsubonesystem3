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

// Package private implements the private page hosting.
package private

import (
	"bytes"
	"github.com/kagucho/tsubonesystem3/db"
	"github.com/kagucho/tsubonesystem3/handler/file"
	"html/template"
	"log"
	"net/http"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
)

// Private is a structure to hold the context of the private page hosting.
type Private struct {
	file      string
	graph     *template.Template
	db        db.DB
	fileError file.FileError
}

// New returns a new private.Private.
func New(share string, db db.DB, fileError file.FileError) (Private, error) {
	graph, parseError := template.ParseFiles(path.Join(share, `graph`))
	if parseError != nil {
		return Private{}, parseError
	}

	return Private{path.Join(share, `public/private`), graph, db, fileError},
		nil
}

func (private Private) ServeHTTP(writer http.ResponseWriter,
	request *http.Request) {
	serve := func() {
		private.fileError.ServeError(writer, http.StatusInternalServerError)
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			log.Println(recovered)
			debug.PrintStack()
		}

		serve()
	}()

	switch request.URL.Path {
	case `/private`:
		request.ParseForm()
		writer.Header().Set(`Content-Language`, `ja`)

		escapedFragment := request.Form[`_escaped_fragment_`]
		if len(escapedFragment) == 0 {
			if request.URL.RawQuery == `` {
				serve = func() {
					http.ServeFile(writer, request, private.file)
				}
			} else {
				request.URL.RawQuery = ``
				serve = func() {
					private.fileError.ServeMovedPermanently(
						writer, request.URL.String())
				}
			}
		} else {
			routeQuery := strings.SplitN(escapedFragment[0], `?`, 2)

			request.URL.RawQuery = ``
			base := request.URL.String()

			graphFunc := graphFuncs[routeQuery[0]]
			if graphFunc == nil {
				graphFunc = graphDefault
			}

			var buffer bytes.Buffer
			if executeError := private.graph.Execute(
				&buffer, graphFunc(private.db, base, routeQuery)); executeError != nil {
				panic(executeError)
			}

			header := writer.Header()
			header.Set(`Content-Length`, strconv.Itoa(buffer.Len()))

			serve = func() {
				buffer.WriteTo(writer)
			}
		}

	case `/private/`:
		request.URL.Path = `/private`
		request.URL.RawQuery = ``
		serve = func() {
			private.fileError.ServeMovedPermanently(
				writer, request.URL.String())
		}

	default:
		serve = func() {
			private.fileError.ServeError(writer, http.StatusNotFound)
		}
	}
}
