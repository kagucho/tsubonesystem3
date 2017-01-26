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

package unchunked

import (
	"sort"
	"testing"
)

type value struct {
	coding codingIndex
	q      uint
}

/*
	RFC 7230 - Hypertext Transfer Protocol (HTTP/1.1): Message Syntax and Routing
	3.2.3.  Whitespace
	https://tools.ietf.org/html/rfc7230#section-3.2.3
*/
func TestIsWhitespace(t *testing.T) {
	t.Parallel()

	for _, test := range [...]struct {
		description string
		character   byte
		expected    bool
	}{
		// OWS            = *( SP / HTAB )
		{`SP`, ' ', true}, {`HTAB`, '\t', true},

		{`!`, '!', false},
	} {
		test := test

		t.Run(string(test.character), func(t *testing.T) {
			t.Parallel()

			if result := isWhitespace(test.character); result != test.expected {
				t.Errorf(`expected %v, got %v`,
					test.expected, result)
			}
		})
	}
}

func TestCodings(t *testing.T) {
	t.Parallel()

	if !sort.StringsAreSorted(codings) {
		t.Error(`expected codings are sorted, they are not`)
	}
}

func TestParseCodings(t *testing.T) {
	t.Parallel()

	for _, coding := range [...]struct {
		coding        string
		expectedIndex codingIndex
	}{
		/*
			4.2.3.  Gzip Coding
			https://tools.ietf.org/html/rfc7230#section-4.2.3
			> The "gzip" coding is an LZ77 coding with a 32-bit
			> Cyclic Redundancy Check (CRC) that is commonly
			> produced by the gzip file compression program
			> [RFC1952].
		*/
		{`GZIP`, codingGzip},

		/*
			RFC 7231 - Hypertext Transfer Protocol (HTTP/1.1): Semantics and Content
			5.3.4.  Accept-Encoding
			https://tools.ietf.org/html/rfc7231#section-5.3.4
			> An "identity" token is used as a synonym for
			> "no encoding" in order to communicate when no
			> encoding is preferred.
		*/
		{`IDENTITY`, codingIdentity},

		/*
			> The asterisk "*" symbol in an Accept-Encoding field
			> matches any available content-coding not explicitly
			> listed in the header field.
		*/
		{`*`, codingAny},

		{`DEFLATE`, codingUnknown}, {`~`, codingUnknown},
	} {
		coding := coding

		for _, delimiter := range [...]string{
			"", ",", ";", " ", "\t",
		} {
			field := coding.coding + delimiter

			t.Run(field, func(t *testing.T) {
				t.Parallel()

				index, cursor := parseCodings(field, 0)

				if index != coding.expectedIndex {
					t.Errorf(`invalid codingIndex; expected %v, got %v`,
						coding.expectedIndex, index)
				}

				if cursor != len(coding.coding) {
					t.Errorf(`invalid cursor; expected %v, got %v`,
						len(coding.coding), cursor)
				}
			})
		}
	}
}

func TestParseQvalue(t *testing.T) {
	t.Parallel()

	for _, qvalue := range [...]struct {
		field     string
		expectedQ uint
	}{
		/*
			5.3.1.  Quality Values
			https://tools.ietf.org/html/rfc7231#section-5.3.1
			qvalue = ( "0" [ "." 0*3DIGIT ] ) / ( "1" [ "." 0*3("0") ] )
		*/
		{`0`, 0x00000000}, {`0.`, 0x00000000},
		{`0.1`, 0x00010000}, {`0.12`, 0x00010200},
		{`0.123`, 0x00010203},
		{`1`, 0x01000000}, {`1.`, 0x01000000},
		{`1.0`, 0x01000000}, {`1.00`, 0x01000000},
		{`1.000`, 0x01000000},
	} {
		for _, delimiter := range [...]struct {
			field          string
			expectedCursor int
		}{
			{``, 0}, {` `, 1}, {`,`, 0},
		} {
			field := qvalue.field + delimiter.field
			expectedQ := qvalue.expectedQ
			expectedCursor := len(qvalue.field) + delimiter.expectedCursor

			t.Run(field, func(t *testing.T) {
				t.Parallel()

				q, cursor := parseQvalue(field, 0)

				if q != expectedQ {
					t.Errorf(`invalid q; expected %v, got %v`,
						expectedQ, q)
				}

				if cursor != expectedCursor {
					t.Errorf(`invalid cursor; expected %v, got %v`,
						expectedCursor, cursor)
				}
			})
		}
	}
}

func TestParseWeight(t *testing.T) {
	t.Parallel()

	for _, test := range [...]struct {
		field          string
		expectedQ      uint
		expectedCursor int
	}{
		// weight = OWS ";" OWS "q=" qvalue
		{` `, 0x01000000, 1}, {` ,`, 0x01000000, 1},
		{` ; q=0`, 0x00000000, 6},
	} {
		test := test

		t.Run(test.field, func(t *testing.T) {
			t.Parallel()

			q, cursor := parseWeight(test.field, 0)

			if q != test.expectedQ {
				t.Errorf(`invalid q; expected %v, got %v`,
					test.expectedQ, q)
			}

			if cursor != test.expectedCursor {
				t.Errorf(`invalid cursor; expected %v, got %v`,
					test.expectedCursor, cursor)
			}
		})
	}
}

/*
	5.3.4.  Accept-Encoding
	https://tools.ietf.org/html/rfc7231#section-5.3.4
	codings          = Accept-Encoding = [ ( "," / ( codings [ weight ] ) ) *( OWS "," [ OWS ( codings [ weight ] ) ] ) ]
*/
func TestParseField(t *testing.T) {
	t.Parallel()

	// [
	t.Run(``, func(t *testing.T) {
		t.Parallel()

		parseField(``, nil)
	})

	// ","
	t.Run(`,`, func(t *testing.T) {
		t.Parallel()

		parseField(`,`, nil)
	})

	expected := [...]value{
		{codingGzip, 0x01000000}, {codingUnknown, 0x00000000},
	}
	for _, test := range [...]string{
		`,, gzip;q=1, deflate;q=0`, `gzip;q=1, deflate;q=0`,
	} {
		test := test

		t.Run(test, func(t *testing.T) {
			t.Parallel()

			var result [len(expected)]value
			count := 0

			parseField(test, func(index codingIndex, q uint) {
				if count < len(result) {
					result[count] = value{index, q}
					count++
				}
			})

			if count != len(result) {
				t.Error(`expected`, len(expected), `codings`)
			} else if result != expected {
				t.Errorf(`expected %v, got %v`,
					expected, result)
			}
		})
	}
}

/*
	RFC 7230 - Hypertext Transfer Protocol (HTTP/1.1): Message Syntax and Routing
	3.2.2.  Field Order
	https://tools.ietf.org/html/rfc7230#section-3.2.2
	> A sender MUST NOT generate multiple header fields with the same field
	> name in a message unless either the entire field value for that header
	> field is defined as a comma-separated list [i.e., #(values)] or the
	> header field is a well-known exception (as noted below).
*/
func TestParse(t *testing.T) {
	t.Parallel()

	expected := [...]value{
		{codingGzip, 0x01000000}, {codingUnknown, 0x00000000},
	}
	var result [len(expected)]value
	count := 0

	parse([]string{`gzip;q=1`, `deflate;q=0`}, func(index codingIndex, q uint) {
		if count < len(result) {
			result[count] = value{index, q}
			count++
		}
	})

	if count != len(result) {
		t.Error(`expected`, len(expected), `codings`)
	} else if result != expected {
		t.Errorf(`expected %v, got %v`, expected, result)
	}
}
