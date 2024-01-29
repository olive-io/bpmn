// Copyright 2023 The bpmn Authors
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
