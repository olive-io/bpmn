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
)

type ParticipantBandKind string
type MessageVisibleKind string
type BPMNDiagram struct {
	Diagram
	BPMNPlaneField      *BPMNPlane       `xml:"http://www.omg.org/spec/BPMN/20100524/DI BPMNPlane"`
	BPMNLabelStyleField []BPMNLabelStyle `xml:"http://www.omg.org/spec/BPMN/20100524/DI BPMNLabelStyle"`
}

func DefaultBPMNDiagram() BPMNDiagram {
	return BPMNDiagram{
		Diagram: DefaultDiagram(),
	}
}

type BPMNDiagramInterface interface {
	Element
	DiagramInterface
	BPMNPlane() (result *BPMNPlane)
	BPMNLabelStyle() (result *[]BPMNLabelStyle)
	SetBPMNPlane(value *BPMNPlane)
	SetBPMNLabelStyle(value []BPMNLabelStyle)
}

func (t *BPMNDiagram) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Diagram.FindBy(f); found {
		return
	}

	return
}
func (t *BPMNDiagram) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := BPMNDiagram(*t)
	return e.EncodeElement(out, start)
}

func (t *BPMNDiagram) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type BPMNDiagramUnmarshaler BPMNDiagram
	out := BPMNDiagramUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = BPMNDiagram(out)
	PostUnmarshal(t, de, &start)
	return nil
}

func (t *BPMNDiagram) BPMNPlane() (result *BPMNPlane) {
	result = t.BPMNPlaneField

	return
}
func (t *BPMNDiagram) SetBPMNPlane(value *BPMNPlane) {
	t.BPMNPlaneField = value
}
func (t *BPMNDiagram) BPMNLabelStyle() (result *[]BPMNLabelStyle) {
	result = &t.BPMNLabelStyleField

	return
}
func (t *BPMNDiagram) SetBPMNLabelStyle(value []BPMNLabelStyle) {
	t.BPMNLabelStyleField = value
}

type BPMNPlane struct {
	Plane
	BpmnElementField *QName `xml:"bpmnElement,attr,omitempty"`
}

func DefaultBPMNPlane() BPMNPlane {
	return BPMNPlane{
		Plane: DefaultPlane(),
	}
}

type BPMNPlaneInterface interface {
	Element
	PlaneInterface
	BpmnElement() (result *QName, present bool)
	SetBpmnElement(value *QName)
}

func (t *BPMNPlane) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Plane.FindBy(f); found {
		return
	}

	return
}
func (t *BPMNPlane) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := BPMNPlane(*t)
	return e.EncodeElement(out, start)
}

func (t *BPMNPlane) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type BPMNPlaneUnmarshaler BPMNPlane
	out := BPMNPlaneUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = BPMNPlane(out)
	PostUnmarshal(t, de, &start)
	return nil
}

func (t *BPMNPlane) BpmnElement() (result *QName, present bool) {
	if t.BpmnElementField != nil {
		present = true
	}
	result = t.BpmnElementField
	return
}
func (t *BPMNPlane) SetBpmnElement(value *QName) {
	t.BpmnElementField = value
}

type BPMNEdge struct {
	LabeledEdge
	BpmnElementField        *QName              `xml:"bpmnElement,attr,omitempty"`
	SourceElementField      *QName              `xml:"sourceElement,attr,omitempty"`
	TargetElementField      *QName              `xml:"targetElement,attr,omitempty"`
	MessageVisibleKindField *MessageVisibleKind `xml:"messageVisibleKind,attr,omitempty"`
	BPMNLabelField          *BPMNLabel          `xml:"http://www.omg.org/spec/BPMN/20100524/DI BPMNLabel"`
}

func DefaultBPMNEdge() BPMNEdge {
	return BPMNEdge{
		LabeledEdge: DefaultLabeledEdge(),
	}
}

type BPMNEdgeInterface interface {
	Element
	LabeledEdgeInterface
	BpmnElement() (result *QName, present bool)
	SourceElement() (result *QName, present bool)
	TargetElement() (result *QName, present bool)
	MessageVisibleKind() (result *MessageVisibleKind, present bool)
	BPMNLabel() (result *BPMNLabel)
	SetBpmnElement(value *QName)
	SetSourceElement(value *QName)
	SetTargetElement(value *QName)
	SetMessageVisibleKind(value *MessageVisibleKind)
	SetBPMNLabel(value *BPMNLabel)
}

