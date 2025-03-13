/*
Copyright 2023 The bpmn Authors

This library is free software; you can redistribute it and/or
modify it under the terms of the GNU Lesser General Public
License as published by the Free Software Foundation; either
version 2.1 of the License, or (at your option) any later version.

This library is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public
License along with this library;
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
