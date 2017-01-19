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

package file

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"strconv"
)

// FileError is a structure to hold the context to serve error files.
type FileError struct {
	fileSystem       http.FileSystem
	movedPermanently *template.Template
}

// NewError returns a new file.FileError.
func NewError(share string) (FileError, error) {
	errorPath := path.Join(share, `error`)

	movedPermanently, parseError :=
		template.ParseFiles(path.Join(errorPath, `301`))
	if parseError != nil {
		return FileError{}, parseError
	}

	return FileError{http.Dir(errorPath), movedPermanently}, nil
}

// ServeError serves the error file corresponding with the given status code.
func (context FileError) ServeError(writer http.ResponseWriter,
	code int) (servedContent bool) {
	var file http.File
	var fileError error
	var fileInfo os.FileInfo

	defer func() {
		if file != nil {
			defer file.Close()
		}

		if recovered := recover(); recovered == nil {
			header := writer.Header()

			if header.Get(`Content-Encoding`) == `` {
				header.Set(`Content-Length`, strconv.FormatInt(fileInfo.Size(), 10))
			}

			header.Set(`Content-Language`, `ja`)
			header.Set(`Content-Type`, `text/html; charset=UTF-8`)

			writer.WriteHeader(code)

			if _, copyError := io.Copy(writer, file); copyError != nil {
				log.Println(copyError)
			}

			servedContent = true
		} else {
			writer.WriteHeader(code)

			log.Println(recovered)
			debug.PrintStack()
		}
	}()

	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				if file != nil {
					file.Close()
					file = nil
				}

				log.Println(recovered)
				debug.PrintStack()

				file, fileError = context.fileSystem.Open(`unknown`)
				if fileError != nil {
					panic(fileError)
				}

				fileInfo, fileError = file.Stat()
				if fileError != nil {
					panic(fileError)
				}
			}
		}()

		file, fileError = context.fileSystem.Open(strconv.Itoa(code))
		if fileError != nil {
			panic(fileError)
		}

		fileInfo, fileError = file.Stat()
		if fileError != nil {
			panic(fileError)
		}
	}()

	return false
}

// ServeMovedPermanently serves "Moved Permanently" status.
func (context FileError) ServeMovedPermanently(writer http.ResponseWriter,
	location string) {
	var buffer bytes.Buffer

	defer func() {
		header := writer.Header()
		header.Set(`Location`, location)

		if recovered := recover(); recovered == nil {
			header.Set(`Content-Language`, `ja`)
			header.Set(`Content-Length`, strconv.Itoa(buffer.Len()))

			writer.WriteHeader(http.StatusMovedPermanently)

			if _, writeError := buffer.WriteTo(writer); writeError != nil {
				log.Println(writeError)
			}
		} else {
			writer.WriteHeader(http.StatusMovedPermanently)

			log.Println(recovered)
			debug.PrintStack()
		}
	}()

	if executeError := context.movedPermanently.Execute(&buffer, location); executeError != nil {
		panic(executeError)
	}
}
