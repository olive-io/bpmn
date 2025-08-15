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
	"sync"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/data"
	"github.com/olive-io/bpmn/v2/pkg/errors"
	"github.com/olive-io/bpmn/v2/pkg/expression"
	"github.com/olive-io/bpmn/v2/pkg/id"
	"github.com/olive-io/bpmn/v2/pkg/tracing"

	_ "github.com/olive-io/bpmn/v2/pkg/expression/expr"
	_ "github.com/olive-io/bpmn/v2/pkg/expression/xpath"
)

// Flow specifies an interface for BPMN flows
type Flow interface {
	// Id returns flow's unique identifier
	Id() id.Id
	// SequenceFlow returns an inbound sequence flow this flow
	// is currently at.
	SequenceFlow() *SequenceFlow
}

type Snapshot struct {
	flowId       id.Id
	sequenceFlow *SequenceFlow
}

func (s *Snapshot) Id() id.Id { return s.flowId }

func (s *Snapshot) SequenceFlow() *SequenceFlow { return s.sequenceFlow }

// flow Represents a flow
type flow struct {
	id                id.Id
	definitions       *schema.Definitions
	current           IFlowNode
	tracer            tracing.ITracer
	flowNodeMapping   *FlowNodeMapping
	flowWaitGroup     *sync.WaitGroup
	idGenerator       id.IGenerator
	actionTransformer ActionTransformer
	terminate         Terminate
	sequenceFlowId    *string
	locator           data.IFlowDataLocator
	retry             *Retry
}

func (f *flow) SequenceFlow() *SequenceFlow {
	if f.sequenceFlowId == nil {
		return nil
	}
	seqFlow, present := f.definitions.FindBy(schema.ExactId(*f.sequenceFlowId).
		And(schema.ElementType((*schema.SequenceFlow)(nil))))
	if present {
		return NewSequenceFlow(seqFlow.(*schema.SequenceFlow), f.definitions)
	}

	return nil
}

func (f *flow) Id() id.Id { return f.id }

func (f *flow) SetTerminate(terminate Terminate) { f.terminate = terminate }

