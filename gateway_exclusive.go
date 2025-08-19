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

package bpmn

import (
	"context"
	"fmt"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/errors"
	"github.com/olive-io/bpmn/v2/pkg/id"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type ExclusiveNoEffectiveSequenceFlows struct {
	*schema.ExclusiveGateway
}

func (e ExclusiveNoEffectiveSequenceFlows) Error() string {
	ownId := "<unnamed>"
	if ownIdPtr, present := e.ExclusiveGateway.Id(); present {
		ownId = *ownIdPtr
	}
	return fmt.Sprintf("No effective sequence flows found in exclusive gateway `%v`", ownId)
}

type ExclusiveGateway struct {
	*Wiring
	element                 *schema.ExclusiveGateway
	mch                     chan imessage
	defaultSequenceFlow     *SequenceFlow
	nonDefaultSequenceFlows []*SequenceFlow
	probing                 map[id.Id]*chan IAction
}

func NewExclusiveGateway(wiring *Wiring, exclusiveGateway *schema.ExclusiveGateway) (gw *ExclusiveGateway, err error) {
	var defaultSequenceFlow *SequenceFlow

	if seqFlow, present := exclusiveGateway.Default(); present {
		if gw, found := wiring.Process.FindBy(schema.ExactId(*seqFlow).
			And(schema.ElementType((*schema.SequenceFlow)(nil)))); found {
			defaultSequenceFlow = new(SequenceFlow)
			*defaultSequenceFlow = MakeSequenceFlow(
				gw.(*schema.SequenceFlow),
				wiring.Process,
			)
		} else {
			err = errors.NotFoundError{
				Expected: fmt.Sprintf("default sequence flow with ID %s", *seqFlow),
			}
			return nil, err
		}
	}

	nonDefaultSequenceFlows := AllSequenceFlows(&wiring.Outgoing,
		func(sequenceFlow *SequenceFlow) bool {
			if defaultSequenceFlow == nil {
				return false
			}
			return *sequenceFlow == *defaultSequenceFlow
		},
	)

	gw = &ExclusiveGateway{
		Wiring:                  wiring,
		element:                 exclusiveGateway,
		mch:                     make(chan imessage, len(wiring.Incoming)*2+1),
		nonDefaultSequenceFlows: nonDefaultSequenceFlows,
		defaultSequenceFlow:     defaultSequenceFlow,
		probing:                 make(map[id.Id]*chan IAction),
	}

	return
}

func (gw *ExclusiveGateway) run(ctx context.Context, sender tracing.ISenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-gw.mch:
			switch m := msg.(type) {
			case gatewayProbingReport:
				if response, ok := gw.probing[m.flowId]; ok {
					if response == nil {
						// Reschedule, there's no next action yet
						go func() {
							gw.mch <- m
						}()
						continue
					}
					delete(gw.probing, m.flowId)
					sfs := make([]*SequenceFlow, 0)
					for _, i := range m.result {
						sfs = append(sfs, gw.nonDefaultSequenceFlows[i])
						break
					}
					switch len(sfs) {
					case 0:
						// no successful non-default sequence flows
						if gw.defaultSequenceFlow == nil {
							// exception (Table 13.2)
							gw.Wiring.Tracer.Send(ErrorTrace{
								Error: ExclusiveNoEffectiveSequenceFlows{
									ExclusiveGateway: gw.element,
								},
							})
						} else {
							// default
							*response <- FlowAction{
								SequenceFlows:      []*SequenceFlow{gw.defaultSequenceFlow},
								UnconditionalFlows: []int{0},
							}
						}
					case 1:
						*response <- FlowAction{
							SequenceFlows:      sfs,
							UnconditionalFlows: []int{0},
						}
					default:
						gw.Wiring.Tracer.Send(ErrorTrace{
							Error: errors.InvalidArgumentError{
								Expected: fmt.Sprintf("maximum 1 outgoing exclusive gateway (%s) flow",
									gw.Wiring.FlowNodeId),
								Actual: len(sfs),
							},
						})
					}
				} else {
					gw.Wiring.Tracer.Send(ErrorTrace{
						Error: errors.InvalidStateError{
							Expected: fmt.Sprintf("probing[%s] is to be present (exclusive gateway %s)",
								m.flowId.String(), gw.Wiring.FlowNodeId),
						},
					})
				}
			case nextActionMessage:
				if _, ok := gw.probing[m.flow.Id()]; ok {
					gw.probing[m.flow.Id()] = &m.response
					// and now we wait until the probe has returned
				} else {
					gw.probing[m.flow.Id()] = nil
					m.response <- ProbeAction{
						SequenceFlows: gw.nonDefaultSequenceFlows,
						ProbeReport: func(indices []int) {
							gw.mch <- gatewayProbingReport{
								result: indices,
								flowId: m.flow.Id(),
							}
						},
					}
				}
			default:
			}
		case <-ctx.Done():
			gw.Tracer.Send(CancellationFlowNodeTrace{Node: gw.element})
			return
		}
	}
}

func (gw *ExclusiveGateway) NextAction(ctx context.Context, flow Flow) chan IAction {
	sender := gw.Tracer.RegisterSender()
	go gw.run(ctx, sender)

	response := make(chan IAction)
	gw.mch <- nextActionMessage{response: response, flow: flow}
	return response
}

func (gw *ExclusiveGateway) Element() schema.FlowNodeInterface {
	return gw.element
}
