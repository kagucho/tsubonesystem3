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
	"bytes"
	htmlTemplate "html/template"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"path"
	"strings"
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

func (context Mail) send(host string, to mail.Address, template string, data interface{}) error {
	from := mail.Address{`TsuboneSystem`, `noreply-kagucho@` + host}

	var buffer bytes.Buffer

	buffer.WriteString("Date: ")
	buffer.WriteString(time.Now().Format(time.RFC1123))

	buffer.WriteString("\r\nFrom: ")
	buffer.WriteString(from.String())

	buffer.WriteString("\r\nTo: ")
	buffer.WriteString(to.String())

	buffer.WriteString("\r\nSubject: ")
	buffer.WriteString(mime.QEncoding.Encode(`UTF-8`, `TsuboneSystem 登録手続き`))

	const boundary = `Copyright(C)2017Kagucho.LicensedUnderAGPL-3.0.`

	buffer.WriteString("\r\n" +
		contentType + ": multipart/alternative; boundary=" + boundary + "\r\n" +
		"\r\n")

	multipartWriter := multipart.NewWriter(&buffer)
	multipartWriter.SetBoundary(boundary)

	textPart, textPartError := multipartWriter.CreatePart(textMIME)
	if textPartError != nil {
		return textPartError
	}

	textEncoding := quotedprintable.NewWriter(textPart)
	textExecuteError := context.text.ExecuteTemplate(textEncoding, template, data)
	textEncoding.Close()

	if textExecuteError != nil {
		return textExecuteError
	}

	htmlPart, htmlPartError := multipartWriter.CreatePart(htmlMIME)
	if htmlPartError != nil {
		return htmlPartError
	}

	htmlEncoding := quotedprintable.NewWriter(htmlPart)
	htmlExecuteError := context.html.ExecuteTemplate(htmlEncoding, template, data)
	htmlEncoding.Close()

	if htmlExecuteError != nil {
		return htmlExecuteError
	}

	/*
		return smtp.SendMail(`smtp.gmail.com:587`,
			smtp.PlainAuth(``, `root.3.173210@gmail.com`, `hlckqmfstwpjqfck`, `smtp.gmail.com`),
			`noreply-tsubonesystem@`+host, []string{address},
			buffer.Bytes())
	*/
	return smtp.SendMail(to.Address[strings.IndexByte(to.Address, '@')+1:]+":25", nil,
		from.Address, []string{to.Address},
		buffer.Bytes())
}
