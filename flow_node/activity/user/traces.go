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

package user

import (
	"context"

	"github.com/olive-io/bpmn/flow_node/activity"
)

type DoOption func(*doResponse)

func WithProperties(properties map[string]any) DoOption {
	return func(rsp *doResponse) {
		rsp.properties = properties
	}
}

func WithErr(err error) DoOption {
	return func(rsp *doResponse) {
		rsp.err = err
	}
}

type doResponse struct {
	properties map[string]any
	err        error
}

type ActiveTrace struct {
	context.Context
	Activity   activity.Activity
	Headers    map[string]any
	Properties map[string]any
	response   chan doResponse
}

func (t *ActiveTrace) TraceInterface() {}

func (t *ActiveTrace) Do(options ...DoOption) {
	var response doResponse
	for _, opt := range options {
		opt(&response)
	}
	t.response <- response
}

func (t *ActiveTrace) Execute() {
	t.response <- doResponse{}
}
