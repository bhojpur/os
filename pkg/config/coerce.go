package config

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
	mapper "github.com/bhojpur/os/pkg/config/data"
	"github.com/bhojpur/os/pkg/config/data/convert"
	"github.com/bhojpur/os/pkg/config/data/mappers"
)

type Converter func(val interface{}) interface{}

type fieldConverter struct {
	mappers.DefaultMapper
	fieldName string
	converter Converter
}

func (f fieldConverter) ToInternal(data map[string]interface{}) error {
	val, ok := data[f.fieldName]
	if !ok {
		return nil
	}
	data[f.fieldName] = f.converter(val)
	return nil
}

type typeConverter struct {
	mappers.DefaultMapper
	converter Converter
	fieldType string
	mappers   mapper.Mappers
}

func (t *typeConverter) ToInternal(data map[string]interface{}) error {
	return t.mappers.ToInternal(data)
}

func (t *typeConverter) ModifySchema(schema *mapper.Schema, schemas *mapper.Schemas) error {
	for name, field := range schema.ResourceFields {
		if field.Type == t.fieldType {
			t.mappers = append(t.mappers, fieldConverter{
				fieldName: name,
				converter: t.converter,
			})
		}
	}
	return nil
}

func NewTypeConverter(fieldType string, converter Converter) mapper.Mapper {
	return &typeConverter{
		fieldType: fieldType,
		converter: converter,
	}
}

func NewToMap() mapper.Mapper {
	return NewTypeConverter("map[string]", func(val interface{}) interface{} {
		if m, ok := val.(map[string]interface{}); ok {
			obj := make(map[string]string, len(m))
			for k, v := range m {
				obj[k] = convert.ToString(v)
			}
			return obj
		}
		return val
	})
}

func NewToSlice() mapper.Mapper {
	return NewTypeConverter("array[string]", func(val interface{}) interface{} {
		if str, ok := val.(string); ok {
			return []string{str}
		}
		return val
	})
}

func NewToBool() mapper.Mapper {
	return NewTypeConverter("boolean", func(val interface{}) interface{} {
		if str, ok := val.(string); ok {
			return str == "true"
		}
		return val
	})
}
