package tracing

import "context"

type Transformer func(trace ITrace) []ITrace

// NewRelay starts a goroutine that relays traces from `in` Tracer to `out` Tracer
// using a given Transformer
//
// This is typically used to compartmentalize tracers.
func NewRelay(ctx context.Context, in, out ITracer, transformer Transformer) {
	ch := in.Subscribe()
	handle := out.RegisterSender()
	go func() {
		for {
			select {
			case <-in.Done():
				handle.Done()
				return
			case <-ctx.Done():
				// wait until `in` Tracer is done
			case trace, ok := <-ch:
				if ok {
					traces := transformer(trace)
					for _, trace := range traces {
						out.Trace(trace)
					}
				}
			}
		}
	}()
}
