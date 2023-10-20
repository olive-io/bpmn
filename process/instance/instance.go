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

package instance

import (
	"context"
	"fmt"
	"sync"

	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/errors"
	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/flow_node/activity"
	"github.com/olive-io/bpmn/flow_node/activity/business_rule"
	"github.com/olive-io/bpmn/flow_node/activity/call"
	"github.com/olive-io/bpmn/flow_node/activity/receive"
	"github.com/olive-io/bpmn/flow_node/activity/script"
	"github.com/olive-io/bpmn/flow_node/activity/send"
	"github.com/olive-io/bpmn/flow_node/activity/service"
	"github.com/olive-io/bpmn/flow_node/activity/subprocess"
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

type Instance struct {
	id                             id.Id
	process                        *schema.Process
	Tracer                         tracing.ITracer
	flowNodeMapping                *flow_node.FlowNodeMapping
	flowWaitGroup                  sync.WaitGroup
	complete                       sync.RWMutex
	idGenerator                    id.IGenerator
	Locator                        *data.FlowDataLocator
	EventIngress                   event.IConsumer
	EventEgress                    event.ISource
	idGeneratorBuilder             id.IGeneratorBuilder
	eventDefinitionInstanceBuilder event.IDefinitionInstanceBuilder
	eventConsumersLock             sync.RWMutex
	eventConsumers                 []event.IConsumer
}

func (instance *Instance) Id() id.Id {
	return instance.id
}

func (instance *Instance) ConsumeEvent(ev event.IEvent) (result event.ConsumptionResult, err error) {
	instance.eventConsumersLock.RLock()
	// We're copying the list of consumers here to ensure that
	// new consumers can subscribe during event forwarding
	eventConsumers := instance.eventConsumers
	instance.eventConsumersLock.RUnlock()
	result, err = event.ForwardEvent(ev, &eventConsumers)
	return
}

func (instance *Instance) RegisterEventConsumer(ev event.IConsumer) (err error) {
	instance.eventConsumersLock.Lock()
	defer instance.eventConsumersLock.Unlock()
	instance.eventConsumers = append(instance.eventConsumers, ev)
	return
}

// Option allows to modify configuration of
// an instance in a flexible fashion (as its just a modification
// function)
//
// It also allows to augment or replace the context.
type Option func(ctx context.Context, instance *Instance) context.Context

// WithTracer overrides instance's tracer
func WithTracer(tracer tracing.ITracer) Option {
	return func(ctx context.Context, instance *Instance) context.Context {
		instance.Tracer = tracer
		return ctx
	}
}

// WithContext will pass a given context to a new instance
// instead of implicitly generated one
func WithContext(newCtx context.Context) Option {
	return func(ctx context.Context, instance *Instance) context.Context {
		return newCtx
	}
}

func WithIdGenerator(builder id.IGeneratorBuilder) Option {
	return func(ctx context.Context, instance *Instance) context.Context {
		instance.idGeneratorBuilder = builder
		return ctx
	}
}

func WithEventIngress(consumer event.IConsumer) Option {
	return func(ctx context.Context, instance *Instance) context.Context {
		instance.EventIngress = consumer
		return ctx
	}
}

func WithEventEgress(source event.ISource) Option {
	return func(ctx context.Context, instance *Instance) context.Context {
		instance.EventEgress = source
		return ctx
	}
}

func WithVariables(variables map[string]any) Option {
	return func(ctx context.Context, instance *Instance) context.Context {
		if instance.Locator == nil {
			instance.Locator = data.NewFlowDataLocator()
		}
		for key, value := range variables {
			instance.Locator.SetVariable(key, value)
		}
		return ctx
	}
}

func WithDataObjects(dataObjects map[string]any) Option {
	return func(ctx context.Context, instance *Instance) context.Context {
		if instance.Locator == nil {
			instance.Locator = data.NewFlowDataLocator()
		}
		for dataObjectId, dataObject := range dataObjects {
			locator, found := instance.Locator.FindIItemAwareLocator(data.LocatorObject)
			if !found {
				locator = data.NewDataObjectContainer()
				instance.Locator.PutIItemAwareLocator(data.LocatorObject, locator)
			}
			container := data.NewContainer(nil)
			container.Put(dataObject)
			locator.PutItemAwareById(dataObjectId, container)
		}
		return ctx
	}
}

func WithEventDefinitionInstanceBuilder(builder event.IDefinitionInstanceBuilder) Option {
	return func(ctx context.Context, instance *Instance) context.Context {
		instance.eventDefinitionInstanceBuilder = builder
		return ctx
	}
}

func (instance *Instance) FlowNodeMapping() *flow_node.FlowNodeMapping {
	return instance.flowNodeMapping
}

func NewInstance(element *schema.Process, definitions *schema.Definitions, options ...Option) (instance *Instance, err error) {
	instance = &Instance{
		process:         element,
		flowNodeMapping: flow_node.NewLockedFlowNodeMapping(),
	}

	ctx := context.Background()

	// Apply options
	for _, option := range options {
		ctx = option(ctx, instance)
	}

	if instance.Tracer == nil {
		instance.Tracer = tracing.NewTracer(ctx)
	}

	if instance.idGeneratorBuilder == nil {
		instance.idGeneratorBuilder = id.DefaultIdGeneratorBuilder
	}

	dataLocator := instance.Locator
	if dataLocator == nil {
		dataLocator = data.NewFlowDataLocator()
	}

	var idGenerator id.IGenerator
	idGenerator, err = instance.idGeneratorBuilder.NewIdGenerator(ctx, instance.Tracer)
	if err != nil {
		return
	}

	instance.idGenerator = idGenerator

	instance.id = idGenerator.New()

	err = instance.EventEgress.RegisterEventConsumer(instance)
	if err != nil {
		return
	}

	var procLocator *data.FlowDataLocator
	procLocator, err = data.NewFlowDataLocatorFromElement(idGenerator, instance.process)
	if err != nil {
		return
	}
	procLocator.Merge(dataLocator)
	instance.Locator = procLocator

	// Flow nodes

	subTracer := tracing.NewTracer(ctx)

	tracing.NewRelay(ctx, subTracer, instance.Tracer, func(trace tracing.ITrace) []tracing.ITrace {
		return []tracing.ITrace{Trace{
			InstanceId: instance.id,
			Trace:      trace,
		}}
	})

	wiringMaker := func(element *schema.FlowNode) (*flow_node.Wiring, error) {
		return flow_node.NewWiring(
			instance.id,
			instance.process,
			definitions,
			element,
			// Event ingress/egress orchestration:
			//
			// Flow nodes will send their message to `instance.EventIngress`
			// (which is typically the model), but consume their messages from
			// `instance`, which is turn a subscriber of `instance.EventEgress`
			// (again, typically, the model).
			//
			// This allows us to use ConsumeEvent on this instance to send
			// events only to the instance (useful for things like event-based
			// process instantiation)
			instance.EventIngress, instance,
			subTracer,
			instance.flowNodeMapping,
			&instance.flowWaitGroup, instance.eventDefinitionInstanceBuilder,
			instance.Locator)
	}

	var wiring *flow_node.Wiring

	for i := range *instance.process.StartEvents() {
		element := &(*instance.process.StartEvents())[i]
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		var startEvent *start.Node
		startEvent, err = start.New(ctx, wiring, element, idGenerator)
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
		var endEvent *end.Node
		endEvent, err = end.New(ctx, wiring, element)
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
		var intermediateCatchEvent *catch.Node
		intermediateCatchEvent, err = catch.New(ctx, wiring, &element.CatchEvent)
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, intermediateCatchEvent)
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
		var harness *activity.Harness
		harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, business_rule.NewBusinessRuleTask(ctx, element))
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
		var harness *activity.Harness
		harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, call.NewCallActivity(ctx, element))
		if err != nil {
			return
		}
		err = instance.flowNodeMapping.RegisterElementToFlowNode(element, harness)
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
		var harness *activity.Harness
		harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, task.NewTask(ctx, element))
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
		var harness *activity.Harness
		serviceTask := service.NewServiceTask(ctx, element)
		harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, serviceTask)
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
		var harness *activity.Harness
		userTask := user.NewUserTask(ctx, element)
		harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, userTask)
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
		var harness *activity.Harness
		scriptTask := receive.NewReceiveTask(ctx, element)
		harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, scriptTask)
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
		var harness *activity.Harness
		scriptTask := script.NewScriptTask(ctx, element)
		harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, scriptTask)
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
		var harness *activity.Harness
		scriptTask := send.NewSendTask(ctx, element)
		harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, scriptTask)
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
		var harness *activity.Harness
		subProcess := subprocess.New(ctx, instance.eventDefinitionInstanceBuilder, idGenerator, subTracer, element)
		harness, err = activity.NewHarness(ctx, wiring, &element.FlowNode, idGenerator, subProcess)
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
		var exclusiveGateway *exclusive.Node
		exclusiveGateway, err = exclusive.New(ctx, wiring, element)
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
		var inclusiveGateway *inclusive.Node
		inclusiveGateway, err = inclusive.New(ctx, wiring, element)
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
		var parallelGateway *parallel.Node
		wiring, err = wiringMaker(&element.FlowNode)
		if err != nil {
			return
		}
		parallelGateway, err = parallel.New(ctx, wiring, element)
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
		var eventBasedGateway *event_based.Node
		eventBasedGateway, err = event_based.New(ctx, wiring, element)
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
	sender := instance.Tracer.RegisterSender()
	go instance.ceaseFlowMonitor(subTracer)(ctx, sender)

	instance.Tracer.Trace(InstantiationTrace{InstanceId: instance.id})

	return
}

