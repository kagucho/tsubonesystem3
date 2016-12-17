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

package authorizer

import (
  `github.com/kagucho/tsubonesystem3/db`
  `github.com/kagucho/tsubonesystem3/handler/apiv0/token/provider`
  `github.com/kagucho/tsubonesystem3/scope`
  `net/http`
  `net/http/httptest`
  `testing`
)

func (authorizer Authorizer) TestAuthorize(t *testing.T) {
  request := httptest.NewRequest(`GET`, `https://kagucho.net/`, nil)

  testError := func(code int, authenticate string, body string) {
    recorder := httptest.NewRecorder()
    authorizer.Authorize(recorder, request, db.DB{}, nil)

    if recorder.Code != code {
      t.Error(`invalid status code; expected `, code, `, got `, recorder.Code)
    }

    if result := recorder.HeaderMap.Get(`WWW-Authenticate`);
       result != authenticate {
      t.Errorf(`invalid WWW-Authenticate in header; expected %q, got %q`,
               result, authenticate)
    }

    if result := recorder.Body.String(); result != body {
      t.Error(`invalid body; expected `, body, `, got `, recorder.Body)
    }
  }

  t.Run(`none`, func(t *testing.T) {
    testError(http.StatusUnauthorized, `Bearer scope=basic`,
      `{"error":"invalid_token","error_description":"expected bearer authentication scheme, got \"\" in Authorization field of request header","error_uri":"https://tools.ietf.org/html/rfc6750#section-2.1","scope":"basic"}`)
  })

  t.Run(`invalidToken`, func(t *testing.T) {
    request.Header.Set(`Authorization`, `Bearer invalid`)
    testError(http.StatusUnauthorized,
      `Bearer error="invalid_token",error_description="expected 3 parts, got 1 parts",error_uri="https://tools.ietf.org/html/rfc7519#section-3.1",scope=basic`,
      `{"error":"invalid_token","error_description":"expected 3 parts, got 1 parts","error_uri":"https://tools.ietf.org/html/rfc7519#section-3.1","scope":"basic"}`)
  })

  t.Run(`invalidScope`, func(t *testing.T) {
    issued, issueError := authorizer.jwt.Issue(`sub`, `invalid`, 1073741824)
    if issueError != nil {
      t.Fatal(issueError)
    }

    request.Header.Set(`Authorization`, `Bearer ` + issued)
    testError(http.StatusUnauthorized,
      `Bearer error="invalid_token",error_description="unknown scope: \"invalid\"",error_uri="https://tools.ietf.org/html/rfc6749#section-7.2",scope=basic`,
      `{"error":"invalid_token","error_description":"unknown scope: \"invalid\"","error_uri":"https://tools.ietf.org/html/rfc6749#section-7.2","scope":"basic"}`)
  })

  t.Run(`insufficientScope`, func(t *testing.T) {
    issued, issueError := authorizer.jwt.Issue(`sub`, ``, 1073741824)
    if issueError != nil {
      t.Fatal(issueError)
    }

    request.Header.Set(`Authorization`, `Bearer ` + issued)
    testError(http.StatusForbidden,
      `Bearer error="insufficient_scope",error_description="The request requires higher privileges than provided by the access token.",error_uri="https://tools.ietf.org/html/rfc6750#section-3.1",scope=basic`,
      `{"error":"insufficient_scope","error_description":"The request requires higher privileges than provided by the access token.","error_uri":"https://tools.ietf.org/html/rfc6750#section-3.1","scope":"basic"}`)
  })

  t.Run(`valid`, func(t *testing.T) {
    issued, issueError := authorizer.jwt.Issue(`sub`, `basic`, 1073741824)
    if issueError != nil {
      t.Fatal(issueError)
    }

    request.Header.Set(`Authorization`, `Bearer ` + issued)
    recorder := httptest.NewRecorder()
    authorizer.Authorize(recorder, request, db.DB{},
                         func(writer http.ResponseWriter,
                              request *http.Request, db db.DB, claim Claim) {
      if claim.Sub != `sub` {
        t.Errorf(`invalid Sub; expected "sub", got %q`, claim.Sub)
      }

      if (claim.Scope != scope.Scope{}.Set(scope.Basic)) {
        t.Error(`invalid scope; expected `, claim.Scope,
                `, got `, scope.Scope{}.Set(scope.Basic))
      }

      writer.WriteHeader(http.StatusOK)
    })

    if recorder.Code != http.StatusOK {
      t.Error(`invalid status code; expected `, http.StatusOK,
              `, got `, recorder.Code)
    }
  })
}

func TestAuthorizer(t *testing.T) {
  token, tokenError := provider.New()
  if tokenError != nil {
    t.Fatal(tokenError)
  }

  var authorizer Authorizer
  if !t.Run(`New`, func(t *testing.T) {
    authorizer = New(&token)
  }) {
    t.FailNow()
  }

  t.Run(`Authorize`, authorizer.TestAuthorize)
}
