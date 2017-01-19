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
	"strings"
)

func isWhitespace(character byte) bool {
	return character == ' ' || character == '\t'
}

type codingIndex uint

const (
	codingAny codingIndex = iota
	codingGzip
	codingIdentity
	codingUnknown
)

var codings = []string{
	codingAny: `*`, codingGzip: `GZIP`, codingIdentity: `IDENTITY`,
}

func parseCodings(field string, begin int) (codingIndex, int) {
	end := begin
	for end < len(field) && field[end] != ',' && field[end] != ';' && !isWhitespace(field[end]) {
		end++
	}

	coding := strings.ToUpper(field[begin:end])
	index := codingIndex(sort.SearchStrings(codings, coding))
	if index < codingIndex(len(codings)) && codings[index] != coding {
		index = codingUnknown
	}

	return index, end
}

const qDigitWidth = 8
const qInitIndex = 3 * qDigitWidth
const qMask = 0x0F0F0F0F

func parseQvalue(field string, cursor int) (uint, int) {
	index := uint(qInitIndex)
	q := uint(field[cursor]) << index

	cursor++
	if cursor >= len(field) {
		goto end
	}

	if field[cursor] != '.' {
		goto skip
	}

	for index > 0 {
		cursor++
		if cursor >= len(field) {
			goto end
		}

		if field[cursor] < '0' {
			goto skip
		}

		index -= qDigitWidth
		q |= uint(field[cursor]) << index
	}

skip:
	for cursor < len(field) && field[cursor] != ',' {
		cursor++
	}

end:
	return q & qMask, cursor
}

func parseWeight(field string, cursor int) (uint, int) {
	for {
		if cursor >= len(field) || field[cursor] == ',' {
			return uint(1 << qInitIndex), cursor
		}

		if field[cursor] == ';' {
			break
		}

		cursor++
	}

	for {
		cursor++

		if (field[cursor] == 'q') {
			break
		}
	}

	return parseQvalue(field, cursor + 2)
}

func parseField(field string, callback func(codingIndex, uint)) {
	cursor := 0

	parseCodingWeight := func() {
		var coding codingIndex
		coding, cursor = parseCodings(field, cursor)

		var q uint
		q, cursor = parseWeight(field, cursor)

		callback(coding, q)
	}

	if len(field) == 0 {
		return
	}

	// ","
	if field[0] == ',' {
		// OWS
		for {
			cursor++
			if cursor >= len(field) {
				return
			}

			if (field[cursor] == ',') {
				break
			}
		}
	} else {
		parseCodingWeight()
	}

	for {
		if cursor >= len(field) {
			return
		}

		// Ignore comma

		for {
			cursor++
			if (!isWhitespace(field[cursor])) {
				break
			}
		}

		parseCodingWeight()
	}
}

func parse(fields []string, callback func(codingIndex, uint)) {
	for _, field := range fields {
		parseField(field, callback)
	}
}
