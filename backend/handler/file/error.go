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

// Error is a structure to hold the context to serve error files.
type Error struct {
	fileSystem       http.FileSystem
	movedPermanently *template.Template
}

// NewError returns a new file.Error.
func NewError(share string) (Error, error) {
	errorPath := path.Join(share, `error`)

	movedPermanently, err :=
		template.ParseFiles(path.Join(errorPath, `301`))
	if err != nil {
		return Error{}, err
	}

	return Error{http.Dir(errorPath), movedPermanently}, nil
}

// ServeError serves the error file corresponding with the given status code.
func (context Error) ServeError(writer http.ResponseWriter,
	code int) (servedContent bool) {
	var err error
	var file http.File
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

			_, err = io.Copy(writer, file)
			if err != nil {
				log.Println(err)
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

				file, err = context.fileSystem.Open(`unknown`)
				if err != nil {
					panic(err)
				}

				fileInfo, err = file.Stat()
				if err != nil {
					panic(err)
				}
			}
		}()

		file, err = context.fileSystem.Open(strconv.Itoa(code))
		if err != nil {
			panic(err)
		}

		fileInfo, err = file.Stat()
		if err != nil {
			panic(err)
		}
	}()

	return false
}

// ServeMovedPermanently serves "Moved Permanently" status.
func (context Error) ServeMovedPermanently(writer http.ResponseWriter,
	location string) {
	var buffer bytes.Buffer

	defer func() {
		header := writer.Header()
		header.Set(`Location`, location)

		if recovered := recover(); recovered == nil {
			header.Set(`Content-Language`, `ja`)
			header.Set(`Content-Length`, strconv.Itoa(buffer.Len()))

			writer.WriteHeader(http.StatusMovedPermanently)

			if _, err := buffer.WriteTo(writer); err != nil {
				log.Println(err)
			}
		} else {
			writer.WriteHeader(http.StatusMovedPermanently)

			log.Println(recovered)
			debug.PrintStack()
		}
	}()

	if err := context.movedPermanently.Execute(&buffer, location); err != nil {
		panic(err)
	}
}
