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

package flow

import (
	"context"
	"fmt"
	"sync"

	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/errors"
	"github.com/olive-io/bpmn/expression"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/sequence_flow"
	"github.com/olive-io/bpmn/tools/id"
	"github.com/olive-io/bpmn/tracing"
)

// Flow Represents a flow
type Flow struct {
	id                id.Id
	definitions       *schema.Definitions
	current           flow_node.IFlowNode
	tracer            tracing.ITracer
	flowNodeMapping   *flow_node.FlowNodeMapping
	flowWaitGroup     *sync.WaitGroup
	idGenerator       id.IGenerator
	actionTransformer flow_node.ActionTransformer
	terminate         flow_node.Terminate
	sequenceFlowId    *string
	locator           data.IFlowDataLocator
	retry             *Retry
}

func (flow *Flow) SequenceFlow() *sequence_flow.SequenceFlow {
	if flow.sequenceFlowId == nil {
		return nil
	}
	seqFlow, present := flow.definitions.FindBy(schema.ExactId(*flow.sequenceFlowId).
		And(schema.ElementType((*schema.SequenceFlow)(nil))))
	if present {
		return sequence_flow.New(seqFlow.(*schema.SequenceFlow), flow.definitions)
	}

	return nil
}

func (flow *Flow) Id() id.Id {
	return flow.id
}

func (flow *Flow) SetTerminate(terminate flow_node.Terminate) {
	flow.terminate = terminate
}

// New creates a new flow from a flow node
//
// The flow does nothing until it is explicitly started.
func New(definitions *schema.Definitions,
	current flow_node.IFlowNode, tracer tracing.ITracer,
	flowNodeMapping *flow_node.FlowNodeMapping, flowWaitGroup *sync.WaitGroup,
	idGenerator id.IGenerator, actionTransformer flow_node.ActionTransformer,
	locator data.IFlowDataLocator,
) *Flow {
	return &Flow{
		id:                idGenerator.New(),
		definitions:       definitions,
		current:           current,
		tracer:            tracer,
		flowNodeMapping:   flowNodeMapping,
		flowWaitGroup:     flowWaitGroup,
		idGenerator:       idGenerator,
		actionTransformer: actionTransformer,
		locator:           locator,
	}
}

func (flow *Flow) executeSequenceFlow(ctx context.Context, sequenceFlow *sequence_flow.SequenceFlow, unconditional bool) (result bool, err error) {
	if unconditional {
		result = true
		return
	}
	if expr, present := sequenceFlow.SequenceFlow.ConditionExpression(); present {
		switch e := expr.Expression.(type) {
		case *schema.FormalExpression:
			var lang string

			if language, present := e.Language(); present {
				lang = *language
			} else {
				lang = *flow.definitions.ExpressionLanguage()
			}

			dataSets := map[string]any{}
			engine := expression.GetEngine(ctx, lang)
			if locator, found := flow.locator.FindIItemAwareLocator(data.LocatorObject); found {
				engine.SetItemAwareLocator(data.LocatorObject, locator)
			}
			if locator, found := flow.locator.FindIItemAwareLocator(data.LocatorHeader); found {
				engine.SetItemAwareLocator(data.LocatorHeader, locator)
			}

			for key, item := range flow.locator.CloneVariables() {
				dataSets[key] = item
			}

			source := *e.TextPayload()
			var compiled expression.ICompiledExpression
			compiled, err = engine.CompileExpression(source)
			if err != nil {
				result = false
				flow.tracer.Trace(tracing.ErrorTrace{Error: err})
				return
			}
			var abstractResult expression.IResult
			abstractResult, err = engine.EvaluateExpression(compiled, dataSets)
			if err != nil {
				result = false
				flow.tracer.Trace(tracing.ErrorTrace{Error: err})
				return
			}
			switch actualResult := abstractResult.(type) {
			case bool:
				result = actualResult
			default:
				err = errors.InvalidArgumentError{
					Expected: fmt.Sprintf("boolean result in conditionExpression (%s)", source),
					Actual:   actualResult,
				}
				result = false
				flow.tracer.Trace(tracing.ErrorTrace{Error: err})
				return
			}
		case *schema.Expression:
			// informal expression, can't execute
			result = true
		}
	} else {
		result = true
	}

	return
}

