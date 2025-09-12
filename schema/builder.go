/*
   Copyright 2025 The  Authors

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

package schema

import (
	"math/rand"
	"time"
)

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
	if p.IdField != nil || *p.IdField == "" {
		p.IdField = NewStringP("Process_" + string(RandBytes(7)))
	}
	if len(builder.ProcessField) == 0 {
		p.IsExecutableField = NewBoolP(true)
	}
	builder.ProcessField = append(builder.ProcessField, p)
	if len(builder.ProcessField) > 1 {
		var collaboration *Collaboration
		if len(builder.CollaborationField) == 0 {
			collaboration = &Collaboration{}
			collaboration.IdField = NewStringP("Collaboration_" + string(RandBytes(7)))
		} else {
			collaboration = &builder.CollaborationField[0]
		}
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
