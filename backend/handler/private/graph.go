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

package private

import (
	"errors"
	"fmt"
	"github.com/kagucho/tsubonesystem3/backend/db"
	"net/url"
	"strconv"
	"strings"
)

type graph struct {
	Properties []property
	URL        string
}

type graphFunc func(db db.DB, base string, routeQuery []string) graph

type property struct {
	Property string
	Content  string
}

func parseQuery(routeQuery []string) map[string]string {
	if len(routeQuery) != 2 {
		return nil
	}

	values := make(map[string]string)

	for _, component := range strings.Split(routeQuery[1], `&`) {
		pair := strings.SplitN(component, `=`, 2)

		if len(pair) == 2 {
			if value, unescapeError := url.QueryUnescape(pair[1]); unescapeError == nil {
				values[pair[0]] = value
			}
		}
	}

	return values
}

var graphFuncs = map[string]graphFunc{
	`club`: graphClub, `clubs`: graphClubs,
	`member`: graphMember, `members`: graphMembers,
	`officer`: graphOfficer, `officers`: graphOfficers,
}

func graphClub(dbInstance db.DB, base string, routeQuery []string) graph {
	var description string
	var title string

	id := parseQuery(routeQuery)[`id`]
	name, err := dbInstance.QueryClubName(id)
	switch err {
	case db.ErrIncorrectIdentity:
		description = `神楽坂一丁目通信局の不明な部門です。`
		title = `TsuboneSystem 不明な部門です`

	case nil:
		description = `神楽坂一丁目通信局の` + name + `の詳細情報です。`
		title = `TsuboneSystem ` + name + `の詳細情報`

	default:
		panic(err)
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

func graphClubs(dbInstance db.DB, base string, routeQuery []string) graph {
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

func graphMember(dbInstance db.DB, base string, routeQuery []string) graph {
	properties := make([]property, 0, 8)

	id := parseQuery(routeQuery)[`id`]
	memberGraph, err := dbInstance.QueryMemberGraph(id)
	switch err {
	case db.ErrIncorrectIdentity:
		properties = append(properties,
			property{
				`og:description`,
				`神楽坂一丁目通信局の不明な局員です。`,
			},
			property{
				`og:title`,
				`TsuboneSystem 不明な局員です`,
			})

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
			properties = append(properties,
				property{`og:profile:gender`, `male`})

		case `女`:
			properties = append(properties,
				property{`og:profile:gender`, `female`})
		}

	default:
		panic(err)
	}

	url := base + `#!member?id=` + url.QueryEscape(id)

	properties = append(properties,
		property{`og:image`, `https://kagucho.net/favicon.ico`},
		property{`og:locale`, `ja_JP`},
		property{`og:type`, `profile`},
		property{`og:url`, url})

	return graph{properties, url}
}

func graphMembers(dbInstance db.DB, base string, routeQuery []string) graph {
	var description string
	var err error
	fragment := `#!members`

	values := parseQuery(routeQuery)

	entrance := values[`entrance`]
	entrancei := 0
	if entrance != `` {
		entrancei, err = strconv.Atoi(values[`entrance`])
		if err == nil && entrancei == 0 {
			err = errors.New(`entrance out of range`)
		}
	}

	if err == nil {
		nickname := values[`nickname`]
		realname := values[`realname`]
		ob := values[`ob`]

		components := make([]string, 0, 3)
		var status db.MemberStatus

		if entrancei != 0 {
			components = append(components,
				fmt.Sprint(`entrance=`, entrancei))
		}

		if nickname != `` {
			components = append(components,
				`nickname=`+url.QueryEscape(nickname))
		}

		if ob == `0` || ob == `1` {
			if ob == `0` {
				status = db.MemberStatusActive
			} else {
				status = db.MemberStatusOB
			}

			components = append(components, `ob=`+ob)
		} else {
			status = db.MemberStatusActive | db.MemberStatusOB
		}

		if realname != `` {
			components = append(components,
				`realname=`+url.QueryEscape(realname))
		}

		if len(components) > 0 {
			fragment += `?` + strings.Join(components, `&`)
		}

		count, err := dbInstance.QueryMembersCount(
			entrancei, nickname, realname, status)
		if err != nil {
			panic(err)
		}

		description = fmt.Sprint(`神楽坂一丁目通信局の局員検索結果です。`,
			count, ` 件`)
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

func graphOfficer(dbInstance db.DB, base string, routeQuery []string) graph {
	var description string
	var title string

	id := parseQuery(routeQuery)[`id`]
	name, err := dbInstance.QueryOfficerName(id)
	switch err {
	case db.ErrIncorrectIdentity:
		description = `神楽坂一丁目通信局の不明な役職です。`
		title = `TsuboneSystem 不明な役職です`

	case nil:
		description = `神楽坂一丁目通信局の` + name + `の詳細情報です。`
		title = `TsuboneSystem ` + name + `の詳細情報`

	default:
		panic(err)
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
