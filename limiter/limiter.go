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

// Package limiter implements a simple rate limiter.
package limiter

import (
	"container/list"
	"sync"
	"time"
)

type execution bool

const (
	executionEnd     = false
	executionUnlimit = true
)

type rate struct {
	count   uint
	element *list.Element
	timer   *time.Timer
	id      string
}

// Limiter is a structure holding the context of the rate limiter. It should be
// initialized with limiter.New.
type Limiter struct {
	execution chan execution
	end       bool
	list      list.List
	rate      map[string]*rate
	mutex     sync.Mutex
}

func (limiter *Limiter) unlimit() {
	listLen := 0

	for {
		if listLen == 0 {
			execution := <-limiter.execution
			if execution == executionEnd {
				return
			}

			limiter.mutex.Lock()
		}

		entry := limiter.list.Front()
		rateEntry := entry.Value.(*rate)

		// Tell the timer is potentially used.
		rateEntry.element = nil

		// Read timer and prevent race condition.
		timer := rateEntry.timer

		limiter.mutex.Unlock()

		select {
		case <-limiter.execution:
			return

		case <-timer.C:
		}

		limiter.mutex.Lock()

		// If there is no more trial after waiting.
		if rateEntry.element == nil {
			delete(limiter.rate, rateEntry.id)
		}

		limiter.list.Remove(entry)
		listLen = limiter.list.Len()

		if listLen == 0 {
			limiter.mutex.Unlock()
		}
	}
}

/*
Challenge records a challenge, associating with the given ID, and returns a bool
telling whether the challenge is successful or not.
*/
func (limiter *Limiter) Challenge(id string) bool {
	const duration = 68719476736
	const limit = 4

	limiter.mutex.Lock()

	var accept bool
	if rateEntry, present := limiter.rate[id]; present {
		if rateEntry.count < limit {
			if rateEntry.element == nil {
				// The timer is potentially used.

				// If the timer is not expired.
				if rateEntry.timer.Stop() {
					// Expire the timer now.
					rateEntry.timer.Reset(0)
					rateEntry.count++
				} else {
					rateEntry.count = 0
				}

				rateEntry.element = limiter.list.PushBack(rateEntry)
				rateEntry.timer = time.NewTimer(duration)
			} else {
				// The timer is not used yet.

				rateEntry.count++
				rateEntry.timer.Reset(duration)
				limiter.list.MoveToBack(rateEntry.element)
			}

			accept = true
		} else {
			accept = false
		}
	} else {
		rateEntry := rate{
			timer: time.NewTimer(duration), id: id,
		}

		limiter.rate[id] = &rateEntry
		rateEntry.element = limiter.list.PushBack(&rateEntry)
		if limiter.list.Len() == 1 {
			limiter.execution <- executionUnlimit
		}

		accept = true
	}

	limiter.mutex.Unlock()

	return accept
}

/*
New returns a new limiter.Limiter. Resources will be holded until End gets
called.
*/
func New() *Limiter {
	instance := Limiter{
		execution: make(chan execution, 1), rate: map[string]*rate{},
	}
	instance.list.Init()

	go instance.unlimit()

	return &instance
}

/*
End releases the resources.

After calling it, calling any functions bound to limiter will result in
an unexpected result.
*/
func (limiter Limiter) End() {
	limiter.execution <- executionEnd
}
