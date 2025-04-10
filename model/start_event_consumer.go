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

package model

import (
	"context"
	"sync"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/logic"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type startEventConsumer struct {
	process              *bpmn.Process
	parallel             bool
	ctx                  context.Context
	consumptionLock      sync.Mutex
	tracer               tracing.ITracer
	events               [][]event.IEvent
	element              schema.CatchEventInterface
	satisfier            *logic.CatchEventSatisfier
	eventInstanceBuilder event.IDefinitionInstanceBuilder
}

func (s *startEventConsumer) NewEventDefinitionInstance(
	def schema.EventDefinitionInterface,
) (definitionInstance event.IDefinitionInstance, err error) {
	instances := s.satisfier.EventDefinitionInstances()
	for i := range *instances {
		if schema.Equal((*instances)[i].EventDefinition(), def) {
			definitionInstance = (*instances)[i]
			break
		}
	}
	return
}

func newStartEventConsumer(
	ctx context.Context,
	tracer tracing.ITracer,
	process *bpmn.Process,
	startEvent *schema.StartEvent,
	eventDefinitionInstanceBuilder event.IDefinitionInstanceBuilder) *startEventConsumer {
	consumer := &startEventConsumer{
		ctx:                  ctx,
		process:              process,
		parallel:             startEvent.ParallelMultiple(),
		tracer:               tracer,
		events:               make([][]event.IEvent, 0, len(startEvent.EventDefinitions())),
		element:              startEvent,
		satisfier:            logic.NewCatchEventSatisfier(startEvent, eventDefinitionInstanceBuilder),
		eventInstanceBuilder: eventDefinitionInstanceBuilder,
	}
	return consumer
}

func (s *startEventConsumer) ConsumeEvent(ev event.IEvent) (result event.ConsumptionResult, err error) {
	s.consumptionLock.Lock()
	defer s.consumptionLock.Unlock()
	defer s.tracer.Trace(EventInstantiationAttemptedTrace{Event: ev, Node: s.element})

	if satisfied, chain := s.satisfier.Satisfy(ev); satisfied {
		// If it's a new chain, add new event buffer
		if chain > len(s.events)-1 {
			s.events = append(s.events, []event.IEvent{ev})
		}
		var inst *bpmn.Instance
		inst, err = s.process.Instantiate(
			bpmn.WithContext(s.ctx),
			bpmn.WithTracer(s.tracer),
			bpmn.WithEventDefinitionInstanceBuilder(event.DefinitionInstanceBuildingChain(
				s, // this will pass-through already existing event definition instance from this execution
				s.eventInstanceBuilder,
			)),
		)
		if err != nil {
			result = event.ConsumptionError
			return
		}
		for _, ev := range s.events[chain] {
			result, err = inst.ConsumeEvent(ev)
			if err != nil {
				result = event.ConsumptionError
				return
			}
		}
		// Remove events buffer
		s.events[chain] = s.events[len(s.events)-1]
		s.events = s.events[:len(s.events)-1]
	} else if chain != logic.EventDidNotMatch {
		// If there was a match
		// If it's a new chain, add new event buffer
		if chain > len(s.events)-1 {
			s.events = append(s.events, []event.IEvent{ev})
		} else {
			s.events[chain] = append(s.events[chain], ev)
		}
	}
	return
}
