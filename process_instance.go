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
	"fmt"
	"sync"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/data"
	"github.com/olive-io/bpmn/v2/pkg/errors"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/id"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type Instance struct {
	*Options

	id                 id.Id
	process            *schema.Process
	flowNodeMapping    *FlowNodeMapping
	flowWaitGroup      sync.WaitGroup
	complete           sync.RWMutex
	eventConsumersLock sync.RWMutex
	eventConsumers     []event.IConsumer
}

func (ins *Instance) Id() id.Id {
	return ins.id
}

func (ins *Instance) Process() *schema.Process { return ins.process }

func (ins *Instance) Tracer() tracing.ITracer { return ins.tracer }

func (ins *Instance) Locator() data.IFlowDataLocator { return ins.locator }

func (ins *Instance) ConsumeEvent(ev event.IEvent) (result event.ConsumptionResult, err error) {
	ins.eventConsumersLock.RLock()
	// We're copying the list of consumers here to ensure that
	// new consumers can subscribe during event forwarding
	eventConsumers := ins.eventConsumers
	ins.eventConsumersLock.RUnlock()
	result, err = event.ForwardEvent(ev, &eventConsumers)
	return
}

func (ins *Instance) RegisterEventConsumer(ev event.IConsumer) (err error) {
	ins.eventConsumersLock.Lock()
	defer ins.eventConsumersLock.Unlock()
	ins.eventConsumers = append(ins.eventConsumers, ev)
	return
}

func (ins *Instance) FlowNodeMapping() *FlowNodeMapping {
	return ins.flowNodeMapping
}

