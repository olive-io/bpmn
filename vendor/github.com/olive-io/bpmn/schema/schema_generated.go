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

package schema

// This file is generated from BPMN 2.0 schema using `make generate`
// DO NOT EDIT
import (
	"encoding/xml"
	"math/big"
)

type AdHocOrdering string
type AssociationDirection string
type ChoreographyLoopType string
type EventBasedGatewayType string
type GatewayDirection string
type Implementation AnyURI
type ItemKind string
type MultiInstanceFlowCondition string
type ProcessType string
type RelationshipDirection string
type TransactionMethod AnyURI
type Definitions struct {
	IdField                     *Id                      `xml:"id,attr,omitempty"`
	NameField                   *string                  `xml:"name,attr,omitempty"`
	TargetNamespaceField        AnyURI                   `xml:"targetNamespace,attr,omitempty"`
	ExpressionLanguageField     *AnyURI                  `xml:"expressionLanguage,attr,omitempty"`
	TypeLanguageField           *AnyURI                  `xml:"typeLanguage,attr,omitempty"`
	ExporterField               *string                  `xml:"exporter,attr,omitempty"`
	ExporterVersionField        *string                  `xml:"exporterVersion,attr,omitempty"`
	ImportField                 []Import                 `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL import"`
	ExtensionField              []Extension              `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL extension"`
	CategoryField               []Category               `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL category"`
	CollaborationField          []Collaboration          `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL collaboration"`
	CorrelationPropertyField    []CorrelationProperty    `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL correlationProperty"`
	DataStoreField              []DataStore              `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataStore"`
	EndPointField               []EndPoint               `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL endPoint"`
	ErrorField                  []Error                  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL error"`
	EscalationField             []Escalation             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL escalation"`
	EventDefinitionField        []EventDefinition        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL eventDefinition"`
	GlobalBusinessRuleTaskField []GlobalBusinessRuleTask `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL globalBusinessRuleTask"`
	GlobalManualTaskField       []GlobalManualTask       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL globalManualTask"`
	GlobalScriptTaskField       []GlobalScriptTask       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL globalScriptTask"`
	GlobalTaskField             []GlobalTask             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL globalTask"`
	GlobalUserTaskField         []GlobalUserTask         `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL globalUserTask"`
	InterfaceField              []Interface              `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL interface"`
	ItemDefinitionField         []ItemDefinition         `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL itemDefinition"`
	MessageField                []Message                `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL message"`
	PartnerEntityField          []PartnerEntity          `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL partnerEntity"`
	PartnerRoleField            []PartnerRole            `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL partnerRole"`
	ProcessField                []Process                `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL process"`
	ResourceField               []Resource               `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL resource"`
	SignalField                 []Signal                 `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL signal"`
	RelationshipField           []Relationship           `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL relationship"`
	DiagramField                *BPMNDiagram             `xml:"http://www.omg.org/spec/BPMN/20100524/DI BPMNDiagram"`
	TextPayloadField            *Payload                 `xml:",chardata"`
}

var defaultDefinitionsExpressionLanguageField AnyURI = "http://www.w3.org/1999/XPath"
var defaultDefinitionsTypeLanguageField AnyURI = "http://www.w3.org/2001/XMLSchema"

func DefaultDefinitions() Definitions {
	return Definitions{
		ExpressionLanguageField: &defaultDefinitionsExpressionLanguageField,
		TypeLanguageField:       &defaultDefinitionsTypeLanguageField,
	}
}

type DefinitionsInterface interface {
	Element
	Id() (result *Id, present bool)
	Name() (result *string, present bool)
	TargetNamespace() (result *AnyURI)
	ExpressionLanguage() (result *AnyURI)
	TypeLanguage() (result *AnyURI)
	Exporter() (result *string, present bool)
	ExporterVersion() (result *string, present bool)
	Imports() (result *[]Import)
	Extensions() (result *[]Extension)
	Categories() (result *[]Category)
	Collaborations() (result *[]Collaboration)
	CorrelationProperties() (result *[]CorrelationProperty)
	DataStores() (result *[]DataStore)
	EndPoints() (result *[]EndPoint)
	Errors() (result *[]Error)
	Escalations() (result *[]Escalation)
	EventDefinitions() (result *[]EventDefinition)
	GlobalBusinessRuleTasks() (result *[]GlobalBusinessRuleTask)
	GlobalManualTasks() (result *[]GlobalManualTask)
	GlobalScriptTasks() (result *[]GlobalScriptTask)
	GlobalTasks() (result *[]GlobalTask)
	GlobalUserTasks() (result *[]GlobalUserTask)
	Interfaces() (result *[]Interface)
	ItemDefinitions() (result *[]ItemDefinition)
	Messages() (result *[]Message)
	PartnerEntities() (result *[]PartnerEntity)
	PartnerRoles() (result *[]PartnerRole)
	Processes() (result *[]Process)
	Resources() (result *[]Resource)
	Signals() (result *[]Signal)
	Relationships() (result *[]Relationship)
	RootElements() []RootElementInterface
	SetId(value *Id)
	SetName(value *string)
	SetTargetNamespace(value AnyURI)
	SetExpressionLanguage(value *AnyURI)
	SetTypeLanguage(value *AnyURI)
	SetExporter(value *string)
	SetExporterVersion(value *string)
	SetImports(value []Import)
	SetExtensions(value []Extension)
	SetCategories(value []Category)
	SetCollaborations(value []Collaboration)
	SetCorrelationProperties(value []CorrelationProperty)
	SetDataStores(value []DataStore)
	SetEndPoints(value []EndPoint)
	SetErrors(value []Error)
	SetEscalations(value []Escalation)
	SetEventDefinitions(value []EventDefinition)
	SetGlobalBusinessRuleTasks(value []GlobalBusinessRuleTask)
	SetGlobalManualTasks(value []GlobalManualTask)
	SetGlobalScriptTasks(value []GlobalScriptTask)
	SetGlobalTasks(value []GlobalTask)
	SetGlobalUserTasks(value []GlobalUserTask)
	SetInterfaces(value []Interface)
	SetItemDefinitions(value []ItemDefinition)
	SetMessages(value []Message)
	SetPartnerEntities(value []PartnerEntity)
	SetPartnerRoles(value []PartnerRole)
	SetProcesses(value []Process)
	SetResources(value []Resource)
	SetSignals(value []Signal)
	SetRelationships(value []Relationship)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Definitions) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Definitions) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Definitions) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}

	for i := range t.ImportField {
		if result, found = t.ImportField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ExtensionField {
		if result, found = t.ExtensionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.CategoryField {
		if result, found = t.CategoryField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.CollaborationField {
		if result, found = t.CollaborationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.CorrelationPropertyField {
		if result, found = t.CorrelationPropertyField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataStoreField {
		if result, found = t.DataStoreField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.EndPointField {
		if result, found = t.EndPointField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ErrorField {
		if result, found = t.ErrorField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.EscalationField {
		if result, found = t.EscalationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.EventDefinitionField {
		if result, found = t.EventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.GlobalBusinessRuleTaskField {
		if result, found = t.GlobalBusinessRuleTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.GlobalManualTaskField {
		if result, found = t.GlobalManualTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.GlobalScriptTaskField {
		if result, found = t.GlobalScriptTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.GlobalTaskField {
		if result, found = t.GlobalTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.GlobalUserTaskField {
		if result, found = t.GlobalUserTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.InterfaceField {
		if result, found = t.InterfaceField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ItemDefinitionField {
		if result, found = t.ItemDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.MessageField {
		if result, found = t.MessageField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.PartnerEntityField {
		if result, found = t.PartnerEntityField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.PartnerRoleField {
		if result, found = t.PartnerRoleField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ProcessField {
		if result, found = t.ProcessField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ResourceField {
		if result, found = t.ResourceField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SignalField {
		if result, found = t.SignalField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.RelationshipField {
		if result, found = t.RelationshipField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *Definitions) Id() (result *Id, present bool) {
	if t.IdField != nil {
		present = true
	}
	result = t.IdField
	return
}
func (t *Definitions) SetId(value *Id) {
	t.IdField = value
}
func (t *Definitions) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *Definitions) SetName(value *string) {
	t.NameField = value
}
func (t *Definitions) TargetNamespace() (result *AnyURI) {
	result = &t.TargetNamespaceField
	return
}
func (t *Definitions) SetTargetNamespace(value AnyURI) {
	t.TargetNamespaceField = value
}
func (t *Definitions) ExpressionLanguage() (result *AnyURI) {
	if t.ExpressionLanguageField == nil {
		result = &defaultDefinitionsExpressionLanguageField
		return
	}
	result = t.ExpressionLanguageField
	return
}
func (t *Definitions) SetExpressionLanguage(value *AnyURI) {
	t.ExpressionLanguageField = value
}
func (t *Definitions) TypeLanguage() (result *AnyURI) {
	if t.TypeLanguageField == nil {
		result = &defaultDefinitionsTypeLanguageField
		return
	}
	result = t.TypeLanguageField
	return
}
func (t *Definitions) SetTypeLanguage(value *AnyURI) {
	t.TypeLanguageField = value
}
func (t *Definitions) Exporter() (result *string, present bool) {
	if t.ExporterField != nil {
		present = true
	}
	result = t.ExporterField
	return
}
func (t *Definitions) SetExporter(value *string) {
	t.ExporterField = value
}
func (t *Definitions) ExporterVersion() (result *string, present bool) {
	if t.ExporterVersionField != nil {
		present = true
	}
	result = t.ExporterVersionField
	return
}
func (t *Definitions) SetExporterVersion(value *string) {
	t.ExporterVersionField = value
}
func (t *Definitions) RootElements() []RootElementInterface {

	result := make([]RootElementInterface, 0)

	for i := range t.CategoryField {
		result = append(result, &t.CategoryField[i])
	}

	for i := range t.CollaborationField {
		result = append(result, &t.CollaborationField[i])
	}

	for i := range t.CorrelationPropertyField {
		result = append(result, &t.CorrelationPropertyField[i])
	}

	for i := range t.DataStoreField {
		result = append(result, &t.DataStoreField[i])
	}

	for i := range t.EndPointField {
		result = append(result, &t.EndPointField[i])
	}

	for i := range t.ErrorField {
		result = append(result, &t.ErrorField[i])
	}

	for i := range t.EscalationField {
		result = append(result, &t.EscalationField[i])
	}

	for i := range t.EventDefinitionField {
		result = append(result, &t.EventDefinitionField[i])
	}

	for i := range t.GlobalBusinessRuleTaskField {
		result = append(result, &t.GlobalBusinessRuleTaskField[i])
	}

	for i := range t.GlobalManualTaskField {
		result = append(result, &t.GlobalManualTaskField[i])
	}

	for i := range t.GlobalScriptTaskField {
		result = append(result, &t.GlobalScriptTaskField[i])
	}

	for i := range t.GlobalTaskField {
		result = append(result, &t.GlobalTaskField[i])
	}

	for i := range t.GlobalUserTaskField {
		result = append(result, &t.GlobalUserTaskField[i])
	}

	for i := range t.InterfaceField {
		result = append(result, &t.InterfaceField[i])
	}

	for i := range t.ItemDefinitionField {
		result = append(result, &t.ItemDefinitionField[i])
	}

	for i := range t.MessageField {
		result = append(result, &t.MessageField[i])
	}

	for i := range t.PartnerEntityField {
		result = append(result, &t.PartnerEntityField[i])
	}

	for i := range t.PartnerRoleField {
		result = append(result, &t.PartnerRoleField[i])
	}

	for i := range t.ProcessField {
		result = append(result, &t.ProcessField[i])
	}

	for i := range t.ResourceField {
		result = append(result, &t.ResourceField[i])
	}

	for i := range t.SignalField {
		result = append(result, &t.SignalField[i])
	}
	return result
}
func (t *Definitions) Imports() (result *[]Import) {
	result = &t.ImportField
	return
}
func (t *Definitions) SetImports(value []Import) {
	t.ImportField = value
}
func (t *Definitions) Extensions() (result *[]Extension) {
	result = &t.ExtensionField
	return
}
func (t *Definitions) SetExtensions(value []Extension) {
	t.ExtensionField = value
}
func (t *Definitions) Categories() (result *[]Category) {
	result = &t.CategoryField
	return
}
func (t *Definitions) SetCategories(value []Category) {
	t.CategoryField = value
}
func (t *Definitions) Collaborations() (result *[]Collaboration) {
	result = &t.CollaborationField
	return
}
func (t *Definitions) SetCollaborations(value []Collaboration) {
	t.CollaborationField = value
}
func (t *Definitions) CorrelationProperties() (result *[]CorrelationProperty) {
	result = &t.CorrelationPropertyField
	return
}
func (t *Definitions) SetCorrelationProperties(value []CorrelationProperty) {
	t.CorrelationPropertyField = value
}
func (t *Definitions) DataStores() (result *[]DataStore) {
	result = &t.DataStoreField
	return
}
func (t *Definitions) SetDataStores(value []DataStore) {
	t.DataStoreField = value
}
func (t *Definitions) EndPoints() (result *[]EndPoint) {
	result = &t.EndPointField
	return
}
func (t *Definitions) SetEndPoints(value []EndPoint) {
	t.EndPointField = value
}
func (t *Definitions) Errors() (result *[]Error) {
	result = &t.ErrorField
	return
}
func (t *Definitions) SetErrors(value []Error) {
	t.ErrorField = value
}
func (t *Definitions) Escalations() (result *[]Escalation) {
	result = &t.EscalationField
	return
}
func (t *Definitions) SetEscalations(value []Escalation) {
	t.EscalationField = value
}
func (t *Definitions) EventDefinitions() (result *[]EventDefinition) {
	result = &t.EventDefinitionField
	return
}
func (t *Definitions) SetEventDefinitions(value []EventDefinition) {
	t.EventDefinitionField = value
}
func (t *Definitions) GlobalBusinessRuleTasks() (result *[]GlobalBusinessRuleTask) {
	result = &t.GlobalBusinessRuleTaskField
	return
}
func (t *Definitions) SetGlobalBusinessRuleTasks(value []GlobalBusinessRuleTask) {
	t.GlobalBusinessRuleTaskField = value
}
func (t *Definitions) GlobalManualTasks() (result *[]GlobalManualTask) {
	result = &t.GlobalManualTaskField
	return
}
func (t *Definitions) SetGlobalManualTasks(value []GlobalManualTask) {
	t.GlobalManualTaskField = value
}
func (t *Definitions) GlobalScriptTasks() (result *[]GlobalScriptTask) {
	result = &t.GlobalScriptTaskField
	return
}
func (t *Definitions) SetGlobalScriptTasks(value []GlobalScriptTask) {
	t.GlobalScriptTaskField = value
}
func (t *Definitions) GlobalTasks() (result *[]GlobalTask) {
	result = &t.GlobalTaskField
	return
}
func (t *Definitions) SetGlobalTasks(value []GlobalTask) {
	t.GlobalTaskField = value
}
func (t *Definitions) GlobalUserTasks() (result *[]GlobalUserTask) {
	result = &t.GlobalUserTaskField
	return
}
func (t *Definitions) SetGlobalUserTasks(value []GlobalUserTask) {
	t.GlobalUserTaskField = value
}
func (t *Definitions) Interfaces() (result *[]Interface) {
	result = &t.InterfaceField
	return
}
func (t *Definitions) SetInterfaces(value []Interface) {
	t.InterfaceField = value
}
func (t *Definitions) ItemDefinitions() (result *[]ItemDefinition) {
	result = &t.ItemDefinitionField
	return
}
func (t *Definitions) SetItemDefinitions(value []ItemDefinition) {
	t.ItemDefinitionField = value
}
func (t *Definitions) Messages() (result *[]Message) {
	result = &t.MessageField
	return
}
func (t *Definitions) SetMessages(value []Message) {
	t.MessageField = value
}
func (t *Definitions) PartnerEntities() (result *[]PartnerEntity) {
	result = &t.PartnerEntityField
	return
}
func (t *Definitions) SetPartnerEntities(value []PartnerEntity) {
	t.PartnerEntityField = value
}
func (t *Definitions) PartnerRoles() (result *[]PartnerRole) {
	result = &t.PartnerRoleField
	return
}
func (t *Definitions) SetPartnerRoles(value []PartnerRole) {
	t.PartnerRoleField = value
}
func (t *Definitions) Processes() (result *[]Process) {
	result = &t.ProcessField
	return
}
func (t *Definitions) SetProcesses(value []Process) {
	t.ProcessField = value
}
func (t *Definitions) Resources() (result *[]Resource) {
	result = &t.ResourceField
	return
}
func (t *Definitions) SetResources(value []Resource) {
	t.ResourceField = value
}
func (t *Definitions) Signals() (result *[]Signal) {
	result = &t.SignalField
	return
}
func (t *Definitions) SetSignals(value []Signal) {
	t.SignalField = value
}
func (t *Definitions) Relationships() (result *[]Relationship) {
	result = &t.RelationshipField
	return
}
func (t *Definitions) SetRelationships(value []Relationship) {
	t.RelationshipField = value
}

type Import struct {
	NamespaceField   AnyURI   `xml:"namespace,attr,omitempty"`
	LocationField    string   `xml:"location,attr,omitempty"`
	ImportTypeField  AnyURI   `xml:"importType,attr,omitempty"`
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultImport() Import {
	return Import{}
}

type ImportInterface interface {
	Element
	Namespace() (result *AnyURI)
	Location() (result *string)
	ImportType() (result *AnyURI)
	SetNamespace(value AnyURI)
	SetLocation(value string)
	SetImportType(value AnyURI)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Import) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Import) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Import) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}

	return
}
func (t *Import) Namespace() (result *AnyURI) {
	result = &t.NamespaceField
	return
}
func (t *Import) SetNamespace(value AnyURI) {
	t.NamespaceField = value
}
func (t *Import) Location() (result *string) {
	result = &t.LocationField
	return
}
func (t *Import) SetLocation(value string) {
	t.LocationField = value
}
func (t *Import) ImportType() (result *AnyURI) {
	result = &t.ImportTypeField
	return
}
func (t *Import) SetImportType(value AnyURI) {
	t.ImportTypeField = value
}

type Activity struct {
	FlowNode
	IsForCompensationField                *bool                             `xml:"isForCompensation,attr,omitempty"`
	StartQuantityField                    *big.Int                          `xml:"startQuantity,attr,omitempty"`
	CompletionQuantityField               *big.Int                          `xml:"completionQuantity,attr,omitempty"`
	DefaultField                          *IdRef                            `xml:"default,attr,omitempty"`
	IoSpecificationField                  *InputOutputSpecification         `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL ioSpecification"`
	PropertyField                         []Property                        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL property"`
	DataInputAssociationField             []DataInputAssociation            `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataInputAssociation"`
	DataOutputAssociationField            []DataOutputAssociation           `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataOutputAssociation"`
	ResourceRoleField                     []ResourceRole                    `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL resourceRole"`
	MultiInstanceLoopCharacteristicsField *MultiInstanceLoopCharacteristics `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL multiInstanceLoopCharacteristics"`
	StandardLoopCharacteristicsField      *StandardLoopCharacteristics      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL standardLoopCharacteristics"`
}

var defaultActivityIsForCompensationField bool = false
var defaultActivityStartQuantityField big.Int = *big.NewInt(1)
var defaultActivityCompletionQuantityField big.Int = *big.NewInt(1)

func DefaultActivity() Activity {
	return Activity{
		FlowNode:                DefaultFlowNode(),
		IsForCompensationField:  &defaultActivityIsForCompensationField,
		StartQuantityField:      &defaultActivityStartQuantityField,
		CompletionQuantityField: &defaultActivityCompletionQuantityField,
	}
}

type ActivityInterface interface {
	Element
	FlowNodeInterface
	IsForCompensation() (result bool)
	StartQuantity() (result *big.Int)
	CompletionQuantity() (result *big.Int)
	Default() (result *IdRef, present bool)
	IoSpecification() (result *InputOutputSpecification, present bool)
	Properties() (result *[]Property)
	DataInputAssociations() (result *[]DataInputAssociation)
	DataOutputAssociations() (result *[]DataOutputAssociation)
	ResourceRoles() (result *[]ResourceRole)
	MultiInstanceLoopCharacteristics() (result *MultiInstanceLoopCharacteristics, present bool)
	StandardLoopCharacteristics() (result *StandardLoopCharacteristics, present bool)
	LoopCharacteristics() LoopCharacteristicsInterface
	SetIsForCompensation(value *bool)
	SetStartQuantity(value *big.Int)
	SetCompletionQuantity(value *big.Int)
	SetDefault(value *IdRef)
	SetIoSpecification(value *InputOutputSpecification)
	SetProperties(value []Property)
	SetDataInputAssociations(value []DataInputAssociation)
	SetDataOutputAssociations(value []DataOutputAssociation)
	SetResourceRoles(value []ResourceRole)
	SetMultiInstanceLoopCharacteristics(value *MultiInstanceLoopCharacteristics)
	SetStandardLoopCharacteristics(value *StandardLoopCharacteristics)
}

func (t *Activity) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.FlowNode.FindBy(f); found {
		return
	}

	if value := t.IoSpecificationField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	for i := range t.PropertyField {
		if result, found = t.PropertyField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataInputAssociationField {
		if result, found = t.DataInputAssociationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataOutputAssociationField {
		if result, found = t.DataOutputAssociationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ResourceRoleField {
		if result, found = t.ResourceRoleField[i].FindBy(f); found {
			return
		}
	}

	if result, found = t.MultiInstanceLoopCharacteristicsField.FindBy(f); found {
		return
	}

	if result, found = t.StandardLoopCharacteristicsField.FindBy(f); found {
		return
	}

	return
}
func (t *Activity) IsForCompensation() (result bool) {
	if t.IsForCompensationField == nil {
		result = defaultActivityIsForCompensationField
		return
	}
	result = *t.IsForCompensationField
	return
}
func (t *Activity) SetIsForCompensation(value *bool) {
	t.IsForCompensationField = value
}
func (t *Activity) StartQuantity() (result *big.Int) {
	if t.StartQuantityField == nil {
		result = &defaultActivityStartQuantityField
		return
	}
	result = t.StartQuantityField
	return
}
func (t *Activity) SetStartQuantity(value *big.Int) {
	t.StartQuantityField = value
}
func (t *Activity) CompletionQuantity() (result *big.Int) {
	if t.CompletionQuantityField == nil {
		result = &defaultActivityCompletionQuantityField
		return
	}
	result = t.CompletionQuantityField
	return
}
func (t *Activity) SetCompletionQuantity(value *big.Int) {
	t.CompletionQuantityField = value
}
func (t *Activity) Default() (result *IdRef, present bool) {
	if t.DefaultField != nil {
		present = true
	}
	result = t.DefaultField
	return
}
func (t *Activity) SetDefault(value *IdRef) {
	t.DefaultField = value
}
func (t *Activity) LoopCharacteristics() LoopCharacteristicsInterface {
	if t.MultiInstanceLoopCharacteristicsField != nil {
		return t.MultiInstanceLoopCharacteristicsField
	}
	if t.StandardLoopCharacteristicsField != nil {
		return t.StandardLoopCharacteristicsField
	}
	return nil
}
func (t *Activity) IoSpecification() (result *InputOutputSpecification, present bool) {
	if t.IoSpecificationField != nil {
		present = true
	}
	result = t.IoSpecificationField
	return
}
func (t *Activity) SetIoSpecification(value *InputOutputSpecification) {
	t.IoSpecificationField = value
}
func (t *Activity) Properties() (result *[]Property) {
	result = &t.PropertyField
	return
}
func (t *Activity) SetProperties(value []Property) {
	t.PropertyField = value
}
func (t *Activity) DataInputAssociations() (result *[]DataInputAssociation) {
	result = &t.DataInputAssociationField
	return
}
func (t *Activity) SetDataInputAssociations(value []DataInputAssociation) {
	t.DataInputAssociationField = value
}
func (t *Activity) DataOutputAssociations() (result *[]DataOutputAssociation) {
	result = &t.DataOutputAssociationField
	return
}
func (t *Activity) SetDataOutputAssociations(value []DataOutputAssociation) {
	t.DataOutputAssociationField = value
}
func (t *Activity) ResourceRoles() (result *[]ResourceRole) {
	result = &t.ResourceRoleField
	return
}
func (t *Activity) SetResourceRoles(value []ResourceRole) {
	t.ResourceRoleField = value
}
func (t *Activity) MultiInstanceLoopCharacteristics() (result *MultiInstanceLoopCharacteristics, present bool) {
	if t.MultiInstanceLoopCharacteristicsField != nil {
		present = true
	}
	result = t.MultiInstanceLoopCharacteristicsField
	return
}
func (t *Activity) SetMultiInstanceLoopCharacteristics(value *MultiInstanceLoopCharacteristics) {
	t.MultiInstanceLoopCharacteristicsField = value
}
func (t *Activity) StandardLoopCharacteristics() (result *StandardLoopCharacteristics, present bool) {
	if t.StandardLoopCharacteristicsField != nil {
		present = true
	}
	result = t.StandardLoopCharacteristicsField
	return
}
func (t *Activity) SetStandardLoopCharacteristics(value *StandardLoopCharacteristics) {
	t.StandardLoopCharacteristicsField = value
}

type AdHocSubProcess struct {
	SubProcess
	CancelRemainingInstancesField *bool          `xml:"cancelRemainingInstances,attr,omitempty"`
	OrderingField                 *AdHocOrdering `xml:"ordering,attr,omitempty"`
	CompletionConditionField      *AnExpression  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL completionCondition"`
	TextPayloadField              *Payload       `xml:",chardata"`
}

var defaultAdHocSubProcessCancelRemainingInstancesField bool = true

func DefaultAdHocSubProcess() AdHocSubProcess {
	return AdHocSubProcess{
		SubProcess:                    DefaultSubProcess(),
		CancelRemainingInstancesField: &defaultAdHocSubProcessCancelRemainingInstancesField,
	}
}

type AdHocSubProcessInterface interface {
	Element
	SubProcessInterface
	CancelRemainingInstances() (result bool)
	Ordering() (result *AdHocOrdering, present bool)
	CompletionCondition() (result *AnExpression, present bool)
	SetCancelRemainingInstances(value *bool)
	SetOrdering(value *AdHocOrdering)
	SetCompletionCondition(value *AnExpression)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *AdHocSubProcess) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *AdHocSubProcess) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *AdHocSubProcess) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.SubProcess.FindBy(f); found {
		return
	}

	if value := t.CompletionConditionField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	return
}
func (t *AdHocSubProcess) CancelRemainingInstances() (result bool) {
	if t.CancelRemainingInstancesField == nil {
		result = defaultAdHocSubProcessCancelRemainingInstancesField
		return
	}
	result = *t.CancelRemainingInstancesField
	return
}
func (t *AdHocSubProcess) SetCancelRemainingInstances(value *bool) {
	t.CancelRemainingInstancesField = value
}
func (t *AdHocSubProcess) Ordering() (result *AdHocOrdering, present bool) {
	if t.OrderingField != nil {
		present = true
	}
	result = t.OrderingField
	return
}
func (t *AdHocSubProcess) SetOrdering(value *AdHocOrdering) {
	t.OrderingField = value
}
func (t *AdHocSubProcess) CompletionCondition() (result *AnExpression, present bool) {
	if t.CompletionConditionField != nil {
		present = true
	}
	result = t.CompletionConditionField
	return
}
func (t *AdHocSubProcess) SetCompletionCondition(value *AnExpression) {
	t.CompletionConditionField = value
}

type Artifact struct {
	BaseElement
}

func DefaultArtifact() Artifact {
	return Artifact{
		BaseElement: DefaultBaseElement(),
	}
}

type ArtifactInterface interface {
	Element
	BaseElementInterface
}

func (t *Artifact) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}

type Assignment struct {
	BaseElement
	FromField        AnExpression `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL from"`
	ToField          AnExpression `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL to"`
	TextPayloadField *Payload     `xml:",chardata"`
}

func DefaultAssignment() Assignment {
	return Assignment{
		BaseElement: DefaultBaseElement(),
	}
}

type AssignmentInterface interface {
	Element
	BaseElementInterface
	From() (result *AnExpression)
	To() (result *AnExpression)
	SetFrom(value AnExpression)
	SetTo(value AnExpression)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Assignment) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Assignment) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Assignment) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	if result, found = t.FromField.FindBy(f); found {
		return
	}

	if result, found = t.ToField.FindBy(f); found {
		return
	}

	return
}
func (t *Assignment) From() (result *AnExpression) {
	result = &t.FromField
	return
}
func (t *Assignment) SetFrom(value AnExpression) {
	t.FromField = value
}
func (t *Assignment) To() (result *AnExpression) {
	result = &t.ToField
	return
}
func (t *Assignment) SetTo(value AnExpression) {
	t.ToField = value
}

type Association struct {
	Artifact
	SourceRefField            QName                 `xml:"sourceRef,attr,omitempty"`
	TargetRefField            QName                 `xml:"targetRef,attr,omitempty"`
	AssociationDirectionField *AssociationDirection `xml:"associationDirection,attr,omitempty"`
	TextPayloadField          *Payload              `xml:",chardata"`
}

var defaultAssociationAssociationDirectionField AssociationDirection = "None"

func DefaultAssociation() Association {
	return Association{
		Artifact:                  DefaultArtifact(),
		AssociationDirectionField: &defaultAssociationAssociationDirectionField,
	}
}

type AssociationInterface interface {
	Element
	ArtifactInterface
	SourceRef() (result *QName)
	TargetRef() (result *QName)
	AssociationDirection() (result *AssociationDirection)
	SetSourceRef(value QName)
	SetTargetRef(value QName)
	SetAssociationDirection(value *AssociationDirection)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Association) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Association) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Association) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Artifact.FindBy(f); found {
		return
	}

	return
}
func (t *Association) SourceRef() (result *QName) {
	result = &t.SourceRefField
	return
}
func (t *Association) SetSourceRef(value QName) {
	t.SourceRefField = value
}
func (t *Association) TargetRef() (result *QName) {
	result = &t.TargetRefField
	return
}
func (t *Association) SetTargetRef(value QName) {
	t.TargetRefField = value
}
func (t *Association) AssociationDirection() (result *AssociationDirection) {
	if t.AssociationDirectionField == nil {
		result = &defaultAssociationAssociationDirectionField
		return
	}
	result = t.AssociationDirectionField
	return
}
func (t *Association) SetAssociationDirection(value *AssociationDirection) {
	t.AssociationDirectionField = value
}

type Auditing struct {
	BaseElement
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultAuditing() Auditing {
	return Auditing{
		BaseElement: DefaultBaseElement(),
	}
}

type AuditingInterface interface {
	Element
	BaseElementInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Auditing) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Auditing) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Auditing) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}

type BaseElement struct {
	IdField                *Id                `xml:"id,attr,omitempty"`
	DocumentationField     []Documentation    `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL documentation"`
	ExtensionElementsField *ExtensionElements `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL extensionElements"`
}

func DefaultBaseElement() BaseElement {
	return BaseElement{}
}

type BaseElementInterface interface {
	Element
	Id() (result *Id, present bool)
	Documentations() (result *[]Documentation)
	ExtensionElements() (result *ExtensionElements, present bool)
	SetId(value *Id)
	SetDocumentations(value []Documentation)
	SetExtensionElements(value *ExtensionElements)
}

func (t *BaseElement) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}

	for i := range t.DocumentationField {
		if result, found = t.DocumentationField[i].FindBy(f); found {
			return
		}
	}

	if value := t.ExtensionElementsField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	return
}
func (t *BaseElement) Id() (result *Id, present bool) {
	if t.IdField != nil {
		present = true
	}
	result = t.IdField
	return
}
func (t *BaseElement) SetId(value *Id) {
	t.IdField = value
}
func (t *BaseElement) Documentations() (result *[]Documentation) {
	result = &t.DocumentationField
	return
}
func (t *BaseElement) SetDocumentations(value []Documentation) {
	t.DocumentationField = value
}
func (t *BaseElement) ExtensionElements() (result *ExtensionElements, present bool) {
	if t.ExtensionElementsField != nil {
		present = true
	}
	result = t.ExtensionElementsField
	return
}
func (t *BaseElement) SetExtensionElements(value *ExtensionElements) {
	t.ExtensionElementsField = value
}

type BaseElementWithMixedContent struct {
	IdField                *Id                `xml:"id,attr,omitempty"`
	DocumentationField     []Documentation    `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL documentation"`
	ExtensionElementsField *ExtensionElements `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL extensionElements"`
}

func DefaultBaseElementWithMixedContent() BaseElementWithMixedContent {
	return BaseElementWithMixedContent{}
}

type BaseElementWithMixedContentInterface interface {
	Element
	Id() (result *Id, present bool)
	Documentations() (result *[]Documentation)
	ExtensionElements() (result *ExtensionElements, present bool)
	SetId(value *Id)
	SetDocumentations(value []Documentation)
	SetExtensionElements(value *ExtensionElements)
}

func (t *BaseElementWithMixedContent) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}

	for i := range t.DocumentationField {
		if result, found = t.DocumentationField[i].FindBy(f); found {
			return
		}
	}

	if value := t.ExtensionElementsField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	return
}
func (t *BaseElementWithMixedContent) Id() (result *Id, present bool) {
	if t.IdField != nil {
		present = true
	}
	result = t.IdField
	return
}
func (t *BaseElementWithMixedContent) SetId(value *Id) {
	t.IdField = value
}
func (t *BaseElementWithMixedContent) Documentations() (result *[]Documentation) {
	result = &t.DocumentationField
	return
}
func (t *BaseElementWithMixedContent) SetDocumentations(value []Documentation) {
	t.DocumentationField = value
}
func (t *BaseElementWithMixedContent) ExtensionElements() (result *ExtensionElements, present bool) {
	if t.ExtensionElementsField != nil {
		present = true
	}
	result = t.ExtensionElementsField
	return
}
func (t *BaseElementWithMixedContent) SetExtensionElements(value *ExtensionElements) {
	t.ExtensionElementsField = value
}

type BoundaryEvent struct {
	CatchEvent
	CancelActivityField *bool    `xml:"cancelActivity,attr,omitempty"`
	AttachedToRefField  QName    `xml:"attachedToRef,attr,omitempty"`
	TextPayloadField    *Payload `xml:",chardata"`
}

var defaultBoundaryEventCancelActivityField bool = true

func DefaultBoundaryEvent() BoundaryEvent {
	return BoundaryEvent{
		CatchEvent:          DefaultCatchEvent(),
		CancelActivityField: &defaultBoundaryEventCancelActivityField,
	}
}

type BoundaryEventInterface interface {
	Element
	CatchEventInterface
	CancelActivity() (result bool)
	AttachedToRef() (result *QName)
	SetCancelActivity(value *bool)
	SetAttachedToRef(value QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *BoundaryEvent) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *BoundaryEvent) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *BoundaryEvent) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.CatchEvent.FindBy(f); found {
		return
	}

	return
}
func (t *BoundaryEvent) CancelActivity() (result bool) {
	if t.CancelActivityField == nil {
		result = defaultBoundaryEventCancelActivityField
		return
	}
	result = *t.CancelActivityField
	return
}
func (t *BoundaryEvent) SetCancelActivity(value *bool) {
	t.CancelActivityField = value
}
func (t *BoundaryEvent) AttachedToRef() (result *QName) {
	result = &t.AttachedToRefField
	return
}
func (t *BoundaryEvent) SetAttachedToRef(value QName) {
	t.AttachedToRefField = value
}

type BusinessRuleTask struct {
	Task
	ImplementationField *Implementation `xml:"implementation,attr,omitempty"`
	TextPayloadField    *Payload        `xml:",chardata"`
}

var defaultBusinessRuleTaskImplementationField Implementation = "##unspecified"

func DefaultBusinessRuleTask() BusinessRuleTask {
	return BusinessRuleTask{
		Task: DefaultTask(),
	}
}

type BusinessRuleTaskInterface interface {
	Element
	TaskInterface
	Implementation() (result *Implementation)
	SetImplementation(value *Implementation)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *BusinessRuleTask) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *BusinessRuleTask) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *BusinessRuleTask) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Task.FindBy(f); found {
		return
	}

	return
}
func (t *BusinessRuleTask) Implementation() (result *Implementation) {
	if t.ImplementationField == nil {
		result = &defaultBusinessRuleTaskImplementationField
		return
	}
	result = t.ImplementationField
	return
}
func (t *BusinessRuleTask) SetImplementation(value *Implementation) {
	t.ImplementationField = value
}

type CallableElement struct {
	RootElement
	NameField                  *string                   `xml:"name,attr,omitempty"`
	SupportedInterfaceRefField []QName                   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL supportedInterfaceRef"`
	IoSpecificationField       *InputOutputSpecification `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL ioSpecification"`
	IoBindingField             []InputOutputBinding      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL ioBinding"`
	TextPayloadField           *Payload                  `xml:",chardata"`
}

func DefaultCallableElement() CallableElement {
	return CallableElement{
		RootElement: DefaultRootElement(),
	}
}

type CallableElementInterface interface {
	Element
	RootElementInterface
	Name() (result *string, present bool)
	SupportedInterfaceRefs() (result *[]QName)
	IoSpecification() (result *InputOutputSpecification, present bool)
	IoBindings() (result *[]InputOutputBinding)
	SetName(value *string)
	SetSupportedInterfaceRefs(value []QName)
	SetIoSpecification(value *InputOutputSpecification)
	SetIoBindings(value []InputOutputBinding)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *CallableElement) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *CallableElement) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *CallableElement) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.RootElement.FindBy(f); found {
		return
	}

	if value := t.IoSpecificationField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	for i := range t.IoBindingField {
		if result, found = t.IoBindingField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *CallableElement) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *CallableElement) SetName(value *string) {
	t.NameField = value
}
func (t *CallableElement) SupportedInterfaceRefs() (result *[]QName) {
	result = &t.SupportedInterfaceRefField
	return
}
func (t *CallableElement) SetSupportedInterfaceRefs(value []QName) {
	t.SupportedInterfaceRefField = value
}
func (t *CallableElement) IoSpecification() (result *InputOutputSpecification, present bool) {
	if t.IoSpecificationField != nil {
		present = true
	}
	result = t.IoSpecificationField
	return
}
func (t *CallableElement) SetIoSpecification(value *InputOutputSpecification) {
	t.IoSpecificationField = value
}
func (t *CallableElement) IoBindings() (result *[]InputOutputBinding) {
	result = &t.IoBindingField
	return
}
func (t *CallableElement) SetIoBindings(value []InputOutputBinding) {
	t.IoBindingField = value
}

type CallActivity struct {
	Activity
	CalledElementField *QName   `xml:"calledElement,attr,omitempty"`
	TextPayloadField   *Payload `xml:",chardata"`
}

func DefaultCallActivity() CallActivity {
	return CallActivity{
		Activity: DefaultActivity(),
	}
}

