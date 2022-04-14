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
	"strings"

	types "github.com/bhojpur/os/pkg/config/data"
)

type Access struct {
	Fields   map[string]string
	Optional bool
}

func (e Access) FromInternal(data map[string]interface{}) {
}

func (e Access) ToInternal(data map[string]interface{}) error {
	return nil
}

func (e Access) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	for name, access := range e.Fields {
		if err := ValidateField(name, schema); err != nil {
			if e.Optional {
				continue
			}
			return err
		}

		field := schema.ResourceFields[name]
		field.Create = strings.Contains(access, "c")
		field.Update = strings.Contains(access, "u")
		field.WriteOnly = strings.Contains(access, "o")

		schema.ResourceFields[name] = field
	}
	return nil
}
