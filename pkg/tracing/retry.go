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
