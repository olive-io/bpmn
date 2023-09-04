// Copyright 2023 Lack (xingyys@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flow_node

import (
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/sequence_flow"
)

type IAction interface {
	action()
}

type ProbeAction struct {
	SequenceFlows []*sequence_flow.SequenceFlow
	// ProbeReport is a function that needs to be called
	// wth sequence flow indices that have successful
	// condition expressions (or none)
	ProbeReport func([]int)
}

func (action ProbeAction) action() {}

type ActionTransformer func(sequenceFlowId *schema.IdRef, action IAction) IAction
type Terminate func(sequenceFlowId *schema.IdRef) chan bool

type FlowAction struct {
	SequenceFlows []*sequence_flow.SequenceFlow
	// Index of sequence flows that should flow without
	// conditionExpression being evaluated
	UnconditionalFlows []int
	// The actions produced by the targets should be processed by
	// this function
	ActionTransformer
	// If supplied channel sends a function that returns true, the flow action
	// is to be terminated if it wasn't already
	Terminate
}

func (action FlowAction) action() {}

type CompleteAction struct{}

func (action CompleteAction) action() {}

type NoAction struct{}

func (action NoAction) action() {}
