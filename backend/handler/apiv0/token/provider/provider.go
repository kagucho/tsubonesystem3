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

// Package provider implements the token provider of API v0.
package provider

import (
	"github.com/kagucho/tsubonesystem3/authority"
	"github.com/kagucho/tsubonesystem3/jwt"
)

// New returns a new jwt.JWT.
func New() (jwt.JWT, error) {
	authority, authorityError := authority.New()
	if authorityError != nil {
		return jwt.JWT{}, authorityError
	}

	return jwt.New(authority), nil
}
