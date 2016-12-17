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

// Package common implements common functions for API v0.
package common

import (
  `encoding/json`
  `log`
  `net/http`
  `reflect`
  `runtime/debug`
  `strconv`
)

// Error is a structure to hold an error to serve.
type Error struct {
  ID string `json:"error"`
  Description string `json:"error_description"`
  URI string `json:"error_uri"`
}

var statusError = map[int]struct{
  id string
  uri string
}{
  http.StatusBadRequest: {
    `invalid_request`,
    `https://tools.ietf.org/html/rfc7231#section-6.5.1`,
  },
  http.StatusNotFound: {
    `not_found`,
    `https://tools.ietf.org/html/rfc7231#section-6.5.4`,
  },
  http.StatusMethodNotAllowed: {
    `method_not_allowed`,
    `https://tools.ietf.org/html/rfc7231#section-6.5.5`,
  },
  http.StatusInternalServerError: {
    `internal_server_error`,
    `https://tools.ietf.org/html/rfc7231#section-6.6.1`,
  },
  http.StatusTooManyRequests: {
    `too_many_requests`,
    `https://tools.ietf.org/html/rfc6585#section-4`,
  },
}

// Recover recovers, logs, and serves Internal Server Error.
func Recover(writer http.ResponseWriter) {
  if recovered := recover(); recovered != nil {
    log.Println(recovered)
    debug.PrintStack()
    ServeErrorDefault(writer, http.StatusInternalServerError)
  }
}

// ServeError serves an error according to the given arguments.
func ServeError(writer http.ResponseWriter,
                id string, description string, uri string, code int) {
  response := func() Error {
    defer Recover(writer)

    var response Error

    if id == `` {
      response.ID = statusError[code].id
    } else {
      response.ID = id
    }

    if description == `` {
      response.Description = http.StatusText(code)
    } else {
      response.Description = description
    }

    if uri == `` {
      response.URI = statusError[code].uri
    } else {
      response.URI = uri
    }

    return response
  }()

  if (response != Error{}) {
    ServeJSON(writer, response, code)
  }
}

/*
  ServeErrorDefault serves an error with the given status code and the default
  messages for the code.
 */
func ServeErrorDefault(writer http.ResponseWriter, code int) {
  ServeError(writer, ``, `` , ``, code)
}

func setHeader(writer http.ResponseWriter, bodyLen int) {
  header := writer.Header()
  if header.Get(`Content-Encoding`) == `` {
    header.Set(`Content-Length`, strconv.Itoa(bodyLen))
  }
  header.Set(`Content-Type`, `application/json; charset=UTF-8`)
}

// ServeJSON writes given data in JSON.
func ServeJSON(writer http.ResponseWriter, data interface{}, code int) {
  bytes := func() []byte {
    defer Recover(writer)

    bytes, marshalError := json.Marshal(data)
    if marshalError != nil {
      panic(marshalError)
    }

    setHeader(writer, len(bytes))

    return bytes
  }()

  writer.WriteHeader(code)
  writer.Write(bytes)
}

// ServeJSONChan writes JSON-chan or abandon if recieved an error with the given
// channel. It also drains the given channels.
func ServeJSONChan(writer http.ResponseWriter, data interface{},
                   dataChan interface{}, errorChan <-chan error, code int) {
  bytesChan := make(chan []byte, 1)

  go func() {
    defer Recover(writer)
    defer close(bytesChan)

    bytes, marshalError := json.Marshal(data)
    if marshalError != nil {
      panic(marshalError)
    }

    bytesChan <- bytes
  }()

loop:
  select {
  case result := <-bytesChan:
    if result != nil {
      setHeader(writer, len(result))
      writer.WriteHeader(code)
      writer.Write(result)
    }

  case result, present := <-errorChan:
    if present {
      ServeErrorDefault(writer, http.StatusInternalServerError)
      log.Println(result)
    } else {
      errorChan = nil
      goto loop
    }
  }

  func() {
    defer func() {
      if recovered := recover(); recovered != nil {
        log.Println(recovered)
        debug.PrintStack()
      }
    }()

    cases := []reflect.SelectCase{
      reflect.SelectCase{
        Dir: reflect.SelectRecv, Chan: reflect.ValueOf(dataChan),
      }, reflect.SelectCase{
        Dir: reflect.SelectRecv, Chan: reflect.ValueOf(errorChan),
      },
    }

    for cases[0].Chan.IsValid() && cases[1].Chan.IsValid() {
      index, value, present := reflect.Select(cases)
      if present {
        if index == 1 {
          log.Println(value)
        }
      } else {
        cases[index].Chan = reflect.Value{}
      }
    }
  }()
}
