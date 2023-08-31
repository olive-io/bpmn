package task

import (
	"context"
	"sync"

	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/flow_node/activity"
	"github.com/olive-io/bpmn/schema"
)

type message interface {
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

type Task struct {
	*flow_node.Wiring
	element       *schema.Task
	runnerChannel chan message
	bodyLock      sync.RWMutex
	body          func(*Task, context.Context) flow_node.IAction
	cancel        context.CancelFunc
}

// SetBody override Task's body with an arbitrary function
//
// Since Task implements Abstract Task, it does nothing by default.
// This allows to add an implementation. Primarily used for testing.
func (node *Task) SetBody(body func(*Task, context.Context) flow_node.IAction) {
	node.bodyLock.Lock()
	defer node.bodyLock.Unlock()
	node.body = body
}

func NewTask(ctx context.Context, startEvent *schema.Task) activity.Constructor {
	return func(wiring *flow_node.Wiring) (node activity.Activity, err error) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		taskNode := &Task{
			Wiring:        wiring,
			element:       startEvent,
			runnerChannel: make(chan message, len(wiring.Incoming)*2+1),
			cancel:        cancel,
		}
		go taskNode.runner(ctx)
		node = taskNode
		return
	}
}

func (node *Task) runner(ctx context.Context) {
	for {
		select {
		case msg := <-node.runnerChannel:
			switch m := msg.(type) {
			case cancelMessage:
				node.cancel()
				m.response <- true
			case nextActionMessage:
				go func() {
					var action flow_node.IAction
					action = flow_node.FlowAction{SequenceFlows: flow_node.AllSequenceFlows(&node.Outgoing)}
					if node.body != nil {
						node.bodyLock.RLock()
						action = node.body(node, ctx)
						node.bodyLock.RUnlock()
					}
					m.response <- action
				}()
			default:
			}
		case <-ctx.Done():
			return
		}
	}
}

func (node *Task) NextAction(flow_interface.T) chan flow_node.IAction {
	response := make(chan flow_node.IAction)
	node.runnerChannel <- nextActionMessage{response: response}
	return response
}

func (node *Task) Element() schema.FlowNodeInterface {
	return node.element
}

func (node *Task) Cancel() <-chan bool {
	response := make(chan bool)
	node.runnerChannel <- cancelMessage{response: response}
	return response
}
