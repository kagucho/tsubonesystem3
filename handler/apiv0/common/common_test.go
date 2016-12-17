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

package common

import (
  `net/http`
  `net/http/httptest`
  `testing`
  `time`
)

func TestRecover(t *testing.T) {
  t.Parallel()

  t.Run(`normal`, func(t *testing.T) {
    t.Parallel()

    recorder := httptest.NewRecorder()

    Recover(recorder)

    if recorder.Code != http.StatusOK {
      t.Error(`invalid status code; expected `, http.StatusOK, `, got `,
              recorder.Code)
    }
  })

  t.Run(`panicked`, func(t *testing.T) {
    t.Parallel()

    recorder := httptest.NewRecorder()

    func() {
      defer Recover(recorder)
      panic(0)
    }()

    if recorder.Code != http.StatusInternalServerError {
      t.Error(`invalid status code; expected `, http.StatusInternalServerError,
              `, got `, recorder.Code)
    }
  })
}

func TestServeError(t *testing.T) {
  t.Parallel()

  for _, test := range []struct{
        testDescription string
        id string
        description string
        uri string
        code int
        expectedBody string
      }{
        {
          `specified`, `id`, `description`, `uri`,
          http.StatusBadRequest,
          `{"error":"id","error_description":"description","error_uri":"uri"}`,
        }, {
          `unspecifiedBadRequest`, ``, ``, ``,
          http.StatusBadRequest,
          `{"error":"invalid_request","error_description":"Bad Request","error_uri":"https://tools.ietf.org/html/rfc7231#section-6.5.1"}`,
        }, {
          `unspecifiedNotFound`, ``, ``, ``,
          http.StatusNotFound,
          `{"error":"not_found","error_description":"Not Found","error_uri":"https://tools.ietf.org/html/rfc7231#section-6.5.4"}`,
        }, {
          `unspecifiedMethodNotAllowed`, ``, ``, ``,
          http.StatusMethodNotAllowed,
          `{"error":"method_not_allowed","error_description":"Method Not Allowed","error_uri":"https://tools.ietf.org/html/rfc7231#section-6.5.5"}`,
        }, {
          `unspecifiedInternalServerError`, ``, ``, ``,
          http.StatusInternalServerError,
          `{"error":"internal_server_error","error_description":"Internal Server Error","error_uri":"https://tools.ietf.org/html/rfc7231#section-6.6.1"}`,
        },
      } {
    test := test

    t.Run(test.testDescription, func(t *testing.T) {
      t.Parallel()

      recorder := httptest.NewRecorder()
      ServeError(recorder, test.id, test.description, test.uri, test.code)

      if result := recorder.Body.String(); result != test.expectedBody {
        t.Errorf("expected %q, got %q", result, recorder.Body)
      }

      if recorder.Code != test.code {
        t.Error(`invalid status code; expected `, test.code,
                `, got `, recorder.Code)
      }
    })
  }
}

func TestServeErrorDefault(t *testing.T) {
  t.Parallel()

  recorder := httptest.NewRecorder()
  ServeErrorDefault(recorder, http.StatusBadRequest)

  const expected = `{"error":"invalid_request","error_description":"Bad Request","error_uri":"https://tools.ietf.org/html/rfc7231#section-6.5.1"}`
  if result := recorder.Body.String(); result != expected {
    t.Errorf(`invalid body; expected %q, got %q`, expected, result)
  }

  if recorder.Code != http.StatusBadRequest {
    t.Error(`invalid status code; expected `, http.StatusBadRequest,
            `, got `, recorder.Code)
  }
}

func TestServeJSON(t *testing.T) {
  const contentType = `application/json; charset=UTF-8`

  t.Parallel()

  t.Run(`invalidData`, func(t *testing.T) {
    t.Parallel()

    recorder := httptest.NewRecorder()
    ServeJSON(recorder, TestServeJSON, http.StatusOK)

    if recorder.Code != http.StatusInternalServerError {
      t.Error(`invalid status code; expected `, http.StatusInternalServerError,
              `, got `, recorder.Code)
    }
  })

  t.Run(`WithContentEncoding`, func(t *testing.T) {
    t.Parallel()

    recorder := httptest.NewRecorder()
    recorder.HeaderMap.Set(`Content-Encoding`, `br`)

    ServeJSON(recorder, "value", http.StatusOK)

    if _, present := recorder.HeaderMap[`Content-Length`]; present {
      t.Error(`invalid Content-Length in header; expected not set, got set`)
    }

    if result := recorder.HeaderMap.Get(`Content-Type`); result != contentType {
      t.Errorf(`invalid Content-Type in header; expected %q, got %q`,
               contentType, result)
    }

    if recorder.Code != http.StatusOK {
      t.Errorf(`invalid status code; expected %q, got %q`,
               http.StatusOK, recorder.Code)
    }

    if result := recorder.Body.String(); result != `"value"` {
      t.Errorf(`invalid body; expected "\"value\"", got %q`, result)
    }
  })

  t.Run(`WithoutContentEncoding`, func(t *testing.T) {
    t.Parallel()

    recorder := httptest.NewRecorder()
    ServeJSON(recorder, "value", http.StatusOK)

    if result := recorder.HeaderMap.Get(`Content-Length`); result != `7` {
      t.Errorf(`invalid Content-Length in header; expected "7", got %q`,
               result)
    }

    if result := recorder.HeaderMap.Get(`Content-Type`); result != contentType {
      t.Errorf(`invalid Content-Type in header; expected %q, got %q`,
               contentType, result)
    }

    if recorder.Code != http.StatusOK {
      t.Error(`invalid status code; expected `, http.StatusOK, `, got `,
              recorder.Code)
    }

    if result := recorder.Body.String(); result != `"value"` {
      t.Errorf(`invalid body; expected "\"value\"", got %q`, result)
    }
  })
}

