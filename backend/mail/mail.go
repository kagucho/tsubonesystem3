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

// Package mail implements a feature to email.
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

// Mail is a structure to hold the context to email. Initialize with mail.New.
type Mail struct {
	templates templates
}

type template struct {
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

// New returns a new mail.Mail.
func New(share string) (Mail, error) {
	var templates templates
	var err error

	htmlBase := path.Join(share, `mail/html`)
	textBase := path.Join(share, `mail/text`)

	for index, file := range [...]string{
		templateConfirmation: `confirmation`,
		templateCreation:     `creation`,
		templateInvitation:   `invitation`,
		templateMessage:      `message`,
	} {
		templates[index].html, err = htmlTemplate.ParseFiles(path.Join(htmlBase, file))
		if err != nil {
			break
		}

		templates[index].text, err = textTemplate.ParseFiles(path.Join(textBase, file))
		if err != nil {
			break
		}
	}

	return Mail{templates}, err
}

func (context Mail) send(host string, recipients []string, toGroup string, tos []mail.Address, subject string, template templateID, data interface{}) (returning error) {
	const boundary = `Copyright(C)2017Kagucho.`

	cmd := exec.Command(`sendmail`, strings.Join(recipients, `,`))

	pipe, pipeErr := cmd.StdinPipe()
	if pipeErr != nil {
		return pipeErr
	}

	if startErr := cmd.Start(); startErr != nil {
		return startErr
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			if killErr := cmd.Process.Kill(); killErr != nil {
				log.Print(killErr)
			}

			if waitErr := cmd.Wait(); waitErr != nil {
				log.Print(waitErr)
			}

			if recoveredErr, ok := recovered.(error); ok {
				returning = recoveredErr
			} else {
				panic(recovered)
			}
		}
	}()

	for _, data := range [...]string{
		"Date: ",
		time.Now().Format(time.RFC1123),
		lineEnding + "From: ",
		(&mail.Address{
			Name: `TsuboneSystem`,
			Address: `noreply-kagucho@` + host,
		}).String(),
		lineEnding,
	} {
		if _, writeErr := pipe.Write([]byte(data)); writeErr != nil {
			panic(writeErr)
		}
	}

	if writeErr := writeTo(pipe, toGroup, tos); writeErr != nil {
		panic(writeErr)
	}

	if subject != `` {
		if _, writeErr := pipe.Write([]byte(lineEnding)); writeErr != nil {
			panic(writeErr)
		}

		if writeErr := writeSubject(pipe, subject); writeErr != nil {
			panic(writeErr)
		}
	}

	if _, writeErr := pipe.Write([]byte(
		lineEnding +
			contentType + ": multipart/alternative; boundary=" + boundary + lineEnding +
			"MIME-Version: 1.0" + lineEnding +
			lineEnding)); writeErr != nil {
		panic(writeErr)
	}

	multipartWriter := multipart.NewWriter(pipe)

	boundaryErr := multipartWriter.SetBoundary(boundary)
	if boundaryErr != nil {
		panic(boundaryErr)
	}

	textPart, textPartErr := multipartWriter.CreatePart(textMIME)
	if textPartErr != nil {
		panic(textPartErr)
	}

	textEncoding := quotedprintable.NewWriter(textPart)
	textExecuteErr := context.templates[template].text.Execute(textEncoding, data)

	if closeErr := textEncoding.Close(); closeErr != nil {
		panic(closeErr)
	}

	if textExecuteErr != nil {
		panic(textExecuteErr)
	}

	htmlPart, htmlPartErr := multipartWriter.CreatePart(htmlMIME)
	if htmlPartErr != nil {
		panic(htmlPartErr)
	}

	htmlEncoding := quotedprintable.NewWriter(htmlPart)
	htmlExecuteErr := context.templates[template].html.Execute(htmlEncoding, data)

	if closeErr := htmlEncoding.Close(); closeErr != nil {
		panic(closeErr)
	}

	if htmlExecuteErr != nil {
		panic(htmlExecuteErr)
	}

	if closeErr := pipe.Close(); closeErr != nil {
		panic(closeErr)
	}

	return cmd.Wait()
}
