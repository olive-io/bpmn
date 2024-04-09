/*
   Copyright 2023 The bpmn Authors

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

package flow

import (
	"github.com/olive-io/bpmn/schema"

	"github.com/olive-io/bpmn/pkg/id"
)

type NewFlowTrace struct {
	FlowId id.Id
}

func (t NewFlowTrace) TraceInterface() {}

type Trace struct {
	Source schema.FlowNodeInterface
	Flows  []Snapshot
}

func (t Trace) TraceInterface() {}

type TerminationTrace struct {
	FlowId id.Id
	Source schema.FlowNodeInterface
}

func (t TerminationTrace) TraceInterface() {}

type CancellationTrace struct {
	FlowId id.Id
}

func (t CancellationTrace) TraceInterface() {}

type CompletionTrace struct {
	Node schema.FlowNodeInterface
}

func (t CompletionTrace) TraceInterface() {}

type CeaseFlowTrace struct{}

func (t CeaseFlowTrace) TraceInterface() {}

type VisitTrace struct {
	Node schema.FlowNodeInterface
}

func (t VisitTrace) TraceInterface() {}

type LeaveTrace struct {
	Node schema.FlowNodeInterface
}

func (t LeaveTrace) TraceInterface() {}
