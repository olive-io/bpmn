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
				//return
			case trace, ok := <-ch:
				if ok {
					traces := transformer(trace)
					for _, t := range traces {
						out.Send(t)
					}
				}
			}
		}
	}()
}
