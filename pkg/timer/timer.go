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

package timer

import (
	"context"
	"time"

	"github.com/qri-io/iso8601"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/clock"
	"github.com/olive-io/bpmn/v2/pkg/errors"
)

func New(ctx context.Context, clock clock.IClock, definition schema.TimerEventDefinition) (ch chan schema.TimerEventDefinition, err error) {
	timeDate, timeDatePresent := definition.TimeDate()
	timeCycle, timeCyclePresent := definition.TimeCycle()
	timeDuration, timeDurationPresent := definition.TimeDuration()
	switch {
	case timeDatePresent && !timeCyclePresent && !timeDurationPresent:
		ch = make(chan schema.TimerEventDefinition)
		var t time.Time
		t, err = iso8601.ParseTime(*timeDate.Expression.TextPayload())
		if err != nil {
			return
		}
		go dateTimeTimer(ctx, clock, t, func() {
			ch <- definition
			close(ch)
		})
	case !timeDatePresent && timeCyclePresent && !timeDurationPresent:
		ch = make(chan schema.TimerEventDefinition)
		var repeatingInterval iso8601.RepeatingInterval
		repeatingInterval, err = iso8601.ParseRepeatingInterval(*timeCycle.Expression.TextPayload())
		if err != nil {
			return
		}
		if repeatingInterval.Interval.Start == nil {
			now := clock.Now()
			repeatingInterval.Interval.Start = &now
		}
		go recurringTimer(ctx, clock, repeatingInterval, func() {
			ch <- definition
		}, func() {
			close(ch)
		})
	case !timeDatePresent && !timeCyclePresent && timeDurationPresent:
		ch = make(chan schema.TimerEventDefinition)
		var duration iso8601.Duration
		duration, err = iso8601.ParseDuration(*timeDuration.Expression.TextPayload())
		if err != nil {
			return
		}
		go dateTimeTimer(ctx, clock, clock.Now().Add(duration.Duration), func() {
			ch <- definition
			close(ch)
		})
	default:
		err = errors.InvalidArgumentError{
			Expected: "one and only one of timeDate, timeCycle or timeDuration must be defined",
			Actual:   definition,
		}
		return
	}
	return
}

func recurringTimer(ctx context.Context, clock clock.IClock, interval iso8601.RepeatingInterval, f func(), final func()) {
	if interval.Interval.Start == nil {
		panic("shouldn't happen, has to be always set, explicitly or by timer.New")
	}
	ch := make(chan struct{})
	go dateTimeTimer(ctx, clock, *interval.Interval.Start, func() {
		ch <- struct{}{}
	})
	select {
	case <-ctx.Done():
		return
	case <-ch:
		break
	}

	repetitions := interval.Repititions

	var endTimer <-chan time.Time
	var timer <-chan time.Time

	t := *interval.Interval.Start

	defer final()

	for {
		if repetitions == 0 {
			return
		}

		timer = clock.Until(t.Add(interval.Interval.Duration.Duration))

		if interval.Interval.End != nil {
			endTimer = clock.Until(*interval.Interval.End)
		}

		select {
		case <-endTimer:
			return
		case <-ctx.Done():
			return
		case t = <-timer:
			if interval.Interval.End == nil || interval.Interval.End.After(clock.Now()) {
				f()
			}
		}

		if repetitions > 0 {
			repetitions--
		}
	}
}

func dateTimeTimer(ctx context.Context, clock clock.IClock, t time.Time, f func()) {
	for {
		timer := clock.Until(t)
		select {
		case <-ctx.Done():
			return
		case <-timer:
			f()
			return
		}
	}
}
