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

package main

import (
	"github.com/kagucho/tsubonesystem3/backend"
	"github.com/kagucho/tsubonesystem3/configuration"
	"log"
	"net"
	"net/http/fcgi"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	backend, backendError := backend.New()
	if backendError != nil {
		log.Panic(backendError)
	}

	listener, listenerError := net.Listen(configuration.FcgiListenNet,
		configuration.FcgiListenAddress)
	if listenerError != nil {
		if closeError := backend.Close(); closeError != nil {
			log.Print(closeError)
		}

		log.Panic(listenerError)
	}

	go func() {
		signalChan := make(chan os.Signal, 2)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		<-signalChan

		if closeError := listener.Close(); closeError != nil {
			log.Print(closeError)
		}

		if closeError := backend.Close(); closeError != nil {
			log.Print(closeError)
		}
	}()

	fcgi.Serve(listener, backend)
}
