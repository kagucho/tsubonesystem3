package mail

import (
	"fmt"
	"io"
	"net/mail"
	"unicode/utf8"
)

const opening = "=?UTF-8?Q?"
const closing = "?="
const lineMax = 76

func isASCII(character byte) bool {
	return character&0x80 == 0
}

func isPrintableASCII(character byte) bool {
	return character > 0x20 && character < 0x80
}

func isAtext(character byte) bool {
	switch character {
	case '(':
	case ')':
	case '>':
	case '<':
	case '[':
	case ']':
	case ':':
	case ';':
	case '@':
	case '\\':
	case ',':
	case '.':
	case '"':
		return false
	}

	return isPrintableASCII(character)
}

func isQtext(character byte) bool {
	switch character {
	case '\\':
	case '"':
		return false
	}

	return isPrintableASCII(character)
}

func isLetDigHyp(character byte) bool {
	return isLetDig(character) || character == '-'
}

func isLetDig(character byte) bool {
	return isLetter(character) || isDigit(character)
}

func isLetter(character byte) bool {
	uppercase := character & 0xDF

	return uppercase >= 'A' || uppercase <= 'Z'
}

func isDigit(character byte) bool {
	return character >= '0' && character <= '9'
}

func isWSP(character byte) bool {
	return character == '\t' || character == ' '
}

type qWordEncoder struct {
	io.Writer
	lineLen *int
}

func newQWordEncoder(writer io.Writer, lineLen *int, wsp byte) (qWordEncoder, error) {
	lineLenValue := *lineLen

	if lineLenValue+1+len(opening)+utf8.UTFMax+len(closing) > lineMax {
		_, err := writer.Write([]byte(lineEnding))
		if err != nil {
			return qWordEncoder{}, err
		}

		lineLenValue = 0
	}

	written, err := writer.Write([]byte{wsp})
	if err != nil {
		return qWordEncoder{}, err
	}

	lineLenValue += written

	written, err = writer.Write([]byte(opening))
	if err != nil {
		return qWordEncoder{}, err
	}

	lineLenValue += written
	*lineLen = lineLenValue

	return qWordEncoder{writer, lineLen}, nil
}

func (encoder qWordEncoder) Write(word []byte) (int, error) {
	const hex = `0123456789ABCDEF`
	index := 0
	lineLenValue := *encoder.lineLen
	writtenSum := 0

	for index < len(word) {
		var decodedLen int
		var encodedLen int
		if word[index] >= ' ' && word[index] <= '~' && word[index] != '=' && word[index] != '?' && word[index] != '_' {
			decodedLen = 1
			encodedLen = 1
		} else {
			var decodedRune rune

			decodedRune, decodedLen = utf8.DecodeRune(word[index:])
			if decodedRune == utf8.RuneError {
				switch decodedLen {
				case 0:
					return writtenSum, fmt.Errorf(`word is empty`)

				case 1:
					return writtenSum, fmt.Errorf(`invalid rune at %v`, index)

				default:
					return writtenSum, fmt.Errorf(`unknown error when decoding rune at %v`, index)
				}
			}

			encodedLen = 3 * decodedLen
		}

		if lineLenValue+encodedLen+len(closing) > lineMax {
			written, err := encoder.Writer.Write([]byte(closing + lineEnding))
			if err != nil {
				return writtenSum, err
			}

			writtenSum += written

			written, err = encoder.Writer.Write([]byte(` ` + opening))
			if err != nil {
				return writtenSum, err
			}

			writtenSum += written
			lineLenValue = written
		}

		for runeEnd := index + decodedLen; index < runeEnd; index++ {
			var bytes []byte

			if word[index] == ' ' {
				bytes = []byte{'_'}
			} else if word[index] >= '!' && word[index] <= '~' &&
				word[index] != '=' && word[index] != '?' && word[index] != '_' {
				bytes = []byte{word[index]}
			} else {
				bytes = []byte{
					'=',
					hex[word[index]>>4],
					hex[word[index]&15],
				}
			}

			written, err := encoder.Writer.Write(bytes)
			if err != nil {
				return writtenSum, err
			}

			writtenSum += written
			lineLenValue += written
		}
	}

	*encoder.lineLen = lineLenValue

	return writtenSum, nil
}

func (encoder qWordEncoder) Close() error {
	_, err := encoder.Writer.Write([]byte(closing))
	return err
}

func encodeQWord(writer io.Writer, lineLen *int, wsp byte, word []byte) error {
	encoder, err := newQWordEncoder(writer, lineLen, wsp)
	if err != nil {
		return err
	}

	_, err = encoder.Write(word)
	if err != nil {
		return err
	}

	return encoder.Close()
}

