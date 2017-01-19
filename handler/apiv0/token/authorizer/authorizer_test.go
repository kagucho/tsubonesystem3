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

package authorizer

import (
	"github.com/kagucho/tsubonesystem3/db"
	`github.com/kagucho/tsubonesystem3/handler/apiv0/common`
	"github.com/kagucho/tsubonesystem3/handler/apiv0/token/provider"
	"github.com/kagucho/tsubonesystem3/scope"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (authorizer Authorizer) TestAuthorize(t *testing.T) {
	request := httptest.NewRequest(`GET`, `https://kagucho.net/`, nil)

	testError := func(code int, authenticate string, body string) {
		recorder := httptest.NewRecorder()
		authorizer.Authorize(recorder, request, db.DB{}, nil)

		if recorder.Code != code {
			t.Error(`invalid status code; expected `, code,
				`, got `, recorder.Code)
		}

		/*
			RFC 6750 - The OAuth 2.0 Authorization Framework: Bearer Token Usage
			3.  The WWW-Authenticate Response Header Field
			https://tools.ietf.org/html/rfc6750#section-3
			> If the protected resource request does not include
			> authentication credentials or does not contain an
			> access token that enables access to the protected
			> resource, the resource server MUST include the HTTP
			> "WWW-Authenticate" response header field; it MAY
			> include it in response to other conditions as well.
			> The "WWW-Authenticate" header field uses the framework
			> defined by HTTP/1.1 [RFC2617].
		*/
		if result := recorder.HeaderMap.Get(`WWW-Authenticate`); result != authenticate {
			t.Errorf(`invalid WWW-Authenticate in header; expected %q, got %q`,
				result, authenticate)
		}

		if result := recorder.Body.String(); result != body {
			t.Error(`invalid body; expected `, body,
				`, got `, recorder.Body)
		}
	}

	t.Run(`none`, func(t *testing.T) {
		/*
			> All challenges defined by this specification MUST use
			> the auth-scheme value "Bearer".  This scheme MUST be
			> followed by one or more auth-param values.  The
			> auth-param attributes used or defined by this
			> specification are as follows.  Other auth-param
			> attributes MAY be used as well.

			3.1.  Error Codes
			https://tools.ietf.org/html/rfc6750#section-3.1

			> The access token provided is expired, revoked,
			> malformed, or invalid for other reasons.
			> The resource SHOULD respond with the HTTP 401
			> (Unauthorized) status code.
		*/
		testError(http.StatusUnauthorized, `Bearer scope=basic`,
			`{"error":"invalid_token","error_description":"expected bearer authentication scheme","error_uri":"https://tools.ietf.org/html/rfc6750#section-2.1","scope":"basic"}
`)
	})

	t.Run(`invalidToken`, func(t *testing.T) {
		/*
			2.1.  Authorization Request Header Field
			https://tools.ietf.org/html/rfc6750#section-2.1
			> When sending the access token in the "Authorization"
			> request header field defined by HTTP/1.1 [RFC2617],
			> the client uses the "Bearer" authentication scheme
			> to transmit the access token.
		*/
		request.Header.Set(`Authorization`, `Bearer invalid`)

		testError(http.StatusUnauthorized,
			`Bearer error="invalid_token",error_description="expected 3 parts, got 1 parts",error_uri="https://tools.ietf.org/html/rfc7519#section-3.1",scope=basic`,
			`{"error":"invalid_token","error_description":"expected 3 parts, got 1 parts","error_uri":"https://tools.ietf.org/html/rfc7519#section-3.1","scope":"basic"}
`)
	})

	t.Run(`invalidScope`, func(t *testing.T) {
		issued, issueError :=
			authorizer.jwt.Issue(`sub`, `invalid`, 1073741824)
		if issueError != nil {
			t.Fatal(issueError)
		}

		request.Header.Set(`Authorization`, `Bearer `+issued)

		testError(http.StatusUnauthorized,
			`Bearer error="invalid_token",error_description="unknown scope: %22invalid%22",error_uri="https://tools.ietf.org/html/rfc6749#section-7.2",scope=basic`,
			`{"error":"invalid_token","error_description":"unknown scope: %22invalid%22","error_uri":"https://tools.ietf.org/html/rfc6749#section-7.2","scope":"basic"}
`)
	})

	t.Run(`insufficientScope`, func(t *testing.T) {
		issued, issueError :=
			authorizer.jwt.Issue(`sub`, ``, 1073741824)
		if issueError != nil {
			t.Fatal(issueError)
		}

		request.Header.Set(`Authorization`, `Bearer `+issued)

		/*
			> The request requires higher privileges than provided
			> by the access token.  The resource server SHOULD
			> respond with the HTTP 403 (Forbidden) status code and
			> MAY include the "scope" attribute with the scope
			> necessary to access the protected resource.
		*/
		testError(http.StatusForbidden,
			`Bearer error="insufficient_scope",error_description="The request requires higher privileges than provided by the access token.",error_uri="https://tools.ietf.org/html/rfc6750#section-3.1",scope=basic`,
			`{"error":"insufficient_scope","error_description":"The request requires higher privileges than provided by the access token.","error_uri":"https://tools.ietf.org/html/rfc6750#section-3.1","scope":"basic"}
`)
	})

	t.Run(`valid`, func(t *testing.T) {
		issued, issueError :=
			authorizer.jwt.Issue(`sub`, `basic`, 1073741824)
		if issueError != nil {
			t.Fatal(issueError)
		}

		request.Header.Set(`Authorization`, `Bearer `+issued)
		recorder := httptest.NewRecorder()
		authorizer.Authorize(recorder, request, db.DB{},
			func(writer http.ResponseWriter,
				request *http.Request, db db.DB, claim Claim) {
				if claim.Sub != `sub` {
					t.Errorf(`invalid Sub; expected "sub", got %q`,
						claim.Sub)
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
	t.Parallel()

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

func TestServeError(t *testing.T) {
	t.Parallel()

	recorder := httptest.NewRecorder()
	serveError(recorder, common.Error{
		/*
			> Values for the "error" and "error_description"
			> attributes (specified in Appendixes A.7 and A.8 of
			> [RFC6749]) MUST NOT include characters outside the set
			> %x20-21 / %x23-5B / %x5D-7E.

			RFC 6749 - The OAuth 2.0 Authorization Framework
			5.2.  Error Response
			https://tools.ietf.org/html/rfc6749#section-5.2

			> Values for the "error" parameter MUST
			> NOT include characters outside the set
			> %x20-21 / %x23-5B / %x5D-7E.
		*/
		"\x00",

/*
			> Values for the "error_description" parameter
			> MUST NOT include characters outside the set
			> %x20-21 / %x23-5B / %x5D-7E.
*/
		"\x20",

		/*
			RFC 6750 - The OAuth 2.0 Authorization Framework: Bearer Token Usage
			3.  The WWW-Authenticate Response Header Field
			https://tools.ietf.org/html/rfc6750#section-3
			> Values for the "error_uri" attribute (specified in
			> Appendix A.9 of [RFC6749]) MUST conform to the
			> URI-reference syntax and thus MUST NOT include
			> characters outside the set %x21 / %x23-5B / %x5D-7E.

			RFC 6749 - The OAuth 2.0 Authorization Framework
			5.2.  Error Response
			https://tools.ietf.org/html/rfc6749#section-5.2
			> Values for the "error_uri" parameter MUST
			> conform to the URI-reference syntax and thus
			> MUST NOT include characters outside the set
			> %x21 / %x23-5B / %x5D-7E.
		*/
		"\x22",
	}, http.StatusUnauthorized)

	const expectedChallenge = `Bearer error="%00",error_description=" ",error_uri="%22",scope=basic`
	if result := recorder.HeaderMap.Get(`WWW-Authenticate`); result != expectedChallenge {
		t.Errorf(`invalid WWW-Authenticate field in header; expected %q, got %q`,
			expectedChallenge, result)
	}

	const expectedCode = http.StatusUnauthorized
	if result := recorder.Code; result != expectedCode {
		t.Errorf(`invalid status code; expected %v, got %v`,
			expectedCode, result)
	}
}
