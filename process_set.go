/*
   Copyright 2025 The  Authors

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
	"github.com/olive-io/bpmn/v2/pkg/id"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type ProcessSet struct {
	*Options
	wg sync.WaitGroup

	executes     []*Process
	waitings     []*schema.Process
	definitions  *schema.Definitions
	complete     sync.RWMutex
	messageFlows map[string]*schema.MessageFlow
	subTracer    tracing.ITracer
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

	// flow nodes
	subTracer := tracing.NewTracer(ctx)

	executes := make([]*Process, 0)
	for _, executeProcess := range executeProcesses {
		var process *Process

		tracing.NewRelay(ctx, subTracer, options.tracer, func(trace tracing.ITrace) []tracing.ITrace {
			return []tracing.ITrace{InstanceTrace{
				InstanceId: process.id,
				Trace:      trace,
			}}
		})

		process, err = NewProcess(executeProcess, definitions, append(opts, WithTracer(subTracer))...)
		if err != nil {
			return nil, fmt.Errorf("create new process: %w", err)
		}
		executes = append(executes, process)
	}

	messageFlows := make(map[string]*schema.MessageFlow)
	for _, collaboration := range *definitions.Collaborations() {
		for _, msg := range *collaboration.MessageFlows() {
			messageFlows[string(msg.SourceRefField)] = &msg
		}
	}

	ps := &ProcessSet{
		Options:      options,
		executes:     executes,
		waitings:     waitingProcesses,
		definitions:  definitions,
		messageFlows: messageFlows,
		subTracer:    subTracer,
	}

	return ps, nil
}

func (ps *ProcessSet) Tracer() tracing.ITracer { return ps.tracer }

func (ps *ProcessSet) Locator() data.IFlowDataLocator { return ps.locator }

func (ps *ProcessSet) StartAll(ctx context.Context) error {

	return nil
}