func (t *BPMNEdge) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.LabeledEdge.FindBy(f); found {
		return
	}

	return
}
func (t *BPMNEdge) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := BPMNEdge(*t)
	return e.EncodeElement(out, start)
}

func (t *BPMNEdge) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type BPMNEdgeUnmarshaler BPMNEdge
	out := BPMNEdgeUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = BPMNEdge(out)
	PostUnmarshal(t, de, &start)
	return nil
}

func (t *BPMNEdge) BpmnElement() (result *QName, present bool) {
	if t.BpmnElementField != nil {
		present = true
	}
	result = t.BpmnElementField
	return
}
func (t *BPMNEdge) SetBpmnElement(value *QName) {
	t.BpmnElementField = value
}
func (t *BPMNEdge) SourceElement() (result *QName, present bool) {
	if t.SourceElementField != nil {
		present = true
	}
	result = t.SourceElementField
	return
}
func (t *BPMNEdge) SetSourceElement(value *QName) {
	t.SourceElementField = value
}
func (t *BPMNEdge) TargetElement() (result *QName, present bool) {
	if t.TargetElementField != nil {
		present = true
	}
	result = t.TargetElementField
	return
}
func (t *BPMNEdge) SetTargetElement(value *QName) {
	t.TargetElementField = value
}
func (t *BPMNEdge) MessageVisibleKind() (result *MessageVisibleKind, present bool) {
	if t.MessageVisibleKindField != nil {
		present = true
	}
	result = t.MessageVisibleKindField
	return
}
func (t *BPMNEdge) SetMessageVisibleKind(value *MessageVisibleKind) {
	t.MessageVisibleKindField = value
}
func (t *BPMNEdge) BPMNLabel() (result *BPMNLabel) {
	result = t.BPMNLabelField

	return
}
func (t *BPMNEdge) SetBPMNLabel(value *BPMNLabel) {
	t.BPMNLabelField = value
}

type BPMNShape struct {
	LabeledShape
	BpmnElementField               *QName               `xml:"bpmnElement,attr,omitempty"`
	IsHorizontalField              *bool                `xml:"isHorizontal,attr,omitempty"`
	IsExpandedField                *bool                `xml:"isExpanded,attr,omitempty"`
	IsMarkerVisibleField           *bool                `xml:"isMarkerVisible,attr,omitempty"`
	IsMessageVisibleField          *bool                `xml:"isMessageVisible,attr,omitempty"`
	ParticipantBandKindField       *ParticipantBandKind `xml:"participantBandKind,attr,omitempty"`
	ChoreographyActivityShapeField *QName               `xml:"choreographyActivityShape,attr,omitempty"`
	BPMNLabelField                 *BPMNLabel           `xml:"http://www.omg.org/spec/BPMN/20100524/DI BPMNLabel"`
}

func DefaultBPMNShape() BPMNShape {
	return BPMNShape{
		LabeledShape: DefaultLabeledShape(),
	}
}

type BPMNShapeInterface interface {
	Element
	LabeledShapeInterface
	BpmnElement() (result *QName, present bool)
	IsHorizontal() (result bool, present bool)
	IsExpanded() (result bool, present bool)
	IsMarkerVisible() (result bool, present bool)
	IsMessageVisible() (result bool, present bool)
	ParticipantBandKind() (result *ParticipantBandKind, present bool)
	ChoreographyActivityShape() (result *QName, present bool)
	BPMNLabel() (result *BPMNLabel)
	SetBpmnElement(value *QName)
	SetIsHorizontal(value *bool)
	SetIsExpanded(value *bool)
	SetIsMarkerVisible(value *bool)
	SetIsMessageVisible(value *bool)
	SetParticipantBandKind(value *ParticipantBandKind)
	SetChoreographyActivityShape(value *QName)
	SetBPMNLabel(value *BPMNLabel)
}

