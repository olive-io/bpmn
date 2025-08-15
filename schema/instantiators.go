/*
Copyright 2023 The bpmn Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
