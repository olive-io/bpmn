// Copyright 2023 Lack (xingyys@gmail.com).
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
