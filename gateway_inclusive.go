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

type gatewayProbingReport struct {
	result []int
	flowId id.Id
}

func (m gatewayProbingReport) message() {}

type flowSync struct {
	response chan IAction
	flow     Flow
}

type inclusiveGateway struct {
	*wiring
	element                 *schema.InclusiveGateway
	mch                     chan imessage
	defaultSequenceFlow     *SequenceFlow
	nonDefaultSequenceFlows []*SequenceFlow
	probing                 *chan IAction
	activated               *flowSync
	awaiting                []id.Id
	arrived                 []id.Id
	sync                    []chan IAction
	flowTracker             *flowTracker
	synchronized            bool
}

func newInclusiveGateway(wr *wiring, element *schema.InclusiveGateway) (gw *inclusiveGateway, err error) {
	var defaultSequenceFlow *SequenceFlow

	if seqFlow, present := element.Default(); present {
		if node, found := wr.process.FindBy(schema.ExactId(*seqFlow).
			And(schema.ElementType((*schema.SequenceFlow)(nil)))); found {
			defaultSequenceFlow = new(SequenceFlow)
			*defaultSequenceFlow = MakeSequenceFlow(
				node.(*schema.SequenceFlow),
				wr.process,
			)
		} else {
			err = errors.NotFoundError{
				Expected: fmt.Sprintf("default sequence flow with ID %s", *seqFlow),
			}
			return nil, err
		}
	}

	nonDefaultSequenceFlows := allSequenceFlows(&wr.outgoing,
		func(sequenceFlow *SequenceFlow) bool {
			if defaultSequenceFlow == nil {
				return false
			}
			return *sequenceFlow == *defaultSequenceFlow
		},
	)

	gw = &inclusiveGateway{
		wiring:                  wr,
		element:                 element,
		mch:                     make(chan imessage, len(wr.incoming)*2+1),
		nonDefaultSequenceFlows: nonDefaultSequenceFlows,
		defaultSequenceFlow:     defaultSequenceFlow,
		flowTracker:             newFlowTracker(wr.tracer, element),
	}
	return
}

func (gw *inclusiveGateway) run(ctx context.Context, sender tracing.ISenderHandle) {
	defer gw.flowTracker.shutdown()
	activity := gw.flowTracker.activity()

	defer sender.Done()

	for {
		select {
		case msg := <-gw.mch:
			switch m := msg.(type) {
			case gatewayProbingReport:
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
						gw.wiring.tracer.Send(ErrorTrace{
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
			gw.tracer.Send(CancellationFlowNodeTrace{Node: gw.element})
			return
		}
	}
}

func (gw *inclusiveGateway) trySync() {
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
			gw.activated.response <- probeAction{
				sequenceFlows: gw.nonDefaultSequenceFlows,
				probeReport: func(indices []int) {
					gw.mch <- gatewayProbingReport{
						result: indices,
						flowId: anId,
					}
				},
			}

			gw.synchronized = true
		}
	}
}

func (gw *inclusiveGateway) NextAction(ctx context.Context, flow Flow) chan IAction {
	sender := gw.tracer.RegisterSender()
	go gw.run(ctx, sender)

	response := make(chan IAction)
	gw.mch <- nextActionMessage{response: response, flow: flow}
	return response
}

func (gw *inclusiveGateway) Element() schema.FlowNodeInterface {
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

func newFlowTracker(tracer tracing.ITracer, element *schema.InclusiveGateway) *flowTracker {
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
	go tracker.run()
	return &tracker
}

func (tracker *flowTracker) run() {
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
