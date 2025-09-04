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
	"sync"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type parallelGateway struct {
	*wiring
	element               *schema.ParallelGateway
	mch                   chan imessage
	reportedIncomingFlows int
	once                  sync.Once
	awaitingActions       []chan IAction
	noOfIncomingFlows     int
}

func newParallelGateway(wr *wiring, element *schema.ParallelGateway) (gw *parallelGateway, err error) {
	gw = &parallelGateway{
		wiring:                wr,
		element:               element,
		mch:                   make(chan imessage, len(wr.incoming)*2+1),
		reportedIncomingFlows: 0,
		awaitingActions:       make([]chan IAction, 0),
		noOfIncomingFlows:     len(wr.incoming),
	}

	return
}

func (gw *parallelGateway) flowWhenReady() {
	if gw.reportedIncomingFlows == gw.noOfIncomingFlows {
		gw.reportedIncomingFlows = 0
		awaitingActions := gw.awaitingActions
		gw.awaitingActions = make([]chan IAction, 0)
		sequences := allSequenceFlows(&gw.outgoing)
		distributeFlows(awaitingActions, sequences)
	}
}

func (gw *parallelGateway) run(ctx context.Context, sender tracing.ISenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-gw.mch:
			switch m := msg.(type) {
			case nextActionMessage:
				gw.reportedIncomingFlows++
				gw.awaitingActions = append(gw.awaitingActions, m.response)
				gw.flowWhenReady()
				gw.tracer.Send(IncomingFlowProcessedTrace{Node: gw.element, Flow: m.flow})
			default:
			}
		case <-ctx.Done():
			gw.tracer.Send(CancellationFlowNodeTrace{Node: gw.element})
			return
		}
	}
}

func (gw *parallelGateway) NextAction(ctx context.Context, flow Flow) chan IAction {
	gw.once.Do(func() {
		sender := gw.tracer.RegisterSender()
		go gw.run(ctx, sender)
	})

	response := make(chan IAction)
	gw.mch <- nextActionMessage{response: response, flow: flow}
	return response
}

func (gw *parallelGateway) Element() schema.FlowNodeInterface {
	return gw.element
}

// IncomingFlowProcessedTrace signals that a particular flow
// has been processed. If any action has been taken, it has already happened
type IncomingFlowProcessedTrace struct {
	Node *schema.ParallelGateway
	Flow Flow
}

func (t IncomingFlowProcessedTrace) Unpack() any { return t.Node }
