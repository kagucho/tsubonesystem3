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

package file

import (
  `bytes`
  `io/ioutil`
  `net/http`
  `net/http/httptest`
  `testing`
)

func assertBodyWithFile(t *testing.T, recorder *httptest.ResponseRecorder,
                        file string) {
  body, bodyError := ioutil.ReadFile(file)
  if bodyError != nil {
    t.Fatal(bodyError)
  }

  if result := recorder.Body.Bytes(); !bytes.Equal(result, body) {
    t.Errorf(`invalid body; expected %q, got %q`, body, result)
  }
}

func TestCustomizedResponseWriter(t *testing.T) {
  t.Parallel()

  fileError, newError := NewError(`test/unknown/na`)
  if newError != nil {
    t.Fatal(newError)
  }

  for _, test := range [...]struct {
        description string
        code int
        body string
      }{
        {`ok`, http.StatusOK, ``},
        {`validError`, http.StatusNotFound, `test/unknown/na/error/404`},
        {`invalidError`, http.StatusBadRequest, ``},
      } {
    test := test

    t.Run(test.description, func(t *testing.T) {
      t.Parallel()

      hijacked := false
      recorder := httptest.NewRecorder()

      customizedWriter :=
        customizedResponseWriter{recorder, &hijacked, fileError}

      customizedWriter.WriteHeader(test.code)

      if expected := test.body != ``; hijacked != expected {
        t.Error(`invalid hijacked flag; expected `, expected,
                `, got `, hijacked)
      }

      if recorder.Code != test.code {
        t.Error(`invalid status code; expected `, test.code,
                `, got `, recorder.Code)
      }

      if test.body == `` {
        if result := recorder.Body.Len(); result != 0 {
          t.Error(`invalid body length; expected 0, got `, result)
        }
      } else {
        assertBodyWithFile(t, recorder, test.body)
      }
    })
  }

  t.Run(`hijacked`, func(t *testing.T) {
    t.Parallel()

    hijacked := true
    recorder := httptest.NewRecorder()
    customizedWriter :=
      customizedResponseWriter{ recorder, &hijacked, fileError }

    if !t.Run(`WriteHeader`, func(t *testing.T) {
      customizedWriter.WriteHeader(http.StatusBadRequest)

      if recorder.Code != http.StatusOK {
        t.Error(`invalid status code; expected `, http.StatusOK,
                `, got `, recorder.Code)
      }
    }) {
      t.FailNow()
    }

    t.Run(`Write`, func(t *testing.T) {
      customizedWriter.Write([]byte(`body`))
      if result := recorder.Body.Len(); result != 0 {
        t.Error(`invalid body length; expected 0, got `, result)
      }
    })
  })
}

func (file File) TestServeHTTP(t *testing.T) {
  for _, test := range [...]struct {
        file string
        lang string
        request string
      }{
        {`test/unknown/file/public/index`, `ja`, `/`},
        {`test/unknown/file/public/index.js`, ``, `/index.js`},
        {`test/unknown/file/public/license`, `ja`, `/license`},
        {`test/unknown/file/error/404`, `ja`, `/index`},
        {`test/unknown/file/error/404`, `ja`, `/invalid`},
      } {
    test := test

    t.Run(test.request, func(t *testing.T) {
      t.Parallel()

      recorder := httptest.NewRecorder()
      request := httptest.NewRequest(
        `GET`, `http://kagucho.net` + test.request, nil)

      file.ServeHTTP(recorder, request)

      assertBodyWithFile(t, recorder, test.file)

      if result := recorder.HeaderMap.Get(`Content-Language`)
         result != test.lang {
        t.Errorf(`invalid Content-Language field in header; expected %q, got %q`,
                 test.lang, result)
      }
    })
  }

  t.Run(`/license/?redundant_query`, func(t *testing.T) {
    t.Parallel()

    recorder := httptest.NewRecorder()
    request := httptest.NewRequest(
      `GET`, `http://kagucho.net/license/?redundant_query`, nil)

    file.ServeHTTP(recorder, request)

    if result := recorder.HeaderMap.Get(`Location`);
       result != `http://kagucho.net/license` {
      t.Errorf(`invalid Location field in header; expected "http://kagucho.net/license", got %q`,
               result)
    }
  })

  t.Run(`/license?redundant_query`, func(t *testing.T) {
    t.Parallel()

    recorder := httptest.NewRecorder()
    request := httptest.NewRequest(
      `GET`, `http://kagucho.net/license?redundant_query`, nil)

    var writer http.ResponseWriter = recorder
    file.ServeHTTP(writer, request)

    if result := recorder.HeaderMap.Get(`Location`)
       result != `http://kagucho.net/license` {
      t.Errorf(`invalid Location field in header; expected "http://kagucho.net/license", got %q`,
               result)
    }
  })
}

func TestFile(t *testing.T) {
  t.Parallel()

  var file File

  if !t.Run(`New`, func(t *testing.T) {
    fileError, newError := NewError(`test/unknown/file`)
    if newError != nil {
      t.Fatal(newError)
    }

    file = New(`test/unknown/file`, fileError)
  }) {
    t.FailNow()
  }

  t.Run(`ServeHTTP`, file.TestServeHTTP)
}
