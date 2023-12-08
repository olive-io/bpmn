// Copyright 2023 The olive Authors
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

package timer

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tools/clock"
	"github.com/qri-io/iso8601"
	"github.com/stretchr/testify/require"
)

func TestTimeDate(t *testing.T) {
	c := clock.NewMock()

	definition := schema.DefaultTimerEventDefinition()
	iso := "2021-05-21T16:43:43+00:00"
	timestamp := schema.AnExpression{}
	err := xml.NewDecoder(bytes.NewBufferString(
		fmt.Sprintf(`<bpmn:expression>%s</bpmn:expression>`, iso),
	)).Decode(&timestamp)
	require.Nil(t, err)
	definition.SetTimeDate(&timestamp)
	timer, err := New(context.Background(), c, definition)
	require.Nil(t, err)
	requireNoMoreMessages(t, timer, false)
	time, err := iso8601.ParseTime(iso)
	require.Nil(t, err)
	c.Set(time)
	<-timer
	requireCompletion(t, timer)
}

func TestTimeDuration(t *testing.T) {
	c := clock.NewMock()

	definition := schema.DefaultTimerEventDefinition()
	iso := "PT30M"
	duration := schema.AnExpression{}
	err := xml.NewDecoder(bytes.NewBufferString(
		fmt.Sprintf(`<bpmn:expression>%s</bpmn:expression>`, iso),
	)).Decode(&duration)
	require.Nil(t, err)
	definition.SetTimeDuration(&duration)
	timer, err := New(context.Background(), c, definition)
	require.Nil(t, err)
	requireNoMoreMessages(t, timer, false)
	dur, err := iso8601.ParseDuration(iso)
	require.Nil(t, err)
	c.Add(dur.Duration)
	<-timer
	requireCompletion(t, timer)
}

func TestTimeCycle(t *testing.T) {
	c := clock.NewMock()

	definition := schema.DefaultTimerEventDefinition()
	iso := "R3/PT30M"
	cycle := schema.AnExpression{}
	err := xml.NewDecoder(bytes.NewBufferString(
		fmt.Sprintf(`<bpmn:expression>%s</bpmn:expression>`, iso),
	)).Decode(&cycle)
	require.Nil(t, err)
	definition.SetTimeCycle(&cycle)
	timer, err := New(context.Background(), c, definition)
	require.NotNil(t, timer)
	require.Nil(t, err)

	requireNoMoreMessages(t, timer, false)

	interval, err := iso8601.ParseRepeatingInterval(iso)
	require.Nil(t, err)

	for i := 0; i < interval.Repititions; i++ {
		c.Add(interval.Interval.Duration.Duration)

		<-timer

		requireNoMoreMessages(t, timer, i+1 == interval.Repititions)
	}

}

func TestTimeCycleNoRep(t *testing.T) {
	c := clock.NewMock()

	definition := schema.DefaultTimerEventDefinition()
	iso := "R0/PT30M"
	cycle := schema.AnExpression{}
	err := xml.NewDecoder(bytes.NewBufferString(
		fmt.Sprintf(`<bpmn:expression>%s</bpmn:expression>`, iso),
	)).Decode(&cycle)
	require.Nil(t, err)
	definition.SetTimeCycle(&cycle)
	timer, err := New(context.Background(), c, definition)
	require.NotNil(t, timer)
	require.Nil(t, err)

	requireNoMoreMessages(t, timer, true)

	interval, err := iso8601.ParseRepeatingInterval(iso)
	require.Nil(t, err)

	c.Add(interval.Interval.Duration.Duration)

	requireCompletion(t, timer)

}

