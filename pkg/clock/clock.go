// Copyright 2023 The bpmn Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
