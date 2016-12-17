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

/*
  Package authority implements HMAC-SHA256 authority.
  This package is kept minimal to ensure that the key is safe.
 */
package authority

import (
  `crypto/hmac`
  `crypto/rand`
  `crypto/sha256`
  `hash`
)

// Authority is the structure to store the key safely.
type Authority struct {
  key [sha256.Size]byte
}

/*
  New returns a new authority.Authority initialized with a cryptographically
  random key.
 */
func New() (Authority, error) {
  var authority Authority
  _, randError := rand.Read(authority.key[:])
  return authority, randError
}

/*
  Alg returns the identifier of the algorithm as descibed in the JSON Web
  Algorithms (JWA) specification.
 */
func (authority Authority) Alg() string {
  return `HS256`
}

// Hash returns hash.Hash initialized with the persistent key.
func (authority Authority) Hash() hash.Hash {
  return hmac.New(sha256.New, authority.key[:])
}