func (t *BPMNShape) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.LabeledShape.FindBy(f); found {
		return
	}

	return
}
func (t *BPMNShape) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := BPMNShape(*t)
	return e.EncodeElement(out, start)
}

func (t *BPMNShape) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type BPMNShapeUnmarshaler BPMNShape
	out := BPMNShapeUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = BPMNShape(out)
	PostUnmarshal(t, de, &start)
	return nil
}

func (t *BPMNShape) BpmnElement() (result *QName, present bool) {
	if t.BpmnElementField != nil {
		present = true
	}
	result = t.BpmnElementField
	return
}
func (t *BPMNShape) SetBpmnElement(value *QName) {
	t.BpmnElementField = value
}
func (t *BPMNShape) IsHorizontal() (result bool, present bool) {
	if t.IsHorizontalField != nil {
		present = true
	}
	result = *t.IsHorizontalField
	return
}
func (t *BPMNShape) SetIsHorizontal(value *bool) {
	t.IsHorizontalField = value
}
func (t *BPMNShape) IsExpanded() (result bool, present bool) {
	if t.IsExpandedField != nil {
		present = true
	}
	result = *t.IsExpandedField
	return
}
func (t *BPMNShape) SetIsExpanded(value *bool) {
	t.IsExpandedField = value
}
func (t *BPMNShape) IsMarkerVisible() (result bool, present bool) {
	if t.IsMarkerVisibleField != nil {
		present = true
	}
	result = *t.IsMarkerVisibleField
	return
}
func (t *BPMNShape) SetIsMarkerVisible(value *bool) {
	t.IsMarkerVisibleField = value
}
func (t *BPMNShape) IsMessageVisible() (result bool, present bool) {
	if t.IsMessageVisibleField != nil {
		present = true
	}
	result = *t.IsMessageVisibleField
	return
}
func (t *BPMNShape) SetIsMessageVisible(value *bool) {
	t.IsMessageVisibleField = value
}
func (t *BPMNShape) ParticipantBandKind() (result *ParticipantBandKind, present bool) {
	if t.ParticipantBandKindField != nil {
		present = true
	}
	result = t.ParticipantBandKindField
	return
}
func (t *BPMNShape) SetParticipantBandKind(value *ParticipantBandKind) {
	t.ParticipantBandKindField = value
}
func (t *BPMNShape) ChoreographyActivityShape() (result *QName, present bool) {
	if t.ChoreographyActivityShapeField != nil {
		present = true
	}
	result = t.ChoreographyActivityShapeField
	return
}
func (t *BPMNShape) SetChoreographyActivityShape(value *QName) {
	t.ChoreographyActivityShapeField = value
}
func (t *BPMNShape) BPMNLabel() (result *BPMNLabel) {
	result = t.BPMNLabelField

	return
}
func (t *BPMNShape) SetBPMNLabel(value *BPMNLabel) {
	t.BPMNLabelField = value
}

type BPMNLabel struct {
	Label
	LabelStyleField *QName `xml:"labelStyle,attr,omitempty"`
}

func DefaultBPMNLabel() BPMNLabel {
	return BPMNLabel{
		Label: DefaultLabel(),
	}
}

type BPMNLabelInterface interface {
	Element
	LabelInterface
	LabelStyle() (result *QName, present bool)
	SetLabelStyle(value *QName)
}

func (t *BPMNLabel) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Label.FindBy(f); found {
		return
	}

	return
}
func (t *BPMNLabel) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := BPMNLabel(*t)
	return e.EncodeElement(out, start)
}

func (t *BPMNLabel) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type BPMNLabelUnmarshaler BPMNLabel
	out := BPMNLabelUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = BPMNLabel(out)
	PostUnmarshal(t, de, &start)
	return nil
}

func (t *BPMNLabel) LabelStyle() (result *QName, present bool) {
	if t.LabelStyleField != nil {
		present = true
	}
	result = t.LabelStyleField
	return
}
func (t *BPMNLabel) SetLabelStyle(value *QName) {
	t.LabelStyleField = value
}

type BPMNLabelStyle struct {
	Style
	FontField *Font `xml:"http://www.omg.org/spec/DD/20100524/DC Font"`
}

