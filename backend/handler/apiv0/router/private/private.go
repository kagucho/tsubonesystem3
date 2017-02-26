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
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/mail"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/member"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/officer"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/party"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/token/authorizer"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/user"
	"github.com/kagucho/tsubonesystem3/backend/scope"
)

var routes = map[string]authorizer.Route{
	`/club/detail`:         {club.DetailServeHTTP, scope.Member},
	`/club/list`:           {club.ListServeHTTP, scope.Member},
	`/mail`:                {mail.ServeHTTP, scope.Member},
	`/member/create`:       {member.CreateServeHTTP, scope.Management},
	`/member/delete`:       {member.DeleteServeHTTP, scope.Management},
	`/member/detail`:       {member.DetailServeHTTP, scope.Member},
	`/member/list`:         {member.ListServeHTTP, scope.Member},
	`/member/listroles`:    {member.ListrolesServeHTTP, scope.Member},
	`/officer/detail`:      {officer.DetailServeHTTP, scope.Member},
	`/officer/list`:        {officer.ListServeHTTP, scope.Member},
	`/party/create`:        {party.CreateServeHTTP, scope.Member},
	`/party/list`:          {party.ListServeHTTP, scope.Member},
	`/party/listnames`:     {party.ListnamesServeHTTP, scope.Member},
	`/party/respond`:       {party.RespondServeHTTP, scope.Member},
	`/user/confirm`:        {user.ConfirmServeHTTP, scope.User},
	`/user/declareob`:      {user.DeclareOBServeHTTP, scope.User},
	`/user/detail`:         {user.DetailServeHTTP, scope.User},
	`/user/update`:         {user.UpdateServeHTTP, scope.User},
	`/user/updatepassword`: {user.UpdatePasswordServeHTTP, scope.User},
}

// GetRoute returns route to the given path.
func GetRoute(path string) authorizer.Route {
	return routes[path]
}
