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

package jwt

import (
	"strconv"
	"time"
)

// Time is a type to represent a time.
type Time struct {
	time.Time
}

/*
UnmarshalJSON sets *marshallable to unmarshalled value from JSON.

This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (marshallable *Time) UnmarshalJSON(data []byte) error {
	parsed, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}

	*marshallable = Time{time.Unix(parsed, 0).In(time.Local)}

	return nil
}

/*
MarshalJSON returns the JSON encoding of the value.

This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (marshallable Time) MarshalJSON() ([]byte, error) {
	// RFC 7519 NumericDate
	return strconv.AppendInt(nil, marshallable.Unix(), 10), nil
}
