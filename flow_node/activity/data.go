/*
   Copyright 2023 The bpmn Authors

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

package activity

import (
	"strings"

	"github.com/olive-io/bpmn/schema"

	"github.com/olive-io/bpmn/data"
)

func FetchTaskDataInput(locator data.IFlowDataLocator, element schema.BaseElementInterface) (headers, properties, dataObjects map[string]any) {
	variables := locator.CloneVariables()
	headers = map[string]any{}
	properties = map[string]any{}
	dataObjects = map[string]any{}
	if extension, found := element.ExtensionElements(); found {
		if header := extension.TaskHeaderField; header != nil {
			fields := header.Header
			for _, field := range fields {
				if field.Type == "" {
					field.Type = schema.ItemTypeString
				}
				value := field.ValueFor()
				headers[field.Name] = value
			}
		}
		if property := extension.PropertiesField; property != nil {
			fields := property.Property
			for _, field := range fields {
				value := field.ValueFor()
				if len(strings.TrimSpace(field.Value)) == 0 {
					value = variables[field.Name]
				}
				properties[field.Name] = value
			}
		}

		awareLocator, found1 := locator.FindIItemAwareLocator(data.LocatorObject)
		if found1 {
			for _, dataInput := range extension.DataInput {
				aWare, ok := awareLocator.FindItemAwareById(dataInput.TargetRef)
				if ok {
					dataObjects[dataInput.Name] = aWare.Get()
				}
			}
		}
	}

	return
}

func ApplyTaskDataOutput(element schema.BaseElementInterface, dataOutputs map[string]any) map[string]data.IItem {
	outputs := map[string]data.IItem{}
	if extension, found := element.ExtensionElements(); found {
		for _, dataOutput := range extension.DataOutput {
			value, ok := dataOutputs[dataOutput.Name]
			if ok {
				outputs[dataOutput.Name] = value
			}
		}
	}
	return outputs
}
