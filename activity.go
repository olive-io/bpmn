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
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/data"
	"github.com/olive-io/bpmn/v2/pkg/errors"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/id"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

var (
	DefaultTaskExecTimeout = time.Second * 30
)

type Type string

const (
	TaskType         Type = "Task"
	ServiceType      Type = "ServiceTask"
	ScriptType       Type = "ScriptTask"
	UserType         Type = "UserTask"
	ManualType       Type = "ManualTask"
	CallType         Type = "CallActivity"
	BusinessRuleType Type = "BusinessRuleTask"
	SendType         Type = "SendTask"
	ReceiveType      Type = "ReceiveTask"
	SubprocessType   Type = "Subprocess"
)

type TypeKey struct{}

// Activity is a generic interface to flow nodes that are activities
type Activity interface {
	IFlowNode
	Type() Type
	// Cancel initiates a cancellation of activity and returns a channel
	// that will signal a boolean (`true` if cancellation was successful,
	// `false` otherwise)
	Cancel() <-chan bool
}

type nextHarnessActionMessage struct {
	flow     Flow
	response chan chan IAction
}

func (m nextHarnessActionMessage) message() {}

type Harness struct {
	*Wiring
	mch                chan imessage
	activity           Activity
	active             int32
	cancellation       sync.Once
	eventConsumers     []event.IConsumer
	eventConsumersLock sync.RWMutex
}

func (node *Harness) ConsumeEvent(ev event.IEvent) (result event.ConsumptionResult, err error) {
	node.eventConsumersLock.RLock()
	defer node.eventConsumersLock.RUnlock()
	if atomic.LoadInt32(&node.active) == 1 {
		result, err = event.ForwardEvent(ev, &node.eventConsumers)
	}
	return
}

func (node *Harness) RegisterEventConsumer(consumer event.IConsumer) (err error) {
	node.eventConsumersLock.Lock()
	defer node.eventConsumersLock.Unlock()
	node.eventConsumers = append(node.eventConsumers, consumer)
	return
}

func (node *Harness) Activity() Activity { return node.activity }

type Constructor = func(*Wiring) (node Activity, err error)

func NewHarness(ctx context.Context, wiring *Wiring, idGenerator id.IGenerator, constructor Constructor) (node *Harness, err error) {
	var activity Activity
	activity, err = constructor(wiring)
	if err != nil {
		return
	}

	boundaryEvents := make([]*schema.BoundaryEvent, 0)

	switch process := wiring.Process.(type) {
	case *schema.Process:
		for i := range *process.BoundaryEvents() {
			boundaryEvent := &(*process.BoundaryEvents())[i]
			if string(*boundaryEvent.AttachedToRef()) == wiring.FlowNodeId {
				boundaryEvents = append(boundaryEvents, boundaryEvent)
			}
		}
	case *schema.SubProcess:
		for i := range *process.BoundaryEvents() {
			boundaryEvent := &(*process.BoundaryEvents())[i]
			if string(*boundaryEvent.AttachedToRef()) == wiring.FlowNodeId {
				boundaryEvents = append(boundaryEvents, boundaryEvent)
			}
		}
	}

	node = &Harness{
		Wiring:   wiring,
		mch:      make(chan imessage, len(wiring.Incoming)*2+1),
		activity: activity,
	}

	err = node.EventEgress.RegisterEventConsumer(node)
	if err != nil {
		return
	}

	for i := range boundaryEvents {
		boundaryEvent := boundaryEvents[i]
		var catchEventFlowNode *Wiring
		catchEventFlowNode, err = wiring.CloneFor(&boundaryEvent.FlowNode)
		if err != nil {
			return
		}
		// this node becomes event egress
		catchEventFlowNode.EventEgress = node

		var catchEvent *CatchEvent
		catchEvent, err = NewCatchEvent(ctx, catchEventFlowNode, &boundaryEvent.CatchEvent)
		if err != nil {
			return
		}

		var actionTransformer ActionTransformer
		if boundaryEvent.CancelActivity() {
			actionTransformer = func(sequenceFlowId *schema.IdRef, action IAction) IAction {
				node.cancellation.Do(func() {
					<-node.activity.Cancel()
				})
				return action
			}
		}
		flowable := newFlow(node.Definitions, catchEvent, node.Tracer, node.FlowNodeMapping, node.FlowWaitGroup, idGenerator, actionTransformer, node.Locator)
		flowable.Start(ctx)
	}
	sender := node.Tracer.RegisterSender()
	go node.run(ctx, sender)
	return
}

