package mail

import (
	"golang.org/x/net/idna"
	"strings"
)

func convertAddress(convert func(string) (string, error), address string) (string, error) {
	if address == `` {
		return ``, nil
	}

	var idnaError error
	splitted := strings.SplitN(address, `@`, 2)

	splitted[1], idnaError = convert(splitted[1])
	if idnaError != nil {
		return ``, idnaError
	}

	return strings.Join(splitted, `@`), nil
}

func AddressToUnicode(address string) (string, error) {
	return convertAddress(idna.ToUnicode, address)
}

func AddressToASCII(address string) (string, error) {
	return convertAddress(idna.ToASCII, address)
}

func ValidateAddressHTML(address string) bool {
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
				return address[index - 1] != '-'
			}

			if address[index] == '.' {
				break
			}

			if !isLetDigHyp(address[index]) {
				return false
			}
		}

		if address[index - 1] == '-' || index > 63 {
			return false
		}
	}
}
