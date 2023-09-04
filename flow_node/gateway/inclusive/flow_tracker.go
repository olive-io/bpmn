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

package inclusive

import (
	"context"
	"sync"

	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tools/id"
	"github.com/olive-io/bpmn/tracing"
)

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
	// This is important because if the node will invoke
	// `activeFlowsInCohort` before the tracker has caught up,
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
	case flow.Trace:
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
	case flow.TerminationTrace:
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
