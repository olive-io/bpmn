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

package subprocess

import (
	"context"
	"fmt"
	"sync"

	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/errors"
	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/flow_node/activity"
	"github.com/olive-io/bpmn/flow_node/activity/business_rule"
	"github.com/olive-io/bpmn/flow_node/activity/call"
	"github.com/olive-io/bpmn/flow_node/activity/receive"
	"github.com/olive-io/bpmn/flow_node/activity/script"
	"github.com/olive-io/bpmn/flow_node/activity/send"
	"github.com/olive-io/bpmn/flow_node/activity/service"
	"github.com/olive-io/bpmn/flow_node/activity/task"
	"github.com/olive-io/bpmn/flow_node/activity/user"
	"github.com/olive-io/bpmn/flow_node/event/catch"
	"github.com/olive-io/bpmn/flow_node/event/end"
	"github.com/olive-io/bpmn/flow_node/event/start"
	"github.com/olive-io/bpmn/flow_node/gateway/event_based"
	"github.com/olive-io/bpmn/flow_node/gateway/exclusive"
	"github.com/olive-io/bpmn/flow_node/gateway/inclusive"
	"github.com/olive-io/bpmn/flow_node/gateway/parallel"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tools/id"
	"github.com/olive-io/bpmn/tracing"
)

type imessage interface {
	message()
}

type nextActionMessage struct {
	response chan flow_node.IAction
}

func (m nextActionMessage) message() {}

type cancelMessage struct {
	response chan bool
}

func (m cancelMessage) message() {}

type SubProcess struct {
	*flow_node.Wiring
	ctx                    context.Context
	cancel                 context.CancelFunc
	element                *schema.SubProcess
	tracer                 tracing.ITracer
	flowNodeMapping        *flow_node.FlowNodeMapping
	flowWaitGroup          sync.WaitGroup
	complete               sync.RWMutex
	idGenerator            id.IGenerator
	eventDefinitionBuilder event.IDefinitionInstanceBuilder
	eventConsumersLock     sync.RWMutex
	eventConsumers         []event.IConsumer
	runnerChannel          chan imessage
}

func New(ctx context.Context,
	eventDefinitionBuilder event.IDefinitionInstanceBuilder,
	idGenerator id.IGenerator,
	tracer tracing.ITracer,
	subProcess *schema.SubProcess) activity.Constructor {
	return func(parentWiring *flow_node.Wiring) (node activity.Activity, err error) {

		flowNodeMapping := flow_node.NewLockedFlowNodeMapping()

		subTracer := tracing.NewTracer(ctx)

		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		process := &SubProcess{
			Wiring:                 parentWiring,
			ctx:                    ctx,
			cancel:                 cancel,
			element:                subProcess,
			tracer:                 subTracer,
			flowNodeMapping:        flowNodeMapping,
			eventDefinitionBuilder: eventDefinitionBuilder,
			idGenerator:            idGenerator,
			runnerChannel:          make(chan imessage, len(parentWiring.Incoming)*2+1),
		}

		locator := parentWiring.Locator
		err = data.ElementToLocator(locator, idGenerator, subProcess)
		if err != nil {
			return
		}

		wiringMaker := func(element *schema.FlowNode) (*flow_node.Wiring, error) {
			return flow_node.NewWiring(
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

		var wiring *flow_node.Wiring

		for i := range *subProcess.StartEvents() {
			element := &(*subProcess.StartEvents())[i]
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			var startEvent *start.Node
			startEvent, err = start.New(ctx, wiring, element, idGenerator)
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
			var endEvent *end.Node
			endEvent, err = end.New(ctx, wiring, element)
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
			var intermediateCatchEvent *catch.Node
			intermediateCatchEvent, err = catch.New(ctx, wiring, &element.CatchEvent)
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
			var harness *activity.Harness
			harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, business_rule.NewBusinessRuleTask(ctx, element))
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
			var harness *activity.Harness
			harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, call.NewCallActivity(ctx, element))
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
			var harness *activity.Harness
			harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, task.NewTask(ctx, element))
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
			var harness *activity.Harness
			serviceTask := service.NewServiceTask(ctx, element)
			harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, serviceTask)
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
			var harness *activity.Harness
			userTask := user.NewUserTask(ctx, element)
			harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, userTask)
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
			var harness *activity.Harness
			scriptTask := receive.NewReceiveTask(ctx, element)
			harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, scriptTask)
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
			var harness *activity.Harness
			scriptTask := script.NewScriptTask(ctx, element)
			harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, scriptTask)
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
			var harness *activity.Harness
			scriptTask := send.NewSendTask(ctx, element)
			harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, scriptTask)
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
			var harness *activity.Harness
			sp := New(ctx, eventDefinitionBuilder, idGenerator, subTracer, element)
			harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, sp)
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
			var exclusiveGateway *exclusive.Node
			exclusiveGateway, err = exclusive.New(ctx, wiring, element)
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
			var inclusiveGateway *inclusive.Node
			inclusiveGateway, err = inclusive.New(ctx, wiring, element)
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
			var parallelGateway *parallel.Node
			wiring, err = wiringMaker(&element.FlowNode)
			if err != nil {
				return
			}
			parallelGateway, err = parallel.New(ctx, wiring, element)
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
			var eventBasedGateway *event_based.Node
			eventBasedGateway, err = event_based.New(ctx, wiring, element)
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
		go process.runner(ctx, tracer)
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
	startEventNode, ok := flowNode.(*start.Node)
	if !ok {
		err = errors.RequirementExpectationError{
			Expected: fmt.Sprintf("start event %s flow node in process %s to be of type start.Node", elementId, processId),
			Actual:   fmt.Sprintf("%T", flowNode),
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
				case flow.TerminationTrace:
					switch flowNode := t.Source.(type) {
					case *schema.StartEvent:
						startEventsActivated = append(startEventsActivated, flowNode)
					default:
					}
				case flow.Trace:
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
			tracer.Trace(flow.CeaseFlowTrace{})
		case <-ctx.Done():

		}

	}
}

func (p *SubProcess) runner(ctx context.Context, out tracing.ITracer) {
	for {
		select {
		case msg := <-p.runnerChannel:
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
						out.Trace(tracing.ErrorTrace{Error: &errors.SubProcessError{
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
							p.Tracer.Trace(flow_node.CancellationTrace{Node: p.element})
							return
						}

						trace = tracing.Unwrap(trace)
						switch trace := trace.(type) {
						case flow.CeaseFlowTrace:
							out.Trace(ProcessLandMarkTrace{Node: p.element})
							break loop
						// case flow.CompletionTrace:
						// ignore end event of sub process
						case flow.TerminationTrace:
							// ignore end event of sub process
						default:
							out.Trace(trace)
						}
					}
					p.tracer.Unsubscribe(traces)

					action := flow_node.FlowAction{SequenceFlows: flow_node.AllSequenceFlows(&p.Wiring.Outgoing)}
					m.response <- action
				}()
			default:
			}
		case <-ctx.Done():
			p.Tracer.Trace(flow_node.CancellationTrace{Node: p.element})
			return
		}
	}
}

func (p *SubProcess) NextAction(flow_interface.T) chan flow_node.IAction {
	response := make(chan flow_node.IAction, 1)
	p.runnerChannel <- nextActionMessage{response: response}
	return response
}

func (p *SubProcess) Element() schema.FlowNodeInterface {
	return p.element
}

func (p *SubProcess) Cancel() <-chan bool {
	response := make(chan bool)
	p.runnerChannel <- cancelMessage{response: response}
	return response
}
