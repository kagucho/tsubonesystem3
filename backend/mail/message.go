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

import "net/url"

func (context Mail) Message(host string, recipients []string, id, nickname, to, subject, body string) error {
	var data struct {
		Base          string
		Body          string
		From          string
		FromReference string
	}

	constructing := url.URL{Scheme: `https`, Host: host}
	data.Base = constructing.String()

	constructing.Path = `/private`
	constructing.Fragment = `!member?id=` + id
	data.FromReference = constructing.String()

	data.Body = body
	data.From = nickname

	return context.send(host, recipients, to, nil, subject, templateMessage, data)
}