func DefaultBPMNLabelStyle() BPMNLabelStyle {
	return BPMNLabelStyle{
		Style: DefaultStyle(),
	}
}

type BPMNLabelStyleInterface interface {
	Element
	StyleInterface
	Font() (result *Font)
	SetFont(value *Font)
}

func (t *BPMNLabelStyle) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Style.FindBy(f); found {
		return
	}

	return
}
func (t *BPMNLabelStyle) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := BPMNLabelStyle(*t)
	return e.EncodeElement(out, start)
}

func (t *BPMNLabelStyle) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type BPMNLabelStyleUnmarshaler BPMNLabelStyle
	out := BPMNLabelStyleUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = BPMNLabelStyle(out)
	PostUnmarshal(t, de, &start)
	return nil
}

func (t *BPMNLabelStyle) Font() (result *Font) {
	result = t.FontField

	return
}
func (t *BPMNLabelStyle) SetFont(value *Font) {
	t.FontField = value
}

type Font struct {
	NameField            *string `xml:"name,attr,omitempty"`
	SizeField            *Double `xml:"size,attr,omitempty"`
	IsBoldField          *bool   `xml:"isBold,attr,omitempty"`
	IsItalicField        *bool   `xml:"isItalic,attr,omitempty"`
	IsUnderlineField     *bool   `xml:"isUnderline,attr,omitempty"`
	IsStrikeThroughField *bool   `xml:"isStrikeThrough,attr,omitempty"`
}

func DefaultFont() Font {
	return Font{}
}

type FontInterface interface {
	Element
	Name() (result *string, present bool)
	Size() (result Double, present bool)
	IsBold() (result bool, present bool)
	IsItalic() (result bool, present bool)
	IsUnderline() (result bool, present bool)
	IsStrikeThrough() (result bool, present bool)
	SetName(value *string)
	SetSize(value *Double)
	SetIsBold(value *bool)
	SetIsItalic(value *bool)
	SetIsUnderline(value *bool)
	SetIsStrikeThrough(value *bool)
}

func (t *Font) FindBy(f ElementPredicate) (result Element, found bool) {
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
func (t *Font) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Font(*t)
	return e.EncodeElement(out, start)
}

func (t *Font) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type FontUnmarshaler Font
	out := FontUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = Font(out)
	PostUnmarshal(t, de, &start)
	return nil
}

func (t *Font) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *Font) SetName(value *string) {
	t.NameField = value
}
func (t *Font) Size() (result Double, present bool) {
	if t.SizeField != nil {
		present = true
	}
	result = *t.SizeField
	return
}
func (t *Font) SetSize(value *Double) {
	t.SizeField = value
}
func (t *Font) IsBold() (result bool, present bool) {
	if t.IsBoldField != nil {
		present = true
	}
	result = *t.IsBoldField
	return
}
func (t *Font) SetIsBold(value *bool) {
	t.IsBoldField = value
}
func (t *Font) IsItalic() (result bool, present bool) {
	if t.IsItalicField != nil {
		present = true
	}
	result = *t.IsItalicField
	return
}
func (t *Font) SetIsItalic(value *bool) {
	t.IsItalicField = value
}
func (t *Font) IsUnderline() (result bool, present bool) {
	if t.IsUnderlineField != nil {
		present = true
	}
	result = *t.IsUnderlineField
	return
}
func (t *Font) SetIsUnderline(value *bool) {
	t.IsUnderlineField = value
}
func (t *Font) IsStrikeThrough() (result bool, present bool) {
	if t.IsStrikeThroughField != nil {
		present = true
	}
	result = *t.IsStrikeThroughField
	return
}
func (t *Font) SetIsStrikeThrough(value *bool) {
	t.IsStrikeThroughField = value
}

type Point struct {
	XField Double `xml:"x,attr,omitempty"`
	YField Double `xml:"y,attr,omitempty"`
}

func DefaultPoint() Point {
	return Point{}
}

type PointInterface interface {
	Element
	X() (result Double)
	Y() (result Double)
	SetX(value Double)
	SetY(value Double)
}

func (t *Point) FindBy(f ElementPredicate) (result Element, found bool) {
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
func (t *Point) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Point(*t)
	return e.EncodeElement(out, start)
}

