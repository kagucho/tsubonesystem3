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

package apiv0

import (
	"bytes"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/util"
	"net/http"
)

type handlerFunc func(http.ResponseWriter, *http.Request, shared)

type handler interface {
	serveHTTP(writer http.ResponseWriter, request *http.Request, shared shared)
}

type field struct {
	name string
	value string
}

type methodMux struct {
	handlers map[string]handlerFunc
	options []field
}

func (mux methodMux) serveHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	handle := mux.handlers[request.Method]
	if handle != nil {
		handle(writer, request, shared)
		return
	}

	buffer := bytes.NewBufferString(`OPTIONS`)
	for method := range mux.handlers {
		buffer.WriteString(`, `)
		buffer.WriteString(method)
	}

	writer.Header().Set(`Allow`, buffer.String())

	if request.Method == `OPTIONS` {
		for _, option := range mux.options {
			writer.Header().Set(option.name, option.value)
		}

		writer.WriteHeader(http.StatusOK)
	} else {
		util.ServeErrorDefault(writer, http.StatusMethodNotAllowed)
	}
}
