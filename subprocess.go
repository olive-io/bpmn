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

// ProcessLandMarkTrace denotes instantiation of a given sub process
type ProcessLandMarkTrace struct {
	Node schema.FlowNodeInterface
}

func (t ProcessLandMarkTrace) Unpack() any { return t.Node }

type SubProcess struct {
	*Wiring
	ctx    context.Context
	cancel context.CancelFunc

	element                *schema.SubProcess
	parentTracer           tracing.ITracer
	tracer                 tracing.ITracer
	flowNodeMapping        *FlowNodeMapping
	flowWaitGroup          sync.WaitGroup
	complete               sync.RWMutex
	idGenerator            id.IGenerator
	eventDefinitionBuilder event.IDefinitionInstanceBuilder
	eventConsumersLock     sync.RWMutex
	eventConsumers         []event.IConsumer
	mch                    chan imessage
}

func NewSubProcess(ctx context.Context,
	eventDefinitionBuilder event.IDefinitionInstanceBuilder,
	idGenerator id.IGenerator,
	tracer tracing.ITracer,
	subProcess *schema.SubProcess) Constructor {
	return func(parentWiring *Wiring) (node Activity, err error) {

		flowNodeMapping := NewLockedFlowNodeMapping()

		subTracer := tracing.NewTracer(ctx)

		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)

		process := &SubProcess{
			Wiring:                 parentWiring,
			ctx:                    ctx,
			cancel:                 cancel,
			element:                subProcess,
			parentTracer:           tracer,
			tracer:                 subTracer,
			flowNodeMapping:        flowNodeMapping,
			eventDefinitionBuilder: eventDefinitionBuilder,
			idGenerator:            idGenerator,
			mch:                    make(chan imessage, len(parentWiring.Incoming)*2+1),
		}

		locator := parentWiring.Locator
		err = data.ElementToLocator(locator, idGenerator, subProcess)
		if err != nil {
			return
		}

		wiringMaker := func(element *schema.FlowNode) (*Wiring, error) {
			return NewWiring(
				parentWiring.ProcessInstanceId,
				subProcess,
				parentWiring.Definitions,
				element,
				parentWiring.EventIngress, process,
				subTracer,
				process.flowNodeMapping,
				&process.flowWaitGroup, process.eventDefinitionBuilder,
				locator)
		}

		var wiring *Wiring

		for i := range *subProcess.StartEvents() {
			element := &(*subProcess.StartEvents())[i]
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var startEvent *StartEvent
			startEvent, err = NewStartEvent(ctx, wiring, element, idGenerator)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, startEvent)
			if err != nil {
				return
			}
		}

		for i := range *subProcess.EndEvents() {
			element := &(*subProcess.EndEvents())[i]
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var endEvent *EndEvent
			endEvent, err = NewEndEvent(ctx, wiring, element)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, endEvent)
			if err != nil {
				return
			}
		}

		for i := range *subProcess.IntermediateCatchEvents() {
			element := &(*subProcess.IntermediateCatchEvents())[i]
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var intermediateCatchEvent *CatchEvent
			intermediateCatchEvent, err = NewCatchEvent(ctx, wiring, &element.CatchEvent)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, intermediateCatchEvent)
			if err != nil {
				return
			}
		}

		for i := range *subProcess.BusinessRuleTasks() {
			element := &(*subProcess.BusinessRuleTasks())[i]
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var harness *Harness
			harness, err = NewHarness(ctx, wiring, idGenerator, NewTask(ctx, element, BusinessRuleActivity))
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, harness)
			if err != nil {
				return
			}
		}

		for i := range *subProcess.CallActivities() {
			element := &(*subProcess.CallActivities())[i]
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var harness *Harness
			harness, err = NewHarness(ctx, wiring, idGenerator, NewTask(ctx, element, CallActivity))
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, harness)
			if err != nil {
				return
			}
		}

		for i := range *subProcess.Tasks() {
			element := &(*subProcess.Tasks())[i]
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var harness *Harness
			harness, err = NewHarness(ctx, wiring, idGenerator, NewTask(ctx, element, TaskActivity))
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, harness)
			if err != nil {
				return
			}
		}

		for i := range *subProcess.ManualTasks() {
			element := &(*subProcess.ManualTasks())[i]
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var harness *Harness
			harness, err = NewHarness(ctx, wiring, idGenerator, NewTask(ctx, element, ManualTaskActivity))
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, harness)
			if err != nil {
				return
			}
		}

		for i := range *subProcess.ServiceTasks() {
			element := &(*subProcess.ServiceTasks())[i]
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var harness *Harness
			serviceTask := NewTask(ctx, element, ServiceTaskActivity)
			harness, err = NewHarness(ctx, wiring, idGenerator, serviceTask)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, harness)
			if err != nil {
				return
			}
		}

		for i := range *subProcess.UserTasks() {
			element := &(*subProcess.UserTasks())[i]
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var harness *Harness
			userTask := NewTask(ctx, element, UserTaskActivity)
			harness, err = NewHarness(ctx, wiring, idGenerator, userTask)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, harness)
			if err != nil {
				return
			}
		}

		for i := range *subProcess.ReceiveTasks() {
			element := &(*subProcess.ReceiveTasks())[i]
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var harness *Harness
			scriptTask := NewTask(ctx, element, ReceiveTaskActivity)
			harness, err = NewHarness(ctx, wiring, idGenerator, scriptTask)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, harness)
			if err != nil {
				return
			}
		}

		for i := range *subProcess.ScriptTasks() {
			element := &(*subProcess.ScriptTasks())[i]
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var harness *Harness
			scriptTask := NewTask(ctx, element, ScriptTaskActivity)
			harness, err = NewHarness(ctx, wiring, idGenerator, scriptTask)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, harness)
			if err != nil {
				return
			}
		}

		for i := range *subProcess.SendTasks() {
			element := &(*subProcess.SendTasks())[i]
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var harness *Harness
			scriptTask := NewTask(ctx, element, SendTaskActivity)
			harness, err = NewHarness(ctx, wiring, idGenerator, scriptTask)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, harness)
			if err != nil {
				return
			}
		}

		for i := range *subProcess.SubProcesses() {
			element := &(*subProcess.SubProcesses())[i]
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var harness *Harness
			sp := NewSubProcess(ctx, eventDefinitionBuilder, idGenerator, subTracer, element)
			harness, err = NewHarness(ctx, wiring, idGenerator, sp)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, harness)
			if err != nil {
				return
			}
		}

		for i := range *subProcess.ExclusiveGateways() {
			element := &(*subProcess.ExclusiveGateways())[i]
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var exclusiveGateway *ExclusiveGateway
			exclusiveGateway, err = NewExclusiveGateway(ctx, wiring, element)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, exclusiveGateway)
			if err != nil {
				return
			}
		}

		for i := range *subProcess.InclusiveGateways() {
			element := &(*subProcess.InclusiveGateways())[i]
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var inclusiveGateway *InclusiveGateway
			inclusiveGateway, err = NewInclusiveGateway(ctx, wiring, element)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, inclusiveGateway)
			if err != nil {
				return
			}
		}

		for i := range *subProcess.ParallelGateways() {
			element := &(*subProcess.ParallelGateways())[i]
			var parallelGateway *ParallelGateway
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			parallelGateway, err = NewParallelGateway(ctx, wiring, element)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, parallelGateway)
			if err != nil {
				return
			}
		}

		for i := range *subProcess.EventBasedGateways() {
			element := &(*subProcess.EventBasedGateways())[i]
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var eventBasedGateway *EventBasedGateway
			eventBasedGateway, err = NewEventBasedGateway(ctx, wiring, element)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, eventBasedGateway)
			if err != nil {
				return
			}
		}

		flowNodeMapping.Finalize()
		// StartAll cease flow monitor
		sender := process.Tracer.RegisterSender()
		go process.ceaseFlowMonitor(subTracer)(ctx, sender)
		node = process
		return
	}
}

