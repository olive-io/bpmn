package id

import (
	"context"

	"github.com/olive-io/bpmn/tracing"
)

type IGeneratorBuilder interface {
	NewIdGenerator(ctx context.Context, tracer tracing.ITracer) (IGenerator, error)
	RestoreIdGenerator(ctx context.Context, serialized []byte, tracer tracing.ITracer) (IGenerator, error)
}

type IGenerator interface {
	Snapshot() ([]byte, error)
	New() Id
}

type Id interface {
	Bytes() []byte
	String() string
}

var DefaultIdGeneratorBuilder IGeneratorBuilder
