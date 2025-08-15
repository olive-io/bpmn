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
	defer s.tracer.Send(EventInstantiationAttemptedTrace{Event: ev, Node: s.element})

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
