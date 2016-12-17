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

// Package server implements the token server of API v0.
package server

import (
  `database/sql`
  `fmt`
  `github.com/kagucho/tsubonesystem3/db`
  `github.com/kagucho/tsubonesystem3/handler/apiv0/common`
  `github.com/kagucho/tsubonesystem3/handler/apiv0/token/provider`
  `github.com/kagucho/tsubonesystem3/handler/apiv0/token/scope`
  `github.com/kagucho/tsubonesystem3/jwt`
  `log`
  `net/http`
  `sync`
  `time`
)

type rate struct {
  count uint
  timer *time.Timer
}

type limiter struct {
  rate map[string]*rate
  mutex sync.Mutex
}

type response struct {
  AccessToken string `json:"access_token"`
  RefreshToken string `json:"refresh_token,omitempty"`
}

// Server is a structure to hold the context of the token server.
type Server struct {
  limiter limiter
  access *jwt.JWT
  refresh jwt.JWT
}

func newLimiter() limiter {
  return limiter{rate: map[string]*rate{}}
}

func (limiter limiter) challenge(username string) bool {
  const duration = 68719476736
  const limit = 8

  limiter.mutex.Lock()

  var accept bool
  if rateEntry, present := limiter.rate[username]; present {
    if rateEntry.count < limit {
      if !rateEntry.timer.Stop() {
        limiter.mutex.Unlock()
        return limiter.challenge(username)
      }

      rateEntry.count++
      rateEntry.timer.Reset(duration)

      accept = true
    } else {
      accept = false
    }
  } else {
    var rateEntry rate
    limiter.rate[username] = &rateEntry
    rateEntry.timer = time.AfterFunc(duration, func() {
      limiter.mutex.Lock()
      delete(limiter.rate, username)
      limiter.mutex.Unlock()
    })

    accept = true
  }

  limiter.mutex.Unlock()

  return accept
}

// New returns a new server.Server.
func New(access *jwt.JWT) (Server, error) {
  refresh, refreshError := provider.New()
  if refreshError != nil {
    return Server{}, refreshError
  }

  return Server{ newLimiter(), access, refresh }, nil
}

func (server Server) ServeHTTP(writer http.ResponseWriter,
                               request *http.Request,
                               db db.DB) {
  serve := func() func() {
    const accessTokenDuration = 2199023255552
    const refreshTokenDuration = 70368744177664
    const refreshDuration = accessTokenDuration * 2

    defer common.Recover(writer)

    if request.Method != `POST` {
      return func() {
        common.ServeError(writer, ``,
          fmt.Sprintf(`expected "POST" method request, got %q method request`,
                      request.Method),
          ``, http.StatusMethodNotAllowed)
      }
    }

    var sub string
    var subScope string
    var refresh bool
    switch grantType := request.PostFormValue(`grant_type`); grantType {
    case `password`:
      sub = request.PostFormValue(`username`)
      if !server.limiter.challenge(sub) {
        return func() {
          common.ServeErrorDefault(writer, http.StatusTooManyRequests)
        }
      }

      subScopeDecoded, scopeError := db.GetScope(
        sub, request.PostFormValue(`password`))
      switch scopeError {
      case nil:
        if subScopeDecoded.IsSetAny() {
          break
        }
        fallthrough

      case sql.ErrNoRows:
        return func() {
          common.ServeError(writer, `invalid_grant`,
                            `invalid username and/or password`,
                            `https://tools.ietf.org/html/rfc6749#section-5.2`,
                            http.StatusBadRequest)
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
        server.refresh.Authenticate(request.PostFormValue(`refresh_token`))
      if authenticateError.IsError() {
        return func() {
          common.ServeError(writer, `invalid_grant`, authenticateError.Error(),
                            authenticateError.URI(), http.StatusBadRequest)
        }
      }

      if claim.Duration < refreshDuration {
        refresh = true
      }

      sub = claim.Sub
      subScope = claim.Scope

    default:
      return func() {
        common.ServeError(writer, `invalid_grant`,
          fmt.Sprintf(
            `expected grant_type "password" or "refresh_token", got %q`,
            grantType),
          `https://tools.ietf.org/html/rfc6749#section-5.2`,
          http.StatusBadRequest)
      }
    }

    accessToken, tokenError := server.access.Issue(sub, subScope,
                                                   accessTokenDuration)
    if tokenError != nil {
      panic(tokenError)
    }

    var refreshToken string
    if refresh {
      refreshToken, tokenError = server.refresh.Issue(sub, subScope,
                                                      refreshTokenDuration)
      if tokenError != nil {
        log.Println(tokenError)
        refreshToken = ``
      }
    } else {
      refreshToken = ``
    }

    return func() {
      common.ServeJSON(writer, response{accessToken, refreshToken},
                       http.StatusOK)
    }
  }()

  if serve != nil {
    serve()
  }
}
