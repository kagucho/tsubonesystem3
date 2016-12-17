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

package jwt

import (
  `encoding/base64`
  `encoding/json`
  `github.com/kagucho/tsubonesystem3/authority`
  `testing`
  `time`
)

func testHeader(t *testing.T, encoded string, alg string) {
  var decoded []byte
  decoded, decodeError := base64.RawURLEncoding.DecodeString(encoded)
  if decodeError != nil {
    t.Error(decodeError)
    return
  }

  var unmarshaled header
  unmarshalError := json.Unmarshal(decoded, &unmarshaled)
  if unmarshalError != nil {
    t.Error(unmarshalError)
    return
  }

  if unmarshaled.Alg != alg {
    t.Errorf(`invalid alg in header; expected %q, got %q`,
             alg, unmarshaled.Alg)
  }
}

func testClaim(t *testing.T, encoded string, sub string, scope string,
               before int64, after int64) {
  decoded, decodeError := base64.RawURLEncoding.DecodeString(encoded)
  if decodeError != nil {
    t.Error(decodeError)
    return
  }

  var unmarshaled claim
  unmarshalError := json.Unmarshal(decoded, &unmarshaled)
  if unmarshalError != nil {
    t.Error(unmarshalError)
    return
  }

  if unmarshaled.Sub != sub {
    t.Errorf(`invalid typ in header; expected %q, got %q`,
             sub, unmarshaled.Sub)
  }

  if unmarshaled.Scope != scope {
    t.Errorf(`invalid scope in header; expected %q, got %q`,
             scope, unmarshaled.Scope)
  }

  if unmarshaled.Exp < before || after < unmarshaled.Exp {
    t.Errorf(`invalid Exp in header; expected an hour later, got %q`,
             time.Unix(unmarshaled.Exp, 0))
  }
}

func TestJWT(t *testing.T) {
  var jwt JWT

  if !t.Run(`New`, func(t *testing.T) {
    authority, authorityError := authority.New()
    if authorityError != nil {
      t.Fatal(authorityError)
    }

    jwt = New(authority)
  }) {
    t.FailNow()
  }

  t.Run(`Authenticate`, jwt.testAuthenticate)
  t.Run(`Issue`, jwt.testIssue)
}
