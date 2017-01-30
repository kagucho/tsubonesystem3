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
	"errors"
	"fmt"
	htmlTemplate "html/template"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/textproto"
	"os/exec"
	"path"
	textTemplate "text/template"
	"time"
)

type Mail struct {
	html *htmlTemplate.Template
	text *textTemplate.Template
}

const contentType = `Content-Type`
const contentTransferEncoding = `Content-Transfer-Encoding`

var contentTypeText = []string{`text/plain;charset=UTF-8`}
var contentTypeHTML = []string{`text/html;charset=UTF-8`}
var contentTransferEncodingCommon = []string{`quoted-printable`}

var mails = []string{`creation`, `confirmation`}

var textMIME = textproto.MIMEHeader{
	contentType:             contentTypeText,
	contentTransferEncoding: contentTransferEncodingCommon,
}

var htmlMIME = textproto.MIMEHeader{
	contentType:             contentTypeHTML,
	contentTransferEncoding: contentTransferEncodingCommon,
}

func New(share string) (Mail, error) {
	var mail Mail
	var parseError error

	htmlBase := path.Join(share, "mail/html")
	textBase := path.Join(share, "mail/text")

	mail.html, parseError = htmlTemplate.ParseFiles(path.Join(htmlBase, mails[0]))
	if parseError != nil {
		return mail, parseError
	}

	mail.text, parseError = textTemplate.ParseFiles(path.Join(textBase, mails[0]))
	if parseError != nil {
		return mail, parseError
	}

	for index := 1; index < len(mails); index++ {
		_, parseError = mail.html.ParseFiles(path.Join(htmlBase, mails[index]))
		if parseError != nil {
			break
		}

		_, parseError = mail.text.ParseFiles(path.Join(textBase, mails[index]))
		if parseError != nil {
			break
		}
	}

	return mail, parseError
}

func (context Mail) send(host string, to mail.Address, template string, data interface{}) (returning error) {
	const boundary = `Copyright(C)2017Kagucho.LicensedUnderAGPL-3.0.`

	cmd := exec.Command("sendmail", "-t")

	pipe, pipeError := cmd.StdinPipe()
	if pipeError != nil {
		return pipeError
	}

	startError := cmd.Start()
	if startError != nil {
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
				returning = errors.New(fmt.Sprint(recovered))
			}
		}
	}()

	for _, data := range [...]string{
		"Date: ", time.Now().Format(time.RFC1123),
		"\r\nFrom: ", (&mail.Address{`TsuboneSystem`, `noreply-kagucho@` + host}).String(),
		"\r\nTo: ", to.String(),
		"\r\nSubject: ", mime.QEncoding.Encode(`UTF-8`, `TsuboneSystem 登録手続き`),
		"\r\n" + contentType + ": multipart/alternative; boundary=" + boundary +
		"\r\nMIME-Version: 1.0" +
		"\r\n" +
		"\r\n",
	} {
		if _, writeError := pipe.Write([]byte(data)); writeError != nil {
			panic(writeError)
		}
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
	textExecuteError := context.text.ExecuteTemplate(textEncoding, template, data)

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
	htmlExecuteError := context.html.ExecuteTemplate(htmlEncoding, template, data)

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
