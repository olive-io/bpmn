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

package parallel

import (
	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/schema"
)

// IncomingFlowProcessedTrace signals that a particular flow
// has been processed. If any action have been taken, it already happened
type IncomingFlowProcessedTrace struct {
	Node *schema.ParallelGateway
	Flow flow_interface.T
}

func (t IncomingFlowProcessedTrace) TraceInterface() {}
