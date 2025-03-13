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

package event

import "sync"

// FanOut is a straightforward Consumer + Source, forwards all consumed
// messages to all subscribers registered at it.
type FanOut struct {
	eventConsumersLock sync.RWMutex
	eventConsumers     []IConsumer
}

func NewFanOut() *FanOut {
	return &FanOut{}
}

func (f *FanOut) ConsumeEvent(ev IEvent) (result ConsumptionResult, err error) {
	f.eventConsumersLock.RLock()
	defer f.eventConsumersLock.RUnlock()
	result, err = ForwardEvent(ev, &f.eventConsumers)
	return
}

func (f *FanOut) RegisterEventConsumer(ev IConsumer) (err error) {
	f.eventConsumersLock.Lock()
	defer f.eventConsumersLock.Unlock()
	f.eventConsumers = append(f.eventConsumers, ev)
	return
}
