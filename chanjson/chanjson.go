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

// Package chanjson implements a JSON Marshaler for channels.
package chanjson

import (
	"bytes"
	"encoding/json"
	"reflect"
)

// ChanJSON is a structure which holds the context of the JSON Marshaler.
type ChanJSON struct {
	wrapped reflect.Value
}

// New returns a new chanjson.ChanJSON
func New(wrapped interface{}) ChanJSON {
	return ChanJSON{reflect.ValueOf(wrapped)}
}

// MarshalJSON returns an array of JSON marshaled from the channel or an error.
func (chanJSON ChanJSON) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 2))
	buffer.WriteByte('[')

	recieved, present := chanJSON.wrapped.Recv()
	if !present {
		buffer.WriteByte(']')
	} else {
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)

		for {
			if resultError := recieved.FieldByName(`Error`); !resultError.IsNil() {
				return nil, resultError.Interface().(error)
			}

			if encodeError := encoder.Encode(recieved.FieldByName(`Value`).Interface()); encodeError != nil {
				return nil, encodeError
			}

			recieved, present = chanJSON.wrapped.Recv()
			if !present {
				break
			}

			bufferBytes := buffer.Bytes()
			bufferBytes[len(bufferBytes) - 1] = ','
		}

		bufferBytes := buffer.Bytes()
		bufferBytes[len(bufferBytes) - 1] = ']'
	}

	return buffer.Bytes(), nil
}
