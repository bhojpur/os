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
	types "github.com/bhojpur/os/pkg/config/data"
	"github.com/bhojpur/os/pkg/config/data/values"
)

type LabelField struct {
	Field string
}

func (e LabelField) FromInternal(data map[string]interface{}) {
	v, ok := values.RemoveValue(data, "labels", "field.bhojpur.net/"+e.Field)
	if ok {
		data[e.Field] = v
	}
}

func (e LabelField) ToInternal(data map[string]interface{}) error {
	v, ok := data[e.Field]
	if ok {
		values.PutValue(data, v, "labels", "field.bhojpur.net/"+e.Field)
	}
	return nil
}

func (e LabelField) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	return ValidateField(e.Field, schema)
}
