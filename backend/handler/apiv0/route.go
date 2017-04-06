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

import "sort"

type route struct {
	prefix string
	handler handler
}

type routeSlice []route

func (slice routeSlice) Len() int {
	return len(slice)
}

func (slice routeSlice) Less(i, j int) bool {
	return slice[i].prefix < slice[j].prefix
}

func (slice routeSlice) Swap(i, j int) {
	temporary := slice[i]
	slice[i] = slice[j]
	slice[j] = temporary
}

func (slice routeSlice) search(path string) int {
	return sort.Search(len(slice), func(index int) bool {
		return slice[index].prefix >= path
	})
}

func (apiv0 APIv0) newRoutes() routeSlice {
	routes := routeSlice{
		{
			`/club`,
			methodMux{
				map[string]handlerFunc{
					`DELETE`: clubDeleteServeHTTP,
					`GET`:    clubGetServeHTTP,
					`HEAD`:   clubGetServeHTTP,
					`PATCH`:   clubPatchServeHTTP,
					`PUT`:    clubPutServeHTTP,
				},
				[]field{
					{`Accept-Patch`, `application/x-www-form-urlencoded`},
					{`Accept-Ranges`, `none`},
				},
			},
		},
		{
			`/clubs`,
			methodMux{
				map[string]handlerFunc{
					`GET`:  clubsGetServeHTTP,
					`HEAD`: clubsGetServeHTTP,
				},
				[]field{{`Accept-Ranges`, `none`}},
			},
		},
		{
			`/mail`,
			methodMux{
				map[string]handlerFunc{
					`DELETE`: mailDeleteServeHTTP,
					`GET`:    mailGetServeHTTP,
					`HEAD`:   mailGetServeHTTP,
					`PATCH`:  mailPatchServeHTTP,
					`PUT`:    mailPutServeHTTP,
				},
				[]field{
					{`Accept-Patch`, `application/x-www-form-urlencoded`},
					{`Accept-Ranges`, `none`},
				},
			},
		},
		{
			`/mails`,
			methodMux{
				map[string]handlerFunc{
					`GET`:  mailsGetServeHTTP,
					`HEAD`: mailsGetServeHTTP,
				},
				[]field{
					{`Accept-Patch`, `application/x-www-form-urlencoded`},
					{`Accept-Ranges`, `none`},
				},
			},
		},
		{
			`/member`,
			methodMux{
				map[string]handlerFunc{
					`DELETE`: memberDeleteServeHTTP,
					`GET`:    memberGetServeHTTP,
					`HEAD`:   memberGetServeHTTP,
					`PATCH`:  memberPatchServeHTTP,
					`PUT`:    memberPutServeHTTP,
				},
				[]field{
					{`Accept-Patch`, `application/x-www-form-urlencoded`},
					{`Accept-Ranges`, `none`},
				},
			},
		},
		{
			`/members`,
			methodMux{
				map[string]handlerFunc{
					`GET`:  membersGetServeHTTP,
					`HEAD`: membersGetServeHTTP,
				},
				[]field{{`Accept-Ranges`, `none`}},
			},
		},
		{
			`/members/mails`,
			methodMux{
				map[string]handlerFunc{
					`GET`:  membersMailsGetServeHTTP,
					`HEAD`: membersMailsGetServeHTTP,
				},
				[]field{{`Accept-Ranges`, `none`}},
			},
		},
		{
			`/officer`,
			methodMux{
				map[string]handlerFunc{
					`DELETE`: officerDeleteServeHTTP,
					`GET`:    officerGetServeHTTP,
					`HEAD`:   officerGetServeHTTP,
					`PATCH`:  officerPatchServeHTTP,
					`PUT`:    officerPutServeHTTP,
				},
				[]field{
					{`Accept-Patch`, `application/x-www-form-urlencoded`},
					{`Accept-Ranges`, `none`},
				},
			},
		},
		{
			`/officers`,
			methodMux{
				map[string]handlerFunc{
					`GET`:  officersGetServeHTTP,
					`HEAD`: officersGetServeHTTP,
				},
				[]field{{`Accept-Ranges`, `none`}},
			},
		},
		{
			`/officers/names`,
			methodMux{
				map[string]handlerFunc{
					`GET`:  officersNamesGetServeHTTP,
					`HEAD`: officersNamesGetServeHTTP,
				},
				[]field{{`Accept-Ranges`, `none`}},
			},
		},
		{
			`/parties`,
			methodMux{
				map[string]handlerFunc{
					`GET`:  partiesGetServeHTTP,
					`HEAD`: partiesGetServeHTTP,
				},
				[]field{{`Accept-Ranges`, `none`}},
			},
		},
		{
			`/party`,
			methodMux{
				map[string]handlerFunc{
					`DELETE`: partyDeleteServeHTTP,
					`GET`:    partyGetServeHTTP,
					`HEAD`:   partyGetServeHTTP,
					`PATCH`:  partyPatchServeHTTP,
					`PUT`:    partyPutServeHTTP,
				},
				[]field{
					{`Accept-Patch`, `application/x-www-form-urlencoded`},
					{`Accept-Ranges`, `none`},
				},
			},
		},
		{
			`/token`,
			methodMux{
				map[string]handlerFunc{
					`POST`: apiv0.tokenServer.serveHTTP,
				},
				[]field{{`Accept-Ranges`, `none`}},
			},
		},
	}

	sort.Sort(routes)

	return routes
}