type CallActivityInterface interface {
	Element
	ActivityInterface
	CalledElement() (result *QName, present bool)
	SetCalledElement(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *CallActivity) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *CallActivity) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *CallActivity) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Activity.FindBy(f); found {
		return
	}

	return
}
func (t *CallActivity) CalledElement() (result *QName, present bool) {
	if t.CalledElementField != nil {
		present = true
	}
	result = t.CalledElementField
	return
}
func (t *CallActivity) SetCalledElement(value *QName) {
	t.CalledElementField = value
}

type CallChoreography struct {
	ChoreographyActivity
	CalledChoreographyRefField  *QName                   `xml:"calledChoreographyRef,attr,omitempty"`
	ParticipantAssociationField []ParticipantAssociation `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL participantAssociation"`
	TextPayloadField            *Payload                 `xml:",chardata"`
}

func DefaultCallChoreography() CallChoreography {
	return CallChoreography{
		ChoreographyActivity: DefaultChoreographyActivity(),
	}
}

type CallChoreographyInterface interface {
	Element
	ChoreographyActivityInterface
	CalledChoreographyRef() (result *QName, present bool)
	ParticipantAssociations() (result *[]ParticipantAssociation)
	SetCalledChoreographyRef(value *QName)
	SetParticipantAssociations(value []ParticipantAssociation)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *CallChoreography) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *CallChoreography) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *CallChoreography) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.ChoreographyActivity.FindBy(f); found {
		return
	}

	for i := range t.ParticipantAssociationField {
		if result, found = t.ParticipantAssociationField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *CallChoreography) CalledChoreographyRef() (result *QName, present bool) {
	if t.CalledChoreographyRefField != nil {
		present = true
	}
	result = t.CalledChoreographyRefField
	return
}
func (t *CallChoreography) SetCalledChoreographyRef(value *QName) {
	t.CalledChoreographyRefField = value
}
func (t *CallChoreography) ParticipantAssociations() (result *[]ParticipantAssociation) {
	result = &t.ParticipantAssociationField
	return
}
func (t *CallChoreography) SetParticipantAssociations(value []ParticipantAssociation) {
	t.ParticipantAssociationField = value
}

type CallConversation struct {
	ConversationNode
	CalledCollaborationRefField *QName                   `xml:"calledCollaborationRef,attr,omitempty"`
	ParticipantAssociationField []ParticipantAssociation `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL participantAssociation"`
	TextPayloadField            *Payload                 `xml:",chardata"`
}

func DefaultCallConversation() CallConversation {
	return CallConversation{
		ConversationNode: DefaultConversationNode(),
	}
}

type CallConversationInterface interface {
	Element
	ConversationNodeInterface
	CalledCollaborationRef() (result *QName, present bool)
	ParticipantAssociations() (result *[]ParticipantAssociation)
	SetCalledCollaborationRef(value *QName)
	SetParticipantAssociations(value []ParticipantAssociation)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *CallConversation) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *CallConversation) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *CallConversation) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.ConversationNode.FindBy(f); found {
		return
	}

	for i := range t.ParticipantAssociationField {
		if result, found = t.ParticipantAssociationField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *CallConversation) CalledCollaborationRef() (result *QName, present bool) {
	if t.CalledCollaborationRefField != nil {
		present = true
	}
	result = t.CalledCollaborationRefField
	return
}
func (t *CallConversation) SetCalledCollaborationRef(value *QName) {
	t.CalledCollaborationRefField = value
}
func (t *CallConversation) ParticipantAssociations() (result *[]ParticipantAssociation) {
	result = &t.ParticipantAssociationField
	return
}
func (t *CallConversation) SetParticipantAssociations(value []ParticipantAssociation) {
	t.ParticipantAssociationField = value
}

type CancelEventDefinition struct {
	EventDefinition
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultCancelEventDefinition() CancelEventDefinition {
	return CancelEventDefinition{
		EventDefinition: DefaultEventDefinition(),
	}
}

type CancelEventDefinitionInterface interface {
	Element
	EventDefinitionInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *CancelEventDefinition) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *CancelEventDefinition) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *CancelEventDefinition) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.EventDefinition.FindBy(f); found {
		return
	}

	return
}

type CatchEvent struct {
	Event
	ParallelMultipleField           *bool                        `xml:"parallelMultiple,attr,omitempty"`
	DataOutputField                 []DataOutput                 `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataOutput"`
	DataOutputAssociationField      []DataOutputAssociation      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataOutputAssociation"`
	OutputSetField                  *OutputSet                   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL outputSet"`
	CancelEventDefinitionField      []CancelEventDefinition      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL cancelEventDefinition"`
	CompensateEventDefinitionField  []CompensateEventDefinition  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL compensateEventDefinition"`
	ConditionalEventDefinitionField []ConditionalEventDefinition `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL conditionalEventDefinition"`
	ErrorEventDefinitionField       []ErrorEventDefinition       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL errorEventDefinition"`
	EscalationEventDefinitionField  []EscalationEventDefinition  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL escalationEventDefinition"`
	LinkEventDefinitionField        []LinkEventDefinition        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL linkEventDefinition"`
	MessageEventDefinitionField     []MessageEventDefinition     `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL messageEventDefinition"`
	SignalEventDefinitionField      []SignalEventDefinition      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL signalEventDefinition"`
	TerminateEventDefinitionField   []TerminateEventDefinition   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL terminateEventDefinition"`
	TimerEventDefinitionField       []TimerEventDefinition       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL timerEventDefinition"`
	EventDefinitionRefField         []QName                      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL eventDefinitionRef"`
}

var defaultCatchEventParallelMultipleField bool = false

func DefaultCatchEvent() CatchEvent {
	return CatchEvent{
		Event:                 DefaultEvent(),
		ParallelMultipleField: &defaultCatchEventParallelMultipleField,
	}
}

type CatchEventInterface interface {
	Element
	EventInterface
	ParallelMultiple() (result bool)
	DataOutputs() (result *[]DataOutput)
	DataOutputAssociations() (result *[]DataOutputAssociation)
	OutputSet() (result *OutputSet, present bool)
	CancelEventDefinitions() (result *[]CancelEventDefinition)
	CompensateEventDefinitions() (result *[]CompensateEventDefinition)
	ConditionalEventDefinitions() (result *[]ConditionalEventDefinition)
	ErrorEventDefinitions() (result *[]ErrorEventDefinition)
	EscalationEventDefinitions() (result *[]EscalationEventDefinition)
	LinkEventDefinitions() (result *[]LinkEventDefinition)
	MessageEventDefinitions() (result *[]MessageEventDefinition)
	SignalEventDefinitions() (result *[]SignalEventDefinition)
	TerminateEventDefinitions() (result *[]TerminateEventDefinition)
	TimerEventDefinitions() (result *[]TimerEventDefinition)
	EventDefinitionRefs() (result *[]QName)
	EventDefinitions() []EventDefinitionInterface
	SetParallelMultiple(value *bool)
	SetDataOutputs(value []DataOutput)
	SetDataOutputAssociations(value []DataOutputAssociation)
	SetOutputSet(value *OutputSet)
	SetCancelEventDefinitions(value []CancelEventDefinition)
	SetCompensateEventDefinitions(value []CompensateEventDefinition)
	SetConditionalEventDefinitions(value []ConditionalEventDefinition)
	SetErrorEventDefinitions(value []ErrorEventDefinition)
	SetEscalationEventDefinitions(value []EscalationEventDefinition)
	SetLinkEventDefinitions(value []LinkEventDefinition)
	SetMessageEventDefinitions(value []MessageEventDefinition)
	SetSignalEventDefinitions(value []SignalEventDefinition)
	SetTerminateEventDefinitions(value []TerminateEventDefinition)
	SetTimerEventDefinitions(value []TimerEventDefinition)
	SetEventDefinitionRefs(value []QName)
}

func (t *CatchEvent) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Event.FindBy(f); found {
		return
	}

	for i := range t.DataOutputField {
		if result, found = t.DataOutputField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataOutputAssociationField {
		if result, found = t.DataOutputAssociationField[i].FindBy(f); found {
			return
		}
	}

	if value := t.OutputSetField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	for i := range t.CancelEventDefinitionField {
		if result, found = t.CancelEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.CompensateEventDefinitionField {
		if result, found = t.CompensateEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ConditionalEventDefinitionField {
		if result, found = t.ConditionalEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ErrorEventDefinitionField {
		if result, found = t.ErrorEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.EscalationEventDefinitionField {
		if result, found = t.EscalationEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.LinkEventDefinitionField {
		if result, found = t.LinkEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.MessageEventDefinitionField {
		if result, found = t.MessageEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SignalEventDefinitionField {
		if result, found = t.SignalEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.TerminateEventDefinitionField {
		if result, found = t.TerminateEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.TimerEventDefinitionField {
		if result, found = t.TimerEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *CatchEvent) ParallelMultiple() (result bool) {
	if t.ParallelMultipleField == nil {
		result = defaultCatchEventParallelMultipleField
		return
	}
	result = *t.ParallelMultipleField
	return
}
func (t *CatchEvent) SetParallelMultiple(value *bool) {
	t.ParallelMultipleField = value
}
func (t *CatchEvent) EventDefinitions() []EventDefinitionInterface {

	result := make([]EventDefinitionInterface, 0)

	for i := range t.CancelEventDefinitionField {
		result = append(result, &t.CancelEventDefinitionField[i])
	}

	for i := range t.CompensateEventDefinitionField {
		result = append(result, &t.CompensateEventDefinitionField[i])
	}

	for i := range t.ConditionalEventDefinitionField {
		result = append(result, &t.ConditionalEventDefinitionField[i])
	}

	for i := range t.ErrorEventDefinitionField {
		result = append(result, &t.ErrorEventDefinitionField[i])
	}

	for i := range t.EscalationEventDefinitionField {
		result = append(result, &t.EscalationEventDefinitionField[i])
	}

	for i := range t.LinkEventDefinitionField {
		result = append(result, &t.LinkEventDefinitionField[i])
	}

	for i := range t.MessageEventDefinitionField {
		result = append(result, &t.MessageEventDefinitionField[i])
	}

	for i := range t.SignalEventDefinitionField {
		result = append(result, &t.SignalEventDefinitionField[i])
	}

	for i := range t.TerminateEventDefinitionField {
		result = append(result, &t.TerminateEventDefinitionField[i])
	}

	for i := range t.TimerEventDefinitionField {
		result = append(result, &t.TimerEventDefinitionField[i])
	}
	return result
}
func (t *CatchEvent) DataOutputs() (result *[]DataOutput) {
	result = &t.DataOutputField
	return
}
func (t *CatchEvent) SetDataOutputs(value []DataOutput) {
	t.DataOutputField = value
}
func (t *CatchEvent) DataOutputAssociations() (result *[]DataOutputAssociation) {
	result = &t.DataOutputAssociationField
	return
}
func (t *CatchEvent) SetDataOutputAssociations(value []DataOutputAssociation) {
	t.DataOutputAssociationField = value
}
func (t *CatchEvent) OutputSet() (result *OutputSet, present bool) {
	if t.OutputSetField != nil {
		present = true
	}
	result = t.OutputSetField
	return
}
func (t *CatchEvent) SetOutputSet(value *OutputSet) {
	t.OutputSetField = value
}
func (t *CatchEvent) CancelEventDefinitions() (result *[]CancelEventDefinition) {
	result = &t.CancelEventDefinitionField
	return
}
func (t *CatchEvent) SetCancelEventDefinitions(value []CancelEventDefinition) {
	t.CancelEventDefinitionField = value
}
func (t *CatchEvent) CompensateEventDefinitions() (result *[]CompensateEventDefinition) {
	result = &t.CompensateEventDefinitionField
	return
}
func (t *CatchEvent) SetCompensateEventDefinitions(value []CompensateEventDefinition) {
	t.CompensateEventDefinitionField = value
}
func (t *CatchEvent) ConditionalEventDefinitions() (result *[]ConditionalEventDefinition) {
	result = &t.ConditionalEventDefinitionField
	return
}
func (t *CatchEvent) SetConditionalEventDefinitions(value []ConditionalEventDefinition) {
	t.ConditionalEventDefinitionField = value
}
func (t *CatchEvent) ErrorEventDefinitions() (result *[]ErrorEventDefinition) {
	result = &t.ErrorEventDefinitionField
	return
}
func (t *CatchEvent) SetErrorEventDefinitions(value []ErrorEventDefinition) {
	t.ErrorEventDefinitionField = value
}
func (t *CatchEvent) EscalationEventDefinitions() (result *[]EscalationEventDefinition) {
	result = &t.EscalationEventDefinitionField
	return
}
func (t *CatchEvent) SetEscalationEventDefinitions(value []EscalationEventDefinition) {
	t.EscalationEventDefinitionField = value
}
func (t *CatchEvent) LinkEventDefinitions() (result *[]LinkEventDefinition) {
	result = &t.LinkEventDefinitionField
	return
}
func (t *CatchEvent) SetLinkEventDefinitions(value []LinkEventDefinition) {
	t.LinkEventDefinitionField = value
}
func (t *CatchEvent) MessageEventDefinitions() (result *[]MessageEventDefinition) {
	result = &t.MessageEventDefinitionField
	return
}
func (t *CatchEvent) SetMessageEventDefinitions(value []MessageEventDefinition) {
	t.MessageEventDefinitionField = value
}
func (t *CatchEvent) SignalEventDefinitions() (result *[]SignalEventDefinition) {
	result = &t.SignalEventDefinitionField
	return
}
func (t *CatchEvent) SetSignalEventDefinitions(value []SignalEventDefinition) {
	t.SignalEventDefinitionField = value
}
func (t *CatchEvent) TerminateEventDefinitions() (result *[]TerminateEventDefinition) {
	result = &t.TerminateEventDefinitionField
	return
}
func (t *CatchEvent) SetTerminateEventDefinitions(value []TerminateEventDefinition) {
	t.TerminateEventDefinitionField = value
}
func (t *CatchEvent) TimerEventDefinitions() (result *[]TimerEventDefinition) {
	result = &t.TimerEventDefinitionField
	return
}
func (t *CatchEvent) SetTimerEventDefinitions(value []TimerEventDefinition) {
	t.TimerEventDefinitionField = value
}
func (t *CatchEvent) EventDefinitionRefs() (result *[]QName) {
	result = &t.EventDefinitionRefField
	return
}
func (t *CatchEvent) SetEventDefinitionRefs(value []QName) {
	t.EventDefinitionRefField = value
}

type Category struct {
	RootElement
	NameField          *string         `xml:"name,attr,omitempty"`
	CategoryValueField []CategoryValue `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL categoryValue"`
	TextPayloadField   *Payload        `xml:",chardata"`
}

func DefaultCategory() Category {
	return Category{
		RootElement: DefaultRootElement(),
	}
}

type CategoryInterface interface {
	Element
	RootElementInterface
	Name() (result *string, present bool)
	CategoryValues() (result *[]CategoryValue)
	SetName(value *string)
	SetCategoryValues(value []CategoryValue)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Category) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Category) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Category) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.RootElement.FindBy(f); found {
		return
	}

	for i := range t.CategoryValueField {
		if result, found = t.CategoryValueField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *Category) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *Category) SetName(value *string) {
	t.NameField = value
}
func (t *Category) CategoryValues() (result *[]CategoryValue) {
	result = &t.CategoryValueField
	return
}
func (t *Category) SetCategoryValues(value []CategoryValue) {
	t.CategoryValueField = value
}

type CategoryValue struct {
	BaseElement
	ValueField       *string  `xml:"value,attr,omitempty"`
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultCategoryValue() CategoryValue {
	return CategoryValue{
		BaseElement: DefaultBaseElement(),
	}
}

type CategoryValueInterface interface {
	Element
	BaseElementInterface
	Value() (result *string, present bool)
	SetValue(value *string)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *CategoryValue) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *CategoryValue) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *CategoryValue) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}
func (t *CategoryValue) Value() (result *string, present bool) {
	if t.ValueField != nil {
		present = true
	}
	result = t.ValueField
	return
}
func (t *CategoryValue) SetValue(value *string) {
	t.ValueField = value
}

type Choreography struct {
	Collaboration
	AdHocSubProcessField        []AdHocSubProcess        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL adHocSubProcess"`
	BoundaryEventField          []BoundaryEvent          `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL boundaryEvent"`
	BusinessRuleTaskField       []BusinessRuleTask       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL businessRuleTask"`
	CallActivityField           []CallActivity           `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL callActivity"`
	CallChoreographyField       []CallChoreography       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL callChoreography"`
	ChoreographyTaskField       []ChoreographyTask       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL choreographyTask"`
	ComplexGatewayField         []ComplexGateway         `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL complexGateway"`
	DataObjectField             []DataObject             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataObject"`
	DataObjectReferenceField    []DataObjectReference    `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataObjectReference"`
	DataStoreReferenceField     []DataStoreReference     `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataStoreReference"`
	EndEventField               []EndEvent               `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL endEvent"`
	EventField                  []Event                  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL event"`
	EventBasedGatewayField      []EventBasedGateway      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL eventBasedGateway"`
	ExclusiveGatewayField       []ExclusiveGateway       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL exclusiveGateway"`
	ImplicitThrowEventField     []ImplicitThrowEvent     `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL implicitThrowEvent"`
	InclusiveGatewayField       []InclusiveGateway       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL inclusiveGateway"`
	IntermediateCatchEventField []IntermediateCatchEvent `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL intermediateCatchEvent"`
	IntermediateThrowEventField []IntermediateThrowEvent `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL intermediateThrowEvent"`
	ManualTaskField             []ManualTask             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL manualTask"`
	ParallelGatewayField        []ParallelGateway        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL parallelGateway"`
	ReceiveTaskField            []ReceiveTask            `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL receiveTask"`
	ScriptTaskField             []ScriptTask             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL scriptTask"`
	SendTaskField               []SendTask               `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL sendTask"`
	SequenceFlowField           []SequenceFlow           `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL sequenceFlow"`
	ServiceTaskField            []ServiceTask            `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL serviceTask"`
	StartEventField             []StartEvent             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL startEvent"`
	SubChoreographyField        []SubChoreography        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL subChoreography"`
	SubProcessField             []SubProcess             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL subProcess"`
	TaskField                   []Task                   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL task"`
	TransactionField            []Transaction            `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL transaction"`
	UserTaskField               []UserTask               `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL userTask"`
	TextPayloadField            *Payload                 `xml:",chardata"`
}

func DefaultChoreography() Choreography {
	return Choreography{
		Collaboration: DefaultCollaboration(),
	}
}

type ChoreographyInterface interface {
	Element
	CollaborationInterface
	AdHocSubProcesses() (result *[]AdHocSubProcess)
	BoundaryEvents() (result *[]BoundaryEvent)
	BusinessRuleTasks() (result *[]BusinessRuleTask)
	CallActivities() (result *[]CallActivity)
	CallChoreographies() (result *[]CallChoreography)
	ChoreographyTasks() (result *[]ChoreographyTask)
	ComplexGateways() (result *[]ComplexGateway)
	DataObjects() (result *[]DataObject)
	DataObjectReferences() (result *[]DataObjectReference)
	DataStoreReferences() (result *[]DataStoreReference)
	EndEvents() (result *[]EndEvent)
	Events() (result *[]Event)
	EventBasedGateways() (result *[]EventBasedGateway)
	ExclusiveGateways() (result *[]ExclusiveGateway)
	ImplicitThrowEvents() (result *[]ImplicitThrowEvent)
	InclusiveGateways() (result *[]InclusiveGateway)
	IntermediateCatchEvents() (result *[]IntermediateCatchEvent)
	IntermediateThrowEvents() (result *[]IntermediateThrowEvent)
	ManualTasks() (result *[]ManualTask)
	ParallelGateways() (result *[]ParallelGateway)
	ReceiveTasks() (result *[]ReceiveTask)
	ScriptTasks() (result *[]ScriptTask)
	SendTasks() (result *[]SendTask)
	SequenceFlows() (result *[]SequenceFlow)
	ServiceTasks() (result *[]ServiceTask)
	StartEvents() (result *[]StartEvent)
	SubChoreographies() (result *[]SubChoreography)
	SubProcesses() (result *[]SubProcess)
	Tasks() (result *[]Task)
	Transactions() (result *[]Transaction)
	UserTasks() (result *[]UserTask)
	FlowElements() []FlowElementInterface
	SetAdHocSubProcesses(value []AdHocSubProcess)
	SetBoundaryEvents(value []BoundaryEvent)
	SetBusinessRuleTasks(value []BusinessRuleTask)
	SetCallActivities(value []CallActivity)
	SetCallChoreographies(value []CallChoreography)
	SetChoreographyTasks(value []ChoreographyTask)
	SetComplexGateways(value []ComplexGateway)
	SetDataObjects(value []DataObject)
	SetDataObjectReferences(value []DataObjectReference)
	SetDataStoreReferences(value []DataStoreReference)
	SetEndEvents(value []EndEvent)
	SetEvents(value []Event)
	SetEventBasedGateways(value []EventBasedGateway)
	SetExclusiveGateways(value []ExclusiveGateway)
	SetImplicitThrowEvents(value []ImplicitThrowEvent)
	SetInclusiveGateways(value []InclusiveGateway)
	SetIntermediateCatchEvents(value []IntermediateCatchEvent)
	SetIntermediateThrowEvents(value []IntermediateThrowEvent)
	SetManualTasks(value []ManualTask)
	SetParallelGateways(value []ParallelGateway)
	SetReceiveTasks(value []ReceiveTask)
	SetScriptTasks(value []ScriptTask)
	SetSendTasks(value []SendTask)
	SetSequenceFlows(value []SequenceFlow)
	SetServiceTasks(value []ServiceTask)
	SetStartEvents(value []StartEvent)
	SetSubChoreographies(value []SubChoreography)
	SetSubProcesses(value []SubProcess)
	SetTasks(value []Task)
	SetTransactions(value []Transaction)
	SetUserTasks(value []UserTask)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Choreography) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Choreography) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Choreography) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Collaboration.FindBy(f); found {
		return
	}

	for i := range t.AdHocSubProcessField {
		if result, found = t.AdHocSubProcessField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.BoundaryEventField {
		if result, found = t.BoundaryEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.BusinessRuleTaskField {
		if result, found = t.BusinessRuleTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.CallActivityField {
		if result, found = t.CallActivityField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.CallChoreographyField {
		if result, found = t.CallChoreographyField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ChoreographyTaskField {
		if result, found = t.ChoreographyTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ComplexGatewayField {
		if result, found = t.ComplexGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataObjectField {
		if result, found = t.DataObjectField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataObjectReferenceField {
		if result, found = t.DataObjectReferenceField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataStoreReferenceField {
		if result, found = t.DataStoreReferenceField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.EndEventField {
		if result, found = t.EndEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.EventField {
		if result, found = t.EventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.EventBasedGatewayField {
		if result, found = t.EventBasedGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ExclusiveGatewayField {
		if result, found = t.ExclusiveGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ImplicitThrowEventField {
		if result, found = t.ImplicitThrowEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.InclusiveGatewayField {
		if result, found = t.InclusiveGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.IntermediateCatchEventField {
		if result, found = t.IntermediateCatchEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.IntermediateThrowEventField {
		if result, found = t.IntermediateThrowEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ManualTaskField {
		if result, found = t.ManualTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ParallelGatewayField {
		if result, found = t.ParallelGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ReceiveTaskField {
		if result, found = t.ReceiveTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ScriptTaskField {
		if result, found = t.ScriptTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SendTaskField {
		if result, found = t.SendTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SequenceFlowField {
		if result, found = t.SequenceFlowField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ServiceTaskField {
		if result, found = t.ServiceTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.StartEventField {
		if result, found = t.StartEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SubChoreographyField {
		if result, found = t.SubChoreographyField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SubProcessField {
		if result, found = t.SubProcessField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.TaskField {
		if result, found = t.TaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.TransactionField {
		if result, found = t.TransactionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.UserTaskField {
		if result, found = t.UserTaskField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *Choreography) FlowElements() []FlowElementInterface {

	result := make([]FlowElementInterface, 0)

	for i := range t.AdHocSubProcessField {
		result = append(result, &t.AdHocSubProcessField[i])
	}

	for i := range t.BoundaryEventField {
		result = append(result, &t.BoundaryEventField[i])
	}

	for i := range t.BusinessRuleTaskField {
		result = append(result, &t.BusinessRuleTaskField[i])
	}

	for i := range t.CallActivityField {
		result = append(result, &t.CallActivityField[i])
	}

	for i := range t.CallChoreographyField {
		result = append(result, &t.CallChoreographyField[i])
	}

	for i := range t.ChoreographyTaskField {
		result = append(result, &t.ChoreographyTaskField[i])
	}

	for i := range t.ComplexGatewayField {
		result = append(result, &t.ComplexGatewayField[i])
	}

	for i := range t.DataObjectField {
		result = append(result, &t.DataObjectField[i])
	}

	for i := range t.DataObjectReferenceField {
		result = append(result, &t.DataObjectReferenceField[i])
	}

	for i := range t.DataStoreReferenceField {
		result = append(result, &t.DataStoreReferenceField[i])
	}

	for i := range t.EndEventField {
		result = append(result, &t.EndEventField[i])
	}

	for i := range t.EventField {
		result = append(result, &t.EventField[i])
	}

	for i := range t.EventBasedGatewayField {
		result = append(result, &t.EventBasedGatewayField[i])
	}

	for i := range t.ExclusiveGatewayField {
		result = append(result, &t.ExclusiveGatewayField[i])
	}

	for i := range t.ImplicitThrowEventField {
		result = append(result, &t.ImplicitThrowEventField[i])
	}

	for i := range t.InclusiveGatewayField {
		result = append(result, &t.InclusiveGatewayField[i])
	}

	for i := range t.IntermediateCatchEventField {
		result = append(result, &t.IntermediateCatchEventField[i])
	}

	for i := range t.IntermediateThrowEventField {
		result = append(result, &t.IntermediateThrowEventField[i])
	}

	for i := range t.ManualTaskField {
		result = append(result, &t.ManualTaskField[i])
	}

	for i := range t.ParallelGatewayField {
		result = append(result, &t.ParallelGatewayField[i])
	}

	for i := range t.ReceiveTaskField {
		result = append(result, &t.ReceiveTaskField[i])
	}

	for i := range t.ScriptTaskField {
		result = append(result, &t.ScriptTaskField[i])
	}

	for i := range t.SendTaskField {
		result = append(result, &t.SendTaskField[i])
	}

	for i := range t.SequenceFlowField {
		result = append(result, &t.SequenceFlowField[i])
	}

	for i := range t.ServiceTaskField {
		result = append(result, &t.ServiceTaskField[i])
	}

	for i := range t.StartEventField {
		result = append(result, &t.StartEventField[i])
	}

	for i := range t.SubChoreographyField {
		result = append(result, &t.SubChoreographyField[i])
	}

	for i := range t.SubProcessField {
		result = append(result, &t.SubProcessField[i])
	}

	for i := range t.TaskField {
		result = append(result, &t.TaskField[i])
	}

	for i := range t.TransactionField {
		result = append(result, &t.TransactionField[i])
	}

	for i := range t.UserTaskField {
		result = append(result, &t.UserTaskField[i])
	}
	return result
}
func (t *Choreography) AdHocSubProcesses() (result *[]AdHocSubProcess) {
	result = &t.AdHocSubProcessField
	return
}
func (t *Choreography) SetAdHocSubProcesses(value []AdHocSubProcess) {
	t.AdHocSubProcessField = value
}
func (t *Choreography) BoundaryEvents() (result *[]BoundaryEvent) {
	result = &t.BoundaryEventField
	return
}
func (t *Choreography) SetBoundaryEvents(value []BoundaryEvent) {
	t.BoundaryEventField = value
}
func (t *Choreography) BusinessRuleTasks() (result *[]BusinessRuleTask) {
	result = &t.BusinessRuleTaskField
	return
}
func (t *Choreography) SetBusinessRuleTasks(value []BusinessRuleTask) {
	t.BusinessRuleTaskField = value
}
func (t *Choreography) CallActivities() (result *[]CallActivity) {
	result = &t.CallActivityField
	return
}
func (t *Choreography) SetCallActivities(value []CallActivity) {
	t.CallActivityField = value
}
func (t *Choreography) CallChoreographies() (result *[]CallChoreography) {
	result = &t.CallChoreographyField
	return
}
func (t *Choreography) SetCallChoreographies(value []CallChoreography) {
	t.CallChoreographyField = value
}
func (t *Choreography) ChoreographyTasks() (result *[]ChoreographyTask) {
	result = &t.ChoreographyTaskField
	return
}
func (t *Choreography) SetChoreographyTasks(value []ChoreographyTask) {
	t.ChoreographyTaskField = value
}
func (t *Choreography) ComplexGateways() (result *[]ComplexGateway) {
	result = &t.ComplexGatewayField
	return
}
func (t *Choreography) SetComplexGateways(value []ComplexGateway) {
	t.ComplexGatewayField = value
}
func (t *Choreography) DataObjects() (result *[]DataObject) {
	result = &t.DataObjectField
	return
}
func (t *Choreography) SetDataObjects(value []DataObject) {
	t.DataObjectField = value
}
func (t *Choreography) DataObjectReferences() (result *[]DataObjectReference) {
	result = &t.DataObjectReferenceField
	return
}
func (t *Choreography) SetDataObjectReferences(value []DataObjectReference) {
	t.DataObjectReferenceField = value
}
func (t *Choreography) DataStoreReferences() (result *[]DataStoreReference) {
	result = &t.DataStoreReferenceField
	return
}
func (t *Choreography) SetDataStoreReferences(value []DataStoreReference) {
	t.DataStoreReferenceField = value
}
func (t *Choreography) EndEvents() (result *[]EndEvent) {
	result = &t.EndEventField
	return
}
func (t *Choreography) SetEndEvents(value []EndEvent) {
	t.EndEventField = value
}
func (t *Choreography) Events() (result *[]Event) {
	result = &t.EventField
	return
}
func (t *Choreography) SetEvents(value []Event) {
	t.EventField = value
}
func (t *Choreography) EventBasedGateways() (result *[]EventBasedGateway) {
	result = &t.EventBasedGatewayField
	return
}
func (t *Choreography) SetEventBasedGateways(value []EventBasedGateway) {
	t.EventBasedGatewayField = value
}
func (t *Choreography) ExclusiveGateways() (result *[]ExclusiveGateway) {
	result = &t.ExclusiveGatewayField
	return
}
func (t *Choreography) SetExclusiveGateways(value []ExclusiveGateway) {
	t.ExclusiveGatewayField = value
}
func (t *Choreography) ImplicitThrowEvents() (result *[]ImplicitThrowEvent) {
	result = &t.ImplicitThrowEventField
	return
}
func (t *Choreography) SetImplicitThrowEvents(value []ImplicitThrowEvent) {
	t.ImplicitThrowEventField = value
}
func (t *Choreography) InclusiveGateways() (result *[]InclusiveGateway) {
	result = &t.InclusiveGatewayField
	return
}
func (t *Choreography) SetInclusiveGateways(value []InclusiveGateway) {
	t.InclusiveGatewayField = value
}
func (t *Choreography) IntermediateCatchEvents() (result *[]IntermediateCatchEvent) {
	result = &t.IntermediateCatchEventField
	return
}
func (t *Choreography) SetIntermediateCatchEvents(value []IntermediateCatchEvent) {
	t.IntermediateCatchEventField = value
}
func (t *Choreography) IntermediateThrowEvents() (result *[]IntermediateThrowEvent) {
	result = &t.IntermediateThrowEventField
	return
}
func (t *Choreography) SetIntermediateThrowEvents(value []IntermediateThrowEvent) {
	t.IntermediateThrowEventField = value
}
func (t *Choreography) ManualTasks() (result *[]ManualTask) {
	result = &t.ManualTaskField
	return
}
func (t *Choreography) SetManualTasks(value []ManualTask) {
	t.ManualTaskField = value
}
func (t *Choreography) ParallelGateways() (result *[]ParallelGateway) {
	result = &t.ParallelGatewayField
	return
}
func (t *Choreography) SetParallelGateways(value []ParallelGateway) {
	t.ParallelGatewayField = value
}
func (t *Choreography) ReceiveTasks() (result *[]ReceiveTask) {
	result = &t.ReceiveTaskField
	return
}
func (t *Choreography) SetReceiveTasks(value []ReceiveTask) {
	t.ReceiveTaskField = value
}
func (t *Choreography) ScriptTasks() (result *[]ScriptTask) {
	result = &t.ScriptTaskField
	return
}
func (t *Choreography) SetScriptTasks(value []ScriptTask) {
	t.ScriptTaskField = value
}
func (t *Choreography) SendTasks() (result *[]SendTask) {
	result = &t.SendTaskField
	return
}
func (t *Choreography) SetSendTasks(value []SendTask) {
	t.SendTaskField = value
}
func (t *Choreography) SequenceFlows() (result *[]SequenceFlow) {
	result = &t.SequenceFlowField
	return
}
func (t *Choreography) SetSequenceFlows(value []SequenceFlow) {
	t.SequenceFlowField = value
}
func (t *Choreography) ServiceTasks() (result *[]ServiceTask) {
	result = &t.ServiceTaskField
	return
}
func (t *Choreography) SetServiceTasks(value []ServiceTask) {
	t.ServiceTaskField = value
}
func (t *Choreography) StartEvents() (result *[]StartEvent) {
	result = &t.StartEventField
	return
}
func (t *Choreography) SetStartEvents(value []StartEvent) {
	t.StartEventField = value
}
func (t *Choreography) SubChoreographies() (result *[]SubChoreography) {
	result = &t.SubChoreographyField
	return
}
func (t *Choreography) SetSubChoreographies(value []SubChoreography) {
	t.SubChoreographyField = value
}
func (t *Choreography) SubProcesses() (result *[]SubProcess) {
	result = &t.SubProcessField
	return
}
func (t *Choreography) SetSubProcesses(value []SubProcess) {
	t.SubProcessField = value
}
func (t *Choreography) Tasks() (result *[]Task) {
	result = &t.TaskField
	return
}
func (t *Choreography) SetTasks(value []Task) {
	t.TaskField = value
}
func (t *Choreography) Transactions() (result *[]Transaction) {
	result = &t.TransactionField
	return
}
func (t *Choreography) SetTransactions(value []Transaction) {
	t.TransactionField = value
}
func (t *Choreography) UserTasks() (result *[]UserTask) {
	result = &t.UserTaskField
	return
}
func (t *Choreography) SetUserTasks(value []UserTask) {
	t.UserTaskField = value
}

type ChoreographyActivity struct {
	FlowNode
	InitiatingParticipantRefField QName                 `xml:"initiatingParticipantRef,attr,omitempty"`
	LoopTypeField                 *ChoreographyLoopType `xml:"loopType,attr,omitempty"`
	ParticipantRefField           []QName               `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL participantRef"`
	CorrelationKeyField           []CorrelationKey      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL correlationKey"`
}

var defaultChoreographyActivityLoopTypeField ChoreographyLoopType = "None"

func DefaultChoreographyActivity() ChoreographyActivity {
	return ChoreographyActivity{
		FlowNode:      DefaultFlowNode(),
		LoopTypeField: &defaultChoreographyActivityLoopTypeField,
	}
}

type ChoreographyActivityInterface interface {
	Element
	FlowNodeInterface
	InitiatingParticipantRef() (result *QName)
	LoopType() (result *ChoreographyLoopType)
	ParticipantRefs() (result *[]QName)
	CorrelationKeys() (result *[]CorrelationKey)
	SetInitiatingParticipantRef(value QName)
	SetLoopType(value *ChoreographyLoopType)
	SetParticipantRefs(value []QName)
	SetCorrelationKeys(value []CorrelationKey)
}

func (t *ChoreographyActivity) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.FlowNode.FindBy(f); found {
		return
	}

	for i := range t.CorrelationKeyField {
		if result, found = t.CorrelationKeyField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *ChoreographyActivity) InitiatingParticipantRef() (result *QName) {
	result = &t.InitiatingParticipantRefField
	return
}
func (t *ChoreographyActivity) SetInitiatingParticipantRef(value QName) {
	t.InitiatingParticipantRefField = value
}
func (t *ChoreographyActivity) LoopType() (result *ChoreographyLoopType) {
	if t.LoopTypeField == nil {
		result = &defaultChoreographyActivityLoopTypeField
		return
	}
	result = t.LoopTypeField
	return
}
func (t *ChoreographyActivity) SetLoopType(value *ChoreographyLoopType) {
	t.LoopTypeField = value
}
func (t *ChoreographyActivity) ParticipantRefs() (result *[]QName) {
	result = &t.ParticipantRefField
	return
}
func (t *ChoreographyActivity) SetParticipantRefs(value []QName) {
	t.ParticipantRefField = value
}
func (t *ChoreographyActivity) CorrelationKeys() (result *[]CorrelationKey) {
	result = &t.CorrelationKeyField
	return
}
func (t *ChoreographyActivity) SetCorrelationKeys(value []CorrelationKey) {
	t.CorrelationKeyField = value
}

type ChoreographyTask struct {
	ChoreographyActivity
	MessageFlowRefField QName    `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL messageFlowRef"`
	TextPayloadField    *Payload `xml:",chardata"`
}

func DefaultChoreographyTask() ChoreographyTask {
	return ChoreographyTask{
		ChoreographyActivity: DefaultChoreographyActivity(),
	}
}

type ChoreographyTaskInterface interface {
	Element
	ChoreographyActivityInterface
	MessageFlowRef() (result *QName)
	SetMessageFlowRef(value QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ChoreographyTask) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ChoreographyTask) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ChoreographyTask) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.ChoreographyActivity.FindBy(f); found {
		return
	}

	return
}
func (t *ChoreographyTask) MessageFlowRef() (result *QName) {
	result = &t.MessageFlowRefField
	return
}
func (t *ChoreographyTask) SetMessageFlowRef(value QName) {
	t.MessageFlowRefField = value
}

type Collaboration struct {
	RootElement
	NameField                    *string                   `xml:"name,attr,omitempty"`
	IsClosedField                *bool                     `xml:"isClosed,attr,omitempty"`
	ParticipantField             []Participant             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL participant"`
	MessageFlowField             []MessageFlow             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL messageFlow"`
	AssociationField             []Association             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL association"`
	GroupField                   []Group                   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL group"`
	TextAnnotationField          []TextAnnotation          `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL textAnnotation"`
	CallConversationField        []CallConversation        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL callConversation"`
	ConversationField            []Conversation            `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL conversation"`
	SubConversationField         []SubConversation         `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL subConversation"`
	ConversationAssociationField []ConversationAssociation `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL conversationAssociation"`
	ParticipantAssociationField  []ParticipantAssociation  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL participantAssociation"`
	MessageFlowAssociationField  []MessageFlowAssociation  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL messageFlowAssociation"`
	CorrelationKeyField          []CorrelationKey          `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL correlationKey"`
	ChoreographyRefField         []QName                   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL choreographyRef"`
	ConversationLinkField        []ConversationLink        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL conversationLink"`
	TextPayloadField             *Payload                  `xml:",chardata"`
}

var defaultCollaborationIsClosedField bool = false

func DefaultCollaboration() Collaboration {
	return Collaboration{
		RootElement:   DefaultRootElement(),
		IsClosedField: &defaultCollaborationIsClosedField,
	}
}

type CollaborationInterface interface {
	Element
	RootElementInterface
	Name() (result *string, present bool)
	IsClosed() (result bool)
	Participants() (result *[]Participant)
	MessageFlows() (result *[]MessageFlow)
	Associations() (result *[]Association)
	Groups() (result *[]Group)
	TextAnnotations() (result *[]TextAnnotation)
	CallConversations() (result *[]CallConversation)
	Conversations() (result *[]Conversation)
	SubConversations() (result *[]SubConversation)
	ConversationAssociations() (result *[]ConversationAssociation)
	ParticipantAssociations() (result *[]ParticipantAssociation)
	MessageFlowAssociations() (result *[]MessageFlowAssociation)
	CorrelationKeys() (result *[]CorrelationKey)
	ChoreographyRefs() (result *[]QName)
	ConversationLinks() (result *[]ConversationLink)
	Artifacts() []ArtifactInterface
	ConversationNodes() []ConversationNodeInterface
	SetName(value *string)
	SetIsClosed(value *bool)
	SetParticipants(value []Participant)
	SetMessageFlows(value []MessageFlow)
	SetAssociations(value []Association)
	SetGroups(value []Group)
	SetTextAnnotations(value []TextAnnotation)
	SetCallConversations(value []CallConversation)
	SetConversations(value []Conversation)
	SetSubConversations(value []SubConversation)
	SetConversationAssociations(value []ConversationAssociation)
	SetParticipantAssociations(value []ParticipantAssociation)
	SetMessageFlowAssociations(value []MessageFlowAssociation)
	SetCorrelationKeys(value []CorrelationKey)
	SetChoreographyRefs(value []QName)
	SetConversationLinks(value []ConversationLink)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Collaboration) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Collaboration) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Collaboration) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.RootElement.FindBy(f); found {
		return
	}

	for i := range t.ParticipantField {
		if result, found = t.ParticipantField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.MessageFlowField {
		if result, found = t.MessageFlowField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.AssociationField {
		if result, found = t.AssociationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.GroupField {
		if result, found = t.GroupField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.TextAnnotationField {
		if result, found = t.TextAnnotationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.CallConversationField {
		if result, found = t.CallConversationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ConversationField {
		if result, found = t.ConversationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SubConversationField {
		if result, found = t.SubConversationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ConversationAssociationField {
		if result, found = t.ConversationAssociationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ParticipantAssociationField {
		if result, found = t.ParticipantAssociationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.MessageFlowAssociationField {
		if result, found = t.MessageFlowAssociationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.CorrelationKeyField {
		if result, found = t.CorrelationKeyField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ConversationLinkField {
		if result, found = t.ConversationLinkField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *Collaboration) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *Collaboration) SetName(value *string) {
	t.NameField = value
}
func (t *Collaboration) IsClosed() (result bool) {
	if t.IsClosedField == nil {
		result = defaultCollaborationIsClosedField
		return
	}
	result = *t.IsClosedField
	return
}
func (t *Collaboration) SetIsClosed(value *bool) {
	t.IsClosedField = value
}
func (t *Collaboration) Artifacts() []ArtifactInterface {

	result := make([]ArtifactInterface, 0)

	for i := range t.AssociationField {
		result = append(result, &t.AssociationField[i])
	}

	for i := range t.GroupField {
		result = append(result, &t.GroupField[i])
	}

	for i := range t.TextAnnotationField {
		result = append(result, &t.TextAnnotationField[i])
	}
	return result
}
func (t *Collaboration) ConversationNodes() []ConversationNodeInterface {

	result := make([]ConversationNodeInterface, 0)

	for i := range t.CallConversationField {
		result = append(result, &t.CallConversationField[i])
	}

	for i := range t.ConversationField {
		result = append(result, &t.ConversationField[i])
	}

	for i := range t.SubConversationField {
		result = append(result, &t.SubConversationField[i])
	}
	return result
}
func (t *Collaboration) Participants() (result *[]Participant) {
	result = &t.ParticipantField
	return
}
func (t *Collaboration) SetParticipants(value []Participant) {
	t.ParticipantField = value
}
func (t *Collaboration) MessageFlows() (result *[]MessageFlow) {
	result = &t.MessageFlowField
	return
}
func (t *Collaboration) SetMessageFlows(value []MessageFlow) {
	t.MessageFlowField = value
}
func (t *Collaboration) Associations() (result *[]Association) {
	result = &t.AssociationField
	return
}
func (t *Collaboration) SetAssociations(value []Association) {
	t.AssociationField = value
}
func (t *Collaboration) Groups() (result *[]Group) {
	result = &t.GroupField
	return
}
func (t *Collaboration) SetGroups(value []Group) {
	t.GroupField = value
}
func (t *Collaboration) TextAnnotations() (result *[]TextAnnotation) {
	result = &t.TextAnnotationField
	return
}
func (t *Collaboration) SetTextAnnotations(value []TextAnnotation) {
	t.TextAnnotationField = value
}
func (t *Collaboration) CallConversations() (result *[]CallConversation) {
	result = &t.CallConversationField
	return
}
func (t *Collaboration) SetCallConversations(value []CallConversation) {
	t.CallConversationField = value
}
func (t *Collaboration) Conversations() (result *[]Conversation) {
	result = &t.ConversationField
	return
}
func (t *Collaboration) SetConversations(value []Conversation) {
	t.ConversationField = value
}
func (t *Collaboration) SubConversations() (result *[]SubConversation) {
	result = &t.SubConversationField
	return
}
func (t *Collaboration) SetSubConversations(value []SubConversation) {
	t.SubConversationField = value
}
func (t *Collaboration) ConversationAssociations() (result *[]ConversationAssociation) {
	result = &t.ConversationAssociationField
	return
}
func (t *Collaboration) SetConversationAssociations(value []ConversationAssociation) {
	t.ConversationAssociationField = value
}
func (t *Collaboration) ParticipantAssociations() (result *[]ParticipantAssociation) {
	result = &t.ParticipantAssociationField
	return
}
func (t *Collaboration) SetParticipantAssociations(value []ParticipantAssociation) {
	t.ParticipantAssociationField = value
}
func (t *Collaboration) MessageFlowAssociations() (result *[]MessageFlowAssociation) {
	result = &t.MessageFlowAssociationField
	return
}
func (t *Collaboration) SetMessageFlowAssociations(value []MessageFlowAssociation) {
	t.MessageFlowAssociationField = value
}
func (t *Collaboration) CorrelationKeys() (result *[]CorrelationKey) {
	result = &t.CorrelationKeyField
	return
}
func (t *Collaboration) SetCorrelationKeys(value []CorrelationKey) {
	t.CorrelationKeyField = value
}
func (t *Collaboration) ChoreographyRefs() (result *[]QName) {
	result = &t.ChoreographyRefField
	return
}
func (t *Collaboration) SetChoreographyRefs(value []QName) {
	t.ChoreographyRefField = value
}
func (t *Collaboration) ConversationLinks() (result *[]ConversationLink) {
	result = &t.ConversationLinkField
	return
}
func (t *Collaboration) SetConversationLinks(value []ConversationLink) {
	t.ConversationLinkField = value
}

type CompensateEventDefinition struct {
	EventDefinition
	WaitForCompletionField *bool    `xml:"waitForCompletion,attr,omitempty"`
	ActivityRefField       *QName   `xml:"activityRef,attr,omitempty"`
	TextPayloadField       *Payload `xml:",chardata"`
}

func DefaultCompensateEventDefinition() CompensateEventDefinition {
	return CompensateEventDefinition{
		EventDefinition: DefaultEventDefinition(),
	}
}

type CompensateEventDefinitionInterface interface {
	Element
	EventDefinitionInterface
	WaitForCompletion() (result bool, present bool)
	ActivityRef() (result *QName, present bool)
	SetWaitForCompletion(value *bool)
	SetActivityRef(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *CompensateEventDefinition) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *CompensateEventDefinition) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *CompensateEventDefinition) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.EventDefinition.FindBy(f); found {
		return
	}

	return
}
func (t *CompensateEventDefinition) WaitForCompletion() (result bool, present bool) {
	if t.WaitForCompletionField != nil {
		present = true
	}
	result = *t.WaitForCompletionField
	return
}
func (t *CompensateEventDefinition) SetWaitForCompletion(value *bool) {
	t.WaitForCompletionField = value
}
func (t *CompensateEventDefinition) ActivityRef() (result *QName, present bool) {
	if t.ActivityRefField != nil {
		present = true
	}
	result = t.ActivityRefField
	return
}
func (t *CompensateEventDefinition) SetActivityRef(value *QName) {
	t.ActivityRefField = value
}

type ComplexBehaviorDefinition struct {
	BaseElement
	ConditionField   FormalExpression    `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL condition"`
	EventField       *ImplicitThrowEvent `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL event"`
	TextPayloadField *Payload            `xml:",chardata"`
}

func DefaultComplexBehaviorDefinition() ComplexBehaviorDefinition {
	return ComplexBehaviorDefinition{
		BaseElement: DefaultBaseElement(),
	}
}

type ComplexBehaviorDefinitionInterface interface {
	Element
	BaseElementInterface
	Condition() (result *FormalExpression)
	Event() (result *ImplicitThrowEvent, present bool)
	SetCondition(value FormalExpression)
	SetEvent(value *ImplicitThrowEvent)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ComplexBehaviorDefinition) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ComplexBehaviorDefinition) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ComplexBehaviorDefinition) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	if result, found = t.ConditionField.FindBy(f); found {
		return
	}

	if value := t.EventField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	return
}
func (t *ComplexBehaviorDefinition) Condition() (result *FormalExpression) {
	result = &t.ConditionField
	return
}
func (t *ComplexBehaviorDefinition) SetCondition(value FormalExpression) {
	t.ConditionField = value
}
func (t *ComplexBehaviorDefinition) Event() (result *ImplicitThrowEvent, present bool) {
	if t.EventField != nil {
		present = true
	}
	result = t.EventField
	return
}
func (t *ComplexBehaviorDefinition) SetEvent(value *ImplicitThrowEvent) {
	t.EventField = value
}

type ComplexGateway struct {
	Gateway
	DefaultField             *IdRef        `xml:"default,attr,omitempty"`
	ActivationConditionField *AnExpression `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL activationCondition"`
	TextPayloadField         *Payload      `xml:",chardata"`
}

func DefaultComplexGateway() ComplexGateway {
	return ComplexGateway{
		Gateway: DefaultGateway(),
	}
}

type ComplexGatewayInterface interface {
	Element
	GatewayInterface
	Default() (result *IdRef, present bool)
	ActivationCondition() (result *AnExpression, present bool)
	SetDefault(value *IdRef)
	SetActivationCondition(value *AnExpression)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ComplexGateway) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ComplexGateway) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ComplexGateway) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Gateway.FindBy(f); found {
		return
	}

	if value := t.ActivationConditionField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	return
}
func (t *ComplexGateway) Default() (result *IdRef, present bool) {
	if t.DefaultField != nil {
		present = true
	}
	result = t.DefaultField
	return
}
func (t *ComplexGateway) SetDefault(value *IdRef) {
	t.DefaultField = value
}
func (t *ComplexGateway) ActivationCondition() (result *AnExpression, present bool) {
	if t.ActivationConditionField != nil {
		present = true
	}
	result = t.ActivationConditionField
	return
}
func (t *ComplexGateway) SetActivationCondition(value *AnExpression) {
	t.ActivationConditionField = value
}

type ConditionalEventDefinition struct {
	EventDefinition
	ConditionField   AnExpression `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL condition"`
	TextPayloadField *Payload     `xml:",chardata"`
}

func DefaultConditionalEventDefinition() ConditionalEventDefinition {
	return ConditionalEventDefinition{
		EventDefinition: DefaultEventDefinition(),
	}
}

type ConditionalEventDefinitionInterface interface {
	Element
	EventDefinitionInterface
	Condition() (result *AnExpression)
	SetCondition(value AnExpression)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ConditionalEventDefinition) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ConditionalEventDefinition) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ConditionalEventDefinition) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.EventDefinition.FindBy(f); found {
		return
	}

	if result, found = t.ConditionField.FindBy(f); found {
		return
	}

	return
}
func (t *ConditionalEventDefinition) Condition() (result *AnExpression) {
	result = &t.ConditionField
	return
}
func (t *ConditionalEventDefinition) SetCondition(value AnExpression) {
	t.ConditionField = value
}

type Conversation struct {
	ConversationNode
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultConversation() Conversation {
	return Conversation{
		ConversationNode: DefaultConversationNode(),
	}
}

type ConversationInterface interface {
	Element
	ConversationNodeInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Conversation) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Conversation) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Conversation) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.ConversationNode.FindBy(f); found {
		return
	}

	return
}

type ConversationAssociation struct {
	BaseElement
	InnerConversationNodeRefField QName    `xml:"innerConversationNodeRef,attr,omitempty"`
	OuterConversationNodeRefField QName    `xml:"outerConversationNodeRef,attr,omitempty"`
	TextPayloadField              *Payload `xml:",chardata"`
}

func DefaultConversationAssociation() ConversationAssociation {
	return ConversationAssociation{
		BaseElement: DefaultBaseElement(),
	}
}

type ConversationAssociationInterface interface {
	Element
	BaseElementInterface
	InnerConversationNodeRef() (result *QName)
	OuterConversationNodeRef() (result *QName)
	SetInnerConversationNodeRef(value QName)
	SetOuterConversationNodeRef(value QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ConversationAssociation) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ConversationAssociation) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ConversationAssociation) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}
func (t *ConversationAssociation) InnerConversationNodeRef() (result *QName) {
	result = &t.InnerConversationNodeRefField
	return
}
func (t *ConversationAssociation) SetInnerConversationNodeRef(value QName) {
	t.InnerConversationNodeRefField = value
}
func (t *ConversationAssociation) OuterConversationNodeRef() (result *QName) {
	result = &t.OuterConversationNodeRefField
	return
}
func (t *ConversationAssociation) SetOuterConversationNodeRef(value QName) {
	t.OuterConversationNodeRefField = value
}

type ConversationLink struct {
	BaseElement
	NameField        *string  `xml:"name,attr,omitempty"`
	SourceRefField   QName    `xml:"sourceRef,attr,omitempty"`
	TargetRefField   QName    `xml:"targetRef,attr,omitempty"`
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultConversationLink() ConversationLink {
	return ConversationLink{
		BaseElement: DefaultBaseElement(),
	}
}

type ConversationLinkInterface interface {
	Element
	BaseElementInterface
	Name() (result *string, present bool)
	SourceRef() (result *QName)
	TargetRef() (result *QName)
	SetName(value *string)
	SetSourceRef(value QName)
	SetTargetRef(value QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ConversationLink) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ConversationLink) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ConversationLink) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}
func (t *ConversationLink) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *ConversationLink) SetName(value *string) {
	t.NameField = value
}
func (t *ConversationLink) SourceRef() (result *QName) {
	result = &t.SourceRefField
	return
}
func (t *ConversationLink) SetSourceRef(value QName) {
	t.SourceRefField = value
}
func (t *ConversationLink) TargetRef() (result *QName) {
	result = &t.TargetRefField
	return
}
func (t *ConversationLink) SetTargetRef(value QName) {
	t.TargetRefField = value
}

type ConversationNode struct {
	BaseElement
	NameField           *string          `xml:"name,attr,omitempty"`
	ParticipantRefField []QName          `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL participantRef"`
	MessageFlowRefField []QName          `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL messageFlowRef"`
	CorrelationKeyField []CorrelationKey `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL correlationKey"`
}

func DefaultConversationNode() ConversationNode {
	return ConversationNode{
		BaseElement: DefaultBaseElement(),
	}
}

type ConversationNodeInterface interface {
	Element
	BaseElementInterface
	Name() (result *string, present bool)
	ParticipantRefs() (result *[]QName)
	MessageFlowRefs() (result *[]QName)
	CorrelationKeys() (result *[]CorrelationKey)
	SetName(value *string)
	SetParticipantRefs(value []QName)
	SetMessageFlowRefs(value []QName)
	SetCorrelationKeys(value []CorrelationKey)
}

func (t *ConversationNode) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	for i := range t.CorrelationKeyField {
		if result, found = t.CorrelationKeyField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *ConversationNode) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *ConversationNode) SetName(value *string) {
	t.NameField = value
}
func (t *ConversationNode) ParticipantRefs() (result *[]QName) {
	result = &t.ParticipantRefField
	return
}
func (t *ConversationNode) SetParticipantRefs(value []QName) {
	t.ParticipantRefField = value
}
func (t *ConversationNode) MessageFlowRefs() (result *[]QName) {
	result = &t.MessageFlowRefField
	return
}
func (t *ConversationNode) SetMessageFlowRefs(value []QName) {
	t.MessageFlowRefField = value
}
func (t *ConversationNode) CorrelationKeys() (result *[]CorrelationKey) {
	result = &t.CorrelationKeyField
	return
}
func (t *ConversationNode) SetCorrelationKeys(value []CorrelationKey) {
	t.CorrelationKeyField = value
}

type CorrelationKey struct {
	BaseElement
	NameField                   *string  `xml:"name,attr,omitempty"`
	CorrelationPropertyRefField []QName  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL correlationPropertyRef"`
	TextPayloadField            *Payload `xml:",chardata"`
}

func DefaultCorrelationKey() CorrelationKey {
	return CorrelationKey{
		BaseElement: DefaultBaseElement(),
	}
}

type CorrelationKeyInterface interface {
	Element
	BaseElementInterface
	Name() (result *string, present bool)
	CorrelationPropertyRefs() (result *[]QName)
	SetName(value *string)
	SetCorrelationPropertyRefs(value []QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *CorrelationKey) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *CorrelationKey) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *CorrelationKey) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}
func (t *CorrelationKey) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *CorrelationKey) SetName(value *string) {
	t.NameField = value
}
func (t *CorrelationKey) CorrelationPropertyRefs() (result *[]QName) {
	result = &t.CorrelationPropertyRefField
	return
}
func (t *CorrelationKey) SetCorrelationPropertyRefs(value []QName) {
	t.CorrelationPropertyRefField = value
}

type CorrelationProperty struct {
	RootElement
	NameField                                   *string                                  `xml:"name,attr,omitempty"`
	TypeField                                   *QName                                   `xml:"type,attr,omitempty"`
	CorrelationPropertyRetrievalExpressionField []CorrelationPropertyRetrievalExpression `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL correlationPropertyRetrievalExpression"`
	TextPayloadField                            *Payload                                 `xml:",chardata"`
}

func DefaultCorrelationProperty() CorrelationProperty {
	return CorrelationProperty{
		RootElement: DefaultRootElement(),
	}
}

type CorrelationPropertyInterface interface {
	Element
	RootElementInterface
	Name() (result *string, present bool)
	Type() (result *QName, present bool)
	CorrelationPropertyRetrievalExpressions() (result *[]CorrelationPropertyRetrievalExpression)
	SetName(value *string)
	SetType(value *QName)
	SetCorrelationPropertyRetrievalExpressions(value []CorrelationPropertyRetrievalExpression)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *CorrelationProperty) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *CorrelationProperty) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *CorrelationProperty) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.RootElement.FindBy(f); found {
		return
	}

	for i := range t.CorrelationPropertyRetrievalExpressionField {
		if result, found = t.CorrelationPropertyRetrievalExpressionField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *CorrelationProperty) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *CorrelationProperty) SetName(value *string) {
	t.NameField = value
}
func (t *CorrelationProperty) Type() (result *QName, present bool) {
	if t.TypeField != nil {
		present = true
	}
	result = t.TypeField
	return
}
func (t *CorrelationProperty) SetType(value *QName) {
	t.TypeField = value
}
func (t *CorrelationProperty) CorrelationPropertyRetrievalExpressions() (result *[]CorrelationPropertyRetrievalExpression) {
	result = &t.CorrelationPropertyRetrievalExpressionField
	return
}
func (t *CorrelationProperty) SetCorrelationPropertyRetrievalExpressions(value []CorrelationPropertyRetrievalExpression) {
	t.CorrelationPropertyRetrievalExpressionField = value
}

type CorrelationPropertyBinding struct {
	BaseElement
	CorrelationPropertyRefField QName            `xml:"correlationPropertyRef,attr,omitempty"`
	DataPathField               FormalExpression `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataPath"`
	TextPayloadField            *Payload         `xml:",chardata"`
}

func DefaultCorrelationPropertyBinding() CorrelationPropertyBinding {
	return CorrelationPropertyBinding{
		BaseElement: DefaultBaseElement(),
	}
}

type CorrelationPropertyBindingInterface interface {
	Element
	BaseElementInterface
	CorrelationPropertyRef() (result *QName)
	DataPath() (result *FormalExpression)
	SetCorrelationPropertyRef(value QName)
	SetDataPath(value FormalExpression)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *CorrelationPropertyBinding) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *CorrelationPropertyBinding) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *CorrelationPropertyBinding) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	if result, found = t.DataPathField.FindBy(f); found {
		return
	}

	return
}
func (t *CorrelationPropertyBinding) CorrelationPropertyRef() (result *QName) {
	result = &t.CorrelationPropertyRefField
	return
}
func (t *CorrelationPropertyBinding) SetCorrelationPropertyRef(value QName) {
	t.CorrelationPropertyRefField = value
}
func (t *CorrelationPropertyBinding) DataPath() (result *FormalExpression) {
	result = &t.DataPathField
	return
}
func (t *CorrelationPropertyBinding) SetDataPath(value FormalExpression) {
	t.DataPathField = value
}

type CorrelationPropertyRetrievalExpression struct {
	BaseElement
	MessageRefField  QName            `xml:"messageRef,attr,omitempty"`
	MessagePathField FormalExpression `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL messagePath"`
	TextPayloadField *Payload         `xml:",chardata"`
}

func DefaultCorrelationPropertyRetrievalExpression() CorrelationPropertyRetrievalExpression {
	return CorrelationPropertyRetrievalExpression{
		BaseElement: DefaultBaseElement(),
	}
}

type CorrelationPropertyRetrievalExpressionInterface interface {
	Element
	BaseElementInterface
	MessageRef() (result *QName)
	MessagePath() (result *FormalExpression)
	SetMessageRef(value QName)
	SetMessagePath(value FormalExpression)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *CorrelationPropertyRetrievalExpression) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *CorrelationPropertyRetrievalExpression) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *CorrelationPropertyRetrievalExpression) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	if result, found = t.MessagePathField.FindBy(f); found {
		return
	}

	return
}
func (t *CorrelationPropertyRetrievalExpression) MessageRef() (result *QName) {
	result = &t.MessageRefField
	return
}
func (t *CorrelationPropertyRetrievalExpression) SetMessageRef(value QName) {
	t.MessageRefField = value
}
func (t *CorrelationPropertyRetrievalExpression) MessagePath() (result *FormalExpression) {
	result = &t.MessagePathField
	return
}
func (t *CorrelationPropertyRetrievalExpression) SetMessagePath(value FormalExpression) {
	t.MessagePathField = value
}

type CorrelationSubscription struct {
	BaseElement
	CorrelationKeyRefField          QName                        `xml:"correlationKeyRef,attr,omitempty"`
	CorrelationPropertyBindingField []CorrelationPropertyBinding `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL correlationPropertyBinding"`
	TextPayloadField                *Payload                     `xml:",chardata"`
}

func DefaultCorrelationSubscription() CorrelationSubscription {
	return CorrelationSubscription{
		BaseElement: DefaultBaseElement(),
	}
}

type CorrelationSubscriptionInterface interface {
	Element
	BaseElementInterface
	CorrelationKeyRef() (result *QName)
	CorrelationPropertyBindings() (result *[]CorrelationPropertyBinding)
	SetCorrelationKeyRef(value QName)
	SetCorrelationPropertyBindings(value []CorrelationPropertyBinding)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *CorrelationSubscription) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *CorrelationSubscription) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *CorrelationSubscription) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	for i := range t.CorrelationPropertyBindingField {
		if result, found = t.CorrelationPropertyBindingField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *CorrelationSubscription) CorrelationKeyRef() (result *QName) {
	result = &t.CorrelationKeyRefField
	return
}
func (t *CorrelationSubscription) SetCorrelationKeyRef(value QName) {
	t.CorrelationKeyRefField = value
}
func (t *CorrelationSubscription) CorrelationPropertyBindings() (result *[]CorrelationPropertyBinding) {
	result = &t.CorrelationPropertyBindingField
	return
}
func (t *CorrelationSubscription) SetCorrelationPropertyBindings(value []CorrelationPropertyBinding) {
	t.CorrelationPropertyBindingField = value
}

type DataAssociation struct {
	BaseElement
	SourceRefField      []IdRef           `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL sourceRef"`
	TargetRefField      IdRef             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL targetRef"`
	TransformationField *FormalExpression `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL transformation"`
	AssignmentField     []Assignment      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL assignment"`
	TextPayloadField    *Payload          `xml:",chardata"`
}

func DefaultDataAssociation() DataAssociation {
	return DataAssociation{
		BaseElement: DefaultBaseElement(),
	}
}

type DataAssociationInterface interface {
	Element
	BaseElementInterface
	SourceRefs() (result *[]IdRef)
	TargetRef() (result *IdRef)
	Transformation() (result *FormalExpression, present bool)
	Assignments() (result *[]Assignment)
	SetSourceRefs(value []IdRef)
	SetTargetRef(value IdRef)
	SetTransformation(value *FormalExpression)
	SetAssignments(value []Assignment)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *DataAssociation) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *DataAssociation) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *DataAssociation) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	if value := t.TransformationField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	for i := range t.AssignmentField {
		if result, found = t.AssignmentField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *DataAssociation) SourceRefs() (result *[]IdRef) {
	result = &t.SourceRefField
	return
}
func (t *DataAssociation) SetSourceRefs(value []IdRef) {
	t.SourceRefField = value
}
func (t *DataAssociation) TargetRef() (result *IdRef) {
	result = &t.TargetRefField
	return
}
func (t *DataAssociation) SetTargetRef(value IdRef) {
	t.TargetRefField = value
}
func (t *DataAssociation) Transformation() (result *FormalExpression, present bool) {
	if t.TransformationField != nil {
		present = true
	}
	result = t.TransformationField
	return
}
func (t *DataAssociation) SetTransformation(value *FormalExpression) {
	t.TransformationField = value
}
func (t *DataAssociation) Assignments() (result *[]Assignment) {
	result = &t.AssignmentField
	return
}
func (t *DataAssociation) SetAssignments(value []Assignment) {
	t.AssignmentField = value
}

type DataInput struct {
	ItemAwareElement
	NameField         *string  `xml:"name,attr,omitempty"`
	IsCollectionField *bool    `xml:"isCollection,attr,omitempty"`
	TextPayloadField  *Payload `xml:",chardata"`
}

var defaultDataInputIsCollectionField bool = false

func DefaultDataInput() DataInput {
	return DataInput{
		ItemAwareElement:  DefaultItemAwareElement(),
		IsCollectionField: &defaultDataInputIsCollectionField,
	}
}

type DataInputInterface interface {
	Element
	ItemAwareElementInterface
	Name() (result *string, present bool)
	IsCollection() (result bool)
	SetName(value *string)
	SetIsCollection(value *bool)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *DataInput) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *DataInput) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *DataInput) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.ItemAwareElement.FindBy(f); found {
		return
	}

	return
}
func (t *DataInput) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *DataInput) SetName(value *string) {
	t.NameField = value
}
func (t *DataInput) IsCollection() (result bool) {
	if t.IsCollectionField == nil {
		result = defaultDataInputIsCollectionField
		return
	}
	result = *t.IsCollectionField
	return
}
func (t *DataInput) SetIsCollection(value *bool) {
	t.IsCollectionField = value
}

type DataInputAssociation struct {
	DataAssociation
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultDataInputAssociation() DataInputAssociation {
	return DataInputAssociation{
		DataAssociation: DefaultDataAssociation(),
	}
}

type DataInputAssociationInterface interface {
	Element
	DataAssociationInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *DataInputAssociation) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *DataInputAssociation) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *DataInputAssociation) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.DataAssociation.FindBy(f); found {
		return
	}

	return
}

