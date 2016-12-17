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

package member

import (
  `database/sql`
  `fmt`
  `github.com/kagucho/tsubonesystem3/db`
  `github.com/kagucho/tsubonesystem3/scope`
  `github.com/kagucho/tsubonesystem3/handler/apiv0/common`
  `github.com/kagucho/tsubonesystem3/handler/apiv0/token/authorizer`
  `net/http`
)

type publicDetail struct {
  Affiliation string `json:"affiliation"`
  Clubs []db.MemberClub `json:"clubs"`
  Entrance uint16 `json:"entrance"`
  Gender string `json:"gender"`
  Mail string `json:"mail"`
  Nickname string `json:"nickname"`
  OB bool `json:"ob"`
  Positions []db.MemberPosition `json:"positions"`
  Realname string `json:"realname"`
}

type privateDetail struct {
  publicDetail
  Tel string `json:"tel"`
}

// DetailServeHTTP serves the detail of the member identified with the given ID
// via HTTP.
func DetailServeHTTP(writer http.ResponseWriter, request *http.Request,
                     dbInstance db.DB, claim authorizer.Claim) {
  serve := func() func() {
    defer common.Recover(writer)

    id := request.FormValue(`id`)
    detail, queryError := dbInstance.QueryMember(id)

    switch queryError {
    case nil:
      public := publicDetail{
        Affiliation: detail.Affiliation,
        Clubs: make([]db.MemberClub, 0),
        Entrance: detail.Entrance,
        Gender: detail.Gender,
        Mail: detail.Mail,
        Nickname: detail.Nickname,
        OB: detail.OB,
        Positions: make([]db.MemberPosition, 0),
        Realname: detail.Realname,
      }

      var chanError error
      chief := false
      for detail.Clubs != nil || detail.ClubsErrors != nil ||
          detail.Positions != nil || detail.PositionsErrors != nil {
        select {
        case result, present := <-detail.Clubs:
          if present {
            public.Clubs = append(public.Clubs, result)
            if result.Chief {
              chief = true
            }
          } else {
            detail.Clubs = nil
          }

        case result, present := <-detail.ClubsErrors:
          if present {
            if chanError == nil {
              chanError = result
            }
          } else {
            detail.ClubsErrors = nil
          }

        case result, present := <-detail.Positions:
          if present {
            public.Positions = append(public.Positions, result)
          } else {
            detail.Positions = nil
          }

        case result, present := <-detail.PositionsErrors:
          if present {
            if chanError == nil {
              chanError = result
            }
          } else {
            detail.PositionsErrors = nil
          }
        }
      }

      if chanError != nil {
        panic(chanError)
      }

      var unmarshalled interface{}
      if chief || len(public.Positions) > 0 ||
         claim.Sub == id || claim.Scope.IsSet(scope.Privacy) {
        unmarshalled = privateDetail{public, detail.Tel}
      } else {
        unmarshalled = public
      }

      return func() {
        common.ServeJSON(writer, unmarshalled, http.StatusOK)
      }

    case sql.ErrNoRows:
      return func() {
        common.ServeError(writer, `invalid_id`,
                          fmt.Sprintf(`unknown ID: %q`, id),
                          ``, http.StatusBadRequest)
      }

    default:
      panic(queryError)
    }
  }()

  if serve != nil {
    serve()
  }
}
