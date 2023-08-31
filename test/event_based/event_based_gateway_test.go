package event_based

import (
	"context"
	"testing"

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow"
	ev "github.com/olive-io/bpmn/flow_node/event/catch"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/test"
	"github.com/olive-io/bpmn/tracing"
	"github.com/stretchr/testify/assert"
)

var testDoc schema.Definitions

func init() {
	test.LoadTestFile("sample/event_based/event_based_gateway.bpmn", &testDoc)
}

func TestEventBasedGateway(t *testing.T) {
	testEventBasedGateway(t, func(reached map[string]int) {
		assert.Equal(t, 1, reached["task1"])
		assert.Equal(t, 1, reached["end"])
	}, event.NewSignalEvent("Sig1"))
	testEventBasedGateway(t, func(reached map[string]int) {
		assert.Equal(t, 1, reached["task2"])
		assert.Equal(t, 1, reached["end"])
	}, event.NewMessageEvent("Msg1", nil))
}

func testEventBasedGateway(t *testing.T, test func(map[string]int), events ...event.Event) {
	processElement := (*testDoc.Processes())[0]
	proc := process.New(&processElement, &testDoc)
	if instance, err := proc.Instantiate(); err == nil {
		traces := instance.Tracer.Subscribe()
		err := instance.StartAll(context.Background())
		if err != nil {
			t.Errorf("failed to run the instance: %s", err)
			return
		}

		resultChan := make(chan string)

		ctx, cancel := context.WithCancel(context.Background())

		go func(ctx context.Context) {
			for {
				select {
				case trace := <-traces:
					trace = tracing.Unwrap(trace)
					switch trace := trace.(type) {
					case ev.ActiveListeningTrace:
						if id, present := trace.Node.Id(); present {
							resultChan <- *id
						}
					case tracing.ErrorTrace:
						t.Errorf("%#v", trace)
						return
					default:
						t.Logf("%#v", trace)
					}
				case <-ctx.Done():
					return
				}
			}
		}(ctx)

		// Wait until both events are ready to listen
		assert.Regexp(t, "(signalEvent|messageEvent)", <-resultChan)
		assert.Regexp(t, "(signalEvent|messageEvent)", <-resultChan)

		cancel()

		ch := make(chan map[string]int)
		go func() {
			reached := make(map[string]int)
			for {
				trace := tracing.Unwrap(<-traces)
				switch trace := trace.(type) {
				case flow.VisitTrace:
					if id, present := trace.Node.Id(); present {
						if counter, ok := reached[*id]; ok {
							reached[*id] = counter + 1
						} else {
							reached[*id] = 1
						}
					} else {
						t.Errorf("can't find element with FlowNodeId %#v", id)
						ch <- reached
						return
					}
				case flow.CeaseFlowTrace:
					ch <- reached
					return
				case tracing.ErrorTrace:
					t.Errorf("%#v", trace)
					ch <- reached
					return
				default:
					t.Logf("%#v", trace)
				}
			}
		}()

		for _, evt := range events {
			_, err := instance.ConsumeEvent(evt)
			if err != nil {
				t.Error(err)
				return
			}
		}

		test(<-ch)

		instance.Tracer.Unsubscribe(traces)
	} else {
		t.Errorf("failed to instantiate the process: %s", err)
		return
	}
}
