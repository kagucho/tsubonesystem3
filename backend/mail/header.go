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
	return character & 0x80 == 0
}

func isPrintableASCII(character byte) bool {
	return character > 0x20 && character < 0x80
}

func isAtext(character byte) bool {
	switch (character) {
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
	switch (character) {
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

	if lineLenValue + 1 + len(opening) + utf8.UTFMax + len(closing) > lineMax {
		_, writeError := writer.Write([]byte(lineEnding))
		if writeError != nil {
			return qWordEncoder{}, writeError
		}

		lineLenValue = 0
	}

	written, writeError := writer.Write([]byte{wsp})
	if writeError != nil {
		return qWordEncoder{}, writeError
	}

	lineLenValue += written

	written, writeError = writer.Write([]byte(opening))
	if writeError != nil {
		return qWordEncoder{}, writeError
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
				switch (decodedLen) {
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

		if lineLenValue + encodedLen + len(closing) > lineMax {
			written, writeError := encoder.Writer.Write([]byte(closing + lineEnding))
			if writeError != nil {
				return writtenSum, writeError
			}

			writtenSum += written

			written, writeError = encoder.Writer.Write([]byte(` ` + opening))
			if writeError != nil {
				return writtenSum, writeError
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
					hex[word[index] >> 4],
					hex[word[index] & 15],
				}
			}

			written, writeError := encoder.Writer.Write(bytes)
			if writeError != nil {
				return writtenSum, writeError
			}

			writtenSum += written
			lineLenValue += written
		}
	}

	*encoder.lineLen = lineLenValue

	return writtenSum, nil
}

func (encoder qWordEncoder) Close() error {
	_, writeError := encoder.Writer.Write([]byte(closing))

	return writeError
}

func encodeQWord(writer io.Writer, lineLen *int, wsp byte, word []byte) error {
	encoder, encoderError := newQWordEncoder(writer, lineLen, wsp)
	if encoderError != nil {
		return encoderError
	}

	_, encoderError = encoder.Write(word)
	if encoderError != nil {
		return encoderError
	}

	return encoder.Close()
}

func writeSubject(writer io.Writer, subject string) error {
	lineLen, writeError := writer.Write([]byte(`Subject:`))
	if writeError != nil {
		return writeError
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
			if ascii && 1 + (index - nextBegin) <= lineMax {
				if toEncodeEnd > toEncodeBegin {
					writeError = encodeQWord(writer, &lineLen, wsp, subjectBytes[toEncodeBegin:toEncodeEnd])
					if writeError != nil {
						return writeError
					}

					wsp = subjectBytes[toEncodeEnd]
				}

				writeError = encodeMultiple(writer, &lineLen,
					[][]byte{{wsp}, subjectBytes[nextBegin:index]})
				if writeError != nil {
					return writeError
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

	if !ascii || 1 + (index - nextBegin) > lineMax {
		writeError = encodeQWord(writer, &lineLen, wsp, subjectBytes[toEncodeBegin:index])
		if writeError != nil {
			return writeError
		}
	} else {
		if toEncodeEnd > toEncodeBegin {
			writeError = encodeQWord(writer, &lineLen, wsp, subjectBytes[toEncodeBegin:toEncodeEnd])
			if writeError != nil {
				return writeError
			}

			wsp = subjectBytes[toEncodeEnd]
		}

		writeError = encodeMultiple(writer, &lineLen,
			[][]byte{{wsp}, subjectBytes[nextBegin:index]})
		if writeError != nil {
			return writeError
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
		if _, writeError := writer.Write([]byte(lineEnding)); writeError != nil {
			return writeError
		}
	}

	lineLenValue = *lineLen

	for _, single := range multiple {
		written, writeError := writer.Write(single)
		if writeError != nil {
			return writeError
		}

		lineLenValue += written
	}

	*lineLen = lineLenValue

	return nil
}

func encodeQuotable(writer io.Writer, lineLen *int, wsp byte, quotable []byte) error {
	if 1 + len(quotable) > lineMax {
		return encodeQWord(writer, lineLen, wsp, quotable)
	}

	for index := 0; index < len(quotable); index++ {
		if !isASCII(quotable[index]) {
			return encodeQWord(writer, lineLen, wsp, quotable)
		} else if !isAtext(quotable[index]) {
			if 3 + len(quotable) < lineMax {
				return encodeMultiple(writer, lineLen,
					[][]byte{{' ', '"'}, quotable, {'"'}})
			} else {
				return encodeQWord(writer, lineLen, wsp, quotable)
			}
		} else if !isQtext(quotable[index]) {
			return fmt.Errorf(`character %q at %v is not quotable`,
				quotable[index], index)
		}
	}

	return encodeMultiple(writer, lineLen, [][]byte{{wsp}, quotable})
}

func writeTo(writer io.Writer, group string, tos []mail.Address) error {
	lineLen, writeError := writer.Write([]byte(`To:`))
	if writeError != nil {
		return writeError
	}

	if group != `` {
		writeError = encodeQuotable(writer, &lineLen, ' ', []byte(group))
		if writeError != nil {
			return writeError
		}

		writeError = encodeMultiple(writer, &lineLen, [][]byte{{':'}})
		if writeError != nil {
			return writeError
		}
	}

	for index, to := range tos {
		if index > 0 {
			writeError = encodeMultiple(writer, &lineLen, [][]byte{{','}})
			if writeError != nil {
				return writeError
			}
		}

		addressBytes := []byte(to.Address)
		if to.Name == `` {
			return encodeMultiple(writer, &lineLen, [][]byte{{' '}, addressBytes})
		}

		nameBytes := []byte(to.Name)
		writeError = encodeQuotable(writer, &lineLen, ' ', nameBytes)
		if writeError != nil {
			return writeError
		}

		writeError = encodeMultiple(writer, &lineLen,
			[][]byte{{' ', '<'}, addressBytes, {'>'}})
		if writeError != nil {
			return writeError
		}
	}

	if group != `` {
		writeError = encodeMultiple(writer, &lineLen, [][]byte{{';'}})
		if writeError != nil {
			return writeError
		}
	}

	return nil
}
