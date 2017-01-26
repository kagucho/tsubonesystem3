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

package unchunked

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type helloHandler struct {
	encoding *string
}

func (handler helloHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	*handler.encoding = writer.Header().Get(`Content-Encoding`)
	writer.Write([]byte{'\n'})
}

func TestUnchunked(t *testing.T) {
	t.Parallel()

	encoding := ``
	handler := helloHandler{encoding: &encoding}
	var unchunked Unchunked

	if !t.Run(`New`, func(t *testing.T) {
		unchunked = New(handler)
	}) {
		t.FailNow()
	}

	t.Run(`ServeHTTP`, func(t *testing.T) {
		for _, test := range [...]struct {
			encoding string
			gzip     bool
		}{
			/*
				RFC 7231 - Hypertext Transfer Protocol (HTTP/1.1): Semantics and Content
				5.3.4.  Accept-Encoding
				https://tools.ietf.org/html/rfc7231#section-5.3.4
				1.  If no Accept-Encoding field is in the
				request, any content-coding is considered
				acceptable by the user agent.
			*/
			{``, true},

			/*
				> If an Accept-Encoding header field is present
				> in a request and none of the available
				> representations for the response have a
				> content-coding that is listed as acceptable,
				> the origin server SHOULD send a response
				> without any content-coding.
			*/
			{`deflate`, false},

			/*
				> The asterisk "*" symbol in an Accept-Encoding
				> field matches any available content-coding
				> not explicitly listed in the header field.

				5.3.1.  Quality Values
				https://tools.ietf.org/html/rfc7231#section-5.3.1
				> The weight is normalized to a real number in
				> the range 0 through 1, where 0.001 is the
				> least preferred and 1 is the most preferred; a
				> value of 0 means "not acceptable".
			*/
			{`*;q=0`, false},

			{`*;q=0.5`, true},

			/*
				> If no "q" parameter is present, the default
				> weight is 1.

				5.3.4.  Accept-Encoding
				https://tools.ietf.org/html/rfc7231#section-5.3.4
				> An "identity" token is used as a synonym
				> for "no encoding" in order to communicate when
				> no encoding is preferred.
			*/
			{`*;q=0.5,identity`, false},

			/*
				4.2.3.  Gzip Coding
				https://tools.ietf.org/html/rfc7230#section-4.2.3
				> The "gzip" coding is an LZ77 coding with a
				> 32-bit Cyclic Redundancy Check (CRC) that is
				> commonly produced by the gzip file compression
				> program [RFC1952].
			*/
			{`gzip;q=0`, false}, {`gzip;q=0.5`, true},
			{`gzip;q=0.5,identity`, false},
			{`gzip;q=0,*;q=1`, false},
		} {
			test := test

			t.Run(test.encoding, func(t *testing.T) {
				recorder := httptest.NewRecorder()

				request := httptest.NewRequest(`GET`, `https://kagucho.net/`, nil)

				if test.encoding != `` {
					request.Header.Set(`Accept-Encoding`, test.encoding)
				}

				unchunked.ServeHTTP(recorder, request)

				if test.gzip {
					const expectedEncoding = `gzip`
					if *handler.encoding != expectedEncoding {
						t.Errorf(`invalid Content-Encoding in field when Write was called; expected %q, got %q`,
							expectedEncoding, handler.encoding)
					}

					const expectedBody = "\x1F\x8B\x08\x00\x00\x09\x6E\x88\x02\xFF\xE2\x02\x04\x00\x00\xFF\xFF\x93\x06\xD7\x32\x01\x00\x00\x00"
					if result := recorder.Body.String(); result != expectedBody {
						t.Errorf(`invalid body; expected %X, got %X`,
							expectedBody, result)
					}
				} else {
					const expectedBody = "\n"
					if result := recorder.Body.String(); result != expectedBody {
						t.Errorf(`invalid body; expected %X, got %X`,
							expectedBody, result)
					}
				}
			})
		}
	})
}
