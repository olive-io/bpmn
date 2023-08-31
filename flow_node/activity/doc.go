package activity

import "github.com/olive-io/bpmn/flow_node"

// Activity is a generic interface to flow nodes that are activities
type Activity interface {
	flow_node.FlowNodeInterface
	// Cancel initiates a cancellation of activity and returns a channel
	// that will signal a boolean (`true` if cancellation was successful,
	// `false` otherwise)
	Cancel() <-chan bool
}
