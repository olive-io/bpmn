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
	"sync/atomic"

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

type subProcess struct {
	wr *wiring

	ctx                    context.Context
	cancel                 context.CancelFunc
	id                     id.Id
	element                *schema.SubProcess
	subTracer              tracing.ITracer
	flowNodeMapping        *FlowNodeMapping
	flowWaitGroup          sync.WaitGroup
	active                 atomic.Bool
	complete               sync.RWMutex
	idGenerator            id.IGenerator
	eventDefinitionBuilder event.IDefinitionInstanceBuilder
	eventConsumersLock     sync.RWMutex
	eventConsumers         []event.IConsumer
	mch                    chan imessage
}

func newSubProcess(eventBuilder event.IDefinitionInstanceBuilder, idGenerator id.IGenerator, subProcessElement *schema.SubProcess) constructor {
	return func(parentWiring *wiring) (act Activity, err error) {

		flowNodeMapping := NewLockedFlowNodeMapping()
		defer flowNodeMapping.Finalize()

		ctx, cancel := context.WithCancel(context.Background())
		subTracer := tracing.NewTracer(ctx)
		process := &subProcess{
			wr:                     parentWiring,
			ctx:                    ctx,
			cancel:                 cancel,
			id:                     idGenerator.New(),
			element:                subProcessElement,
			subTracer:              subTracer,
			eventDefinitionBuilder: eventBuilder,
			idGenerator:            idGenerator,
			flowNodeMapping:        flowNodeMapping,
			mch:                    make(chan imessage, len(parentWiring.incoming)*2+1),
		}

		locator := parentWiring.locator
		err = data.ElementToLocator(locator, idGenerator, subProcessElement)
		if err != nil {
			return
		}

		wiringMaker := func(element *schema.FlowNode) (*wiring, error) {
			return newWiring(
				parentWiring.processInstanceId,
				subProcessElement,
				parentWiring.definitions,
				element,
				parentWiring.eventIngress, process,
				subTracer,
				process.flowNodeMapping,
				&process.flowWaitGroup, process.eventDefinitionBuilder,
				locator)
		}

		var wr *wiring

		for i := range *subProcessElement.StartEvents() {
			element := &(*subProcessElement.StartEvents())[i]
			wr, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var startEventNode *startEvent
			startEventNode, err = newStartEvent(wr, element, idGenerator)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, startEventNode)
			if err != nil {
				return
			}
		}

		for i := range *subProcessElement.EndEvents() {
			element := &(*subProcessElement.EndEvents())[i]
			wr, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var endEventNode *endEvent
			endEventNode, err = newEndEvent(wr, element)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, endEventNode)
			if err != nil {
				return
			}
		}

		for i := range *subProcessElement.IntermediateCatchEvents() {
			element := &(*subProcessElement.IntermediateCatchEvents())[i]
			wr, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var intermediateCatchEvent *catchEvent
			intermediateCatchEvent, err = newCatchEvent(wr, &element.CatchEvent)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, intermediateCatchEvent)
			if err != nil {
				return
			}
		}

		for i := range *subProcessElement.BusinessRuleTasks() {
			element := &(*subProcessElement.BusinessRuleTasks())[i]
			wr, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var node *harness
			node, err = newHarness(wr, idGenerator, newTask(element, BusinessRuleActivity))
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, node)
			if err != nil {
				return
			}
		}

		for i := range *subProcessElement.CallActivities() {
			element := &(*subProcessElement.CallActivities())[i]
			wr, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var node *harness
			node, err = newHarness(wr, idGenerator, newTask(element, CallActivity))
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, node)
			if err != nil {
				return
			}
		}

		for i := range *subProcessElement.Tasks() {
			element := &(*subProcessElement.Tasks())[i]
			wr, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var node *harness
			node, err = newHarness(wr, idGenerator, newTask(element, TaskActivity))
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, node)
			if err != nil {
				return
			}
		}

		for i := range *subProcessElement.ManualTasks() {
			element := &(*subProcessElement.ManualTasks())[i]
			wr, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var node *harness
			node, err = newHarness(wr, idGenerator, newTask(element, ManualTaskActivity))
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, node)
			if err != nil {
				return
			}
		}

		for i := range *subProcessElement.ServiceTasks() {
			element := &(*subProcessElement.ServiceTasks())[i]
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
			err = flowNodeMapping.RegisterElementToFlowNode(element, node)
			if err != nil {
				return
			}
		}

		for i := range *subProcessElement.UserTasks() {
			element := &(*subProcessElement.UserTasks())[i]
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
			err = flowNodeMapping.RegisterElementToFlowNode(element, node)
			if err != nil {
				return
			}
		}

		for i := range *subProcessElement.ReceiveTasks() {
			element := &(*subProcessElement.ReceiveTasks())[i]
			wr, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var node *harness
			scriptTask := newTask(element, ReceiveTaskActivity)
			node, err = newHarness(wr, idGenerator, scriptTask)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, node)
			if err != nil {
				return
			}
		}

		for i := range *subProcessElement.ScriptTasks() {
			element := &(*subProcessElement.ScriptTasks())[i]
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
			err = flowNodeMapping.RegisterElementToFlowNode(element, node)
			if err != nil {
				return
			}
		}

		for i := range *subProcessElement.SendTasks() {
			element := &(*subProcessElement.SendTasks())[i]
			wr, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var node *harness
			scriptTask := newTask(element, SendTaskActivity)
			node, err = newHarness(wr, idGenerator, scriptTask)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, node)
			if err != nil {
				return
			}
		}

		for i := range *subProcessElement.SubProcesses() {
			element := &(*subProcessElement.SubProcesses())[i]
			wr, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var node *harness
			sp := newSubProcess(eventBuilder, idGenerator, element)
			node, err = newHarness(wr, idGenerator, sp)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, node)
			if err != nil {
				return
			}
		}

		for i := range *subProcessElement.ExclusiveGateways() {
			element := &(*subProcessElement.ExclusiveGateways())[i]
			wr, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var node *exclusiveGateway
			node, err = newExclusiveGateway(wr, element)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, node)
			if err != nil {
				return
			}
		}

		for i := range *subProcessElement.InclusiveGateways() {
			element := &(*subProcessElement.InclusiveGateways())[i]
			wr, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var node *inclusiveGateway
			node, err = newInclusiveGateway(wr, element)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, node)
			if err != nil {
				return
			}
		}

		for i := range *subProcessElement.ParallelGateways() {
			element := &(*subProcessElement.ParallelGateways())[i]
			wr, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var node *parallelGateway
			node, err = newParallelGateway(wr, element)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, node)
			if err != nil {
				return
			}
		}

		for i := range *subProcessElement.EventBasedGateways() {
			element := &(*subProcessElement.EventBasedGateways())[i]
			wr, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var node *eventBasedGateway
			node, err = newEventBasedGateway(wr, element)
			if err != nil {
				return
			}
			err = flowNodeMapping.RegisterElementToFlowNode(element, node)
			if err != nil {
				return
			}
		}

		act = process
		return
	}
}