// StartWith explicitly starts the instance by triggering a given start event
func (instance *Instance) StartWith(ctx context.Context, startEvent schema.StartEventInterface) (err error) {
	flowNode, found := instance.flowNodeMapping.ResolveElementToFlowNode(startEvent)
	elementId := "<unnamed>"
	if idPtr, present := startEvent.Id(); present {
		elementId = *idPtr
	}
	processId := "<unnamed>"
	if idPtr, present := instance.process.Id(); present {
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

// StartAll explicitly starts the instance by triggering all start events, if any
func (instance *Instance) StartAll(ctx context.Context) (err error) {
	for i := range *instance.process.StartEvents() {
		err = instance.StartWith(ctx, &(*instance.process.StartEvents())[i])
		if err != nil {
			return
		}
	}
	return
}

func (instance *Instance) ceaseFlowMonitor(tracer tracing.ITracer) func(ctx context.Context, sender tracing.ISenderHandle) {
	// Subscribing to traces early as otherwise events produced
	// after the goroutine below is started are not going to be
	// sent to it.
	traces := tracer.Subscribe()
	instance.complete.Lock()
	return func(ctx context.Context, sender tracing.ISenderHandle) {
		defer sender.Done()
		defer instance.complete.Unlock()

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
			if len(startEventsActivated) == len(*instance.process.StartEvents()) {
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
			instance.flowWaitGroup.Wait()
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

// WaitUntilComplete waits until the instance is complete.
// Returns true if the instance was complete, false if the context signalled `Done`
func (instance *Instance) WaitUntilComplete(ctx context.Context) (complete bool) {
	signal := make(chan bool)
	go func() {
		instance.complete.Lock()
		defer instance.complete.Unlock()
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
