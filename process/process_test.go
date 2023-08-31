package process

import (
	"context"
	"testing"

	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/test"
	"github.com/olive-io/bpmn/tracing"
	"github.com/stretchr/testify/assert"
)

var defaultDefinitions = schema.DefaultDefinitions()

var sampleDoc schema.Definitions

func init() {
	test.LoadTestFile("sample/process/sample.bpmn", &sampleDoc)
}

func TestExplicitInstantiation(t *testing.T) {
	if proc, found := sampleDoc.FindBy(schema.ExactId("sample")); found {
		process := New(proc.(*schema.Process), &defaultDefinitions)
		inst, err := process.Instantiate()
		assert.Nil(t, err)
		assert.NotNil(t, inst)
	} else {
		t.Fatalf("Can't find process `sample`")
	}
}

func TestCancellation(t *testing.T) {
	if proc, found := sampleDoc.FindBy(schema.ExactId("sample")); found {
		ctx, cancel := context.WithCancel(context.Background())

		process := New(proc.(*schema.Process), &defaultDefinitions, WithContext(ctx))

		tracer := tracing.NewTracer(ctx)
		traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 128))

		inst, err := process.Instantiate(instance.WithContext(ctx), instance.WithTracer(tracer))
		assert.Nil(t, err)
		assert.NotNil(t, inst)

		cancel()

		cancelledFlowNodes := make([]schema.FlowNodeInterface, 0)

		for trace := range traces {
			trace = tracing.Unwrap(trace)
			switch trace := trace.(type) {
			case flow_node.CancellationTrace:
				cancelledFlowNodes = append(cancelledFlowNodes, trace.Node)
			default:
			}
		}

		assert.NotEmpty(t, cancelledFlowNodes)
	} else {
		t.Fatalf("Can't find process `sample`")
	}
}
