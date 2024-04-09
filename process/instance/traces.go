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

package instance

import (
	"github.com/olive-io/bpmn/pkg/id"
	"github.com/olive-io/bpmn/tracing"
)

// InstantiationTrace denotes instantiation of a given process
type InstantiationTrace struct {
	InstanceId id.Id
}

func (i InstantiationTrace) TraceInterface() {}

// Trace wraps any trace with process instance id
type Trace struct {
	InstanceId id.Id
	Trace      tracing.ITrace
}

func (t Trace) Unwrap() tracing.ITrace {
	return t.Trace
}

func (t Trace) TraceInterface() {}
