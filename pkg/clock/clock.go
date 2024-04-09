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

// IClock is a generic interface for clocks
type IClock interface {
	// Now returns current time
	Now() time.Time
	// After returns a channel that will send one and only one
	// timestamp once it waited for the given duration
	//
	// Implementation of this method needs to guarantee that
	// the returned time is either equal or greater than the given
	// one plus duration
	After(time.Duration) <-chan time.Time
	// Until returns a channel that will send one and only one
	// timestamp once the provided time has occurred
	//
	// Implementation of this method needs to guarantee that
	// the returned time is either equal or greater than the given
	// one
	Until(time.Time) <-chan time.Time
	// Changes returns a channel that will send a message when
	// the clock detects a change (in case of a real clock, can
	// be a significant drift forward, or a drift backward; for
	// mock clock, an explicit change)
	Changes() <-chan time.Time
}

type contextKey string

func (c contextKey) String() string {
	return "clock package context key " + string(c)
}

// FromContext retrieves a Clock from a given context,
// if there's any. If there's none, it'll create a Host
// clock
func FromContext(ctx context.Context) (c IClock, err error) {
	val := ctx.Value(contextKey("clock"))
	if val == nil {
		return Host(ctx)
	}
	c = val.(IClock)
	return
}

// ToContext saves Clock into a given context, returning a new one
func ToContext(ctx context.Context, clock IClock) context.Context {
	return context.WithValue(ctx, contextKey("clock"), clock)
}
