package data

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
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/bhojpur/host/pkg/common/name"
	"github.com/bhojpur/os/pkg/config/data/convert"
)

type SchemaCollection struct {
	Data []Schema
}

type SchemasInitFunc func(*Schemas) *Schemas

type MappersFactory func() []Mapper

type Schemas struct {
	sync.Mutex
	processingTypes    map[reflect.Type]*Schema
	typeNames          map[reflect.Type]string
	schemas            map[string]*Schema
	mappers            map[string][]Mapper
	DefaultMappers     MappersFactory
	DefaultPostMappers MappersFactory
	errors             []error
}

func NewSchemas() *Schemas {
	return &Schemas{
		processingTypes: map[reflect.Type]*Schema{},
		typeNames:       map[reflect.Type]string{},
		schemas:         map[string]*Schema{},
		mappers:         map[string][]Mapper{},
	}
}

func (s *Schemas) Init(initFunc SchemasInitFunc) *Schemas {
	return initFunc(s)
}

func (s *Schemas) Err() error {
	return NewErrors(s.errors...)
}

func (s *Schemas) AddSchemas(schema *Schemas) *Schemas {
	for _, schema := range schema.Schemas() {
		s.AddSchema(*schema)
	}
	return s
}

func (s *Schemas) RemoveSchema(schema Schema) *Schemas {
	s.Lock()
	defer s.Unlock()
	return s.doRemoveSchema(schema)
}

func (s *Schemas) doRemoveSchema(schema Schema) *Schemas {
	delete(s.schemas, schema.ID)
	return s
}

func (s *Schemas) AddSchema(schema Schema) *Schemas {
	s.Lock()
	defer s.Unlock()
	return s.doAddSchema(schema)
}

func (s *Schemas) doAddSchema(schema Schema) *Schemas {
	s.setupDefaults(&schema)
	s.schemas[schema.ID] = &schema
	return s
}

func (s *Schemas) setupDefaults(schema *Schema) {
	schema.Type = "/meta/schemas/schema"
	if schema.ID == "" {
		s.errors = append(s.errors, fmt.Errorf("ID is not set on schema: %v", schema))
		return
	}
	if schema.PluralName == "" {
		schema.PluralName = name.GuessPluralName(schema.ID)
	}
	if schema.CodeName == "" {
		schema.CodeName = convert.Capitalize(schema.ID)
	}
	if schema.CodeNamePlural == "" {
		schema.CodeNamePlural = name.GuessPluralName(schema.CodeName)
	}
}

func (s *Schemas) AddMapper(schemaID string, mapper Mapper) *Schemas {
	s.mappers[schemaID] = append(s.mappers[schemaID], mapper)
	return s
}

func (s *Schemas) Schemas() map[string]*Schema {
	return s.schemas
}

func (s *Schemas) Schema(name string) *Schema {
	return s.doSchema(name, true)
}

func (s *Schemas) doSchema(name string, lock bool) *Schema {
	if lock {
		s.Lock()
	}
	schema := s.schemas[name]
	if lock {
		s.Unlock()
	}

	if schema != nil {
		return schema
	}

	for _, check := range s.schemas {
		if strings.EqualFold(check.ID, name) || strings.EqualFold(check.PluralName, name) {
			return check
		}
	}

	return nil
}

type Errors []error

func (e Errors) Err() error {
	return NewErrors(e...)
}

func (e Errors) Error() string {
	buf := bytes.NewBuffer(nil)
	for _, err := range e {
		if buf.Len() > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(err.Error())
	}

	return buf.String()
}

func NewErrors(inErrors ...error) error {
	var errors []error
	for _, err := range inErrors {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) == 0 {
		return nil
	} else if len(errors) == 1 {
		return errors[0]
	}
	return Errors(errors)
}
