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

package bpmn

import (
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/id"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type ErrorTrace struct {
	Error error
}

func (t ErrorTrace) TraceInterface() {}

type WarningTrace struct {
	Warning any
}

func (t WarningTrace) TraceInterface() {}

type NewFlowTrace struct {
	FlowId id.Id
}

func (t NewFlowTrace) TraceInterface() {}

type FlowTrace struct {
	Source schema.FlowNodeInterface
	Flows  []Snapshot
}

func (t FlowTrace) TraceInterface() {}

type TerminationTrace struct {
	FlowId id.Id
	Source schema.FlowNodeInterface
}

func (t TerminationTrace) TraceInterface() {}

type CancellationTrace struct {
	FlowId id.Id
}

func (t CancellationTrace) TraceInterface() {}

type CompletionTrace struct {
	Node schema.FlowNodeInterface
}

func (t CompletionTrace) TraceInterface() {}

type CeaseFlowTrace struct{}

func (t CeaseFlowTrace) TraceInterface() {}

type VisitTrace struct {
	Node schema.FlowNodeInterface
}

func (t VisitTrace) TraceInterface() {}

type LeaveTrace struct {
	Node schema.FlowNodeInterface
}

func (t LeaveTrace) TraceInterface() {}

type CancellationFlowNodeTrace struct {
	Node schema.FlowNodeInterface
}

func (t CancellationFlowNodeTrace) TraceInterface() {}

type NewFlowNodeTrace struct {
	Node schema.FlowNodeInterface
}

func (t NewFlowNodeTrace) TraceInterface() {}

// ProcessTrace wraps any trace within a given process
type ProcessTrace struct {
	Process *schema.Process
	Trace   tracing.ITrace
}

func (t ProcessTrace) Unwrap() tracing.ITrace {
	return t.Trace
}

func (t ProcessTrace) TraceInterface() {}

// InstantiationTrace denotes instantiation of a given process
type InstantiationTrace struct {
	InstanceId id.Id
}

func (i InstantiationTrace) TraceInterface() {}

// InstanceTrace wraps any trace with process instance id
type InstanceTrace struct {
	InstanceId id.Id
	Trace      tracing.ITrace
}

func (t InstanceTrace) Unwrap() tracing.ITrace {
	return t.Trace
}

func (t InstanceTrace) TraceInterface() {}
