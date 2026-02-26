/*
Copyright 2026 The bpmn Authors

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

package bpmn

import (
	"context"
	"fmt"

	"github.com/olive-io/bpmn/v2/pkg/id"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

const (
	defaultIDGeneratorScopeProcess = "process"
)

var buildDefaultIDGenerator = func(ctx context.Context, tracer tracing.ITracer) (id.IGenerator, error) {
	return id.GetSno().NewIdGenerator(ctx, tracer)
}

func ensureDefaultIDGenerator(ctx context.Context, tracer tracing.ITracer, scope string) id.IGenerator {
	idGenerator, err := buildDefaultIDGenerator(ctx, tracer)
	if err != nil {
		tracer.Send(id.WarningTrace{Warning: fmt.Errorf("init sno id generator failed in %s, fallback to local generator: %w", scope, err)})
		return id.NewFallbackGenerator()
	}
	return idGenerator
}