func (flow *Flow) handleSequenceFlow(ctx context.Context, sequenceFlow *sequence_flow.SequenceFlow, unconditional bool,
	actionTransformer flow_node.ActionTransformer, terminate flow_node.Terminate) (flowed bool) {
	ok, err := flow.executeSequenceFlow(ctx, sequenceFlow, unconditional)
	if err != nil {
		flow.tracer.Trace(tracing.ErrorTrace{Error: err})
		return
	}
	if !ok {
		return
	}

	target, err := sequenceFlow.Target()
	if err != nil {
		flow.tracer.Trace(tracing.ErrorTrace{Error: err})
		return
	}

	if flowNode, found := flow.flowNodeMapping.ResolveElementToFlowNode(target); found {
		if idPtr, present := sequenceFlow.Id(); present {
			flow.sequenceFlowId = idPtr
		} else {
			flow.tracer.Trace(tracing.ErrorTrace{
				Error: errors.NotFoundError{Expected: fmt.Sprintf("id for sequence flow %#v", sequenceFlow)},
			})
		}
		flow.current = flowNode
		flow.terminate = terminate
		flow.tracer.Trace(VisitTrace{Node: flow.current.Element()})
		flow.actionTransformer = actionTransformer
		flowed = true
		return
	}

	flow.tracer.Trace(tracing.ErrorTrace{
		Error: errors.NotFoundError{Expected: fmt.Sprintf("flow node for element %#v", target)},
	})
	return
}

// handleAdditionalSequenceFlow returns a new flowId (if it will flow), flow start function and a flag
// that indicates whether it'll flow.
//
// The reason behind the way this function works is that we don't want these additional sequence flows
// to flow until after we logged the fact that they'll flow (using FlowTrace). This makes it consistent
// with the behaviour of handleSequenceFlow since it doesn't really flow anywhere but simply sets the current
// flow to continue flowing.
//
// Having this FlowTrace consistency is extremely important because otherwise there will be cases when
// the FlowTrace will come through *after* traces created by newly started flows, completely messing up
// the order of traces which is expected to be linear.
func (flow *Flow) handleAdditionalSequenceFlow(ctx context.Context, sequenceFlow *sequence_flow.SequenceFlow,
	unconditional bool, actionTransformer flow_node.ActionTransformer,
	terminate flow_node.Terminate) (flowId id.Id, f func(), flowed bool) {
	ok, err := flow.executeSequenceFlow(ctx, sequenceFlow, unconditional)
	if err != nil {
		flow.tracer.Trace(tracing.ErrorTrace{Error: err})
		return
	}
	if !ok {
		return
	}
	target, err := sequenceFlow.Target()
	if err != nil {
		flow.tracer.Trace(tracing.ErrorTrace{Error: err})
		return
	}
	if flowNode, found := flow.flowNodeMapping.ResolveElementToFlowNode(target); found {
		flowId = flow.idGenerator.New()
		f = func() {
			flowable := New(flow.definitions, flowNode, flow.tracer, flow.flowNodeMapping,
				flow.flowWaitGroup, flow.idGenerator, actionTransformer, flow.locator)
			flowable.id = flowId // important: override id with pre-generated one
			if idPtr, present := sequenceFlow.Id(); present {
				flowable.sequenceFlowId = idPtr
			} else {
				flow.tracer.Trace(tracing.ErrorTrace{
					Error: errors.NotFoundError{Expected: fmt.Sprintf("id for sequence flow %#v", sequenceFlow)},
				})
			}
			flowable.terminate = terminate
			flowable.Start(ctx)
		}
		flowed = true
	} else {
		flow.tracer.Trace(tracing.ErrorTrace{
			Error: errors.NotFoundError{Expected: fmt.Sprintf("flow node for element %#v", target)},
		})
	}
	return
}

func (flow *Flow) termination() chan bool {
	if flow.terminate == nil {
		return nil
	}
	return flow.terminate(flow.sequenceFlowId)
}