func (t *Point) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type PointUnmarshaler Point
	out := PointUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = Point(out)
	PostUnmarshal(t, de, &start)
	return nil
}

func (t *Point) X() (result Double) {
	result = *&t.XField
	return
}
func (t *Point) SetX(value Double) {
	t.XField = value
}
func (t *Point) Y() (result Double) {
	result = *&t.YField
	return
}
func (t *Point) SetY(value Double) {
	t.YField = value
}

type Bounds struct {
	XField      Double `xml:"x,attr,omitempty"`
	YField      Double `xml:"y,attr,omitempty"`
	WidthField  Double `xml:"width,attr,omitempty"`
	HeightField Double `xml:"height,attr,omitempty"`
}

func DefaultBounds() Bounds {
	return Bounds{}
}

type BoundsInterface interface {
	Element
	X() (result Double)
	Y() (result Double)
	Width() (result Double)
	Height() (result Double)
	SetX(value Double)
	SetY(value Double)
	SetWidth(value Double)
	SetHeight(value Double)
}

func (t *Bounds) FindBy(f ElementPredicate) (result Element, found bool) {
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
func (t *Bounds) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Bounds(*t)
	return e.EncodeElement(out, start)
}

func (t *Bounds) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type BoundsUnmarshaler Bounds
	out := BoundsUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = Bounds(out)
	PostUnmarshal(t, de, &start)
	return nil
}

func (t *Bounds) X() (result Double) {
	result = *&t.XField
	return
}
func (t *Bounds) SetX(value Double) {
	t.XField = value
}
func (t *Bounds) Y() (result Double) {
	result = *&t.YField
	return
}
func (t *Bounds) SetY(value Double) {
	t.YField = value
}
func (t *Bounds) Width() (result Double) {
	result = *&t.WidthField
	return
}
func (t *Bounds) SetWidth(value Double) {
	t.WidthField = value
}
func (t *Bounds) Height() (result Double) {
	result = *&t.HeightField
	return
}
func (t *Bounds) SetHeight(value Double) {
	t.HeightField = value
}

type DiagramElement struct {
	IdField        *Id          `xml:"id,attr,omitempty"`
	ExtensionField *DIExtension `xml:"http://www.omg.org/spec/DD/20100524/DI extension"`
}

func DefaultDiagramElement() DiagramElement {
	return DiagramElement{}
}

type DiagramElementInterface interface {
	Element
	Id() (result *Id, present bool)
	Extension() (result *DIExtension)
	SetId(value *Id)
	SetExtension(value *DIExtension)
}

func (t *DiagramElement) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}

	if result, found = t.ExtensionField.FindBy(f); found {
		return
	}

	return
}
func (t *DiagramElement) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := DiagramElement(*t)
	return e.EncodeElement(out, start)
}

func (t *DiagramElement) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type DiagramElementUnmarshaler DiagramElement
	out := DiagramElementUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = DiagramElement(out)
	PostUnmarshal(t, de, &start)
	return nil
}

func (t *DiagramElement) Id() (result *Id, present bool) {
	if t.IdField != nil {
		present = true
	}
	result = t.IdField
	return
}
func (t *DiagramElement) SetId(value *Id) {
	t.IdField = value
}
func (t *DiagramElement) Extension() (result *DIExtension) {
	result = t.ExtensionField
	return
}
func (t *DiagramElement) SetExtension(value *DIExtension) {
	t.ExtensionField = value
}

type Diagram struct {
	NameField          *string `xml:"name,attr,omitempty"`
	DocumentationField *string `xml:"documentation,attr,omitempty"`
	ResolutionField    *Double `xml:"resolution,attr,omitempty"`
	IdField            *Id     `xml:"id,attr,omitempty"`
}

func DefaultDiagram() Diagram {
	return Diagram{}
}

type DiagramInterface interface {
	Element
	Name() (result *string, present bool)
	Documentation() (result *string, present bool)
	Resolution() (result Double, present bool)
	Id() (result *Id, present bool)
	SetName(value *string)
	SetDocumentation(value *string)
	SetResolution(value *Double)
	SetId(value *Id)
}

