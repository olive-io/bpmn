/*
   Copyright 2023 The bpmn Authors

   This program is offered under a commercial and under the AGPL license.
   For AGPL licensing, see below.

   AGPL licensing:
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package activity

import (
	"context"

	"github.com/olive-io/bpmn/schema"

	"github.com/olive-io/bpmn/flow_node"
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
	select {
	case <-t.done:
		return
	default:
	}

	var response DoResponse
	for _, opt := range options {
		opt(&response)
	}

	t.response <- response
	close(t.done)
}
