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

package scope

import (
  `github.com/kagucho/tsubonesystem3/scope`
  `testing`
)

var tests = []struct {
  description string
  encoded string
  decoded scope.Scope
}{
  {
    `basic`, `basic`, scope.Scope{}.Set(scope.Basic),
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
      if decoded, decodeError := Decode(test.encoded); decodeError != nil {
        t.Error(decodeError)
      } else if decoded != test.decoded {
        t.Errorf(`expected %b; got %b`, test.decoded, decoded)
      }
    })
  }

  t.Run(`empty`, func(t *testing.T) {
    t.Parallel()

    decoded, decodeError := Decode(``)
    if decodeError != nil {
      t.Fatal(decodeError)
    }

    if decoded.IsSetAny() {
      t.Errorf(`expected nothing is set, got %b`, decoded)
    }
  })

  t.Run(`unknown`, func(t *testing.T) {
    t.Parallel()

    decoded, decodeError := Decode(`unknown`);
    if decodeError == nil {
      t.Error(`expected error "unknown scope: \"unknown\"", got (nil)`)
    } else if decodeError.Error() != `unknown scope: "unknown"` {
      t.Errorf(`expected error "unknown scope: \"unknown\"", got %q`,
               decodeError)
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
      if encoded, encodeError := Encode(test.decoded); encodeError != nil {
        t.Error(encodeError)
      } else if encoded != test.encoded {
        t.Errorf(`expected %q; got %q`, test.encoded, encoded)
      }
    })
  }
}