// newFlow creates a new flow from a flow node
//
// The flow does nothing until it is explicitly started.
func newFlow(definitions *schema.Definitions,
	current IFlowNode, tracer tracing.ITracer,
	flowNodeMapping *FlowNodeMapping, flowWaitGroup *sync.WaitGroup,
	idGenerator id.IGenerator, actionTransformer ActionTransformer,
	locator data.IFlowDataLocator,
) *flow {
	return &flow{
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

func (f *flow) executeSequenceFlow(ctx context.Context, sequenceFlow *SequenceFlow, unconditional bool) (result bool, err error) {
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
				lang = *f.definitions.ExpressionLanguage()
			}

			properties := map[string]any{}
			engine := expression.GetEngine(ctx, lang)
			if locator, found := f.locator.FindIItemAwareLocator(data.LocatorObject); found {
				engine.SetItemAwareLocator(data.LocatorObject, locator)
			}
			if locator, found := f.locator.FindIItemAwareLocator(data.LocatorHeader); found {
				engine.SetItemAwareLocator(data.LocatorHeader, locator)
			}
			if locator, found := f.locator.FindIItemAwareLocator(data.LocatorProperty); found {
				engine.SetItemAwareLocator(data.LocatorProperty, locator)
			}

			for key, item := range f.locator.CloneVariables() {
				properties[key] = item
			}

			source := *e.TextPayload()
			var compiled expression.ICompiledExpression
			compiled, err = engine.CompileExpression(source)
			if err != nil {
				result = false
				f.tracer.Send(ErrorTrace{Error: err})
				return
			}
			var abstractResult expression.IResult
			abstractResult, err = engine.EvaluateExpression(compiled, properties)
			if err != nil {
				result = false
				f.tracer.Send(ErrorTrace{Error: err})
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
				f.tracer.Send(ErrorTrace{Error: err})
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

func (f *flow) handleSequenceFlow(ctx context.Context, sequenceFlow *SequenceFlow, unconditional bool,
	actionTransformer ActionTransformer, terminate Terminate) (flowed bool) {
	ok, err := f.executeSequenceFlow(ctx, sequenceFlow, unconditional)
	if err != nil {
		f.tracer.Send(ErrorTrace{Error: err})
		return
	}
	if !ok {
		return
	}

	f.tracer.Send(LeaveTrace{Node: f.current.Element()})
	target, err := sequenceFlow.Target()
	if err != nil {
		f.tracer.Send(ErrorTrace{Error: err})
		return
	}

	if flowNode, found := f.flowNodeMapping.ResolveElementToFlowNode(target); found {
		if idPtr, present := sequenceFlow.Id(); present {
			f.sequenceFlowId = idPtr
		} else {
			f.tracer.Send(ErrorTrace{
				Error: errors.NotFoundError{Expected: fmt.Sprintf("id for sequence flow %#v", sequenceFlow)},
			})
		}
		f.current = flowNode
		f.terminate = terminate
		node := f.current.Element()
		f.tracer.Send(VisitTrace{Node: node})
		f.actionTransformer = actionTransformer
		flowed = true
		return
	}

	f.tracer.Send(ErrorTrace{
		Error: errors.NotFoundError{Expected: fmt.Sprintf("flow node for element %#v", target)},
	})
	return
}

// handleAdditionalSequenceFlow returns a new flowId (if it will flow), flow start function and a flag
// that indicates whether it'll flow.
//
// The reason behind the way this function works is that we don't want these additional sequence flows
// to flow until after we logged the fact that they'll flow (using FlowTrace). This makes it consistent
// with the behavior of handleSequenceFlow since it doesn't really flow anywhere but simply sets the current
// flow to continue flowing.
//
// Having this FlowTrace consistency is extremely important because otherwise there will be cases when
// the FlowTrace will come through *after* traces created by newly started flows, completely messing up
// the order of traces which is expected to be linear.
func (f *flow) handleAdditionalSequenceFlow(ctx context.Context, sequenceFlow *SequenceFlow,
	unconditional bool, actionTransformer ActionTransformer,
	terminate Terminate) (flowId id.Id, handle func(), flowed bool) {
	ok, err := f.executeSequenceFlow(ctx, sequenceFlow, unconditional)
	if err != nil {
		f.tracer.Send(ErrorTrace{Error: err})
		return
	}
	if !ok {
		return
	}
	target, err := sequenceFlow.Target()
	if err != nil {
		f.tracer.Send(ErrorTrace{Error: err})
		return
	}
	if flowNode, found := f.flowNodeMapping.ResolveElementToFlowNode(target); found {
		flowId = f.idGenerator.New()
		handle = func() {
			flowable := newFlow(f.definitions, flowNode, f.tracer, f.flowNodeMapping,
				f.flowWaitGroup, f.idGenerator, actionTransformer, f.locator)
			flowable.id = flowId // important: override id with pre-generated one
			if idPtr, present := sequenceFlow.Id(); present {
				flowable.sequenceFlowId = idPtr
			} else {
				f.tracer.Send(ErrorTrace{
					Error: errors.NotFoundError{Expected: fmt.Sprintf("id for sequence flow %#v", sequenceFlow)},
				})
			}
			flowable.terminate = terminate
			flowable.Start(ctx)
		}
		flowed = true
	} else {
		f.tracer.Send(ErrorTrace{
			Error: errors.NotFoundError{Expected: fmt.Sprintf("flow node for element %#v", target)},
		})
	}
	return
}

func (f *flow) termination() chan bool {
	if f.terminate == nil {
		return nil
	}
	return f.terminate(f.sequenceFlowId)
}

// Start starts the flow
func (f *flow) Start(ctx context.Context) {
	f.flowWaitGroup.Add(1)
	sender := f.tracer.RegisterSender()
	go func() {
		defer sender.Done()
		f.tracer.Send(NewFlowTrace{FlowId: f.id})
		defer f.flowWaitGroup.Done()
		f.tracer.Send(VisitTrace{Node: f.current.Element()})
		for {
		await:
			select {
			case <-ctx.Done():
				f.tracer.Send(CancellationFlowTrace{
					FlowId: f.Id(),
				})
				return
			case terminate := <-f.termination():
				if terminate {
					f.tracer.Send(TerminationTrace{
						FlowId: f.Id(),
						Source: f.current.Element(),
					})
					return
				} else {
					goto await
				}
			case action := <-f.current.NextAction(f):
				if f.actionTransformer != nil {
					action = f.actionTransformer(f.sequenceFlowId, action)
				}
				switch a := action.(type) {
				case ProbeAction:
					results := make([]int, 0)
					for i, seqFlow := range a.SequenceFlows {
						if result, err := f.executeSequenceFlow(ctx, seqFlow, false); err == nil {
							if result {
								results = append(results, i)
							}
						} else {
							f.tracer.Send(ErrorTrace{Error: err})
						}
					}
					a.ProbeReport(results)
				case FlowAction:
					if res := a.Response; res != nil {
						if res.Err != nil {
							source := f.current.Element()
							fid, _ := source.Id()
							f.tracer.Send(ErrorTrace{Error: errors.TaskExecError{
								Id:     *fid,
								Reason: res.Err.Error(),
							}})

							if res.Handler != nil {
								select {
								case handler := <-res.Handler:
									switch handler.Mode {
									case RetryMode:
										if f.retry == nil {
											retry := &Retry{}
											if extension, present := source.ExtensionElements(); present {
												if taskDefinition := extension.TaskDefinitionField; taskDefinition != nil {
													retry.Reset(taskDefinition.Retries)
												}
											}
											f.retry = retry
										}

										f.retry.Reset(handler.Retries)
										if f.retry.IsContinue() {
											f.retry.Step()
											goto await
										}
										return
									case SkipMode:
									case ExitMode:
										return
									}
								case <-ctx.Done():
									return
								}
							}
						}

						if len(res.DataObjects) > 0 {
							locator, found := f.locator.FindIItemAwareLocator(data.LocatorObject)
							if !found {
								locator = data.NewDataObjectContainer()
								f.locator.PutIItemAwareLocator(data.LocatorObject, locator)
							}
							for name, do := range res.DataObjects {
								aware, found1 := locator.FindItemAwareByName(name)
								if !found1 {
									aware = data.NewContainer(nil)
									locator.PutItemAwareById(name, aware)
								}
								aware.Put(do)
							}
						}
						if len(res.Variables) > 0 {
							for key, value := range res.Variables {
								f.locator.SetVariable(key, value)
							}
						}
					}

					sequences := a.SequenceFlows
					if len(sequences) > 0 {
						unconditional := make([]bool, len(sequences))
						for _, index := range a.UnconditionalFlows {
							unconditional[index] = true
						}
						source := f.current.Element()

						current := sequences[0]
						effectiveFlows := make([]Snapshot, 0)

						flowed := f.handleSequenceFlow(ctx, current, unconditional[0], a.ActionTransformer, a.Terminate)

						if flowed {
							effectiveFlows = append(effectiveFlows, Snapshot{sequenceFlow: current, flowId: f.Id()})
						}

						rest := sequences[1:]
						flowFuncs := make([]func(), 0)
						for i, sequenceFlow := range rest {
							flowId, flowFunc, flowed := f.handleAdditionalSequenceFlow(ctx, sequenceFlow, unconditional[i+1],
								a.ActionTransformer, a.Terminate)
							if flowed {
								effectiveFlows = append(effectiveFlows, Snapshot{sequenceFlow: sequenceFlow, flowId: flowId})
								flowFuncs = append(flowFuncs, flowFunc)
							}
						}

						if len(effectiveFlows) > 0 {
							f.tracer.Send(FlowTrace{
								Source: source,
								Flows:  effectiveFlows,
							})
							for _, flowFunc := range flowFuncs {
								flowFunc()
							}
						} else {
							// no flows to continue with, abort
							f.tracer.Send(TerminationTrace{
								FlowId: f.Id(),
								Source: source,
							})
							return
						}

					} else {
						// nowhere to flow, abort
						return
					}
				case CompleteAction:
					f.tracer.Send(CompletionTrace{
						Node: f.current.Element(),
					})
					f.tracer.Send(TerminationTrace{
						FlowId: f.Id(),
						Source: f.current.Element(),
					})
					return
				case NoAction:
					f.tracer.Send(TerminationTrace{
						FlowId: f.Id(),
						Source: f.current.Element(),
					})
					return
				default:
				}
			}
		}

	}()
}
