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

package jwt

import (
  `strings`
  `testing`
  `time`
)

func (context JWT) testIssue(t *testing.T) {
  const duration = 1073741824

  before := time.Now().Add(duration).Unix()

  issued, issueError := context.Issue(`sub`, `scope`, duration)
  if issueError != nil {
    t.Fatal(issueError)
  }

  after := time.Now().Add(duration).Unix()

  splited := strings.Split(issued, `.`)

  splitedLen := len(splited)
  if splitedLen != 3 {
    t.Error(`expected 3 elements, found %q elements`, splitedLen)
  }

  testHeader(t, splited[0], context.authority.Alg())
  testClaim(t, splited[1], `sub`, `scope`, before, after)
}