type DataObject struct {
	FlowElement
	ItemAware
	IsCollectionField *bool    `xml:"isCollection,attr,omitempty"`
	TextPayloadField  *Payload `xml:",chardata"`
}

var defaultDataObjectIsCollectionField bool = false

func DefaultDataObject() DataObject {
	return DataObject{
		FlowElement:       DefaultFlowElement(),
		ItemAware:         DefaultItemAware(),
		IsCollectionField: &defaultDataObjectIsCollectionField,
	}
}

type DataObjectInterface interface {
	Element
	FlowElementInterface
	ItemAwareInterface
	IsCollection() (result bool)
	SetIsCollection(value *bool)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *DataObject) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *DataObject) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *DataObject) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.FlowElement.FindBy(f); found {
		return
	}
	if result, found = t.ItemAware.FindBy(f); found {
		return
	}

	return
}
func (t *DataObject) IsCollection() (result bool) {
	if t.IsCollectionField == nil {
		result = defaultDataObjectIsCollectionField
		return
	}
	result = *t.IsCollectionField
	return
}
func (t *DataObject) SetIsCollection(value *bool) {
	t.IsCollectionField = value
}

type DataObjectReference struct {
	FlowElement
	ItemAware
	DataObjectRefField *IdRef   `xml:"dataObjectRef,attr,omitempty"`
	TextPayloadField   *Payload `xml:",chardata"`
}

func DefaultDataObjectReference() DataObjectReference {
	return DataObjectReference{
		FlowElement: DefaultFlowElement(),
		ItemAware:   DefaultItemAware(),
	}
}

type DataObjectReferenceInterface interface {
	Element
	FlowElementInterface
	ItemAwareInterface
	DataObjectRef() (result *IdRef, present bool)
	SetDataObjectRef(value *IdRef)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *DataObjectReference) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *DataObjectReference) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *DataObjectReference) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.FlowElement.FindBy(f); found {
		return
	}
	if result, found = t.ItemAware.FindBy(f); found {
		return
	}

	return
}
func (t *DataObjectReference) DataObjectRef() (result *IdRef, present bool) {
	if t.DataObjectRefField != nil {
		present = true
	}
	result = t.DataObjectRefField
	return
}
func (t *DataObjectReference) SetDataObjectRef(value *IdRef) {
	t.DataObjectRefField = value
}

type DataOutput struct {
	ItemAwareElement
	NameField         *string  `xml:"name,attr,omitempty"`
	IsCollectionField *bool    `xml:"isCollection,attr,omitempty"`
	TextPayloadField  *Payload `xml:",chardata"`
}

var defaultDataOutputIsCollectionField bool = false

func DefaultDataOutput() DataOutput {
	return DataOutput{
		ItemAwareElement:  DefaultItemAwareElement(),
		IsCollectionField: &defaultDataOutputIsCollectionField,
	}
}

type DataOutputInterface interface {
	Element
	ItemAwareElementInterface
	Name() (result *string, present bool)
	IsCollection() (result bool)
	SetName(value *string)
	SetIsCollection(value *bool)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *DataOutput) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *DataOutput) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *DataOutput) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.ItemAwareElement.FindBy(f); found {
		return
	}

	return
}
func (t *DataOutput) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *DataOutput) SetName(value *string) {
	t.NameField = value
}
func (t *DataOutput) IsCollection() (result bool) {
	if t.IsCollectionField == nil {
		result = defaultDataOutputIsCollectionField
		return
	}
	result = *t.IsCollectionField
	return
}
func (t *DataOutput) SetIsCollection(value *bool) {
	t.IsCollectionField = value
}

type DataOutputAssociation struct {
	DataAssociation
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultDataOutputAssociation() DataOutputAssociation {
	return DataOutputAssociation{
		DataAssociation: DefaultDataAssociation(),
	}
}

type DataOutputAssociationInterface interface {
	Element
	DataAssociationInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *DataOutputAssociation) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *DataOutputAssociation) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *DataOutputAssociation) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.DataAssociation.FindBy(f); found {
		return
	}

	return
}

type DataState struct {
	BaseElement
	NameField        *string  `xml:"name,attr,omitempty"`
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultDataState() DataState {
	return DataState{
		BaseElement: DefaultBaseElement(),
	}
}

type DataStateInterface interface {
	Element
	BaseElementInterface
	Name() (result *string, present bool)
	SetName(value *string)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *DataState) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *DataState) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *DataState) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}
func (t *DataState) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *DataState) SetName(value *string) {
	t.NameField = value
}

type DataStore struct {
	RootElement
	ItemAware
	NameField        *string  `xml:"name,attr,omitempty"`
	CapacityField    *big.Int `xml:"capacity,attr,omitempty"`
	IsUnlimitedField *bool    `xml:"isUnlimited,attr,omitempty"`
	TextPayloadField *Payload `xml:",chardata"`
}

var defaultDataStoreIsUnlimitedField bool = true

func DefaultDataStore() DataStore {
	return DataStore{
		RootElement:      DefaultRootElement(),
		ItemAware:        DefaultItemAware(),
		IsUnlimitedField: &defaultDataStoreIsUnlimitedField,
	}
}

type DataStoreInterface interface {
	Element
	RootElementInterface
	ItemAwareInterface
	Name() (result *string, present bool)
	Capacity() (result *big.Int, present bool)
	IsUnlimited() (result bool)
	SetName(value *string)
	SetCapacity(value *big.Int)
	SetIsUnlimited(value *bool)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *DataStore) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *DataStore) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *DataStore) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.RootElement.FindBy(f); found {
		return
	}
	if result, found = t.ItemAware.FindBy(f); found {
		return
	}

	return
}
func (t *DataStore) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *DataStore) SetName(value *string) {
	t.NameField = value
}
func (t *DataStore) Capacity() (result *big.Int, present bool) {
	if t.CapacityField != nil {
		present = true
	}
	result = t.CapacityField
	return
}
func (t *DataStore) SetCapacity(value *big.Int) {
	t.CapacityField = value
}
func (t *DataStore) IsUnlimited() (result bool) {
	if t.IsUnlimitedField == nil {
		result = defaultDataStoreIsUnlimitedField
		return
	}
	result = *t.IsUnlimitedField
	return
}
func (t *DataStore) SetIsUnlimited(value *bool) {
	t.IsUnlimitedField = value
}

type DataStoreReference struct {
	FlowElement
	ItemAware
	DataStoreRefField *QName   `xml:"dataStoreRef,attr,omitempty"`
	TextPayloadField  *Payload `xml:",chardata"`
}

func DefaultDataStoreReference() DataStoreReference {
	return DataStoreReference{
		FlowElement: DefaultFlowElement(),
		ItemAware:   DefaultItemAware(),
	}
}

type DataStoreReferenceInterface interface {
	Element
	FlowElementInterface
	ItemAwareInterface
	DataStoreRef() (result *QName, present bool)
	SetDataStoreRef(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *DataStoreReference) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *DataStoreReference) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *DataStoreReference) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.FlowElement.FindBy(f); found {
		return
	}
	if result, found = t.ItemAware.FindBy(f); found {
		return
	}

	return
}
func (t *DataStoreReference) DataStoreRef() (result *QName, present bool) {
	if t.DataStoreRefField != nil {
		present = true
	}
	result = t.DataStoreRefField
	return
}
func (t *DataStoreReference) SetDataStoreRef(value *QName) {
	t.DataStoreRefField = value
}

type Documentation struct {
	IdField          *Id      `xml:"id,attr,omitempty"`
	TextFormatField  *string  `xml:"textFormat,attr,omitempty"`
	TextPayloadField *Payload `xml:",chardata"`
}

var defaultDocumentationTextFormatField string = "text/plain"

func DefaultDocumentation() Documentation {
	return Documentation{
		TextFormatField: &defaultDocumentationTextFormatField,
	}
}

type DocumentationInterface interface {
	Element
	Id() (result *Id, present bool)
	TextFormat() (result *string)
	SetId(value *Id)
	SetTextFormat(value *string)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Documentation) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Documentation) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Documentation) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}

	return
}
func (t *Documentation) Id() (result *Id, present bool) {
	if t.IdField != nil {
		present = true
	}
	result = t.IdField
	return
}
func (t *Documentation) SetId(value *Id) {
	t.IdField = value
}
func (t *Documentation) TextFormat() (result *string) {
	if t.TextFormatField == nil {
		result = &defaultDocumentationTextFormatField
		return
	}
	result = t.TextFormatField
	return
}
func (t *Documentation) SetTextFormat(value *string) {
	t.TextFormatField = value
}

type EndEvent struct {
	ThrowEvent
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultEndEvent() EndEvent {
	return EndEvent{
		ThrowEvent: DefaultThrowEvent(),
	}
}

type EndEventInterface interface {
	Element
	ThrowEventInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *EndEvent) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *EndEvent) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *EndEvent) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.ThrowEvent.FindBy(f); found {
		return
	}

	return
}

type EndPoint struct {
	RootElement
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultEndPoint() EndPoint {
	return EndPoint{
		RootElement: DefaultRootElement(),
	}
}

type EndPointInterface interface {
	Element
	RootElementInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *EndPoint) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *EndPoint) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *EndPoint) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.RootElement.FindBy(f); found {
		return
	}

	return
}

type Error struct {
	RootElement
	NameField         *string  `xml:"name,attr,omitempty"`
	ErrorCodeField    *string  `xml:"errorCode,attr,omitempty"`
	StructureRefField *QName   `xml:"structureRef,attr,omitempty"`
	TextPayloadField  *Payload `xml:",chardata"`
}

func DefaultError() Error {
	return Error{
		RootElement: DefaultRootElement(),
	}
}

type ErrorInterface interface {
	Element
	RootElementInterface
	Name() (result *string, present bool)
	ErrorCode() (result *string, present bool)
	StructureRef() (result *QName, present bool)
	SetName(value *string)
	SetErrorCode(value *string)
	SetStructureRef(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Error) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Error) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Error) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.RootElement.FindBy(f); found {
		return
	}

	return
}
func (t *Error) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *Error) SetName(value *string) {
	t.NameField = value
}
func (t *Error) ErrorCode() (result *string, present bool) {
	if t.ErrorCodeField != nil {
		present = true
	}
	result = t.ErrorCodeField
	return
}
func (t *Error) SetErrorCode(value *string) {
	t.ErrorCodeField = value
}
func (t *Error) StructureRef() (result *QName, present bool) {
	if t.StructureRefField != nil {
		present = true
	}
	result = t.StructureRefField
	return
}
func (t *Error) SetStructureRef(value *QName) {
	t.StructureRefField = value
}

