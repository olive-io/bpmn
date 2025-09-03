/*
Copyright 2023 The bpmn Authors

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

package id

import (
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type mockTracer struct {
	traces      []tracing.ITrace
	channels    []chan tracing.ITrace
	done        chan struct{}
	senderCount int
}

func newMockTracer() *mockTracer {
	return &mockTracer{
		traces:   make([]tracing.ITrace, 0),
		channels: make([]chan tracing.ITrace, 0),
		done:     make(chan struct{}),
	}
}

func (m *mockTracer) Send(trace tracing.ITrace) {
	m.traces = append(m.traces, trace)
	for _, ch := range m.channels {
		select {
		case ch <- trace:
		default:
		}
	}
}

func (m *mockTracer) Subscribe() chan tracing.ITrace {
	ch := make(chan tracing.ITrace)
	m.channels = append(m.channels, ch)
	return ch
}

func (m *mockTracer) SubscribeChannel(channel chan tracing.ITrace) chan tracing.ITrace {
	m.channels = append(m.channels, channel)
	return channel
}

func (m *mockTracer) Unsubscribe(channel chan tracing.ITrace) {
	for i, ch := range m.channels {
		if ch == channel {
			m.channels = append(m.channels[:i], m.channels[i+1:]...)
			break
		}
	}
}

func (m *mockTracer) RegisterSender() tracing.ISenderHandle {
	m.senderCount++
	return &mockSenderHandle{tracer: m}
}

func (m *mockTracer) Done() chan struct{} {
	return m.done
}

type mockSenderHandle struct {
	tracer *mockTracer
	done   bool
}

func (m *mockSenderHandle) Done() {
	if !m.done {
		m.done = true
		m.tracer.senderCount--
		if m.tracer.senderCount <= 0 {
			close(m.tracer.done)
		}
	}
}