/*
  Copyright (C) 2016  Kagucho <kagucho.net@gmail.com>

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU Affero General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU Affero General Public License for more details.

  You should have received a copy of the GNU Affero General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

// Package file implements a file server.
package file

import (
  `net/http`
  `path`
)

// File is a structure to keep the context of the file server.
type File struct {
  publicRoot string
  fileError FileError
}

// New returns a new file server.
func New(share string, fileError FileError) File {
  return File{path.Join(share, `public`), fileError}
}

type customizedResponseWriter struct {
  http.ResponseWriter
  hijacked *bool
  fileError FileError
}

func (writer customizedResponseWriter) Write(bytes []byte) (int, error) {
  if *writer.hijacked {
    return 0, nil
  }

  return writer.ResponseWriter.Write(bytes)
}

func (writer customizedResponseWriter) WriteHeader(code int) {
  if code < 400 {
    writer.ResponseWriter.WriteHeader(code)
  } else if !*writer.hijacked {
    *writer.hijacked =
      writer.fileError.ServeError(writer.ResponseWriter, code)
  }
}

// ServeHTTP serves files via HTTP.
func (file File) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
  var relative string
  if request.URL.Path == `/` {
    relative = `index`
    writer.Header().Set(`Content-Language`, `ja`)
  } else if request.URL.Path[len(request.URL.Path) - 1] == '/' {
    request.URL.Path = request.URL.Path[:len(request.URL.Path) - 1]
    request.URL.RawQuery = ``
    file.fileError.ServeMovedPermanently(writer, request.URL.String())
    return
  } else if request.URL.RawQuery != `` {
    request.URL.RawQuery = ``
    file.fileError.ServeMovedPermanently(writer, request.URL.String())
    return
  } else if request.URL.Path == `/index` {
    file.fileError.ServeError(writer, http.StatusNotFound)
    return
  } else {
    relative = request.URL.Path
    if path.Ext(relative) == `` {
      writer.Header().Set(`Content-Language`, `ja`)
    }
  }

  hijacked := false
  writerCustomized := customizedResponseWriter{
    writer, &hijacked, file.fileError,
  }

  http.ServeFile(writerCustomized, request,
                 path.Join(file.publicRoot, relative))
}
