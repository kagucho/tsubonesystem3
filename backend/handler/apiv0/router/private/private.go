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

// Package private implements router for private Web pages.
package private

import (
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/club"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/member"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/officer"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/token/authorizer"
	"github.com/kagucho/tsubonesystem3/backend/scope"
)

var routes = map[string]authorizer.Route{
	`/club/detail`:           {club.DetailServeHTTP, scope.Basic},
	`/club/list`:             {club.ListServeHTTP, scope.Basic},
	`/club/listname`:         {club.ListNameServeHTTP, scope.Basic},
	`/member/create`:         {member.CreateServeHTTP, scope.Basic},
	`/member/declareob`:      {member.DeclareOBServeHTTP, scope.Basic},
	`/member/delete`:         {member.DeleteServeHTTP, scope.Management},
	`/member/detail`:         {member.DetailServeHTTP, scope.Basic},
	`/member/list`:           {member.ListServeHTTP, scope.Basic},
	`/member/update`:         {member.UpdateServeHTTP, scope.Basic},
	`/member/updatepassword`: {member.UpdatePasswordServeHTTP, scope.Basic},
	`/officer/detail`:        {officer.DetailServeHTTP, scope.Basic},
	`/officer/list`:          {officer.ListServeHTTP, scope.Basic},
}

// GetRoute returns route to the given path.
func GetRoute(path string) authorizer.Route {
	return routes[path]
}
