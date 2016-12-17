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

package server

import (
  `encoding/json`
  `fmt`
  `github.com/kagucho/tsubonesystem3/db`
  `github.com/kagucho/tsubonesystem3/handler/apiv0/token/provider`
  `net/http`
  `net/http/httptest`
  `regexp`
  `strings`
  `testing`)

func (server Server) TestServeHTTP(t *testing.T) {
  db, dbError := db.Open()
  if dbError != nil {
    t.Fatal(dbError)
  }

  record := func(t *testing.T, method string, body string) *httptest.ResponseRecorder {
    recorder := httptest.NewRecorder()

    request := httptest.NewRequest(
      method, `http://kagucho.net/api/v0/token`, strings.NewReader(body))
    request.Header.Set(`Content-Type`, `application/x-www-form-urlencoded`)

    server.ServeHTTP(recorder, request, db)

    return recorder
  }

  func() {
    defer db.Close()

    for _, test := range [...]struct{
          description string
          requestMethod string
          requestBody string
          code int
          regexp bool
          pattern string
        }{
          {
            `invalidGrantType`, `POST`,
            `grant_type=invalid&username=1stDisplayId&password=1stPassword`,
            http.StatusBadRequest, false,
            `{"error":"invalid_grant","error_description":"expected grant_type \"password\" or \"refresh_token\", got \"invalid\"","error_uri":"https://tools.ietf.org/html/rfc6749#section-5.2"}`,
          }, {
            `invalidUsername`, `POST`,
            `grant_type=password&password=1stPassword`,
            http.StatusBadRequest, false,
            `{"error":"invalid_grant","error_description":"invalid username and/or password","error_uri":"https://tools.ietf.org/html/rfc6749#section-5.2"}`,
          }, {
            `invalidPassword`, `POST`,
            `grant_type=password&username=1stDisplayId`,
            http.StatusBadRequest, false,
            `{"error":"invalid_grant","error_description":"invalid username and/or password","error_uri":"https://tools.ietf.org/html/rfc6749#section-5.2"}`,
          }, {
            `invalidMethod`, `GET`, ``,
            http.StatusMethodNotAllowed, false,
             `{"error":"method_not_allowed","error_description":"expected \"POST\" method request, got \"GET\" method request","error_uri":"https://tools.ietf.org/html/rfc7231#section-6.5.5"}`,
          }, {
             `invalidRefresh`, `POST`, `grant_type=refresh_token`,
             http.StatusBadRequest, false,
            `{"error":"invalid_grant","error_description":"expected 3 parts, got 1 parts","error_uri":"https://tools.ietf.org/html/rfc7519#section-3.1"}`,
          }, {
             `valid`, `POST`,
             `grant_type=password&username=1stDisplayId&password=1stPassword`,
             http.StatusOK, true,
             `{"access_token":"[^"]+","refresh_token":"[^"]+"}`,
          },
        } {
      test := test

      t.Run(test.description, func(t *testing.T) {
        recorder := record(t, test.requestMethod, test.requestBody)

        if test.regexp {
          result := recorder.Body.Bytes();
          if matched, matchError := regexp.Match(test.pattern, result);
             matchError != nil {
            t.Error(matchError)
          } else if !matched {
            t.Error(`invalid body; expected to match with `, test.pattern,
                    `, got `, result)
          }
        } else {
          if result := recorder.Body.String(); result != test.pattern {
            t.Error(`invalid body; expected `, test.pattern, `, got `, result)
          }
        }

        if recorder.Code != test.code {
          t.Error(`invalid status code; expected `, test.code,
                  `, got `, recorder.Code)
        }
      })
    }

    t.Run(`limit`, func(t *testing.T) {
      for count := 0; count < 9; count++ {
        recorder := record(t, `POST`, `grant_type=password&username=limit`)

        const expectedBody = `{"error":"invalid_grant","error_description":"invalid username and/or password","error_uri":"https://tools.ietf.org/html/rfc6749#section-5.2"}`
        if result := recorder.Body.String(); result != expectedBody {
          t.Error(`invalid body; expected `, expectedBody, `, got `, result)
        }

        if recorder.Code != http.StatusBadRequest {
          t.Error(`invalid status code; expected `, http.StatusBadRequest,
                  `, got `, recorder.Code)
        }
      }

      recorder := record(t, `POST`, `grant_type=password&username=limit`)

      const expectedBody = `{"error":"too_many_requests","error_description":"Too Many Requests","error_uri":"https://tools.ietf.org/html/rfc6585#section-4"}`
      if result := recorder.Body.String(); result != expectedBody {
        t.Errorf(`invalid body; expected `, expectedBody, `, got `, result)
      }

      if recorder.Code != http.StatusTooManyRequests {
        t.Error(`invalid status code; expected `, http.StatusTooManyRequests,
                `, got `, recorder.Code)
      }
    })
  }()

  t.Run(`refresh`, func(t *testing.T) {
    t.Parallel()

    refreshToken, tokenError := server.refresh.Issue(
                                   `1stDisplayId`, `basic`, 70368744177664)
    if tokenError != nil {
      t.Fatal(tokenError)
    }

    recorder := record(t, `POST`,
                       fmt.Sprint(`grant_type=refresh_token&refresh_token=`,
                                  refreshToken))

    var response response
    decoder := json.NewDecoder(recorder.Body)
    if decodeError := decoder.Decode(&response); decodeError != nil {
      t.Error(decodeError)
    } else if response.RefreshToken != `` {
      t.Errorf(`invalid refresh_token; expected empty "", got %q`,
               response.RefreshToken)
    }
  })

  t.Run(`expiringRefresh`, func(t *testing.T) {
    t.Parallel()

    refreshToken, tokenError := server.refresh.Issue(
                                    `1stDisplayId`, `basic`, 4398046511103)
    if tokenError != nil {
      t.Fatal(tokenError)
    }

    recorder := record(t, `POST`,
                       fmt.Sprint(`grant_type=refresh_token&refresh_token=`,
                                  refreshToken))

    var response response
    decoder := json.NewDecoder(recorder.Body)
    if decodeError := decoder.Decode(&response); decodeError != nil {
      t.Error(decodeError)
    } else if response.RefreshToken == `` {
      t.Error(`expected refresh_token, got empty ""`)
    }
  })

  t.Run(`internalServerError`, func(t *testing.T) {
    t.Parallel()

    recorder := record(t, `POST`, `grant_type=password&username=1stDisplayId`)

    if recorder.Code != http.StatusInternalServerError {
      t.Error(`invalid status code; expected `,
              http.StatusInternalServerError, `, got `, recorder.Code)
    }
  })
}

func TestServer(t *testing.T) {
  t.Parallel()

  token, tokenError := provider.New()
  if tokenError != nil {
    t.Fatal(tokenError)
  }

  var server Server
  if !t.Run(`New`, func(t *testing.T) {
    var newError error
    server, newError = New(&token)
    if newError != nil {
      t.Error(newError)
    }
  }) {
    t.FailNow()
  }

  t.Run(`ServeHTTP`, server.TestServeHTTP)
}