type ErrorEventDefinition struct {
	EventDefinition
	ErrorRefField    *QName   `xml:"errorRef,attr,omitempty"`
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultErrorEventDefinition() ErrorEventDefinition {
	return ErrorEventDefinition{
		EventDefinition: DefaultEventDefinition(),
	}
}

type ErrorEventDefinitionInterface interface {
	Element
	EventDefinitionInterface
	ErrorRef() (result *QName, present bool)
	SetErrorRef(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ErrorEventDefinition) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ErrorEventDefinition) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ErrorEventDefinition) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.EventDefinition.FindBy(f); found {
		return
	}

	return
}
func (t *ErrorEventDefinition) ErrorRef() (result *QName, present bool) {
	if t.ErrorRefField != nil {
		present = true
	}
	result = t.ErrorRefField
	return
}
func (t *ErrorEventDefinition) SetErrorRef(value *QName) {
	t.ErrorRefField = value
}

type Escalation struct {
	RootElement
	NameField           *string  `xml:"name,attr,omitempty"`
	EscalationCodeField *string  `xml:"escalationCode,attr,omitempty"`
	StructureRefField   *QName   `xml:"structureRef,attr,omitempty"`
	TextPayloadField    *Payload `xml:",chardata"`
}

func DefaultEscalation() Escalation {
	return Escalation{
		RootElement: DefaultRootElement(),
	}
}

type EscalationInterface interface {
	Element
	RootElementInterface
	Name() (result *string, present bool)
	EscalationCode() (result *string, present bool)
	StructureRef() (result *QName, present bool)
	SetName(value *string)
	SetEscalationCode(value *string)
	SetStructureRef(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Escalation) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Escalation) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Escalation) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.RootElement.FindBy(f); found {
		return
	}

	return
}
func (t *Escalation) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *Escalation) SetName(value *string) {
	t.NameField = value
}
func (t *Escalation) EscalationCode() (result *string, present bool) {
	if t.EscalationCodeField != nil {
		present = true
	}
	result = t.EscalationCodeField
	return
}
func (t *Escalation) SetEscalationCode(value *string) {
	t.EscalationCodeField = value
}
func (t *Escalation) StructureRef() (result *QName, present bool) {
	if t.StructureRefField != nil {
		present = true
	}
	result = t.StructureRefField
	return
}
func (t *Escalation) SetStructureRef(value *QName) {
	t.StructureRefField = value
}

type EscalationEventDefinition struct {
	EventDefinition
	EscalationRefField *QName   `xml:"escalationRef,attr,omitempty"`
	TextPayloadField   *Payload `xml:",chardata"`
}

func DefaultEscalationEventDefinition() EscalationEventDefinition {
	return EscalationEventDefinition{
		EventDefinition: DefaultEventDefinition(),
	}
}

type EscalationEventDefinitionInterface interface {
	Element
	EventDefinitionInterface
	EscalationRef() (result *QName, present bool)
	SetEscalationRef(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *EscalationEventDefinition) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *EscalationEventDefinition) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *EscalationEventDefinition) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.EventDefinition.FindBy(f); found {
		return
	}

	return
}
func (t *EscalationEventDefinition) EscalationRef() (result *QName, present bool) {
	if t.EscalationRefField != nil {
		present = true
	}
	result = t.EscalationRefField
	return
}
func (t *EscalationEventDefinition) SetEscalationRef(value *QName) {
	t.EscalationRefField = value
}

type Event struct {
	FlowNode
	PropertyField []Property `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL property"`
}

func DefaultEvent() Event {
	return Event{
		FlowNode: DefaultFlowNode(),
	}
}

type EventInterface interface {
	Element
	FlowNodeInterface
	Properties() (result *[]Property)
	SetProperties(value []Property)
}

func (t *Event) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.FlowNode.FindBy(f); found {
		return
	}

	for i := range t.PropertyField {
		if result, found = t.PropertyField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *Event) Properties() (result *[]Property) {
	result = &t.PropertyField
	return
}
func (t *Event) SetProperties(value []Property) {
	t.PropertyField = value
}

type EventBasedGateway struct {
	Gateway
	InstantiateField      *bool                  `xml:"instantiate,attr,omitempty"`
	EventGatewayTypeField *EventBasedGatewayType `xml:"eventGatewayType,attr,omitempty"`
	TextPayloadField      *Payload               `xml:",chardata"`
}

var defaultEventBasedGatewayInstantiateField bool = false
var defaultEventBasedGatewayEventGatewayTypeField EventBasedGatewayType = "Exclusive"

func DefaultEventBasedGateway() EventBasedGateway {
	return EventBasedGateway{
		Gateway:               DefaultGateway(),
		InstantiateField:      &defaultEventBasedGatewayInstantiateField,
		EventGatewayTypeField: &defaultEventBasedGatewayEventGatewayTypeField,
	}
}

type EventBasedGatewayInterface interface {
	Element
	GatewayInterface
	Instantiate() (result bool)
	EventGatewayType() (result *EventBasedGatewayType)
	SetInstantiate(value *bool)
	SetEventGatewayType(value *EventBasedGatewayType)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *EventBasedGateway) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *EventBasedGateway) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *EventBasedGateway) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Gateway.FindBy(f); found {
		return
	}

	return
}
func (t *EventBasedGateway) Instantiate() (result bool) {
	if t.InstantiateField == nil {
		result = defaultEventBasedGatewayInstantiateField
		return
	}
	result = *t.InstantiateField
	return
}
func (t *EventBasedGateway) SetInstantiate(value *bool) {
	t.InstantiateField = value
}
func (t *EventBasedGateway) EventGatewayType() (result *EventBasedGatewayType) {
	if t.EventGatewayTypeField == nil {
		result = &defaultEventBasedGatewayEventGatewayTypeField
		return
	}
	result = t.EventGatewayTypeField
	return
}
func (t *EventBasedGateway) SetEventGatewayType(value *EventBasedGatewayType) {
	t.EventGatewayTypeField = value
}

type EventDefinition struct {
	RootElement
}

func DefaultEventDefinition() EventDefinition {
	return EventDefinition{
		RootElement: DefaultRootElement(),
	}
}

type EventDefinitionInterface interface {
	Element
	RootElementInterface
}

func (t *EventDefinition) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.RootElement.FindBy(f); found {
		return
	}

	return
}

type ExclusiveGateway struct {
	Gateway
	DefaultField     *IdRef   `xml:"default,attr,omitempty"`
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultExclusiveGateway() ExclusiveGateway {
	return ExclusiveGateway{
		Gateway: DefaultGateway(),
	}
}

type ExclusiveGatewayInterface interface {
	Element
	GatewayInterface
	Default() (result *IdRef, present bool)
	SetDefault(value *IdRef)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ExclusiveGateway) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ExclusiveGateway) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ExclusiveGateway) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Gateway.FindBy(f); found {
		return
	}

	return
}
func (t *ExclusiveGateway) Default() (result *IdRef, present bool) {
	if t.DefaultField != nil {
		present = true
	}
	result = t.DefaultField
	return
}
func (t *ExclusiveGateway) SetDefault(value *IdRef) {
	t.DefaultField = value
}

type Expression struct {
	BaseElementWithMixedContent
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultExpression() Expression {
	return Expression{
		BaseElementWithMixedContent: DefaultBaseElementWithMixedContent(),
	}
}

type ExpressionInterface interface {
	Element
	BaseElementWithMixedContentInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Expression) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Expression) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Expression) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElementWithMixedContent.FindBy(f); found {
		return
	}

	return
}

type Extension struct {
	DefinitionField     *QName          `xml:"definition,attr,omitempty"`
	MustUnderstandField *bool           `xml:"mustUnderstand,attr,omitempty"`
	DocumentationField  []Documentation `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL documentation"`
	TextPayloadField    *Payload        `xml:",chardata"`
}

var defaultExtensionMustUnderstandField bool = false

func DefaultExtension() Extension {
	return Extension{
		MustUnderstandField: &defaultExtensionMustUnderstandField,
	}
}

type ExtensionInterface interface {
	Element
	Definition() (result *QName, present bool)
	MustUnderstand() (result *bool)
	Documentations() (result *[]Documentation)
	SetDefinition(value *QName)
	SetMustUnderstand(value *bool)
	SetDocumentations(value []Documentation)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Extension) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Extension) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Extension) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}

	for i := range t.DocumentationField {
		if result, found = t.DocumentationField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *Extension) Definition() (result *QName, present bool) {
	if t.DefinitionField != nil {
		present = true
	}
	result = t.DefinitionField
	return
}
func (t *Extension) SetDefinition(value *QName) {
	t.DefinitionField = value
}
func (t *Extension) MustUnderstand() (result *bool) {
	if t.MustUnderstandField == nil {
		result = &defaultExtensionMustUnderstandField
		return
	}
	result = t.MustUnderstandField
	return
}
func (t *Extension) SetMustUnderstand(value *bool) {
	t.MustUnderstandField = value
}
func (t *Extension) Documentations() (result *[]Documentation) {
	result = &t.DocumentationField
	return
}
func (t *Extension) SetDocumentations(value []Documentation) {
	t.DocumentationField = value
}

type ExtensionElements struct {
	DataObjectBody      *ExtensionDataObjectBody `xml:"http://olive.io/spec/BPMN/MODEL dataObjectBody"`
	TaskDefinitionField *TaskDefinition          `xml:"http://olive.io/spec/BPMN/MODEL taskDefinition"`
	TaskHeaderField     *TaskHeader              `xml:"http://olive.io/spec/BPMN/MODEL taskHeaders"`
	PropertiesField     *Properties              `xml:"http://olive.io/spec/BPMN/MODEL properties"`
	ScriptField         *ExtensionScript         `xml:"http://olive.io/spec/BPMN/MODEL script"`
	CalledElement       *ExtensionCalledElement  `xml:"http://olive.io/spec/BPMN/MODEL calledElement"`
	CalledDecision      *ExtensionCalledDecision `xml:"http://olive.io/spec/BPMN/MODEL calledDecision"`
	DataInput           []ExtensionAssociation   `xml:"http://olive.io/spec/BPMN/MODEL dataInput"`
	DataOutput          []ExtensionAssociation   `xml:"http://olive.io/spec/BPMN/MODEL dataOutput"`
	TextPayloadField    *Payload                 `xml:",chardata"`
}

func DefaultExtensionElements() ExtensionElements {
	return ExtensionElements{}
}

type ExtensionElementsInterface interface {
	Element

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ExtensionElements) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ExtensionElements) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ExtensionElements) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}

	return
}

type FlowElement struct {
	BaseElement
	NameField             *string     `xml:"name,attr,omitempty"`
	AuditingField         *Auditing   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL auditing"`
	MonitoringField       *Monitoring `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL monitoring"`
	CategoryValueRefField []QName     `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL categoryValueRef"`
}

func DefaultFlowElement() FlowElement {
	return FlowElement{
		BaseElement: DefaultBaseElement(),
	}
}

type FlowElementInterface interface {
	Element
	BaseElementInterface
	Name() (result *string, present bool)
	Auditing() (result *Auditing, present bool)
	Monitoring() (result *Monitoring, present bool)
	CategoryValueRefs() (result *[]QName)
	SetName(value *string)
	SetAuditing(value *Auditing)
	SetMonitoring(value *Monitoring)
	SetCategoryValueRefs(value []QName)
}

func (t *FlowElement) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	if value := t.AuditingField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	if value := t.MonitoringField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	return
}
func (t *FlowElement) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *FlowElement) SetName(value *string) {
	t.NameField = value
}
func (t *FlowElement) Auditing() (result *Auditing, present bool) {
	if t.AuditingField != nil {
		present = true
	}
	result = t.AuditingField
	return
}
func (t *FlowElement) SetAuditing(value *Auditing) {
	t.AuditingField = value
}
func (t *FlowElement) Monitoring() (result *Monitoring, present bool) {
	if t.MonitoringField != nil {
		present = true
	}
	result = t.MonitoringField
	return
}
func (t *FlowElement) SetMonitoring(value *Monitoring) {
	t.MonitoringField = value
}
func (t *FlowElement) CategoryValueRefs() (result *[]QName) {
	result = &t.CategoryValueRefField
	return
}
func (t *FlowElement) SetCategoryValueRefs(value []QName) {
	t.CategoryValueRefField = value
}

type FlowNode struct {
	FlowElement
	IncomingField []QName `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL incoming"`
	OutgoingField []QName `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL outgoing"`
}

func DefaultFlowNode() FlowNode {
	return FlowNode{
		FlowElement: DefaultFlowElement(),
	}
}

type FlowNodeInterface interface {
	Element
	FlowElementInterface
	Incomings() (result *[]QName)
	Outgoings() (result *[]QName)
	SetIncomings(value []QName)
	SetOutgoings(value []QName)
}

func (t *FlowNode) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.FlowElement.FindBy(f); found {
		return
	}

	return
}
func (t *FlowNode) Incomings() (result *[]QName) {
	result = &t.IncomingField
	return
}
func (t *FlowNode) SetIncomings(value []QName) {
	t.IncomingField = value
}
func (t *FlowNode) Outgoings() (result *[]QName) {
	result = &t.OutgoingField
	return
}
func (t *FlowNode) SetOutgoings(value []QName) {
	t.OutgoingField = value
}

type FormalExpression struct {
	Expression
	LanguageField           *AnyURI  `xml:"language,attr,omitempty"`
	EvaluatesToTypeRefField *QName   `xml:"evaluatesToTypeRef,attr,omitempty"`
	TextPayloadField        *Payload `xml:",chardata"`
}

func DefaultFormalExpression() FormalExpression {
	return FormalExpression{
		Expression: DefaultExpression(),
	}
}

type FormalExpressionInterface interface {
	Element
	ExpressionInterface
	Language() (result *AnyURI, present bool)
	EvaluatesToTypeRef() (result *QName, present bool)
	SetLanguage(value *AnyURI)
	SetEvaluatesToTypeRef(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *FormalExpression) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *FormalExpression) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *FormalExpression) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Expression.FindBy(f); found {
		return
	}

	return
}
func (t *FormalExpression) Language() (result *AnyURI, present bool) {
	if t.LanguageField != nil {
		present = true
	}
	result = t.LanguageField
	return
}
func (t *FormalExpression) SetLanguage(value *AnyURI) {
	t.LanguageField = value
}
func (t *FormalExpression) EvaluatesToTypeRef() (result *QName, present bool) {
	if t.EvaluatesToTypeRefField != nil {
		present = true
	}
	result = t.EvaluatesToTypeRefField
	return
}
func (t *FormalExpression) SetEvaluatesToTypeRef(value *QName) {
	t.EvaluatesToTypeRefField = value
}

type Gateway struct {
	FlowNode
	GatewayDirectionField *GatewayDirection `xml:"gatewayDirection,attr,omitempty"`
	TextPayloadField      *Payload          `xml:",chardata"`
}

var defaultGatewayGatewayDirectionField GatewayDirection = "Unspecified"

func DefaultGateway() Gateway {
	return Gateway{
		FlowNode:              DefaultFlowNode(),
		GatewayDirectionField: &defaultGatewayGatewayDirectionField,
	}
}

type GatewayInterface interface {
	Element
	FlowNodeInterface
	GatewayDirection() (result *GatewayDirection)
	SetGatewayDirection(value *GatewayDirection)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Gateway) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Gateway) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Gateway) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.FlowNode.FindBy(f); found {
		return
	}

	return
}
func (t *Gateway) GatewayDirection() (result *GatewayDirection) {
	if t.GatewayDirectionField == nil {
		result = &defaultGatewayGatewayDirectionField
		return
	}
	result = t.GatewayDirectionField
	return
}
func (t *Gateway) SetGatewayDirection(value *GatewayDirection) {
	t.GatewayDirectionField = value
}

type GlobalBusinessRuleTask struct {
	GlobalTask
	ImplementationField *Implementation `xml:"implementation,attr,omitempty"`
	TextPayloadField    *Payload        `xml:",chardata"`
}

var defaultGlobalBusinessRuleTaskImplementationField Implementation = "##unspecified"

func DefaultGlobalBusinessRuleTask() GlobalBusinessRuleTask {
	return GlobalBusinessRuleTask{
		GlobalTask: DefaultGlobalTask(),
	}
}

type GlobalBusinessRuleTaskInterface interface {
	Element
	GlobalTaskInterface
	Implementation() (result *Implementation)
	SetImplementation(value *Implementation)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *GlobalBusinessRuleTask) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *GlobalBusinessRuleTask) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *GlobalBusinessRuleTask) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.GlobalTask.FindBy(f); found {
		return
	}

	return
}
func (t *GlobalBusinessRuleTask) Implementation() (result *Implementation) {
	if t.ImplementationField == nil {
		result = &defaultGlobalBusinessRuleTaskImplementationField
		return
	}
	result = t.ImplementationField
	return
}
func (t *GlobalBusinessRuleTask) SetImplementation(value *Implementation) {
	t.ImplementationField = value
}

type GlobalChoreographyTask struct {
	Choreography
	InitiatingParticipantRefField *QName   `xml:"initiatingParticipantRef,attr,omitempty"`
	TextPayloadField              *Payload `xml:",chardata"`
}

func DefaultGlobalChoreographyTask() GlobalChoreographyTask {
	return GlobalChoreographyTask{
		Choreography: DefaultChoreography(),
	}
}

type GlobalChoreographyTaskInterface interface {
	Element
	ChoreographyInterface
	InitiatingParticipantRef() (result *QName, present bool)
	SetInitiatingParticipantRef(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *GlobalChoreographyTask) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *GlobalChoreographyTask) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *GlobalChoreographyTask) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Choreography.FindBy(f); found {
		return
	}

	return
}
func (t *GlobalChoreographyTask) InitiatingParticipantRef() (result *QName, present bool) {
	if t.InitiatingParticipantRefField != nil {
		present = true
	}
	result = t.InitiatingParticipantRefField
	return
}
func (t *GlobalChoreographyTask) SetInitiatingParticipantRef(value *QName) {
	t.InitiatingParticipantRefField = value
}

type GlobalConversation struct {
	Collaboration
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultGlobalConversation() GlobalConversation {
	return GlobalConversation{
		Collaboration: DefaultCollaboration(),
	}
}

type GlobalConversationInterface interface {
	Element
	CollaborationInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *GlobalConversation) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *GlobalConversation) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *GlobalConversation) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Collaboration.FindBy(f); found {
		return
	}

	return
}

type GlobalManualTask struct {
	GlobalTask
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultGlobalManualTask() GlobalManualTask {
	return GlobalManualTask{
		GlobalTask: DefaultGlobalTask(),
	}
}

type GlobalManualTaskInterface interface {
	Element
	GlobalTaskInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *GlobalManualTask) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *GlobalManualTask) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *GlobalManualTask) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.GlobalTask.FindBy(f); found {
		return
	}

	return
}

type GlobalScriptTask struct {
	GlobalTask
	ScriptLanguageField *AnyURI  `xml:"scriptLanguage,attr,omitempty"`
	ScriptField         *Script  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL script"`
	TextPayloadField    *Payload `xml:",chardata"`
}

func DefaultGlobalScriptTask() GlobalScriptTask {
	return GlobalScriptTask{
		GlobalTask: DefaultGlobalTask(),
	}
}

type GlobalScriptTaskInterface interface {
	Element
	GlobalTaskInterface
	ScriptLanguage() (result *AnyURI, present bool)
	Script() (result *Script, present bool)
	SetScriptLanguage(value *AnyURI)
	SetScript(value *Script)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *GlobalScriptTask) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *GlobalScriptTask) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *GlobalScriptTask) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.GlobalTask.FindBy(f); found {
		return
	}

	if value := t.ScriptField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	return
}
func (t *GlobalScriptTask) ScriptLanguage() (result *AnyURI, present bool) {
	if t.ScriptLanguageField != nil {
		present = true
	}
	result = t.ScriptLanguageField
	return
}
func (t *GlobalScriptTask) SetScriptLanguage(value *AnyURI) {
	t.ScriptLanguageField = value
}
func (t *GlobalScriptTask) Script() (result *Script, present bool) {
	if t.ScriptField != nil {
		present = true
	}
	result = t.ScriptField
	return
}
func (t *GlobalScriptTask) SetScript(value *Script) {
	t.ScriptField = value
}

type GlobalTask struct {
	CallableElement
	ResourceRoleField []ResourceRole `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL resourceRole"`
	TextPayloadField  *Payload       `xml:",chardata"`
}

func DefaultGlobalTask() GlobalTask {
	return GlobalTask{
		CallableElement: DefaultCallableElement(),
	}
}

type GlobalTaskInterface interface {
	Element
	CallableElementInterface
	ResourceRoles() (result *[]ResourceRole)
	SetResourceRoles(value []ResourceRole)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *GlobalTask) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *GlobalTask) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *GlobalTask) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.CallableElement.FindBy(f); found {
		return
	}

	for i := range t.ResourceRoleField {
		if result, found = t.ResourceRoleField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *GlobalTask) ResourceRoles() (result *[]ResourceRole) {
	result = &t.ResourceRoleField
	return
}
func (t *GlobalTask) SetResourceRoles(value []ResourceRole) {
	t.ResourceRoleField = value
}

type GlobalUserTask struct {
	GlobalTask
	ImplementationField *Implementation `xml:"implementation,attr,omitempty"`
	RenderingField      []Rendering     `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL rendering"`
	TextPayloadField    *Payload        `xml:",chardata"`
}

var defaultGlobalUserTaskImplementationField Implementation = "##unspecified"

func DefaultGlobalUserTask() GlobalUserTask {
	return GlobalUserTask{
		GlobalTask: DefaultGlobalTask(),
	}
}

type GlobalUserTaskInterface interface {
	Element
	GlobalTaskInterface
	Implementation() (result *Implementation)
	Renderings() (result *[]Rendering)
	SetImplementation(value *Implementation)
	SetRenderings(value []Rendering)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *GlobalUserTask) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *GlobalUserTask) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *GlobalUserTask) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.GlobalTask.FindBy(f); found {
		return
	}

	for i := range t.RenderingField {
		if result, found = t.RenderingField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *GlobalUserTask) Implementation() (result *Implementation) {
	if t.ImplementationField == nil {
		result = &defaultGlobalUserTaskImplementationField
		return
	}
	result = t.ImplementationField
	return
}
func (t *GlobalUserTask) SetImplementation(value *Implementation) {
	t.ImplementationField = value
}
func (t *GlobalUserTask) Renderings() (result *[]Rendering) {
	result = &t.RenderingField
	return
}
func (t *GlobalUserTask) SetRenderings(value []Rendering) {
	t.RenderingField = value
}

type Group struct {
	Artifact
	CategoryValueRefField *QName   `xml:"categoryValueRef,attr,omitempty"`
	TextPayloadField      *Payload `xml:",chardata"`
}

func DefaultGroup() Group {
	return Group{
		Artifact: DefaultArtifact(),
	}
}

type GroupInterface interface {
	Element
	ArtifactInterface
	CategoryValueRef() (result *QName, present bool)
	SetCategoryValueRef(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Group) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Group) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Group) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Artifact.FindBy(f); found {
		return
	}

	return
}
func (t *Group) CategoryValueRef() (result *QName, present bool) {
	if t.CategoryValueRefField != nil {
		present = true
	}
	result = t.CategoryValueRefField
	return
}
func (t *Group) SetCategoryValueRef(value *QName) {
	t.CategoryValueRefField = value
}

type HumanPerformer struct {
	Performer
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultHumanPerformer() HumanPerformer {
	return HumanPerformer{
		Performer: DefaultPerformer(),
	}
}

type HumanPerformerInterface interface {
	Element
	PerformerInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *HumanPerformer) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *HumanPerformer) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *HumanPerformer) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Performer.FindBy(f); found {
		return
	}

	return
}

type ImplicitThrowEvent struct {
	ThrowEvent
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultImplicitThrowEvent() ImplicitThrowEvent {
	return ImplicitThrowEvent{
		ThrowEvent: DefaultThrowEvent(),
	}
}

type ImplicitThrowEventInterface interface {
	Element
	ThrowEventInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ImplicitThrowEvent) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ImplicitThrowEvent) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ImplicitThrowEvent) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.ThrowEvent.FindBy(f); found {
		return
	}

	return
}

type InclusiveGateway struct {
	Gateway
	DefaultField     *IdRef   `xml:"default,attr,omitempty"`
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultInclusiveGateway() InclusiveGateway {
	return InclusiveGateway{
		Gateway: DefaultGateway(),
	}
}

type InclusiveGatewayInterface interface {
	Element
	GatewayInterface
	Default() (result *IdRef, present bool)
	SetDefault(value *IdRef)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *InclusiveGateway) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *InclusiveGateway) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *InclusiveGateway) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Gateway.FindBy(f); found {
		return
	}

	return
}
func (t *InclusiveGateway) Default() (result *IdRef, present bool) {
	if t.DefaultField != nil {
		present = true
	}
	result = t.DefaultField
	return
}
func (t *InclusiveGateway) SetDefault(value *IdRef) {
	t.DefaultField = value
}

type InputSet struct {
	BaseElement
	NameField                    *string  `xml:"name,attr,omitempty"`
	DataInputRefsField           []IdRef  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataInputRefs"`
	OptionalInputRefsField       []IdRef  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL optionalInputRefs"`
	WhileExecutingInputRefsField []IdRef  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL whileExecutingInputRefs"`
	OutputSetRefsField           []IdRef  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL outputSetRefs"`
	TextPayloadField             *Payload `xml:",chardata"`
}

func DefaultInputSet() InputSet {
	return InputSet{
		BaseElement: DefaultBaseElement(),
	}
}

type InputSetInterface interface {
	Element
	BaseElementInterface
	Name() (result *string, present bool)
	DataInputRefses() (result *[]IdRef)
	OptionalInputRefses() (result *[]IdRef)
	WhileExecutingInputRefses() (result *[]IdRef)
	OutputSetRefses() (result *[]IdRef)
	SetName(value *string)
	SetDataInputRefses(value []IdRef)
	SetOptionalInputRefses(value []IdRef)
	SetWhileExecutingInputRefses(value []IdRef)
	SetOutputSetRefses(value []IdRef)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *InputSet) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *InputSet) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *InputSet) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}
func (t *InputSet) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *InputSet) SetName(value *string) {
	t.NameField = value
}
func (t *InputSet) DataInputRefses() (result *[]IdRef) {
	result = &t.DataInputRefsField
	return
}
func (t *InputSet) SetDataInputRefses(value []IdRef) {
	t.DataInputRefsField = value
}
func (t *InputSet) OptionalInputRefses() (result *[]IdRef) {
	result = &t.OptionalInputRefsField
	return
}
func (t *InputSet) SetOptionalInputRefses(value []IdRef) {
	t.OptionalInputRefsField = value
}
func (t *InputSet) WhileExecutingInputRefses() (result *[]IdRef) {
	result = &t.WhileExecutingInputRefsField
	return
}
func (t *InputSet) SetWhileExecutingInputRefses(value []IdRef) {
	t.WhileExecutingInputRefsField = value
}
func (t *InputSet) OutputSetRefses() (result *[]IdRef) {
	result = &t.OutputSetRefsField
	return
}
func (t *InputSet) SetOutputSetRefses(value []IdRef) {
	t.OutputSetRefsField = value
}

type Interface struct {
	RootElement
	NameField              string      `xml:"name,attr,omitempty"`
	ImplementationRefField *QName      `xml:"implementationRef,attr,omitempty"`
	OperationField         []Operation `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL operation"`
	TextPayloadField       *Payload    `xml:",chardata"`
}

func DefaultInterface() Interface {
	return Interface{
		RootElement: DefaultRootElement(),
	}
}

type InterfaceInterface interface {
	Element
	RootElementInterface
	Name() (result *string)
	ImplementationRef() (result *QName, present bool)
	Operations() (result *[]Operation)
	SetName(value string)
	SetImplementationRef(value *QName)
	SetOperations(value []Operation)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Interface) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Interface) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Interface) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.RootElement.FindBy(f); found {
		return
	}

	for i := range t.OperationField {
		if result, found = t.OperationField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *Interface) Name() (result *string) {
	result = &t.NameField
	return
}
func (t *Interface) SetName(value string) {
	t.NameField = value
}
func (t *Interface) ImplementationRef() (result *QName, present bool) {
	if t.ImplementationRefField != nil {
		present = true
	}
	result = t.ImplementationRefField
	return
}
func (t *Interface) SetImplementationRef(value *QName) {
	t.ImplementationRefField = value
}
func (t *Interface) Operations() (result *[]Operation) {
	result = &t.OperationField
	return
}
func (t *Interface) SetOperations(value []Operation) {
	t.OperationField = value
}

type IntermediateCatchEvent struct {
	CatchEvent
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultIntermediateCatchEvent() IntermediateCatchEvent {
	return IntermediateCatchEvent{
		CatchEvent: DefaultCatchEvent(),
	}
}

type IntermediateCatchEventInterface interface {
	Element
	CatchEventInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *IntermediateCatchEvent) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *IntermediateCatchEvent) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *IntermediateCatchEvent) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.CatchEvent.FindBy(f); found {
		return
	}

	return
}

type IntermediateThrowEvent struct {
	ThrowEvent
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultIntermediateThrowEvent() IntermediateThrowEvent {
	return IntermediateThrowEvent{
		ThrowEvent: DefaultThrowEvent(),
	}
}

type IntermediateThrowEventInterface interface {
	Element
	ThrowEventInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *IntermediateThrowEvent) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *IntermediateThrowEvent) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *IntermediateThrowEvent) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.ThrowEvent.FindBy(f); found {
		return
	}

	return
}

type InputOutputBinding struct {
	BaseElement
	OperationRefField  QName    `xml:"operationRef,attr,omitempty"`
	InputDataRefField  IdRef    `xml:"inputDataRef,attr,omitempty"`
	OutputDataRefField IdRef    `xml:"outputDataRef,attr,omitempty"`
	TextPayloadField   *Payload `xml:",chardata"`
}

func DefaultInputOutputBinding() InputOutputBinding {
	return InputOutputBinding{
		BaseElement: DefaultBaseElement(),
	}
}

type InputOutputBindingInterface interface {
	Element
	BaseElementInterface
	OperationRef() (result *QName)
	InputDataRef() (result *IdRef)
	OutputDataRef() (result *IdRef)
	SetOperationRef(value QName)
	SetInputDataRef(value IdRef)
	SetOutputDataRef(value IdRef)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *InputOutputBinding) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *InputOutputBinding) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *InputOutputBinding) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}
func (t *InputOutputBinding) OperationRef() (result *QName) {
	result = &t.OperationRefField
	return
}
func (t *InputOutputBinding) SetOperationRef(value QName) {
	t.OperationRefField = value
}
func (t *InputOutputBinding) InputDataRef() (result *IdRef) {
	result = &t.InputDataRefField
	return
}
func (t *InputOutputBinding) SetInputDataRef(value IdRef) {
	t.InputDataRefField = value
}
func (t *InputOutputBinding) OutputDataRef() (result *IdRef) {
	result = &t.OutputDataRefField
	return
}
func (t *InputOutputBinding) SetOutputDataRef(value IdRef) {
	t.OutputDataRefField = value
}

type InputOutputSpecification struct {
	BaseElement
	DataInputField   []DataInput  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataInput"`
	DataOutputField  []DataOutput `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataOutput"`
	InputSetField    []InputSet   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL inputSet"`
	OutputSetField   []OutputSet  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL outputSet"`
	TextPayloadField *Payload     `xml:",chardata"`
}

func DefaultInputOutputSpecification() InputOutputSpecification {
	return InputOutputSpecification{
		BaseElement: DefaultBaseElement(),
	}
}

type InputOutputSpecificationInterface interface {
	Element
	BaseElementInterface
	DataInputs() (result *[]DataInput)
	DataOutputs() (result *[]DataOutput)
	InputSets() (result *[]InputSet)
	OutputSets() (result *[]OutputSet)
	SetDataInputs(value []DataInput)
	SetDataOutputs(value []DataOutput)
	SetInputSets(value []InputSet)
	SetOutputSets(value []OutputSet)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *InputOutputSpecification) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *InputOutputSpecification) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *InputOutputSpecification) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	for i := range t.DataInputField {
		if result, found = t.DataInputField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataOutputField {
		if result, found = t.DataOutputField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.InputSetField {
		if result, found = t.InputSetField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.OutputSetField {
		if result, found = t.OutputSetField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *InputOutputSpecification) DataInputs() (result *[]DataInput) {
	result = &t.DataInputField
	return
}
func (t *InputOutputSpecification) SetDataInputs(value []DataInput) {
	t.DataInputField = value
}
func (t *InputOutputSpecification) DataOutputs() (result *[]DataOutput) {
	result = &t.DataOutputField
	return
}
func (t *InputOutputSpecification) SetDataOutputs(value []DataOutput) {
	t.DataOutputField = value
}
func (t *InputOutputSpecification) InputSets() (result *[]InputSet) {
	result = &t.InputSetField
	return
}
func (t *InputOutputSpecification) SetInputSets(value []InputSet) {
	t.InputSetField = value
}
func (t *InputOutputSpecification) OutputSets() (result *[]OutputSet) {
	result = &t.OutputSetField
	return
}
func (t *InputOutputSpecification) SetOutputSets(value []OutputSet) {
	t.OutputSetField = value
}

type ItemAware struct {
	ItemSubjectRefField *QName     `xml:"itemSubjectRef,attr,omitempty"`
	DataStateField      *DataState `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataState"`
}

func DefaultItemAware() ItemAware {
	return ItemAware{}
}

type ItemAwareInterface interface {
	Element
	ItemSubjectRef() (result *QName, present bool)
	DataState() (result *DataState, present bool)
	SetItemSubjectRef(value *QName)
	SetDataState(value *DataState)
}

func (t *ItemAware) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}

	if value := t.DataStateField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	return
}
func (t *ItemAware) ItemSubjectRef() (result *QName, present bool) {
	if t.ItemSubjectRefField != nil {
		present = true
	}
	result = t.ItemSubjectRefField
	return
}
func (t *ItemAware) SetItemSubjectRef(value *QName) {
	t.ItemSubjectRefField = value
}
func (t *ItemAware) DataState() (result *DataState, present bool) {
	if t.DataStateField != nil {
		present = true
	}
	result = t.DataStateField
	return
}
func (t *ItemAware) SetDataState(value *DataState) {
	t.DataStateField = value
}

type ItemAwareElement struct {
	BaseElement
	ItemAware
}

func DefaultItemAwareElement() ItemAwareElement {
	return ItemAwareElement{
		BaseElement: DefaultBaseElement(),
		ItemAware:   DefaultItemAware(),
	}
}

type ItemAwareElementInterface interface {
	Element
	BaseElementInterface
	ItemAwareInterface
}

func (t *ItemAwareElement) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}
	if result, found = t.ItemAware.FindBy(f); found {
		return
	}

	return
}

type ItemDefinition struct {
	RootElement
	StructureRefField *QName    `xml:"structureRef,attr,omitempty"`
	IsCollectionField *bool     `xml:"isCollection,attr,omitempty"`
	ItemKindField     *ItemKind `xml:"itemKind,attr,omitempty"`
	TextPayloadField  *Payload  `xml:",chardata"`
}

var defaultItemDefinitionIsCollectionField bool = false
var defaultItemDefinitionItemKindField ItemKind = "Information"

func DefaultItemDefinition() ItemDefinition {
	return ItemDefinition{
		RootElement:       DefaultRootElement(),
		IsCollectionField: &defaultItemDefinitionIsCollectionField,
		ItemKindField:     &defaultItemDefinitionItemKindField,
	}
}

type ItemDefinitionInterface interface {
	Element
	RootElementInterface
	StructureRef() (result *QName, present bool)
	IsCollection() (result bool)
	ItemKind() (result *ItemKind)
	SetStructureRef(value *QName)
	SetIsCollection(value *bool)
	SetItemKind(value *ItemKind)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ItemDefinition) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ItemDefinition) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ItemDefinition) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.RootElement.FindBy(f); found {
		return
	}

	return
}
func (t *ItemDefinition) StructureRef() (result *QName, present bool) {
	if t.StructureRefField != nil {
		present = true
	}
	result = t.StructureRefField
	return
}
func (t *ItemDefinition) SetStructureRef(value *QName) {
	t.StructureRefField = value
}
func (t *ItemDefinition) IsCollection() (result bool) {
	if t.IsCollectionField == nil {
		result = defaultItemDefinitionIsCollectionField
		return
	}
	result = *t.IsCollectionField
	return
}
func (t *ItemDefinition) SetIsCollection(value *bool) {
	t.IsCollectionField = value
}
func (t *ItemDefinition) ItemKind() (result *ItemKind) {
	if t.ItemKindField == nil {
		result = &defaultItemDefinitionItemKindField
		return
	}
	result = t.ItemKindField
	return
}
func (t *ItemDefinition) SetItemKind(value *ItemKind) {
	t.ItemKindField = value
}

type Lane struct {
	BaseElement
	NameField                *string      `xml:"name,attr,omitempty"`
	PartitionElementRefField *QName       `xml:"partitionElementRef,attr,omitempty"`
	PartitionElementField    *BaseElement `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL partitionElement"`
	FlowNodeRefField         []IdRef      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL flowNodeRef"`
	ChildLaneSetField        *LaneSet     `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL childLaneSet"`
	TextPayloadField         *Payload     `xml:",chardata"`
}

func DefaultLane() Lane {
	return Lane{
		BaseElement: DefaultBaseElement(),
	}
}

type LaneInterface interface {
	Element
	BaseElementInterface
	Name() (result *string, present bool)
	PartitionElementRef() (result *QName, present bool)
	PartitionElement() (result *BaseElement, present bool)
	FlowNodeRefs() (result *[]IdRef)
	ChildLaneSet() (result *LaneSet, present bool)
	SetName(value *string)
	SetPartitionElementRef(value *QName)
	SetPartitionElement(value *BaseElement)
	SetFlowNodeRefs(value []IdRef)
	SetChildLaneSet(value *LaneSet)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Lane) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Lane) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Lane) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	if value := t.PartitionElementField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	if value := t.ChildLaneSetField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	return
}
func (t *Lane) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *Lane) SetName(value *string) {
	t.NameField = value
}
func (t *Lane) PartitionElementRef() (result *QName, present bool) {
	if t.PartitionElementRefField != nil {
		present = true
	}
	result = t.PartitionElementRefField
	return
}
func (t *Lane) SetPartitionElementRef(value *QName) {
	t.PartitionElementRefField = value
}
func (t *Lane) PartitionElement() (result *BaseElement, present bool) {
	if t.PartitionElementField != nil {
		present = true
	}
	result = t.PartitionElementField
	return
}
func (t *Lane) SetPartitionElement(value *BaseElement) {
	t.PartitionElementField = value
}
func (t *Lane) FlowNodeRefs() (result *[]IdRef) {
	result = &t.FlowNodeRefField
	return
}
func (t *Lane) SetFlowNodeRefs(value []IdRef) {
	t.FlowNodeRefField = value
}
func (t *Lane) ChildLaneSet() (result *LaneSet, present bool) {
	if t.ChildLaneSetField != nil {
		present = true
	}
	result = t.ChildLaneSetField
	return
}
func (t *Lane) SetChildLaneSet(value *LaneSet) {
	t.ChildLaneSetField = value
}

