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
	"github.com/kagucho/tsubonesystem3/db"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/common"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/context"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/router/private"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/router/public"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/token/authorizer"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/token/backend"
	"github.com/kagucho/tsubonesystem3/mail"
	"net/http"
)

// Apiv0 is a structure to hold the context of API v0.
type Apiv0 struct {
	context    context.Context
	public     public.Public
}

// New returns a new Apiv0.
func New(db db.DB, mail mail.Mail) (Apiv0, error) {
	token, tokenError := backend.New()
	if tokenError != nil {
		return Apiv0{}, tokenError
	}

	return Apiv0{context.Context{db, mail, token}, public.New()}, nil
}

// ServeHTTP servs API v0 via HTTP.
func (apiv0 Apiv0) ServeHTTP(writer http.ResponseWriter,
	request *http.Request) {
	route := func() func() {
		defer common.Recover(writer)

		publicRoute := apiv0.public.GetRoute(request.URL.Path)
		if publicRoute != nil {
			return func() {
				publicRoute(writer, request, apiv0.context)
			}
		}

		privateRoute := private.GetRoute(request.URL.Path)
		if privateRoute.Handle == nil {
			return func() {
				common.ServeErrorDefault(writer, http.StatusNotFound)
			}
		}

		return func() {
			authorizer.Authorize(
				writer, request, apiv0.context, privateRoute)
		}
	}()

	if route != nil {
		route()
	}
}
