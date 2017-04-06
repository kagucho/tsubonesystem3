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

// Package encoding provides encoding features.
package encoding

import (
	"bytes"
	"encoding/json"
	"github.com/kagucho/tsubonesystem3/jwt"
	"log"
	"reflect"
	"time"
)

// Time is a type representing a time and marshallable as JSON.
type Time struct {
	jwt.Time
}

/*
ZeroString is a string which should be marshalled as null if it has zero value.
*/
type ZeroString string

/*
ZeroUint16 is a uint16 which should be marshalled as null if it has zero value.
*/
type ZeroUint16 uint16

// ParseQueryTime parses time in URL query.
func ParseQueryTime(encoded string) (Time, error) {
	var time Time
	err := time.UnmarshalJSON([]byte(encoded))
	return time, err
}

func (specific Time) Generic() time.Time {
	return specific.Time.Time
}

/*
MarshalJSON returns the JSON encoding of the value.

This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (zeroString ZeroString) MarshalJSON() ([]byte, error) {
	return marshalJSONNullFromZero(string(zeroString))
}

/*
MarshalJSON returns the JSON encoding of the value.

This implements an interface used in encoding/encoding.

json - The Go Programming Language
Example (CustomMarshalJSON)
https://golang.org/pkg/encoding/json/#example__customMarshalJSON
*/
func (zeroUint16 ZeroUint16) MarshalJSON() ([]byte, error) {
	return marshalJSONNullFromZero(uint16(zeroUint16))
}

/*
MarshalJSONArray marshals a JSON array from values returned by the given
function.

The given function returns the value, error, and a boolean telling whether
a new value is present.

If the error is not nil, it stops marshalling and calls the function as long as
it says a new value is present. That behavior is convenient to marshal channels,
which are often needed to be drained to release resources.
*/
func MarshalJSONArray(callback func() (interface{}, error, bool)) ([]byte, error) {
	defer func() {
		_, err, _ := callback()
		if err != nil {
			log.Print(err)
		}
	}()

	buffer := bytes.NewBuffer(make([]byte, 0, 2))
	buffer.WriteByte('[')

	recievedValue, err, present := callback()
	if !present {
		buffer.WriteByte(']')
	} else {
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)

		for {
			if err != nil {
				return nil, err
			}

			if err := encoder.Encode(recievedValue); err != nil {
				return nil, err
			}

			recievedValue, err, present = callback()
			if !present {
				break
			}

			bufferBytes := buffer.Bytes()
			bufferBytes[len(bufferBytes)-1] = ','
		}

		bufferBytes := buffer.Bytes()
		bufferBytes[len(bufferBytes)-1] = ']'
	}

	return buffer.Bytes(), nil
}

/*
MarshalJSONObject marshals a JSON object from values returned by the given
function.

The given function returns the key, value, error, and a boolean telling whether
a new value is present.

If the error is not nil, it stops marshalling and calls the function as long as
it says a new value is present. That behavior is convenient to marshal channels,
which are often needed to be drained to release resources.
*/
func MarshalJSONObject(callback func() (string, interface{}, error, bool)) ([]byte, error) {
	defer func() {
		_, _, err, _ := callback()
		if err != nil {
			log.Print(err)
		}
	}()

	buffer := bytes.NewBuffer(make([]byte, 0, 2))
	buffer.WriteByte('{')

	key, recievedValue, err, present := callback()
	if !present {
		buffer.WriteByte('}')
	} else {
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)

		for {
			if err != nil {
				return nil, err
			}

			buffer.WriteByte('"')
			buffer.WriteString(key)
			buffer.WriteString(`":`)

			if err := encoder.Encode(recievedValue); err != nil {
				return nil, err
			}

			key, recievedValue, err, present = callback()
			if !present {
				break
			}

			bufferBytes := buffer.Bytes()
			bufferBytes[len(bufferBytes)-1] = ','
		}

		bufferBytes := buffer.Bytes()
		bufferBytes[len(bufferBytes)-1] = '}'
	}

	return buffer.Bytes(), nil
}

func NewTime(unwrapped time.Time) Time {
	return Time{jwt.Time{unwrapped}}
}

func marshalJSONNullFromZero(unmarshalled interface{}) ([]byte, error) {
	if unmarshalled == reflect.Zero(reflect.ValueOf(unmarshalled).Type()).Interface() {
		return []byte(`null`), nil
	}

	return json.Marshal(unmarshalled)
}
