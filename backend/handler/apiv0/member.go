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

package apiv0

import (
	"bytes"
	"encoding/json"
	"github.com/kagucho/tsubonesystem3/backend/db"
	"github.com/kagucho/tsubonesystem3/backend/encoding"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/util"
	"github.com/kagucho/tsubonesystem3/backend/mail"
	"github.com/kagucho/tsubonesystem3/backend/scope"
	"log"
	"net/http"
	netMail "net/mail"
	"strconv"
)

type publicDetail struct {
	Affiliation encoding.ZeroString `json:"affiliation"`
	Clubs       json.RawMessage     `json:"clubs"`
	Entrance    encoding.ZeroUint16 `json:"entrance"`
	Gender      encoding.ZeroString `json:"gender"`
	Mail        string              `json:"mail"`
	Nickname    string              `json:"nickname"`
	OB          bool                `json:"ob"`
	Positions   json.RawMessage     `json:"positions"`
	Realname    encoding.ZeroString `json:"realname"`
}

type privateDetail struct {
	publicDetail
	Confirmed bool                `json:"confirmed"`
	Tel       encoding.ZeroString `json:"tel"`
}

func memberDeleteServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path == `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	if (authorize(writer, request, shared, scope.Management) == claim{}) {
		return
	}

	switch err := shared.DB.DeleteMember(request.URL.Path[1:]); err {
	case db.ErrIncorrectIdentity:
		util.ServeErrorDefault(writer, http.StatusNotFound)

	case db.ErrMemberIsOfficer:
		util.ServeError(writer,
			util.Error{Description: `member is officer`},
			http.StatusUnprocessableEntity)

	case nil:
		util.ServeJSON(writer, struct{}{}, http.StatusOK)

	default:
		panic(err)
	}
}

func memberTransformClubs(querier db.MemberClubQuerier) (bool, json.RawMessage) {
	chief := false

	clubs := querier.Query()
	clubsEncoded := bytes.NewBuffer(make(json.RawMessage, 0, 2))
	clubsEncoded.WriteByte('[')
	clubsEncoder := json.NewEncoder(clubsEncoded)

	club, present := <-clubs
	if !present {
		clubsEncoded.WriteByte(']')
		return chief, clubsEncoded.Bytes()
	}

	for {
		if club.Error != nil {
			panic(club.Error)
		}

		if err := clubsEncoder.Encode(club.MemberClub); err != nil {
			for club := range clubs {
				if club.Error != nil {
					log.Print(club.Error)
				}
			}

			panic(err)
		}

		if club.Chief {
			chief = true
		}

		club, present = <-clubs
		clubsEncodedBytes := clubsEncoded.Bytes()

		if !present {
			clubsEncodedBytes[len(clubsEncodedBytes)-1] = ']'

			return chief, clubsEncodedBytes
		}

		clubsEncodedBytes[len(clubsEncodedBytes)-1] = ','
	}
}

func memberTransformPositions(querier db.PositionQuerier) (bool, json.RawMessage) {
	positions := querier.Query()
	positionsEncoded := bytes.NewBuffer(make(json.RawMessage, 0, 2))
	positionsEncoded.WriteByte('[')
	positionsEncoder := json.NewEncoder(positionsEncoded)

	position, present := <-positions
	if !present {
		positionsEncoded.WriteByte(']')
		return false, positionsEncoded.Bytes()
	}

	for {
		if position.Error != nil {
			panic(position.Error)
		}

		if err := positionsEncoder.Encode(position.ID); err != nil {
			for position := range positions {
				if position.Error != nil {
					log.Print(position.Error)
				}
			}

			panic(err)
		}

		position, present = <-positions
		positionsEncodedBytes := positionsEncoded.Bytes()

		if !present {
			positionsEncodedBytes[len(positionsEncodedBytes)-1] = ']'

			return true, positionsEncodedBytes
		}

		positionsEncodedBytes[len(positionsEncodedBytes)-1] = ','
	}
}

func memberGetServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	var authorized claim
	var id string

	if request.URL.Path == `` {
		authorized = authorize(writer, request, shared, scope.User)
		id = authorized.sub
	} else {
		authorized = authorize(writer, request, shared, scope.Member)
		id = request.URL.Path[1:]
	}

	if (authorized == claim{}) {
		return
	}

	switch detail, err := shared.DB.QueryMemberDetail(id); err {
	case nil:
		var chief bool
		var positionsPresent bool
		var clubs json.RawMessage
		var positions json.RawMessage

		func() {
			defer func() {
				if err := detail.End(); err != nil {
					log.Print(err)
				}
			}()

			chief, clubs = memberTransformClubs(detail.Clubs)
			positionsPresent, positions = memberTransformPositions(detail.Positions)
		}()

		mail, err := mail.AddressToUnicode(detail.Mail)
		if err != nil {
			panic(err)
		}

		public := publicDetail{
			Affiliation: detail.Affiliation,
			Clubs:       clubs,
			Entrance:    detail.Entrance,
			Gender:      detail.Gender,
			Mail:        mail,
			Nickname:    detail.Nickname,
			OB:          detail.OB,
			Positions:   positions,
			Realname:    detail.Realname,
		}

		var unmarshalled interface{}
		if id == authorized.sub || chief || positionsPresent || authorized.scope.IsSet(scope.Privacy) {
			unmarshalled = privateDetail{public, detail.Confirmed, detail.Tel}
		} else {
			unmarshalled = public
		}

		util.ServeJSON(writer, unmarshalled, http.StatusOK)

	case db.ErrIncorrectIdentity:
		util.ServeErrorDefault(writer, http.StatusNotFound)

	default:
		panic(err)
	}
}

func memberPatchServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	var authorized claim
	var id string

	if request.URL.Path == `` {
		authorized = authorize(writer, request, shared, scope.User)
		id = authorized.sub
	} else {
		authorized = authorize(writer, request, shared, scope.Management)
		id = request.URL.Path[1:]
	}

	if (authorized == claim{}) {
		return
	}

	confirm := authorized.tmp
	if !confirm {
		if token := request.FormValue(`token`); token != `` {
			authenticated, err := shared.Token.AuthenticateMail(token)
			if err.IsError() {
				util.ServeError(writer,
					util.Error{Description: `bad token`},
					http.StatusBadRequest)

				return
			}

			if id != authenticated.Sub {
				util.ServeError(writer,
					util.Error{Description: `incorrect token`},
					http.StatusBadRequest)

				return
			}

			confirm = true
		}
	}

	entrance := 0
	if entranceString := request.PostFormValue(`entrance`); entranceString != `` {
		var err error

		entrance, err = strconv.Atoi(entranceString)
		if err != nil {
			description := `syntax error in entrance`
			status := http.StatusBadRequest

			if err == strconv.ErrRange {
				description = `entrance out of range`
				status = http.StatusUnprocessableEntity
			}

			util.ServeError(writer,
				util.Error{Description: description},
				status)

			return
		}
	}

	address := request.PostFormValue(`mail`)
	if address != `` {
		var err error

		address, err = mail.AddressToASCII(address)
		if err != nil {
			util.ServeError(writer,
				util.Error{Description: `invalid mail`},
				http.StatusBadRequest)
		}
	}

	ob := false
	switch request.PostFormValue(`ob`) {
	case ``:

	case `1`:
		ob = true

	default:
		util.ServeError(writer,
			util.Error{Description: `invalid ob`},
			http.StatusBadRequest)
	}

	
	password := request.PostFormValue(`new_password`)
	if password != `` && !authorized.tmp {
		// FIXME: should I add rate limiter?
		err := shared.DB.Authenticate(id, request.PostFormValue(`current_password`));
		if err == db.ErrIncorrectIdentity {
			util.ServeError(writer,
				util.Error{Description: `incorrect identity`},
				http.StatusUnprocessableEntity)
			return
		} else if err != nil {
			panic(err)
		}
	}

	affiliation := request.PostFormValue(`affiliation`)
	clubs := request.PostFormValue(`clubs`)
	gender := request.PostFormValue(`gender`)
	nickname := request.PostFormValue(`nickname`)
	realname := request.PostFormValue(`realname`)
	tel := request.PostFormValue(`tel`)

	switch err := shared.DB.UpdateMember(id, confirm, ob, password,
		affiliation, clubs, entrance, gender,
		address, nickname, realname, tel); err {
	case db.ErrIncorrectIdentity:
		// FIXME: should be StatusUnprocessableEntity if member is found
		// but clubs is not.
		util.ServeErrorDefault(writer, http.StatusNotFound)

	case db.ErrInvalid:
		util.ServeErrorDefault(writer, http.StatusUnprocessableEntity)

	case nil:
		if address != `` {
			token, tokenErr := shared.Token.IssueMail(id)
			if tokenErr != nil {
				panic(tokenErr)
			}

			nickname, queryErr := shared.DB.QueryMemberNickname(id)
			if queryErr != nil {
				panic(queryErr)
			}

			mailErr := shared.Mail.SendConfirmation(request.Host,
				netMail.Address{
					Name:    nickname,
					Address: address,
				}, id, token)
			if mailErr != nil {
				util.ServeError(writer,
					util.Error{
						ID:          `mail_failure`,
						Description: `failed to mail`,
					}, http.StatusOK)

				log.Print(mailErr)

				return
			}
		}

		var response interface{}

		if password != `` && authorized.tmp {
			const scope = `user member`

			accessToken, accessErr := shared.Token.IssueAccess(id, scope)
			if accessErr != nil {
				accessToken = ``
			}

			refreshToken, refreshErr := shared.Token.IssueRefresh(id, scope)
			if refreshErr != nil {
				refreshToken = ``
			}

			response = struct {
				AccessToken  string `json:"access_token"`
				RefreshToken string `json:"refresh_token,omitempty"`
				Scope        string `json:"scope"`
			}{accessToken, refreshToken, scope}
		} else {
			response = struct{}{}
		}

		util.ServeJSON(writer, response, http.StatusOK)

	default:
		panic(err)
	}
}

func memberPutServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if len(request.URL.Path) < 2 {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	if (authorize(writer, request, shared, scope.Management) == claim{}) {
		return
	}

	address := request.FormValue(`mail`)
	id := request.URL.Path[1:]
	nickname := request.FormValue(`nickname`)

	asciiAddress, err := mail.AddressToASCII(address)
	if err != nil {
		util.ServeError(writer,
			util.Error{Description: `invalid mail`},
			http.StatusBadRequest)

		return
	}

	switch err := shared.DB.InsertMember(id, asciiAddress, nickname); err {
	case db.ErrBadOmission:
		util.ServeError(writer,
			util.Error{Description: `id, mail, and nickname are required`},
			http.StatusUnprocessableEntity)

	case db.ErrDupEntry:
		// FIXME: should be StatusMethodNotAllowed
		util.ServeError(writer,
			util.Error{Description: `duplicate id`},
			http.StatusUnprocessableEntity)

	case db.ErrInvalid:
		util.ServeErrorDefault(writer, http.StatusUnprocessableEntity)

	case nil:
		token, err := shared.Token.IssueTmpUserAccess(id)
		if err != nil {
			panic(err)
		}

		err = shared.Mail.SendCreation(request.Host,
			netMail.Address{Name: nickname, Address: asciiAddress},
			id, token)
		if err != nil {
			util.ServeMailError(writer)
			return
		}

		util.ServeJSON(writer, struct{}{}, http.StatusCreated)

	default:
		panic(err)
	}
}

func membersGetServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path != `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	if (authorize(writer, request, shared, scope.Member) == claim{}) {
		return
	}

	util.ServeJSON(writer, shared.DB.QueryMembers(), http.StatusOK)
}

func membersMailsGetServeHTTP(writer http.ResponseWriter, request *http.Request, shared shared) {
	if request.URL.Path != `` {
		util.ServeErrorDefault(writer, http.StatusNotFound)
		return
	}

	if (authorize(writer, request, shared, scope.Member) == claim{}) {
		return
	}

	util.ServeJSON(writer, shared.DB.QueryMemberMails(), http.StatusOK)
}
