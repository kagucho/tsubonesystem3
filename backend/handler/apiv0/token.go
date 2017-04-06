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
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/token/backend"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/util"
	"github.com/kagucho/tsubonesystem3/limiter"
	"net/http"
)

type tokenServer limiter.Limiter

func newTokenServer() *tokenServer {
	return (*tokenServer)(limiter.New())
}

func (server tokenServer) end() {
	limiter.Limiter(server).End()
}

func (server *tokenServer) serveHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.Method != `POST` {
		util.ServeError(writer,
			util.Error{
				Description: `expected 'POST' method request`,
			}, http.StatusMethodNotAllowed)

		return
	}

	var sub string
	var subScope string
	var refresh bool
	switch grantType := request.PostFormValue(`grant_type`); grantType {
	case `password`:
		sub = request.PostFormValue(`username`)
		if !(*limiter.Limiter)(server).Challenge(sub) {
			util.ServeErrorDefault(writer, http.StatusTooManyRequests)

			return
		}

		subScopeDecoded, err := shared.DB.GetScope(
			sub, request.PostFormValue(`password`))
		if err == db.ErrIncorrectIdentity {
			util.ServeError(writer,
				util.Error{
					ID:          `invalid_grant`,
					Description: `invalid username and/or password`,
					URI:         `https://tools.ietf.org/html/rfc6749#section-5.2`,
				}, http.StatusBadRequest)

			return
		} else if err != nil {
			panic(err)
		}

		subScope, err = token.EncodeScope(subScopeDecoded)
		if err != nil {
			panic(err)
		}

		refresh = true

	case `refresh_token`:
		claim, err := shared.Token.AuthenticateRefresh(
			request.PostFormValue(`refresh_token`))
		if err.IsError() {
			util.ServeError(writer,
				util.Error{
					ID:          `invalid_grant`,
					Description: err.Error(),
					URI:         err.URI(),
				}, http.StatusBadRequest)

			return
		}

		if backend.RefreshRequiresRenew(claim) {
			refresh = true
		}

		sub = claim.Sub
		subScope = claim.Scope

	default:
		util.ServeError(writer,
			util.Error{
				ID:          `invalid_grant`,
				Description: `expected grant_type 'password' or 'refresh_token'`,
				URI:         `https://tools.ietf.org/html/rfc6749#section-5.2`,
			}, http.StatusBadRequest)

		return
	}

	accessToken, err := shared.Token.IssueAccess(sub, subScope)
	if err != nil {
		panic(err)
	}

	var refreshToken string
	if refresh {
		refreshToken, err = shared.Token.IssueRefresh(sub, subScope)
		if err != nil {
			refreshToken = ``
		}
	} else {
		refreshToken = ``
	}

	util.ServeJSON(writer,
		struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token,omitempty"`
			Scope        string `json:"scope"`
		}{accessToken, refreshToken, subScope},
		http.StatusOK)
}
