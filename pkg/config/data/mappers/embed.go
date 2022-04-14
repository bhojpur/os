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
)

type Embed struct {
	Field          string
	Optional       bool
	ReadOnly       bool
	Ignore         []string
	ignoreOverride bool
	embeddedFields []string
	EmptyValueOk   bool
}

func (e *Embed) FromInternal(data map[string]interface{}) {
	sub, _ := data[e.Field].(map[string]interface{})
	for _, fieldName := range e.embeddedFields {
		if v, ok := sub[fieldName]; ok {
			data[fieldName] = v
		}
	}
	delete(data, e.Field)
}

func (e *Embed) ToInternal(data map[string]interface{}) error {
	if data == nil {
		return nil
	}

	sub := map[string]interface{}{}
	for _, fieldName := range e.embeddedFields {
		if v, ok := data[fieldName]; ok {
			sub[fieldName] = v
		}

		delete(data, fieldName)
	}
	if len(sub) == 0 {
		if e.EmptyValueOk {
			data[e.Field] = nil
		}
		return nil
	}
	data[e.Field] = sub
	return nil
}

func (e *Embed) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	err := ValidateField(e.Field, schema)
	if err != nil {
		if e.Optional {
			return nil
		}
		return err
	}

	e.embeddedFields = []string{}

	embeddedSchemaID := schema.ResourceFields[e.Field].Type
	embeddedSchema := schemas.Schema(embeddedSchemaID)
	if embeddedSchema == nil {
		if e.Optional {
			return nil
		}
		return fmt.Errorf("failed to find schema %s for embedding", embeddedSchemaID)
	}

	deleteField := true
outer:
	for name, field := range embeddedSchema.ResourceFields {
		for _, ignore := range e.Ignore {
			if ignore == name {
				continue outer
			}
		}

		if name == e.Field {
			deleteField = false
		} else {
			if !e.ignoreOverride {
				if _, ok := schema.ResourceFields[name]; ok {
					return fmt.Errorf("embedding field %s on %s will overwrite the field %s",
						e.Field, schema.ID, name)
				}
			}
		}

		if e.ReadOnly {
			field.Create = false
			field.Update = false
		}

		schema.ResourceFields[name] = field
		e.embeddedFields = append(e.embeddedFields, name)
	}

	if deleteField {
		delete(schema.ResourceFields, e.Field)
	}

	return nil
}