func NewInstance(process *schema.Process, definitions *schema.Definitions, options *Options) (instance *Instance, err error) {
	instance = &Instance{
		Options:         options,
		process:         process,
		flowNodeMapping: NewLockedFlowNodeMapping(),
	}

	if instance.ctx == nil {
		instance.ctx = context.Background()
	}
	ctx := instance.ctx

	if instance.tracer == nil {
		instance.tracer = tracing.NewTracer(ctx)
	}

	if instance.idGeneratorBuilder == nil {
		instance.idGeneratorBuilder = id.DefaultIdGeneratorBuilder
	}

	if instance.locator == nil {
		instance.locator = data.NewFlowDataLocator()
	}

	var idGenerator id.IGenerator
	idGenerator, err = instance.idGeneratorBuilder.NewIdGenerator(ctx, instance.tracer)
	if err != nil {
		return
	}

	instance.idGenerator = idGenerator

	instance.id = idGenerator.New()

	err = instance.eventEgress.RegisterEventConsumer(instance)
	if err != nil {
		return
	}

	err = data.ElementToLocator(instance.locator, idGenerator, instance.process)
	if err != nil {
		return
	}

	// flow nodes
	subTracer := tracing.NewTracer(ctx)

	tracing.NewRelay(ctx, subTracer, instance.tracer, func(trace tracing.ITrace) []tracing.ITrace {
		return []tracing.ITrace{InstanceTrace{
			InstanceId: instance.id,
			Trace:      trace,
		}}
	})

	wiringMaker := func(element *schema.FlowNode) (*Wiring, error) {
		return NewWiring(
			instance.id,
			instance.process,
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
			instance.eventIngress, instance,
			subTracer,
			instance.flowNodeMapping,
			&instance.flowWaitGroup, instance.eventDefinitionInstanceBuilder,
			instance.locator)
	}

	var wiring *Wiring

	for i := range *instance.process.StartEvents() {
		element := &(*instance.process.StartEvents())[i]
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var startEvent *StartEvent
		startEvent, err = NewStartEvent(ctx, wiring, element, idGenerator)
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, startEvent)
		if err != nil {
			return
		}
	}

	for i := range *instance.process.EndEvents() {
		element := &(*instance.process.EndEvents())[i]
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var endEvent *EndEvent
		endEvent, err = NewEndEvent(ctx, wiring, element)
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, endEvent)
		if err != nil {
			return
		}
	}

	for i := range *instance.process.IntermediateCatchEvents() {
		element := &(*instance.process.IntermediateCatchEvents())[i]
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var intermediateCatchEvent *CatchEvent
		intermediateCatchEvent, err = NewCatchEvent(ctx, wiring, &element.CatchEvent)
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, intermediateCatchEvent)
		if err != nil {
			return
		}
	}

	for i := range *instance.process.Tasks() {
		element := &(*instance.process.Tasks())[i]
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var harness *Harness
		harness, err = NewHarness(ctx, wiring, idGenerator, NewTask(ctx, element))
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, harness)
		if err != nil {
			return
		}
	}

	for i := range *instance.process.BusinessRuleTasks() {
		element := &(*instance.process.BusinessRuleTasks())[i]
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var harness *Harness
		harness, err = NewHarness(ctx, wiring, idGenerator, NewBusinessRuleTask(ctx, element))
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, harness)
		if err != nil {
			return
		}
	}

	for i := range *instance.process.CallActivities() {
		element := &(*instance.process.CallActivities())[i]
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var harness *Harness
		harness, err = NewHarness(ctx, wiring, idGenerator, NewCallActivity(ctx, element))
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, harness)
		if err != nil {
			return
		}
	}

	for i := range *instance.process.ManualTasks() {
		element := &(*instance.process.ManualTasks())[i]
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var harness *Harness
		harness, err = NewHarness(ctx, wiring, idGenerator, NewManualTask(ctx, element))
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, harness)
		if err != nil {
			return
		}
	}

	for i := range *instance.process.ServiceTasks() {
		element := &(*instance.process.ServiceTasks())[i]
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var harness *Harness
		serviceTask := NewServiceTask(ctx, element)
		harness, err = NewHarness(ctx, wiring, idGenerator, serviceTask)
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, harness)
		if err != nil {
			return
		}
	}

	for i := range *instance.process.UserTasks() {
		element := &(*instance.process.UserTasks())[i]
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var harness *Harness
		userTask := NewUserTask(ctx, element)
		harness, err = NewHarness(ctx, wiring, idGenerator, userTask)
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, harness)
		if err != nil {
			return
		}
	}

	for i := range *instance.process.ReceiveTasks() {
		element := &(*instance.process.ReceiveTasks())[i]
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var harness *Harness
		scriptTask := NewReceiveTask(ctx, element)
		harness, err = NewHarness(ctx, wiring, idGenerator, scriptTask)
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, harness)
		if err != nil {
			return
		}
	}

	for i := range *instance.process.ScriptTasks() {
		element := &(*instance.process.ScriptTasks())[i]
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var harness *Harness
		scriptTask := NewScriptTask(ctx, element)
		harness, err = NewHarness(ctx, wiring, idGenerator, scriptTask)
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, harness)
		if err != nil {
			return
		}
	}

	for i := range *instance.process.SendTasks() {
		element := &(*instance.process.SendTasks())[i]
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var harness *Harness
		scriptTask := NewSendTask(ctx, element)
		harness, err = NewHarness(ctx, wiring, idGenerator, scriptTask)
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, harness)
		if err != nil {
			return
		}
	}

	for i := range *instance.process.SubProcesses() {
		element := &(*instance.process.SubProcesses())[i]
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var harness *Harness
		subProcess := NewSubProcess(ctx, instance.eventDefinitionInstanceBuilder, idGenerator, subTracer, element)
		harness, err = NewHarness(ctx, wiring, idGenerator, subProcess)
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, harness)
		if err != nil {
			return
		}
	}

	for i := range *instance.process.ExclusiveGateways() {
		element := &(*instance.process.ExclusiveGateways())[i]
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var exclusiveGateway *ExclusiveGateway
		exclusiveGateway, err = NewExclusiveGateway(ctx, wiring, element)
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, exclusiveGateway)
		if err != nil {
			return
		}
	}

	for i := range *instance.process.InclusiveGateways() {
		element := &(*instance.process.InclusiveGateways())[i]
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var inclusiveGateway *InclusiveGateway
		inclusiveGateway, err = NewInclusiveGateway(ctx, wiring, element)
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, inclusiveGateway)
		if err != nil {
			return
		}
	}

	for i := range *instance.process.ParallelGateways() {
		element := &(*instance.process.ParallelGateways())[i]
		var parallelGateway *ParallelGateway
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		parallelGateway, err = NewParallelGateway(ctx, wiring, element)
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, parallelGateway)
		if err != nil {
			return
		}
	}

	for i := range *instance.process.EventBasedGateways() {
		element := &(*instance.process.EventBasedGateways())[i]
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var eventBasedGateway *EventBasedGateway
		eventBasedGateway, err = NewEventBasedGateway(ctx, wiring, element)
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, eventBasedGateway)
		if err != nil {
			return
		}
	}

	instance.flowNodeMapping.Finalize()

	// StartAll cease flow monitor
	sender := instance.tracer.RegisterSender()
	go instance.ceaseFlowMonitor(subTracer)(ctx, sender)

	instance.tracer.Trace(InstantiationTrace{InstanceId: instance.id})

	return
}

