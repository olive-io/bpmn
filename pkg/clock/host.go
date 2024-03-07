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
