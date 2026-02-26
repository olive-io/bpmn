/*
Copyright 2025 The bpmn Authors

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

import (
	"math"
	"math/rand"
	"sort"
	"time"
)

const (
	autoLayoutStartX           = 96.0
	autoLayoutStartY           = 96.0
	autoLayoutColumnGap        = 180.0
	autoLayoutRowGap           = 120.0
	autoLayoutProcessGap       = 180.0
	autoLayoutMinProcessHeight = 160.0
)

type AutoLayoutConfig struct {
	StartX     float64
	StartY     float64
	ColumnGap  float64
	RowGap     float64
	ProcessGap float64
}

func DefaultAutoLayoutConfig() *AutoLayoutConfig {
	return &AutoLayoutConfig{
		StartX:     autoLayoutStartX,
		StartY:     autoLayoutStartY,
		ColumnGap:  autoLayoutColumnGap,
		RowGap:     autoLayoutRowGap,
		ProcessGap: autoLayoutProcessGap,
	}
}

const (
	letterBytes   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandBytes(n int) []byte {
	source := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, source.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = source.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return b
}

const (
	DefaultExpressionLanguage = "https://github.com/expr-lang/expr"
	DefaultExporter           = "Olive Bpmn Modeler"
	DefaultTargetNamespace    = "http://bpmn.io/schema/bpmn"
)

type DefinitionBuilder struct {
	*Definitions
}

func NewDefinitionsBuilder() *DefinitionBuilder {
	def := DefaultDefinitions()
	def.ExpressionLanguageField = NewStringP(DefaultExpressionLanguage)
	def.ExporterField = NewStringP(DefaultExporter)
	def.ExporterVersionField = NewStringP("1.0.0")
	def.TargetNamespaceField = DefaultTargetNamespace
	def.IdField = NewStringP("Definitions_" + string(RandBytes(7)))

	builder := &DefinitionBuilder{Definitions: &def}
	return builder
}

func (builder *DefinitionBuilder) SetId(id string) *DefinitionBuilder {
	builder.Definitions.SetId(&id)
	return builder
}

func (builder *DefinitionBuilder) SetVersion(version string) *DefinitionBuilder {
	builder.SetExporterVersion(&version)
	return builder
}

func (builder *DefinitionBuilder) AddProcess(p Process) *DefinitionBuilder {
	if len(builder.ProcessField) == 0 {
		builder.ProcessField = make([]Process, 0)
	}
	if p.IdField == nil || *p.IdField == "" {
		p.IdField = NewStringP("Process_" + string(RandBytes(7)))
	}
	if len(builder.ProcessField) == 0 {
		p.IsExecutableField = NewBoolP(true)
	}
	builder.ProcessField = append(builder.ProcessField, p)
	if len(builder.ProcessField) > 1 {
		if len(builder.CollaborationField) == 0 {
			collaboration := Collaboration{}
			collaboration.IdField = NewStringP("Collaboration_" + string(RandBytes(7)))
			builder.CollaborationField = append(builder.CollaborationField, collaboration)
		}
		collaboration := &builder.CollaborationField[0]
		start, end := len(collaboration.ParticipantField), len(builder.ProcessField)
		for i := start; i < end; i++ {
			pf := builder.ProcessField[i]
			pid := *pf.IdField
			participant := Participant{}
			participant.IdField = NewStringP("Participant_" + string(RandBytes(7)))
			participant.SetProcessRef(NewQName(pid))
			collaboration.ParticipantField = append(collaboration.ParticipantField, participant)
		}
	}
	return builder
}

type processNodeLayout struct {
	id     string
	node   FlowNodeInterface
	order  int
	level  int
	x      float64
	y      float64
	width  float64
	height float64
}

type processFlowEdge struct {
	id     string
	source string
	target string
}

type autoLayoutBounds struct {
	x      float64
	y      float64
	width  float64
	height float64
}

func (builder *DefinitionBuilder) AutoLayout(cfg *AutoLayoutConfig) {
	if len(builder.ProcessField) == 0 {
		builder.DiagramField = nil
		return
	}

	diagram := DefaultBPMNDiagram()
	diagramID := Id("BPMNDiagram_" + string(RandBytes(7)))
	diagram.SetId(&diagramID)

	plane := DefaultBPMNPlane()
	planeID := Id("BPMNPlane_" + string(RandBytes(7)))
	plane.SetId(&planeID)

	if len(builder.ProcessField) > 1 && len(builder.CollaborationField) > 0 {
		if cid, ok := builder.CollaborationField[0].Id(); ok && cid != nil && *cid != "" {
			plane.SetBpmnElement(NewQName(*cid))
		}
	} else {
		if pid, ok := builder.ProcessField[0].Id(); ok && pid != nil && *pid != "" {
			plane.SetBpmnElement(NewQName(*pid))
		}
	}

	shapes := make([]BPMNShape, 0)
	edges := make([]BPMNEdge, 0)

	currentY := cfg.StartY
	for i := range builder.ProcessField {
		processShapes, processEdges, processHeight := buildProcessLayout(&builder.ProcessField[i], cfg, currentY)
		shapes = append(shapes, processShapes...)
		edges = append(edges, processEdges...)
		currentY += processHeight + cfg.ProcessGap
	}

	plane.BPMNShapeFields = shapes
	plane.BPMNEdgeFields = edges
	diagram.SetBPMNPlane(&plane)
	builder.DiagramField = &diagram
}

func buildProcessLayout(process *Process, cfg *AutoLayoutConfig, startY float64) (shapes []BPMNShape, edges []BPMNEdge, processHeight float64) {
	nodes := collectProcessFlowNodes(process)
	if len(nodes) == 0 {
		return nil, nil, autoLayoutMinProcessHeight
	}

	flowEdges := collectProcessFlowEdges(process)
	levelMap := computeFlowNodeLevels(process, nodes, flowEdges)
	rowMap := computeFlowNodeRows(nodes, levelMap, flowEdges)
	boundsByNode := make(map[string]autoLayoutBounds, len(nodes))
	for _, node := range nodes {
		level := levelMap[node.id]
		row := rowMap[node.id]
		node.level = level
		node.x = cfg.StartX + float64(level)*cfg.ColumnGap
		rowCenterY := startY + float64(row)*cfg.RowGap
		node.y = rowCenterY - node.height/2
		boundsByNode[node.id] = autoLayoutBounds{
			x:      node.x,
			y:      node.y,
			width:  node.width,
			height: node.height,
		}
	}

	for _, node := range nodes {
		shape := DefaultBPMNShape()
		shapeID := Id("Shape_" + string(RandBytes(7)))
		shape.SetId(&shapeID)
		shape.SetBpmnElement(NewQName(node.id))
		shape.SetBounds(newBounds(node.x, node.y, node.width, node.height))
		shapes = append(shapes, shape)
	}

	for _, flowEdge := range flowEdges {
		source, okSource := boundsByNode[flowEdge.source]
		target, okTarget := boundsByNode[flowEdge.target]
		if !okSource || !okTarget {
			continue
		}

		edge := DefaultBPMNEdge()
		edgeID := Id("Edge_" + string(RandBytes(7)))
		edge.SetId(&edgeID)

		edge.SetBpmnElement(NewQName(flowEdge.id))

		edge.SetSourceElement(NewQName(flowEdge.source))

		edge.SetTargetElement(NewQName(flowEdge.target))

		edge.SetWaypoints(buildAlignedWaypoints(source, target))
		edges = append(edges, edge)
	}

	maxBottom := startY
	for _, node := range nodes {
		bottom := node.y + node.height
		if bottom > maxBottom {
			maxBottom = bottom
		}
	}
	processHeight = maxBottom - startY
	if processHeight < autoLayoutMinProcessHeight {
		processHeight = autoLayoutMinProcessHeight
	}
	return
}

func collectProcessFlowNodes(process *Process) []*processNodeLayout {
	flowElements := process.FlowElements()
	nodes := make([]*processNodeLayout, 0)
	seen := make(map[string]struct{})

	for i := range flowElements {
		node, ok := flowElements[i].(FlowNodeInterface)
		if !ok {
			continue
		}
		id, present := node.Id()
		if !present || id == nil || *id == "" {
			continue
		}
		if _, ok := seen[*id]; ok {
			continue
		}
		seen[*id] = struct{}{}
		width, height := flowNodeDefaultSize(node)
		nodeOrder := len(nodes)
		nodes = append(nodes, &processNodeLayout{
			id:     *id,
			node:   node,
			order:  nodeOrder,
			width:  width,
			height: height,
		})
	}

	return nodes
}

func collectProcessFlowEdges(process *Process) []processFlowEdge {
	edges := make([]processFlowEdge, 0, len(process.SequenceFlowField))
	for i := range process.SequenceFlowField {
		flow := &process.SequenceFlowField[i]
		source := string(flow.SourceRefField)
		target := string(flow.TargetRefField)
		if source == "" || target == "" {
			continue
		}
		flowID := source + "_" + target
		if id, present := flow.Id(); present && id != nil && *id != "" {
			flowID = *id
		}
		edges = append(edges, processFlowEdge{id: flowID, source: source, target: target})
	}
	return edges
}

func computeFlowNodeLevels(process *Process, nodes []*processNodeLayout, edges []processFlowEdge) map[string]int {
	levels := make(map[string]int, len(nodes))
	known := make(map[string]struct{}, len(nodes))
	for _, node := range nodes {
		known[node.id] = struct{}{}
		levels[node.id] = 0
	}

	for i := range process.StartEventField {
		if id, present := process.StartEventField[i].Id(); present && id != nil {
			levels[*id] = 0
		}
	}

	iterations := len(nodes)
	for i := 0; i < iterations; i++ {
		updated := false
		for _, edge := range edges {
			if _, ok := known[edge.source]; !ok {
				continue
			}
			if _, ok := known[edge.target]; !ok {
				continue
			}
			next := levels[edge.source] + 1
			if next > levels[edge.target] {
				levels[edge.target] = next
				updated = true
			}
		}
		if !updated {
			break
		}
	}

	return levels
}

func computeFlowNodeRows(nodes []*processNodeLayout, levels map[string]int, edges []processFlowEdge) map[string]int {
	rows := make(map[string]int, len(nodes))
	incoming := make(map[string][]string)
	maxLevel := 0
	for _, edge := range edges {
		incoming[edge.target] = append(incoming[edge.target], edge.source)
	}
	for _, node := range nodes {
		if levels[node.id] > maxLevel {
			maxLevel = levels[node.id]
		}
	}

	for level := 0; level <= maxLevel; level++ {
		levelNodes := make([]*processNodeLayout, 0)
		for _, node := range nodes {
			if levels[node.id] == level {
				levelNodes = append(levelNodes, node)
			}
		}

		sort.Slice(levelNodes, func(i, j int) bool {
			left := desiredRow(levelNodes[i], incoming, rows)
			right := desiredRow(levelNodes[j], incoming, rows)
			if left == right {
				return levelNodes[i].order < levelNodes[j].order
			}
			return left < right
		})

		occupied := make(map[int]struct{})
		for _, node := range levelNodes {
			r := int(math.Round(desiredRow(node, incoming, rows)))
			if r < 0 {
				r = 0
			}
			for {
				if _, exists := occupied[r]; !exists {
					break
				}
				r++
			}
			rows[node.id] = r
			occupied[r] = struct{}{}
		}
	}

	return rows
}

func desiredRow(node *processNodeLayout, incoming map[string][]string, rows map[string]int) float64 {
	predecessors := incoming[node.id]
	if len(predecessors) == 0 {
		return 0
	}

	total := 0.0
	count := 0
	for _, source := range predecessors {
		row, ok := rows[source]
		if !ok {
			continue
		}
		total += float64(row)
		count++
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

func flowNodeDefaultSize(node FlowNodeInterface) (width, height float64) {
	switch node.(type) {
	case *StartEvent, *EndEvent, *BoundaryEvent, *IntermediateCatchEvent, *IntermediateThrowEvent, *Event, *ImplicitThrowEvent:
		return 36, 36
	case *ExclusiveGateway, *InclusiveGateway, *ParallelGateway, *EventBasedGateway, *ComplexGateway:
		return 50, 50
	case *SubProcess, *AdHocSubProcess, *Transaction:
		return 120, 100
	default:
		return 100, 80
	}
}

func newBounds(x, y, width, height float64) *Bounds {
	bounds := DefaultBounds()
	bounds.SetX(x)
	bounds.SetY(y)
	bounds.SetWidth(width)
	bounds.SetHeight(height)
	return &bounds
}

func newPoint(x, y float64) Point {
	point := DefaultPoint()
	point.SetX(x)
	point.SetY(y)
	return point
}

func buildAlignedWaypoints(source, target autoLayoutBounds) []Point {
	startX := source.x + source.width
	startY := source.y + source.height/2
	endX := target.x
	endY := target.y + target.height/2

	if math.Abs(startY-endY) < 0.001 {
		return []Point{newPoint(startX, startY), newPoint(endX, endY)}
	}

	midX := (startX + endX) / 2
	return []Point{
		newPoint(startX, startY),
		newPoint(midX, startY),
		newPoint(midX, endY),
		newPoint(endX, endY),
	}
}

// Out returns *Definitions and reset internal Definitions values
func (builder *DefinitionBuilder) Out() *Definitions {
	out := &Definitions{}
	*out = *(builder.Definitions)
	b := NewDefinitionsBuilder()
	*builder = *b
	return out
}

type ProcessBuilder struct {
	*Process
	ptr FlowNodeInterface
}

func NewProcessBuilder() *ProcessBuilder {
	pro := DefaultProcess()
	pro.IdField = NewStringP("Process_" + string(RandBytes(7)))
	start := StartEvent{}
	start.IdField = NewStringP("Event_" + string(RandBytes(7)))
	start.OutgoingField = make([]QName, 0)
	pro.StartEventField = append(pro.StartEventField, start)
	builder := &ProcessBuilder{Process: &pro, ptr: &start}
	return builder
}

func (builder *ProcessBuilder) AddActivity(act ActivityInterface) *ProcessBuilder {
	actId, _ := act.Id()
	if actId == nil || *actId == "" {
		act.SetId(NewStringP("Activity_" + string(RandBytes(7))))
	}
	b := builder.link(act)
	switch tv := act.(type) {
	case *Task:
		b.TaskField = append(b.TaskField, *tv)
	case *BusinessRuleTask:
		b.BusinessRuleTaskField = append(b.BusinessRuleTaskField, *tv)
	case *UserTask:
		b.UserTaskField = append(b.UserTaskField, *tv)
	case *CallActivity:
		b.CallActivityField = append(b.CallActivityField, *tv)
	case *ManualTask:
		b.ManualTaskField = append(b.ManualTaskField, *tv)
	case *SendTask:
		b.SendTaskField = append(b.SendTaskField, *tv)
	case *ScriptTask:
		b.ScriptTaskField = append(b.ScriptTaskField, *tv)
	case *ServiceTask:
		b.ServiceTaskField = append(b.ServiceTaskField, *tv)
	case *ReceiveTask:
		b.ReceiveTaskField = append(b.ReceiveTaskField, *tv)
	case *SubProcess:
		b.SubProcessField = append(b.SubProcessField, *tv)
	}
	return b
}

func (builder *ProcessBuilder) link(node FlowNodeInterface) *ProcessBuilder {
	sourceId, _ := builder.ptr.Id()
	targetId, _ := node.Id()
	sf := &SequenceFlow{}
	sid := "Flow_" + string(RandBytes(7))
	sf.IdField = NewStringP(sid)
	sf.SourceRefField = *sourceId
	sf.TargetRefField = *targetId
	outgoings := builder.ptr.Outgoings()
	if outgoings == nil {
		names := make([]QName, 0)
		outgoings = &names
	}
	*outgoings = append(*outgoings, QName(sid))
	builder.ptr.SetOutgoings(*outgoings)
	e, found := builder.FindBy(ExactId(*sourceId))
	if found {
		vv, ok := e.(FlowNodeInterface)
		if ok {
			vv.SetOutgoings(*outgoings)
		}
	}

	incomings := node.Incomings()
	if incomings == nil {
		names := make([]QName, 0)
		incomings = &names
	}
	*incomings = append(*incomings, QName(sid))
	node.SetIncomings(*incomings)
	builder.SequenceFlowField = append(builder.SequenceFlowField, *sf)
	builder.ptr = node
	return builder
}

func (builder *ProcessBuilder) Out() *Process {
	endEvent := DefaultEndEvent()
	endEvent.IdField = NewStringP("Event_" + string(RandBytes(7)))
	builder.link(&endEvent)
	builder.EndEventField = append(builder.EndEventField, endEvent)

	out := &Process{}
	*out = *(builder.Process)
	b := NewProcessBuilder()
	*builder = *b
	return out
}
