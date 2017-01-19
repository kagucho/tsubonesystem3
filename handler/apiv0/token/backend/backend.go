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

// Package server implements the token server of API v0.
//
// The implementation comforms to the specification of authorization server
// in RFC 6749 - The OAuth 2.0 Authorization Framework.
// https://tools.ietf.org/html/rfc6749#section-1.1
package backend

import (
	"github.com/kagucho/tsubonesystem3/handler/apiv0/token/provider"
	"github.com/kagucho/tsubonesystem3/jwt"
)

// Backend is a structure to hold the context of the token backend.
type Backend struct {
	access  jwt.JWT
	refresh jwt.JWT
}

const accessTokenDuration = 2199023255552
const refreshTokenDuration = 70368744177664
const refreshDuration = accessTokenDuration * 2

// New returns a new backend.Backend.
func New() (Backend, error) {
	access, accessError := provider.New()
	if accessError != nil {
		return Backend{}, accessError
	}

	refresh, refreshError := provider.New()
	if refreshError != nil {
		return Backend{}, refreshError
	}

	return Backend{access, refresh}, nil
}

func (backend Backend) Authenticate(token string) (jwt.Claim, jwt.Error) {
	return backend.access.Authenticate(token)
}

func (backend Backend) AuthenticateRefresh(token string) (jwt.Claim, jwt.Error) {
	return backend.refresh.Authenticate(token)
}

func RefreshRequiresRenew(claim jwt.Claim) bool {
	return claim.Duration < refreshDuration
}

func (backend Backend) IssueTemporaryAccessUpdater(sub string) (string, error) {
	return backend.access.Issue(sub, "update", 0, true)
}

func (backend Backend) IssueAccess(sub, scope string) (string, error) {
	return backend.access.Issue(sub, scope, accessTokenDuration, false)
}

func (backend Backend) IssueRefresh(sub, scope string) (string, error) {
	return backend.refresh.Issue(sub, scope, refreshTokenDuration, false)
}
