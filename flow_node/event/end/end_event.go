package end

import (
	"context"

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tracing"
)

type message interface {
	message()
}

type nextActionMessage struct {
	response chan flow_node.Action
}

func (m nextActionMessage) message() {}

type Node struct {
	*flow_node.Wiring
	element              *schema.EndEvent
	activated            bool
	completed            bool
	runnerChannel        chan message
	startEventsActivated []*schema.StartEvent
}

func New(ctx context.Context, wiring *flow_node.Wiring, endEvent *schema.EndEvent) (node *Node, err error) {
	node = &Node{
		Wiring:               wiring,
		element:              endEvent,
		activated:            false,
		completed:            false,
		runnerChannel:        make(chan message, len(wiring.Incoming)*2+1),
		startEventsActivated: make([]*schema.StartEvent, 0),
	}
	sender := node.Tracer.RegisterSender()
	go node.runner(ctx, sender)
	return
}

func (node *Node) runner(ctx context.Context, sender tracing.SenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-node.runnerChannel:
			switch m := msg.(type) {
			case nextActionMessage:
				if !node.activated {
					node.activated = true
				}
				// If the node already completed, then we essentially fuse it
				if node.completed {
					m.response <- flow_node.CompleteAction{}
					continue
				}

				if _, err := node.EventIngress.ConsumeEvent(
					event.MakeEndEvent(node.element),
				); err == nil {
					node.completed = true
					m.response <- flow_node.CompleteAction{}
				} else {
					node.Wiring.Tracer.Trace(tracing.ErrorTrace{Error: err})
				}
			default:
			}
		case <-ctx.Done():
			node.Tracer.Trace(flow_node.CancellationTrace{Node: node.element})
			return
		}
	}
}

func (node *Node) NextAction(flow_interface.T) chan flow_node.Action {
	response := make(chan flow_node.Action)
	node.runnerChannel <- nextActionMessage{response: response}
	return response
}

func (node *Node) Element() schema.FlowNodeInterface {
	return node.element
}
