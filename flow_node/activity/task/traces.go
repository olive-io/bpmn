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

package task

import (
	"context"

	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/flow_node/activity"
)

type ActiveTrace struct {
	Context  context.Context
	Activity activity.Activity
	response chan flow_node.IAction
}

func (t *ActiveTrace) TraceInterface() {}

func (t *ActiveTrace) Do(action flow_node.IAction) {
	t.response <- action
}

func (t *ActiveTrace) Execute() {
	node := t.Activity.(*Task)
	t.response <- flow_node.FlowAction{SequenceFlows: flow_node.AllSequenceFlows(&node.Wiring.Outgoing)}
}