// Start starts the flow
func (flow *Flow) Start(ctx context.Context) {
	flow.flowWaitGroup.Add(1)
	sender := flow.tracer.RegisterSender()
	go func() {
		defer sender.Done()
		flow.tracer.Trace(NewFlowTrace{FlowId: flow.id})
		defer flow.flowWaitGroup.Done()
		flow.tracer.Trace(VisitTrace{Node: flow.current.Element()})
		for {
		await:
			select {
			case <-ctx.Done():
				flow.tracer.Trace(CancellationTrace{
					FlowId: flow.Id(),
				})
				return
			case terminate := <-flow.termination():
				if terminate {
					flow.tracer.Trace(TerminationTrace{
						FlowId: flow.Id(),
						Source: flow.current.Element(),
					})
					return
				} else {
					goto await
				}
			case action := <-flow.current.NextAction(flow):
				if flow.actionTransformer != nil {
					action = flow.actionTransformer(flow.sequenceFlowId, action)
				}
				switch a := action.(type) {
				case flow_node.ProbeAction:
					results := make([]int, 0)
					for i, seqFlow := range a.SequenceFlows {
						if result, err := flow.executeSequenceFlow(ctx, seqFlow, false); err == nil {
							if result {
								results = append(results, i)
							}
						} else {
							flow.tracer.Trace(tracing.ErrorTrace{Error: err})
						}
					}
					a.ProbeReport(results)
				case flow_node.FlowAction:
					if res := a.Response; res != nil {
						if res.Err != nil {
							source := flow.current.Element()
							fid, _ := source.Id()
							flow.tracer.Trace(tracing.ErrorTrace{Error: errors.TaskExecError{
								Id:     *fid,
								Reason: res.Err.Error(),
							}})

							if res.Handler == nil {
								return
							}

							select {
							case handler := <-res.Handler:
								switch handler.Mode {
								case flow_node.HandleRetry:
									if flow.retry == nil {
										retry := &Retry{}
										if extension, present := source.ExtensionElements(); present {
											if taskDefinition := extension.TaskDefinitionField; taskDefinition != nil {
												retry.Reset(taskDefinition.Retries)
											}
										}
										flow.retry = retry
									}

									flow.retry.Reset(handler.Retries)
									if flow.retry.IsContinue() {
										flow.retry.Step()
										goto await
									}
									return
								case flow_node.HandleSkip:
								case flow_node.HandleExit:
									return
								}
							case <-ctx.Done():
								return
							}
						}

						if len(res.DataObjects) > 0 {
							locator, found := flow.locator.FindIItemAwareLocator(data.LocatorObject)
							if found {
								for name, do := range res.DataObjects {
									aware, found := locator.FindItemAwareByName(name)
									if !found {
										aware = data.NewContainer(nil)
									}
									aware.Put(do)
								}
							}
						}
						if len(res.Variables) > 0 {
							for key, value := range res.Variables {
								flow.locator.SetVariable(key, value)
							}
						}
					}

					sequenceFlows := a.SequenceFlows
					if len(a.SequenceFlows) > 0 {
						unconditional := make([]bool, len(a.SequenceFlows))
						for _, index := range a.UnconditionalFlows {
							unconditional[index] = true
						}
						source := flow.current.Element()

						current := sequenceFlows[0]
						effectiveFlows := make([]Snapshot, 0)

						flowed := flow.handleSequenceFlow(ctx, current, unconditional[0], a.ActionTransformer, a.Terminate)

						if flowed {
							effectiveFlows = append(effectiveFlows, Snapshot{sequenceFlow: current, flowId: flow.Id()})
						}

						rest := sequenceFlows[1:]
						flowFuncs := make([]func(), 0)
						for i, sequenceFlow := range rest {
							flowId, flowFunc, flowed := flow.handleAdditionalSequenceFlow(ctx, sequenceFlow, unconditional[i+1],
								a.ActionTransformer, a.Terminate)
							if flowed {
								effectiveFlows = append(effectiveFlows, Snapshot{sequenceFlow: sequenceFlow, flowId: flowId})
								flowFuncs = append(flowFuncs, flowFunc)
							}
						}

						if len(effectiveFlows) > 0 {
							flow.tracer.Trace(Trace{
								Source: source,
								Flows:  effectiveFlows,
							})
							for _, flowFunc := range flowFuncs {
								flowFunc()
							}
						} else {
							// no flows to continue with, abort
							flow.tracer.Trace(TerminationTrace{
								FlowId: flow.Id(),
								Source: source,
							})
							return
						}

					} else {
						// nowhere to flow, abort
						return
					}
				case flow_node.CompleteAction:
					flow.tracer.Trace(CompletionTrace{
						Node: flow.current.Element(),
					})
					flow.tracer.Trace(TerminationTrace{
						FlowId: flow.Id(),
						Source: flow.current.Element(),
					})
					return
				case flow_node.NoAction:
					flow.tracer.Trace(TerminationTrace{
						FlowId: flow.Id(),
						Source: flow.current.Element(),
					})
					return
				default:
				}
			}
		}

	}()
}
