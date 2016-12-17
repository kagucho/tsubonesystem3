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

package apiv0

import (
  `encoding/json`
  `github.com/kagucho/tsubonesystem3/db`
  `net/http`
  `net/http/httptest`
  `strings`
  `testing`
)

func (apiv0 Apiv0) TestServeHTTP(t *testing.T) {
  var token string

  t.Run(`unauthorized`, func(t *testing.T) {
    t.Parallel()

    recorder := httptest.NewRecorder()
    request := httptest.NewRequest(`GET`, `http://kagucho.net/private`, nil)
    apiv0.ServeHTTP(recorder, request)

    if recorder.Code != http.StatusUnauthorized {
      t.Errorf(`expected %v, got %v`, http.StatusUnauthorized, recorder.Code)
    }
  })

  if !t.Run(`public`, func(t *testing.T) {
    recorder := httptest.NewRecorder()

    request := httptest.NewRequest(`POST`, `http://kagucho.net/token`,
      strings.NewReader(`grant_type=password&username=1stDisplayId&password=1stPassword`))

    request.Header[`Content-Type`] =
        []string{`application/x-www-form-urlencoded`}

    apiv0.ServeHTTP(recorder, request)

    response := struct{AccessToken string `json:"access_token"`}{}

    decoder := json.NewDecoder(recorder.Body)
    if decodeError := decoder.Decode(&response);
       decodeError != nil {
      t.Fatal(decodeError)
    }

    token = `Bearer ` + response.AccessToken
  }) {
    t.FailNow()
  }

  t.Run(`private`, func(t *testing.T) {
    t.Parallel()

    recorder := httptest.NewRecorder()
    request := httptest.NewRequest(
      `GET`, `http://kagucho.net/member/detail`, nil)

    request.Header.Set(`Authorization`, token)
    apiv0.ServeHTTP(recorder, request)

    if recorder.Code != http.StatusBadRequest {
      t.Errorf(`expected %v, got %v`, http.StatusBadRequest, recorder.Code)
    }
  })

  t.Run(`invalid`, func(t *testing.T) {
    t.Parallel()

    recorder := httptest.NewRecorder()
    request := httptest.NewRequest(
      `GET`, `http://kagucho.net/invalid_path`, nil)

    request.Header.Set(`Authorization`, token)
    apiv0.ServeHTTP(recorder, request)

    if recorder.Code != http.StatusNotFound {
      t.Errorf(`expected %v, got %v`, http.StatusNotFound, recorder.Code)
    }
  })
}

func TestApiv0(t *testing.T) {
  var apiv0 Apiv0

  db, dbError := db.Open()
  if dbError != nil {
    t.Fatal(dbError)
  }

  defer db.Close()

  if !t.Run(`New`, func(t *testing.T) {
    var newError error

    apiv0, newError = New(db)
    if newError != nil {
      t.Error(newError)
    }
  }) {
    t.FailNow()
  }

  t.Run(`ServeHTTP`, apiv0.TestServeHTTP)
}
