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

// Package token implements the token context of API v0.
package token

import (
  `github.com/kagucho/tsubonesystem3/handler/apiv0/token/authorizer`
  `github.com/kagucho/tsubonesystem3/handler/apiv0/token/provider`
  `github.com/kagucho/tsubonesystem3/handler/apiv0/token/server`
)

// New returns the token server and the token authorizer.
func New() (server.Server, authorizer.Authorizer, error) {
  token, tokenError := provider.New()
  if tokenError != nil {
    return server.Server{}, authorizer.Authorizer{}, tokenError
  }

  serverInstance, serverError := server.New(&token)
  if serverError != nil {
    return server.Server{}, authorizer.Authorizer{}, serverError
  }

  return serverInstance, authorizer.New(&token), nil
}
