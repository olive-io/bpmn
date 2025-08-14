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
	"github.com/olive-io/bpmn/v2/pkg/errors"
	"github.com/olive-io/bpmn/v2/pkg/id"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type InclusiveNoEffectiveSequenceFlows struct {
	*schema.InclusiveGateway
}

func (e InclusiveNoEffectiveSequenceFlows) Error() string {
	ownId := "<unnamed>"
	if ownIdPtr, present := e.InclusiveGateway.Id(); present {
		ownId = *ownIdPtr
	}
	return fmt.Sprintf("No effective sequence flows found in exclusive gateway `%v`", ownId)
}

type probingReport struct {
	result []int
	flowId id.Id
}

func (m probingReport) message() {}

type flowSync struct {
	response chan IAction
	flow     Flow
}

type InclusiveGateway struct {
	*Wiring
	ctx                     context.Context
	element                 *schema.InclusiveGateway
	mch                     chan imessage
	defaultSequenceFlow     *SequenceFlow
	nonDefaultSequenceFlows []*SequenceFlow
	probing                 *chan IAction
	activated               *flowSync
	awaiting                []id.Id
	arrived                 []id.Id
	sync                    []chan IAction
	*flowTracker
	synchronized bool
}

func NewInclusiveGateway(ctx context.Context, wiring *Wiring, inclusiveGateway *schema.InclusiveGateway) (gw *InclusiveGateway, err error) {
	var defaultSequenceFlow *SequenceFlow

	if seqFlow, present := inclusiveGateway.Default(); present {
		if node, found := wiring.Process.FindBy(schema.ExactId(*seqFlow).
			And(schema.ElementType((*schema.SequenceFlow)(nil)))); found {
			defaultSequenceFlow = new(SequenceFlow)
			*defaultSequenceFlow = MakeSequenceFlow(
				node.(*schema.SequenceFlow),
				wiring.Process,
			)
		} else {
			err = errors.NotFoundError{
				Expected: fmt.Sprintf("default sequence flow with ID %s", *seqFlow),
			}
			return nil, err
		}
	}

	nonDefaultSequenceFlows := AllSequenceFlows(&wiring.Outgoing,
		func(sequenceFlow *SequenceFlow) bool {
			if defaultSequenceFlow == nil {
				return false
			}
			return *sequenceFlow == *defaultSequenceFlow
		},
	)

	gw = &InclusiveGateway{
		Wiring:                  wiring,
		ctx:                     ctx,
		element:                 inclusiveGateway,
		mch:                     make(chan imessage, len(wiring.Incoming)*2+1),
		nonDefaultSequenceFlows: nonDefaultSequenceFlows,
		defaultSequenceFlow:     defaultSequenceFlow,
		flowTracker:             newFlowTracker(ctx, wiring.Tracer, inclusiveGateway),
	}
	return
}

func (gw *InclusiveGateway) run(ctx context.Context, sender tracing.ISenderHandle) {
	defer gw.flowTracker.shutdown()
	activity := gw.flowTracker.activity()

	defer sender.Done()

	for {
		select {
		case msg := <-gw.mch:
			switch m := msg.(type) {
			case probingReport:
				response := gw.probing
				if response == nil {
					// Reschedule, there's no next action yet
					go func() {
						gw.mch <- m
					}()
					continue
				}
				gw.probing = nil
				sfs := make([]*SequenceFlow, 0)
				for _, i := range m.result {
					sfs = append(sfs, gw.nonDefaultSequenceFlows[i])
				}

				switch len(sfs) {
				case 0:
					// no successful non-default sequence flows
					if gw.defaultSequenceFlow == nil {
						// exception (Table 13.2)
						gw.Wiring.Tracer.Send(ErrorTrace{
							Error: InclusiveNoEffectiveSequenceFlows{
								InclusiveGateway: gw.element,
							},
						})
					} else {
						distributeFlows(gw.sync, []*SequenceFlow{gw.defaultSequenceFlow})
					}
				default:
					distributeFlows(gw.sync, sfs)
				}
				gw.synchronized = false
				gw.activated = nil
			case nextActionMessage:
				if gw.synchronized {
					if m.flow.Id() == gw.activated.flow.Id() {
						// Activating flow returned
						gw.sync = append(gw.sync, m.response)
						gw.probing = &m.response
						// and now we wait until the probe has returned
					}
				} else {
					if gw.activated == nil {
						// Haven't been activated yet
						gw.activated = &flowSync{response: m.response, flow: m.flow}
						gw.awaiting = gw.flowTracker.activeFlowsInCohort(m.flow.Id())
						gw.arrived = []id.Id{m.flow.Id()}
						gw.sync = make([]chan IAction, 0)
					} else {
						// Already activated
						gw.arrived = append(gw.arrived, m.flow.Id())
						gw.sync = append(gw.sync, m.response)
					}
					gw.trySync()
				}

			default:
			}
		case <-activity:
			if !gw.synchronized && gw.activated != nil {
				gw.awaiting = gw.flowTracker.activeFlowsInCohort(gw.activated.flow.Id())
				gw.trySync()
			}
		case <-ctx.Done():
			gw.Tracer.Send(CancellationFlowNodeTrace{Node: gw.element})
			return
		}
	}
}

