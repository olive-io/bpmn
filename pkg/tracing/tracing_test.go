/*
   Copyright 2025 The  Authors

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

package tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type myTrace struct{}

func (t *myTrace) Unpack() any {
	return struct{}{}
}

type myTracer struct {
	ITrace
}

func (t *myTracer) Unwrap() ITrace {
	return t.ITrace
}

func TestUnwrap(t *testing.T) {
	tr1 := &myTrace{}
	mt := &myTracer{ITrace: tr1}

	unwrapped := Unwrap(mt)
	assert.Equal(t, unwrapped, tr1)
	assert.Equal(t, unwrapped, Unwrap(unwrapped))
}

func TestNewTracer(t *testing.T) {
	trr := NewTracer(context.Background())

	out := make(chan int, 1)
	sub := trr.Subscribe()
	defer trr.Unsubscribe(sub)
	go func() {
		select {
		case <-sub:
			out <- 1
		}
	}()

	tr := &myTrace{}
	trr.Send(tr)

	assert.Equal(t, <-out, 1)
}
