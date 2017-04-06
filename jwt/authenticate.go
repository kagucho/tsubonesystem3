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

package jwt

import (
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

// Claim is a structure to hold the processed and validated claim
type Claim struct {
	Sub      string
	Scope    string
	Duration time.Duration
	Tmp      bool
}

// Error is a structure to hold the error and URI for the description.
type Error struct {
	error
	uri string
}

// IsError returns whether it indicates an error or not.
func (jwtError Error) IsError() bool {
	return jwtError.error != nil
}

// URI returns the URI which describes the error.
func (jwtError Error) URI() string {
	return jwtError.uri
}

func (context JWT) validateHeader(encoded string) Error {
	decoder := json.NewDecoder(base64.NewDecoder(
		base64.RawURLEncoding, strings.NewReader(encoded)))

	var decoded header
	if err := decoder.Decode(&decoded); err != nil {
		return Error{err, `https://tools.ietf.org/html/rfc7159`}
	}

	if decoder.More() {
		return Error{
			errors.New(`header contains something superfluous`),
			`https://tools.ietf.org/html/rfc7159`,
		}
	}

	if decoded.Alg != context.authority.Alg() {
		return Error{
			fmt.Errorf(`expected alg "HS256", header says %q`,
				decoded.Alg),
			`https://tools.ietf.org/html/rfc7515#section-4.1.1`,
		}
	}

	return Error{}
}

func validClaim(encoded string) (Claim, Error) {
	decoder := json.NewDecoder(base64.NewDecoder(
		base64.RawURLEncoding, strings.NewReader(encoded)))

	var decoded claim
	if err := decoder.Decode(&decoded); err != nil {
		return Claim{},
			Error{
				err,
				`https://tools.ietf.org/html/rfc7159`,
			}
	}

	if decoder.More() {
		return Claim{},
			Error{
				errors.New(`claim contains something superfluous`),
				`https://tools.ietf.org/html/rfc7159`,
			}
	}

	now := time.Now()
	duration := decoded.Exp.Sub(now)
	if duration < 0 {
		return Claim{},
			Error{
				fmt.Errorf(`claim is expired; it is %v, claim says expired in %v`,
					now, decoded.Exp),
				`https://tools.ietf.org/html/rfc7519#section-4.1.4`,
			}
	}

	return Claim{decoded.Sub, decoded.Scope, duration, decoded.Tmp}, Error{}
}

// Authenticate returns the authenticated claim of the given JWT.
func (context JWT) Authenticate(jwt string) (Claim, Error) {
	splited := strings.Split(jwt, `.`)
	if len(splited) != 3 {
		return Claim{},
			Error{
				fmt.Errorf(`expected 3 parts, got %v parts`, len(splited)),
				`https://tools.ietf.org/html/rfc7519#section-3.1`,
			}
	}

	if invalid := context.validateHeader(splited[0]); invalid.error != nil {
		return Claim{}, invalid
	}

	claim, claimErr := validClaim(splited[1])
	if claimErr.error != nil {
		return Claim{}, claimErr
	}

	signature, decodeErr := base64.RawURLEncoding.DecodeString(splited[2])
	if decodeErr != nil {
		return Claim{},
			Error{
				decodeErr,
				`https://tools.ietf.org/html/rfc4648#section-5`,
			}
	}

	hash := context.authority.Hash()
	io.WriteString(hash, splited[0])
	hash.Write([]byte{'.'})
	io.WriteString(hash, splited[1])

	if !hmac.Equal(hash.Sum(nil), signature) {
		return Claim{},
			Error{
				errors.New(`invalid signature`),
				`https://tools.ietf.org/html/rfc7515#section-5`,
			}
	}

	return claim, Error{}
}
