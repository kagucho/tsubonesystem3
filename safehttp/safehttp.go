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

package safehttp

import "net/http"

/*
ResponseWriter implements recoverable http.ResponseWriter. It must be
initialized with safehttp.New.
*/
type ResponseWriter struct {
	header  http.Header
	wrapped http.ResponseWriter
	written bool
}

// NewWriter returns a new safehttp.ResponseWriter.
func NewWriter(wrapped http.ResponseWriter) ResponseWriter {
	return ResponseWriter{header: make(http.Header), wrapped: wrapped}
}

// Header implements Header of http.ResponseWriter.
func (writer *ResponseWriter) Header() http.Header {
	return writer.header
}

/*
Recover recovers the writer. If the header is not written, it calls the given
function to allow to recover. Otherwise, it does nothing.
*/
func (writer *ResponseWriter) Recover(callback func()) {
	if !writer.written {
		callback()
	}
}

// Write implements Write of http.ResponseWriter.
func (writer *ResponseWriter) Write(bytes []byte) (int, error) {
	writer.prepare()
	return writer.wrapped.Write(bytes)
}

// WriteHeader implements WriteHeader of http.ResponseWriter.
func (writer *ResponseWriter) WriteHeader(code int) {
	writer.prepare()
	writer.wrapped.WriteHeader(code)
}

func (writer *ResponseWriter) prepare() {
	if !writer.written {
		writer.written = true

		header := writer.wrapped.Header()
		for key, value := range writer.header {
			header[key] = value
		}
	}
}