type ResponseRecorderChan struct {
  *httptest.ResponseRecorder
  buffer chan []byte
}

func (recorder ResponseRecorderChan) Write(buffer []byte) (int, error) {
  recorder.buffer <- buffer
  return recorder.ResponseRecorder.Write(buffer)
}

func TestServeJSONChan(t *testing.T) {
  t.Parallel()

  t.Run(`select`, func(t *testing.T) {
    t.Parallel()

    t.Run(`whileMarshalling`, func(t *testing.T) {
      t.Parallel()
      t.Run(`0`, func(t *testing.T) {
        t.Parallel()

        recorder := httptest.NewRecorder()

        dataChan := make(chan int)

        errorChan := make(chan error)
        close(errorChan)

        time.AfterFunc(134217728, func() {
          close(dataChan)
        })

        ServeJSONChan(recorder, dataChan, dataChan, errorChan, http.StatusOK)

        if recorder.Code != http.StatusOK {
          t.Error(`invalid status code; expected `, http.StatusOK,
                  `, got `, recorder.Code)
        }

        if result := recorder.Body.String(); result != `[]` {
          t.Errorf(`invalid body; expected "[]", got %q`, result)
        }
      })

      t.Run(`1`, func(t *testing.T) {
        t.Parallel()

        recorder := httptest.NewRecorder()

        dataChan := make(chan int)
        defer close(dataChan)

        errorChan := make(chan error, 1)
        errorChan <- nil
        close(errorChan)

        ServeJSONChan(recorder, dataChan, nil, errorChan, http.StatusOK)

        if recorder.Code != http.StatusInternalServerError {
          t.Error(`invalid status code; expected `,
                  http.StatusInternalServerError, `, got `, recorder.Code)
        }
      })
    })

    t.Run(`afterMarshalling`, func(t *testing.T) {
      t.Parallel()

      recorder :=
        ResponseRecorderChan{ httptest.NewRecorder(), make(chan []byte) }

      dataChan := make(chan int)
      errorChan := make(chan error)

      go func() {
        defer close(dataChan)
        defer close(errorChan)

        <-recorder.buffer

        dataChan <- 0
        errorChan <- nil
        dataChan <- 0
      }()

      ServeJSONChan(recorder, nil, dataChan, errorChan, http.StatusOK)

      if recorder.Code != http.StatusOK {
        t.Error(`invalid status code; expected `, http.StatusOK,
                `, got `, recorder.Code)
      }

      if result := recorder.Body.String(); result != `null` {
        t.Errorf(`invalid body; expected "null", got %q`, result)
      }
    })
  })

  t.Run(`invalid`, func(t *testing.T) {
    t.Parallel()

    t.Run(`data`, func(t *testing.T) {
      t.Parallel()

      recorder := httptest.NewRecorder()

      errorChan := make(chan error)
      defer close(errorChan)

      ServeJSONChan(recorder, ServeJSONChan, nil, errorChan, http.StatusOK)
      if recorder.Code != http.StatusInternalServerError {
        t.Error(`invalid status code; expected `,
                http.StatusInternalServerError, `, got `, recorder.Code)
      }
    })

    t.Run(`chan`, func(t *testing.T) {
      t.Parallel()

      recorder := httptest.NewRecorder()

      errorChan := make(chan error)
      defer close(errorChan)

      ServeJSONChan(recorder, nil, ServeJSONChan, nil, http.StatusOK)

      if recorder.Code != http.StatusOK {
        t.Error(`invalid status code; expected `, http.StatusOK,
                `, got `, recorder.Code)
      }

      if result := recorder.Body.String(); result != `null` {
        t.Errorf(`invalid body; expected "null", got %q`, result)
      }
    })
  })
}

func TestSetHeader(t *testing.T) {
  t.Parallel()

  testType := func(t *testing.T, recorder *httptest.ResponseRecorder) {
    const expected = `application/json; charset=UTF-8`
    if result := recorder.HeaderMap.Get(`Content-Type`); result != expected {
      t.Errorf(`invalid Content-Type in header; expected %q, got %q`,
               expected, result)
    }
  }

  t.Run(`WithContentEncoding`, func(t *testing.T) {
    recorder := httptest.NewRecorder()
    recorder.HeaderMap.Set(`Content-Encoding`, `br`)

    setHeader(recorder, 0)

    if _, present := recorder.HeaderMap[`Content-Length`]; present {
      t.Error(`invalid Content-Length; expected not set, got set`)
    }

    testType(t, recorder)
  })

  t.Run(`WithoutContentEncoding`, func(t *testing.T) {
    recorder := httptest.NewRecorder()

    setHeader(recorder, 0)

    const expected = `0`
    if result := recorder.HeaderMap.Get(`Content-Length`);
       result != expected {
      t.Errorf(`invalid Content-Length; expected %q, got %q`, expected, result)
    }

    testType(t, recorder)
  })
}
