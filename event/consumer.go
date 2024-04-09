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

package event

import (
	"github.com/hashicorp/go-multierror"
)

// ConsumptionResult Result of event consumption
type ConsumptionResult int

const (
	// Consumed Consumer has successfully consumed the event
	Consumed ConsumptionResult = iota
	// ConsumptionError Consumer had an unexpected error consuming the event
	//
	// The error is returned as a second value.
	ConsumptionError
	// PartiallyConsumed Consumer with multiple sub-consumers was able to
	// have at least some of the sub-consumers consume the event
	// but there were some errors (the errors are returned as
	// `multierror` list in the second value)
	PartiallyConsumed
	// DropEventConsumer Event consumer should no longer receive events
	//
	// This is a mechanism for "unsubscribing" from an event source
	DropEventConsumer
)

// IConsumer Process event consumer interface
type IConsumer interface {
	ConsumeEvent(IEvent) (ConsumptionResult, error)
}

// VoidConsumer Process event consumer that does nothing and returns EventConsumed result
type VoidConsumer struct{}

func (t VoidConsumer) ConsumeEvent(ev IEvent) (result ConsumptionResult, err error) {
	result = Consumed
	return
}

// ForwardEvent Forwards a process event to a list of consumers
//
// This function handles full and partial failure of consumers to consume the event.
// If none of them errors out, the result will be EventConsumed.
// If some do, the result will be EventPartiallyConsumed and err will be *multierror.Errors
// If all do, the result will be EventConsumptionError and err will be *multierror.Errors
func ForwardEvent(ev IEvent, eventConsumers *[]IConsumer) (result ConsumptionResult, err error) {
	var errors *multierror.Error
	for _, consumer := range *eventConsumers {
		result, consumerError := consumer.ConsumeEvent(ev)
		if result == ConsumptionError && consumerError != nil {
			errors = multierror.Append(errors, consumerError)
		}
	}
	switch {
	case errors != nil && errors.Len() > 0 && errors.Len() < len(*eventConsumers):
		// if there were errors, but not in all deliveries
		result = PartiallyConsumed
		err = errors
	case errors != nil && errors.Len() > 0:
		// if there were errors everywhere
		result = ConsumptionError
		err = errors
	default:
		// if there were no errors
		result = Consumed
	}
	return
}
