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

package common

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorEncode(t *testing.T) {
	t.Parallel()

	for _, test := range [...]struct {
		decoded string
		encoded string
	}{
		{"\x00", `%00`}, {"\x20", "\x20"},
		{"\x22", `%22`}, {`%`, `%25`},
		{"\x5C", `%5C`}, {"\x80", `%80`},
	} {
		test := test

		t.Run(test.decoded, func(t *testing.T) {
			t.Parallel()

			if result := ErrorEncode(test.decoded); result != test.encoded {
				t.Errorf(`expected %q, got %q`,
					test.encoded, result)
			}
		})
	}
}

func TestServeError(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		testDescription string
		response        Error
		code            int
		expectedBody    string
	}{
		{
			`escape`,
			Error{
				/*
					5.2.  Error Response
					https://tools.ietf.org/html/rfc6749#section-5.2

					> Values for the "error" parameter MUST
					> NOT include characters outside the set
					> %x20-21 / %x23-5B / %x5D-7E.
				*/
				"\x00",

				/*
					> Values for the "error_description"
					> parameter MUST NOT include characters
					> outside the set %x20-21 / %x23-5B /
					> %x5D-7E.
				*/
				"\x20",

				/*
					> Values for the "error_uri" parameter
					> MUST conform to the URI-reference
					> syntax and thus MUST NOT include
					> characters outside the set %x21 /
					> %x23-5B / %x5D-7E.
				*/
				"\x22",
			},

			http.StatusBadRequest,
			`{"error":"%00","error_description":" ","error_uri":"%22"}
`,
		}, {
			`unspecifiedBadRequest`, Error{},
			http.StatusBadRequest,

			/*
				> invalid_request
				>       The request is missing a required
				>       parameter, includes an unsupported
				>       parameter value (other than grant type),
				>       repeats a parameter, includes multiple
				>       credentials, utilizes more than one
				>       mechanism for authenticating the client,
				>       or is otherwise malformed.
			*/
			`{"error":"invalid_request","error_description":"Bad Request","error_uri":"https://tools.ietf.org/html/rfc7231#section-6.5.1"}
`,
		}, {
			`unspecifiedNotFound`, Error{},
			http.StatusNotFound,
			`{"error":"not_found","error_description":"Not Found","error_uri":"https://tools.ietf.org/html/rfc7231#section-6.5.4"}
`,
		}, {
			`unspecifiedMethodNotAllowed`, Error{},
			http.StatusMethodNotAllowed,

			/*
				4.1.2.1.  Error Response
				https://tools.ietf.org/html/rfc6749#section-4.1.2.1

				> unsupported_response_type
				>       The authorization server does not
				>       support obtaining an authorization code
				>       using this method.
			*/
			`{"error":"unsupported_response_type","error_description":"Method Not Allowed","error_uri":"https://tools.ietf.org/html/rfc7231#section-6.5.5"}
`,
		}, {
			/*
				> server_error
				>       The authorization server encountered an
				>       unexpected condition that prevented it
				>       from fulfilling the request.
			*/
			`unspecifiedInternalServerError`, Error{},
			http.StatusInternalServerError,
			`{"error":"server_error","error_description":"Internal Server Error","error_uri":"https://tools.ietf.org/html/rfc7231#section-6.6.1"}
`,
		},
	} {
		test := test

		t.Run(test.testDescription, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			ServeError(recorder, test.response, test.code)

			if result := recorder.Body.String(); result != test.expectedBody {
				t.Errorf("expected %q, got %q",
					result, recorder.Body)
			}

			if recorder.Code != test.code {
				t.Error(`invalid status code; expected `, test.code,
					`, got `, recorder.Code)
			}
		})
	}
}

func TestServeJSON(t *testing.T) {
	/*
		RFC 4627 - The application/json Media Type for JavaScript Object Notation (JSON)
		The MIME media type for JSON text is application/json.

		6. IANA Considerations
		https://tools.ietf.org/html/rfc4627#section-6
		> The MIME media type for JSON text is application/json.
	*/
	const contentType = `application/json`

	t.Parallel()

	t.Run(`invalidData`, func(t *testing.T) {
		t.Parallel()

		recorder := httptest.NewRecorder()

		defer func() {
			if recover() == nil {
				t.Error(`expected panicking`)
			}
		}()

		ServeJSON(recorder, TestServeJSON, http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	ServeJSON(recorder, "value", http.StatusOK)

	if result := recorder.HeaderMap.Get(`Content-Type`); result != contentType {
		t.Errorf(`invalid Content-Type in header; expected %q, got %q`,
			contentType, result)
	}

	if recorder.Code != http.StatusOK {
		t.Errorf(`invalid status code; expected %q, got %q`,
			http.StatusOK, recorder.Code)
	}

	if result := recorder.Body.String(); result != "\"value\"\n" {
		t.Errorf(`invalid body; expected "\"value\"", got %q`, result)
	}
}

type ResponseRecorderChan struct {
	*httptest.ResponseRecorder
	buffer chan []byte
}

func (recorder ResponseRecorderChan) Write(buffer []byte) (int, error) {
	recorder.buffer <- buffer
	return recorder.ResponseRecorder.Write(buffer)
}
