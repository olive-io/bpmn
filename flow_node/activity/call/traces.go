// Copyright 2023 The olive Authors
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

package call

import (
	"context"

	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/flow_node/activity"
	"github.com/olive-io/bpmn/schema"
)

type DoOption func(*doResponse)

func WithObjects(dataObjects map[string]any) DoOption {
	return func(rsp *doResponse) {
		rsp.dataObjects = dataObjects
	}
}

func WithProperties(properties map[string]any) DoOption {
	return func(rsp *doResponse) {
		rsp.properties = properties
	}
}

func WithErrHandle(err error, ch <-chan flow_node.ErrHandler) DoOption {
	return func(rsp *doResponse) {
		rsp.err = err
		rsp.handlerCh = ch
	}
}

func WithErr(err error) DoOption {
	return func(rsp *doResponse) {
		rsp.err = err
		rsp.handlerCh = nil
	}
}

type doResponse struct {
	dataObjects map[string]any
	properties  map[string]any
	err         error
	handlerCh   <-chan flow_node.ErrHandler
}

type ActiveTrace struct {
	Context       context.Context
	Activity      activity.Activity
	CalledElement *schema.ExtensionCalledElement
	response      chan doResponse
}

func (t *ActiveTrace) TraceInterface() {}

func (t *ActiveTrace) Do(options ...DoOption) {
	var response doResponse
	for _, opt := range options {
		opt(&response)
	}
	t.response <- response
}

func (t *ActiveTrace) GetActivity() activity.Activity {
	return t.Activity
}

func (t *ActiveTrace) Execute() {
	t.response <- doResponse{}
}
