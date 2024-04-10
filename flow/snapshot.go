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

package flow

import (
	"github.com/olive-io/bpmn/pkg/id"
	"github.com/olive-io/bpmn/sequence_flow"
)

type Snapshot struct {
	flowId       id.Id
	sequenceFlow *sequence_flow.SequenceFlow
}

func (s *Snapshot) Id() id.Id {
	return s.flowId
}

func (s *Snapshot) SequenceFlow() *sequence_flow.SequenceFlow {
	return s.sequenceFlow
}
