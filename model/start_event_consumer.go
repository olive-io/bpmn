package model

import (
	"context"
	"sync"

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tools/logic"
	"github.com/olive-io/bpmn/tracing"
)

type startEventConsumer struct {
	process              *process.Process
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
	process *process.Process,
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
	defer s.tracer.Trace(EventInstantiationAttemptedTrace{Event: ev, Element: s.element})

	if satisfied, chain := s.satisfier.Satisfy(ev); satisfied {
		// If it's a new chain, add new event buffer
		if chain > len(s.events)-1 {
			s.events = append(s.events, []event.IEvent{ev})
		}
		var inst *instance.Instance
		inst, err = s.process.Instantiate(
			instance.WithContext(s.ctx),
			instance.WithTracer(s.tracer),
			instance.WithEventDefinitionInstanceBuilder(event.DefinitionInstanceBuildingChain(
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
