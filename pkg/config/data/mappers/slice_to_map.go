package mappers

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"

	types "github.com/bhojpur/os/pkg/config/data"
	"github.com/bhojpur/os/pkg/config/data/definition"
)

type SliceToMap struct {
	Field string
	Key   string
}

func (s SliceToMap) FromInternal(data map[string]interface{}) {
	datas, _ := data[s.Field].([]interface{})
	result := map[string]interface{}{}

	for _, item := range datas {
		if mapItem, ok := item.(map[string]interface{}); ok {
			name, _ := mapItem[s.Key].(string)
			delete(mapItem, s.Key)
			result[name] = mapItem
		}
	}

	if len(result) > 0 {
		data[s.Field] = result
	}
}

func (s SliceToMap) ToInternal(data map[string]interface{}) error {
	datas, _ := data[s.Field].(map[string]interface{})
	var result []interface{}

	for name, item := range datas {
		mapItem, _ := item.(map[string]interface{})
		if mapItem != nil {
			mapItem[s.Key] = name
			result = append(result, mapItem)
		}
	}

	if len(result) > 0 {
		data[s.Field] = result
	} else if datas != nil {
		data[s.Field] = result
	}

	return nil
}

func (s SliceToMap) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	err := ValidateField(s.Field, schema)
	if err != nil {
		return err
	}

	subSchema, subFieldName, _, _, err := getField(schema, schemas, fmt.Sprintf("%s/%s", s.Field, s.Key))
	if err != nil {
		return err
	}

	field := schema.ResourceFields[s.Field]
	if !definition.IsArrayType(field.Type) {
		return fmt.Errorf("field %s on %s is not an array", s.Field, schema.ID)
	}

	field.Type = "map[" + definition.SubType(field.Type) + "]"
	schema.ResourceFields[s.Field] = field

	delete(subSchema.ResourceFields, subFieldName)

	return nil
}
