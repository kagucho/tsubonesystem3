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
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
)

type unchunkedResponseWriterState struct {
	buffer bytes.Buffer
	code int
}

type unchunkedResponseWriter struct {
	http.ResponseWriter
	state *unchunkedResponseWriterState
}

type identityUnchunkedResponseWriter struct {
	unchunkedResponseWriter
}

type gzipUnchunkedResponseWriter struct{
	unchunkedResponseWriter
	gzip *gzip.Writer
}

func (writer unchunkedResponseWriter) finalize() {
	writer.ResponseWriter.Header().Set(`Content-Length`,
		strconv.Itoa(writer.state.buffer.Len()))
	writer.ResponseWriter.WriteHeader(writer.state.code)
	writer.state.buffer.WriteTo(writer.ResponseWriter)
}

func (writer identityUnchunkedResponseWriter) Write(body []byte) (int, error) {
	var virtualWriter io.Writer
	if writer.Header().Get(`Content-Length`) == `` {
		virtualWriter = &writer.state.buffer
	} else {
		virtualWriter = writer.ResponseWriter
		writer.state.code = 0
	}

	return virtualWriter.Write(body)
}

func (writer identityUnchunkedResponseWriter) WriteHeader(code int) {
	if writer.Header().Get(`Content-Length`) == `` {
		writer.state.code = code
	} else {
		writer.ResponseWriter.WriteHeader(code)
		writer.state.code = 0
	}
}

func (writer identityUnchunkedResponseWriter) finalize() {
	if writer.unchunkedResponseWriter.state.code != 0 {
		writer.unchunkedResponseWriter.finalize()
	}
}

func (writer gzipUnchunkedResponseWriter) Write(body []byte) (int, error) {
	return writer.gzip.Write(body)
}

func (writer gzipUnchunkedResponseWriter) WriteHeader(code int) {
	writer.state.code = code
}

func (writer gzipUnchunkedResponseWriter) finalize() {
	writer.gzip.Close()
	writer.unchunkedResponseWriter.finalize()
}
