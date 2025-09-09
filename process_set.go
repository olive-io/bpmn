/*
   Copyright 2025 The bpmn Authors

   This program is offered under a commercial and under the AGPL license.
   For AGPL licensing, see below.

   AGPL licensing:
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package bpmn

import (
	"context"
	"fmt"
	"sync"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/data"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/id"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type ProcessSet struct {
	*Options
	wg sync.WaitGroup

	sourceOptions []Option

	executes []*Process
	waitings []*schema.Process

	definitions *schema.Definitions

	sourceMessageFlows map[string]*schema.MessageFlow

	cmu     sync.RWMutex
	catchCh map[string]chan struct{}

	subTracer tracing.ITracer

	mch  chan imessage
	done chan struct{}
}

func NewProcessSet(executeProcesses, waitingProcesses []*schema.Process, definitions *schema.Definitions, opts ...Option) (*ProcessSet, error) {
	options := NewOptions(opts...)

	var err error
	ctx := options.ctx
	tracer := options.tracer
	if tracer == nil {
		options.tracer = tracing.NewTracer(ctx)
		opts = append(opts, WithTracer(options.tracer))
	}

	if options.idGenerator == nil {
		options.idGenerator, err = id.GetSno().NewIdGenerator(ctx, tracer)
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithIdGenerator(options.idGenerator))
	}
	if options.locator == nil {
		options.locator = data.NewFlowDataLocator()
		opts = append(opts, WithLocator(options.locator))
	}

	executes := make([]*Process, 0)
	for _, executeProcess := range executeProcesses {
		var process *Process

		subTracer := tracing.NewTracer(ctx)
		tracing.NewRelay(ctx, subTracer, options.tracer, func(trace tracing.ITrace) []tracing.ITrace {
			return []tracing.ITrace{trace}
		})

		process, err = NewProcess(executeProcess, definitions, append(opts, WithTracer(subTracer))...)
		if err != nil {
			return nil, fmt.Errorf("create new process: %w", err)
		}
		executes = append(executes, process)
	}

	sourceMessageFlows := make(map[string]*schema.MessageFlow)
	for _, collaboration := range *definitions.Collaborations() {
		for _, msg := range *collaboration.MessageFlows() {
			sourceMessageFlows[string(msg.SourceRefField)] = &msg
		}
	}

	ps := &ProcessSet{
		Options:            options,
		sourceOptions:      opts,
		executes:           executes,
		waitings:           waitingProcesses,
		definitions:        definitions,
		sourceMessageFlows: sourceMessageFlows,

		catchCh: make(map[string]chan struct{}),

		mch:  make(chan imessage, len(executes)+1),
		done: make(chan struct{}, 1),
	}

	return ps, nil
}

func (ps *ProcessSet) Tracer() tracing.ITracer { return ps.tracer }

func (ps *ProcessSet) Locator() data.IFlowDataLocator { return ps.locator }

func (ps *ProcessSet) StartAll(ctx context.Context) error {
	go ps.run(ctx)

	for _, process := range ps.executes {
		err := process.StartAll(ctx)
		if err != nil {
			return fmt.Errorf("start process %s: %w", process.Id().String(), err)
		}

		ps.wg.Add(1)
		go ps.tracerProcess(ctx, process, &ps.wg)
	}

	return nil
}

// WaitUntilComplete waits until the instance is complete.
// Returns true if the instance was complete, false if the context signaled `Done`
func (ps *ProcessSet) WaitUntilComplete(ctx context.Context) (complete bool) {
	go func() {
		ps.wg.Wait()
		close(ps.done)
	}()
	select {
	case <-ctx.Done():
		complete = false
	case <-ps.done:
		complete = true
	}
	return
}

func (ps *ProcessSet) run(ctx context.Context) {
	for {
		select {
		case ch := <-ps.mch:
			switch msg := ch.(type) {
			case throwMessage:
				sourceRef, ok := ps.sourceMessageFlows[msg.Id]
				if ok {
					startFlowNode, waitingProcess, found := ps.resolveWaitingProcessAndEvent(string(sourceRef.TargetRefField))
					if found {
						// flow nodes
						subTracer := tracing.NewTracer(ctx)
						tracing.NewRelay(ctx, subTracer, ps.tracer, func(trace tracing.ITrace) []tracing.ITrace {
							return []tracing.ITrace{trace}
						})

						process, err := NewProcess(waitingProcess, ps.definitions, append(ps.sourceOptions, WithTracer(subTracer))...)
						if err != nil {
							ps.tracer.Send(ErrorTrace{Error: err})
							continue
						}

						err = process.StartWith(ctx, startFlowNode)
						if err != nil {
							ps.tracer.Send(ErrorTrace{Error: err})
							continue
						}
						ps.wg.Add(1)
						go ps.tracerProcess(ctx, process, &ps.wg)
					}
					trigger, found := ps.triggerCatch(string(sourceRef.TargetRefField))
					if found {
						trigger()
					}
				}
			}
		case <-ps.done:
			ps.tracer.Send(CeaseProcessSetTrace{Definitions: ps.definitions})
			return
		case <-ctx.Done():
			return
		}
	}
}

func (ps *ProcessSet) tracerProcess(ctx context.Context, process *Process, wg *sync.WaitGroup) {
	defer wg.Done()

	traces := process.Tracer().Subscribe()
	defer process.tracer.Unsubscribe(traces)

LOOP:
	for {
		var trace tracing.ITrace
		select {
		case trace = <-traces:
		case <-ctx.Done():
			return
		default:
			continue
		}

		trace = tracing.Unwrap(trace)
		switch msg := trace.(type) {
		case FlowTrace:
			switch evt := msg.Source.(type) {
			case *schema.ThrowEvent:
				eventId, ok := evt.Id()
				if ok {
					ps.mch <- throwMessage{Id: *eventId}
				}
			}
		case ActiveListeningTrace:
			ready := make(chan struct{}, 1)
			ps.registerCatch(msg.Node, ready)
			go func() {
				wg.Add(1)
				defer wg.Done()

				for {
					select {
					case <-ready:
						for _, eventDefinition := range msg.Node.SignalEventDefinitionField {
							ref, ok := eventDefinition.SignalRef()
							if ok {
								process.ConsumeEvent(event.NewSignalEvent(string(*ref)))
							}
						}
						for _, eventDefinition := range msg.Node.MessageEventDefinitionField {
							ref, ok := eventDefinition.MessageRef()
							if ok {
								process.ConsumeEvent(event.NewMessageEvent(string(*ref), (*string)(eventDefinition.OperationRefField)))
							}
						}
						return
					case <-ctx.Done():
						return
					}
				}
			}()
		case CeaseFlowTrace:
			break LOOP
		}
	}

	return
}

func (ps *ProcessSet) registerCatch(node *schema.CatchEvent, ch chan struct{}) {
	idPtr, ok := node.Id()
	if !ok {
		return
	}
	ps.cmu.Lock()
	defer ps.cmu.Unlock()
	ps.catchCh[*idPtr] = ch
}

func (ps *ProcessSet) triggerCatch(id string) (func(), bool) {
	ps.cmu.RLock()
	defer ps.cmu.RUnlock()

	ch, ok := ps.catchCh[id]
	if !ok {
		return nil, false
	}

	cancel := func() {
		ps.cmu.Lock()
		close(ch)
		delete(ps.catchCh, id)
		ps.cmu.Unlock()
	}
	return cancel, true
}

func (ps *ProcessSet) resolveWaitingProcessAndEvent(idRef string) (schema.FlowNodeInterface, *schema.Process, bool) {
	for _, process := range ps.waitings {
		element, found := process.FindBy(schema.ExactId(idRef))
		if found {
			node, ok := element.(schema.FlowNodeInterface)
			if ok {
				return node, process, true
			}
		}
	}
	return nil, nil, false
}

type throwMessage struct {
	Id string
}

func (msg throwMessage) message() {}
