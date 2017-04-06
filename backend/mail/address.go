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
	"golang.org/x/net/idna"
	"strings"
)

func convertAddress(convert func(string) (string, error), address string) (string, error) {
	if address == `` {
		return ``, nil
	}

	var err error
	splitted := strings.SplitN(address, `@`, 2)

	splitted[1], err = convert(splitted[1])
	if err != nil {
		return ``, err
	}

	return strings.Join(splitted, `@`), nil
}

/*
AddressToUnicode returns the Unicode representation of the given email
address.
*/
func AddressToUnicode(address string) (string, error) {
	return convertAddress(idna.ToUnicode, address)
}

// AddressToASCII returns the ASCII representation of the given email address.
func AddressToASCII(address string) (string, error) {
	return convertAddress(idna.ToASCII, address)
}

/*
ValidateAddress returns a boolean telling whether the given email address is
valid or not.
*/
func ValidateAddress(address string) bool {
	// https://html.spec.whatwg.org/multipage/forms.html#valid-e-mail-address
	index := 0
	for {
		if index >= len(address) {
			return false
		}

		if address[index] == '@' {
			break
		}

		if !isAtext(address[index]) && address[index] != '.' {
			return false
		}

		index++
	}

	for {
		index++
		if index >= len(address) {
			return false
		}

		address = address[index:]

		if !isLetDig(address[0]) {
			return false
		}

		for index = 1; ; index++ {
			if index >= len(address) {
				return address[index-1] != '-'
			}

			if address[index] == '.' {
				break
			}

			if !isLetDigHyp(address[index]) {
				return false
			}
		}

		if address[index-1] == '-' || index > 63 {
			return false
		}
	}
}