type LaneSet struct {
	BaseElement
	NameField        *string  `xml:"name,attr,omitempty"`
	LaneField        []Lane   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL lane"`
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultLaneSet() LaneSet {
	return LaneSet{
		BaseElement: DefaultBaseElement(),
	}
}

type LaneSetInterface interface {
	Element
	BaseElementInterface
	Name() (result *string, present bool)
	Lanes() (result *[]Lane)
	SetName(value *string)
	SetLanes(value []Lane)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *LaneSet) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *LaneSet) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *LaneSet) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	for i := range t.LaneField {
		if result, found = t.LaneField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *LaneSet) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *LaneSet) SetName(value *string) {
	t.NameField = value
}
func (t *LaneSet) Lanes() (result *[]Lane) {
	result = &t.LaneField
	return
}
func (t *LaneSet) SetLanes(value []Lane) {
	t.LaneField = value
}

type LinkEventDefinition struct {
	EventDefinition
	NameField        string   `xml:"name,attr,omitempty"`
	SourceField      []QName  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL source"`
	TargetField      *QName   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL target"`
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultLinkEventDefinition() LinkEventDefinition {
	return LinkEventDefinition{
		EventDefinition: DefaultEventDefinition(),
	}
}

type LinkEventDefinitionInterface interface {
	Element
	EventDefinitionInterface
	Name() (result *string)
	Sources() (result *[]QName)
	Target() (result *QName, present bool)
	SetName(value string)
	SetSources(value []QName)
	SetTarget(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *LinkEventDefinition) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *LinkEventDefinition) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *LinkEventDefinition) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.EventDefinition.FindBy(f); found {
		return
	}

	return
}
func (t *LinkEventDefinition) Name() (result *string) {
	result = &t.NameField
	return
}
func (t *LinkEventDefinition) SetName(value string) {
	t.NameField = value
}
func (t *LinkEventDefinition) Sources() (result *[]QName) {
	result = &t.SourceField
	return
}
func (t *LinkEventDefinition) SetSources(value []QName) {
	t.SourceField = value
}
func (t *LinkEventDefinition) Target() (result *QName, present bool) {
	if t.TargetField != nil {
		present = true
	}
	result = t.TargetField
	return
}
func (t *LinkEventDefinition) SetTarget(value *QName) {
	t.TargetField = value
}

type LoopCharacteristics struct {
	BaseElement
}

func DefaultLoopCharacteristics() LoopCharacteristics {
	return LoopCharacteristics{
		BaseElement: DefaultBaseElement(),
	}
}

type LoopCharacteristicsInterface interface {
	Element
	BaseElementInterface
}

func (t *LoopCharacteristics) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}

type ManualTask struct {
	Task
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultManualTask() ManualTask {
	return ManualTask{
		Task: DefaultTask(),
	}
}

type ManualTaskInterface interface {
	Element
	TaskInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ManualTask) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ManualTask) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ManualTask) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Task.FindBy(f); found {
		return
	}

	return
}

type Message struct {
	RootElement
	NameField        *string  `xml:"name,attr,omitempty"`
	ItemRefField     *QName   `xml:"itemRef,attr,omitempty"`
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultMessage() Message {
	return Message{
		RootElement: DefaultRootElement(),
	}
}

type MessageInterface interface {
	Element
	RootElementInterface
	Name() (result *string, present bool)
	ItemRef() (result *QName, present bool)
	SetName(value *string)
	SetItemRef(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Message) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Message) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Message) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.RootElement.FindBy(f); found {
		return
	}

	return
}
func (t *Message) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *Message) SetName(value *string) {
	t.NameField = value
}
func (t *Message) ItemRef() (result *QName, present bool) {
	if t.ItemRefField != nil {
		present = true
	}
	result = t.ItemRefField
	return
}
func (t *Message) SetItemRef(value *QName) {
	t.ItemRefField = value
}

type MessageEventDefinition struct {
	EventDefinition
	MessageRefField   *QName   `xml:"messageRef,attr,omitempty"`
	OperationRefField *QName   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL operationRef"`
	TextPayloadField  *Payload `xml:",chardata"`
}

func DefaultMessageEventDefinition() MessageEventDefinition {
	return MessageEventDefinition{
		EventDefinition: DefaultEventDefinition(),
	}
}

type MessageEventDefinitionInterface interface {
	Element
	EventDefinitionInterface
	MessageRef() (result *QName, present bool)
	OperationRef() (result *QName, present bool)
	SetMessageRef(value *QName)
	SetOperationRef(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *MessageEventDefinition) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *MessageEventDefinition) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *MessageEventDefinition) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.EventDefinition.FindBy(f); found {
		return
	}

	return
}
func (t *MessageEventDefinition) MessageRef() (result *QName, present bool) {
	if t.MessageRefField != nil {
		present = true
	}
	result = t.MessageRefField
	return
}
func (t *MessageEventDefinition) SetMessageRef(value *QName) {
	t.MessageRefField = value
}
func (t *MessageEventDefinition) OperationRef() (result *QName, present bool) {
	if t.OperationRefField != nil {
		present = true
	}
	result = t.OperationRefField
	return
}
func (t *MessageEventDefinition) SetOperationRef(value *QName) {
	t.OperationRefField = value
}

type MessageFlow struct {
	BaseElement
	NameField        *string  `xml:"name,attr,omitempty"`
	SourceRefField   QName    `xml:"sourceRef,attr,omitempty"`
	TargetRefField   QName    `xml:"targetRef,attr,omitempty"`
	MessageRefField  *QName   `xml:"messageRef,attr,omitempty"`
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultMessageFlow() MessageFlow {
	return MessageFlow{
		BaseElement: DefaultBaseElement(),
	}
}

type MessageFlowInterface interface {
	Element
	BaseElementInterface
	Name() (result *string, present bool)
	SourceRef() (result *QName)
	TargetRef() (result *QName)
	MessageRef() (result *QName, present bool)
	SetName(value *string)
	SetSourceRef(value QName)
	SetTargetRef(value QName)
	SetMessageRef(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *MessageFlow) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *MessageFlow) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *MessageFlow) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}
func (t *MessageFlow) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *MessageFlow) SetName(value *string) {
	t.NameField = value
}
func (t *MessageFlow) SourceRef() (result *QName) {
	result = &t.SourceRefField
	return
}
func (t *MessageFlow) SetSourceRef(value QName) {
	t.SourceRefField = value
}
func (t *MessageFlow) TargetRef() (result *QName) {
	result = &t.TargetRefField
	return
}
func (t *MessageFlow) SetTargetRef(value QName) {
	t.TargetRefField = value
}
func (t *MessageFlow) MessageRef() (result *QName, present bool) {
	if t.MessageRefField != nil {
		present = true
	}
	result = t.MessageRefField
	return
}
func (t *MessageFlow) SetMessageRef(value *QName) {
	t.MessageRefField = value
}

type MessageFlowAssociation struct {
	BaseElement
	InnerMessageFlowRefField QName    `xml:"innerMessageFlowRef,attr,omitempty"`
	OuterMessageFlowRefField QName    `xml:"outerMessageFlowRef,attr,omitempty"`
	TextPayloadField         *Payload `xml:",chardata"`
}

func DefaultMessageFlowAssociation() MessageFlowAssociation {
	return MessageFlowAssociation{
		BaseElement: DefaultBaseElement(),
	}
}

type MessageFlowAssociationInterface interface {
	Element
	BaseElementInterface
	InnerMessageFlowRef() (result *QName)
	OuterMessageFlowRef() (result *QName)
	SetInnerMessageFlowRef(value QName)
	SetOuterMessageFlowRef(value QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *MessageFlowAssociation) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *MessageFlowAssociation) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *MessageFlowAssociation) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}
func (t *MessageFlowAssociation) InnerMessageFlowRef() (result *QName) {
	result = &t.InnerMessageFlowRefField
	return
}
func (t *MessageFlowAssociation) SetInnerMessageFlowRef(value QName) {
	t.InnerMessageFlowRefField = value
}
func (t *MessageFlowAssociation) OuterMessageFlowRef() (result *QName) {
	result = &t.OuterMessageFlowRefField
	return
}
func (t *MessageFlowAssociation) SetOuterMessageFlowRef(value QName) {
	t.OuterMessageFlowRefField = value
}

type Monitoring struct {
	BaseElement
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultMonitoring() Monitoring {
	return Monitoring{
		BaseElement: DefaultBaseElement(),
	}
}

type MonitoringInterface interface {
	Element
	BaseElementInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Monitoring) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Monitoring) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Monitoring) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}

type MultiInstanceLoopCharacteristics struct {
	LoopCharacteristics
	IsSequentialField              *bool                       `xml:"isSequential,attr,omitempty"`
	BehaviorField                  *MultiInstanceFlowCondition `xml:"behavior,attr,omitempty"`
	OneBehaviorEventRefField       *QName                      `xml:"oneBehaviorEventRef,attr,omitempty"`
	NoneBehaviorEventRefField      *QName                      `xml:"noneBehaviorEventRef,attr,omitempty"`
	LoopCardinalityField           *AnExpression               `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL loopCardinality"`
	LoopDataInputRefField          *QName                      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL loopDataInputRef"`
	LoopDataOutputRefField         *QName                      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL loopDataOutputRef"`
	InputDataItemField             *DataInput                  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL inputDataItem"`
	OutputDataItemField            *DataOutput                 `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL outputDataItem"`
	ComplexBehaviorDefinitionField []ComplexBehaviorDefinition `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL complexBehaviorDefinition"`
	CompletionConditionField       *AnExpression               `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL completionCondition"`
	TextPayloadField               *Payload                    `xml:",chardata"`
}

var defaultMultiInstanceLoopCharacteristicsIsSequentialField bool = false
var defaultMultiInstanceLoopCharacteristicsBehaviorField MultiInstanceFlowCondition = "All"

func DefaultMultiInstanceLoopCharacteristics() MultiInstanceLoopCharacteristics {
	return MultiInstanceLoopCharacteristics{
		LoopCharacteristics: DefaultLoopCharacteristics(),
		IsSequentialField:   &defaultMultiInstanceLoopCharacteristicsIsSequentialField,
		BehaviorField:       &defaultMultiInstanceLoopCharacteristicsBehaviorField,
	}
}

type MultiInstanceLoopCharacteristicsInterface interface {
	Element
	LoopCharacteristicsInterface
	IsSequential() (result bool)
	Behavior() (result *MultiInstanceFlowCondition)
	OneBehaviorEventRef() (result *QName, present bool)
	NoneBehaviorEventRef() (result *QName, present bool)
	LoopCardinality() (result *AnExpression, present bool)
	LoopDataInputRef() (result *QName, present bool)
	LoopDataOutputRef() (result *QName, present bool)
	InputDataItem() (result *DataInput, present bool)
	OutputDataItem() (result *DataOutput, present bool)
	ComplexBehaviorDefinitions() (result *[]ComplexBehaviorDefinition)
	CompletionCondition() (result *AnExpression, present bool)
	SetIsSequential(value *bool)
	SetBehavior(value *MultiInstanceFlowCondition)
	SetOneBehaviorEventRef(value *QName)
	SetNoneBehaviorEventRef(value *QName)
	SetLoopCardinality(value *AnExpression)
	SetLoopDataInputRef(value *QName)
	SetLoopDataOutputRef(value *QName)
	SetInputDataItem(value *DataInput)
	SetOutputDataItem(value *DataOutput)
	SetComplexBehaviorDefinitions(value []ComplexBehaviorDefinition)
	SetCompletionCondition(value *AnExpression)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *MultiInstanceLoopCharacteristics) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *MultiInstanceLoopCharacteristics) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *MultiInstanceLoopCharacteristics) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.LoopCharacteristics.FindBy(f); found {
		return
	}

	if value := t.LoopCardinalityField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	if value := t.InputDataItemField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	if value := t.OutputDataItemField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	for i := range t.ComplexBehaviorDefinitionField {
		if result, found = t.ComplexBehaviorDefinitionField[i].FindBy(f); found {
			return
		}
	}

	if value := t.CompletionConditionField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	return
}
func (t *MultiInstanceLoopCharacteristics) IsSequential() (result bool) {
	if t.IsSequentialField == nil {
		result = defaultMultiInstanceLoopCharacteristicsIsSequentialField
		return
	}
	result = *t.IsSequentialField
	return
}
func (t *MultiInstanceLoopCharacteristics) SetIsSequential(value *bool) {
	t.IsSequentialField = value
}
func (t *MultiInstanceLoopCharacteristics) Behavior() (result *MultiInstanceFlowCondition) {
	if t.BehaviorField == nil {
		result = &defaultMultiInstanceLoopCharacteristicsBehaviorField
		return
	}
	result = t.BehaviorField
	return
}
func (t *MultiInstanceLoopCharacteristics) SetBehavior(value *MultiInstanceFlowCondition) {
	t.BehaviorField = value
}
func (t *MultiInstanceLoopCharacteristics) OneBehaviorEventRef() (result *QName, present bool) {
	if t.OneBehaviorEventRefField != nil {
		present = true
	}
	result = t.OneBehaviorEventRefField
	return
}
func (t *MultiInstanceLoopCharacteristics) SetOneBehaviorEventRef(value *QName) {
	t.OneBehaviorEventRefField = value
}
func (t *MultiInstanceLoopCharacteristics) NoneBehaviorEventRef() (result *QName, present bool) {
	if t.NoneBehaviorEventRefField != nil {
		present = true
	}
	result = t.NoneBehaviorEventRefField
	return
}
func (t *MultiInstanceLoopCharacteristics) SetNoneBehaviorEventRef(value *QName) {
	t.NoneBehaviorEventRefField = value
}
func (t *MultiInstanceLoopCharacteristics) LoopCardinality() (result *AnExpression, present bool) {
	if t.LoopCardinalityField != nil {
		present = true
	}
	result = t.LoopCardinalityField
	return
}
func (t *MultiInstanceLoopCharacteristics) SetLoopCardinality(value *AnExpression) {
	t.LoopCardinalityField = value
}
func (t *MultiInstanceLoopCharacteristics) LoopDataInputRef() (result *QName, present bool) {
	if t.LoopDataInputRefField != nil {
		present = true
	}
	result = t.LoopDataInputRefField
	return
}
func (t *MultiInstanceLoopCharacteristics) SetLoopDataInputRef(value *QName) {
	t.LoopDataInputRefField = value
}
func (t *MultiInstanceLoopCharacteristics) LoopDataOutputRef() (result *QName, present bool) {
	if t.LoopDataOutputRefField != nil {
		present = true
	}
	result = t.LoopDataOutputRefField
	return
}
func (t *MultiInstanceLoopCharacteristics) SetLoopDataOutputRef(value *QName) {
	t.LoopDataOutputRefField = value
}
func (t *MultiInstanceLoopCharacteristics) InputDataItem() (result *DataInput, present bool) {
	if t.InputDataItemField != nil {
		present = true
	}
	result = t.InputDataItemField
	return
}
func (t *MultiInstanceLoopCharacteristics) SetInputDataItem(value *DataInput) {
	t.InputDataItemField = value
}
func (t *MultiInstanceLoopCharacteristics) OutputDataItem() (result *DataOutput, present bool) {
	if t.OutputDataItemField != nil {
		present = true
	}
	result = t.OutputDataItemField
	return
}
func (t *MultiInstanceLoopCharacteristics) SetOutputDataItem(value *DataOutput) {
	t.OutputDataItemField = value
}
func (t *MultiInstanceLoopCharacteristics) ComplexBehaviorDefinitions() (result *[]ComplexBehaviorDefinition) {
	result = &t.ComplexBehaviorDefinitionField
	return
}
func (t *MultiInstanceLoopCharacteristics) SetComplexBehaviorDefinitions(value []ComplexBehaviorDefinition) {
	t.ComplexBehaviorDefinitionField = value
}
func (t *MultiInstanceLoopCharacteristics) CompletionCondition() (result *AnExpression, present bool) {
	if t.CompletionConditionField != nil {
		present = true
	}
	result = t.CompletionConditionField
	return
}
func (t *MultiInstanceLoopCharacteristics) SetCompletionCondition(value *AnExpression) {
	t.CompletionConditionField = value
}

type Operation struct {
	BaseElement
	NameField              string   `xml:"name,attr,omitempty"`
	ImplementationRefField *QName   `xml:"implementationRef,attr,omitempty"`
	InMessageRefField      QName    `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL inMessageRef"`
	OutMessageRefField     *QName   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL outMessageRef"`
	ErrorRefField          []QName  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL errorRef"`
	TextPayloadField       *Payload `xml:",chardata"`
}

func DefaultOperation() Operation {
	return Operation{
		BaseElement: DefaultBaseElement(),
	}
}

type OperationInterface interface {
	Element
	BaseElementInterface
	Name() (result *string)
	ImplementationRef() (result *QName, present bool)
	InMessageRef() (result *QName)
	OutMessageRef() (result *QName, present bool)
	ErrorRefs() (result *[]QName)
	SetName(value string)
	SetImplementationRef(value *QName)
	SetInMessageRef(value QName)
	SetOutMessageRef(value *QName)
	SetErrorRefs(value []QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Operation) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Operation) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Operation) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}
func (t *Operation) Name() (result *string) {
	result = &t.NameField
	return
}
func (t *Operation) SetName(value string) {
	t.NameField = value
}
func (t *Operation) ImplementationRef() (result *QName, present bool) {
	if t.ImplementationRefField != nil {
		present = true
	}
	result = t.ImplementationRefField
	return
}
func (t *Operation) SetImplementationRef(value *QName) {
	t.ImplementationRefField = value
}
func (t *Operation) InMessageRef() (result *QName) {
	result = &t.InMessageRefField
	return
}
func (t *Operation) SetInMessageRef(value QName) {
	t.InMessageRefField = value
}
func (t *Operation) OutMessageRef() (result *QName, present bool) {
	if t.OutMessageRefField != nil {
		present = true
	}
	result = t.OutMessageRefField
	return
}
func (t *Operation) SetOutMessageRef(value *QName) {
	t.OutMessageRefField = value
}
func (t *Operation) ErrorRefs() (result *[]QName) {
	result = &t.ErrorRefField
	return
}
func (t *Operation) SetErrorRefs(value []QName) {
	t.ErrorRefField = value
}

type OutputSet struct {
	BaseElement
	NameField                     *string  `xml:"name,attr,omitempty"`
	DataOutputRefsField           []IdRef  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataOutputRefs"`
	OptionalOutputRefsField       []IdRef  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL optionalOutputRefs"`
	WhileExecutingOutputRefsField []IdRef  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL whileExecutingOutputRefs"`
	InputSetRefsField             []IdRef  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL inputSetRefs"`
	TextPayloadField              *Payload `xml:",chardata"`
}

func DefaultOutputSet() OutputSet {
	return OutputSet{
		BaseElement: DefaultBaseElement(),
	}
}

type OutputSetInterface interface {
	Element
	BaseElementInterface
	Name() (result *string, present bool)
	DataOutputRefses() (result *[]IdRef)
	OptionalOutputRefses() (result *[]IdRef)
	WhileExecutingOutputRefses() (result *[]IdRef)
	InputSetRefses() (result *[]IdRef)
	SetName(value *string)
	SetDataOutputRefses(value []IdRef)
	SetOptionalOutputRefses(value []IdRef)
	SetWhileExecutingOutputRefses(value []IdRef)
	SetInputSetRefses(value []IdRef)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *OutputSet) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *OutputSet) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *OutputSet) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}
func (t *OutputSet) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *OutputSet) SetName(value *string) {
	t.NameField = value
}
func (t *OutputSet) DataOutputRefses() (result *[]IdRef) {
	result = &t.DataOutputRefsField
	return
}
func (t *OutputSet) SetDataOutputRefses(value []IdRef) {
	t.DataOutputRefsField = value
}
func (t *OutputSet) OptionalOutputRefses() (result *[]IdRef) {
	result = &t.OptionalOutputRefsField
	return
}
func (t *OutputSet) SetOptionalOutputRefses(value []IdRef) {
	t.OptionalOutputRefsField = value
}
func (t *OutputSet) WhileExecutingOutputRefses() (result *[]IdRef) {
	result = &t.WhileExecutingOutputRefsField
	return
}
func (t *OutputSet) SetWhileExecutingOutputRefses(value []IdRef) {
	t.WhileExecutingOutputRefsField = value
}
func (t *OutputSet) InputSetRefses() (result *[]IdRef) {
	result = &t.InputSetRefsField
	return
}
func (t *OutputSet) SetInputSetRefses(value []IdRef) {
	t.InputSetRefsField = value
}

type ParallelGateway struct {
	Gateway
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultParallelGateway() ParallelGateway {
	return ParallelGateway{
		Gateway: DefaultGateway(),
	}
}

type ParallelGatewayInterface interface {
	Element
	GatewayInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ParallelGateway) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ParallelGateway) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ParallelGateway) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Gateway.FindBy(f); found {
		return
	}

	return
}

type Participant struct {
	BaseElement
	NameField                    *string                  `xml:"name,attr,omitempty"`
	ProcessRefField              *QName                   `xml:"processRef,attr,omitempty"`
	InterfaceRefField            []QName                  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL interfaceRef"`
	EndPointRefField             []QName                  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL endPointRef"`
	ParticipantMultiplicityField *ParticipantMultiplicity `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL participantMultiplicity"`
	TextPayloadField             *Payload                 `xml:",chardata"`
}

func DefaultParticipant() Participant {
	return Participant{
		BaseElement: DefaultBaseElement(),
	}
}

type ParticipantInterface interface {
	Element
	BaseElementInterface
	Name() (result *string, present bool)
	ProcessRef() (result *QName, present bool)
	InterfaceRefs() (result *[]QName)
	EndPointRefs() (result *[]QName)
	ParticipantMultiplicity() (result *ParticipantMultiplicity, present bool)
	SetName(value *string)
	SetProcessRef(value *QName)
	SetInterfaceRefs(value []QName)
	SetEndPointRefs(value []QName)
	SetParticipantMultiplicity(value *ParticipantMultiplicity)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Participant) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Participant) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Participant) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	if value := t.ParticipantMultiplicityField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	return
}
func (t *Participant) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *Participant) SetName(value *string) {
	t.NameField = value
}
func (t *Participant) ProcessRef() (result *QName, present bool) {
	if t.ProcessRefField != nil {
		present = true
	}
	result = t.ProcessRefField
	return
}
func (t *Participant) SetProcessRef(value *QName) {
	t.ProcessRefField = value
}
func (t *Participant) InterfaceRefs() (result *[]QName) {
	result = &t.InterfaceRefField
	return
}
func (t *Participant) SetInterfaceRefs(value []QName) {
	t.InterfaceRefField = value
}
func (t *Participant) EndPointRefs() (result *[]QName) {
	result = &t.EndPointRefField
	return
}
func (t *Participant) SetEndPointRefs(value []QName) {
	t.EndPointRefField = value
}
func (t *Participant) ParticipantMultiplicity() (result *ParticipantMultiplicity, present bool) {
	if t.ParticipantMultiplicityField != nil {
		present = true
	}
	result = t.ParticipantMultiplicityField
	return
}
func (t *Participant) SetParticipantMultiplicity(value *ParticipantMultiplicity) {
	t.ParticipantMultiplicityField = value
}

type ParticipantAssociation struct {
	BaseElement
	InnerParticipantRefField QName    `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL innerParticipantRef"`
	OuterParticipantRefField QName    `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL outerParticipantRef"`
	TextPayloadField         *Payload `xml:",chardata"`
}

func DefaultParticipantAssociation() ParticipantAssociation {
	return ParticipantAssociation{
		BaseElement: DefaultBaseElement(),
	}
}

type ParticipantAssociationInterface interface {
	Element
	BaseElementInterface
	InnerParticipantRef() (result *QName)
	OuterParticipantRef() (result *QName)
	SetInnerParticipantRef(value QName)
	SetOuterParticipantRef(value QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ParticipantAssociation) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ParticipantAssociation) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ParticipantAssociation) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}
func (t *ParticipantAssociation) InnerParticipantRef() (result *QName) {
	result = &t.InnerParticipantRefField
	return
}
func (t *ParticipantAssociation) SetInnerParticipantRef(value QName) {
	t.InnerParticipantRefField = value
}
func (t *ParticipantAssociation) OuterParticipantRef() (result *QName) {
	result = &t.OuterParticipantRefField
	return
}
func (t *ParticipantAssociation) SetOuterParticipantRef(value QName) {
	t.OuterParticipantRefField = value
}

type ParticipantMultiplicity struct {
	BaseElement
	MinimumField     *int32   `xml:"minimum,attr,omitempty"`
	MaximumField     *int32   `xml:"maximum,attr,omitempty"`
	TextPayloadField *Payload `xml:",chardata"`
}

var defaultParticipantMultiplicityMinimumField int32 = 0
var defaultParticipantMultiplicityMaximumField int32 = 1

func DefaultParticipantMultiplicity() ParticipantMultiplicity {
	return ParticipantMultiplicity{
		BaseElement:  DefaultBaseElement(),
		MinimumField: &defaultParticipantMultiplicityMinimumField,
		MaximumField: &defaultParticipantMultiplicityMaximumField,
	}
}

type ParticipantMultiplicityInterface interface {
	Element
	BaseElementInterface
	Minimum() (result int32)
	Maximum() (result int32)
	SetMinimum(value *int32)
	SetMaximum(value *int32)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ParticipantMultiplicity) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ParticipantMultiplicity) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ParticipantMultiplicity) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}
func (t *ParticipantMultiplicity) Minimum() (result int32) {
	if t.MinimumField == nil {
		result = defaultParticipantMultiplicityMinimumField
		return
	}
	result = *t.MinimumField
	return
}
func (t *ParticipantMultiplicity) SetMinimum(value *int32) {
	t.MinimumField = value
}
func (t *ParticipantMultiplicity) Maximum() (result int32) {
	if t.MaximumField == nil {
		result = defaultParticipantMultiplicityMaximumField
		return
	}
	result = *t.MaximumField
	return
}
func (t *ParticipantMultiplicity) SetMaximum(value *int32) {
	t.MaximumField = value
}

type PartnerEntity struct {
	RootElement
	NameField           *string  `xml:"name,attr,omitempty"`
	ParticipantRefField []QName  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL participantRef"`
	TextPayloadField    *Payload `xml:",chardata"`
}

func DefaultPartnerEntity() PartnerEntity {
	return PartnerEntity{
		RootElement: DefaultRootElement(),
	}
}

type PartnerEntityInterface interface {
	Element
	RootElementInterface
	Name() (result *string, present bool)
	ParticipantRefs() (result *[]QName)
	SetName(value *string)
	SetParticipantRefs(value []QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *PartnerEntity) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *PartnerEntity) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *PartnerEntity) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.RootElement.FindBy(f); found {
		return
	}

	return
}
func (t *PartnerEntity) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *PartnerEntity) SetName(value *string) {
	t.NameField = value
}
func (t *PartnerEntity) ParticipantRefs() (result *[]QName) {
	result = &t.ParticipantRefField
	return
}
func (t *PartnerEntity) SetParticipantRefs(value []QName) {
	t.ParticipantRefField = value
}

type PartnerRole struct {
	RootElement
	NameField           *string  `xml:"name,attr,omitempty"`
	ParticipantRefField []QName  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL participantRef"`
	TextPayloadField    *Payload `xml:",chardata"`
}

func DefaultPartnerRole() PartnerRole {
	return PartnerRole{
		RootElement: DefaultRootElement(),
	}
}

type PartnerRoleInterface interface {
	Element
	RootElementInterface
	Name() (result *string, present bool)
	ParticipantRefs() (result *[]QName)
	SetName(value *string)
	SetParticipantRefs(value []QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *PartnerRole) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *PartnerRole) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *PartnerRole) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.RootElement.FindBy(f); found {
		return
	}

	return
}
func (t *PartnerRole) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *PartnerRole) SetName(value *string) {
	t.NameField = value
}
func (t *PartnerRole) ParticipantRefs() (result *[]QName) {
	result = &t.ParticipantRefField
	return
}
func (t *PartnerRole) SetParticipantRefs(value []QName) {
	t.ParticipantRefField = value
}

type Performer struct {
	ResourceRole
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultPerformer() Performer {
	return Performer{
		ResourceRole: DefaultResourceRole(),
	}
}

type PerformerInterface interface {
	Element
	ResourceRoleInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Performer) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Performer) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Performer) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.ResourceRole.FindBy(f); found {
		return
	}

	return
}

type PotentialOwner struct {
	HumanPerformer
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultPotentialOwner() PotentialOwner {
	return PotentialOwner{
		HumanPerformer: DefaultHumanPerformer(),
	}
}

type PotentialOwnerInterface interface {
	Element
	HumanPerformerInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *PotentialOwner) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *PotentialOwner) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *PotentialOwner) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.HumanPerformer.FindBy(f); found {
		return
	}

	return
}

type Process struct {
	CallableElement
	ProcessTypeField                  *ProcessType              `xml:"processType,attr,omitempty"`
	IsClosedField                     *bool                     `xml:"isClosed,attr,omitempty"`
	IsExecutableField                 *bool                     `xml:"isExecutable,attr,omitempty"`
	DefinitionalCollaborationRefField *QName                    `xml:"definitionalCollaborationRef,attr,omitempty"`
	AuditingField                     *Auditing                 `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL auditing"`
	MonitoringField                   *Monitoring               `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL monitoring"`
	PropertyField                     []Property                `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL property"`
	LaneSetField                      []LaneSet                 `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL laneSet"`
	AdHocSubProcessField              []AdHocSubProcess         `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL adHocSubProcess"`
	BoundaryEventField                []BoundaryEvent           `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL boundaryEvent"`
	BusinessRuleTaskField             []BusinessRuleTask        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL businessRuleTask"`
	CallActivityField                 []CallActivity            `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL callActivity"`
	CallChoreographyField             []CallChoreography        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL callChoreography"`
	ChoreographyTaskField             []ChoreographyTask        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL choreographyTask"`
	ComplexGatewayField               []ComplexGateway          `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL complexGateway"`
	DataObjectField                   []DataObject              `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataObject"`
	DataObjectReferenceField          []DataObjectReference     `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataObjectReference"`
	DataStoreReferenceField           []DataStoreReference      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataStoreReference"`
	EndEventField                     []EndEvent                `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL endEvent"`
	EventField                        []Event                   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL event"`
	EventBasedGatewayField            []EventBasedGateway       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL eventBasedGateway"`
	ExclusiveGatewayField             []ExclusiveGateway        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL exclusiveGateway"`
	ImplicitThrowEventField           []ImplicitThrowEvent      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL implicitThrowEvent"`
	InclusiveGatewayField             []InclusiveGateway        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL inclusiveGateway"`
	IntermediateCatchEventField       []IntermediateCatchEvent  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL intermediateCatchEvent"`
	IntermediateThrowEventField       []IntermediateThrowEvent  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL intermediateThrowEvent"`
	ManualTaskField                   []ManualTask              `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL manualTask"`
	ParallelGatewayField              []ParallelGateway         `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL parallelGateway"`
	ReceiveTaskField                  []ReceiveTask             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL receiveTask"`
	ScriptTaskField                   []ScriptTask              `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL scriptTask"`
	SendTaskField                     []SendTask                `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL sendTask"`
	SequenceFlowField                 []SequenceFlow            `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL sequenceFlow"`
	ServiceTaskField                  []ServiceTask             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL serviceTask"`
	StartEventField                   []StartEvent              `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL startEvent"`
	SubChoreographyField              []SubChoreography         `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL subChoreography"`
	SubProcessField                   []SubProcess              `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL subProcess"`
	TaskField                         []Task                    `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL task"`
	TransactionField                  []Transaction             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL transaction"`
	UserTaskField                     []UserTask                `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL userTask"`
	AssociationField                  []Association             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL association"`
	GroupField                        []Group                   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL group"`
	TextAnnotationField               []TextAnnotation          `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL textAnnotation"`
	ResourceRoleField                 []ResourceRole            `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL resourceRole"`
	CorrelationSubscriptionField      []CorrelationSubscription `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL correlationSubscription"`
	SupportsField                     []QName                   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL supports"`
	TextPayloadField                  *Payload                  `xml:",chardata"`
}

var defaultProcessProcessTypeField ProcessType = "None"
var defaultProcessIsClosedField bool = false

func DefaultProcess() Process {
	return Process{
		CallableElement:  DefaultCallableElement(),
		ProcessTypeField: &defaultProcessProcessTypeField,
		IsClosedField:    &defaultProcessIsClosedField,
	}
}

