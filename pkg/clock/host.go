/*
   Copyright 2023 The bpmn Authors

   This program is offered under a commercial and under the AGPL license.
   For AGPL licensing, see below.

   AGPL licensing:
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package clock

import (
	"context"
	"time"
)

type host struct {
	changes chan time.Time
}

func (h *host) Now() time.Time {
	return time.Now()
}

func (h *host) After(duration time.Duration) <-chan time.Time {
	expected := h.Now().Add(duration)
	in := time.After(duration)
	out := make(chan time.Time, 1)
	go func() {
		c := in
		for {
			t := <-c
			if t.Equal(expected) || t.After(expected) {
				out <- t
				return
			} else {
				c = time.After(expected.Sub(t))
			}
		}
	}()
	return out
}

func (h *host) Until(t time.Time) <-chan time.Time {
	return time.After(time.Until(t))
}

func (h *host) Changes() <-chan time.Time {
	return h.changes
}

// Host is a clock source that uses time package as a source
// of time.
func Host(ctx context.Context) (c IClock, err error) {
	changes := make(chan time.Time)
	c = &host{changes: changes}
	err = changeMonitor(ctx, changes)
	return
}
