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

// Package common implements common functions for API v0.
package common

import (
	"encoding/json"
	"net/http"
)

// Error is a structure to hold an error to serve.
//
// The JSON-encoded structure is conforming to RFC 6749 - The OAuth 2.0
// Authorization Framework 5.2.  Error Response.
// https://tools.ietf.org/html/rfc6749#section-5.2
type Error struct {
	ID          string `json:"error",omitempty`
	Description string `json:"error_description",omitempty`
	URI         string `json:"error_uri",omitempty`
}

var statusError = map[int]struct {
	id  string
	uri string
}{
	http.StatusBadRequest: {
		`invalid_request`,
		`https://tools.ietf.org/html/rfc7231#section-6.5.1`,
	},
	http.StatusNotFound: {
		`not_found`,
		`https://tools.ietf.org/html/rfc7231#section-6.5.4`,
	},
	http.StatusMethodNotAllowed: {
		`unsupported_response_type`,
		`https://tools.ietf.org/html/rfc7231#section-6.5.5`,
	},
	http.StatusInternalServerError: {
		`server_error`,
		`https://tools.ietf.org/html/rfc7231#section-6.6.1`,
	},
	http.StatusTooManyRequests: {
		`too_many_requests`,
		`https://tools.ietf.org/html/rfc6585#section-4`,
	},
}

// ErrorEncode encodes the error to a string conforming to
// RFC 6749 - The OAuth 2.0 Authorization Framework 5.2.  Error Response
// https://tools.ietf.org/html/rfc6749#section-5.2
func ErrorEncode(decoded string) string {
	const hex = `0123456789ABCDEF`
	encoded := make([]byte, 0, len(decoded)*3)

	for index := 0; index < len(decoded); index++ {
		if decoded[index] < 0x20 || decoded[index] == 0x22 ||
			decoded[index] == '%' ||
			decoded[index] == 0x5c || decoded[index] > 0x7e {
			encoded = append(encoded, '%',
				hex[decoded[index]>>4],
				hex[decoded[index]&15])
		} else {
			encoded = append(encoded, decoded[index])
		}
	}

	return string(encoded)
}

// ServeError serves an error according to the given arguments.
// The error response is conforming to RFC 6749  - The OAuth 2.0
// Authorization Framework 5.2.  Error Response.
// https://tools.ietf.org/html/rfc6749#section-4.2.2.1
func ServeError(writer http.ResponseWriter, response Error, code int) {
	if response.ID == `` {
		response.ID = statusError[code].id
	} else {
		response.ID = ErrorEncode(response.ID)
	}

	if response.Description == `` {
		response.Description = http.StatusText(code)
	} else {
		response.Description = ErrorEncode(response.Description)
	}

	if response.URI == `` {
		response.URI = statusError[code].uri
	} else {
		response.URI = ErrorEncode(response.URI)
	}

	ServeJSON(writer, response, code)
}

// ServeErrorDefault serves an error with the given status code and the default
// messages for the code.
func ServeErrorDefault(writer http.ResponseWriter, code int) {
	ServeError(writer, Error{}, code)
}

func ServeMailError(writer http.ResponseWriter) {
	ServeError(writer,
		Error{ID: `mail_failure`, Description: `failed to mail`},
		http.StatusOK)
}

// ServeJSON writes given data in JSON.
func ServeJSON(writer http.ResponseWriter, data interface{}, code int) {
	writer.Header().Set(`Content-Type`, `application/json`)
	writer.WriteHeader(code)

	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)

	if encodeError := encoder.Encode(data); encodeError != nil {
		panic(encodeError)
	}
}
