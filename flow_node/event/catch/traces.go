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

package catch

import (
	"github.com/olive-io/bpmn/schema"

	"github.com/olive-io/bpmn/event"
)

type ActiveListeningTrace struct {
	Node *schema.CatchEvent
}

func (t ActiveListeningTrace) TraceInterface() {}

// EventObservedTrace signals the fact that a particular event
// has been in fact observed by the node
type EventObservedTrace struct {
	Node  *schema.CatchEvent
	Event event.IEvent
}

func (t EventObservedTrace) TraceInterface() {}
