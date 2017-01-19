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
	textTemplate "text/template"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"net/url"
	"path"
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

var textMIME = textproto.MIMEHeader{
	contentType: contentTypeText,
	contentTransferEncoding: contentTransferEncodingCommon,
}

var htmlMIME = textproto.MIMEHeader{
	contentType: contentTypeHTML,
	contentTransferEncoding: contentTransferEncodingCommon,
}

func New(share string) (Mail, error) {
	var mail Mail
	var parseError error

	mail.html, parseError = htmlTemplate.ParseFiles(path.Join(share, "mail/html/creation"))
	if parseError == nil {
		mail.text, parseError = textTemplate.ParseFiles(path.Join(share, "mail/text/creation"))
	}

	return mail, parseError
}

func (context Mail) MailCreation(id, address, nickname, host string, token string) error {
	var data struct {
		Base string
		Register string
	}

	constructing := url.URL{Scheme: `https`, Host: host}
	data.Base = constructing.String()

	constructing.Path = `/private`
	constructing.Fragment = `token=`+token
	data.Register = constructing.String()

	var buffer bytes.Buffer

	buffer.WriteString("Date: ")
	buffer.WriteString(time.Now().Format(time.RFC1123))

	buffer.WriteString("\r\n"+
		"From: TsuboneSystem <noreply-tsubonesystem@kagucho.net>\r\n"+
		"To: ")

	buffer.WriteString((&mail.Address{nickname, address}).String())

	buffer.WriteString("\r\nSubject: ")
	buffer.WriteString(mime.QEncoding.Encode(`UTF-8`, `TsuboneSystem 登録手続き`))

	const boundary = `Copyright(C)2017Kagucho.LicensedUnderAGPL-3.0.`

	buffer.WriteString("\r\n"+
		contentType+": multipart/alternative; boundary="+boundary+"\r\n"+
		"\r\n")

	multipartWriter := multipart.NewWriter(&buffer)
	multipartWriter.SetBoundary(boundary)

	textPart, textPartError := multipartWriter.CreatePart(textMIME)
	if textPartError != nil {
		return textPartError
	}

	textEncoding := quotedprintable.NewWriter(textPart)
	textExecuteError := context.text.Execute(textEncoding, data)
	textEncoding.Close()

	if textExecuteError != nil {
		return textExecuteError
	}

	htmlPart, htmlPartError := multipartWriter.CreatePart(htmlMIME)
	if htmlPartError != nil {
		return htmlPartError
	}

	htmlEncoding := quotedprintable.NewWriter(htmlPart)
	htmlExecuteError := context.html.Execute(htmlEncoding, data)
	htmlEncoding.Close()

	if htmlExecuteError != nil {
		return htmlExecuteError
	}

	log.Print(smtp.SendMail(`smtp.gmail.com:587`,
		smtp.PlainAuth(``, `root.3.173210@gmail.com`, `hlckqmfstwpjqfck`, `smtp.gmail.com`),
		`noreply-tsubonesystem@kagucho.net`, []string{address},
		buffer.Bytes()))

	return nil
}
