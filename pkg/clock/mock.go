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
	"sort"
	"sync"
	"time"
)

type after struct {
	time.Time
	ch chan time.Time
}

type afters []after

func (a afters) Len() int {
	return len(a)
}

func (a afters) Less(i, j int) bool {
	return a[i].Time.UnixNano() < a[j].Time.UnixNano()
}

func (a afters) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Mock clock is a fake clock that doesn't change
// unless explicitly instructed to.
type Mock struct {
	sync.RWMutex
	now     time.Time
	changes chan time.Time
	timers  afters
}

func (m *Mock) Now() time.Time {
	m.RLock()
	defer m.RUnlock()
	now := m.now
	return now
}

func (m *Mock) After(duration time.Duration) <-chan time.Time {
	m.Lock()
	defer m.Unlock()
	ch := make(chan time.Time, 1)
	if duration.Nanoseconds() <= 0 {
		ch <- m.now
		close(ch)
		return ch
	}
	m.timers = append(m.timers, after{Time: m.now.Add(duration), ch: ch})
	return ch
}

func (m *Mock) Until(t time.Time) <-chan time.Time {
	m.Lock()
	defer m.Unlock()
	ch := make(chan time.Time, 1)
	if m.now.Equal(t) || m.now.After(t) {
		ch <- m.now
		close(ch)
		return ch
	}
	m.timers = append(m.timers, after{Time: t, ch: ch})
	return ch
}

func (m *Mock) Changes() <-chan time.Time {
	return m.changes
}

// NewMockAt creates a new Mock clock at a specific
// point in time
func NewMockAt(t time.Time) *Mock {
	source := &Mock{
		now:     t,
		changes: make(chan time.Time, 1),
	}
	return source
}

// NewMock creates a new Mock clock, set at the start of
// UNIX time
func NewMock() *Mock {
	return NewMockAt(time.Unix(0, 0))
}

func (m *Mock) Set(t time.Time) {
	m.Lock()
	defer m.Unlock()
	m.lockedSet(t)
}

func (m *Mock) Add(duration time.Duration) {
	m.Lock()
	defer m.Unlock()
	m.lockedSet(m.now.Add(duration))
}

// lockedSet should only be called when Mock is locked
func (m *Mock) lockedSet(t time.Time) {
	after := make([]after, 0, len(m.timers))
	sort.Sort(m.timers)
	for i := range m.timers {
		if m.timers[i].Equal(t) || m.timers[i].Before(t) {
			m.timers[i].ch <- t
		} else {
			after = append(after, m.timers[i])
		}
	}
	select {
	case m.changes <- t:
		// delivered changes notification
	default:
		// drop old time (we have a buffer of one)
		<-m.changes
		// push out new time
		m.changes <- t
	}
	m.timers = after
	m.now = t
}
