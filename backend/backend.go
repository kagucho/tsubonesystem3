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

package backend

import (
	"github.com/kagucho/tsubonesystem3/backend/db"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0"
	"github.com/kagucho/tsubonesystem3/backend/handler/file"
	"github.com/kagucho/tsubonesystem3/backend/handler/private"
	"github.com/kagucho/tsubonesystem3/backend/mail"
	"github.com/kagucho/tsubonesystem3/unchunked"
	"github.com/kardianos/osext"
	"log"
	"net/http"
	"path"
)

type Backend struct {
	unchunked.Unchunked
	db db.DB
}

var share string

func init() {
	executable, executableError := osext.ExecutableFolder()
	if executableError != nil {
		log.Panic(executableError)
	}

	share = path.Join(executable, `../share/tsubonesystem3`)
}

func New() (Backend, error) {
	fileError, fileErrorError := file.NewError(share)
	if fileErrorError != nil {
		return Backend{}, fileErrorError
	}

	mail, mailError := mail.New(share)
	if mailError != nil {
		return Backend{}, mailError
	}

	db, dbError := db.Prepare()
	if dbError != nil {
		return Backend{}, dbError
	}

	apiv0, apiv0Error := apiv0.New(db, mail)
	if apiv0Error != nil {
		if closeError := db.Close(); closeError != nil {
			log.Print(closeError)
		}

		return Backend{}, apiv0Error
	}

	private, privateError := private.New(share, db, fileError)
	if privateError != nil {
		if closeError := db.Close(); closeError != nil {
			log.Print(closeError)
		}

		return Backend{}, privateError
	}

	file := file.New(share, fileError)

	var mux http.ServeMux

	mux.Handle(`/api/v0/`, http.StripPrefix(`/api/v0`, apiv0))
	mux.Handle(`/private`, private)
	mux.Handle(`/`, file)

	return Backend{unchunked.New(mux.ServeHTTP), db}, nil
}

func (backend Backend) Close() error {
	return backend.db.Close()
}
