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

package flow_interface

import (
	"github.com/olive-io/bpmn/pkg/id"
	"github.com/olive-io/bpmn/sequence_flow"
)

// T specifies an interface for BPMN flows
type T interface {
	// Id returns flow's unique identifier
	Id() id.Id
	// SequenceFlow returns an inbound sequence flow this flow
	// is currently at.
	SequenceFlow() *sequence_flow.SequenceFlow
}
