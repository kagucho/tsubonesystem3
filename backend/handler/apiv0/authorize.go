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
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/token"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/util"
	"github.com/kagucho/tsubonesystem3/backend/scope"
	"net/http"
	"strings"
)

type claim struct {
	sub   string
	scope scope.Scope
	tmp   bool
}

func authorize(writer http.ResponseWriter, request *http.Request, shared shared, scope uint) claim {
	authorization := request.Header.Get("Authorization")
	const prefix = "Bearer "

	if !strings.HasPrefix(authorization, prefix) {
		encodedScope := token.EncodeScopeIndex(scope)
		writer.Header().Set(`WWW-Authenticate`, `Bearer scope=`+encodedScope)
		util.ServeJSON(writer,
			token.Error{
				Error: util.Error{
					ID:          `invalid_token`,
					Description: `expected bearer authentication scheme`,
					URI:         `https://tools.ietf.org/html/rfc6750#section-2.1`,
				},
				Scope: encodedScope,
			}, http.StatusUnauthorized)

		return claim{}
	}

	authenticated, authenticateErr :=
		shared.Token.Authenticate(authorization[len(prefix):])
	if authenticateErr.IsError() {
		token.ServeError(writer,
			token.Error{
				Error: util.Error{
					ID:          `invalid_token`,
					Description: authenticateErr.Error(),
					URI:         authenticateErr.URI(),
				},
				Scope: token.EncodeScopeIndex(scope),
			}, http.StatusUnauthorized)

		return claim{}
	}

	if authenticated.Tmp {
		temporary, queryErr := shared.DB.QueryMemberTmp(authenticated.Sub)
		if queryErr == db.ErrIncorrectIdentity || !temporary {
			token.ServeError(writer,
				token.Error{
					Error: util.Error{
						ID:          `invalid_token`,
						Description: `invalid token`,
						URI:         `https://tools.ietf.org/html/rfc6749#section-7.2`,
					},
					Scope: token.EncodeScopeIndex(scope),
				}, http.StatusUnauthorized)

			return claim{}
		}

		if queryErr != nil {
			panic(queryErr)
		}
	}

	decodedScope, scopeErr := token.DecodeScope(authenticated.Scope)
	if scopeErr != nil {
		token.ServeError(writer,
			token.Error{
				Error: util.Error{
					ID:          `invalid_token`,
					Description: scopeErr.Error(),
					URI:         `https://tools.ietf.org/html/rfc6749#section-7.2`,
				},
				Scope: token.EncodeScopeIndex(scope),
			}, http.StatusUnauthorized)

		return claim{}
	}

	if !decodedScope.IsSet(scope) {
		token.ServeError(writer,
			token.Error{
				Error: util.Error{
					ID:          `insufficient_scope`,
					Description: `The request requires higher privileges than provided by the access token.`,
					URI:         `https://tools.ietf.org/html/rfc6750#section-3.1`,
				},
				Scope: token.EncodeScopeIndex(scope),
			}, http.StatusForbidden)

		return claim{}
	}

	return claim{authenticated.Sub, decodedScope, authenticated.Tmp}
}