func (t *Diagram) FindBy(f ElementPredicate) (result Element, found bool) {
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
func (t *Diagram) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Diagram(*t)
	return e.EncodeElement(out, start)
}

func (t *Diagram) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type DiagramUnmarshaler Diagram
	out := DiagramUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = Diagram(out)
	PostUnmarshal(t, de, &start)
	return nil
}

func (t *Diagram) Name() (result *string, present bool) {
	if t.NameField != nil {
		present = true
	}
	result = t.NameField
	return
}
func (t *Diagram) SetName(value *string) {
	t.NameField = value
}
func (t *Diagram) Documentation() (result *string, present bool) {
	if t.DocumentationField != nil {
		present = true
	}
	result = t.DocumentationField
	return
}
func (t *Diagram) SetDocumentation(value *string) {
	t.DocumentationField = value
}
func (t *Diagram) Resolution() (result Double, present bool) {
	if t.ResolutionField != nil {
		present = true
	}
	result = *t.ResolutionField
	return
}
func (t *Diagram) SetResolution(value *Double) {
	t.ResolutionField = value
}
func (t *Diagram) Id() (result *Id, present bool) {
	if t.IdField != nil {
		present = true
	}
	result = t.IdField
	return
}
func (t *Diagram) SetId(value *Id) {
	t.IdField = value
}

type Node struct {
	DiagramElement
}

func DefaultNode() Node {
	return Node{
		DiagramElement: DefaultDiagramElement(),
	}
}

type NodeInterface interface {
	Element
	DiagramElementInterface
}

func (t *Node) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.DiagramElement.FindBy(f); found {
		return
	}

	return
}
func (t *Node) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Node(*t)
	return e.EncodeElement(out, start)
}

func (t *Node) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type NodeUnmarshaler Node
	out := NodeUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = Node(out)
	PostUnmarshal(t, de, &start)
	return nil
}

type Edge struct {
	DiagramElement
	WaypointField []Point `xml:"http://www.omg.org/spec/DD/20100524/DI waypoint"`
}

func DefaultEdge() Edge {
	return Edge{
		DiagramElement: DefaultDiagramElement(),
	}
}

type EdgeInterface interface {
	Element
	DiagramElementInterface
	Waypoints() (result *[]Point)
	SetWaypoints(value []Point)
}

func (t *Edge) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.DiagramElement.FindBy(f); found {
		return
	}

	return
}
func (t *Edge) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Edge(*t)
	return e.EncodeElement(out, start)
}

func (t *Edge) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type EdgeUnmarshaler Edge
	out := EdgeUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = Edge(out)
	PostUnmarshal(t, de, &start)
	return nil
}

func (t *Edge) Waypoints() (result *[]Point) {
	result = &t.WaypointField
	return
}
func (t *Edge) SetWaypoints(value []Point) {
	t.WaypointField = value
}

type LabeledEdge struct {
	Edge
}

func DefaultLabeledEdge() LabeledEdge {
	return LabeledEdge{
		Edge: DefaultEdge(),
	}
}

type LabeledEdgeInterface interface {
	Element
	EdgeInterface
}

func (t *LabeledEdge) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Edge.FindBy(f); found {
		return
	}

	return
}
func (t *LabeledEdge) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := LabeledEdge(*t)
	return e.EncodeElement(out, start)
}

func (t *LabeledEdge) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type LabeledEdgeUnmarshaler LabeledEdge
	out := LabeledEdgeUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = LabeledEdge(out)
	PostUnmarshal(t, de, &start)
	return nil
}

type Shape struct {
	Node
	BoundsField *Bounds `xml:"http://www.omg.org/spec/DD/20100524/DC Bounds"`
}

func DefaultShape() Shape {
	return Shape{
		Node: DefaultNode(),
	}
}

type ShapeInterface interface {
	Element
	NodeInterface
	Bounds() (result *Bounds)
	SetBounds(value *Bounds)
}

func (t *Shape) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Node.FindBy(f); found {
		return
	}

	return
}
func (t *Shape) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Shape(*t)
	return e.EncodeElement(out, start)
}

