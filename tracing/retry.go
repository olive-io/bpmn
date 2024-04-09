/*
   Copyright 2023 The bpmn Authors

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
