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

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type ParallelGateway struct {
	*Wiring
	ctx                   context.Context
	element               *schema.ParallelGateway
	mch                   chan imessage
	reportedIncomingFlows int
	awaitingActions       []chan IAction
	noOfIncomingFlows     int
}

func NewParallelGateway(ctx context.Context, wiring *Wiring, parallelGateway *schema.ParallelGateway) (gateway *ParallelGateway, err error) {
	gateway = &ParallelGateway{
		Wiring:                wiring,
		ctx:                   ctx,
		element:               parallelGateway,
		mch:                   make(chan imessage, len(wiring.Incoming)*2+1),
		reportedIncomingFlows: 0,
		awaitingActions:       make([]chan IAction, 0),
		noOfIncomingFlows:     len(wiring.Incoming),
	}

	return
}

func (gw *ParallelGateway) flowWhenReady() {
	if gw.reportedIncomingFlows == gw.noOfIncomingFlows {
		gw.reportedIncomingFlows = 0
		awaitingActions := gw.awaitingActions
		gw.awaitingActions = make([]chan IAction, 0)
		sequences := AllSequenceFlows(&gw.Outgoing)
		distributeFlows(awaitingActions, sequences)
	}
}

func (gw *ParallelGateway) run(ctx context.Context, sender tracing.ISenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-gw.mch:
			switch m := msg.(type) {
			case nextActionMessage:
				gw.reportedIncomingFlows++
				gw.awaitingActions = append(gw.awaitingActions, m.response)
				gw.flowWhenReady()
				gw.Tracer.Send(IncomingFlowProcessedTrace{Node: gw.element, Flow: m.flow})
			default:
			}
		case <-ctx.Done():
			gw.Tracer.Send(CancellationFlowNodeTrace{Node: gw.element})
			return
		}
	}
}

func (gw *ParallelGateway) NextAction(flow Flow) chan IAction {
	sender := gw.Tracer.RegisterSender()
	go gw.run(gw.ctx, sender)

	response := make(chan IAction)
	gw.mch <- nextActionMessage{response: response, flow: flow}
	return response
}

func (gw *ParallelGateway) Element() schema.FlowNodeInterface {
	return gw.element
}

// IncomingFlowProcessedTrace signals that a particular flow
// has been processed. If any action has been taken, it has already happened
type IncomingFlowProcessedTrace struct {
	Node *schema.ParallelGateway
	Flow Flow
}

func (t IncomingFlowProcessedTrace) Unpack() any { return t.Node }
