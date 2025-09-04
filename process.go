/*
Copyright 2023 The bpmn Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bpmn

import (
	"context"
	"fmt"
	"sync"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/data"
	"github.com/olive-io/bpmn/v2/pkg/errors"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/id"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type Options struct {
	ctx                            context.Context
	idGenerator                    id.IGenerator
	locator                        data.IFlowDataLocator
	eventIngress                   event.IConsumer
	eventEgress                    event.ISource
	tracer                         tracing.ITracer
	eventDefinitionInstanceBuilder event.IDefinitionInstanceBuilder
}

func NewOptions(opts ...Option) *Options {
	var options Options

	for _, opt := range opts {
		opt(&options)
	}

	if options.ctx == nil {
		options.ctx = context.Background()
	}

	if options.eventIngress == nil {
		options.eventIngress = event.NewFanOut()
	}

	if options.eventEgress == nil {
		options.eventEgress = event.NewFanOut()
	}

	if options.tracer == nil {
		options.tracer = tracing.NewTracer(options.ctx)
	}

	if options.locator == nil {
		options.locator = data.NewFlowDataLocator()
	}

	if options.idGenerator == nil {
		var err error
		options.idGenerator, err = id.GetSno().NewIdGenerator(options.ctx, options.tracer)
		if err != nil {
			panic(err)
		}
	}
	if options.eventDefinitionInstanceBuilder == nil {
		options.eventDefinitionInstanceBuilder = event.WrappingDefinitionInstanceBuilder
	}

	return &options
}

// Option allows to modify configuration of
// an instance in a flexible fashion (as it's just a modification
// function)
type Option func(*Options)

// WithTracer overrides instance's tracer
func WithTracer(tracer tracing.ITracer) Option {
	return func(opt *Options) {
		opt.tracer = tracer
	}
}

// WithContext will pass a given context to a new instance
// instead of implicitly generated one
func WithContext(ctx context.Context) Option {
	return func(opt *Options) {
		opt.ctx = ctx
	}
}

func WithEventIngress(consumer event.IConsumer) Option {
	return func(opt *Options) {
		opt.eventIngress = consumer
	}
}

func WithEventEgress(source event.ISource) Option {
	return func(opt *Options) {
		opt.eventEgress = source
	}
}

func WithIdGenerator(idGenerator id.IGenerator) Option {
	return func(opt *Options) {
		opt.idGenerator = idGenerator
	}
}

func WithLocator(locator data.IFlowDataLocator) Option {
	return func(opt *Options) {
		opt.locator = locator
	}
}

func WithVariables(variables map[string]any) Option {
	return func(opt *Options) {
		if opt.locator == nil {
			opt.locator = data.NewFlowDataLocator()
		}
		for key, value := range variables {
			opt.locator.SetVariable(key, value)
		}
	}
}

func WithProcessEventDefinitionInstanceBuilder(eventDefinitionInstanceBuilder event.IDefinitionInstanceBuilder) Option {
	return func(o *Options) {
		o.eventDefinitionInstanceBuilder = eventDefinitionInstanceBuilder
	}
}

func WithDataObjects(dataObjects map[string]any) Option {
	return func(opt *Options) {
		if opt.locator == nil {
			opt.locator = data.NewFlowDataLocator()
		}
		for dataObjectId, dataObject := range dataObjects {
			locator, found := opt.locator.FindIItemAwareLocator(data.LocatorObject)
			if !found {
				locator = data.NewDataObjectContainer()
				opt.locator.PutIItemAwareLocator(data.LocatorObject, locator)
			}
			container := data.NewContainer(nil)
			container.Put(dataObject)
			locator.PutItemAwareById(dataObjectId, container)
		}
	}
}

type Process struct {
	*Options

	id                 id.Id
	element            *schema.Process
	flowNodeMapping    *FlowNodeMapping
	flowWaitGroup      sync.WaitGroup
	complete           sync.RWMutex
	eventConsumersLock sync.RWMutex
	eventConsumers     []event.IConsumer
}

func (p *Process) Id() id.Id { return p.id }

func (p *Process) Element() *schema.Process { return p.element }

func (p *Process) Tracer() tracing.ITracer { return p.tracer }

func (p *Process) Locator() data.IFlowDataLocator { return p.locator }

func (p *Process) ConsumeEvent(ev event.IEvent) (result event.ConsumptionResult, err error) {
	p.eventConsumersLock.RLock()
	// We're copying the list of consumers here to ensure that
	// new consumers can subscribe during event forwarding
	eventConsumers := p.eventConsumers
	p.eventConsumersLock.RUnlock()
	result, err = event.ForwardEvent(ev, &eventConsumers)
	return
}

func (p *Process) RegisterEventConsumer(ev event.IConsumer) (err error) {
	p.eventConsumersLock.Lock()
	defer p.eventConsumersLock.Unlock()
	p.eventConsumers = append(p.eventConsumers, ev)
	return
}

func (p *Process) FlowNodeMapping() *FlowNodeMapping { return p.flowNodeMapping }

func NewProcess(processElem *schema.Process, definitions *schema.Definitions, opts ...Option) (process *Process, err error) {
	options := NewOptions(opts...)

	process = &Process{
		Options:         options,
		element:         processElem,
		flowNodeMapping: NewLockedFlowNodeMapping(),
	}
	idGenerator := options.idGenerator

	process.id = idGenerator.New()

	err = process.eventEgress.RegisterEventConsumer(process)
	if err != nil {
		return
	}

	err = data.ElementToLocator(process.locator, idGenerator, processElem)
	if err != nil {
		return
	}

	ctx := process.ctx

	// flow nodes
	subTracer := tracing.NewTracer(ctx)

	tracing.NewRelay(ctx, subTracer, process.tracer, func(trace tracing.ITrace) []tracing.ITrace {
		return []tracing.ITrace{InstanceTrace{
			InstanceId: process.id,
			Trace:      trace,
		}}
	})

	wiringMaker := func(element *schema.FlowNode) (*wiring, error) {
		return newWiring(
			process.id,
			processElem,
			definitions,
			element,
			// Event ingress/egress orchestration:
			//
			// flow nodes will send their message to `instance.EventIngress`
			// (which is typically the model), but consume their messages from
			// `instance`, which is turn a subscriber of `instance.EventEgress`
			// (again, typically, the model).
			//
			// This allows us to use ConsumeEvent on this instance to send
			// events only to the instance (useful for things like event-based
			// process instantiation)
			process.eventIngress, process,
			subTracer,
			process.flowNodeMapping,
			&process.flowWaitGroup,
			process.eventDefinitionInstanceBuilder,
			process.locator)
	}

	var wr *wiring

	for i := range *processElem.StartEvents() {
		element := &(*processElem.StartEvents())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var node *startEvent
		node, err = newStartEvent(wr, element, idGenerator)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, node)
		if err != nil {
			return
		}
	}

	for i := range *processElem.EndEvents() {
		element := &(*processElem.EndEvents())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var node *endEvent
		node, err = newEndEvent(wr, element)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, node)
		if err != nil {
			return
		}
	}

	for i := range *processElem.IntermediateCatchEvents() {
		element := &(*processElem.IntermediateCatchEvents())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var intermediateCatchEvent *catchEvent
		intermediateCatchEvent, err = newCatchEvent(wr, &element.CatchEvent)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, intermediateCatchEvent)
		if err != nil {
			return
		}
	}

	for i := range *processElem.IntermediateThrowEvents() {
		element := &(*processElem.IntermediateThrowEvents())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var intermediateThrowEvent *throwEvent
		intermediateThrowEvent, err = newThrowEvent(wr, &element.ThrowEvent)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, intermediateThrowEvent)
		if err != nil {
			return
		}
	}

	for i := range *processElem.Tasks() {
		element := &(*processElem.Tasks())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var node *harness
		task := newTask(element, TaskActivity)
		node, err = newHarness(wr, idGenerator, task)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, node)
		if err != nil {
			return
		}
	}

	for i := range *processElem.BusinessRuleTasks() {
		element := &(*processElem.BusinessRuleTasks())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var node *harness
		businessRule := newTask(element, BusinessRuleActivity)
		node, err = newHarness(wr, idGenerator, businessRule)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, node)
		if err != nil {
			return
		}
	}

	for i := range *processElem.CallActivities() {
		element := &(*processElem.CallActivities())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var node *harness
		callAct := newTask(element, CallActivity)
		node, err = newHarness(wr, idGenerator, callAct)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, node)
		if err != nil {
			return
		}
	}

	for i := range *processElem.ManualTasks() {
		element := &(*processElem.ManualTasks())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var node *harness
		manualTask := newTask(element, ManualTaskActivity)
		node, err = newHarness(wr, idGenerator, manualTask)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, node)
		if err != nil {
			return
		}
	}

	for i := range *processElem.ServiceTasks() {
		element := &(*processElem.ServiceTasks())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var node *harness
		serviceTask := newTask(element, ServiceTaskActivity)
		node, err = newHarness(wr, idGenerator, serviceTask)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, node)
		if err != nil {
			return
		}
	}

	for i := range *processElem.UserTasks() {
		element := &(*processElem.UserTasks())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var node *harness
		userTask := newTask(element, UserTaskActivity)
		node, err = newHarness(wr, idGenerator, userTask)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, node)
		if err != nil {
			return
		}
	}

	for i := range *processElem.ReceiveTasks() {
		element := &(*processElem.ReceiveTasks())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var node *harness
		receiveTask := newTask(element, ReceiveTaskActivity)
		node, err = newHarness(wr, idGenerator, receiveTask)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, node)
		if err != nil {
			return
		}
	}

	for i := range *processElem.ScriptTasks() {
		element := &(*processElem.ScriptTasks())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var node *harness
		scriptTask := newTask(element, ScriptTaskActivity)
		node, err = newHarness(wr, idGenerator, scriptTask)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, node)
		if err != nil {
			return
		}
	}

	for i := range *processElem.SendTasks() {
		element := &(*processElem.SendTasks())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var node *harness
		sendTask := newTask(element, SendTaskActivity)
		node, err = newHarness(wr, idGenerator, sendTask)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, node)
		if err != nil {
			return
		}
	}

	for i := range *processElem.SubProcesses() {
		element := &(*processElem.SubProcesses())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var node *harness
		subProcess := newSubProcess(process.eventDefinitionInstanceBuilder, idGenerator, element)
		node, err = newHarness(wr, idGenerator, subProcess)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, node)
		if err != nil {
			return
		}
	}

	for i := range *processElem.ExclusiveGateways() {
		element := &(*processElem.ExclusiveGateways())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var node *exclusiveGateway
		node, err = newExclusiveGateway(wr, element)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, node)
		if err != nil {
			return
		}
	}

	for i := range *processElem.InclusiveGateways() {
		element := &(*processElem.InclusiveGateways())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var node *inclusiveGateway
		node, err = newInclusiveGateway(wr, element)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, node)
		if err != nil {
			return
		}
	}

	for i := range *processElem.ParallelGateways() {
		element := &(*processElem.ParallelGateways())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var node *parallelGateway
		node, err = newParallelGateway(wr, element)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, node)
		if err != nil {
			return
		}
	}

	for i := range *processElem.EventBasedGateways() {
		element := &(*processElem.EventBasedGateways())[i]
		wr, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var node *eventBasedGateway
		node, err = newEventBasedGateway(wr, element)
		if err != nil {
			return
		}
		err = process.flowNodeMapping.RegisterElementToFlowNode(element, node)
		if err != nil {
			return
		}
	}

	process.flowNodeMapping.Finalize()

	// StartAll cease flow monitor
	sender := process.tracer.RegisterSender()
	go process.ceaseFlowMonitor(subTracer)(ctx, sender)
	process.tracer.Send(InstantiationTrace{InstanceId: process.id})

	return
}

// StartWith explicitly starts the instance by triggering a given start event
func (p *Process) StartWith(ctx context.Context, element schema.StartEventInterface) (err error) {
	flowNode, found := p.flowNodeMapping.ResolveElementToFlowNode(element)
	elementId := "<unnamed>"
	if idPtr, present := element.Id(); present {
		elementId = *idPtr
	}
	processId := "<unnamed>"
	if idPtr, present := p.element.Id(); present {
		processId = *idPtr
	}
	if !found {
		err = errors.NotFoundError{Expected: fmt.Sprintf("start event %s in process %s", elementId, processId)}
		return
	}
	startEventNode, ok := flowNode.(*startEvent)
	if !ok {
		err = errors.RequirementExpectationError{
			Expected: fmt.Sprintf("start event %s flow node in process %s to be of type start.Node", elementId, processId),
			Actual:   fmt.Sprintf("%sFlow", flowNode),
		}
		return
	}
	startEventNode.Trigger(ctx)
	return
}

// StartAll explicitly starts the instance by triggering all start events, if any
func (p *Process) StartAll(ctx context.Context) (err error) {
	for i := range *p.element.StartEvents() {
		err = p.StartWith(ctx, &(*p.element.StartEvents())[i])
		if err != nil {
			return
		}
	}
	return
}

func (p *Process) ceaseFlowMonitor(tracer tracing.ITracer) func(ctx context.Context, sender tracing.ISenderHandle) {
	// Subscribing to traces early as otherwise events produced
	// after the goroutine below is started are not going to be
	// sent to it.
	traces := tracer.Subscribe()
	p.complete.Lock()
	return func(ctx context.Context, sender tracing.ISenderHandle) {
		defer sender.Done()
		defer p.complete.Unlock()

		/* 13.4.6 End Events:

		The Process instance is [...] completed, if
		and only if the following two conditions
		hold:

		(1) All start nodes of the Process have been
		visited. More precisely, all StartAll Events
		have been triggered (1.1), and for all
		starting Event-Based Gateways, one of the
		associated Events has been triggered (1.2).

		(2) There is no token remaining within the
		Process instance
		*/
		startEventsActivated := make([]*schema.StartEvent, 0)

		// So, at first, we wait for (1.1) to occur
		// [(1.2) will be addded when we actually support them]

		for {
			if len(startEventsActivated) == len(*p.element.StartEvents()) {
				break
			}

			select {
			case trace := <-traces:
				trace = tracing.Unwrap(trace)
				switch t := trace.(type) {
				case TerminationTrace:
					switch flowNode := t.Source.(type) {
					case *schema.StartEvent:
						startEventsActivated = append(startEventsActivated, flowNode)
					default:
					}
				case FlowTrace:
					switch flowNode := t.Source.(type) {
					case *schema.StartEvent:
						startEventsActivated = append(startEventsActivated, flowNode)
					default:
					}
				default:
				}
			case <-ctx.Done():
				tracer.Unsubscribe(traces)
				return
			}
		}

		tracer.Unsubscribe(traces)

		// Then, we're waiting for (2) to occur
		waitIsOver := make(chan struct{})
		go func() {
			p.flowWaitGroup.Wait()
			close(waitIsOver)
		}()
		select {
		case <-waitIsOver:
			// Send out a cease flow trace
			tracer.Send(CeaseFlowTrace{Process: p.element})
		case <-ctx.Done():
		}
	}
}

// WaitUntilComplete waits until the instance is complete.
// Returns true if the instance was complete, false if the context signaled `Done`
func (p *Process) WaitUntilComplete(ctx context.Context) (complete bool) {
	signal := make(chan bool)
	go func() {
		p.complete.Lock()
		defer p.complete.Unlock()
		signal <- true
	}()
	select {
	case <-ctx.Done():
		complete = false
	case <-signal:
		complete = true
	}
	return
}