func (gw *InclusiveGateway) trySync() {
	if !gw.synchronized && len(gw.arrived) >= len(gw.awaiting) {
		// Have we got everybody?
		matches := 0
		for i := range gw.arrived {
			for j := range gw.awaiting {
				if gw.awaiting[j] == gw.arrived[i] {
					matches++
				}
			}
		}
		if matches == len(gw.awaiting) {
			anId := gw.activated.flow.Id()
			// Probe outgoing sequence flow using the first flow
			gw.activated.response <- ProbeAction{
				SequenceFlows: gw.nonDefaultSequenceFlows,
				ProbeReport: func(indices []int) {
					gw.mch <- probingReport{
						result: indices,
						flowId: anId,
					}
				},
			}

			gw.synchronized = true
		}
	}
}

func (gw *InclusiveGateway) NextAction(flow Flow) chan IAction {
	sender := gw.Tracer.RegisterSender()
	go gw.run(gw.ctx, sender)

	response := make(chan IAction)
	gw.mch <- nextActionMessage{response: response, flow: flow}
	return response
}

func (gw *InclusiveGateway) Element() schema.FlowNodeInterface {
	return gw.element
}

type flowTracker struct {
	traces     <-chan tracing.ITrace
	shutdownCh chan bool
	flows      map[id.Id]schema.Id
	activityCh chan struct{}
	lock       sync.RWMutex
	element    *schema.InclusiveGateway
}

func (tracker *flowTracker) activity() <-chan struct{} {
	return tracker.activityCh
}

func newFlowTracker(ctx context.Context, tracer tracing.ITracer, element *schema.InclusiveGateway) *flowTracker {
	tracker := flowTracker{
		traces:     tracer.Subscribe(),
		shutdownCh: make(chan bool),
		flows:      make(map[id.Id]schema.Id),
		activityCh: make(chan struct{}),
		element:    element,
	}
	// Lock the tracker until it has caught up enough
	// to see the incoming flow for the node
	tracker.lock.Lock()
	go tracker.run(ctx)
	return &tracker
}

func (tracker *flowTracker) run(ctx context.Context) {
	// As per note in the constructor, we're starting in a locked mode
	locked := true
	// Flag for notifying the node about activity
	notify := false
	// Indicates whether the tracker has observed a flow
	// that reaches the node that uses this tracker.
	// This is important because if the node invokes
	// `activeFlowsInCohort`, before the tracker has caught up,
	// it'll return an empty list, and the node will assume that
	// there's no other flow to wait for, and will proceed (which
	// is incorrect)
	reachedNode := false
	for {
		select {
		case trace := <-tracker.traces:
			locked, notify, reachedNode = tracker.handleTrace(locked, trace, notify, reachedNode)
			// continue draining
			continue
		case <-tracker.shutdownCh:
			if locked {
				tracker.lock.Unlock()
			}
			return
		case <-ctx.Done():
			if locked {
				tracker.lock.Unlock()
			}
			return
		default:
			// Nothing else is coming in, unlock if locked
			if locked && reachedNode {
				tracker.lock.Unlock()
				if notify {
					tracker.activityCh <- struct{}{}
					notify = false
				}
				locked = false
			}
			// and now proceed with the second select to wait
			// for an event without doing busy work (this `default` clause)
		}
		select {
		case trace := <-tracker.traces:
			locked, notify, reachedNode = tracker.handleTrace(locked, trace, notify, reachedNode)
		case <-tracker.shutdownCh:
			if locked {
				tracker.lock.Unlock()
			}
			return
		case <-ctx.Done():
			if locked {
				tracker.lock.Unlock()
			}
			return
		}

	}
}

func (tracker *flowTracker) handleTrace(locked bool, trace tracing.ITrace, notify bool, reachedNode bool) (bool, bool, bool) {
	trace = tracing.Unwrap(trace)
	if !locked {
		// Lock tracker records until messages are drained
		tracker.lock.Lock()
		locked = true
	}
	switch t := trace.(type) {
	case FlowTrace:
		for _, snapshot := range t.Flows {
			// If we haven't reached the node
			if !reachedNode {
				// Try and see if this flow is the one that goes into it
				targetId := snapshot.SequenceFlow().TargetRef()
				if idPtr, present := tracker.element.Id(); present {
					reachedNode = *idPtr == *targetId
				}
			}
			if idPtr, present := t.Source.Id(); present {
				_, ok := tracker.flows[snapshot.Id()]
				_, isInclusive := t.Source.(*schema.InclusiveGateway)
				if !ok || isInclusive {
					tracker.flows[snapshot.Id()] = *idPtr
				}
			}
		}
		notify = true
	case TerminationTrace:
		delete(tracker.flows, t.FlowId)
		notify = true
	}
	return locked, notify, reachedNode
}

func (tracker *flowTracker) shutdown() {
	close(tracker.shutdownCh)
}

func (tracker *flowTracker) activeFlowsInCohort(flowId id.Id) (result []id.Id) {
	result = make([]id.Id, 0)
	tracker.lock.RLock()
	defer tracker.lock.RUnlock()
	if location, ok := tracker.flows[flowId]; ok {
		for k, v := range tracker.flows {
			if v == location {
				result = append(result, k)
			}
		}
	}
	return
}
