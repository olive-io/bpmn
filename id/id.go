package id

import (
	"context"

	"github.com/olive-io/bpmn/tracing"
)

type GeneratorBuilder interface {
	NewIdGenerator(ctx context.Context, tracer tracing.Tracer) (Generator, error)
	RestoreIdGenerator(ctx context.Context, serialized []byte, tracer tracing.Tracer) (Generator, error)
}

type Generator interface {
	Snapshot() ([]byte, error)
	New() Id
}

type Id interface {
	Bytes() []byte
	String() string
}

var DefaultIdGeneratorBuilder GeneratorBuilder
