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
	"github.com/bhojpur/os/pkg/config/data/convert"
)

type SliceMerge struct {
	From             []string
	To               string
	IgnoreDefinition bool
}

func (s SliceMerge) FromInternal(data map[string]interface{}) {
	var result []interface{}
	for _, name := range s.From {
		val, ok := data[name]
		if !ok {
			continue
		}
		result = append(result, convert.ToInterfaceSlice(val)...)
	}

	if result != nil {
		data[s.To] = result
	}
}

func (s SliceMerge) ToInternal(data map[string]interface{}) error {
	return nil
}

func (s SliceMerge) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	if s.IgnoreDefinition {
		return nil
	}

	for _, from := range s.From {
		if err := ValidateField(from, schema); err != nil {
			return err
		}
		if from != s.To {
			delete(schema.ResourceFields, from)
		}
	}

	return ValidateField(s.To, schema)
}
