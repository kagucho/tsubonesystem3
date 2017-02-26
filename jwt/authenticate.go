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
	Sub       string
	Scope     string
	Duration  time.Duration
	Tmp       bool
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
	if decodeError := decoder.Decode(&decoded); decodeError != nil {
		return Error{decodeError, `https://tools.ietf.org/html/rfc7159`}
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
	if decodeError := decoder.Decode(&decoded); decodeError != nil {
		return Claim{},
			Error{
				decodeError,
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

	exp := time.Unix(decoded.Exp, 0)
	now := time.Now()
	duration := exp.Sub(now)
	if duration < 0 {
		return Claim{},
			Error{
				fmt.Errorf(`claim is expired; it is %v, claim says expired in %v`,
					now, exp),
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

	claim, claimError := validClaim(splited[1])
	if claimError.error != nil {
		return Claim{}, claimError
	}

	signature, decodeError := base64.RawURLEncoding.DecodeString(splited[2])
	if decodeError != nil {
		return Claim{},
			Error{
				decodeError,
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