type ProcessInterface interface {
	Element
	CallableElementInterface
	ProcessType() (result *ProcessType)
	IsClosed() (result bool)
	IsExecutable() (result bool, present bool)
	DefinitionalCollaborationRef() (result *QName, present bool)
	Auditing() (result *Auditing, present bool)
	Monitoring() (result *Monitoring, present bool)
	Properties() (result *[]Property)
	LaneSets() (result *[]LaneSet)
	AdHocSubProcesses() (result *[]AdHocSubProcess)
	BoundaryEvents() (result *[]BoundaryEvent)
	BusinessRuleTasks() (result *[]BusinessRuleTask)
	CallActivities() (result *[]CallActivity)
	CallChoreographies() (result *[]CallChoreography)
	ChoreographyTasks() (result *[]ChoreographyTask)
	ComplexGateways() (result *[]ComplexGateway)
	DataObjects() (result *[]DataObject)
	DataObjectReferences() (result *[]DataObjectReference)
	DataStoreReferences() (result *[]DataStoreReference)
	EndEvents() (result *[]EndEvent)
	Events() (result *[]Event)
	EventBasedGateways() (result *[]EventBasedGateway)
	ExclusiveGateways() (result *[]ExclusiveGateway)
	ImplicitThrowEvents() (result *[]ImplicitThrowEvent)
	InclusiveGateways() (result *[]InclusiveGateway)
	IntermediateCatchEvents() (result *[]IntermediateCatchEvent)
	IntermediateThrowEvents() (result *[]IntermediateThrowEvent)
	ManualTasks() (result *[]ManualTask)
	ParallelGateways() (result *[]ParallelGateway)
	ReceiveTasks() (result *[]ReceiveTask)
	ScriptTasks() (result *[]ScriptTask)
	SendTasks() (result *[]SendTask)
	SequenceFlows() (result *[]SequenceFlow)
	ServiceTasks() (result *[]ServiceTask)
	StartEvents() (result *[]StartEvent)
	SubChoreographies() (result *[]SubChoreography)
	SubProcesses() (result *[]SubProcess)
	Tasks() (result *[]Task)
	Transactions() (result *[]Transaction)
	UserTasks() (result *[]UserTask)
	Associations() (result *[]Association)
	Groups() (result *[]Group)
	TextAnnotations() (result *[]TextAnnotation)
	ResourceRoles() (result *[]ResourceRole)
	CorrelationSubscriptions() (result *[]CorrelationSubscription)
	Supportses() (result *[]QName)
	FlowElements() []FlowElementInterface
	Artifacts() []ArtifactInterface
	SetProcessType(value *ProcessType)
	SetIsClosed(value *bool)
	SetIsExecutable(value *bool)
	SetDefinitionalCollaborationRef(value *QName)
	SetAuditing(value *Auditing)
	SetMonitoring(value *Monitoring)
	SetProperties(value []Property)
	SetLaneSets(value []LaneSet)
	SetAdHocSubProcesses(value []AdHocSubProcess)
	SetBoundaryEvents(value []BoundaryEvent)
	SetBusinessRuleTasks(value []BusinessRuleTask)
	SetCallActivities(value []CallActivity)
	SetCallChoreographies(value []CallChoreography)
	SetChoreographyTasks(value []ChoreographyTask)
	SetComplexGateways(value []ComplexGateway)
	SetDataObjects(value []DataObject)
	SetDataObjectReferences(value []DataObjectReference)
	SetDataStoreReferences(value []DataStoreReference)
	SetEndEvents(value []EndEvent)
	SetEvents(value []Event)
	SetEventBasedGateways(value []EventBasedGateway)
	SetExclusiveGateways(value []ExclusiveGateway)
	SetImplicitThrowEvents(value []ImplicitThrowEvent)
	SetInclusiveGateways(value []InclusiveGateway)
	SetIntermediateCatchEvents(value []IntermediateCatchEvent)
	SetIntermediateThrowEvents(value []IntermediateThrowEvent)
	SetManualTasks(value []ManualTask)
	SetParallelGateways(value []ParallelGateway)
	SetReceiveTasks(value []ReceiveTask)
	SetScriptTasks(value []ScriptTask)
	SetSendTasks(value []SendTask)
	SetSequenceFlows(value []SequenceFlow)
	SetServiceTasks(value []ServiceTask)
	SetStartEvents(value []StartEvent)
	SetSubChoreographies(value []SubChoreography)
	SetSubProcesses(value []SubProcess)
	SetTasks(value []Task)
	SetTransactions(value []Transaction)
	SetUserTasks(value []UserTask)
	SetAssociations(value []Association)
	SetGroups(value []Group)
	SetTextAnnotations(value []TextAnnotation)
	SetResourceRoles(value []ResourceRole)
	SetCorrelationSubscriptions(value []CorrelationSubscription)
	SetSupportses(value []QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Process) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Process) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Process) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.CallableElement.FindBy(f); found {
		return
	}

	if value := t.AuditingField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	if value := t.MonitoringField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	for i := range t.PropertyField {
		if result, found = t.PropertyField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.LaneSetField {
		if result, found = t.LaneSetField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.AdHocSubProcessField {
		if result, found = t.AdHocSubProcessField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.BoundaryEventField {
		if result, found = t.BoundaryEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.BusinessRuleTaskField {
		if result, found = t.BusinessRuleTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.CallActivityField {
		if result, found = t.CallActivityField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.CallChoreographyField {
		if result, found = t.CallChoreographyField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ChoreographyTaskField {
		if result, found = t.ChoreographyTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ComplexGatewayField {
		if result, found = t.ComplexGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataObjectField {
		if result, found = t.DataObjectField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataObjectReferenceField {
		if result, found = t.DataObjectReferenceField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataStoreReferenceField {
		if result, found = t.DataStoreReferenceField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.EndEventField {
		if result, found = t.EndEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.EventField {
		if result, found = t.EventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.EventBasedGatewayField {
		if result, found = t.EventBasedGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ExclusiveGatewayField {
		if result, found = t.ExclusiveGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ImplicitThrowEventField {
		if result, found = t.ImplicitThrowEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.InclusiveGatewayField {
		if result, found = t.InclusiveGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.IntermediateCatchEventField {
		if result, found = t.IntermediateCatchEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.IntermediateThrowEventField {
		if result, found = t.IntermediateThrowEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ManualTaskField {
		if result, found = t.ManualTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ParallelGatewayField {
		if result, found = t.ParallelGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ReceiveTaskField {
		if result, found = t.ReceiveTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ScriptTaskField {
		if result, found = t.ScriptTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SendTaskField {
		if result, found = t.SendTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SequenceFlowField {
		if result, found = t.SequenceFlowField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ServiceTaskField {
		if result, found = t.ServiceTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.StartEventField {
		if result, found = t.StartEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SubChoreographyField {
		if result, found = t.SubChoreographyField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SubProcessField {
		if result, found = t.SubProcessField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.TaskField {
		if result, found = t.TaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.TransactionField {
		if result, found = t.TransactionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.UserTaskField {
		if result, found = t.UserTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.AssociationField {
		if result, found = t.AssociationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.GroupField {
		if result, found = t.GroupField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.TextAnnotationField {
		if result, found = t.TextAnnotationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ResourceRoleField {
		if result, found = t.ResourceRoleField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.CorrelationSubscriptionField {
		if result, found = t.CorrelationSubscriptionField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *Process) ProcessType() (result *ProcessType) {
	if t.ProcessTypeField == nil {
		result = &defaultProcessProcessTypeField
		return
	}
	result = t.ProcessTypeField
	return
}
func (t *Process) SetProcessType(value *ProcessType) {
	t.ProcessTypeField = value
}
func (t *Process) IsClosed() (result bool) {
	if t.IsClosedField == nil {
		result = defaultProcessIsClosedField
		return
	}
	result = *t.IsClosedField
	return
}
func (t *Process) SetIsClosed(value *bool) {
	t.IsClosedField = value
}
func (t *Process) IsExecutable() (result bool, present bool) {
	if t.IsExecutableField != nil {
		present = true
	}
	result = *t.IsExecutableField
	return
}
func (t *Process) SetIsExecutable(value *bool) {
	t.IsExecutableField = value
}
func (t *Process) DefinitionalCollaborationRef() (result *QName, present bool) {
	if t.DefinitionalCollaborationRefField != nil {
		present = true
	}
	result = t.DefinitionalCollaborationRefField
	return
}
func (t *Process) SetDefinitionalCollaborationRef(value *QName) {
	t.DefinitionalCollaborationRefField = value
}
func (t *Process) FlowElements() []FlowElementInterface {

	result := make([]FlowElementInterface, 0)

	for i := range t.AdHocSubProcessField {
		result = append(result, &t.AdHocSubProcessField[i])
	}

	for i := range t.BoundaryEventField {
		result = append(result, &t.BoundaryEventField[i])
	}

	for i := range t.BusinessRuleTaskField {
		result = append(result, &t.BusinessRuleTaskField[i])
	}

	for i := range t.CallActivityField {
		result = append(result, &t.CallActivityField[i])
	}

	for i := range t.CallChoreographyField {
		result = append(result, &t.CallChoreographyField[i])
	}

	for i := range t.ChoreographyTaskField {
		result = append(result, &t.ChoreographyTaskField[i])
	}

	for i := range t.ComplexGatewayField {
		result = append(result, &t.ComplexGatewayField[i])
	}

	for i := range t.DataObjectField {
		result = append(result, &t.DataObjectField[i])
	}

	for i := range t.DataObjectReferenceField {
		result = append(result, &t.DataObjectReferenceField[i])
	}

	for i := range t.DataStoreReferenceField {
		result = append(result, &t.DataStoreReferenceField[i])
	}

	for i := range t.EndEventField {
		result = append(result, &t.EndEventField[i])
	}

	for i := range t.EventField {
		result = append(result, &t.EventField[i])
	}

	for i := range t.EventBasedGatewayField {
		result = append(result, &t.EventBasedGatewayField[i])
	}

	for i := range t.ExclusiveGatewayField {
		result = append(result, &t.ExclusiveGatewayField[i])
	}

	for i := range t.ImplicitThrowEventField {
		result = append(result, &t.ImplicitThrowEventField[i])
	}

	for i := range t.InclusiveGatewayField {
		result = append(result, &t.InclusiveGatewayField[i])
	}

	for i := range t.IntermediateCatchEventField {
		result = append(result, &t.IntermediateCatchEventField[i])
	}

	for i := range t.IntermediateThrowEventField {
		result = append(result, &t.IntermediateThrowEventField[i])
	}

	for i := range t.ManualTaskField {
		result = append(result, &t.ManualTaskField[i])
	}

	for i := range t.ParallelGatewayField {
		result = append(result, &t.ParallelGatewayField[i])
	}

	for i := range t.ReceiveTaskField {
		result = append(result, &t.ReceiveTaskField[i])
	}

	for i := range t.ScriptTaskField {
		result = append(result, &t.ScriptTaskField[i])
	}

	for i := range t.SendTaskField {
		result = append(result, &t.SendTaskField[i])
	}

	for i := range t.SequenceFlowField {
		result = append(result, &t.SequenceFlowField[i])
	}

	for i := range t.ServiceTaskField {
		result = append(result, &t.ServiceTaskField[i])
	}

	for i := range t.StartEventField {
		result = append(result, &t.StartEventField[i])
	}

	for i := range t.SubChoreographyField {
		result = append(result, &t.SubChoreographyField[i])
	}

	for i := range t.SubProcessField {
		result = append(result, &t.SubProcessField[i])
	}

	for i := range t.TaskField {
		result = append(result, &t.TaskField[i])
	}

	for i := range t.TransactionField {
		result = append(result, &t.TransactionField[i])
	}

	for i := range t.UserTaskField {
		result = append(result, &t.UserTaskField[i])
	}
	return result
}
func (t *Process) Artifacts() []ArtifactInterface {

	result := make([]ArtifactInterface, 0)

	for i := range t.AssociationField {
		result = append(result, &t.AssociationField[i])
	}

	for i := range t.GroupField {
		result = append(result, &t.GroupField[i])
	}

	for i := range t.TextAnnotationField {
		result = append(result, &t.TextAnnotationField[i])
	}
	return result
}
func (t *Process) Auditing() (result *Auditing, present bool) {
	if t.AuditingField != nil {
		present = true
	}
	result = t.AuditingField
	return
}
func (t *Process) SetAuditing(value *Auditing) {
	t.AuditingField = value
}
func (t *Process) Monitoring() (result *Monitoring, present bool) {
	if t.MonitoringField != nil {
		present = true
	}
	result = t.MonitoringField
	return
}
func (t *Process) SetMonitoring(value *Monitoring) {
	t.MonitoringField = value
}
func (t *Process) Properties() (result *[]Property) {
	result = &t.PropertyField
	return
}
func (t *Process) SetProperties(value []Property) {
	t.PropertyField = value
}
func (t *Process) LaneSets() (result *[]LaneSet) {
	result = &t.LaneSetField
	return
}
func (t *Process) SetLaneSets(value []LaneSet) {
	t.LaneSetField = value
}
func (t *Process) AdHocSubProcesses() (result *[]AdHocSubProcess) {
	result = &t.AdHocSubProcessField
	return
}
func (t *Process) SetAdHocSubProcesses(value []AdHocSubProcess) {
	t.AdHocSubProcessField = value
}
func (t *Process) BoundaryEvents() (result *[]BoundaryEvent) {
	result = &t.BoundaryEventField
	return
}
func (t *Process) SetBoundaryEvents(value []BoundaryEvent) {
	t.BoundaryEventField = value
}
func (t *Process) BusinessRuleTasks() (result *[]BusinessRuleTask) {
	result = &t.BusinessRuleTaskField
	return
}
func (t *Process) SetBusinessRuleTasks(value []BusinessRuleTask) {
	t.BusinessRuleTaskField = value
}
func (t *Process) CallActivities() (result *[]CallActivity) {
	result = &t.CallActivityField
	return
}
func (t *Process) SetCallActivities(value []CallActivity) {
	t.CallActivityField = value
}
func (t *Process) CallChoreographies() (result *[]CallChoreography) {
	result = &t.CallChoreographyField
	return
}
func (t *Process) SetCallChoreographies(value []CallChoreography) {
	t.CallChoreographyField = value
}
func (t *Process) ChoreographyTasks() (result *[]ChoreographyTask) {
	result = &t.ChoreographyTaskField
	return
}
func (t *Process) SetChoreographyTasks(value []ChoreographyTask) {
	t.ChoreographyTaskField = value
}
func (t *Process) ComplexGateways() (result *[]ComplexGateway) {
	result = &t.ComplexGatewayField
	return
}
func (t *Process) SetComplexGateways(value []ComplexGateway) {
	t.ComplexGatewayField = value
}
func (t *Process) DataObjects() (result *[]DataObject) {
	result = &t.DataObjectField
	return
}
func (t *Process) SetDataObjects(value []DataObject) {
	t.DataObjectField = value
}
func (t *Process) DataObjectReferences() (result *[]DataObjectReference) {
	result = &t.DataObjectReferenceField
	return
}
func (t *Process) SetDataObjectReferences(value []DataObjectReference) {
	t.DataObjectReferenceField = value
}
func (t *Process) DataStoreReferences() (result *[]DataStoreReference) {
	result = &t.DataStoreReferenceField
	return
}
func (t *Process) SetDataStoreReferences(value []DataStoreReference) {
	t.DataStoreReferenceField = value
}
func (t *Process) EndEvents() (result *[]EndEvent) {
	result = &t.EndEventField
	return
}
func (t *Process) SetEndEvents(value []EndEvent) {
	t.EndEventField = value
}
func (t *Process) Events() (result *[]Event) {
	result = &t.EventField
	return
}
func (t *Process) SetEvents(value []Event) {
	t.EventField = value
}
func (t *Process) EventBasedGateways() (result *[]EventBasedGateway) {
	result = &t.EventBasedGatewayField
	return
}
func (t *Process) SetEventBasedGateways(value []EventBasedGateway) {
	t.EventBasedGatewayField = value
}
func (t *Process) ExclusiveGateways() (result *[]ExclusiveGateway) {
	result = &t.ExclusiveGatewayField
	return
}
func (t *Process) SetExclusiveGateways(value []ExclusiveGateway) {
	t.ExclusiveGatewayField = value
}
func (t *Process) ImplicitThrowEvents() (result *[]ImplicitThrowEvent) {
	result = &t.ImplicitThrowEventField
	return
}
func (t *Process) SetImplicitThrowEvents(value []ImplicitThrowEvent) {
	t.ImplicitThrowEventField = value
}
func (t *Process) InclusiveGateways() (result *[]InclusiveGateway) {
	result = &t.InclusiveGatewayField
	return
}
func (t *Process) SetInclusiveGateways(value []InclusiveGateway) {
	t.InclusiveGatewayField = value
}
func (t *Process) IntermediateCatchEvents() (result *[]IntermediateCatchEvent) {
	result = &t.IntermediateCatchEventField
	return
}
func (t *Process) SetIntermediateCatchEvents(value []IntermediateCatchEvent) {
	t.IntermediateCatchEventField = value
}
func (t *Process) IntermediateThrowEvents() (result *[]IntermediateThrowEvent) {
	result = &t.IntermediateThrowEventField
	return
}
func (t *Process) SetIntermediateThrowEvents(value []IntermediateThrowEvent) {
	t.IntermediateThrowEventField = value
}
func (t *Process) ManualTasks() (result *[]ManualTask) {
	result = &t.ManualTaskField
	return
}
func (t *Process) SetManualTasks(value []ManualTask) {
	t.ManualTaskField = value
}
func (t *Process) ParallelGateways() (result *[]ParallelGateway) {
	result = &t.ParallelGatewayField
	return
}
func (t *Process) SetParallelGateways(value []ParallelGateway) {
	t.ParallelGatewayField = value
}
func (t *Process) ReceiveTasks() (result *[]ReceiveTask) {
	result = &t.ReceiveTaskField
	return
}
func (t *Process) SetReceiveTasks(value []ReceiveTask) {
	t.ReceiveTaskField = value
}
func (t *Process) ScriptTasks() (result *[]ScriptTask) {
	result = &t.ScriptTaskField
	return
}
func (t *Process) SetScriptTasks(value []ScriptTask) {
	t.ScriptTaskField = value
}
func (t *Process) SendTasks() (result *[]SendTask) {
	result = &t.SendTaskField
	return
}
func (t *Process) SetSendTasks(value []SendTask) {
	t.SendTaskField = value
}
func (t *Process) SequenceFlows() (result *[]SequenceFlow) {
	result = &t.SequenceFlowField
	return
}
func (t *Process) SetSequenceFlows(value []SequenceFlow) {
	t.SequenceFlowField = value
}
func (t *Process) ServiceTasks() (result *[]ServiceTask) {
	result = &t.ServiceTaskField
	return
}
func (t *Process) SetServiceTasks(value []ServiceTask) {
	t.ServiceTaskField = value
}
func (t *Process) StartEvents() (result *[]StartEvent) {
	result = &t.StartEventField
	return
}
func (t *Process) SetStartEvents(value []StartEvent) {
	t.StartEventField = value
}
func (t *Process) SubChoreographies() (result *[]SubChoreography) {
	result = &t.SubChoreographyField
	return
}
func (t *Process) SetSubChoreographies(value []SubChoreography) {
	t.SubChoreographyField = value
}
func (t *Process) SubProcesses() (result *[]SubProcess) {
	result = &t.SubProcessField
	return
}
func (t *Process) SetSubProcesses(value []SubProcess) {
	t.SubProcessField = value
}
func (t *Process) Tasks() (result *[]Task) {
	result = &t.TaskField
	return
}
func (t *Process) SetTasks(value []Task) {
	t.TaskField = value
}
func (t *Process) Transactions() (result *[]Transaction) {
	result = &t.TransactionField
	return
}
func (t *Process) SetTransactions(value []Transaction) {
	t.TransactionField = value
}
func (t *Process) UserTasks() (result *[]UserTask) {
	result = &t.UserTaskField
	return
}
func (t *Process) SetUserTasks(value []UserTask) {
	t.UserTaskField = value
}
func (t *Process) Associations() (result *[]Association) {
	result = &t.AssociationField
	return
}
func (t *Process) SetAssociations(value []Association) {
	t.AssociationField = value
}
func (t *Process) Groups() (result *[]Group) {
	result = &t.GroupField
	return
}
func (t *Process) SetGroups(value []Group) {
	t.GroupField = value
}
func (t *Process) TextAnnotations() (result *[]TextAnnotation) {
	result = &t.TextAnnotationField
	return
}
func (t *Process) SetTextAnnotations(value []TextAnnotation) {
	t.TextAnnotationField = value
}
func (t *Process) ResourceRoles() (result *[]ResourceRole) {
	result = &t.ResourceRoleField
	return
}
func (t *Process) SetResourceRoles(value []ResourceRole) {
	t.ResourceRoleField = value
}
func (t *Process) CorrelationSubscriptions() (result *[]CorrelationSubscription) {
	result = &t.CorrelationSubscriptionField
	return
}
func (t *Process) SetCorrelationSubscriptions(value []CorrelationSubscription) {
	t.CorrelationSubscriptionField = value
}
func (t *Process) Supportses() (result *[]QName) {
	result = &t.SupportsField
	return
}
func (t *Process) SetSupportses(value []QName) {
	t.SupportsField = value
}

type Property struct {
	ItemAwareElement
	NameField        *string  `xml:"name,attr,omitempty"`
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultProperty() Property {
	return Property{
		ItemAwareElement: DefaultItemAwareElement(),
	}
}

type PropertyInterface interface {
	Element
	ItemAwareElementInterface
	Name() (result *string, present bool)
	SetName(value *string)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Property) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Property) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Property) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.ItemAwareElement.FindBy(f); found {
		return
	}

	return
}
func (t *Property) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *Property) SetName(value *string) {
	t.NameField = value
}

type ReceiveTask struct {
	Task
	ImplementationField *Implementation `xml:"implementation,attr,omitempty"`
	InstantiateField    *bool           `xml:"instantiate,attr,omitempty"`
	MessageRefField     *QName          `xml:"messageRef,attr,omitempty"`
	OperationRefField   *QName          `xml:"operationRef,attr,omitempty"`
	TextPayloadField    *Payload        `xml:",chardata"`
}

var defaultReceiveTaskImplementationField Implementation = "##WebService"
var defaultReceiveTaskInstantiateField bool = false

func DefaultReceiveTask() ReceiveTask {
	return ReceiveTask{
		Task:             DefaultTask(),
		InstantiateField: &defaultReceiveTaskInstantiateField,
	}
}

type ReceiveTaskInterface interface {
	Element
	TaskInterface
	Implementation() (result *Implementation)
	Instantiate() (result bool)
	MessageRef() (result *QName, present bool)
	OperationRef() (result *QName, present bool)
	SetImplementation(value *Implementation)
	SetInstantiate(value *bool)
	SetMessageRef(value *QName)
	SetOperationRef(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ReceiveTask) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ReceiveTask) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ReceiveTask) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Task.FindBy(f); found {
		return
	}

	return
}
func (t *ReceiveTask) Implementation() (result *Implementation) {
	if t.ImplementationField == nil {
		result = &defaultReceiveTaskImplementationField
		return
	}
	result = t.ImplementationField
	return
}
func (t *ReceiveTask) SetImplementation(value *Implementation) {
	t.ImplementationField = value
}
func (t *ReceiveTask) Instantiate() (result bool) {
	if t.InstantiateField == nil {
		result = defaultReceiveTaskInstantiateField
		return
	}
	result = *t.InstantiateField
	return
}
func (t *ReceiveTask) SetInstantiate(value *bool) {
	t.InstantiateField = value
}
func (t *ReceiveTask) MessageRef() (result *QName, present bool) {
	if t.MessageRefField != nil {
		present = true
	}
	result = t.MessageRefField
	return
}
func (t *ReceiveTask) SetMessageRef(value *QName) {
	t.MessageRefField = value
}
func (t *ReceiveTask) OperationRef() (result *QName, present bool) {
	if t.OperationRefField != nil {
		present = true
	}
	result = t.OperationRefField
	return
}
func (t *ReceiveTask) SetOperationRef(value *QName) {
	t.OperationRefField = value
}

type Relationship struct {
	BaseElement
	TypeField        string                 `xml:"type,attr,omitempty"`
	DirectionField   *RelationshipDirection `xml:"direction,attr,omitempty"`
	SourceField      []QName                `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL source"`
	TargetField      []QName                `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL target"`
	TextPayloadField *Payload               `xml:",chardata"`
}

func DefaultRelationship() Relationship {
	return Relationship{
		BaseElement: DefaultBaseElement(),
	}
}

type RelationshipInterface interface {
	Element
	BaseElementInterface
	Type() (result *string)
	Direction() (result *RelationshipDirection, present bool)
	Sources() (result *[]QName)
	Targets() (result *[]QName)
	SetType(value string)
	SetDirection(value *RelationshipDirection)
	SetSources(value []QName)
	SetTargets(value []QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Relationship) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Relationship) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Relationship) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}
func (t *Relationship) Type() (result *string) {
	result = &t.TypeField
	return
}
func (t *Relationship) SetType(value string) {
	t.TypeField = value
}
func (t *Relationship) Direction() (result *RelationshipDirection, present bool) {
	if t.DirectionField != nil {
		present = true
	}
	result = t.DirectionField
	return
}
func (t *Relationship) SetDirection(value *RelationshipDirection) {
	t.DirectionField = value
}
func (t *Relationship) Sources() (result *[]QName) {
	result = &t.SourceField
	return
}
func (t *Relationship) SetSources(value []QName) {
	t.SourceField = value
}
func (t *Relationship) Targets() (result *[]QName) {
	result = &t.TargetField
	return
}
func (t *Relationship) SetTargets(value []QName) {
	t.TargetField = value
}

type Rendering struct {
	BaseElement
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultRendering() Rendering {
	return Rendering{
		BaseElement: DefaultBaseElement(),
	}
}

type RenderingInterface interface {
	Element
	BaseElementInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Rendering) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Rendering) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Rendering) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}

type Resource struct {
	RootElement
	NameField              string              `xml:"name,attr,omitempty"`
	ResourceParameterField []ResourceParameter `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL resourceParameter"`
	TextPayloadField       *Payload            `xml:",chardata"`
}

func DefaultResource() Resource {
	return Resource{
		RootElement: DefaultRootElement(),
	}
}

type ResourceInterface interface {
	Element
	RootElementInterface
	Name() (result *string)
	ResourceParameters() (result *[]ResourceParameter)
	SetName(value string)
	SetResourceParameters(value []ResourceParameter)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Resource) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Resource) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Resource) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.RootElement.FindBy(f); found {
		return
	}

	for i := range t.ResourceParameterField {
		if result, found = t.ResourceParameterField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *Resource) Name() (result *string) {
	result = &t.NameField
	return
}
func (t *Resource) SetName(value string) {
	t.NameField = value
}
func (t *Resource) ResourceParameters() (result *[]ResourceParameter) {
	result = &t.ResourceParameterField
	return
}
func (t *Resource) SetResourceParameters(value []ResourceParameter) {
	t.ResourceParameterField = value
}

type ResourceAssignmentExpression struct {
	BaseElement
	ExpressionField  Expression `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL expression"`
	TextPayloadField *Payload   `xml:",chardata"`
}

func DefaultResourceAssignmentExpression() ResourceAssignmentExpression {
	return ResourceAssignmentExpression{
		BaseElement: DefaultBaseElement(),
	}
}

type ResourceAssignmentExpressionInterface interface {
	Element
	BaseElementInterface
	Expression() (result *Expression)
	SetExpression(value Expression)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ResourceAssignmentExpression) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ResourceAssignmentExpression) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ResourceAssignmentExpression) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	if result, found = t.ExpressionField.FindBy(f); found {
		return
	}

	return
}
func (t *ResourceAssignmentExpression) Expression() (result *Expression) {
	result = &t.ExpressionField
	return
}
func (t *ResourceAssignmentExpression) SetExpression(value Expression) {
	t.ExpressionField = value
}

type ResourceParameter struct {
	BaseElement
	NameField        *string  `xml:"name,attr,omitempty"`
	TypeField        *QName   `xml:"type,attr,omitempty"`
	IsRequiredField  *bool    `xml:"isRequired,attr,omitempty"`
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultResourceParameter() ResourceParameter {
	return ResourceParameter{
		BaseElement: DefaultBaseElement(),
	}
}

type ResourceParameterInterface interface {
	Element
	BaseElementInterface
	Name() (result *string, present bool)
	Type() (result *QName, present bool)
	IsRequired() (result bool, present bool)
	SetName(value *string)
	SetType(value *QName)
	SetIsRequired(value *bool)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ResourceParameter) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ResourceParameter) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ResourceParameter) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}
func (t *ResourceParameter) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *ResourceParameter) SetName(value *string) {
	t.NameField = value
}
func (t *ResourceParameter) Type() (result *QName, present bool) {
	if t.TypeField != nil {
		present = true
	}
	result = t.TypeField
	return
}
func (t *ResourceParameter) SetType(value *QName) {
	t.TypeField = value
}
func (t *ResourceParameter) IsRequired() (result bool, present bool) {
	if t.IsRequiredField != nil {
		present = true
	}
	result = *t.IsRequiredField
	return
}
func (t *ResourceParameter) SetIsRequired(value *bool) {
	t.IsRequiredField = value
}

type ResourceParameterBinding struct {
	BaseElement
	ParameterRefField QName      `xml:"parameterRef,attr,omitempty"`
	ExpressionField   Expression `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL expression"`
	TextPayloadField  *Payload   `xml:",chardata"`
}

func DefaultResourceParameterBinding() ResourceParameterBinding {
	return ResourceParameterBinding{
		BaseElement: DefaultBaseElement(),
	}
}

type ResourceParameterBindingInterface interface {
	Element
	BaseElementInterface
	ParameterRef() (result *QName)
	Expression() (result *Expression)
	SetParameterRef(value QName)
	SetExpression(value Expression)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ResourceParameterBinding) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ResourceParameterBinding) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ResourceParameterBinding) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	if result, found = t.ExpressionField.FindBy(f); found {
		return
	}

	return
}
func (t *ResourceParameterBinding) ParameterRef() (result *QName) {
	result = &t.ParameterRefField
	return
}
func (t *ResourceParameterBinding) SetParameterRef(value QName) {
	t.ParameterRefField = value
}
func (t *ResourceParameterBinding) Expression() (result *Expression) {
	result = &t.ExpressionField
	return
}
func (t *ResourceParameterBinding) SetExpression(value Expression) {
	t.ExpressionField = value
}

type ResourceRole struct {
	BaseElement
	NameField                         *string                       `xml:"name,attr,omitempty"`
	ResourceRefField                  QName                         `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL resourceRef"`
	ResourceParameterBindingField     []ResourceParameterBinding    `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL resourceParameterBinding"`
	ResourceAssignmentExpressionField *ResourceAssignmentExpression `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL resourceAssignmentExpression"`
	TextPayloadField                  *Payload                      `xml:",chardata"`
}

func DefaultResourceRole() ResourceRole {
	return ResourceRole{
		BaseElement: DefaultBaseElement(),
	}
}

type ResourceRoleInterface interface {
	Element
	BaseElementInterface
	Name() (result *string, present bool)
	ResourceRef() (result *QName)
	ResourceParameterBindings() (result *[]ResourceParameterBinding)
	ResourceAssignmentExpression() (result *ResourceAssignmentExpression, present bool)
	SetName(value *string)
	SetResourceRef(value QName)
	SetResourceParameterBindings(value []ResourceParameterBinding)
	SetResourceAssignmentExpression(value *ResourceAssignmentExpression)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ResourceRole) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ResourceRole) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ResourceRole) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	for i := range t.ResourceParameterBindingField {
		if result, found = t.ResourceParameterBindingField[i].FindBy(f); found {
			return
		}
	}

	if value := t.ResourceAssignmentExpressionField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	return
}
func (t *ResourceRole) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *ResourceRole) SetName(value *string) {
	t.NameField = value
}
func (t *ResourceRole) ResourceRef() (result *QName) {
	result = &t.ResourceRefField
	return
}
func (t *ResourceRole) SetResourceRef(value QName) {
	t.ResourceRefField = value
}
func (t *ResourceRole) ResourceParameterBindings() (result *[]ResourceParameterBinding) {
	result = &t.ResourceParameterBindingField
	return
}
func (t *ResourceRole) SetResourceParameterBindings(value []ResourceParameterBinding) {
	t.ResourceParameterBindingField = value
}
func (t *ResourceRole) ResourceAssignmentExpression() (result *ResourceAssignmentExpression, present bool) {
	if t.ResourceAssignmentExpressionField != nil {
		present = true
	}
	result = t.ResourceAssignmentExpressionField
	return
}
func (t *ResourceRole) SetResourceAssignmentExpression(value *ResourceAssignmentExpression) {
	t.ResourceAssignmentExpressionField = value
}

type RootElement struct {
	BaseElement
}

func DefaultRootElement() RootElement {
	return RootElement{
		BaseElement: DefaultBaseElement(),
	}
}

type RootElementInterface interface {
	Element
	BaseElementInterface
}

func (t *RootElement) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.BaseElement.FindBy(f); found {
		return
	}

	return
}

type ScriptTask struct {
	Task
	ScriptFormatField *string  `xml:"scriptFormat,attr,omitempty"`
	ScriptField       *Script  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL script"`
	TextPayloadField  *Payload `xml:",chardata"`
}

func DefaultScriptTask() ScriptTask {
	return ScriptTask{
		Task: DefaultTask(),
	}
}

type ScriptTaskInterface interface {
	Element
	TaskInterface
	ScriptFormat() (result *string, present bool)
	Script() (result *Script, present bool)
	SetScriptFormat(value *string)
	SetScript(value *Script)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ScriptTask) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ScriptTask) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ScriptTask) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Task.FindBy(f); found {
		return
	}

	if value := t.ScriptField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	return
}
func (t *ScriptTask) ScriptFormat() (result *string, present bool) {
	if t.ScriptFormatField != nil {
		present = true
	}
	result = t.ScriptFormatField
	return
}
func (t *ScriptTask) SetScriptFormat(value *string) {
	t.ScriptFormatField = value
}
func (t *ScriptTask) Script() (result *Script, present bool) {
	if t.ScriptField != nil {
		present = true
	}
	result = t.ScriptField
	return
}
func (t *ScriptTask) SetScript(value *Script) {
	t.ScriptField = value
}

type Script struct {
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultScript() Script {
	return Script{}
}

type ScriptInterface interface {
	Element

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Script) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Script) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Script) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}

	return
}

type SendTask struct {
	Task
	ImplementationField *Implementation `xml:"implementation,attr,omitempty"`
	MessageRefField     *QName          `xml:"messageRef,attr,omitempty"`
	OperationRefField   *QName          `xml:"operationRef,attr,omitempty"`
	TextPayloadField    *Payload        `xml:",chardata"`
}

var defaultSendTaskImplementationField Implementation = "##WebService"

func DefaultSendTask() SendTask {
	return SendTask{
		Task: DefaultTask(),
	}
}

type SendTaskInterface interface {
	Element
	TaskInterface
	Implementation() (result *Implementation)
	MessageRef() (result *QName, present bool)
	OperationRef() (result *QName, present bool)
	SetImplementation(value *Implementation)
	SetMessageRef(value *QName)
	SetOperationRef(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *SendTask) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *SendTask) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *SendTask) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Task.FindBy(f); found {
		return
	}

	return
}
func (t *SendTask) Implementation() (result *Implementation) {
	if t.ImplementationField == nil {
		result = &defaultSendTaskImplementationField
		return
	}
	result = t.ImplementationField
	return
}
func (t *SendTask) SetImplementation(value *Implementation) {
	t.ImplementationField = value
}
func (t *SendTask) MessageRef() (result *QName, present bool) {
	if t.MessageRefField != nil {
		present = true
	}
	result = t.MessageRefField
	return
}
func (t *SendTask) SetMessageRef(value *QName) {
	t.MessageRefField = value
}
func (t *SendTask) OperationRef() (result *QName, present bool) {
	if t.OperationRefField != nil {
		present = true
	}
	result = t.OperationRefField
	return
}
func (t *SendTask) SetOperationRef(value *QName) {
	t.OperationRefField = value
}

type SequenceFlow struct {
	FlowElement
	SourceRefField           IdRef         `xml:"sourceRef,attr,omitempty"`
	TargetRefField           IdRef         `xml:"targetRef,attr,omitempty"`
	IsImmediateField         *bool         `xml:"isImmediate,attr,omitempty"`
	ConditionExpressionField *AnExpression `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL conditionExpression"`
	TextPayloadField         *Payload      `xml:",chardata"`
}

func DefaultSequenceFlow() SequenceFlow {
	return SequenceFlow{
		FlowElement: DefaultFlowElement(),
	}
}

type SequenceFlowInterface interface {
	Element
	FlowElementInterface
	SourceRef() (result *IdRef)
	TargetRef() (result *IdRef)
	IsImmediate() (result *bool, present bool)
	ConditionExpression() (result *AnExpression, present bool)
	SetSourceRef(value IdRef)
	SetTargetRef(value IdRef)
	SetIsImmediate(value *bool)
	SetConditionExpression(value *AnExpression)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *SequenceFlow) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *SequenceFlow) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *SequenceFlow) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.FlowElement.FindBy(f); found {
		return
	}

	if value := t.ConditionExpressionField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	return
}
func (t *SequenceFlow) SourceRef() (result *IdRef) {
	result = &t.SourceRefField
	return
}
func (t *SequenceFlow) SetSourceRef(value IdRef) {
	t.SourceRefField = value
}
func (t *SequenceFlow) TargetRef() (result *IdRef) {
	result = &t.TargetRefField
	return
}
func (t *SequenceFlow) SetTargetRef(value IdRef) {
	t.TargetRefField = value
}
func (t *SequenceFlow) IsImmediate() (result *bool, present bool) {
	if t.IsImmediateField != nil {
		present = true
	}
	result = t.IsImmediateField
	return
}
func (t *SequenceFlow) SetIsImmediate(value *bool) {
	t.IsImmediateField = value
}
func (t *SequenceFlow) ConditionExpression() (result *AnExpression, present bool) {
	if t.ConditionExpressionField != nil {
		present = true
	}
	result = t.ConditionExpressionField
	return
}
func (t *SequenceFlow) SetConditionExpression(value *AnExpression) {
	t.ConditionExpressionField = value
}

type ServiceTask struct {
	Task
	ImplementationField *Implementation `xml:"implementation,attr,omitempty"`
	OperationRefField   *QName          `xml:"operationRef,attr,omitempty"`
	TextPayloadField    *Payload        `xml:",chardata"`
}

var defaultServiceTaskImplementationField Implementation = "##WebService"

func DefaultServiceTask() ServiceTask {
	return ServiceTask{
		Task: DefaultTask(),
	}
}

type ServiceTaskInterface interface {
	Element
	TaskInterface
	Implementation() (result *Implementation)
	OperationRef() (result *QName, present bool)
	SetImplementation(value *Implementation)
	SetOperationRef(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *ServiceTask) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *ServiceTask) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *ServiceTask) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Task.FindBy(f); found {
		return
	}

	return
}
func (t *ServiceTask) Implementation() (result *Implementation) {
	if t.ImplementationField == nil {
		result = &defaultServiceTaskImplementationField
		return
	}
	result = t.ImplementationField
	return
}
func (t *ServiceTask) SetImplementation(value *Implementation) {
	t.ImplementationField = value
}
func (t *ServiceTask) OperationRef() (result *QName, present bool) {
	if t.OperationRefField != nil {
		present = true
	}
	result = t.OperationRefField
	return
}
func (t *ServiceTask) SetOperationRef(value *QName) {
	t.OperationRefField = value
}

type Signal struct {
	RootElement
	NameField         *string  `xml:"name,attr,omitempty"`
	StructureRefField *QName   `xml:"structureRef,attr,omitempty"`
	TextPayloadField  *Payload `xml:",chardata"`
}

func DefaultSignal() Signal {
	return Signal{
		RootElement: DefaultRootElement(),
	}
}

type SignalInterface interface {
	Element
	RootElementInterface
	Name() (result *string, present bool)
	StructureRef() (result *QName, present bool)
	SetName(value *string)
	SetStructureRef(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Signal) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Signal) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Signal) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.RootElement.FindBy(f); found {
		return
	}

	return
}
func (t *Signal) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *Signal) SetName(value *string) {
	t.NameField = value
}
func (t *Signal) StructureRef() (result *QName, present bool) {
	if t.StructureRefField != nil {
		present = true
	}
	result = t.StructureRefField
	return
}
func (t *Signal) SetStructureRef(value *QName) {
	t.StructureRefField = value
}

type SignalEventDefinition struct {
	EventDefinition
	SignalRefField   *QName   `xml:"signalRef,attr,omitempty"`
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultSignalEventDefinition() SignalEventDefinition {
	return SignalEventDefinition{
		EventDefinition: DefaultEventDefinition(),
	}
}

type SignalEventDefinitionInterface interface {
	Element
	EventDefinitionInterface
	SignalRef() (result *QName, present bool)
	SetSignalRef(value *QName)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *SignalEventDefinition) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *SignalEventDefinition) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *SignalEventDefinition) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.EventDefinition.FindBy(f); found {
		return
	}

	return
}
func (t *SignalEventDefinition) SignalRef() (result *QName, present bool) {
	if t.SignalRefField != nil {
		present = true
	}
	result = t.SignalRefField
	return
}
func (t *SignalEventDefinition) SetSignalRef(value *QName) {
	t.SignalRefField = value
}

type StandardLoopCharacteristics struct {
	LoopCharacteristics
	TestBeforeField    *bool        `xml:"testBefore,attr,omitempty"`
	LoopMaximumField   *big.Int     `xml:"loopMaximum,attr,omitempty"`
	LoopConditionField AnExpression `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL loopCondition"`
	TextPayloadField   *Payload     `xml:",chardata"`
}

var defaultStandardLoopCharacteristicsTestBeforeField bool = false

func DefaultStandardLoopCharacteristics() StandardLoopCharacteristics {
	return StandardLoopCharacteristics{
		LoopCharacteristics: DefaultLoopCharacteristics(),
		TestBeforeField:     &defaultStandardLoopCharacteristicsTestBeforeField,
	}
}

type StandardLoopCharacteristicsInterface interface {
	Element
	LoopCharacteristicsInterface
	TestBefore() (result bool)
	LoopMaximum() (result *big.Int, present bool)
	LoopCondition() (result *AnExpression)
	SetTestBefore(value *bool)
	SetLoopMaximum(value *big.Int)
	SetLoopCondition(value AnExpression)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *StandardLoopCharacteristics) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *StandardLoopCharacteristics) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *StandardLoopCharacteristics) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.LoopCharacteristics.FindBy(f); found {
		return
	}

	if result, found = t.LoopConditionField.FindBy(f); found {
		return
	}

	return
}
func (t *StandardLoopCharacteristics) TestBefore() (result bool) {
	if t.TestBeforeField == nil {
		result = defaultStandardLoopCharacteristicsTestBeforeField
		return
	}
	result = *t.TestBeforeField
	return
}
func (t *StandardLoopCharacteristics) SetTestBefore(value *bool) {
	t.TestBeforeField = value
}
func (t *StandardLoopCharacteristics) LoopMaximum() (result *big.Int, present bool) {
	if t.LoopMaximumField != nil {
		present = true
	}
	result = t.LoopMaximumField
	return
}
func (t *StandardLoopCharacteristics) SetLoopMaximum(value *big.Int) {
	t.LoopMaximumField = value
}
func (t *StandardLoopCharacteristics) LoopCondition() (result *AnExpression) {
	result = &t.LoopConditionField
	return
}
func (t *StandardLoopCharacteristics) SetLoopCondition(value AnExpression) {
	t.LoopConditionField = value
}

type StartEvent struct {
	CatchEvent
	IsInterruptingField *bool    `xml:"isInterrupting,attr,omitempty"`
	TextPayloadField    *Payload `xml:",chardata"`
}

var defaultStartEventIsInterruptingField bool = true

func DefaultStartEvent() StartEvent {
	return StartEvent{
		CatchEvent:          DefaultCatchEvent(),
		IsInterruptingField: &defaultStartEventIsInterruptingField,
	}
}

type StartEventInterface interface {
	Element
	CatchEventInterface
	IsInterrupting() (result bool)
	SetIsInterrupting(value *bool)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *StartEvent) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *StartEvent) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *StartEvent) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.CatchEvent.FindBy(f); found {
		return
	}

	return
}
func (t *StartEvent) IsInterrupting() (result bool) {
	if t.IsInterruptingField == nil {
		result = defaultStartEventIsInterruptingField
		return
	}
	result = *t.IsInterruptingField
	return
}
func (t *StartEvent) SetIsInterrupting(value *bool) {
	t.IsInterruptingField = value
}

