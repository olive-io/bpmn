package tracing

// Trace is an interface for actual data traces
type Trace interface {
	TraceInterface()
}

// SenderHandle is an interface for registered senders
type SenderHandle interface {
	// Done indicates that the sender has terminated
	Done()
}

type Tracer interface {
	// Subscribe creates a new unbuffered channel and subscribes it to
	// traces from the Tracer
	//
	// Note that this channel should be continuously read from until unsubscribed
	// from, otherwise, the Tracer will block.
	Subscribe() chan Trace

	// SubscribeChannel subscribe a channel to traces from the Tracer
	//
	// Note that this channel should be continuously read from (modulo buffering),
	// otherwise, the Tracer will block.
	SubscribeChannel(channel chan Trace) chan Trace

	// Unsubscribe removes channel from subscription list
	Unsubscribe(channel chan Trace)

	// Trace sends in a trace to a tracer
	Trace(trace Trace)

	// RegisterSender register a sender for terminate purposes
	//
	// Once Sender is being terminated, before closing subscription channels,
	// it'll wait until all senders call SenderHandle.Done
	RegisterSender() SenderHandle

	// Done returns a channel that is closed when the tracer is done and terminated
	Done() chan struct{}
}

// WrappedTrace is a trace that wraps another trace.
//
// The purpose of it is to allow components to produce traces that will
// be wrapped into additional context, without being aware of it.
//
// Typically this would be done by creating a NewTraceTransformingTracer tracer
// over the original one and passing it to such components.
//
// Consumers looking for individual traces should use Unwrap to retrieve
// the original trace (as opposed to the wrapped one)
type WrappedTrace interface {
	Trace
	// Unwrap returns a wrapped trace
	Unwrap() Trace
}

// Unwrap will recursively unwrap a trace if wrapped,
// or return the trace as is if it isn't wrapped.
func Unwrap(trace Trace) Trace {
	for {
		if unwrapped, ok := trace.(WrappedTrace); ok {
			trace = unwrapped.Unwrap()
		} else {
			return trace
		}
	}
}
