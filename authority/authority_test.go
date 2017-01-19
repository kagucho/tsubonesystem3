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

package authority

import "testing"

func (authority Authority) TestAlg(t *testing.T) {
	t.Parallel()

	/*
		RFC 7518 - JSON Web Algorithms (JWA)
		3.2.  HMAC with SHA-2 Functions
		https://tools.ietf.org/html/rfc7518#section-3.2

		> The following "alg" (algorithm) Header Parameter values are
		> used to indicate that the JWS Signature is an HMAC value
		> computed using the corresponding algorithm:
	*/
	if alg := authority.Alg(); alg != `HS256` {
		t.Errorf(`expected "HS256", got %q`, alg)
	}
}

func (authority Authority) TestHash(t *testing.T) {
	t.Parallel()
	authority.Hash()
}

func TestAuthority(t *testing.T) {
	var authority Authority

	if !t.Run(`New`, func(t *testing.T) {
		var authorityError error
		authority, authorityError = New()
		if authorityError != nil {
			t.Error(authorityError)
		}
	}) {
		t.FailNow()
	}

	t.Run(`Alg`, authority.TestAlg)

	t.Run(`Hash`, authority.TestHash)
}
