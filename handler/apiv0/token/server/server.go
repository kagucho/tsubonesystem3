package server

import (
	"database/sql"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/common"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/context"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/token/backend"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/token/scope"
	"log"
	"net/http"
)

type Server struct {
	*limiter
}

type response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope string `json:"scope"`
}

func New() Server {
	return Server{newLimiter()}
}

// ServeHTTP serves tokens.
func (server Server) ServeHTTP(writer http.ResponseWriter, request *http.Request, context context.Context) {
	serve := func() func() {
		defer common.Recover(writer)

		if request.Method != `POST` {
			return func() {
				common.ServeError(writer,
					common.Error{
						Description: `expected 'POST' method request`,
					}, http.StatusMethodNotAllowed)
			}
		}

		var sub string
		var subScope string
		var refresh bool
		switch grantType := request.PostFormValue(`grant_type`); grantType {
		case `password`:
			sub = request.PostFormValue(`username`)
			if !server.challenge(sub) {
				return func() {
					common.ServeErrorDefault(writer,
						http.StatusTooManyRequests)
				}
			}

			subScopeDecoded, scopeError := context.DB.GetScope(
				sub, request.PostFormValue(`password`))
			switch scopeError {
			case nil:
				if subScopeDecoded.IsSetAny() {
					break
				}
				fallthrough

			case sql.ErrNoRows:
				return func() {
					common.ServeError(writer,
						common.Error{
							`invalid_grant`,
							`invalid username and/or password`,
							`https://tools.ietf.org/html/rfc6749#section-5.2`,
						}, http.StatusBadRequest)
				}

			default:
				panic(scopeError)
			}

			subScope, scopeError = scope.Encode(subScopeDecoded)
			if scopeError != nil {
				panic(scopeError)
			}

			refresh = true

		case `refresh_token`:
			claim, authenticateError :=
				context.Token.AuthenticateRefresh(
					request.PostFormValue(`refresh_token`))
			if authenticateError.IsError() {
				return func() {
					common.ServeError(writer,
						common.Error{
							`invalid_grant`,
							authenticateError.Error(),
							authenticateError.URI(),
						}, http.StatusBadRequest)
				}
			}

			if backend.RefreshRequiresRenew(claim) {
				refresh = true
			}

			sub = claim.Sub
			subScope = claim.Scope

		default:
			return func() {
				common.ServeError(writer,
					common.Error{
						`invalid_grant`,
						`expected grant_type 'password' or 'refresh_token'`,
						`https://tools.ietf.org/html/rfc6749#section-5.2`,
					}, http.StatusBadRequest)
			}
		}

		accessToken, tokenError := context.Token.IssueAccess(sub, subScope)
		if tokenError != nil {
			panic(tokenError)
		}

		var refreshToken string
		if refresh {
			refreshToken, tokenError = context.Token.IssueRefresh(sub, subScope)
			if tokenError != nil {
				log.Println(tokenError)
				refreshToken = ``
			}
		} else {
			refreshToken = ``
		}

		return func() {
			common.ServeJSON(writer,
				response{
					AccessToken:  accessToken,
					RefreshToken: refreshToken,
					Scope:        subScope,
				}, http.StatusOK)
		}
	}()

	if serve != nil {
		serve()
	}
}
