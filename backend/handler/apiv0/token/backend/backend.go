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

// Package backend implements the backend of the token server of API v0.
package backend

import (
	"github.com/kagucho/tsubonesystem3/authority"
	"github.com/kagucho/tsubonesystem3/jwt"
)

// Backend is a structure to hold the context of the token backend.
type Backend struct {
	access  jwt.JWT
	mail    jwt.JWT
	refresh jwt.JWT
}

const accessTokenDuration = 2199023255552
const refreshTokenDuration = 70368744177664
const refreshDuration = accessTokenDuration * 2
const tmpDuration = 70368744177664

func newJWT() (jwt.JWT, error) {
	authority, err := authority.New()
	if err != nil {
		return jwt.JWT{}, err
	}

	return jwt.New(authority), nil
}

// New returns a new backend.Backend.
func New() (Backend, error) {
	access, accessErr := newJWT()
	if accessErr != nil {
		return Backend{}, accessErr
	}

	mail, mailErr := newJWT()
	if mailErr != nil {
		return Backend{}, mailErr
	}

	refresh, refreshErr := newJWT()
	if refreshErr != nil {
		return Backend{}, refreshErr
	}

	return Backend{access, mail, refresh}, nil
}

// Authenticate returns a claim authenticated with the given access token.
func (backend Backend) Authenticate(token string) (jwt.Claim, jwt.Error) {
	return backend.access.Authenticate(token)
}

/*
AuthenticateMail returns a claim authenticated with the given token embedded
in an email.
*/
func (backend Backend) AuthenticateMail(token string) (jwt.Claim, jwt.Error) {
	return backend.mail.Authenticate(token)
}

/*
AuthenticateRefresh returns a claim authenticated with the given refresh token.
*/
func (backend Backend) AuthenticateRefresh(token string) (jwt.Claim, jwt.Error) {
	return backend.refresh.Authenticate(token)
}

/*
RefreshRequiresRenew returns a bool telling the refresh token is required to
renew.
*/
func RefreshRequiresRenew(claim jwt.Claim) bool {
	return claim.Duration < refreshDuration
}

// IssueMail returns a token to embed in an email.
func (backend Backend) IssueMail(sub string) (string, error) {
	return backend.mail.Issue(sub, ``, tmpDuration, false)
}

/*
IssueTmpUserAccess returns an access token which is valid until the user fills
his information, including credentials.
*/
func (backend Backend) IssueTmpUserAccess(sub string) (string, error) {
	return backend.access.Issue(sub, `user`, tmpDuration, true)
}

// IssueAccess returns an access token.
func (backend Backend) IssueAccess(sub, scope string) (string, error) {
	return backend.access.Issue(sub, scope, accessTokenDuration, false)
}

// IssueRefresh returns a refresh token.
func (backend Backend) IssueRefresh(sub, scope string) (string, error) {
	return backend.refresh.Issue(sub, scope, refreshTokenDuration, false)
}
