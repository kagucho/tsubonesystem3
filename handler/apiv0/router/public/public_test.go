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

package public

import (
  `github.com/kagucho/tsubonesystem3/handler/apiv0/token`
  `testing`
)

func TestPublic(t *testing.T) {
  tokenServer, _, tokenError := token.New()
  if tokenError != nil {
    t.Fatal(tokenError)
  }

  var public Public
  if !t.Run(`New`, func(t *testing.T) {
    public = New(tokenServer)
  }) {
    t.FailNow()
  }

  t.Run(`Route`, func(t *testing.T) {
    if public.GetRoute(`/token`) == nil {
      t.Error(`no /member/find`)
    }

    if public.GetRoute(`/invalid`) != nil {
      t.Error(`invalid /invalid; expected (nil)`)
    }
  })
}