func (node *Harness) run(ctx context.Context, sender tracing.ISenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-node.mch:
			switch m := msg.(type) {
			case nextHarnessActionMessage:
				atomic.StoreInt32(&node.active, 1)
				node.Tracer.Trace(ActiveBoundaryTrace{Start: true, Node: node.activity.Element()})
				in := node.activity.NextAction(m.flow)
				out := make(chan IAction, 1)
				go func(ctx context.Context) {
					select {
					case out <- <-in:
						atomic.StoreInt32(&node.active, 0)
						node.Tracer.Trace(ActiveBoundaryTrace{Start: false, Node: node.activity.Element()})
					case <-ctx.Done():
						return
					}
				}(ctx)
				m.response <- out
			default:
			}
		case <-ctx.Done():
			node.Tracer.Trace(CancellationFlowNodeTrace{Node: node.activity.Element()})
			return
		}
	}
}

func (node *Harness) NextAction(flow Flow) chan IAction {
	response := make(chan chan IAction, 1)
	node.mch <- nextHarnessActionMessage{flow: flow, response: response}
	return <-response
}

func (node *Harness) Element() schema.FlowNodeInterface { return node.activity.Element() }

type ActiveBoundaryTrace struct {
	Start bool
	Node  schema.FlowNodeInterface
}

func (b ActiveBoundaryTrace) Element() any { return b.Node }

type DoOption func(*DoResponse)

func DoWithObjects(dataObjects map[string]any) DoOption {
	return func(rsp *DoResponse) {
		rsp.DataObjects = dataObjects
	}
}

func DoWithResults(results map[string]any) DoOption {
	return func(rsp *DoResponse) {
		rsp.Results = results
	}
}

func DoWithErrHandle(err error, ch <-chan ErrHandler) DoOption {
	return func(rsp *DoResponse) {
		rsp.Err = err
		rsp.HandlerCh = ch
	}
}

func DoWithValue(key, value any) DoOption {
	return func(rsp *DoResponse) {
		rsp.Context = context.WithValue(rsp.Context, key, value)
	}
}

func DoWithErr(err error) DoOption {
	return func(rsp *DoResponse) {
		rsp.Err = err
		rsp.HandlerCh = nil
	}
}

type DoResponse struct {
	Context     context.Context
	DataObjects map[string]any
	Results     map[string]any
	Err         error
	HandlerCh   <-chan ErrHandler
}

func newDoOption(opts ...DoOption) *DoResponse {
	var rsp DoResponse
	for _, opt := range opts {
		opt(&rsp)
	}

	return &rsp
}

type TaskTraceBuilder struct {
	t TaskTrace
}

func NewTaskTraceBuilder() *TaskTraceBuilder {
	trace := newTaskTrace()
	return &TaskTraceBuilder{t: *trace}
}

func (b *TaskTraceBuilder) Context(ctx context.Context) *TaskTraceBuilder {
	b.t.ctx = ctx
	return b
}

func (b *TaskTraceBuilder) Value(key, value any) *TaskTraceBuilder {
	b.t.ctx = context.WithValue(b.t.ctx, key, value)
	return b
}

func (b *TaskTraceBuilder) Timeout(timeout time.Duration) *TaskTraceBuilder {
	b.t.timeout = timeout
	return b
}

func (b *TaskTraceBuilder) Activity(activity Activity) *TaskTraceBuilder {
	b.t.activity = activity
	return b
}

func (b *TaskTraceBuilder) DataObjects(dataObjects map[string]any) *TaskTraceBuilder {
	b.t.dataObjects = dataObjects
	return b
}

func (b *TaskTraceBuilder) Headers(headers map[string]string) *TaskTraceBuilder {
	b.t.headers = headers
	return b
}

func (b *TaskTraceBuilder) Properties(properties map[string]any) *TaskTraceBuilder {
	b.t.properties = properties
	return b
}

func (b *TaskTraceBuilder) Build() *TaskTrace {
	go b.t.process()
	return &b.t
}

// TaskTrace describes common channel handler for all tasks
type TaskTrace struct {
	ctx         context.Context
	timeout     time.Duration
	activity    Activity
	headers     map[string]string
	properties  map[string]any
	dataObjects map[string]any
	forward     chan DoResponse
	response    chan DoResponse
	done        chan struct{}
}