// StartWith explicitly starts the instance by triggering a given start event
func (ins *Instance) StartWith(ctx context.Context, startEvent schema.StartEventInterface) (err error) {
	flowNode, found := ins.flowNodeMapping.ResolveElementToFlowNode(startEvent)
	elementId := "<unnamed>"
	if idPtr, present := startEvent.Id(); present {
		elementId = *idPtr
	}
	processId := "<unnamed>"
	if idPtr, present := ins.process.Id(); present {
		processId = *idPtr
	}
	if !found {
		err = errors.NotFoundError{Expected: fmt.Sprintf("start event %s in process %s", elementId, processId)}
		return
	}
	startEventNode, ok := flowNode.(*StartEvent)
	if !ok {
		err = errors.RequirementExpectationError{
			Expected: fmt.Sprintf("start event %s flow node in process %s to be of type start.Node", elementId, processId),
			Actual:   fmt.Sprintf("%Flow", flowNode),
		}
		return
	}
	startEventNode.Trigger()
	return
}

// StartAll explicitly starts the instance by triggering all start events, if any
func (ins *Instance) StartAll() (err error) {
	ctx := ins.ctx
	for i := range *ins.process.StartEvents() {
		err = ins.StartWith(ctx, &(*ins.process.StartEvents())[i])
		if err != nil {
			return
		}
	}
	return
}

func (ins *Instance) ceaseFlowMonitor(tracer tracing.ITracer) func(ctx context.Context, sender tracing.ISenderHandle) {
	// Subscribing to traces early as otherwise events produced
	// after the goroutine below is started are not going to be
	// sent to it.
	traces := tracer.Subscribe()
	ins.complete.Lock()
	return func(ctx context.Context, sender tracing.ISenderHandle) {
		defer sender.Done()
		defer ins.complete.Unlock()

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
			if len(startEventsActivated) == len(*ins.process.StartEvents()) {
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
			ins.flowWaitGroup.Wait()
			close(waitIsOver)
		}()
		select {
		case <-waitIsOver:
			// Send out a cease flow trace
			tracer.Trace(CeaseFlowTrace{Process: ins.process})
		case <-ctx.Done():
		}
	}
}

// WaitUntilComplete waits until the instance is complete.
// Returns true if the instance was complete, false if the context signalled `Done`
func (ins *Instance) WaitUntilComplete(ctx context.Context) (complete bool) {
	signal := make(chan bool)
	go func() {
		ins.complete.Lock()
		defer ins.complete.Unlock()
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