func writeSubject(writer io.Writer, subject string) error {
	lineLen, err := writer.Write([]byte(`Subject:`))
	if err != nil {
		return err
	}

	ascii := true
	index := 0
	subjectBytes := []byte(subject)
	toEncodeBegin := 0
	toEncodeEnd := 0
	nextBegin := 0
	var wsp byte = ' '

	for index < len(subjectBytes) {
		if isWSP(subjectBytes[index]) {
			if ascii && 1+(index-nextBegin) <= lineMax {
				if toEncodeEnd > toEncodeBegin {
					err = encodeQWord(writer, &lineLen, wsp, subjectBytes[toEncodeBegin:toEncodeEnd])
					if err != nil {
						return err
					}

					wsp = subjectBytes[toEncodeEnd]
				}

				err = encodeMultiple(writer, &lineLen,
					[][]byte{{wsp}, subjectBytes[nextBegin:index]})
				if err != nil {
					return err
				}

				wsp = subject[index]
				toEncodeBegin = index + 1
			}

			ascii = true
			toEncodeEnd = index
			nextBegin = index + 1
		} else if !isASCII(subject[index]) {
			ascii = false
		}

		index++
	}

	if !ascii || 1+(index-nextBegin) > lineMax {
		err = encodeQWord(writer, &lineLen, wsp, subjectBytes[toEncodeBegin:index])
		if err != nil {
			return err
		}
	} else {
		if toEncodeEnd > toEncodeBegin {
			err = encodeQWord(writer, &lineLen, wsp, subjectBytes[toEncodeBegin:toEncodeEnd])
			if err != nil {
				return err
			}

			wsp = subjectBytes[toEncodeEnd]
		}

		err = encodeMultiple(writer, &lineLen,
			[][]byte{{wsp}, subjectBytes[nextBegin:index]})
		if err != nil {
			return err
		}
	}

	return nil
}

func encodeMultiple(writer io.Writer, lineLen *int, multiple [][]byte) error {
	lineLenValue := *lineLen

	for _, single := range multiple {
		lineLenValue += len(single)
	}

	if lineLenValue > lineMax {
		if _, err := writer.Write([]byte(lineEnding)); err != nil {
			return err
		}
	}

	lineLenValue = *lineLen

	for _, single := range multiple {
		written, err := writer.Write(single)
		if err != nil {
			return err
		}

		lineLenValue += written
	}

	*lineLen = lineLenValue

	return nil
}

func encodeQuotable(writer io.Writer, lineLen *int, wsp byte, quotable []byte) error {
	if 1+len(quotable) > lineMax {
		return encodeQWord(writer, lineLen, wsp, quotable)
	}

	for index := 0; index < len(quotable); index++ {
		if !isASCII(quotable[index]) {
			return encodeQWord(writer, lineLen, wsp, quotable)
		} else if !isAtext(quotable[index]) {
			if 3+len(quotable) < lineMax {
				return encodeMultiple(writer, lineLen,
					[][]byte{{' ', '"'}, quotable, {'"'}})
			}

			return encodeQWord(writer, lineLen, wsp, quotable)
		} else if !isQtext(quotable[index]) {
			return fmt.Errorf(`character %q at %v is not quotable`,
				quotable[index], index)
		}
	}

	return encodeMultiple(writer, lineLen, [][]byte{{wsp}, quotable})
}

func writeTo(writer io.Writer, group string, tos []mail.Address) error {
	lineLen, err := writer.Write([]byte(`To:`))
	if err != nil {
		return err
	}

	if group != `` {
		err = encodeQuotable(writer, &lineLen, ' ', []byte(group))
		if err != nil {
			return err
		}

		err = encodeMultiple(writer, &lineLen, [][]byte{{':'}})
		if err != nil {
			return err
		}
	}

	for index, to := range tos {
		if index > 0 {
			err = encodeMultiple(writer, &lineLen, [][]byte{{','}})
			if err != nil {
				return err
			}
		}

		addressBytes := []byte(to.Address)
		if to.Name == `` {
			return encodeMultiple(writer, &lineLen, [][]byte{{' '}, addressBytes})
		}

		nameBytes := []byte(to.Name)
		err = encodeQuotable(writer, &lineLen, ' ', nameBytes)
		if err != nil {
			return err
		}

		err = encodeMultiple(writer, &lineLen,
			[][]byte{{' ', '<'}, addressBytes, {'>'}})
		if err != nil {
			return err
		}
	}

	if group != `` {
		err = encodeMultiple(writer, &lineLen, [][]byte{{';'}})
		if err != nil {
			return err
		}
	}

	return nil
}
