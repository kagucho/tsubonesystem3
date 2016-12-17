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

// Package public implements a router for the public Web pages.
package public

import `github.com/kagucho/tsubonesystem3/handler/apiv0/token/server`

// Public is a structure to hold the context of the public router.
type Public struct {
  routes map[string]Route
}

// New returns a new Public.
func New(token server.Server) Public {
  return Public{map[string]Route{`/token`: token.ServeHTTP}}
}

// GetRoute returns the route for the given path.
func (public Public) GetRoute(path string) Route {
  return public.routes[path]
}
