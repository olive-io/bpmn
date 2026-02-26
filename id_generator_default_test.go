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
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/olive-io/bpmn/schema"
	pkgid "github.com/olive-io/bpmn/v2/pkg/id"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

func TestNewOptionsFallbackGeneratorSendsWarningTrace(t *testing.T) {
	original := buildDefaultIDGenerator
	defer func() { buildDefaultIDGenerator = original }()

	buildDefaultIDGenerator = func(ctx context.Context, tracer tracing.ITracer) (pkgid.IGenerator, error) {
		return nil, errors.New("sno unavailable")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tracer := tracing.NewTracer(ctx)
	traces := tracer.Subscribe()
	defer tracer.Unsubscribe(traces)

	options := NewOptions(WithContext(ctx), WithTracer(tracer))
	require.NotNil(t, options.idGenerator)

	warning, ok := waitWarningTrace(traces, time.Second)
	require.True(t, ok)
	assert.Contains(t, warning, "fallback to local generator")
	assert.Contains(t, warning, "process")
}

func TestNewProcessSetFallbackGeneratorSendsWarningTrace(t *testing.T) {
	original := buildDefaultIDGenerator
	defer func() { buildDefaultIDGenerator = original }()

	buildDefaultIDGenerator = func(ctx context.Context, tracer tracing.ITracer) (pkgid.IGenerator, error) {
		return nil, errors.New("sno unavailable")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tracer := tracing.NewTracer(ctx)
	traces := tracer.Subscribe()
	defer tracer.Unsubscribe(traces)

	definitions := schema.DefaultDefinitions()
	p := schema.NewProcessBuilder().Out()
	p.IsExecutableField = schema.NewBoolP(true)
	definitions.ProcessField = []schema.Process{*p}

	executes := []*schema.Process{&definitions.ProcessField[0]}
	_, err := NewProcessSet(executes, nil, &definitions, WithContext(ctx), WithTracer(tracer))
	require.NoError(t, err)

	warning, ok := waitWarningTrace(traces, time.Second)
	require.True(t, ok)
	assert.Contains(t, warning, "fallback to local generator")
	assert.Contains(t, warning, "process")
}

func waitWarningTrace(traces <-chan tracing.ITrace, timeout time.Duration) (string, bool) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case trace := <-traces:
			trace = tracing.Unwrap(trace)
			if warning, ok := trace.(pkgid.WarningTrace); ok {
				return fmt.Sprint(warning.Warning), true
			}
		case <-timer.C:
			return "", false
		}
	}
}
