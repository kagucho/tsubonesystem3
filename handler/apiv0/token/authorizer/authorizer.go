/*
  Copyright (C) 2016  Kagucho <kagucho.net@gmail.com>

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU Affero General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU Affero General Public License for more details.

  You should have received a copy of the GNU Affero General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

// Package authorizer implements a token authorizer.
package authorizer

import (
  `fmt`
  `github.com/kagucho/tsubonesystem3/handler/apiv0/common`
  tokenScope `github.com/kagucho/tsubonesystem3/handler/apiv0/token/scope`
  `github.com/kagucho/tsubonesystem3/db`
  `github.com/kagucho/tsubonesystem3/scope`
  `github.com/kagucho/tsubonesystem3/jwt`
  `net/http`
  `strings`
)

// Authorizer is a structure to hold the context of the token authorizer.
type Authorizer struct {
  jwt *jwt.JWT
}

// Claim is a structure to hold the authorized claim.
type Claim struct {
  Sub string
  Scope scope.Scope
}

// New returns a new authorizer.Authorizer.
func New(jwt *jwt.JWT) Authorizer {
  return Authorizer{jwt}
}

func escape(unescaped string) string {
  return strings.Replace(strings.Replace(unescaped, `\`, `\\`, -1),
                         `"`, `\"`, -1)
}

func serveErrorBody(writer http.ResponseWriter, id string, description string,
                    uri string, code int) {
  common.ServeJSON(writer,
                   struct{
                     common.Error
                     Scope string `json:"scope"`
                   }{common.Error{id, description, uri}, `basic`},
                   code)
}

func serveError(writer http.ResponseWriter, id string, description string,
                uri string, code int) {
  writer.Header().Set(`WWW-Authenticate`,
    fmt.Sprintf(
      `Bearer error="%s",error_description="%s",error_uri="%s",scope=basic`,
      escape(id), escape(description), escape(uri)))

  serveErrorBody(writer, id, description, uri, code)
}

// Authorize authorizes the appropriate user to the given page according to
// the token included in the request.
func (authorizer Authorizer) Authorize(
  writer http.ResponseWriter, request *http.Request, db db.DB,
  handler func(writer http.ResponseWriter, request *http.Request, db db.DB,
               claim Claim)) {
  serve := func() func() {
    defer common.Recover(writer)

    authorization := request.Header.Get("Authorization")
    const prefix = "Bearer "
    if !strings.HasPrefix(authorization, prefix) {
      writer.Header().Set(`WWW-Authenticate`, `Bearer scope=basic`)
      return func() {
        serveErrorBody(writer, `invalid_token`,
                       fmt.Sprintf(`expected bearer authentication scheme, got "%s" in Authorization field of request header`,
                                   authorization),
                       `https://tools.ietf.org/html/rfc6750#section-2.1`,
                       http.StatusUnauthorized)
      }
    }

    claim, authenticateError :=
      authorizer.jwt.Authenticate(authorization[len(prefix):])
    if authenticateError.IsError() {
      return func() {
        serveError(writer, `invalid_token`, authenticateError.Error(),
                   authenticateError.URI(), http.StatusUnauthorized)
      }
    }

    decodedScope, scopeError := tokenScope.Decode(claim.Scope)
    if scopeError != nil {
      return func() {
        serveError(writer, `invalid_token`, scopeError.Error(),
                   `https://tools.ietf.org/html/rfc6749#section-7.2`,
                   http.StatusUnauthorized)
      }
    }

    if !decodedScope.IsSet(scope.Basic) {
      return func() {
        serveError(writer, `insufficient_scope`,
                   `The request requires higher privileges than provided by the access token.`,
                   `https://tools.ietf.org/html/rfc6750#section-3.1`,
                   http.StatusForbidden)
      }
    }

    return func() {
      handler(writer, request, db, Claim{claim.Sub, decodedScope})
    }
  }()

  if serve != nil {
    serve()
  }
}
