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
	"strings"
)

// SendConfirmation sends an email to confirm the email address.
func (context Mail) SendConfirmation(host string, address mail.Address, id string, token string) error {
	var data struct {
		Base         string
		Confirmation string
	}

	constructing := url.URL{Scheme: `https`, Host: host}
	data.Base = constructing.String()

	constructing.Path = `/private`
	constructing.Fragment = strings.Join([]string{`!member?id=`, id, `&confirm=`, token}, ``)
	data.Confirmation = constructing.String()

	return context.send(host, []string{`-t`}, ``, []mail.Address{address},
		`TsuboneSystem メール確認`, templateConfirmation, data)
}
