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
)

type ReadOnly struct {
	Field     string
	Optional  bool
	SubFields bool
}

func (r ReadOnly) FromInternal(data map[string]interface{}) {
}

func (r ReadOnly) ToInternal(data map[string]interface{}) error {
	return nil
}

func (r ReadOnly) readOnly(field types.Field, schema *types.Schema, schemas *types.Schemas) types.Field {
	field.Create = false
	field.Update = false

	if r.SubFields {
		subSchema := schemas.Schema(field.Type)
		if subSchema != nil {
			for name, field := range subSchema.ResourceFields {
				field.Create = false
				field.Update = false
				subSchema.ResourceFields[name] = field
			}
		}
	}

	return field
}

func (r ReadOnly) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	if r.Field == "*" {
		for name, field := range schema.ResourceFields {
			schema.ResourceFields[name] = r.readOnly(field, schema, schemas)
		}
		return nil
	}

	if err := ValidateField(r.Field, schema); err != nil {
		if r.Optional {
			return nil
		}
		return err
	}

	field := schema.ResourceFields[r.Field]
	schema.ResourceFields[r.Field] = r.readOnly(field, schema, schemas)

	return nil
}
