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

// Package authorizer implements an authorizer based on token specificated by
// RFC 6749 - The OAuth 2.0 Authorization Framework.
// https://tools.ietf.org/html/rfc6749
package authorizer

import (
	"fmt"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/common"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/context"
	tokenScope "github.com/kagucho/tsubonesystem3/backend/handler/apiv0/token/scope"
	"github.com/kagucho/tsubonesystem3/backend/scope"
	"net/http"
	"strings"
)

// Claim is a structure to hold the authorized claim.
type Claim struct {
	Sub   string
	Scope scope.Scope
	Tmp   bool
}

type Route struct {
	Handle func(http.ResponseWriter, *http.Request, context.Context, Claim)
	Scope  uint
}

type oauthError struct {
	common.Error
	Scope string `json:"scope"`
}

func serveError(writer http.ResponseWriter, response oauthError, code int) {
	response.Error = common.Error{
		common.ErrorEncode(response.ID),
		common.ErrorEncode(response.Description),
		common.ErrorEncode(response.URI),
	}

	writer.Header().Set(`WWW-Authenticate`,
		fmt.Sprintf(
			`Bearer error="%s",error_description="%s",error_uri="%s",scope=%s`,
			response.ID, response.Description, response.URI, response.Scope))

	common.ServeJSON(writer, response, code)
}

// Authorize authorizes the appropriate user to the given page according to
// the token included in the request.
func Authorize(writer http.ResponseWriter, request *http.Request, context context.Context, route Route) {
	authorization := request.Header.Get("Authorization")
	const prefix = "Bearer "

	if !strings.HasPrefix(authorization, prefix) {
		scope := tokenScope.Table[route.Scope]
		writer.Header().Set(`WWW-Authenticate`, `Bearer scope=`+scope)
		common.ServeJSON(writer,
			oauthError{
				common.Error{
					`invalid_token`,
					`expected bearer authentication scheme`,
					`https://tools.ietf.org/html/rfc6750#section-2.1`,
				}, tokenScope.Table[route.Scope],
			}, http.StatusUnauthorized)

		return
	}

	claim, authenticateError :=
		context.Token.Authenticate(authorization[len(prefix):])
	if authenticateError.IsError() {
		serveError(writer,
			oauthError{
				common.Error{
					`invalid_token`,
					authenticateError.Error(),
					authenticateError.URI(),
				}, tokenScope.Table[route.Scope],
			}, http.StatusUnauthorized)

		return
	}

	if claim.Tmp {
		temporary, queryError := context.DB.QueryMemberTmp(claim.Sub)
		if queryError != nil {
			serveError(writer,
				oauthError{
					common.Error{
						`invalid_token`,
						`invalid ID`,
						`https://tools.ietf.org/html/rfc6749#section-7.2`,
					}, tokenScope.Table[route.Scope],
				}, http.StatusUnauthorized)

			return
		}

		if !temporary {
			serveError(writer,
				oauthError{
					common.Error{
						`invalid_token`,
						`token is expired`,
						`https://tools.ietf.org/html/rfc6749#section-7.2`,
					}, tokenScope.Table[route.Scope],
				}, http.StatusUnauthorized)

			return
		}
	}

	decodedScope, scopeError := tokenScope.Decode(claim.Scope)
	if scopeError != nil {
		serveError(writer,
			oauthError{
				common.Error{
					`invalid_token`,
					scopeError.Error(),
					`https://tools.ietf.org/html/rfc6749#section-7.2`,
				}, tokenScope.Table[route.Scope],
			}, http.StatusUnauthorized)

		return
	}

	if !decodedScope.IsSet(route.Scope) {
		serveError(writer,
			oauthError{
				common.Error{
					`insufficient_scope`,
					`The request requires higher privileges than provided by the access token.`,
					`https://tools.ietf.org/html/rfc6750#section-3.1`,
				}, tokenScope.Table[route.Scope],
			}, http.StatusForbidden)

		return
	}

	route.Handle(writer, request, context,
		Claim{claim.Sub, decodedScope, claim.Tmp})
}
