/*
  Copyright (C) 2016  Kagucho <kagucho.net@gmail.com>

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU Affero General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU Affero General Public License for more details.

  You should have received a copy of the GNU Affero General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

/*
  Command tsubonesystem3 implements a Web service for Kagucho, a club in Tokyo
  University of Science.
 */
package main

import (
  `github.com/kagucho/tsubonesystem3/configuration`
  `github.com/kagucho/tsubonesystem3/db`
  `github.com/kagucho/tsubonesystem3/handler/apiv0`
  `github.com/kagucho/tsubonesystem3/handler/private`
  `github.com/kagucho/tsubonesystem3/handler/file`
  `github.com/kardianos/osext`
  `log`
  `math/rand`
  `net`
  `net/http`
  `net/http/fcgi`
  `os`
  `os/signal`
  `path`
  `syscall`
  `time`
)

func getShare() (string, error) {
  executable, executableError := osext.ExecutableFolder()
  if executableError != nil {
    return executable, executableError
  }

  return path.Join(executable, `../share/tsubonesystem3`), nil
}

func main() {
  share, shareError := getShare()
  if shareError != nil {
    log.Fatalln(shareError)
  }

  fileError, fileErrorError := file.NewError(share)
  if fileErrorError != nil {
    log.Fatalln(fileErrorError)
  }

  db, dbError := db.Open()
  if dbError != nil {
    log.Fatalln(dbError)
  }

  apiv0, apiv0Error := apiv0.New(db)
  if apiv0Error != nil {
    if closeError := db.Close(); closeError != nil {
      log.Println(closeError)
    }

    log.Fatalln(apiv0Error)
  }

  private, privateError := private.New(share, db, fileError)
  if privateError != nil {
    if closeError := db.Close(); closeError != nil {
      log.Println(closeError)
    }

    log.Fatalln(apiv0Error)
  }

  file := file.New(share, fileError)

  listener, listenerError := net.Listen(configuration.ListenNet,
                                        configuration.ListenAddress)
  if listenerError != nil {
    if closeError := db.Close(); closeError != nil {
      log.Println(closeError)
    }

    log.Fatalln(listenerError)
  }

  go func() {
    signalChan := make(chan os.Signal, 2)
    signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
    <-signalChan

    if closeError := listener.Close(); closeError != nil {
      log.Println(closeError)
    }

    if closeError := db.Close(); closeError != nil {
      log.Println(closeError)
    }
  }()

  // random function is used in TODO: update comment
  rand.Seed(time.Now().Unix())

  http.Handle(`/api/v0/`, http.StripPrefix(`/api/v0`, apiv0))
  http.Handle(`/private`, private)
  http.Handle(`/`, file)

  fcgi.Serve(listener, nil)
}
