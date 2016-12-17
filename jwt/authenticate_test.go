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

package jwt

import (
  `bytes`
  `encoding/base64`
  `regexp`
  `testing`
  `time`
)

func (jwt JWT) testAuthenticate(t *testing.T) {
  for _, test := range [...]struct{
        description string
        paramSignature int
        paramMessage string
        valid bool
        expectedErrorMessage string
        expectedErrorURI string
      }{
        {
          `manyElements`, 2,
          `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJzdWIiLCJleHAiOjQyOTQ5NjcyOTYsImp0aSI6IjU1NzcwMDY3OTE5NDc3Nzk0MTAifQ`,
          false, `expected 3 parts, got 4 parts`,
          `https://tools.ietf.org/html/rfc7519#section-3.1`,
        }, {
          `malformedHeader`, 1,
          `bWFsZm9ybWVk.eyJzdWIiOiJzdWIiLCJleHAiOjQyOTQ5NjcyOTYsImp0aSI6IjU1NzcwMDY3OTE5NDc3Nzk0MTAifQ`,
          false, `invalid character 'm' looking for beginning of value`,
          `https://tools.ietf.org/html/rfc7159`,
        }, {
          `multipleHeader`, 1,
          `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJzdWIiLCJleHAiOjQyOTQ5NjcyOTYsImp0aSI6IjU1NzcwMDY3OTE5NDc3Nzk0MTAifQ`,
          false, `header contains something superfluous`,
          `https://tools.ietf.org/html/rfc7159`,
        }, {
          `nonHS256Header`, 1,
          `eyJ0eXAiOiJKV1QiLCJhbGciOiJpbnZhbGlkIn0.eyJzdWIiOiJzdWIiLCJleHAiOjQyOTQ5NjcyOTYsImp0aSI6IjU1NzcwMDY3OTE5NDc3Nzk0MTAifQ`,
          false, `expected alg "HS256", header says "invalid"`,
          `https://tools.ietf.org/html/rfc7515#section-4.1.1`,
        }, {
          `malformedClaim`, 1,
          `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.bWFsZm9ybWVk`,
          false, `invalid character 'm' looking for beginning of value`,
          `https://tools.ietf.org/html/rfc7159`,
        }, {
          `multipleClaim`, 1,
           `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJzdWIiLCJleHAiOjQyOTQ5NjcyOTYsImp0aSI6IjU1NzcwMDY3OTE5NDc3Nzk0MTAifXsic3ViIjoic3ViIiwiZXhwIjo0Mjk0OTY3Mjk2LCJqdGkiOiI1NTc3MDA2NzkxOTQ3Nzc5NDEwIn0`,
          false, `claim contains something superfluous`,
          `https://tools.ietf.org/html/rfc7159`,
        }, {
          `invalidExpClaim`, 1,
          `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJzdWIiLCJleHAiOjAsImp0aSI6IjU1NzcwMDY3OTE5NDc3Nzk0MTAifQ`,
          false, `claim is expired; it is .+, claim says expired in .+`,
          `https://tools.ietf.org/html/rfc7519#section-4.1.4`,
        }, {
          `malformedSignature`, 0,
          `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJzdWIiLCJleHAiOjQyOTQ5NjcyOTYsImp0aSI6IjU1NzcwMDY3OTE5NDc3Nzk0MTAifQ.%`,
          false, `illegal base64 data at input byte 0`,
          `https://tools.ietf.org/html/rfc4648#section-5`,
        }, {
          `invalidSignature`, 0,
          `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJzdWIiLCJleHAiOjQyOTQ5NjcyOTYsImp0aSI6IjU1NzcwMDY3OTE5NDc3Nzk0MTAifQ.AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA`,
          false, `invalid signature`,
          `https://tools.ietf.org/html/rfc7515#section-5`,
        }, {
          `valid`, 1,
          `eyJ0eXAiOiJpbnZhbGlkIiwiYWxnIjoiSFMyNTYifQ.eyJzdWIiOiJzdWIiLCJleHAiOjQyOTQ5NjcyOTYsImp0aSI6IjU1NzcwMDY3OTE5NDc3Nzk0MTAifQ`,
          true, ``, ``,
        },
      } {
    test := test

    t.Run(test.description, func(t *testing.T) {
      sum := bytes.NewBuffer([]byte(test.paramMessage))

      hash := jwt.authority.Hash()
      hash.Write(sum.Bytes())

      var signature [43]byte
      base64.RawURLEncoding.Encode(signature[:], hash.Sum(nil))

      for count := 0;
          count < test.paramSignature;
          count++ {
        sum.Write([]byte {'.'})
        sum.Write(signature[:])
      }

      exp := time.Unix(4294967296, 0)
      before := exp.Sub(time.Now())
      authenticated, authenticateError := jwt.Authenticate(sum.String())
      after := exp.Sub(time.Now())
      if !authenticateError.IsError() {
        if test.expectedErrorMessage != `` {
          t.Errorf(`expected error %q, got no error`, test.expectedErrorMessage)
        }
      } else if test.expectedErrorMessage == `` {
        t.Errorf(`expected error (nil), got %q`, authenticateError.Error())
      } else {
        message := authenticateError.Error();
        matched, matchError := regexp.MatchString(test.expectedErrorMessage,
                                                  message);
        if matchError != nil {
          t.Error(matchError)
        } else if !matched {
          t.Errorf(`invalid error; expected to match with %q, got %q`,
                   test.expectedErrorMessage, message)
        }
      }

      if uri := authenticateError.URI(); uri != test.expectedErrorURI {
        t.Errorf(`invalid error uri; expected %q, got %q`,
                 test.expectedErrorURI, uri)
      }

      if test.valid {
        if authenticated.Sub != `sub` {
          t.Errorf(`invalid subject; expected "sub", got %q`, authenticated)
        }

        if authenticated.Duration < after || authenticated.Duration > before {
          t.Errorf(`invalid duration: %v`, authenticated.Duration)
        }
      } else {
        if authenticated.Sub != `` {
          t.Errorf(`invalid subject; expected empty (""), got %q`,
                   authenticated.Sub)
        }

        if authenticated.Duration != 0 {
          t.Errorf(`invalid duration; expected empty (0), got %q`,
                   authenticated.Duration)
        }
      }
    })
  }
}
