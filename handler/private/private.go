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

// Package private implements private page hosting
package private

import (
  `bytes`
  `database/sql`
  `errors`
  `fmt`
  `github.com/kagucho/tsubonesystem3/db`
  `github.com/kagucho/tsubonesystem3/handler/file`
  `html/template`
  `log`
  `net/http`
  `net/url`
  `path`
  `runtime/debug`
  `strconv`
  `strings`
)

type Private struct {
  file string
  graph *template.Template
  db db.DB
  fileError file.FileError
}

type graph struct {
  Properties []property
  URL string
}

type graphFunc func(db db.DB, base string, routeQuery []string) graph

type property struct {
  Property string
  Content string
}

var graphFuncs map[string]graphFunc = map[string]graphFunc {
  `club`: graphClub, `clubs`: graphClubs,
  `member`: graphMember, `members`: graphMembers,
  `officer`: graphOfficer, `officers`: graphOfficers,
}

func New(share string, db db.DB, fileError file.FileError) (Private, error) {
  graph, parseError := template.ParseFiles(path.Join(share, `graph`))
  if parseError != nil {
    return Private{}, parseError
  }

  return Private{ path.Join(share, `public/private`), graph, db, fileError },
         nil
}

func parseQuery(routeQuery []string) map[string]string {
  if len(routeQuery) != 2 {
    return nil
  }

  values := make(map[string]string)

  for _, component := range strings.Split(routeQuery[1], `&`) {
    pair := strings.SplitN(component, `=`, 2)

    if len(pair) == 2 {
      if value, unescapeError := url.QueryUnescape(pair[1]);
         unescapeError == nil {
        values[pair[0]] = value
      }
    }
  }

  return values
}

func graphClub(db db.DB, base string, routeQuery []string) graph {
  var description string
  var title string

  id := parseQuery(routeQuery)[`id`]
  name, queryError := db.QueryClubName(id)
  switch queryError {
  case sql.ErrNoRows:
    description = `神楽坂一丁目通信局の不明な部門です。`
    title = `TsuboneSystem 不明な部門です`

  case nil:
    description = `神楽坂一丁目通信局の` + name + `の詳細情報です。`
    title = `TsuboneSystem ` + name + `の詳細情報`

  default:
    panic(queryError)
  }

  url := base + `#!club?id=` + url.QueryEscape(id)

  return graph{
    []property{
      {`og:description`, description},
      {`og:image`, `https://kagucho.net/favicon.ico`},
      {`og:locale`, `ja_JP`},
      {`og:title`, title},
      {`og:type`, `website`},
      {`og:url`, url},
    }, url,
  }
}

func graphClubs(db db.DB, base string, routeQuery []string) graph {
  url := base + "#!clubs"

  return graph{
    []property{
      {`og:description`, `神楽坂一丁目通信局の部門一覧です。`},
      {`og:image`, `https://kagucho.net/favicon.ico`},
      {`og:locale`, `ja_JP`},
      {`og:title`, `TsuboneSystem 部門一覧`},
      {`og:type`, `website`},
      {`og:url`, url},
    }, url,
  }
}

func graphMember(db db.DB, base string, routeQuery[]string) graph {
  properties := make([]property, 0, 8)

  id := parseQuery(routeQuery)[`id`]
  memberGraph, queryError := db.QueryMemberGraph(id)
  switch queryError {
  case sql.ErrNoRows:
    properties = append(properties,
      property{`og:description`, `神楽坂一丁目通信局の不明な局員です。`},
      property{`og:title`, `TsuboneSystem 不明な局員です`})

  case nil:
    properties = append(properties,
      property{ 
        `og:description`,
        `神楽坂一丁目通信局の局員、` + memberGraph.Nickname + `の詳細情報です。`,
      }, property{
        `og:title`, `TsuboneSystem ` + memberGraph.Nickname + `の詳細情報`,
      }, property{
        `og:profile:username`, id,
      })

    switch memberGraph.Gender {
    case `男`:
      properties = append(properties, property{`og:profile:gender`, `male`})

    case `女`:
      properties = append(properties, property{`og:profile:gender`, `female`})
    }

  default:
    panic(queryError)
  }

  url := base + `#!member?id=` + url.QueryEscape(id)

  properties = append(properties,
    property{`og:image`, `https://kagucho.net/favicon.ico`},
    property{`og:locale`, `ja_JP`},
    property{`og:type`, `profile`},
    property{`og:url`, url})

  return graph{ properties, url }
}

