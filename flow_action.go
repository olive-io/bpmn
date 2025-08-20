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
	"github.com/olive-io/bpmn/v2/pkg/data"
)

type IAction interface {
	action()
}

type probeAction struct {
	sequenceFlows []*SequenceFlow
	// probeReport is a function that needs to be called
	// wth sequence flow indices that have successful
	// condition expressions (or none)
	probeReport func([]int)
}

func (action probeAction) action() {}

type ActionTransformer func(sequenceFlowId *schema.IdRef, action IAction) IAction
type Terminate func(sequenceFlowId *schema.IdRef) chan bool

type FlowActionResponse struct {
	dataObjects map[string]data.IItem
	variables   map[string]data.IItem
	err         error
	handler     <-chan ErrHandler
}

type ErrHandleMode int

const (
	RetryMode ErrHandleMode = iota + 1
	SkipMode
	ExitMode
)

type ErrHandler struct {
	Mode    ErrHandleMode
	Retries int32
}

type flowAction struct {
	response *FlowActionResponse

	sequenceFlows []*SequenceFlow
	// Index of sequence flows that should flow without
	// conditionExpression being evaluated
	unconditionalFlows []int
	// The actions produced by the targets should be produced by
	// this function
	actionTransformer ActionTransformer
	// If a supplied channel sends a function that returns true, the flow action
	// is to be terminated if it wasn't already
	terminate Terminate
}

func (action flowAction) action() {}

type completeAction struct{}

func (action completeAction) action() {}

type noAction struct{}

func (action noAction) action() {}
