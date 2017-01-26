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

// Package scope implements a decoder and encoder of scope described in
// RFC 6749 - The OAuth 2.0 Authorization Framework.
// https://tools.ietf.org/html/rfc6749
package scope

import (
	"fmt"
	"github.com/kagucho/tsubonesystem3/backend/scope"
	"strings"
)

var Table = []string{
	scope.Update:     `update`,
	scope.Basic:      `basic`,
	scope.Privacy:    `privacy`,
	scope.Management: `management`,
}

// Decode returns the decoded scope.
func Decode(encoded string) (scope.Scope, error) {
	decoded := scope.Scope{}

	if encoded != `` {
	Found:
		for _, splitted := range strings.Split(encoded, ` `) {
			for index, scopeString := range Table {
				if scopeString == splitted {
					decoded = decoded.Set(uint(index))
					continue Found
				}
			}

			return scope.Scope{},
				fmt.Errorf(`unknown scope: %q`, splitted)
		}
	}

	return decoded, nil
}

// Encode returns the encoded scope.
func Encode(decoded scope.Scope) (string, error) {
	scopes := make([]string, 0, len(Table))

	for index, scopeString := range Table {
		if decoded.IsSet(uint(index)) {
			scopes = append(scopes, scopeString)
		}
	}

	return strings.Join(scopes, ` `), nil
}
