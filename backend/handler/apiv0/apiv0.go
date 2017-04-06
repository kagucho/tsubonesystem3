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

// Package apiv0 implements API v0.
package apiv0

import (
	"github.com/kagucho/tsubonesystem3/backend/db"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/token/backend"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/util"
	"github.com/kagucho/tsubonesystem3/backend/mail"
	"github.com/kagucho/tsubonesystem3/safehttp"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

type shared struct {
	DB    db.DB
	Mail  mail.Mail
	Token backend.Backend // FIXME: why exported?
}

/*
APIv0 is a structure to hold the context of API v0. It should be initialized
with apiv0.New.
*/
type APIv0 struct {
	shared       shared
	routes       routeSlice
	tokenServer  *tokenServer
}

/*
New returns a new apiv0.APIv0. End must be called before disposing returned
apiv0.APIv0.
*/
func New(db db.DB, mail mail.Mail) (APIv0, error) {
	token, err := backend.New()
	if err != nil {
		return APIv0{}, err
	}

	apiv0 := APIv0{
		shared: shared{db, mail, token},
		tokenServer: newTokenServer(),
	}

	apiv0.routes = apiv0.newRoutes()

	return apiv0, nil
}

/*
End releases the resources.

After calling it, any functions bound to the apiv0.APIv0 will result in
unexpected result.
*/
func (apiv0 APIv0) End() {
	apiv0.tokenServer.end()
}

// ServeHTTP serves API v0 via HTTP.
func (apiv0 APIv0) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	safeWriter := safehttp.NewWriter(writer)

	defer safeWriter.Recover(func() {
		util.ServeErrorDefault(writer, http.StatusInternalServerError)

		log.Print(recover())
		debug.PrintStack()
	})

	index := apiv0.routes.search(request.URL.Path)
	if index > len(apiv0.routes) || apiv0.routes[index].prefix != request.URL.Path {
		index--
		if index < 0 || !strings.HasPrefix(request.URL.Path, apiv0.routes[index].prefix) || request.URL.Path[len(apiv0.routes[index].prefix)] != '/' {
			util.ServeErrorDefault(&safeWriter, http.StatusNotFound)
			return
		}
	}

	request.URL.Path = strings.TrimPrefix(request.URL.Path, apiv0.routes[index].prefix)
	apiv0.routes[index].handler.serveHTTP(&safeWriter, request, apiv0.shared)
}
