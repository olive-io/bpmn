package catch_event

import (
	"context"
	"testing"

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow"
	ev "github.com/olive-io/bpmn/flow_node/event/catch"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/test"
	"github.com/olive-io/bpmn/tracing"

	"github.com/stretchr/testify/assert"
)

func TestSignalEvent(t *testing.T) {
	testEvent(t, "sample/catch_event/intermediate_catch_event.bpmn", "signalCatch", nil, false, event.NewSignalEvent("global_sig1"))
}

func TestMessageEvent(t *testing.T) {
	testEvent(t, "sample/catch_event/intermediate_catch_event.bpmn", "messageCatch", nil, false, event.NewMessageEvent("msg", nil))
}

func TestMultipleEvent(t *testing.T) {
	// either
	testEvent(t, "sample/catch_event/intermediate_catch_event_multiple.bpmn", "multipleCatch", nil, false, event.NewMessageEvent("msg", nil))
	// or
	testEvent(t, "sample/catch_event/intermediate_catch_event_multiple.bpmn", "multipleCatch", nil, false, event.NewSignalEvent("global_sig1"))
}

func TestMultipleParallelEvent(t *testing.T) {
	// both
	testEvent(t, "sample/catch_event/intermediate_catch_event_multiple_parallel.bpmn", "multipleParallelCatch",
		nil, false, event.NewMessageEvent("msg", nil), event.NewSignalEvent("global_sig1"))
	// either
	testEvent(t, "sample/catch_event/intermediate_catch_event_multiple_parallel.bpmn", "multipleParallelCatch", nil, true, event.NewMessageEvent("msg", nil))
	testEvent(t, "sample/catch_event/intermediate_catch_event_multiple_parallel.bpmn", "multipleParallelCatch", nil, true, event.NewSignalEvent("global_sig1"))
}

type eventInstance struct {
	id string
}

func (e eventInstance) EventDefinition() schema.EventDefinitionInterface {
	definition := schema.DefaultEventDefinition()
	return &definition
}

type eventDefinitionInstanceBuilder struct{}

func (e eventDefinitionInstanceBuilder) NewEventDefinitionInstance(def schema.EventDefinitionInterface) (event.DefinitionInstance, error) {
	switch d := def.(type) {
	case *schema.TimerEventDefinition:
		id, _ := d.Id()
		return eventInstance{id: *id}, nil
	case *schema.ConditionalEventDefinition:
		id, _ := d.Id()
		return eventInstance{id: *id}, nil
	default:
		return event.WrapEventDefinition(d), nil
	}
}

func TestTimerEvent(t *testing.T) {
	i := eventInstance{id: "timer_1"}
	b := eventDefinitionInstanceBuilder{}
	testEvent(t, "sample/catch_event/intermediate_catch_event.bpmn", "timerCatch", &b, false, event.MakeTimerEvent(i))
}

func TestConditionalEvent(t *testing.T) {
	i := eventInstance{id: "conditional_1"}
	b := eventDefinitionInstanceBuilder{}
	testEvent(t, "sample/catch_event/intermediate_catch_event.bpmn", "conditionalCatch", &b, false, event.MakeTimerEvent(i))
}

func testEvent(t *testing.T, filename string, nodeId string, eventDefinitionInstanceBuilder event.DefinitionInstanceBuilder, eventObservationOnly bool, events ...event.Event) {
	var testDoc schema.Definitions
	test.LoadTestFile(filename, &testDoc)
	processElement := (*testDoc.Processes())[0]
	proc := process.New(&processElement, &testDoc, process.WithEventDefinitionInstanceBuilder(eventDefinitionInstanceBuilder))

	tracer := tracing.NewTracer(context.Background())
	traces := tracer.SubscribeChannel(make(chan tracing.Trace, 64))

	if inst, err := proc.Instantiate(instance.WithTracer(tracer)); err == nil {
		err := inst.StartAll(context.Background())
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
		resultChan := make(chan bool)
		go func() {
			for {
				trace := tracing.Unwrap(<-traces)
				switch trace := trace.(type) {
				case ev.ActiveListeningTrace:
					if id, present := trace.Node.Id(); present {
						if *id == nodeId {
							// listening
							resultChan <- true
							return
						}

					}
				case tracing.ErrorTrace:
					t.Errorf("%#v", trace)
					resultChan <- false
					return
				default:
					t.Logf("%#v", trace)
				}
			}
		}()

		assert.True(t, <-resultChan)

		go func() {
			defer inst.Tracer.Unsubscribe(traces)
			eventsToObserve := events
			for {
				trace := tracing.Unwrap(<-traces)
				switch trace := trace.(type) {
				case ev.EventObservedTrace:
					if eventObservationOnly {
						for i := range eventsToObserve {
							if eventsToObserve[i] == trace.Event {
								eventsToObserve[i] = eventsToObserve[len(eventsToObserve)-1]
								eventsToObserve = eventsToObserve[:len(eventsToObserve)-1]
								break
							}
						}
						if len(eventsToObserve) == 0 {
							resultChan <- true
							return
						}
					}
				case flow.Trace:
					if id, present := trace.Source.Id(); present {
						if *id == nodeId {
							// success!
							resultChan <- true
							return
						}

					}
				case tracing.ErrorTrace:
					t.Errorf("%#v", trace)
					resultChan <- false
					return
				default:
					t.Logf("%#v", trace)
				}
			}
		}()

		for _, evt := range events {
			_, err = inst.ConsumeEvent(evt)
			assert.Nil(t, err)
		}

		assert.True(t, <-resultChan)

	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}

}
