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

package activity

import "github.com/olive-io/bpmn/flow_node"

type Type string

const (
	TaskType         Type = "Task"
	ServiceType      Type = "ServiceTask"
	ScriptType       Type = "ScriptTask"
	UserType         Type = "UserTask"
	CallType         Type = "CallActivity"
	BusinessRuleType Type = "BusinessRuleTask"
	SendType         Type = "SendTask"
	ReceiveType      Type = "ReceiveTask"
	SubprocessType   Type = "Subprocess"
)

// Activity is a generic interface to flow nodes that are activities
type Activity interface {
	flow_node.IFlowNode
	Type() Type
	// Cancel initiates a cancellation of activity and returns a channel
	// that will signal a boolean (`true` if cancellation was successful,
	// `false` otherwise)
	Cancel() <-chan bool
}
