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
	"path/filepath"
)

// Backend is the structure to hold the context of the backend.
type Backend struct {
	unchunked.Unchunked
	apiv0 apiv0.APIv0
	db    db.DB
}

var share string

func init() {
	executable, err := osext.ExecutableFolder()
	if err != nil {
		log.Panic(err)
	}

	share = filepath.Join(executable, `../share/tsubonesystem3`)
}

// New returns a new backend.Backend.
func New() (Backend, error) {
	log.Print("TsuboneSystem3  Copyright (C) 2017  Kagucho <kagucho.net@gmail.com>")
	log.Print("This program comes with ABSOLUTELY NO WARRANTY.")
	log.Print("This is free software, and you are welcome to redistribute it")
	log.Print("under certain conditions; see `/license' for details.")

	fileError, fileErrorErr := file.NewError(share)
	if fileErrorErr != nil {
		return Backend{}, fileErrorErr
	}

	mail, mailErr := mail.New(share)
	if mailErr != nil {
		return Backend{}, mailErr
	}

	db, dbErr := db.New()
	if dbErr != nil {
		return Backend{}, dbErr
	}

	apiv0, apiv0Err := apiv0.New(db, mail)
	if apiv0Err != nil {
		if closeErr := db.Close(); closeErr != nil {
			log.Print(closeErr)
		}

		return Backend{}, apiv0Err
	}

	private, privateErr := private.New(share, db, fileError)
	if privateErr != nil {
		if closeErr := db.Close(); closeErr != nil {
			log.Print(closeErr)
		}

		return Backend{}, privateErr
	}

	file := file.New(share, fileError)

	var mux http.ServeMux

	mux.Handle(`/api/v0/`, http.StripPrefix(`/api/v0`, apiv0))
	mux.Handle(`/private`, private)
	mux.Handle(`/`, file)

	return Backend{unchunked.New(mux.ServeHTTP), apiv0, db}, nil
}

/*
End releases all resources.

This function must be called before disposing backend.Backend returned by
backend.New.

After calling this, calling functions bound to backend.Backend will result in an
unexpected result.
*/
func (backend Backend) End() error {
	backend.apiv0.End()
	return backend.db.Close()
}
