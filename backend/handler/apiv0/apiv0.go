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

// Package apiv0 implements API v0.
package apiv0

import (
	"github.com/kagucho/tsubonesystem3/backend/db"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/common"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/context"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/router/private"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/router/public"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/token/authorizer"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/token/backend"
	"github.com/kagucho/tsubonesystem3/backend/mail"
	"log"
	"net/http"
	"runtime/debug"
)

// APIv0 is a structure to hold the context of API v0.
type APIv0 struct {
	context context.Context
	public  public.Public
}

type apiv0Writer struct {
	http.ResponseWriter
	header http.Header
	written *bool
}

func (writer apiv0Writer) Header() http.Header {
	return writer.header
}

func (writer apiv0Writer) Write(bytes []byte) (int, error) {
	writer.prepare()
	return writer.ResponseWriter.Write(bytes)
}

func (writer apiv0Writer) WriteHeader(code int) {
	writer.prepare()
	writer.ResponseWriter.WriteHeader(code)
}

func (writer apiv0Writer) prepare() {
	if !*writer.written {
		*writer.written = true

		header := writer.ResponseWriter.Header()
		for key, value := range writer.header {
			header[key] = value
		}
	}
}

// New returns a new APIv0.
func New(db db.DB, mail mail.Mail) (APIv0, error) {
	token, tokenError := backend.New()
	if tokenError != nil {
		return APIv0{}, tokenError
	}

	return APIv0{context.Context{db, mail, token}, public.New()}, nil
}

// ServeHTTP servs API v0 via HTTP.
func (apiv0 APIv0) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	written := false
	apiv0Writer := apiv0Writer{writer, http.Header{}, &written}

	defer func() {
		if !written {
			common.ServeErrorDefault(writer, http.StatusInternalServerError)

			log.Print(recover())
			debug.PrintStack()
		}
	}()

	publicRoute := apiv0.public.GetRoute(request.URL.Path)
	if publicRoute != nil {
		publicRoute(apiv0Writer, request, apiv0.context)

		return
	}

	privateRoute := private.GetRoute(request.URL.Path)
	if privateRoute.Handle == nil {
		common.ServeErrorDefault(apiv0Writer, http.StatusNotFound)

		return
	}

	authorizer.Authorize(apiv0Writer, request, apiv0.context, privateRoute)
}