func newTaskTrace() *TaskTrace {
	trace := TaskTrace{
		ctx:         context.TODO(),
		timeout:     DefaultTaskExecTimeout,
		headers:     make(map[string]string),
		properties:  make(map[string]any),
		dataObjects: make(map[string]any),
		forward:     make(chan DoResponse, 1),
		response:    make(chan DoResponse, 1),
		done:        make(chan struct{}, 1),
	}
	return &trace
}

func (t *TaskTrace) Element() any { return t.activity }

func (t *TaskTrace) Context() context.Context {
	return t.ctx
}

func (t *TaskTrace) GetActivity() Activity {
	return t.activity
}

func (t *TaskTrace) GetDataObjects() map[string]any {
	return t.dataObjects
}

func (t *TaskTrace) GetHeaders() map[string]string {
	return t.headers
}

func (t *TaskTrace) GetProperties() map[string]any {
	return t.properties
}

func (t *TaskTrace) out() <-chan DoResponse {
	return t.response
}

func (t *TaskTrace) Do(options ...DoOption) {
	select {
	case <-t.done:
		return
	default:
	}

	response := newDoOption(options...)
	t.forward <- *response
}

func (t *TaskTrace) process() {
	duration := t.timeout
	if duration < time.Second {
		duration = time.Second
	}

	select {
	case <-t.done:
		return
	case <-t.ctx.Done():
		rsp := newDoOption(DoWithErr(t.ctx.Err()))
		t.response <- *rsp
	case <-time.After(duration):
		var tid string
		if v, ok := t.activity.Element().Id(); ok {
			tid = *v
		}

		rsp := newDoOption(DoWithErr(errors.TaskExecError{Id: tid, Reason: "timed out"}))
		t.response <- *rsp
	case rsp := <-t.forward:
		t.response <- rsp
	}

	select {
	case <-t.done:
	default:
		close(t.done)
	}
}

func FetchTaskDataInput(locator data.IFlowDataLocator, element schema.BaseElementInterface) (headers map[string]string, properties, dataObjects map[string]any) {
	variables := locator.CloneVariables()
	headers = map[string]string{}
	properties = map[string]any{}
	dataObjects = map[string]any{}
	if extension, found := element.ExtensionElements(); found {
		if header := extension.TaskHeaderField; header != nil {
			fields := header.Header
			for _, field := range fields {
				headers[field.Name] = field.Value
			}
		}
		if property := extension.PropertiesField; property != nil {
			fields := property.Property
			for _, field := range fields {
				value := field.ValueFor()
				if len(strings.TrimSpace(field.Value)) == 0 {
					if vv, ok := variables[field.Name]; ok {
						value = vv
					}
				}
				properties[field.Name] = value
			}
		}

		awareLocator, found1 := locator.FindIItemAwareLocator(data.LocatorObject)
		if found1 {
			for _, dataInput := range extension.DataInput {
				aWare, ok := awareLocator.FindItemAwareById(dataInput.TargetRef)
				if ok {
					dataObjects[dataInput.Name] = aWare.Get()
				}
			}
		}
	}

	return
}

// FetchTaskTimeout returns timeout by schema.BaseElementInterface
func FetchTaskTimeout(element schema.BaseElementInterface) time.Duration {
	var timeout time.Duration
	if extension, found := element.ExtensionElements(); found {
		if field := extension.TaskDefinitionField; field != nil {
			timeout, _ = time.ParseDuration(field.Timeout)
		}
	}
	if timeout == 0 {
		timeout = DefaultTaskExecTimeout
	}
	return timeout
}

func ApplyTaskDataOutput(element schema.BaseElementInterface, dataOutputs map[string]any) map[string]data.IItem {
	outputs := map[string]data.IItem{}
	if extension, found := element.ExtensionElements(); found {
		for _, dataOutput := range extension.DataOutput {
			value, ok := dataOutputs[dataOutput.Name]
			if ok {
				outputs[dataOutput.Name] = value
			}
		}
	}
	return outputs
}

func ApplyTaskResult(element schema.BaseElementInterface, results map[string]any) map[string]data.IItem {
	outputs := map[string]data.IItem{}
	if extension, found := element.ExtensionElements(); found {
		if field := extension.ResultsField; field != nil {
			for _, item := range extension.ResultsField.Item {
				value, ok := results[item.Name]
				if ok {
					outputs[item.Name] = value
				}
			}
			return outputs
		}
	}

	for key, value := range results {
		outputs[key] = value
	}
	return outputs
}
