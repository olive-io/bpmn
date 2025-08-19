/*
Copyright 2025 The bpmn Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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

func TestNewRelay(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	in := NewTracer(ctx)
	out := NewTracer(ctx)

	NewRelay(ctx, in, out, func(trace ITrace) []ITrace {
		return []ITrace{trace}
	})

	sender := out.Subscribe()

	mt := &myTrace{}
	in.Send(mt)

	recv := <-sender

	assert.Equal(t, mt, recv)
}
