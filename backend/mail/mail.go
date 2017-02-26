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
	htmlTemplate "html/template"
	"log"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/textproto"
	"os/exec"
	"path"
	"strings"
	textTemplate "text/template"
	"time"
)

type Mail struct {
	templates templates
}

type template struct{
	html *htmlTemplate.Template
	text *textTemplate.Template
}

type templates [templateTotal]template

const contentType = `Content-Type`
const contentTransferEncoding = `Content-Transfer-Encoding`
const lineEnding = "\r\n"

var contentTypeText = []string{`text/plain;charset=UTF-8`}
var contentTypeHTML = []string{`text/html;charset=UTF-8`}
var contentTransferEncodingCommon = []string{`quoted-printable`}

type templateID uint

const (
	templateConfirmation templateID = iota
	templateCreation
	templateInvitation
	templateMessage

	templateTotal
)

var textMIME = textproto.MIMEHeader{
	contentType:             contentTypeText,
	contentTransferEncoding: contentTransferEncodingCommon,
}

var htmlMIME = textproto.MIMEHeader{
	contentType:             contentTypeHTML,
	contentTransferEncoding: contentTransferEncodingCommon,
}

func New(share string) (Mail, error) {
	var templates templates
	var parseError error

	htmlBase := path.Join(share, `mail/html`)
	textBase := path.Join(share, `mail/text`)

	for index, file := range [...]string{
		templateConfirmation: `confirmation`,
		templateCreation:     `creation`,
		templateInvitation:   `invitation`,
		templateMessage:      `broadcast`,
	} {
		templates[index].html, parseError = htmlTemplate.ParseFiles(path.Join(htmlBase, file))
		if parseError != nil {
			break
		}

		templates[index].text, parseError = textTemplate.ParseFiles(path.Join(textBase, file))
		if parseError != nil {
			break
		}
	}

	return Mail{templates}, parseError
}

func (context Mail) send(host string, recipients []string, toGroup string, tos []mail.Address, subject string, template templateID, data interface{}) (returning error) {
	const boundary = `Copyright(C)2017Kagucho.`

	cmd := exec.Command(`sendmail`, strings.Join(recipients, `,`))

	pipe, pipeError := cmd.StdinPipe()
	if pipeError != nil {
		return pipeError
	}

	if startError := cmd.Start(); startError != nil {
		return startError
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			if killError := cmd.Process.Kill(); killError != nil {
				log.Print(killError)
			}

			if waitError := cmd.Wait(); waitError != nil {
				log.Print(waitError)
			}

			if recoveredError, ok := recovered.(error); ok {
				returning = recoveredError
			} else {
				panic(recovered)
			}
		}
	}()

	for _, data := range [...]string{
		"Date: ", time.Now().Format(time.RFC1123),
		lineEnding + "From: ", (&mail.Address{`TsuboneSystem`, `noreply-kagucho@` + host}).String(),
		lineEnding,
	} {
		if _, writeError := pipe.Write([]byte(data)); writeError != nil {
			panic(writeError)
		}
	}

	if writeError := writeTo(pipe, toGroup, tos); writeError != nil {
		panic(writeError)
	}

	if subject != `` {
		if _, writeError := pipe.Write([]byte(lineEnding)); writeError != nil {
			panic(writeError)
		}

		if writeError := writeSubject(pipe, subject); writeError != nil {
			panic(writeError)
		}
	}

	if _, writeError := pipe.Write([]byte(
		lineEnding +
		contentType + ": multipart/alternative; boundary=" + boundary + lineEnding +
		"MIME-Version: 1.0" + lineEnding +
		lineEnding)); writeError != nil {
		panic(writeError)
	}

	multipartWriter := multipart.NewWriter(pipe)

	boundaryError := multipartWriter.SetBoundary(boundary)
	if boundaryError != nil {
		panic(boundaryError)
	}

	textPart, textPartError := multipartWriter.CreatePart(textMIME)
	if textPartError != nil {
		panic(textPartError)
	}

	textEncoding := quotedprintable.NewWriter(textPart)
	textExecuteError := context.templates[template].text.Execute(textEncoding, data)

	if closeError := textEncoding.Close(); closeError != nil {
		panic(closeError)
	}

	if textExecuteError != nil {
		panic(textExecuteError)
	}

	htmlPart, htmlPartError := multipartWriter.CreatePart(htmlMIME)
	if htmlPartError != nil {
		panic(htmlPartError)
	}

	htmlEncoding := quotedprintable.NewWriter(htmlPart)
	htmlExecuteError := context.templates[template].html.Execute(htmlEncoding, data)

	if closeError := htmlEncoding.Close(); closeError != nil {
		panic(closeError)
	}

	if htmlExecuteError != nil {
		panic(htmlExecuteError)
	}

	if closeError := pipe.Close(); closeError != nil {
		panic(closeError)
	}

	return cmd.Wait()
}
