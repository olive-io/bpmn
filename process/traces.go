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

package process

import (
	"github.com/olive-io/bpmn/schema"

	"github.com/olive-io/bpmn/tracing"
)

// Trace wraps any trace within a given process
type Trace struct {
	Process *schema.Process
	Trace   tracing.ITrace
}

func (t Trace) Unwrap() tracing.ITrace {
	return t.Trace
}

func (t Trace) TraceInterface() {}
