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

// Package token implements a some utility related to token.
package token

import (
	"fmt"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/util"
	"github.com/kagucho/tsubonesystem3/backend/scope"
	"net/http"
	"strings"
)

/*
Error is a strucutre representing a error when authenticating with token,
with the expected scope.
*/
type Error struct {
	util.Error
	Scope string `json:"scope"`
}

var table = []string{
	scope.Management: `management`,
	scope.Member:     `member`,
	scope.Privacy:    `privacy`,
	scope.User:       `user`,
}

// DecodeScope returns the decoded scope.
func DecodeScope(encoded string) (scope.Scope, error) {
	decoded := scope.Scope{}

	if encoded != `` {
	Found:
		for _, splitted := range strings.Split(encoded, ` `) {
			for index, scopeString := range table {
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

// EncodeScope returns the encoded scope.
func EncodeScope(decoded scope.Scope) (string, error) {
	scopes := make([]string, 0, len(table))

	for index, scopeString := range table {
		if decoded.IsSet(uint(index)) {
			scopes = append(scopes, scopeString)
		}
	}

	return strings.Join(scopes, ` `), nil
}

// EncodeScopeIndex encodes the scope identified by the given index.
func EncodeScopeIndex(decoded uint) string {
	return table[decoded]
}

// ServeError serves an error when authenticating with the bearer token.
func ServeError(writer http.ResponseWriter, response Error, code int) {
	response.Error = util.Error{
		ID:          util.ErrorEncode(response.ID),
		Description: util.ErrorEncode(response.Description),
		URI:         util.ErrorEncode(response.URI),
	}

	writer.Header().Set(`WWW-Authenticate`,
		fmt.Sprintf(
			`Bearer error="%s",error_description="%s",error_uri="%s",scope=%s`,
			response.ID, response.Description, response.URI, response.Scope))

	util.ServeJSON(writer, response, code)
}
