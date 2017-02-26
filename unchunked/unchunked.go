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

// Package unchunked implements an unchunked response writer wrapper.
package unchunked

import (
	"compress/gzip"
	"net/http"
)

// Unchunked is a structure which implements an unchunked http.ResponseWriter.
type Unchunked struct {
	handle http.HandlerFunc
}

// New returns a new unchunked.Unchunked.
func New(handle http.HandlerFunc) Unchunked {
	return Unchunked{handle}
}

// ServeHTTP serves an unchunked response.
func (unchunked Unchunked) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var gzipEncode bool

	if accepted, present := request.Header[`Accept-Encoding`]; present {
		const unsetGzipQ = 2 << qInitIndex
		gzipQ := uint(unsetGzipQ)
		identityQ := uint(0)

		parse(accepted, func(coding codingIndex, q uint) {
			switch coding {
			case codingAny:
				if gzipQ == unsetGzipQ {
					gzipQ = q
				}

			case codingGzip:
				gzipQ = q

			case codingIdentity:
				identityQ = q
			}
		})

		gzipEncode = gzipQ >= identityQ &&
			gzipQ > 0 && gzipQ != unsetGzipQ
	} else {
		gzipEncode = true
	}

	commonWriter := unchunkedResponseWriter{
		writer, &unchunkedResponseWriterState{code: http.StatusOK},
	}

	if gzipEncode {
		writer.Header().Set(`Content-Encoding`, `gzip`)
		gzipWriter := gzipUnchunkedResponseWriter{
			unchunkedResponseWriter: commonWriter,
		}

		var gzipError error
		gzipWriter.gzip, gzipError = gzip.NewWriterLevel(
			&gzipWriter.unchunkedResponseWriter.state.buffer,
			gzip.BestCompression)
		if gzipError != nil {
			panic(gzipError)
		}

		unchunked.handle(gzipWriter, request)
		gzipWriter.finalize()
	} else {
		identityWriter := identityUnchunkedResponseWriter{commonWriter}

		unchunked.handle(identityWriter, request)
		identityWriter.finalize()
	}
}
