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

package encoding

import (
	"errors"
	"testing"
)

type result struct {
	Error error
	Value interface{}
}

func TestChanJSON(t *testing.T) {
	t.Parallel()

	for _, test := range [...]struct {
		description string
		results     []result
		json        string
		jsonError   string
	}{
		/*
			RFC 7159 - The JavaScript Object Notation (JSON) Data Interchange Format
			5.  Arrays
			https://tools.ietf.org/html/rfc7159#section-5
			> An array structure is represented as square brackets
			> surrounding zero or more values (or elements).
		*/
		{`none`, []result{}, `[]`, ``},

		{`Error`, []result{{Error: errors.New(`error`)}}, ``, `error`},
		{`invalid`, []result{{Value: TestChanJSON}}, ``, `json: unsupported type: func(*testing.T)`},

		// > Elements are separated by commas.
		{`multiple`, []result{{}, {}}, `[null,null]`, ``},
	} {
		test := test

		t.Run(test.description, func(t *testing.T) {
			resultChan := make(chan result)

			go func() {
				defer close(resultChan)

				for _, result := range test.results {
					resultChan <- result
				}
			}()

			var chanJSON ChanJSON

			if !t.Run(`New`, func(t *testing.T) {
				chanJSON = New(resultChan)
			}) {
				t.FailNow()
			}

			t.Run(`MarshalJSON`, func(t *testing.T) {
				json, err := chanJSON.MarshalJSON()
				if test.jsonError == `` {
					if err != nil {
						t.Error(err)
					}

					if jsonString := string(json); jsonString != test.json {
						t.Error(`expected `, test.json,
							`, got `, jsonString)
					}
				} else {
					if errorString := err.Error(); errorString != test.jsonError {
						t.Errorf(`invalid error; expected %v, got %v`,
							test.jsonError, errorString)
					}
				}
			})
		})
	}
}