func graphMembers(dbInstance db.DB, base string, routeQuery []string) graph {
  var description string
  fragment := `#!members`

  values := parseQuery(routeQuery)

  entrance := values[`entrance`]
  entrancei := 0
  var entranceError error
  if entrance != `` {
    entrancei, entranceError = strconv.Atoi(values[`entrance`])
    if entranceError == nil && !db.ValidateMemberEntrance(entrancei) {
      entranceError = errors.New(`entrance out of range`)
    }
  }

  if entranceError == nil {
    nickname := values[`nickname`]
    realname := values[`realname`]
    ob := values[`ob`]

    components := make([]string, 0, 3)
    var status db.MemberStatus

    if entrancei != 0 {
      components = append(components, fmt.Sprint(`entrance=`, entrancei))
    }

    if nickname != `` {
      components = append(components, `nickname=` + url.QueryEscape(nickname))
    }

    if ob == `0` || ob == `1` {
      if ob == `0` {
        status = db.MemberStatusActive
      } else {
        status = db.MemberStatusOB
      }

      components = append(components, `ob=` + ob)
    } else {
      status = db.MemberStatusActive | db.MemberStatusOB
    }

    if realname != `` {
      components = append(components, `realname=` + url.QueryEscape(realname))
    }

    if len(components) > 0 {
      fragment += `?` + strings.Join(components, `&`)
    }

    count, queryError := dbInstance.QueryMembersCount(
      entrancei, nickname, realname, status)
    if queryError != nil {
      panic(queryError)
    }

    description = fmt.Sprint(`神楽坂一丁目通信局の局員検索結果です。`, count, ` 件`)
  } else {
    description = `神楽坂一丁目通信局の局員検索結果です。不正な入学年度の指定がありました。`
  }

  url := base + fragment
  return graph{
    []property{
      {`og:description`, description},
      {`og:image`, `https://kagucho.net/favicon.ico`},
      {`og:locale`, `ja_JP`},
      {`og:title`, `TsuboneSystem 局員検索`},
      {`og:type`, `website`},
      {`og:url`, url},
    }, url,
  }
}

func graphOfficer(db db.DB, base string, routeQuery []string) graph {
  var description string
  var title string

  id := parseQuery(routeQuery)[`id`]
  name, queryError := db.QueryOfficerName(id)
  switch queryError {
  case sql.ErrNoRows:
    description = `神楽坂一丁目通信局の不明な役職です。`
    title = `TsuboneSystem 不明な役職です`

  case nil:
    description = `神楽坂一丁目通信局の` + name + `の詳細情報です。`
    title = `TsuboneSystem ` + name + `の詳細情報`

  default:
    panic(queryError)
  }

  url := base + `#!officer?id=` + url.QueryEscape(id)

  return graph{
    []property{
      {`og:description`, description},
      {`og:image`, `https://kagucho.net/favicon.ico`},
      {`og:locale`, `ja_JP`},
      {`og:title`, title},
      {`og:type`, `website`},
      {`og:url`, url},
    }, url,
  }
}

func graphOfficers(db db.DB, base string, routeQuery []string) graph {
  url := base + `#!officers`

  return graph{
    []property{
      {`og:description`, `神楽坂一丁目通信局の役員一覧です。`},
      {`og:image`, `https://kagucho.net/favicon.ico`},
      {`og:locale`, `ja_JP`},
      {`og:title`, `TsuboneSystem 役員一覧`},
      {`og:type`, `website`},
      {`og:url`, url},
    }, url,
  }
}

func graphDefault(db db.DB, base string, routeQuery []string) graph {
  return graph{
    []property{
      {`og:description`, `神楽坂一丁目通信局内で利用しているWebサービスです。`},
      {`og:image`, `https://kagucho.net/favicon.ico`},
      {`og:locale`, `ja_JP`},
      {`og:title`, `TsuboneSystem`},
      {`og:type`, `website`},
      {`og:url`, base},
    }, base,
  }
}

func (private Private) ServeHTTP(writer http.ResponseWriter,
                                 request *http.Request) {
  serve := func() {
    private.fileError.ServeError(writer, http.StatusInternalServerError)
  }

  defer func() {
    if recovered := recover(); recovered != nil {
      log.Println(recovered)
      debug.PrintStack()
    }

    serve()
  }()

  switch request.URL.Path {
  case `/private`:
    request.ParseForm()
    writer.Header().Set(`Content-Language`, `ja`)

    escapedFragment := request.Form[`_escaped_fragment_`]
    if len(escapedFragment) == 0 {
      if request.URL.RawQuery == `` {
        serve = func() {
          http.ServeFile(writer, request, private.file)
        }
      } else {
        request.URL.RawQuery = ``
        serve = func() {
          private.fileError.ServeMovedPermanently(writer, request.URL.String())
        }
      }
    } else {
      routeQuery := strings.SplitN(escapedFragment[0], `?`, 2)

      request.URL.RawQuery = ``
      base := request.URL.String()

      graphFunc := graphFuncs[routeQuery[0]]
      if graphFunc == nil {
        graphFunc = graphDefault
      }

      var buffer bytes.Buffer
      if executeError := private.graph.Execute(
           &buffer, graphFunc(private.db, base, routeQuery))
         executeError != nil {
        panic(executeError)
      }

      header := writer.Header()
      header.Set(`Content-Length`, strconv.Itoa(buffer.Len()))

      serve = func() {
        buffer.WriteTo(writer)
      }
    }

  case `/private/`:
    request.URL.Path = `/private`
    request.URL.RawQuery = ``
    serve = func() {
      private.fileError.ServeMovedPermanently(writer, request.URL.String())
    }

  default:
    serve = func() {
      private.fileError.ServeError(writer, http.StatusNotFound)
    }
  }
}
