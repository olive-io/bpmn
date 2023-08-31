package start

import (
	"context"

	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/id"
	"github.com/olive-io/bpmn/logic"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tracing"
)

type message interface {
	message()
}

type nextActionMessage struct {
	response chan flow_node.IAction
}

func (m nextActionMessage) message() {}

type startMessage struct{}

func (m startMessage) message() {}

type eventMessage struct {
	event event.IEvent
}

func (m eventMessage) message() {}

type Node struct {
	*flow_node.Wiring
	element          *schema.StartEvent
	runnerChannel    chan message
	activated        bool
	idGenerator      id.IGenerator
	itemAwareLocator data.IItemAwareLocator
	satisfier        *logic.CatchEventSatisfier
}

func New(ctx context.Context, wiring *flow_node.Wiring, startEvent *schema.StartEvent,
	idGenerator id.IGenerator, itemAwareLocator data.IItemAwareLocator,
) (node *Node, err error) {
	eventDefinitions := startEvent.EventDefinitions()
	eventInstances := make([]event.IDefinitionInstance, len(eventDefinitions))

	for i, eventDefinition := range eventDefinitions {
		var instance event.IDefinitionInstance
		instance, err = wiring.EventDefinitionInstanceBuilder.NewEventDefinitionInstance(eventDefinition)
		if err != nil {
			return
		}
		eventInstances[i] = instance
	}

	node = &Node{
		Wiring:           wiring,
		element:          startEvent,
		runnerChannel:    make(chan message, len(wiring.Incoming)*2+1),
		activated:        false,
		idGenerator:      idGenerator,
		itemAwareLocator: itemAwareLocator,
		satisfier:        logic.NewCatchEventSatisfier(startEvent, wiring.EventDefinitionInstanceBuilder),
	}
	sender := node.Tracer.RegisterSender()
	go node.runner(ctx, sender)
	err = node.EventEgress.RegisterEventConsumer(node)
	if err != nil {
		return
	}
	return
}

func (node *Node) runner(ctx context.Context, sender tracing.ISenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-node.runnerChannel:
			switch m := msg.(type) {
			case nextActionMessage:
				if !node.activated {
					node.activated = true
					m.response <- flow_node.FlowAction{SequenceFlows: flow_node.AllSequenceFlows(&node.Outgoing)}
				} else {
					m.response <- flow_node.CompleteAction{}
				}
			case startMessage:
				node.flow(ctx)
			case eventMessage:
				if !node.activated {
					if satisfied, _ := node.satisfier.Satisfy(m.event); satisfied {
						node.flow(ctx)
					}
				}
			default:
			}
		case <-ctx.Done():
			node.Tracer.Trace(flow_node.CancellationTrace{Node: node.element})
			return
		}
	}
}

func (node *Node) flow(ctx context.Context) {
	newFlow := flow.New(node.Wiring.Definitions, node, node.Wiring.Tracer,
		node.Wiring.FlowNodeMapping, node.Wiring.FlowWaitGroup, node.idGenerator, nil,
		node.itemAwareLocator,
	)
	newFlow.Start(ctx)
}

func (node *Node) ConsumeEvent(ev event.IEvent) (result event.ConsumptionResult, err error) {
	node.runnerChannel <- eventMessage{event: ev}
	result = event.Consumed
	return
}

func (node *Node) Trigger() {
	node.runnerChannel <- startMessage{}
}

func (node *Node) NextAction(flow_interface.T) chan flow_node.IAction {
	response := make(chan flow_node.IAction)
	node.runnerChannel <- nextActionMessage{response: response}
	return response
}

func (node *Node) Element() schema.FlowNodeInterface {
	return node.element
}