func (p *SubProcess) ConsumeEvent(ev event.IEvent) (result event.ConsumptionResult, err error) {
	p.eventConsumersLock.RLock()
	// We're copying the list of consumers here to ensure that
	// new consumers can subscribe during event forwarding
	eventConsumers := p.eventConsumers
	p.eventConsumersLock.RUnlock()
	result, err = event.ForwardEvent(ev, &eventConsumers)
	return
}

func (p *SubProcess) RegisterEventConsumer(ev event.IConsumer) (err error) {
	p.eventConsumersLock.Lock()
	defer p.eventConsumersLock.Unlock()
	p.eventConsumers = append(p.eventConsumers, ev)
	return
}

// startWith explicitly starts the subprocess by triggering a given start event
func (p *SubProcess) startWith(ctx context.Context, startEvent schema.StartEventInterface) (err error) {
	flowNode, found := p.flowNodeMapping.ResolveElementToFlowNode(startEvent)
	elementId := "<unnamed>"
	if idPtr, present := startEvent.Id(); present {
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
	startEventNode, ok := flowNode.(*StartEvent)
	if !ok {
		err = errors.RequirementExpectationError{
			Expected: fmt.Sprintf("start event %s flow node in process %s to be of type start.Node", elementId, processId),
			Actual:   fmt.Sprintf("%sFlow", flowNode),
		}
		return
	}
	startEventNode.Trigger()
	return
}

// startAll explicitly starts the subprocess by triggering all start events, if any
func (p *SubProcess) startAll(ctx context.Context) (err error) {
	for i := range *p.element.StartEvents() {
		err = p.startWith(ctx, &(*p.element.StartEvents())[i])
		if err != nil {
			return
		}
	}
	return
}

func (p *SubProcess) ceaseFlowMonitor(tracer tracing.ITracer) func(ctx context.Context, sender tracing.ISenderHandle) {
	// Subscribing to traces early as otherwise events produced
	// after the goroutine below is started are not going to be
	// sent to it.
	traces := tracer.Subscribe()
	p.complete.Lock()
	return func(ctx context.Context, sender tracing.ISenderHandle) {
		defer sender.Done()
		defer p.complete.Unlock()

		startEventsActivated := make([]*schema.StartEvent, 0)

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

func (p *SubProcess) run(ctx context.Context, out tracing.ITracer) {
	for {
		select {
		case msg := <-p.mch:
			switch m := msg.(type) {
			case cancelMessage:
				p.cancel()
				m.response <- true
			case nextActionMessage:
				go func() {
					if err := p.startAll(ctx); err != nil {
						subProcessId := ""
						if pid, present := p.element.Id(); present {
							subProcessId = *pid
						}
						out.Send(ErrorTrace{Error: &errors.SubProcessError{
							Id:     subProcessId,
							Reason: err.Error(),
						}})
						return
					}

					traces := p.tracer.Subscribe()
				loop:
					for {
						var trace tracing.ITrace
						select {
						case trace = <-traces:
						case <-ctx.Done():
							p.Tracer.Send(CancellationFlowNodeTrace{Node: p.element})
							return
						}

						trace = tracing.Unwrap(trace)
						switch tr := trace.(type) {
						case CeaseFlowTrace:
							out.Send(ProcessLandMarkTrace{Node: p.element})
							break loop
						case CompletionTrace:
							// ignore end event of subprocess
						case TerminationTrace:
							// ignore end event of subprocess
						default:
							out.Send(tr)
						}
					}
					p.tracer.Unsubscribe(traces)

					action := FlowAction{SequenceFlows: AllSequenceFlows(&p.Wiring.Outgoing)}
					m.response <- action
				}()
			default:
			}
		case <-ctx.Done():
			p.Tracer.Send(CancellationFlowNodeTrace{Node: p.element})
			return
		}
	}
}

func (p *SubProcess) NextAction(Flow) chan IAction {
	go p.run(p.ctx, p.parentTracer)

	response := make(chan IAction, 1)
	p.mch <- nextActionMessage{response: response}
	return response
}

func (p *SubProcess) Element() schema.FlowNodeInterface {
	return p.element
}

func (p *SubProcess) Type() ActivityType {
	return SubprocessActivity
}

func (p *SubProcess) Cancel() <-chan bool {
	response := make(chan bool)
	p.mch <- cancelMessage{response: response}
	return response
}