func TestTimeCycleStartDate(t *testing.T) {
	c := clock.NewMock()

	definition := schema.DefaultTimerEventDefinition()
	date := "2021-05-21T16:43:43+00:00"
	iso := fmt.Sprintf("R3/%s/PT30M", date)
	cycle := schema.AnExpression{}
	err := xml.NewDecoder(bytes.NewBufferString(
		fmt.Sprintf(`<bpmn:expression>%s</bpmn:expression>`, iso),
	)).Decode(&cycle)
	require.Nil(t, err)
	definition.SetTimeCycle(&cycle)
	timer, err := New(context.Background(), c, definition)
	require.NotNil(t, timer)
	require.Nil(t, err)

	requireNoMoreMessages(t, timer, false)

	interval, err := iso8601.ParseRepeatingInterval(iso)
	require.Nil(t, err)

	c.Add(interval.Interval.Duration.Duration)
	requireNoMoreMessages(t, timer, false)

	c.Set(*interval.Interval.Start)
	requireNoMoreMessages(t, timer, false)

	for i := 0; i < interval.Repititions; i++ {
		c.Add(interval.Interval.Duration.Duration)

		<-timer

		requireNoMoreMessages(t, timer, i+1 == interval.Repititions)
	}

}

func TestTimeCycleIndefinitely(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c := clock.NewMock()

	definition := schema.DefaultTimerEventDefinition()
	iso := "R/PT30M"
	cycle := schema.AnExpression{}
	err := xml.NewDecoder(bytes.NewBufferString(
		fmt.Sprintf(`<bpmn:expression>%s</bpmn:expression>`, iso),
	)).Decode(&cycle)
	require.Nil(t, err)
	definition.SetTimeCycle(&cycle)
	timer, err := New(ctx, c, definition)
	require.NotNil(t, timer)
	require.Nil(t, err)

	requireNoMoreMessages(t, timer, false)

	interval, err := iso8601.ParseRepeatingInterval(iso)
	require.Nil(t, err)

	// Do some number of iterations to show that it repeats indefinitely
	for i := 0; i < 3; i++ {
		c.Add(interval.Interval.Duration.Duration)

		<-timer

		requireNoMoreMessages(t, timer, i == 2)
	}

	cancel()

}

func TestTimeCycleEndDate(t *testing.T) {
	c := clock.NewMock()

	definition := schema.DefaultTimerEventDefinition()
	date := "2021-05-21T16:43:43+00:00"
	iso := fmt.Sprintf("R/PT30M/%s", date)
	cycle := schema.AnExpression{}
	err := xml.NewDecoder(bytes.NewBufferString(
		fmt.Sprintf(`<bpmn:expression>%s</bpmn:expression>`, iso),
	)).Decode(&cycle)
	require.Nil(t, err)
	definition.SetTimeCycle(&cycle)
	timer, err := New(context.Background(), c, definition)
	require.NotNil(t, timer)
	require.Nil(t, err)

	select {
	case <-timer:
		require.FailNow(t, "shouldn't happen")
	default:
	}

	interval, err := iso8601.ParseRepeatingInterval(iso)
	require.Nil(t, err)

	// Do some number of iterations to show that it repeats indefinitely
	for i := 0; i < 3; i++ {
		c.Add(interval.Interval.Duration.Duration)

		<-timer

		requireNoMoreMessages(t, timer, i == 2)
	}

	// Shift to the end
	c.Set(*interval.Interval.End)
	// Add a duration
	c.Add(interval.Interval.Duration.Duration)

	// No more repetitions
	requireCompletion(t, timer)
}

// requireCompletion tests whether timer receives anything but channel
// closure event; if it does, it'll fail the test
func requireCompletion(t *testing.T, timer chan schema.TimerEventDefinition) {
	select {
	case _, ok := <-timer:
		// only allow channel closure
		require.False(t, ok)
	default:
	}
}

// requireNoMoreMessages tests whether timer receives anything; if it does,
// it'll fail the test, unless `last` is set to true, then it will behave exactly
// like requireCompletion.
func requireNoMoreMessages(t *testing.T, timer chan schema.TimerEventDefinition, last bool) {
	if last {
		requireCompletion(t, timer)
	} else {
		select {
		case <-timer:
			require.FailNow(t, "no more messages expected")
		default:
		}
	}
}