type SubChoreography struct {
	ChoreographyActivity
	AdHocSubProcessField        []AdHocSubProcess        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL adHocSubProcess"`
	BoundaryEventField          []BoundaryEvent          `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL boundaryEvent"`
	BusinessRuleTaskField       []BusinessRuleTask       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL businessRuleTask"`
	CallActivityField           []CallActivity           `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL callActivity"`
	CallChoreographyField       []CallChoreography       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL callChoreography"`
	ChoreographyTaskField       []ChoreographyTask       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL choreographyTask"`
	ComplexGatewayField         []ComplexGateway         `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL complexGateway"`
	DataObjectField             []DataObject             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataObject"`
	DataObjectReferenceField    []DataObjectReference    `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataObjectReference"`
	DataStoreReferenceField     []DataStoreReference     `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataStoreReference"`
	EndEventField               []EndEvent               `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL endEvent"`
	EventField                  []Event                  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL event"`
	EventBasedGatewayField      []EventBasedGateway      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL eventBasedGateway"`
	ExclusiveGatewayField       []ExclusiveGateway       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL exclusiveGateway"`
	ImplicitThrowEventField     []ImplicitThrowEvent     `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL implicitThrowEvent"`
	InclusiveGatewayField       []InclusiveGateway       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL inclusiveGateway"`
	IntermediateCatchEventField []IntermediateCatchEvent `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL intermediateCatchEvent"`
	IntermediateThrowEventField []IntermediateThrowEvent `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL intermediateThrowEvent"`
	ManualTaskField             []ManualTask             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL manualTask"`
	ParallelGatewayField        []ParallelGateway        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL parallelGateway"`
	ReceiveTaskField            []ReceiveTask            `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL receiveTask"`
	ScriptTaskField             []ScriptTask             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL scriptTask"`
	SendTaskField               []SendTask               `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL sendTask"`
	SequenceFlowField           []SequenceFlow           `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL sequenceFlow"`
	ServiceTaskField            []ServiceTask            `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL serviceTask"`
	StartEventField             []StartEvent             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL startEvent"`
	SubChoreographyField        []SubChoreography        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL subChoreography"`
	SubProcessField             []SubProcess             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL subProcess"`
	TaskField                   []Task                   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL task"`
	TransactionField            []Transaction            `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL transaction"`
	UserTaskField               []UserTask               `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL userTask"`
	AssociationField            []Association            `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL association"`
	GroupField                  []Group                  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL group"`
	TextAnnotationField         []TextAnnotation         `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL textAnnotation"`
	TextPayloadField            *Payload                 `xml:",chardata"`
}

func DefaultSubChoreography() SubChoreography {
	return SubChoreography{
		ChoreographyActivity: DefaultChoreographyActivity(),
	}
}

type SubChoreographyInterface interface {
	Element
	ChoreographyActivityInterface
	AdHocSubProcesses() (result *[]AdHocSubProcess)
	BoundaryEvents() (result *[]BoundaryEvent)
	BusinessRuleTasks() (result *[]BusinessRuleTask)
	CallActivities() (result *[]CallActivity)
	CallChoreographies() (result *[]CallChoreography)
	ChoreographyTasks() (result *[]ChoreographyTask)
	ComplexGateways() (result *[]ComplexGateway)
	DataObjects() (result *[]DataObject)
	DataObjectReferences() (result *[]DataObjectReference)
	DataStoreReferences() (result *[]DataStoreReference)
	EndEvents() (result *[]EndEvent)
	Events() (result *[]Event)
	EventBasedGateways() (result *[]EventBasedGateway)
	ExclusiveGateways() (result *[]ExclusiveGateway)
	ImplicitThrowEvents() (result *[]ImplicitThrowEvent)
	InclusiveGateways() (result *[]InclusiveGateway)
	IntermediateCatchEvents() (result *[]IntermediateCatchEvent)
	IntermediateThrowEvents() (result *[]IntermediateThrowEvent)
	ManualTasks() (result *[]ManualTask)
	ParallelGateways() (result *[]ParallelGateway)
	ReceiveTasks() (result *[]ReceiveTask)
	ScriptTasks() (result *[]ScriptTask)
	SendTasks() (result *[]SendTask)
	SequenceFlows() (result *[]SequenceFlow)
	ServiceTasks() (result *[]ServiceTask)
	StartEvents() (result *[]StartEvent)
	SubChoreographies() (result *[]SubChoreography)
	SubProcesses() (result *[]SubProcess)
	Tasks() (result *[]Task)
	Transactions() (result *[]Transaction)
	UserTasks() (result *[]UserTask)
	Associations() (result *[]Association)
	Groups() (result *[]Group)
	TextAnnotations() (result *[]TextAnnotation)
	FlowElements() []FlowElementInterface
	Artifacts() []ArtifactInterface
	SetAdHocSubProcesses(value []AdHocSubProcess)
	SetBoundaryEvents(value []BoundaryEvent)
	SetBusinessRuleTasks(value []BusinessRuleTask)
	SetCallActivities(value []CallActivity)
	SetCallChoreographies(value []CallChoreography)
	SetChoreographyTasks(value []ChoreographyTask)
	SetComplexGateways(value []ComplexGateway)
	SetDataObjects(value []DataObject)
	SetDataObjectReferences(value []DataObjectReference)
	SetDataStoreReferences(value []DataStoreReference)
	SetEndEvents(value []EndEvent)
	SetEvents(value []Event)
	SetEventBasedGateways(value []EventBasedGateway)
	SetExclusiveGateways(value []ExclusiveGateway)
	SetImplicitThrowEvents(value []ImplicitThrowEvent)
	SetInclusiveGateways(value []InclusiveGateway)
	SetIntermediateCatchEvents(value []IntermediateCatchEvent)
	SetIntermediateThrowEvents(value []IntermediateThrowEvent)
	SetManualTasks(value []ManualTask)
	SetParallelGateways(value []ParallelGateway)
	SetReceiveTasks(value []ReceiveTask)
	SetScriptTasks(value []ScriptTask)
	SetSendTasks(value []SendTask)
	SetSequenceFlows(value []SequenceFlow)
	SetServiceTasks(value []ServiceTask)
	SetStartEvents(value []StartEvent)
	SetSubChoreographies(value []SubChoreography)
	SetSubProcesses(value []SubProcess)
	SetTasks(value []Task)
	SetTransactions(value []Transaction)
	SetUserTasks(value []UserTask)
	SetAssociations(value []Association)
	SetGroups(value []Group)
	SetTextAnnotations(value []TextAnnotation)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *SubChoreography) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *SubChoreography) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *SubChoreography) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.ChoreographyActivity.FindBy(f); found {
		return
	}

	for i := range t.AdHocSubProcessField {
		if result, found = t.AdHocSubProcessField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.BoundaryEventField {
		if result, found = t.BoundaryEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.BusinessRuleTaskField {
		if result, found = t.BusinessRuleTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.CallActivityField {
		if result, found = t.CallActivityField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.CallChoreographyField {
		if result, found = t.CallChoreographyField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ChoreographyTaskField {
		if result, found = t.ChoreographyTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ComplexGatewayField {
		if result, found = t.ComplexGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataObjectField {
		if result, found = t.DataObjectField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataObjectReferenceField {
		if result, found = t.DataObjectReferenceField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataStoreReferenceField {
		if result, found = t.DataStoreReferenceField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.EndEventField {
		if result, found = t.EndEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.EventField {
		if result, found = t.EventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.EventBasedGatewayField {
		if result, found = t.EventBasedGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ExclusiveGatewayField {
		if result, found = t.ExclusiveGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ImplicitThrowEventField {
		if result, found = t.ImplicitThrowEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.InclusiveGatewayField {
		if result, found = t.InclusiveGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.IntermediateCatchEventField {
		if result, found = t.IntermediateCatchEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.IntermediateThrowEventField {
		if result, found = t.IntermediateThrowEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ManualTaskField {
		if result, found = t.ManualTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ParallelGatewayField {
		if result, found = t.ParallelGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ReceiveTaskField {
		if result, found = t.ReceiveTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ScriptTaskField {
		if result, found = t.ScriptTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SendTaskField {
		if result, found = t.SendTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SequenceFlowField {
		if result, found = t.SequenceFlowField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ServiceTaskField {
		if result, found = t.ServiceTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.StartEventField {
		if result, found = t.StartEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SubChoreographyField {
		if result, found = t.SubChoreographyField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SubProcessField {
		if result, found = t.SubProcessField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.TaskField {
		if result, found = t.TaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.TransactionField {
		if result, found = t.TransactionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.UserTaskField {
		if result, found = t.UserTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.AssociationField {
		if result, found = t.AssociationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.GroupField {
		if result, found = t.GroupField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.TextAnnotationField {
		if result, found = t.TextAnnotationField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *SubChoreography) FlowElements() []FlowElementInterface {

	result := make([]FlowElementInterface, 0)

	for i := range t.AdHocSubProcessField {
		result = append(result, &t.AdHocSubProcessField[i])
	}

	for i := range t.BoundaryEventField {
		result = append(result, &t.BoundaryEventField[i])
	}

	for i := range t.BusinessRuleTaskField {
		result = append(result, &t.BusinessRuleTaskField[i])
	}

	for i := range t.CallActivityField {
		result = append(result, &t.CallActivityField[i])
	}

	for i := range t.CallChoreographyField {
		result = append(result, &t.CallChoreographyField[i])
	}

	for i := range t.ChoreographyTaskField {
		result = append(result, &t.ChoreographyTaskField[i])
	}

	for i := range t.ComplexGatewayField {
		result = append(result, &t.ComplexGatewayField[i])
	}

	for i := range t.DataObjectField {
		result = append(result, &t.DataObjectField[i])
	}

	for i := range t.DataObjectReferenceField {
		result = append(result, &t.DataObjectReferenceField[i])
	}

	for i := range t.DataStoreReferenceField {
		result = append(result, &t.DataStoreReferenceField[i])
	}

	for i := range t.EndEventField {
		result = append(result, &t.EndEventField[i])
	}

	for i := range t.EventField {
		result = append(result, &t.EventField[i])
	}

	for i := range t.EventBasedGatewayField {
		result = append(result, &t.EventBasedGatewayField[i])
	}

	for i := range t.ExclusiveGatewayField {
		result = append(result, &t.ExclusiveGatewayField[i])
	}

	for i := range t.ImplicitThrowEventField {
		result = append(result, &t.ImplicitThrowEventField[i])
	}

	for i := range t.InclusiveGatewayField {
		result = append(result, &t.InclusiveGatewayField[i])
	}

	for i := range t.IntermediateCatchEventField {
		result = append(result, &t.IntermediateCatchEventField[i])
	}

	for i := range t.IntermediateThrowEventField {
		result = append(result, &t.IntermediateThrowEventField[i])
	}

	for i := range t.ManualTaskField {
		result = append(result, &t.ManualTaskField[i])
	}

	for i := range t.ParallelGatewayField {
		result = append(result, &t.ParallelGatewayField[i])
	}

	for i := range t.ReceiveTaskField {
		result = append(result, &t.ReceiveTaskField[i])
	}

	for i := range t.ScriptTaskField {
		result = append(result, &t.ScriptTaskField[i])
	}

	for i := range t.SendTaskField {
		result = append(result, &t.SendTaskField[i])
	}

	for i := range t.SequenceFlowField {
		result = append(result, &t.SequenceFlowField[i])
	}

	for i := range t.ServiceTaskField {
		result = append(result, &t.ServiceTaskField[i])
	}

	for i := range t.StartEventField {
		result = append(result, &t.StartEventField[i])
	}

	for i := range t.SubChoreographyField {
		result = append(result, &t.SubChoreographyField[i])
	}

	for i := range t.SubProcessField {
		result = append(result, &t.SubProcessField[i])
	}

	for i := range t.TaskField {
		result = append(result, &t.TaskField[i])
	}

	for i := range t.TransactionField {
		result = append(result, &t.TransactionField[i])
	}

	for i := range t.UserTaskField {
		result = append(result, &t.UserTaskField[i])
	}
	return result
}
func (t *SubChoreography) Artifacts() []ArtifactInterface {

	result := make([]ArtifactInterface, 0)

	for i := range t.AssociationField {
		result = append(result, &t.AssociationField[i])
	}

	for i := range t.GroupField {
		result = append(result, &t.GroupField[i])
	}

	for i := range t.TextAnnotationField {
		result = append(result, &t.TextAnnotationField[i])
	}
	return result
}
func (t *SubChoreography) AdHocSubProcesses() (result *[]AdHocSubProcess) {
	result = &t.AdHocSubProcessField
	return
}
func (t *SubChoreography) SetAdHocSubProcesses(value []AdHocSubProcess) {
	t.AdHocSubProcessField = value
}
func (t *SubChoreography) BoundaryEvents() (result *[]BoundaryEvent) {
	result = &t.BoundaryEventField
	return
}
func (t *SubChoreography) SetBoundaryEvents(value []BoundaryEvent) {
	t.BoundaryEventField = value
}
func (t *SubChoreography) BusinessRuleTasks() (result *[]BusinessRuleTask) {
	result = &t.BusinessRuleTaskField
	return
}
func (t *SubChoreography) SetBusinessRuleTasks(value []BusinessRuleTask) {
	t.BusinessRuleTaskField = value
}
func (t *SubChoreography) CallActivities() (result *[]CallActivity) {
	result = &t.CallActivityField
	return
}
func (t *SubChoreography) SetCallActivities(value []CallActivity) {
	t.CallActivityField = value
}
func (t *SubChoreography) CallChoreographies() (result *[]CallChoreography) {
	result = &t.CallChoreographyField
	return
}
func (t *SubChoreography) SetCallChoreographies(value []CallChoreography) {
	t.CallChoreographyField = value
}
func (t *SubChoreography) ChoreographyTasks() (result *[]ChoreographyTask) {
	result = &t.ChoreographyTaskField
	return
}
func (t *SubChoreography) SetChoreographyTasks(value []ChoreographyTask) {
	t.ChoreographyTaskField = value
}
func (t *SubChoreography) ComplexGateways() (result *[]ComplexGateway) {
	result = &t.ComplexGatewayField
	return
}
func (t *SubChoreography) SetComplexGateways(value []ComplexGateway) {
	t.ComplexGatewayField = value
}
func (t *SubChoreography) DataObjects() (result *[]DataObject) {
	result = &t.DataObjectField
	return
}
func (t *SubChoreography) SetDataObjects(value []DataObject) {
	t.DataObjectField = value
}
func (t *SubChoreography) DataObjectReferences() (result *[]DataObjectReference) {
	result = &t.DataObjectReferenceField
	return
}
func (t *SubChoreography) SetDataObjectReferences(value []DataObjectReference) {
	t.DataObjectReferenceField = value
}
func (t *SubChoreography) DataStoreReferences() (result *[]DataStoreReference) {
	result = &t.DataStoreReferenceField
	return
}
func (t *SubChoreography) SetDataStoreReferences(value []DataStoreReference) {
	t.DataStoreReferenceField = value
}
func (t *SubChoreography) EndEvents() (result *[]EndEvent) {
	result = &t.EndEventField
	return
}
func (t *SubChoreography) SetEndEvents(value []EndEvent) {
	t.EndEventField = value
}
func (t *SubChoreography) Events() (result *[]Event) {
	result = &t.EventField
	return
}
func (t *SubChoreography) SetEvents(value []Event) {
	t.EventField = value
}
func (t *SubChoreography) EventBasedGateways() (result *[]EventBasedGateway) {
	result = &t.EventBasedGatewayField
	return
}
func (t *SubChoreography) SetEventBasedGateways(value []EventBasedGateway) {
	t.EventBasedGatewayField = value
}
func (t *SubChoreography) ExclusiveGateways() (result *[]ExclusiveGateway) {
	result = &t.ExclusiveGatewayField
	return
}
func (t *SubChoreography) SetExclusiveGateways(value []ExclusiveGateway) {
	t.ExclusiveGatewayField = value
}
func (t *SubChoreography) ImplicitThrowEvents() (result *[]ImplicitThrowEvent) {
	result = &t.ImplicitThrowEventField
	return
}
func (t *SubChoreography) SetImplicitThrowEvents(value []ImplicitThrowEvent) {
	t.ImplicitThrowEventField = value
}
func (t *SubChoreography) InclusiveGateways() (result *[]InclusiveGateway) {
	result = &t.InclusiveGatewayField
	return
}
func (t *SubChoreography) SetInclusiveGateways(value []InclusiveGateway) {
	t.InclusiveGatewayField = value
}
func (t *SubChoreography) IntermediateCatchEvents() (result *[]IntermediateCatchEvent) {
	result = &t.IntermediateCatchEventField
	return
}
func (t *SubChoreography) SetIntermediateCatchEvents(value []IntermediateCatchEvent) {
	t.IntermediateCatchEventField = value
}
func (t *SubChoreography) IntermediateThrowEvents() (result *[]IntermediateThrowEvent) {
	result = &t.IntermediateThrowEventField
	return
}
func (t *SubChoreography) SetIntermediateThrowEvents(value []IntermediateThrowEvent) {
	t.IntermediateThrowEventField = value
}
func (t *SubChoreography) ManualTasks() (result *[]ManualTask) {
	result = &t.ManualTaskField
	return
}
func (t *SubChoreography) SetManualTasks(value []ManualTask) {
	t.ManualTaskField = value
}
func (t *SubChoreography) ParallelGateways() (result *[]ParallelGateway) {
	result = &t.ParallelGatewayField
	return
}
func (t *SubChoreography) SetParallelGateways(value []ParallelGateway) {
	t.ParallelGatewayField = value
}
func (t *SubChoreography) ReceiveTasks() (result *[]ReceiveTask) {
	result = &t.ReceiveTaskField
	return
}
func (t *SubChoreography) SetReceiveTasks(value []ReceiveTask) {
	t.ReceiveTaskField = value
}
func (t *SubChoreography) ScriptTasks() (result *[]ScriptTask) {
	result = &t.ScriptTaskField
	return
}
func (t *SubChoreography) SetScriptTasks(value []ScriptTask) {
	t.ScriptTaskField = value
}
func (t *SubChoreography) SendTasks() (result *[]SendTask) {
	result = &t.SendTaskField
	return
}
func (t *SubChoreography) SetSendTasks(value []SendTask) {
	t.SendTaskField = value
}
func (t *SubChoreography) SequenceFlows() (result *[]SequenceFlow) {
	result = &t.SequenceFlowField
	return
}
func (t *SubChoreography) SetSequenceFlows(value []SequenceFlow) {
	t.SequenceFlowField = value
}
func (t *SubChoreography) ServiceTasks() (result *[]ServiceTask) {
	result = &t.ServiceTaskField
	return
}
func (t *SubChoreography) SetServiceTasks(value []ServiceTask) {
	t.ServiceTaskField = value
}
func (t *SubChoreography) StartEvents() (result *[]StartEvent) {
	result = &t.StartEventField
	return
}
func (t *SubChoreography) SetStartEvents(value []StartEvent) {
	t.StartEventField = value
}
func (t *SubChoreography) SubChoreographies() (result *[]SubChoreography) {
	result = &t.SubChoreographyField
	return
}
func (t *SubChoreography) SetSubChoreographies(value []SubChoreography) {
	t.SubChoreographyField = value
}
func (t *SubChoreography) SubProcesses() (result *[]SubProcess) {
	result = &t.SubProcessField
	return
}
func (t *SubChoreography) SetSubProcesses(value []SubProcess) {
	t.SubProcessField = value
}
func (t *SubChoreography) Tasks() (result *[]Task) {
	result = &t.TaskField
	return
}
func (t *SubChoreography) SetTasks(value []Task) {
	t.TaskField = value
}
func (t *SubChoreography) Transactions() (result *[]Transaction) {
	result = &t.TransactionField
	return
}
func (t *SubChoreography) SetTransactions(value []Transaction) {
	t.TransactionField = value
}
func (t *SubChoreography) UserTasks() (result *[]UserTask) {
	result = &t.UserTaskField
	return
}
func (t *SubChoreography) SetUserTasks(value []UserTask) {
	t.UserTaskField = value
}
func (t *SubChoreography) Associations() (result *[]Association) {
	result = &t.AssociationField
	return
}
func (t *SubChoreography) SetAssociations(value []Association) {
	t.AssociationField = value
}
func (t *SubChoreography) Groups() (result *[]Group) {
	result = &t.GroupField
	return
}
func (t *SubChoreography) SetGroups(value []Group) {
	t.GroupField = value
}
func (t *SubChoreography) TextAnnotations() (result *[]TextAnnotation) {
	result = &t.TextAnnotationField
	return
}
func (t *SubChoreography) SetTextAnnotations(value []TextAnnotation) {
	t.TextAnnotationField = value
}

type SubConversation struct {
	ConversationNode
	CallConversationField []CallConversation `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL callConversation"`
	ConversationField     []Conversation     `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL conversation"`
	SubConversationField  []SubConversation  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL subConversation"`
	TextPayloadField      *Payload           `xml:",chardata"`
}

func DefaultSubConversation() SubConversation {
	return SubConversation{
		ConversationNode: DefaultConversationNode(),
	}
}

type SubConversationInterface interface {
	Element
	ConversationNodeInterface
	CallConversations() (result *[]CallConversation)
	Conversations() (result *[]Conversation)
	SubConversations() (result *[]SubConversation)
	ConversationNodes() []ConversationNodeInterface
	SetCallConversations(value []CallConversation)
	SetConversations(value []Conversation)
	SetSubConversations(value []SubConversation)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *SubConversation) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *SubConversation) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *SubConversation) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.ConversationNode.FindBy(f); found {
		return
	}

	for i := range t.CallConversationField {
		if result, found = t.CallConversationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ConversationField {
		if result, found = t.ConversationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SubConversationField {
		if result, found = t.SubConversationField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *SubConversation) ConversationNodes() []ConversationNodeInterface {

	result := make([]ConversationNodeInterface, 0)

	for i := range t.CallConversationField {
		result = append(result, &t.CallConversationField[i])
	}

	for i := range t.ConversationField {
		result = append(result, &t.ConversationField[i])
	}

	for i := range t.SubConversationField {
		result = append(result, &t.SubConversationField[i])
	}
	return result
}
func (t *SubConversation) CallConversations() (result *[]CallConversation) {
	result = &t.CallConversationField
	return
}
func (t *SubConversation) SetCallConversations(value []CallConversation) {
	t.CallConversationField = value
}
func (t *SubConversation) Conversations() (result *[]Conversation) {
	result = &t.ConversationField
	return
}
func (t *SubConversation) SetConversations(value []Conversation) {
	t.ConversationField = value
}
func (t *SubConversation) SubConversations() (result *[]SubConversation) {
	result = &t.SubConversationField
	return
}
func (t *SubConversation) SetSubConversations(value []SubConversation) {
	t.SubConversationField = value
}

type SubProcess struct {
	Activity
	TriggeredByEventField       *bool                    `xml:"triggeredByEvent,attr,omitempty"`
	LaneSetField                []LaneSet                `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL laneSet"`
	AdHocSubProcessField        []AdHocSubProcess        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL adHocSubProcess"`
	BoundaryEventField          []BoundaryEvent          `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL boundaryEvent"`
	BusinessRuleTaskField       []BusinessRuleTask       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL businessRuleTask"`
	CallActivityField           []CallActivity           `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL callActivity"`
	CallChoreographyField       []CallChoreography       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL callChoreography"`
	ChoreographyTaskField       []ChoreographyTask       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL choreographyTask"`
	ComplexGatewayField         []ComplexGateway         `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL complexGateway"`
	DataObjectField             []DataObject             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataObject"`
	DataObjectReferenceField    []DataObjectReference    `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataObjectReference"`
	DataStoreReferenceField     []DataStoreReference     `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataStoreReference"`
	EndEventField               []EndEvent               `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL endEvent"`
	EventField                  []Event                  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL event"`
	EventBasedGatewayField      []EventBasedGateway      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL eventBasedGateway"`
	ExclusiveGatewayField       []ExclusiveGateway       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL exclusiveGateway"`
	ImplicitThrowEventField     []ImplicitThrowEvent     `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL implicitThrowEvent"`
	InclusiveGatewayField       []InclusiveGateway       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL inclusiveGateway"`
	IntermediateCatchEventField []IntermediateCatchEvent `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL intermediateCatchEvent"`
	IntermediateThrowEventField []IntermediateThrowEvent `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL intermediateThrowEvent"`
	ManualTaskField             []ManualTask             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL manualTask"`
	ParallelGatewayField        []ParallelGateway        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL parallelGateway"`
	ReceiveTaskField            []ReceiveTask            `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL receiveTask"`
	ScriptTaskField             []ScriptTask             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL scriptTask"`
	SendTaskField               []SendTask               `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL sendTask"`
	SequenceFlowField           []SequenceFlow           `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL sequenceFlow"`
	ServiceTaskField            []ServiceTask            `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL serviceTask"`
	StartEventField             []StartEvent             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL startEvent"`
	SubChoreographyField        []SubChoreography        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL subChoreography"`
	SubProcessField             []SubProcess             `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL subProcess"`
	TaskField                   []Task                   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL task"`
	TransactionField            []Transaction            `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL transaction"`
	UserTaskField               []UserTask               `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL userTask"`
	AssociationField            []Association            `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL association"`
	GroupField                  []Group                  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL group"`
	TextAnnotationField         []TextAnnotation         `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL textAnnotation"`
	TextPayloadField            *Payload                 `xml:",chardata"`
}

var defaultSubProcessTriggeredByEventField bool = false

func DefaultSubProcess() SubProcess {
	return SubProcess{
		Activity:              DefaultActivity(),
		TriggeredByEventField: &defaultSubProcessTriggeredByEventField,
	}
}

type SubProcessInterface interface {
	Element
	ActivityInterface
	TriggeredByEvent() (result bool)
	LaneSets() (result *[]LaneSet)
	AdHocSubProcesses() (result *[]AdHocSubProcess)
	BoundaryEvents() (result *[]BoundaryEvent)
	BusinessRuleTasks() (result *[]BusinessRuleTask)
	CallActivities() (result *[]CallActivity)
	CallChoreographies() (result *[]CallChoreography)
	ChoreographyTasks() (result *[]ChoreographyTask)
	ComplexGateways() (result *[]ComplexGateway)
	DataObjects() (result *[]DataObject)
	DataObjectReferences() (result *[]DataObjectReference)
	DataStoreReferences() (result *[]DataStoreReference)
	EndEvents() (result *[]EndEvent)
	Events() (result *[]Event)
	EventBasedGateways() (result *[]EventBasedGateway)
	ExclusiveGateways() (result *[]ExclusiveGateway)
	ImplicitThrowEvents() (result *[]ImplicitThrowEvent)
	InclusiveGateways() (result *[]InclusiveGateway)
	IntermediateCatchEvents() (result *[]IntermediateCatchEvent)
	IntermediateThrowEvents() (result *[]IntermediateThrowEvent)
	ManualTasks() (result *[]ManualTask)
	ParallelGateways() (result *[]ParallelGateway)
	ReceiveTasks() (result *[]ReceiveTask)
	ScriptTasks() (result *[]ScriptTask)
	SendTasks() (result *[]SendTask)
	SequenceFlows() (result *[]SequenceFlow)
	ServiceTasks() (result *[]ServiceTask)
	StartEvents() (result *[]StartEvent)
	SubChoreographies() (result *[]SubChoreography)
	SubProcesses() (result *[]SubProcess)
	Tasks() (result *[]Task)
	Transactions() (result *[]Transaction)
	UserTasks() (result *[]UserTask)
	Associations() (result *[]Association)
	Groups() (result *[]Group)
	TextAnnotations() (result *[]TextAnnotation)
	FlowElements() []FlowElementInterface
	Artifacts() []ArtifactInterface
	SetTriggeredByEvent(value *bool)
	SetLaneSets(value []LaneSet)
	SetAdHocSubProcesses(value []AdHocSubProcess)
	SetBoundaryEvents(value []BoundaryEvent)
	SetBusinessRuleTasks(value []BusinessRuleTask)
	SetCallActivities(value []CallActivity)
	SetCallChoreographies(value []CallChoreography)
	SetChoreographyTasks(value []ChoreographyTask)
	SetComplexGateways(value []ComplexGateway)
	SetDataObjects(value []DataObject)
	SetDataObjectReferences(value []DataObjectReference)
	SetDataStoreReferences(value []DataStoreReference)
	SetEndEvents(value []EndEvent)
	SetEvents(value []Event)
	SetEventBasedGateways(value []EventBasedGateway)
	SetExclusiveGateways(value []ExclusiveGateway)
	SetImplicitThrowEvents(value []ImplicitThrowEvent)
	SetInclusiveGateways(value []InclusiveGateway)
	SetIntermediateCatchEvents(value []IntermediateCatchEvent)
	SetIntermediateThrowEvents(value []IntermediateThrowEvent)
	SetManualTasks(value []ManualTask)
	SetParallelGateways(value []ParallelGateway)
	SetReceiveTasks(value []ReceiveTask)
	SetScriptTasks(value []ScriptTask)
	SetSendTasks(value []SendTask)
	SetSequenceFlows(value []SequenceFlow)
	SetServiceTasks(value []ServiceTask)
	SetStartEvents(value []StartEvent)
	SetSubChoreographies(value []SubChoreography)
	SetSubProcesses(value []SubProcess)
	SetTasks(value []Task)
	SetTransactions(value []Transaction)
	SetUserTasks(value []UserTask)
	SetAssociations(value []Association)
	SetGroups(value []Group)
	SetTextAnnotations(value []TextAnnotation)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *SubProcess) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *SubProcess) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *SubProcess) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Activity.FindBy(f); found {
		return
	}

	for i := range t.LaneSetField {
		if result, found = t.LaneSetField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.AdHocSubProcessField {
		if result, found = t.AdHocSubProcessField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.BoundaryEventField {
		if result, found = t.BoundaryEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.BusinessRuleTaskField {
		if result, found = t.BusinessRuleTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.CallActivityField {
		if result, found = t.CallActivityField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.CallChoreographyField {
		if result, found = t.CallChoreographyField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ChoreographyTaskField {
		if result, found = t.ChoreographyTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ComplexGatewayField {
		if result, found = t.ComplexGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataObjectField {
		if result, found = t.DataObjectField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataObjectReferenceField {
		if result, found = t.DataObjectReferenceField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataStoreReferenceField {
		if result, found = t.DataStoreReferenceField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.EndEventField {
		if result, found = t.EndEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.EventField {
		if result, found = t.EventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.EventBasedGatewayField {
		if result, found = t.EventBasedGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ExclusiveGatewayField {
		if result, found = t.ExclusiveGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ImplicitThrowEventField {
		if result, found = t.ImplicitThrowEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.InclusiveGatewayField {
		if result, found = t.InclusiveGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.IntermediateCatchEventField {
		if result, found = t.IntermediateCatchEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.IntermediateThrowEventField {
		if result, found = t.IntermediateThrowEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ManualTaskField {
		if result, found = t.ManualTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ParallelGatewayField {
		if result, found = t.ParallelGatewayField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ReceiveTaskField {
		if result, found = t.ReceiveTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ScriptTaskField {
		if result, found = t.ScriptTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SendTaskField {
		if result, found = t.SendTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SequenceFlowField {
		if result, found = t.SequenceFlowField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ServiceTaskField {
		if result, found = t.ServiceTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.StartEventField {
		if result, found = t.StartEventField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SubChoreographyField {
		if result, found = t.SubChoreographyField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SubProcessField {
		if result, found = t.SubProcessField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.TaskField {
		if result, found = t.TaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.TransactionField {
		if result, found = t.TransactionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.UserTaskField {
		if result, found = t.UserTaskField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.AssociationField {
		if result, found = t.AssociationField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.GroupField {
		if result, found = t.GroupField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.TextAnnotationField {
		if result, found = t.TextAnnotationField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *SubProcess) TriggeredByEvent() (result bool) {
	if t.TriggeredByEventField == nil {
		result = defaultSubProcessTriggeredByEventField
		return
	}
	result = *t.TriggeredByEventField
	return
}
func (t *SubProcess) SetTriggeredByEvent(value *bool) {
	t.TriggeredByEventField = value
}
func (t *SubProcess) FlowElements() []FlowElementInterface {

	result := make([]FlowElementInterface, 0)

	for i := range t.AdHocSubProcessField {
		result = append(result, &t.AdHocSubProcessField[i])
	}

	for i := range t.BoundaryEventField {
		result = append(result, &t.BoundaryEventField[i])
	}

	for i := range t.BusinessRuleTaskField {
		result = append(result, &t.BusinessRuleTaskField[i])
	}

	for i := range t.CallActivityField {
		result = append(result, &t.CallActivityField[i])
	}

	for i := range t.CallChoreographyField {
		result = append(result, &t.CallChoreographyField[i])
	}

	for i := range t.ChoreographyTaskField {
		result = append(result, &t.ChoreographyTaskField[i])
	}

	for i := range t.ComplexGatewayField {
		result = append(result, &t.ComplexGatewayField[i])
	}

	for i := range t.DataObjectField {
		result = append(result, &t.DataObjectField[i])
	}

	for i := range t.DataObjectReferenceField {
		result = append(result, &t.DataObjectReferenceField[i])
	}

	for i := range t.DataStoreReferenceField {
		result = append(result, &t.DataStoreReferenceField[i])
	}

	for i := range t.EndEventField {
		result = append(result, &t.EndEventField[i])
	}

	for i := range t.EventField {
		result = append(result, &t.EventField[i])
	}

	for i := range t.EventBasedGatewayField {
		result = append(result, &t.EventBasedGatewayField[i])
	}

	for i := range t.ExclusiveGatewayField {
		result = append(result, &t.ExclusiveGatewayField[i])
	}

	for i := range t.ImplicitThrowEventField {
		result = append(result, &t.ImplicitThrowEventField[i])
	}

	for i := range t.InclusiveGatewayField {
		result = append(result, &t.InclusiveGatewayField[i])
	}

	for i := range t.IntermediateCatchEventField {
		result = append(result, &t.IntermediateCatchEventField[i])
	}

	for i := range t.IntermediateThrowEventField {
		result = append(result, &t.IntermediateThrowEventField[i])
	}

	for i := range t.ManualTaskField {
		result = append(result, &t.ManualTaskField[i])
	}

	for i := range t.ParallelGatewayField {
		result = append(result, &t.ParallelGatewayField[i])
	}

	for i := range t.ReceiveTaskField {
		result = append(result, &t.ReceiveTaskField[i])
	}

	for i := range t.ScriptTaskField {
		result = append(result, &t.ScriptTaskField[i])
	}

	for i := range t.SendTaskField {
		result = append(result, &t.SendTaskField[i])
	}

	for i := range t.SequenceFlowField {
		result = append(result, &t.SequenceFlowField[i])
	}

	for i := range t.ServiceTaskField {
		result = append(result, &t.ServiceTaskField[i])
	}

	for i := range t.StartEventField {
		result = append(result, &t.StartEventField[i])
	}

	for i := range t.SubChoreographyField {
		result = append(result, &t.SubChoreographyField[i])
	}

	for i := range t.SubProcessField {
		result = append(result, &t.SubProcessField[i])
	}

	for i := range t.TaskField {
		result = append(result, &t.TaskField[i])
	}

	for i := range t.TransactionField {
		result = append(result, &t.TransactionField[i])
	}

	for i := range t.UserTaskField {
		result = append(result, &t.UserTaskField[i])
	}
	return result
}
func (t *SubProcess) Artifacts() []ArtifactInterface {

	result := make([]ArtifactInterface, 0)

	for i := range t.AssociationField {
		result = append(result, &t.AssociationField[i])
	}

	for i := range t.GroupField {
		result = append(result, &t.GroupField[i])
	}

	for i := range t.TextAnnotationField {
		result = append(result, &t.TextAnnotationField[i])
	}
	return result
}
func (t *SubProcess) LaneSets() (result *[]LaneSet) {
	result = &t.LaneSetField
	return
}
func (t *SubProcess) SetLaneSets(value []LaneSet) {
	t.LaneSetField = value
}
func (t *SubProcess) AdHocSubProcesses() (result *[]AdHocSubProcess) {
	result = &t.AdHocSubProcessField
	return
}
func (t *SubProcess) SetAdHocSubProcesses(value []AdHocSubProcess) {
	t.AdHocSubProcessField = value
}
func (t *SubProcess) BoundaryEvents() (result *[]BoundaryEvent) {
	result = &t.BoundaryEventField
	return
}
func (t *SubProcess) SetBoundaryEvents(value []BoundaryEvent) {
	t.BoundaryEventField = value
}
func (t *SubProcess) BusinessRuleTasks() (result *[]BusinessRuleTask) {
	result = &t.BusinessRuleTaskField
	return
}
func (t *SubProcess) SetBusinessRuleTasks(value []BusinessRuleTask) {
	t.BusinessRuleTaskField = value
}
func (t *SubProcess) CallActivities() (result *[]CallActivity) {
	result = &t.CallActivityField
	return
}
func (t *SubProcess) SetCallActivities(value []CallActivity) {
	t.CallActivityField = value
}
func (t *SubProcess) CallChoreographies() (result *[]CallChoreography) {
	result = &t.CallChoreographyField
	return
}
func (t *SubProcess) SetCallChoreographies(value []CallChoreography) {
	t.CallChoreographyField = value
}
func (t *SubProcess) ChoreographyTasks() (result *[]ChoreographyTask) {
	result = &t.ChoreographyTaskField
	return
}
func (t *SubProcess) SetChoreographyTasks(value []ChoreographyTask) {
	t.ChoreographyTaskField = value
}
func (t *SubProcess) ComplexGateways() (result *[]ComplexGateway) {
	result = &t.ComplexGatewayField
	return
}
func (t *SubProcess) SetComplexGateways(value []ComplexGateway) {
	t.ComplexGatewayField = value
}
func (t *SubProcess) DataObjects() (result *[]DataObject) {
	result = &t.DataObjectField
	return
}
func (t *SubProcess) SetDataObjects(value []DataObject) {
	t.DataObjectField = value
}
func (t *SubProcess) DataObjectReferences() (result *[]DataObjectReference) {
	result = &t.DataObjectReferenceField
	return
}
func (t *SubProcess) SetDataObjectReferences(value []DataObjectReference) {
	t.DataObjectReferenceField = value
}
func (t *SubProcess) DataStoreReferences() (result *[]DataStoreReference) {
	result = &t.DataStoreReferenceField
	return
}
func (t *SubProcess) SetDataStoreReferences(value []DataStoreReference) {
	t.DataStoreReferenceField = value
}
func (t *SubProcess) EndEvents() (result *[]EndEvent) {
	result = &t.EndEventField
	return
}
func (t *SubProcess) SetEndEvents(value []EndEvent) {
	t.EndEventField = value
}
func (t *SubProcess) Events() (result *[]Event) {
	result = &t.EventField
	return
}
func (t *SubProcess) SetEvents(value []Event) {
	t.EventField = value
}
func (t *SubProcess) EventBasedGateways() (result *[]EventBasedGateway) {
	result = &t.EventBasedGatewayField
	return
}
func (t *SubProcess) SetEventBasedGateways(value []EventBasedGateway) {
	t.EventBasedGatewayField = value
}
func (t *SubProcess) ExclusiveGateways() (result *[]ExclusiveGateway) {
	result = &t.ExclusiveGatewayField
	return
}
func (t *SubProcess) SetExclusiveGateways(value []ExclusiveGateway) {
	t.ExclusiveGatewayField = value
}
func (t *SubProcess) ImplicitThrowEvents() (result *[]ImplicitThrowEvent) {
	result = &t.ImplicitThrowEventField
	return
}
func (t *SubProcess) SetImplicitThrowEvents(value []ImplicitThrowEvent) {
	t.ImplicitThrowEventField = value
}
func (t *SubProcess) InclusiveGateways() (result *[]InclusiveGateway) {
	result = &t.InclusiveGatewayField
	return
}
func (t *SubProcess) SetInclusiveGateways(value []InclusiveGateway) {
	t.InclusiveGatewayField = value
}
func (t *SubProcess) IntermediateCatchEvents() (result *[]IntermediateCatchEvent) {
	result = &t.IntermediateCatchEventField
	return
}
func (t *SubProcess) SetIntermediateCatchEvents(value []IntermediateCatchEvent) {
	t.IntermediateCatchEventField = value
}
func (t *SubProcess) IntermediateThrowEvents() (result *[]IntermediateThrowEvent) {
	result = &t.IntermediateThrowEventField
	return
}
func (t *SubProcess) SetIntermediateThrowEvents(value []IntermediateThrowEvent) {
	t.IntermediateThrowEventField = value
}
func (t *SubProcess) ManualTasks() (result *[]ManualTask) {
	result = &t.ManualTaskField
	return
}
func (t *SubProcess) SetManualTasks(value []ManualTask) {
	t.ManualTaskField = value
}
func (t *SubProcess) ParallelGateways() (result *[]ParallelGateway) {
	result = &t.ParallelGatewayField
	return
}
func (t *SubProcess) SetParallelGateways(value []ParallelGateway) {
	t.ParallelGatewayField = value
}
func (t *SubProcess) ReceiveTasks() (result *[]ReceiveTask) {
	result = &t.ReceiveTaskField
	return
}
func (t *SubProcess) SetReceiveTasks(value []ReceiveTask) {
	t.ReceiveTaskField = value
}
func (t *SubProcess) ScriptTasks() (result *[]ScriptTask) {
	result = &t.ScriptTaskField
	return
}
func (t *SubProcess) SetScriptTasks(value []ScriptTask) {
	t.ScriptTaskField = value
}
func (t *SubProcess) SendTasks() (result *[]SendTask) {
	result = &t.SendTaskField
	return
}
func (t *SubProcess) SetSendTasks(value []SendTask) {
	t.SendTaskField = value
}
func (t *SubProcess) SequenceFlows() (result *[]SequenceFlow) {
	result = &t.SequenceFlowField
	return
}
func (t *SubProcess) SetSequenceFlows(value []SequenceFlow) {
	t.SequenceFlowField = value
}
func (t *SubProcess) ServiceTasks() (result *[]ServiceTask) {
	result = &t.ServiceTaskField
	return
}
func (t *SubProcess) SetServiceTasks(value []ServiceTask) {
	t.ServiceTaskField = value
}
func (t *SubProcess) StartEvents() (result *[]StartEvent) {
	result = &t.StartEventField
	return
}
func (t *SubProcess) SetStartEvents(value []StartEvent) {
	t.StartEventField = value
}
func (t *SubProcess) SubChoreographies() (result *[]SubChoreography) {
	result = &t.SubChoreographyField
	return
}
func (t *SubProcess) SetSubChoreographies(value []SubChoreography) {
	t.SubChoreographyField = value
}
func (t *SubProcess) SubProcesses() (result *[]SubProcess) {
	result = &t.SubProcessField
	return
}
func (t *SubProcess) SetSubProcesses(value []SubProcess) {
	t.SubProcessField = value
}
func (t *SubProcess) Tasks() (result *[]Task) {
	result = &t.TaskField
	return
}
func (t *SubProcess) SetTasks(value []Task) {
	t.TaskField = value
}
func (t *SubProcess) Transactions() (result *[]Transaction) {
	result = &t.TransactionField
	return
}
func (t *SubProcess) SetTransactions(value []Transaction) {
	t.TransactionField = value
}
func (t *SubProcess) UserTasks() (result *[]UserTask) {
	result = &t.UserTaskField
	return
}
func (t *SubProcess) SetUserTasks(value []UserTask) {
	t.UserTaskField = value
}
func (t *SubProcess) Associations() (result *[]Association) {
	result = &t.AssociationField
	return
}
func (t *SubProcess) SetAssociations(value []Association) {
	t.AssociationField = value
}
func (t *SubProcess) Groups() (result *[]Group) {
	result = &t.GroupField
	return
}
func (t *SubProcess) SetGroups(value []Group) {
	t.GroupField = value
}
func (t *SubProcess) TextAnnotations() (result *[]TextAnnotation) {
	result = &t.TextAnnotationField
	return
}
func (t *SubProcess) SetTextAnnotations(value []TextAnnotation) {
	t.TextAnnotationField = value
}

type Task struct {
	Activity
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultTask() Task {
	return Task{
		Activity: DefaultActivity(),
	}
}

type TaskInterface interface {
	Element
	ActivityInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Task) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Task) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Task) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Activity.FindBy(f); found {
		return
	}

	return
}

type TerminateEventDefinition struct {
	EventDefinition
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultTerminateEventDefinition() TerminateEventDefinition {
	return TerminateEventDefinition{
		EventDefinition: DefaultEventDefinition(),
	}
}

type TerminateEventDefinitionInterface interface {
	Element
	EventDefinitionInterface

	TextPayload() *string
	SetTextPayload(string)
}

func (t *TerminateEventDefinition) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *TerminateEventDefinition) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *TerminateEventDefinition) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.EventDefinition.FindBy(f); found {
		return
	}

	return
}

type TextAnnotation struct {
	Artifact
	TextFormatField  *string  `xml:"textFormat,attr,omitempty"`
	TextField        *Text    `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL text"`
	TextPayloadField *Payload `xml:",chardata"`
}

var defaultTextAnnotationTextFormatField string = "text/plain"

func DefaultTextAnnotation() TextAnnotation {
	return TextAnnotation{
		Artifact:        DefaultArtifact(),
		TextFormatField: &defaultTextAnnotationTextFormatField,
	}
}

type TextAnnotationInterface interface {
	Element
	ArtifactInterface
	TextFormat() (result *string)
	Text() (result *Text, present bool)
	SetTextFormat(value *string)
	SetText(value *Text)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *TextAnnotation) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *TextAnnotation) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *TextAnnotation) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Artifact.FindBy(f); found {
		return
	}

	if value := t.TextField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	return
}
func (t *TextAnnotation) TextFormat() (result *string) {
	if t.TextFormatField == nil {
		result = &defaultTextAnnotationTextFormatField
		return
	}
	result = t.TextFormatField
	return
}
func (t *TextAnnotation) SetTextFormat(value *string) {
	t.TextFormatField = value
}
func (t *TextAnnotation) Text() (result *Text, present bool) {
	if t.TextField != nil {
		present = true
	}
	result = t.TextField
	return
}
func (t *TextAnnotation) SetText(value *Text) {
	t.TextField = value
}

type Text struct {
	TextPayloadField *Payload `xml:",chardata"`
}

func DefaultText() Text {
	return Text{}
}

type TextInterface interface {
	Element

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Text) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Text) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Text) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}

	return
}

type ThrowEvent struct {
	Event
	DataInputField                  []DataInput                  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataInput"`
	DataInputAssociationField       []DataInputAssociation       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL dataInputAssociation"`
	InputSetField                   *InputSet                    `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL inputSet"`
	CancelEventDefinitionField      []CancelEventDefinition      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL cancelEventDefinition"`
	CompensateEventDefinitionField  []CompensateEventDefinition  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL compensateEventDefinition"`
	ConditionalEventDefinitionField []ConditionalEventDefinition `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL conditionalEventDefinition"`
	ErrorEventDefinitionField       []ErrorEventDefinition       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL errorEventDefinition"`
	EscalationEventDefinitionField  []EscalationEventDefinition  `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL escalationEventDefinition"`
	LinkEventDefinitionField        []LinkEventDefinition        `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL linkEventDefinition"`
	MessageEventDefinitionField     []MessageEventDefinition     `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL messageEventDefinition"`
	SignalEventDefinitionField      []SignalEventDefinition      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL signalEventDefinition"`
	TerminateEventDefinitionField   []TerminateEventDefinition   `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL terminateEventDefinition"`
	TimerEventDefinitionField       []TimerEventDefinition       `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL timerEventDefinition"`
	EventDefinitionRefField         []QName                      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL eventDefinitionRef"`
}

func DefaultThrowEvent() ThrowEvent {
	return ThrowEvent{
		Event: DefaultEvent(),
	}
}

type ThrowEventInterface interface {
	Element
	EventInterface
	DataInputs() (result *[]DataInput)
	DataInputAssociations() (result *[]DataInputAssociation)
	InputSet() (result *InputSet, present bool)
	CancelEventDefinitions() (result *[]CancelEventDefinition)
	CompensateEventDefinitions() (result *[]CompensateEventDefinition)
	ConditionalEventDefinitions() (result *[]ConditionalEventDefinition)
	ErrorEventDefinitions() (result *[]ErrorEventDefinition)
	EscalationEventDefinitions() (result *[]EscalationEventDefinition)
	LinkEventDefinitions() (result *[]LinkEventDefinition)
	MessageEventDefinitions() (result *[]MessageEventDefinition)
	SignalEventDefinitions() (result *[]SignalEventDefinition)
	TerminateEventDefinitions() (result *[]TerminateEventDefinition)
	TimerEventDefinitions() (result *[]TimerEventDefinition)
	EventDefinitionRefs() (result *[]QName)
	EventDefinitions() []EventDefinitionInterface
	SetDataInputs(value []DataInput)
	SetDataInputAssociations(value []DataInputAssociation)
	SetInputSet(value *InputSet)
	SetCancelEventDefinitions(value []CancelEventDefinition)
	SetCompensateEventDefinitions(value []CompensateEventDefinition)
	SetConditionalEventDefinitions(value []ConditionalEventDefinition)
	SetErrorEventDefinitions(value []ErrorEventDefinition)
	SetEscalationEventDefinitions(value []EscalationEventDefinition)
	SetLinkEventDefinitions(value []LinkEventDefinition)
	SetMessageEventDefinitions(value []MessageEventDefinition)
	SetSignalEventDefinitions(value []SignalEventDefinition)
	SetTerminateEventDefinitions(value []TerminateEventDefinition)
	SetTimerEventDefinitions(value []TimerEventDefinition)
	SetEventDefinitionRefs(value []QName)
}

func (t *ThrowEvent) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Event.FindBy(f); found {
		return
	}

	for i := range t.DataInputField {
		if result, found = t.DataInputField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.DataInputAssociationField {
		if result, found = t.DataInputAssociationField[i].FindBy(f); found {
			return
		}
	}

	if value := t.InputSetField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	for i := range t.CancelEventDefinitionField {
		if result, found = t.CancelEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.CompensateEventDefinitionField {
		if result, found = t.CompensateEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ConditionalEventDefinitionField {
		if result, found = t.ConditionalEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.ErrorEventDefinitionField {
		if result, found = t.ErrorEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.EscalationEventDefinitionField {
		if result, found = t.EscalationEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.LinkEventDefinitionField {
		if result, found = t.LinkEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.MessageEventDefinitionField {
		if result, found = t.MessageEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.SignalEventDefinitionField {
		if result, found = t.SignalEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.TerminateEventDefinitionField {
		if result, found = t.TerminateEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	for i := range t.TimerEventDefinitionField {
		if result, found = t.TimerEventDefinitionField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *ThrowEvent) EventDefinitions() []EventDefinitionInterface {

	result := make([]EventDefinitionInterface, 0)

	for i := range t.CancelEventDefinitionField {
		result = append(result, &t.CancelEventDefinitionField[i])
	}

	for i := range t.CompensateEventDefinitionField {
		result = append(result, &t.CompensateEventDefinitionField[i])
	}

	for i := range t.ConditionalEventDefinitionField {
		result = append(result, &t.ConditionalEventDefinitionField[i])
	}

	for i := range t.ErrorEventDefinitionField {
		result = append(result, &t.ErrorEventDefinitionField[i])
	}

	for i := range t.EscalationEventDefinitionField {
		result = append(result, &t.EscalationEventDefinitionField[i])
	}

	for i := range t.LinkEventDefinitionField {
		result = append(result, &t.LinkEventDefinitionField[i])
	}

	for i := range t.MessageEventDefinitionField {
		result = append(result, &t.MessageEventDefinitionField[i])
	}

	for i := range t.SignalEventDefinitionField {
		result = append(result, &t.SignalEventDefinitionField[i])
	}

	for i := range t.TerminateEventDefinitionField {
		result = append(result, &t.TerminateEventDefinitionField[i])
	}

	for i := range t.TimerEventDefinitionField {
		result = append(result, &t.TimerEventDefinitionField[i])
	}
	return result
}
func (t *ThrowEvent) DataInputs() (result *[]DataInput) {
	result = &t.DataInputField
	return
}
func (t *ThrowEvent) SetDataInputs(value []DataInput) {
	t.DataInputField = value
}
func (t *ThrowEvent) DataInputAssociations() (result *[]DataInputAssociation) {
	result = &t.DataInputAssociationField
	return
}
func (t *ThrowEvent) SetDataInputAssociations(value []DataInputAssociation) {
	t.DataInputAssociationField = value
}
func (t *ThrowEvent) InputSet() (result *InputSet, present bool) {
	if t.InputSetField != nil {
		present = true
	}
	result = t.InputSetField
	return
}
func (t *ThrowEvent) SetInputSet(value *InputSet) {
	t.InputSetField = value
}
func (t *ThrowEvent) CancelEventDefinitions() (result *[]CancelEventDefinition) {
	result = &t.CancelEventDefinitionField
	return
}
func (t *ThrowEvent) SetCancelEventDefinitions(value []CancelEventDefinition) {
	t.CancelEventDefinitionField = value
}
func (t *ThrowEvent) CompensateEventDefinitions() (result *[]CompensateEventDefinition) {
	result = &t.CompensateEventDefinitionField
	return
}
func (t *ThrowEvent) SetCompensateEventDefinitions(value []CompensateEventDefinition) {
	t.CompensateEventDefinitionField = value
}
func (t *ThrowEvent) ConditionalEventDefinitions() (result *[]ConditionalEventDefinition) {
	result = &t.ConditionalEventDefinitionField
	return
}
func (t *ThrowEvent) SetConditionalEventDefinitions(value []ConditionalEventDefinition) {
	t.ConditionalEventDefinitionField = value
}
func (t *ThrowEvent) ErrorEventDefinitions() (result *[]ErrorEventDefinition) {
	result = &t.ErrorEventDefinitionField
	return
}
func (t *ThrowEvent) SetErrorEventDefinitions(value []ErrorEventDefinition) {
	t.ErrorEventDefinitionField = value
}
func (t *ThrowEvent) EscalationEventDefinitions() (result *[]EscalationEventDefinition) {
	result = &t.EscalationEventDefinitionField
	return
}
func (t *ThrowEvent) SetEscalationEventDefinitions(value []EscalationEventDefinition) {
	t.EscalationEventDefinitionField = value
}
func (t *ThrowEvent) LinkEventDefinitions() (result *[]LinkEventDefinition) {
	result = &t.LinkEventDefinitionField
	return
}
func (t *ThrowEvent) SetLinkEventDefinitions(value []LinkEventDefinition) {
	t.LinkEventDefinitionField = value
}
func (t *ThrowEvent) MessageEventDefinitions() (result *[]MessageEventDefinition) {
	result = &t.MessageEventDefinitionField
	return
}
func (t *ThrowEvent) SetMessageEventDefinitions(value []MessageEventDefinition) {
	t.MessageEventDefinitionField = value
}
func (t *ThrowEvent) SignalEventDefinitions() (result *[]SignalEventDefinition) {
	result = &t.SignalEventDefinitionField
	return
}
func (t *ThrowEvent) SetSignalEventDefinitions(value []SignalEventDefinition) {
	t.SignalEventDefinitionField = value
}
func (t *ThrowEvent) TerminateEventDefinitions() (result *[]TerminateEventDefinition) {
	result = &t.TerminateEventDefinitionField
	return
}
func (t *ThrowEvent) SetTerminateEventDefinitions(value []TerminateEventDefinition) {
	t.TerminateEventDefinitionField = value
}
func (t *ThrowEvent) TimerEventDefinitions() (result *[]TimerEventDefinition) {
	result = &t.TimerEventDefinitionField
	return
}
func (t *ThrowEvent) SetTimerEventDefinitions(value []TimerEventDefinition) {
	t.TimerEventDefinitionField = value
}
func (t *ThrowEvent) EventDefinitionRefs() (result *[]QName) {
	result = &t.EventDefinitionRefField
	return
}
func (t *ThrowEvent) SetEventDefinitionRefs(value []QName) {
	t.EventDefinitionRefField = value
}

type TimerEventDefinition struct {
	EventDefinition
	TimeDateField     *AnExpression `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL timeDate"`
	TimeDurationField *AnExpression `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL timeDuration"`
	TimeCycleField    *AnExpression `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL timeCycle"`
	TextPayloadField  *Payload      `xml:",chardata"`
}

func DefaultTimerEventDefinition() TimerEventDefinition {
	return TimerEventDefinition{
		EventDefinition: DefaultEventDefinition(),
	}
}

type TimerEventDefinitionInterface interface {
	Element
	EventDefinitionInterface
	TimeDate() (result *AnExpression, present bool)
	TimeDuration() (result *AnExpression, present bool)
	TimeCycle() (result *AnExpression, present bool)
	SetTimeDate(value *AnExpression)
	SetTimeDuration(value *AnExpression)
	SetTimeCycle(value *AnExpression)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *TimerEventDefinition) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *TimerEventDefinition) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *TimerEventDefinition) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.EventDefinition.FindBy(f); found {
		return
	}

	if value := t.TimeDateField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	if value := t.TimeDurationField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	if value := t.TimeCycleField; value != nil {
		if result, found = value.FindBy(f); found {
			return
		}
	}

	return
}
func (t *TimerEventDefinition) TimeDate() (result *AnExpression, present bool) {
	if t.TimeDateField != nil {
		present = true
	}
	result = t.TimeDateField
	return
}
func (t *TimerEventDefinition) SetTimeDate(value *AnExpression) {
	t.TimeDateField = value
}
func (t *TimerEventDefinition) TimeDuration() (result *AnExpression, present bool) {
	if t.TimeDurationField != nil {
		present = true
	}
	result = t.TimeDurationField
	return
}
func (t *TimerEventDefinition) SetTimeDuration(value *AnExpression) {
	t.TimeDurationField = value
}
func (t *TimerEventDefinition) TimeCycle() (result *AnExpression, present bool) {
	if t.TimeCycleField != nil {
		present = true
	}
	result = t.TimeCycleField
	return
}
func (t *TimerEventDefinition) SetTimeCycle(value *AnExpression) {
	t.TimeCycleField = value
}

type Transaction struct {
	SubProcess
	MethodField      *TransactionMethod `xml:"method,attr,omitempty"`
	TextPayloadField *Payload           `xml:",chardata"`
}

var defaultTransactionMethodField TransactionMethod = "##Compensate"

func DefaultTransaction() Transaction {
	return Transaction{
		SubProcess: DefaultSubProcess(),
	}
}

type TransactionInterface interface {
	Element
	SubProcessInterface
	Method() (result *TransactionMethod)
	SetMethod(value *TransactionMethod)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *Transaction) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *Transaction) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *Transaction) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.SubProcess.FindBy(f); found {
		return
	}

	return
}
func (t *Transaction) Method() (result *TransactionMethod) {
	if t.MethodField == nil {
		result = &defaultTransactionMethodField
		return
	}
	result = t.MethodField
	return
}
func (t *Transaction) SetMethod(value *TransactionMethod) {
	t.MethodField = value
}

type UserTask struct {
	Task
	ImplementationField *Implementation `xml:"implementation,attr,omitempty"`
	RenderingField      []Rendering     `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL rendering"`
	TextPayloadField    *Payload        `xml:",chardata"`
}

var defaultUserTaskImplementationField Implementation = "##unspecified"

func DefaultUserTask() UserTask {
	return UserTask{
		Task: DefaultTask(),
	}
}

type UserTaskInterface interface {
	Element
	TaskInterface
	Implementation() (result *Implementation)
	Renderings() (result *[]Rendering)
	SetImplementation(value *Implementation)
	SetRenderings(value []Rendering)

	TextPayload() *string
	SetTextPayload(string)
}

func (t *UserTask) TextPayload() *string {
	s := t.TextPayloadField.String()
	return &s
}

func (t *UserTask) SetTextPayload(text string) {
	payload := Payload(text)
	t.TextPayloadField = &payload
}

func (t *UserTask) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Task.FindBy(f); found {
		return
	}

	for i := range t.RenderingField {
		if result, found = t.RenderingField[i].FindBy(f); found {
			return
		}
	}

	return
}
func (t *UserTask) Implementation() (result *Implementation) {
	if t.ImplementationField == nil {
		result = &defaultUserTaskImplementationField
		return
	}
	result = t.ImplementationField
	return
}
func (t *UserTask) SetImplementation(value *Implementation) {
	t.ImplementationField = value
}
func (t *UserTask) Renderings() (result *[]Rendering) {
	result = &t.RenderingField
	return
}
func (t *UserTask) SetRenderings(value []Rendering) {
	t.RenderingField = value
}

func (t *Definitions) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Definitions(*t)
	return e.EncodeElement(out, start)
}

func (t *Definitions) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Import) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Import(*t)
	return e.EncodeElement(out, start)
}

func (t *Import) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Activity) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Activity(*t)
	return e.EncodeElement(out, start)
}

func (t *Activity) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *AdHocSubProcess) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := AdHocSubProcess(*t)
	return e.EncodeElement(out, start)
}

func (t *AdHocSubProcess) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Artifact) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Artifact(*t)
	return e.EncodeElement(out, start)
}

func (t *Artifact) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Assignment) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Assignment(*t)
	return e.EncodeElement(out, start)
}

func (t *Assignment) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Association) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Association(*t)
	return e.EncodeElement(out, start)
}

func (t *Association) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Auditing) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Auditing(*t)
	return e.EncodeElement(out, start)
}

func (t *Auditing) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *BaseElement) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := BaseElement(*t)
	return e.EncodeElement(out, start)
}

func (t *BaseElement) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *BaseElementWithMixedContent) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := BaseElementWithMixedContent(*t)
	return e.EncodeElement(out, start)
}

func (t *BaseElementWithMixedContent) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *BoundaryEvent) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := BoundaryEvent(*t)
	return e.EncodeElement(out, start)
}

func (t *BoundaryEvent) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *BusinessRuleTask) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := BusinessRuleTask(*t)
	return e.EncodeElement(out, start)
}

func (t *BusinessRuleTask) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *CallableElement) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := CallableElement(*t)
	return e.EncodeElement(out, start)
}

func (t *CallableElement) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *CallActivity) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := CallActivity(*t)
	return e.EncodeElement(out, start)
}

func (t *CallActivity) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *CallChoreography) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := CallChoreography(*t)
	return e.EncodeElement(out, start)
}

func (t *CallChoreography) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *CallConversation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := CallConversation(*t)
	return e.EncodeElement(out, start)
}

func (t *CallConversation) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *CancelEventDefinition) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := CancelEventDefinition(*t)
	return e.EncodeElement(out, start)
}

func (t *CancelEventDefinition) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *CatchEvent) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := CatchEvent(*t)
	return e.EncodeElement(out, start)
}

func (t *CatchEvent) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Category) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Category(*t)
	return e.EncodeElement(out, start)
}

func (t *Category) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *CategoryValue) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := CategoryValue(*t)
	return e.EncodeElement(out, start)
}

func (t *CategoryValue) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Choreography) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Choreography(*t)
	return e.EncodeElement(out, start)
}

func (t *Choreography) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ChoreographyActivity) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ChoreographyActivity(*t)
	return e.EncodeElement(out, start)
}

func (t *ChoreographyActivity) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ChoreographyTask) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ChoreographyTask(*t)
	return e.EncodeElement(out, start)
}

func (t *ChoreographyTask) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Collaboration) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Collaboration(*t)
	return e.EncodeElement(out, start)
}

func (t *Collaboration) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *CompensateEventDefinition) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := CompensateEventDefinition(*t)
	return e.EncodeElement(out, start)
}

func (t *CompensateEventDefinition) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ComplexBehaviorDefinition) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ComplexBehaviorDefinition(*t)
	return e.EncodeElement(out, start)
}

func (t *ComplexBehaviorDefinition) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ComplexGateway) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ComplexGateway(*t)
	return e.EncodeElement(out, start)
}

func (t *ComplexGateway) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ConditionalEventDefinition) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ConditionalEventDefinition(*t)
	return e.EncodeElement(out, start)
}

func (t *ConditionalEventDefinition) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Conversation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Conversation(*t)
	return e.EncodeElement(out, start)
}

func (t *Conversation) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ConversationAssociation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ConversationAssociation(*t)
	return e.EncodeElement(out, start)
}

func (t *ConversationAssociation) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ConversationLink) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ConversationLink(*t)
	return e.EncodeElement(out, start)
}