func (p *subProcess) ConsumeEvent(ev event.IEvent) (result event.ConsumptionResult, err error) {
	p.eventConsumersLock.RLock()
	// We're copying the list of consumers here to ensure that
	// new consumers can subscribe during event forwarding
	eventConsumers := p.eventConsumers
	p.eventConsumersLock.RUnlock()
	result, err = event.ForwardEvent(ev, &eventConsumers)
	return
}

func (p *subProcess) RegisterEventConsumer(ev event.IConsumer) (err error) {
	p.eventConsumersLock.Lock()
	defer p.eventConsumersLock.Unlock()
	p.eventConsumers = append(p.eventConsumers, ev)
	return
}

// startWith explicitly starts the subprocess by triggering a given start event
func (p *subProcess) startWith(ctx context.Context, element schema.StartEventInterface) (err error) {
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

// startAll explicitly starts the subprocess by triggering all start events, if any
func (p *subProcess) startAll(ctx context.Context) (err error) {
	for i := range *p.element.StartEvents() {
		err = p.startWith(ctx, &(*p.element.StartEvents())[i])
		if err != nil {
			return
		}
	}
	return
}

func (p *subProcess) ceaseFlowMonitor(tracer tracing.ITracer) func(ctx context.Context, sender tracing.ISenderHandle) {
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

func (p *subProcess) run(ctx context.Context, out tracing.ITracer) {
	defer p.cancel()
	for {
		select {
		case msg := <-p.mch:
			switch m := msg.(type) {
			case cancelMessage:
				m.response <- true
				return
			case nextActionMessage:
				if p.active.Load() {
					continue
				}

				p.active.Store(true)
				go func() {
					defer p.active.Store(false)

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

					traces := p.subTracer.Subscribe()
					defer p.subTracer.Unsubscribe(traces)
				loop:
					for {
						var trace tracing.ITrace
						select {
						case trace = <-traces:
						case <-ctx.Done():
							p.wr.tracer.Send(CancellationFlowNodeTrace{Node: p.element})
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

					action := FlowAction{SequenceFlows: allSequenceFlows(&p.wr.outgoing)}
					m.response <- action
				}()
			default:
			}
		case <-ctx.Done():
			p.wr.tracer.Send(CancellationFlowNodeTrace{Node: p.element})
			return
		}
	}
}

func (p *subProcess) NextAction(ctx context.Context, flow Flow) chan IAction {
	// flow nodes
	// StartAll cease flow monitor
	sender := p.subTracer.RegisterSender()
	go p.ceaseFlowMonitor(p.wr.tracer)(ctx, sender)

	go p.run(ctx, p.wr.tracer)

	response := make(chan IAction, 1)
	p.mch <- nextActionMessage{response: response}
	return response
}

func (p *subProcess) Element() schema.FlowNodeInterface {
	return p.element
}

func (p *subProcess) Type() ActivityType {
	return SubprocessActivity
}

func (p *subProcess) Cancel() <-chan bool {
	response := make(chan bool)
	p.mch <- cancelMessage{response: response}
	return response
}
