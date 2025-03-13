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

package schema

// InstantiatingFlowNodes returns a list of flow nodes that
// can instantiate the process.
func (t *Process) InstantiatingFlowNodes() (result []FlowNodeInterface) {
	result = make([]FlowNodeInterface, 0)

	for i := range *t.StartEvents() {
		startEvent := &(*t.StartEvents())[i]
		// Start event that observes some events
		if len(startEvent.EventDefinitions()) > 0 {
			result = append(result, startEvent)
		}
	}

	for i := range *t.EventBasedGateways() {
		gateway := &(*t.EventBasedGateways())[i]
		// Event-based gateways with `instantiate` set to true
		// and no incoming sequence flows
		if gateway.Instantiate() && len(*gateway.Incomings()) == 0 {
			result = append(result, gateway)
		}
	}

	for i := range *t.ReceiveTasks() {
		task := &(*t.ReceiveTasks())[i]
		// Event-based gateways with `instantiate` set to true
		// and no incoming sequence flows
		if task.Instantiate() && len(*task.Incomings()) == 0 {
			result = append(result, task)
		}
	}

	return
}
