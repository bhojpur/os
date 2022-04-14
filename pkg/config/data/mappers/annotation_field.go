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
	"encoding/json"

	types "github.com/bhojpur/os/pkg/config/data"
	"github.com/bhojpur/os/pkg/config/data/convert"
	"github.com/bhojpur/os/pkg/config/data/values"
)

type AnnotationField struct {
	Field            string
	Object           bool
	List             bool
	IgnoreDefinition bool
}

func (e AnnotationField) FromInternal(data map[string]interface{}) {
	v, ok := values.RemoveValue(data, "annotations", "field.bhojpur.net/"+e.Field)
	if ok {
		if e.Object {
			data := map[string]interface{}{}
			//ignore error
			if err := json.Unmarshal([]byte(convert.ToString(v)), &data); err == nil {
				v = data
			}
		}
		if e.List {
			var data []interface{}
			if err := json.Unmarshal([]byte(convert.ToString(v)), &data); err == nil {
				v = data
			}
		}

		data[e.Field] = v
	}
}

func (e AnnotationField) ToInternal(data map[string]interface{}) error {
	v, ok := data[e.Field]
	if ok {
		if e.Object || e.List {
			if bytes, err := json.Marshal(v); err == nil {
				v = string(bytes)
			}
		}
		values.PutValue(data, convert.ToString(v), "annotations", "field.bhojpur.net/"+e.Field)
	}
	values.RemoveValue(data, e.Field)
	return nil
}

func (e AnnotationField) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	if e.IgnoreDefinition {
		return nil
	}
	return ValidateField(e.Field, schema)
}
