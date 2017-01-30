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

// Package jwt implements JWT (JSON Web Token).
package jwt

import (
	"hash"
	"math/rand"
	"time"
)

// Authority is the interface for the authority of JWS (JSON Web Signature).
type Authority interface {
	Hash() hash.Hash
	Alg() string
}

// JWT is the structure to hold the context of the JWT issuer and signer.
type JWT struct {
	authority Authority
}

func init() {
	// math/rand is used for Jti.
	rand.Seed(time.Now().Unix())
}

// New returns a new jwt.JWT.
func New(authority Authority) JWT {
	return JWT{authority}
}

type header struct {
	Alg string `json:"alg"`
}

type claim struct {
	Sub   string `json:"sub"`
	Scope string `json:"scope",omitempty`
	Exp   int64  `json:"exp",omitempty`
	Tmp   bool   `json:"tmp",omitempty`
	Jti   string `json:"jti"`
}
