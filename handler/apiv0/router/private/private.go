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

// Package private implements router for private Web pages.
package private

import `github.com/kagucho/tsubonesystem3/handler/apiv0/club`
import `github.com/kagucho/tsubonesystem3/handler/apiv0/member`
import `github.com/kagucho/tsubonesystem3/handler/apiv0/officer`

var routes = map[string]Route{
  `/club/detail`: club.DetailServeHTTP,
  `/club/list`: club.ListServeHTTP,
  `/club/listname`: club.ListNameServeHTTP,
  `/member/detail`: member.DetailServeHTTP,
  `/member/list`: member.ListServeHTTP,
  `/officer/detail`: officer.DetailServeHTTP,
  `/officer/list`: officer.ListServeHTTP,
}

// GetRoute returns route to the given path.
func GetRoute(path string) Route {
  return routes[path]
}
