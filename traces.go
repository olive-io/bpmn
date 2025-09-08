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
	Node   schema.FlowNodeInterface
}

func (t CancellationFlowTrace) Unpack() any { return t.Node }

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

type CeaseProcessSetTrace struct {
	Definitions *schema.Definitions
}

func (t CeaseProcessSetTrace) Unpack() any { return t.Definitions }