func (t *Shape) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type ShapeUnmarshaler Shape
	out := ShapeUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = Shape(out)
	PostUnmarshal(t, de, &start)
	return nil
}

func (t *Shape) Bounds() (result *Bounds) {
	result = t.BoundsField

	return
}
func (t *Shape) SetBounds(value *Bounds) {
	t.BoundsField = value
}

type LabeledShape struct {
	Shape
}

func DefaultLabeledShape() LabeledShape {
	return LabeledShape{
		Shape: DefaultShape(),
	}
}

type LabeledShapeInterface interface {
	Element
	ShapeInterface
}

func (t *LabeledShape) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Shape.FindBy(f); found {
		return
	}

	return
}
func (t *LabeledShape) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := LabeledShape(*t)
	return e.EncodeElement(out, start)
}

func (t *LabeledShape) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type LabeledShapeUnmarshaler LabeledShape
	out := LabeledShapeUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = LabeledShape(out)
	PostUnmarshal(t, de, &start)
	return nil
}

type Label struct {
	Node
	BoundsField *Bounds `xml:"http://www.omg.org/spec/DD/20100524/DC Bounds"`
}

func DefaultLabel() Label {
	return Label{
		Node: DefaultNode(),
	}
}

type LabelInterface interface {
	Element
	NodeInterface
	Bounds() (result *Bounds)
	SetBounds(value *Bounds)
}

func (t *Label) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Node.FindBy(f); found {
		return
	}

	return
}
func (t *Label) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Label(*t)
	return e.EncodeElement(out, start)
}

func (t *Label) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type LabelUnmarshaler Label
	out := LabelUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = Label(out)
	PostUnmarshal(t, de, &start)
	return nil
}

func (t *Label) Bounds() (result *Bounds) {
	result = t.BoundsField

	return
}
func (t *Label) SetBounds(value *Bounds) {
	t.BoundsField = value
}

type Plane struct {
	Node
	DiagramElementField []DiagramElement `xml:"http://www.omg.org/spec/DD/20100524/DI DiagramElement"`
	BPMNShapeFields     []BPMNShape      `xml:"http://www.omg.org/spec/BPMN/20100524/DI BPMNShape"`
	BPMNEdgeFields      []BPMNEdge       `xml:"http://www.omg.org/spec/BPMN/20100524/DI BPMNEdge"`
}

func DefaultPlane() Plane {
	return Plane{
		Node: DefaultNode(),
	}
}

type PlaneInterface interface {
	Element
	NodeInterface
	DiagramElement() (result *[]DiagramElement)
	SetDiagramElement(value []DiagramElement)
}

func (t *Plane) FindBy(f ElementPredicate) (result Element, found bool) {
	if t == nil {
		return
	}
	if f(t) {
		result = t
		found = true
		return
	}
	if result, found = t.Node.FindBy(f); found {
		return
	}

	return
}
func (t *Plane) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Plane(*t)
	return e.EncodeElement(out, start)
}

func (t *Plane) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type PlaneUnmarshaler Plane
	out := PlaneUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = Plane(out)
	PostUnmarshal(t, de, &start)
	return nil
}

func (t *Plane) DiagramElement() (result *[]DiagramElement) {
	result = &t.DiagramElementField

	return
}
func (t *Plane) SetDiagramElement(value []DiagramElement) {
	t.DiagramElementField = value
}

type Style struct {
	IdField *Id `xml:"id,attr,omitempty"`
}

func DefaultStyle() Style {
	return Style{}
}

type StyleInterface interface {
	Element
	Id() (result *Id, present bool)
	SetId(value *Id)
}

func (t *Style) FindBy(f ElementPredicate) (result Element, found bool) {
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
func (t *Style) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	PreMarshal(t, e, &start)
	out := Style(*t)
	return e.EncodeElement(out, start)
}

func (t *Style) UnMarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type StyleUnmarshaler Style
	out := StyleUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return nil
	}
	*t = Style(out)
	PostUnmarshal(t, de, &start)
	return nil
}

func (t *Style) Id() (result *Id, present bool) {
	if t.IdField != nil {
		present = true
	}
	result = t.IdField
	return
}
func (t *Style) SetId(value *Id) {
	t.IdField = value
}
