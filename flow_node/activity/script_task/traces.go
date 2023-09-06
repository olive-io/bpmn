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

package script_task

import (
	"context"

	"github.com/olive-io/bpmn/schema"
)

type submitResponse struct {
	variables map[string]any
	err       error
}

type ScriptExecTrace struct {
	context.Context
	Element     schema.FlowNodeInterface
	DataObjects map[string]any
	Headers     map[string]any
	Properties  map[string]any
	response    chan submitResponse
}

func (t *ScriptExecTrace) TraceInterface() {}

func (t *ScriptExecTrace) Execute(variables map[string]any, err error) {
	t.response <- submitResponse{variables: variables, err: err}
}
