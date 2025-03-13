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

package id

import (
	"context"

	"github.com/olive-io/bpmn/v2/pkg/tracing"
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
