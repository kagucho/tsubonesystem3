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

package member

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/kagucho/tsubonesystem3/db"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/common"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/context"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/token/authorizer"
	"github.com/kagucho/tsubonesystem3/scope"
	"net/http"
)

type jsonMessage []byte

type publicDetail struct {
	Affiliation string              `json:"affiliation,omitempty"`
	Clubs       jsonMessage         `json:"clubs"`
	Entrance    uint16              `json:"entrance,omitempty"`
	Gender      string              `json:"gender,omitempty"`
	Mail        string              `json:"mail"`
	Nickname    string              `json:"nickname"`
	OB          bool                `json:"ob"`
	Positions   jsonMessage         `json:"positions"`
	Realname    string              `json:"realname,omitempty"`
}

type privateDetail struct {
	publicDetail
	Tel string `json:"tel,omitempty"`
}

func (message jsonMessage) MarshalJSON() ([]byte, error) {
	return message, nil
}

func position(member db.Member) (bool, jsonMessage, jsonMessage) {
	chief := false

	clubs := bytes.NewBuffer(make(jsonMessage, 0, 2))
	clubs.WriteByte('[')
	clubsEncoder := json.NewEncoder(clubs)
	clubsPresent := false

	positions := bytes.NewBuffer(make(jsonMessage, 0, 2))
	positions.WriteByte('[')
	positionsEncoder := json.NewEncoder(positions)
	positionsPresent := false

	for member.Clubs != nil || member.Positions != nil {
		select {
		case result, present := <-member.Clubs:
			if present {
				if result.Error != nil {
					panic(result.Error)
				}

				if clubsPresent {
					clubsBytes := clubs.Bytes()
					clubsBytes[len(clubsBytes) - 1] = ','
				} else {
					clubsPresent = true
				}

				if encodeError := clubsEncoder.Encode(result.Value); encodeError != nil {
					panic(encodeError)
				}

				if result.Value.Chief {
					chief = true
				}
			} else {
				if clubsPresent {
					clubsBytes := clubs.Bytes()
					clubsBytes[len(clubsBytes) - 1] = ']'
				} else {
					clubs.WriteByte(']')
				}

				member.Clubs = nil
			}

		case result, present := <-member.Positions:
			if present {
				if result.Error != nil {
					panic(result.Error)
				}

				if positionsPresent {
					positionsBytes := positions.Bytes()
					positionsBytes[len(positionsBytes) - 1] = ','
				} else {
					positionsPresent = true
				}

				if encodeError := positionsEncoder.Encode(result.Value); encodeError != nil {
					panic(encodeError)
				}
			} else {
				if positionsPresent {
					positionsBytes := positions.Bytes()
					positionsBytes[len(positionsBytes) - 1] = ']'
				} else {
					positions.WriteByte(']')
				}

				member.Positions = nil
			}
		}
	}

	return chief || positionsPresent, clubs.Bytes(), positions.Bytes()
}

// DetailServeHTTP serves the detail of the member identified with the given ID
// via HTTP.
func DetailServeHTTP(writer http.ResponseWriter, request *http.Request,
	context context.Context, claim authorizer.Claim) {
	serve := func() func() {
		defer common.Recover(writer)

		id := request.FormValue(`id`)
		member, queryError := context.DB.QueryMember(id)

		switch queryError {
		case nil:
			inPosition, clubs, positions := position(member)

			public := publicDetail{
				Affiliation: member.Affiliation,
				Clubs:       clubs,
				Entrance:    member.Entrance,
				Gender:      member.Gender,
				Mail:        member.Mail,
				Nickname:    member.Nickname,
				OB:          member.OB,
				Positions:   positions,
				Realname:    member.Realname,
			}

			var unmarshalled interface{}
			if inPosition || claim.Sub == id || claim.Scope.IsSet(scope.Privacy) {
				unmarshalled = privateDetail{public, member.Tel}
			} else {
				unmarshalled = public
			}

			return func() {
				common.ServeJSON(writer, unmarshalled,
					http.StatusOK)
			}

		case sql.ErrNoRows:
			return func() {
				common.ServeError(writer,
					common.Error{
						ID: `invalid_id`,
						Description: `invalid ID`,
					}, http.StatusBadRequest)
			}

		default:
			panic(queryError)
		}
	}()

	if serve != nil {
		serve()
	}
}
