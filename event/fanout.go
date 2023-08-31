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
