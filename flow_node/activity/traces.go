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

package activity

import (
	"context"

	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/schema"
)

type ActiveBoundaryTrace struct {
	Start bool
	Node  schema.FlowNodeInterface
}

func (b ActiveBoundaryTrace) TraceInterface() {}

type DoOption func(*DoResponse)

func WithObjects(dataObjects map[string]any) DoOption {
	return func(rsp *DoResponse) {
		rsp.DataObjects = dataObjects
	}
}

func WithProperties(properties map[string]any) DoOption {
	return func(rsp *DoResponse) {
		rsp.Properties = properties
	}
}

func WithErrHandle(err error, ch <-chan flow_node.ErrHandler) DoOption {
	return func(rsp *DoResponse) {
		rsp.Err = err
		rsp.HandlerCh = ch
	}
}

func WithValue(key, value any) DoOption {
	return func(rsp *DoResponse) {
		rsp.Context = context.WithValue(rsp.Context, key, value)
	}
}

func WithErr(err error) DoOption {
	return func(rsp *DoResponse) {
		rsp.Err = err
		rsp.HandlerCh = nil
	}
}

type DoResponse struct {
	Context     context.Context
	DataObjects map[string]any
	Properties  map[string]any
	Err         error
	HandlerCh   <-chan flow_node.ErrHandler
}

type TraceBuilder struct {
	t Trace
}

func NewTraceBuilder() *TraceBuilder {
	trace := Trace{
		ctx:         context.TODO(),
		dataObjects: make(map[string]any),
		headers:     make(map[string]any),
		properties:  make(map[string]any),
		response:    make(chan DoResponse, 1),
		done:        make(chan struct{}, 1),
	}
	return &TraceBuilder{t: trace}
}

func (b *TraceBuilder) Context(ctx context.Context) *TraceBuilder {
	b.t.ctx = ctx
	return b
}

func (b *TraceBuilder) Value(key, value any) *TraceBuilder {
	b.t.ctx = context.WithValue(b.t.ctx, key, value)
	return b
}

func (b *TraceBuilder) Activity(activity Activity) *TraceBuilder {
	b.t.activity = activity
	return b
}

func (b *TraceBuilder) DataObjects(dataObjects map[string]any) *TraceBuilder {
	b.t.dataObjects = dataObjects
	return b
}

func (b *TraceBuilder) Headers(headers map[string]any) *TraceBuilder {
	b.t.headers = headers
	return b
}

func (b *TraceBuilder) Properties(properties map[string]any) *TraceBuilder {
	b.t.properties = properties
	return b
}

func (b *TraceBuilder) Response(ch chan DoResponse) *TraceBuilder {
	b.t.response = ch
	return b
}

func (b *TraceBuilder) Build() *Trace {
	return &b.t
}

// Trace describes common channel handler for all tasks
type Trace struct {
	ctx         context.Context
	activity    Activity
	dataObjects map[string]any
	headers     map[string]any
	properties  map[string]any
	response    chan DoResponse
	done        chan struct{}
}

func (t *Trace) TraceInterface() {}

func (t *Trace) Context() context.Context {
	return t.ctx
}

func (t *Trace) GetActivity() Activity {
	return t.activity
}

func (t *Trace) GetDataObjects() map[string]any {
	return t.dataObjects
}

func (t *Trace) GetHeaders() map[string]any {
	return t.headers
}

func (t *Trace) GetProperties() map[string]any {
	return t.properties
}

func (t *Trace) Do(options ...DoOption) {
	var response DoResponse
	for _, opt := range options {
		opt(&response)
	}
	select {
	case <-t.done:
	default:
		t.response <- response
		close(t.done)
	}
}