func (t *ConversationLink) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ConversationNode) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ConversationNode(*t)
	return e.EncodeElement(out, start)
}

func (t *ConversationNode) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *CorrelationKey) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := CorrelationKey(*t)
	return e.EncodeElement(out, start)
}

func (t *CorrelationKey) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *CorrelationProperty) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := CorrelationProperty(*t)
	return e.EncodeElement(out, start)
}

func (t *CorrelationProperty) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *CorrelationPropertyBinding) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := CorrelationPropertyBinding(*t)
	return e.EncodeElement(out, start)
}

func (t *CorrelationPropertyBinding) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *CorrelationPropertyRetrievalExpression) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := CorrelationPropertyRetrievalExpression(*t)
	return e.EncodeElement(out, start)
}

func (t *CorrelationPropertyRetrievalExpression) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *CorrelationSubscription) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := CorrelationSubscription(*t)
	return e.EncodeElement(out, start)
}

func (t *CorrelationSubscription) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *DataAssociation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := DataAssociation(*t)
	return e.EncodeElement(out, start)
}

func (t *DataAssociation) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *DataInput) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := DataInput(*t)
	return e.EncodeElement(out, start)
}

func (t *DataInput) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *DataInputAssociation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := DataInputAssociation(*t)
	return e.EncodeElement(out, start)
}

func (t *DataInputAssociation) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *DataObject) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := DataObject(*t)
	return e.EncodeElement(out, start)
}

func (t *DataObject) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *DataObjectReference) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := DataObjectReference(*t)
	return e.EncodeElement(out, start)
}

func (t *DataObjectReference) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *DataOutput) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := DataOutput(*t)
	return e.EncodeElement(out, start)
}

func (t *DataOutput) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *DataOutputAssociation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := DataOutputAssociation(*t)
	return e.EncodeElement(out, start)
}

func (t *DataOutputAssociation) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *DataState) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := DataState(*t)
	return e.EncodeElement(out, start)
}

func (t *DataState) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *DataStore) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := DataStore(*t)
	return e.EncodeElement(out, start)
}

func (t *DataStore) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *DataStoreReference) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := DataStoreReference(*t)
	return e.EncodeElement(out, start)
}

func (t *DataStoreReference) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Documentation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Documentation(*t)
	return e.EncodeElement(out, start)
}

func (t *Documentation) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *EndEvent) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := EndEvent(*t)
	return e.EncodeElement(out, start)
}

func (t *EndEvent) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *EndPoint) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := EndPoint(*t)
	return e.EncodeElement(out, start)
}

func (t *EndPoint) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Error) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Error(*t)
	return e.EncodeElement(out, start)
}

func (t *Error) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ErrorEventDefinition) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ErrorEventDefinition(*t)
	return e.EncodeElement(out, start)
}

func (t *ErrorEventDefinition) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Escalation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Escalation(*t)
	return e.EncodeElement(out, start)
}

func (t *Escalation) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *EscalationEventDefinition) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := EscalationEventDefinition(*t)
	return e.EncodeElement(out, start)
}

func (t *EscalationEventDefinition) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Event) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Event(*t)
	return e.EncodeElement(out, start)
}

func (t *Event) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *EventBasedGateway) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := EventBasedGateway(*t)
	return e.EncodeElement(out, start)
}

func (t *EventBasedGateway) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *EventDefinition) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := EventDefinition(*t)
	return e.EncodeElement(out, start)
}

func (t *EventDefinition) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ExclusiveGateway) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ExclusiveGateway(*t)
	return e.EncodeElement(out, start)
}

func (t *ExclusiveGateway) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Expression) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Expression(*t)
	return e.EncodeElement(out, start)
}

func (t *Expression) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Extension) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Extension(*t)
	return e.EncodeElement(out, start)
}

func (t *Extension) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ExtensionElements) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ExtensionElements(*t)
	return e.EncodeElement(out, start)
}

func (t *ExtensionElements) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *FlowElement) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := FlowElement(*t)
	return e.EncodeElement(out, start)
}

func (t *FlowElement) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *FlowNode) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := FlowNode(*t)
	return e.EncodeElement(out, start)
}

func (t *FlowNode) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *FormalExpression) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := FormalExpression(*t)
	return e.EncodeElement(out, start)
}

func (t *FormalExpression) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Gateway) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Gateway(*t)
	return e.EncodeElement(out, start)
}

func (t *Gateway) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *GlobalBusinessRuleTask) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := GlobalBusinessRuleTask(*t)
	return e.EncodeElement(out, start)
}

func (t *GlobalBusinessRuleTask) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *GlobalChoreographyTask) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := GlobalChoreographyTask(*t)
	return e.EncodeElement(out, start)
}

func (t *GlobalChoreographyTask) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *GlobalConversation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := GlobalConversation(*t)
	return e.EncodeElement(out, start)
}

func (t *GlobalConversation) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *GlobalManualTask) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := GlobalManualTask(*t)
	return e.EncodeElement(out, start)
}

func (t *GlobalManualTask) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *GlobalScriptTask) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := GlobalScriptTask(*t)
	return e.EncodeElement(out, start)
}

func (t *GlobalScriptTask) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *GlobalTask) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := GlobalTask(*t)
	return e.EncodeElement(out, start)
}

func (t *GlobalTask) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *GlobalUserTask) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := GlobalUserTask(*t)
	return e.EncodeElement(out, start)
}

func (t *GlobalUserTask) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Group) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Group(*t)
	return e.EncodeElement(out, start)
}

func (t *Group) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *HumanPerformer) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := HumanPerformer(*t)
	return e.EncodeElement(out, start)
}

func (t *HumanPerformer) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ImplicitThrowEvent) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ImplicitThrowEvent(*t)
	return e.EncodeElement(out, start)
}

func (t *ImplicitThrowEvent) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *InclusiveGateway) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := InclusiveGateway(*t)
	return e.EncodeElement(out, start)
}

func (t *InclusiveGateway) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *InputSet) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := InputSet(*t)
	return e.EncodeElement(out, start)
}

func (t *InputSet) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Interface) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Interface(*t)
	return e.EncodeElement(out, start)
}

func (t *Interface) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *IntermediateCatchEvent) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := IntermediateCatchEvent(*t)
	return e.EncodeElement(out, start)
}

func (t *IntermediateCatchEvent) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *IntermediateThrowEvent) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := IntermediateThrowEvent(*t)
	return e.EncodeElement(out, start)
}

func (t *IntermediateThrowEvent) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *InputOutputBinding) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := InputOutputBinding(*t)
	return e.EncodeElement(out, start)
}

func (t *InputOutputBinding) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *InputOutputSpecification) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := InputOutputSpecification(*t)
	return e.EncodeElement(out, start)
}

func (t *InputOutputSpecification) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ItemDefinition) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ItemDefinition(*t)
	return e.EncodeElement(out, start)
}

func (t *ItemDefinition) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Lane) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Lane(*t)
	return e.EncodeElement(out, start)
}

func (t *Lane) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *LaneSet) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := LaneSet(*t)
	return e.EncodeElement(out, start)
}

func (t *LaneSet) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *LinkEventDefinition) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := LinkEventDefinition(*t)
	return e.EncodeElement(out, start)
}

func (t *LinkEventDefinition) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *LoopCharacteristics) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := LoopCharacteristics(*t)
	return e.EncodeElement(out, start)
}

func (t *LoopCharacteristics) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ManualTask) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ManualTask(*t)
	return e.EncodeElement(out, start)
}

func (t *ManualTask) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Message) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Message(*t)
	return e.EncodeElement(out, start)
}

func (t *Message) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *MessageEventDefinition) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := MessageEventDefinition(*t)
	return e.EncodeElement(out, start)
}

func (t *MessageEventDefinition) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *MessageFlow) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := MessageFlow(*t)
	return e.EncodeElement(out, start)
}

func (t *MessageFlow) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *MessageFlowAssociation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := MessageFlowAssociation(*t)
	return e.EncodeElement(out, start)
}

func (t *MessageFlowAssociation) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Monitoring) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Monitoring(*t)
	return e.EncodeElement(out, start)
}

func (t *Monitoring) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *MultiInstanceLoopCharacteristics) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := MultiInstanceLoopCharacteristics(*t)
	return e.EncodeElement(out, start)
}

func (t *MultiInstanceLoopCharacteristics) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Operation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Operation(*t)
	return e.EncodeElement(out, start)
}

func (t *Operation) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *OutputSet) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := OutputSet(*t)
	return e.EncodeElement(out, start)
}

func (t *OutputSet) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ParallelGateway) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ParallelGateway(*t)
	return e.EncodeElement(out, start)
}

func (t *ParallelGateway) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Participant) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Participant(*t)
	return e.EncodeElement(out, start)
}

func (t *Participant) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ParticipantAssociation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ParticipantAssociation(*t)
	return e.EncodeElement(out, start)
}

func (t *ParticipantAssociation) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ParticipantMultiplicity) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ParticipantMultiplicity(*t)
	return e.EncodeElement(out, start)
}

func (t *ParticipantMultiplicity) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *PartnerEntity) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := PartnerEntity(*t)
	return e.EncodeElement(out, start)
}

func (t *PartnerEntity) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *PartnerRole) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := PartnerRole(*t)
	return e.EncodeElement(out, start)
}

func (t *PartnerRole) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Performer) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Performer(*t)
	return e.EncodeElement(out, start)
}

func (t *Performer) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *PotentialOwner) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := PotentialOwner(*t)
	return e.EncodeElement(out, start)
}

func (t *PotentialOwner) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Process) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Process(*t)
	return e.EncodeElement(out, start)
}

func (t *Process) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Property) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Property(*t)
	return e.EncodeElement(out, start)
}

func (t *Property) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ReceiveTask) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ReceiveTask(*t)
	return e.EncodeElement(out, start)
}

func (t *ReceiveTask) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Relationship) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Relationship(*t)
	return e.EncodeElement(out, start)
}

func (t *Relationship) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Rendering) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Rendering(*t)
	return e.EncodeElement(out, start)
}

func (t *Rendering) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Resource) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Resource(*t)
	return e.EncodeElement(out, start)
}

func (t *Resource) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ResourceAssignmentExpression) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ResourceAssignmentExpression(*t)
	return e.EncodeElement(out, start)
}

func (t *ResourceAssignmentExpression) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ResourceParameter) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ResourceParameter(*t)
	return e.EncodeElement(out, start)
}

func (t *ResourceParameter) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ResourceParameterBinding) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ResourceParameterBinding(*t)
	return e.EncodeElement(out, start)
}

func (t *ResourceParameterBinding) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ResourceRole) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ResourceRole(*t)
	return e.EncodeElement(out, start)
}

func (t *ResourceRole) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *RootElement) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := RootElement(*t)
	return e.EncodeElement(out, start)
}

func (t *RootElement) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ScriptTask) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ScriptTask(*t)
	return e.EncodeElement(out, start)
}

func (t *ScriptTask) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Script) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Script(*t)
	return e.EncodeElement(out, start)
}

func (t *Script) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *SendTask) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := SendTask(*t)
	return e.EncodeElement(out, start)
}

func (t *SendTask) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *SequenceFlow) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := SequenceFlow(*t)
	return e.EncodeElement(out, start)
}

func (t *SequenceFlow) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ServiceTask) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ServiceTask(*t)
	return e.EncodeElement(out, start)
}

func (t *ServiceTask) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Signal) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Signal(*t)
	return e.EncodeElement(out, start)
}

func (t *Signal) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *SignalEventDefinition) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := SignalEventDefinition(*t)
	return e.EncodeElement(out, start)
}

func (t *SignalEventDefinition) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *StandardLoopCharacteristics) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := StandardLoopCharacteristics(*t)
	return e.EncodeElement(out, start)
}

func (t *StandardLoopCharacteristics) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *StartEvent) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := StartEvent(*t)
	return e.EncodeElement(out, start)
}

func (t *StartEvent) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *SubChoreography) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := SubChoreography(*t)
	return e.EncodeElement(out, start)
}

func (t *SubChoreography) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *SubConversation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := SubConversation(*t)
	return e.EncodeElement(out, start)
}

func (t *SubConversation) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *SubProcess) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := SubProcess(*t)
	return e.EncodeElement(out, start)
}

func (t *SubProcess) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Task) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Task(*t)
	return e.EncodeElement(out, start)
}

func (t *Task) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *TerminateEventDefinition) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := TerminateEventDefinition(*t)
	return e.EncodeElement(out, start)
}

func (t *TerminateEventDefinition) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *TextAnnotation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := TextAnnotation(*t)
	return e.EncodeElement(out, start)
}

func (t *TextAnnotation) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Text) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Text(*t)
	return e.EncodeElement(out, start)
}

func (t *Text) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *ThrowEvent) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := ThrowEvent(*t)
	return e.EncodeElement(out, start)
}

func (t *ThrowEvent) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *TimerEventDefinition) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := TimerEventDefinition(*t)
	return e.EncodeElement(out, start)
}

func (t *TimerEventDefinition) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *Transaction) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Transaction(*t)
	return e.EncodeElement(out, start)
}

func (t *Transaction) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}

func (t *UserTask) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := UserTask(*t)
	return e.EncodeElement(out, start)
}

func (t *UserTask) UnMarshalXML(de *xml.Decoder, start *xml.StartElement) error {
	PreUnmarshal(t, de, start)
	return de.DecodeElement(t, start)
}
