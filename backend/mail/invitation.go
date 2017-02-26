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

package mail

import (
	"net/url"
	"time"
)

func (context Mail) Invite(host string, recipients []string, title string, start, end time.Time, place, inviteds string, due time.Time, details string) error {
	var data struct {
		Title    string
		Datetime string
		Place    string
		Inviteds string
		Due      string
		Details  string
		URL      string
		Base     string
	}

	constructing := url.URL{Scheme: `https`, Host: host}
	data.Base = constructing.String()

	constructing.Path = `/private`
	constructing.Fragment = `!party?title=` + url.QueryEscape(title)
	data.URL = constructing.String()

	data.Title = title
	data.Datetime = start.Format(time.Stamp) + ` - ` + end.Format(time.Stamp)
	data.Place = place
	data.Inviteds = inviteds
	data.Due = due.Format(time.Stamp)
	data.Details = details

	return context.send(host, recipients, inviteds, nil, `TsuboneSystem パーティー招待: ` + title, templateInvitation, data)
}
