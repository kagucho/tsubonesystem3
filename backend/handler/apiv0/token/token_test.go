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

package token

import (
	"github.com/kagucho/tsubonesystem3/backend/scope"
	"testing"
)

var tests = []struct {
	description string
	encoded     string
	decoded     scope.Scope
}{
	{
		`basic`,

		/*
			RFC 6749 - The OAuth 2.0 Authorization Framework
			A.4.  "scope" Syntax
			https://tools.ietf.org/html/rfc6749#appendix-A.4

			> The "scope" element is defined in Section 3.3:

			>   scope-token = 1*NQCHAR
			>   scope       = scope-token *( SP scope-token )
		*/
		`basic`,

		scope.Scope{}.Set(scope.Basic),
	}, {
		`privacy`, `privacy`, scope.Scope{}.Set(scope.Privacy),
	}, {
		`management`, `management`, scope.Scope{}.Set(scope.Management),
	}, {
		`multiple`, `basic privacy`,
		scope.Scope{}.Set(scope.Basic).Set(scope.Privacy),
	},
}

func TestDecode(t *testing.T) {
	t.Parallel()

	for _, test := range tests {
		test := test

		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			if decoded, err := Decode(test.encoded); err != nil {
				t.Error(err)
			} else if decoded != test.decoded {
				t.Errorf(`expected %b; got %b`, test.decoded, decoded)
			}
		})
	}

	t.Run(`empty`, func(t *testing.T) {
		t.Parallel()

		decoded, err := Decode(``)
		if decodeError != nil {
			t.Fatal(err)
		}

		if decoded.IsSetAny() {
			t.Errorf(`expected nothing is set, got %b`, decoded)
		}
	})

	t.Run(`unknown`, func(t *testing.T) {
		t.Parallel()

		decoded, err := Decode(`unknown`)
		if err == nil {
			t.Error(`expected error "unknown scope: \"unknown\"", got (nil)`)
		} else if err.Error() != `unknown scope: "unknown"` {
			t.Errorf(`expected error "unknown scope: \"unknown\"", got %q`,
				err)
		}

		if decoded.IsSetAny() {
			t.Errorf(`expected nothing is set, got %q`, decoded)
		}
	})
}

func TestEncode(t *testing.T) {
	t.Parallel()

	for _, test := range tests {
		test := test

		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			if encoded, err := Encode(test.decoded); err != nil {
				t.Error(err)
			} else if encoded != test.encoded {
				t.Errorf(`expected %q; got %q`,
					test.encoded, encoded)
			}
		})
	}
}
