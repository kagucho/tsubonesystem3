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
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUnchunkedResponseWriter(t *testing.T) {
	t.Parallel()

	testRecoreded := func(t *testing.T, recorder *httptest.ResponseRecorder, code int) {
		if recorder.Code != code {
			t.Errorf(`invalid status code; expected %v, got %v`,
				code, recorder.Code)
		}

		if result := recorder.HeaderMap.Get(`Content-Length`); result != `1` {
			t.Errorf(`invalid Content-Length field in header; expected "1", got %q`,
				result)
		}

		if result := recorder.Body.String(); result != "\n" {
			t.Errorf(`invalid body; expected "\n", got %q`, result)
		}
	}

	t.Run(`WithContentLength`, func(t *testing.T) {
		testCommon := func(t *testing.T, writeHeader bool) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			recorder.HeaderMap.Set(`Content-Length`, `1`)

			unchunked := identityUnchunkedResponseWriter{
				unchunkedResponseWriter{
					recorder,
					&unchunkedResponseWriterState{
						code: http.StatusOK,
					},
				},
			}

			if writeHeader {
				unchunked.WriteHeader(http.StatusBadRequest)
			}

			unchunked.Write([]byte{'\n'})

			if writeHeader {
				testRecoreded(t, recorder, http.StatusBadRequest)
			} else {
				testRecoreded(t, recorder, http.StatusOK)
			}

			unchunked.finalize()
		}

		t.Run(`WithWriteHeader`, func(t *testing.T) {
			testCommon(t, true)
		})

		t.Run(`WithoutWriteHeader`, func(t *testing.T) {
			testCommon(t, false)
		})
	})

	t.Run(`WithoutContentLength`, func(t *testing.T) {
		testCommon := func(t *testing.T, writeHeader bool) {
			t.Parallel()

			recorder := httptest.NewRecorder()

			unchunked := identityUnchunkedResponseWriter{
				unchunkedResponseWriter{
					recorder,
					&unchunkedResponseWriterState{
						code: http.StatusOK,
					},
				},
			}

			if writeHeader {
				unchunked.WriteHeader(http.StatusBadRequest)
			}

			unchunked.Write([]byte{'\n'})
			unchunked.finalize()

			if writeHeader {
				testRecoreded(t, recorder, http.StatusBadRequest)
			} else {
				testRecoreded(t, recorder, http.StatusOK)
			}
		}

		t.Run(`WithWriteHeader`, func(t *testing.T) {
			testCommon(t, true)
		})

		t.Run(`WithoutWriteHeader`, func(t *testing.T) {
			testCommon(t, false)
		})
	})
}

func TestGzipUnchunkedResponseWriter(t *testing.T) {
	recorder := httptest.NewRecorder()

	unchunked := gzipUnchunkedResponseWriter{
		unchunkedResponseWriter: unchunkedResponseWriter{
			recorder,
			&unchunkedResponseWriterState{code: http.StatusOK},
		},
	}

	unchunked.gzip, _ = gzip.NewWriterLevel(
		&unchunked.unchunkedResponseWriter.state.buffer,
		gzip.BestCompression)

	unchunked.WriteHeader(http.StatusBadRequest)
	unchunked.Write([]byte{'\n'})
	unchunked.finalize()

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`invalid status code; expected %v, got %v`,
			http.StatusBadRequest, recorder.Code)
	}

	const expectedLength = `25`
	if result := recorder.HeaderMap.Get(`Content-Length`); result != expectedLength {
		t.Errorf(`invalid Content-Length field in header; expected %q, got %q`,
			expectedLength, result)
	}

	const expectedBody = "\x1F\x8B\x08\x00\x00\x09\x6E\x88\x02\xFF\xE2\x02\x04\x00\x00\xFF\xFF\x93\x06\xD7\x32\x01\x00\x00\x00"
	if result := recorder.Body.String(); result != expectedBody {
		t.Errorf(`invalid body; expected %X, got %X`,
			expectedBody, result)
	}
}
