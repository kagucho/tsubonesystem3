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
	"net/mail"
	"net/url"
)

// SendCreation sends an email ton continue the creation of a member.
func (context Mail) SendCreation(host string, address mail.Address, id, token string) error {
	var data struct {
		Base     string
		Register string
	}

	constructing := url.URL{Scheme: `https`, Host: host}
	data.Base = constructing.String()

	constructing.Path = `/private`
	constructing.Fragment = `!member?id=` + id + `&fill=` + token
	data.Register = constructing.String()

	return context.send(host, []string{`-t`}, ``, []mail.Address{address},
		`TsuboneSystem 登録手続き`, templateCreation, data)
}
