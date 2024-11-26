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
	"github.com/olive-io/bpmn/v2/pkg/data"
)

type IAction interface {
	action()
}

type ProbeAction struct {
	SequenceFlows []*SequenceFlow
	// ProbeReport is a function that needs to be called
	// wth sequence flow indices that have successful
	// condition expressions (or none)
	ProbeReport func([]int)
}

func (action ProbeAction) action() {}

type ActionTransformer func(sequenceFlowId *schema.IdRef, action IAction) IAction
type Terminate func(sequenceFlowId *schema.IdRef) chan bool

type FlowActionResponse struct {
	DataObjects map[string]data.IItem
	Variables   map[string]data.IItem
	Err         error
	Handler     <-chan ErrHandler
}

type ErrHandleMode int

const (
	HandleRetry ErrHandleMode = iota + 1
	HandleSkip
	HandleExit
)

type ErrHandler struct {
	Mode    ErrHandleMode
	Retries int32
}

type FlowAction struct {
	Response *FlowActionResponse

	SequenceFlows []*SequenceFlow
	// Index of sequence flows that should flow without
	// conditionExpression being evaluated
	UnconditionalFlows []int
	// The actions produced by the targets should be processed by
	// this function
	ActionTransformer ActionTransformer
	// If supplied channel sends a function that returns true, the flow action
	// is to be terminated if it wasn't already
	Terminate Terminate
}

func (action FlowAction) action() {}

type CompleteAction struct{}

func (action CompleteAction) action() {}

type NoAction struct{}

func (action NoAction) action() {}
