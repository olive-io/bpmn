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

func (t ErrorTrace) Unpack() any { return t.Error }

type NewFlowTrace struct {
	FlowId id.Id
}

func (t NewFlowTrace) Unpack() any { return t.FlowId }

type FlowTrace struct {
	Source schema.FlowNodeInterface
	Flows  []Snapshot
}

func (t FlowTrace) Unpack() any { return t.Source }

type TerminationTrace struct {
	FlowId id.Id
	Source schema.FlowNodeInterface
}

func (t TerminationTrace) Unpack() any { return t.Source }

type CancellationFlowTrace struct {
	FlowId id.Id
}

func (t CancellationFlowTrace) Unpack() any { return t.FlowId }

type CompletionTrace struct {
	Node schema.FlowNodeInterface
}

func (t CompletionTrace) Unpack() any { return t.Node }

type CeaseFlowTrace struct {
	Process schema.Element
}

func (t CeaseFlowTrace) Unpack() any { return t.Process }

type VisitTrace struct {
	Node schema.FlowNodeInterface
}

func (t VisitTrace) Unpack() any { return t.Node }

type LeaveTrace struct {
	Node schema.FlowNodeInterface
}

func (t LeaveTrace) Unpack() any { return t.Node }

type CancellationFlowNodeTrace struct {
	Node schema.FlowNodeInterface
}

func (t CancellationFlowNodeTrace) Unpack() any { return t.Node }

type NewFlowNodeTrace struct {
	Node schema.FlowNodeInterface
}

func (t NewFlowNodeTrace) Unpack() any { return t.Node }

// ProcessTrace wraps any trace within a given process
type ProcessTrace struct {
	Process *schema.Process
	Trace   tracing.ITrace
}

func (t ProcessTrace) Unwrap() tracing.ITrace {
	return t.Trace
}

func (t ProcessTrace) Unpack() any { return t.Process }

// InstantiationTrace denotes instantiation of a given process
type InstantiationTrace struct {
	InstanceId id.Id
}

func (i InstantiationTrace) Unpack() any { return i.InstanceId }

// InstanceTrace wraps any trace with process instance id
type InstanceTrace struct {
	InstanceId id.Id
	Trace      tracing.ITrace
}

func (t InstanceTrace) Unwrap() tracing.ITrace {
	return t.Trace
}

func (t InstanceTrace) Unpack() any { return t.InstanceId }
