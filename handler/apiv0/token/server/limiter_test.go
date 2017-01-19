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

package server

import (
	"testing"
	"time"
)

/*
	This function is an exhaustive test.
	1. Creating new timer
	2. Incrementing a used rate
	3. Updating the order of expiration of an used rate
	4. Incrementing an unused rate
	5. Updating the order of expiration of a unused rate
	6. Expiring timer
	7. Unlimiter blocking for another entry when all entries got removed.
*/
func TestLimiter(t *testing.T) {
	var limiter *limiter

	if !t.Run(`newLimiter`, func(t *testing.T) {
		limiter = newLimiter()
	}) {
		t.FailNow()
	}

	// 1 and 2
	for count := 0; count <= 3; count++ {
		if !limiter.challenge(`used`) {
			t.Errorf(`failed trial %v of "used", expected success`,
				count)
		}
	}

	// 3, give the unlimiter some chance to use "used"
	time.Sleep(8388608)

	// 1 and 3
	if !limiter.challenge(`reorderUsed`) {
		t.Error(`failed trial 1 of "reorderUsed", expected success`)
	}

	// 2 and 3
	if !limiter.challenge(`used`) {
		t.Error(`failed trial 4 of "used", expected success`)
	}

	// 3 and 4
	limiter.rate[`used`].timer.Reset(0)

	// 3, give the unlimiter some chance to expire "used"
	time.Sleep(8388608)

	// 3
	if limiter.challenge(`used`) {
		t.Error(`succeeded trial of "used" before the expiration of "reorderUsed"`)
	}

	// 1 and 4
	for count := 0; count <= 3; count++ {
		if !limiter.challenge(`unused`) {
			t.Errorf(`failed trial %v of "unused", expected success`, count)
		}
	}

	// 1 and 5
	if !limiter.challenge(`reorderUnused`) {
		t.Error(`failed trial 1 of "reorderUnused", expected success`)
	}

	// 4 and 5
	if !limiter.challenge(`unused`) {
		t.Error(`failed trial 4 of "unused", expected success`)
	}

	// Removing "reorderUsed", which has kept "unused" unused.
	limiter.rate[`reorderUsed`].timer.Reset(0)

	// 5
	limiter.rate[`unused`].timer.Reset(0)

	/*
		5.
		Give the unlimiter some chance to expire "used", "reorderUsed",
		and eventually "unused".
	*/
	time.Sleep(8388608)

	// 6
	if !limiter.challenge(`used`) {
		t.Error(`failed trial of "used" after expiring it and its predecessor, expected success`)
	}

	// 5
	if limiter.challenge(`unused`) {
		t.Error(`succeeded trial of "unused" before the expiration of "reorderUnused", expected failure`)
	}

	// 6 and 7
	limiter.rate[`reorderUnused`].timer.Reset(0)
	limiter.rate[`used`].timer.Reset(0)

	// 6 and 7
	// Give the unlimiter some chance to expire "reorderUnused" and "unused".
	time.Sleep(8388608)

	// 6 and 7
	if !limiter.challenge(`unused`) {
		t.Error(`failed trial of "unused" after expiring it and its predecessor, expected success`)
	}
}
