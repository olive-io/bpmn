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

package tracing

// ITrace is an interface for actual data traces
// This is used for observability of systems, such as letting
// flow nodes, tests and other components to know what happened
// in other parts of the system.
type ITrace interface {
	Unpack() any
}

// ISenderHandle is an interface for registered senders
type ISenderHandle interface {
	// Done indicates that the sender has terminated
	Done()
}

type ITracer interface {
	// Subscribe creates a new unbuffered channel and subscribes it to
	// traces from the Tracer
	//
	// Note that this channel should be continuously read from until unsubscribed
	// from, otherwise, the Tracer will block.
	Subscribe() chan ITrace

	// SubscribeChannel subscribe a channel to traces from the Tracer
	//
	// Note that this channel should be continuously read from (modulo buffering),
	// otherwise, the Tracer will block.
	SubscribeChannel(channel chan ITrace) chan ITrace

	// Unsubscribe removes channel from subscription list
	Unsubscribe(channel chan ITrace)

	// Send sends in a trace to a tracer
	Send(trace ITrace)

	// RegisterSender register a sender for terminate purposes
	//
	// Once Sender is being terminated, before closing subscription channels,
	// it'll wait until all senders call SenderHandle.Done
	RegisterSender() ISenderHandle

	// Done returns a channel that is closed when the tracer is done and terminated
	Done() chan struct{}
}

// ITraceW is a trace that wraps another trace.
//
// The purpose of it is to allow components to produce traces that will
// be wrapped into additional context, without being aware of it.
//
// Typically, this would be done by creating a NewTraceTransformingTracer tracer
// over the original one and passing it to such components.
//
// Consumers looking for individual traces should use Unwrap to retrieve
// the original trace (as opposed to the wrapped one)
type ITraceW interface {
	ITrace
	// Unwrap returns a wrapped trace
	Unwrap() ITrace
}

// Unwrap will recursively unwrap a trace if wrapped,
// or return the trace as is if it isn't wrapped.
func Unwrap(trace ITrace) ITrace {
	for {
		if unwrapped, ok := trace.(ITraceW); ok {
			trace = unwrapped.Unwrap()
		} else {
			return trace
		}
	}
}
