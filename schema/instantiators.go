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
