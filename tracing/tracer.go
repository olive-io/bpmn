// Copyright 2023 The olive Authors
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

package tracing

import (
	"context"
	"sync"
)

type subscription struct {
	channel chan ITrace
	ok      chan struct{}
}

type unSubscription struct {
	channel chan ITrace
	ok      chan struct{}
}

type tracer struct {
	traces         chan ITrace
	subscription   chan subscription
	unSubscription chan unSubscription
	terminate      chan struct{}
	done           chan struct{}
	subscribers    []chan ITrace
	senders        sync.WaitGroup
}

func NewTracer(ctx context.Context) ITracer {
	t := tracer{
		traces:         make(chan ITrace),
		subscription:   make(chan subscription),
		unSubscription: make(chan unSubscription),
		terminate:      make(chan struct{}),
		done:           make(chan struct{}),
		subscribers:    make([]chan ITrace, 0),
	}
	go t.runner(ctx)
	return &t
}

func (t *tracer) runner(ctx context.Context) {
	var termination sync.Once
	defer close(t.done)

	for {
		select {
		case sch := <-t.subscription:
			t.subscribers = append(t.subscribers, sch.channel)
			sch.ok <- struct{}{}
		case unsch := <-t.unSubscription:
			pos := -1
			for i := range t.subscribers {
				if t.subscribers[i] == unsch.channel {
					pos = i
					break
				}
			}

			if pos >= 0 {
				l := len(t.subscribers) - 1
				// remove subscriber by replacing it with the last one
				t.subscribers[pos] = t.subscribers[l]
				t.subscribers[l] = nil
				// and truncating the list of subscribers
				t.subscribers = t.subscribers[:l]
				// (as we don't care about the order)
				unsch.ok <- struct{}{}
			}
		case trace := <-t.traces:
			for _, subscriber := range t.subscribers {
				subscriber <- trace
			}
		case <-ctx.Done():
			// Start a termination waiting routine (only once)
			termination.Do(func() {
				go func() {
					// Wait until all senders have terminated
					t.senders.Wait()
					// Send an internal termination message
					t.terminate <- struct{}{}
				}()
			})
			// Let tracer continue to work for now
		case <-t.terminate:
			for _, subscriber := range t.subscribers {
				close(subscriber)
			}
			return
		}
	}
}

func (t *tracer) Subscribe() chan ITrace {
	return t.SubscribeChannel(make(chan ITrace))
}

func (t *tracer) SubscribeChannel(channel chan ITrace) chan ITrace {
	okCh := make(chan struct{})
	sub := subscription{channel: channel, ok: okCh}
	t.subscription <- sub
	<-okCh
	return channel
}

func (t *tracer) Unsubscribe(channel chan ITrace) {
	okChan := make(chan struct{})
	unsub := unSubscription{channel: channel, ok: okChan}
loop:
	for {
		select {
		// If the tracer is done, it's as good as if we're unsubscribed
		case <-t.Done():
			return
		case <-channel:
			continue loop
		case t.unSubscription <- unsub:
			continue loop
		case <-okChan:
			return
		}
	}
}

func (t *tracer) Trace(trace ITrace) {
	t.traces <- trace
}

func (t *tracer) RegisterSender() ISenderHandle {
	t.senders.Add(1)
	return &t.senders
}

func (t *tracer) Done() chan struct{} {
	return t.done
}
