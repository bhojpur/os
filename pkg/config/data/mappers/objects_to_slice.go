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
	"github.com/sirupsen/logrus"
)

type MaybeStringer interface {
	MaybeString() interface{}
}

type StringerFactory func() MaybeStringer
type ToObject func(interface{}) (interface{}, error)

type ObjectsToSlice struct {
	Field     string
	NewObject StringerFactory
	ToObject  ToObject
}

func (p ObjectsToSlice) FromInternal(data map[string]interface{}) {
	if data == nil {
		return
	}

	objs, ok := data[p.Field]
	if !ok {
		return
	}

	var result []interface{}
	for _, obj := range convert.ToMapSlice(objs) {
		target := p.NewObject()
		if err := convert.ToObj(obj, target); err != nil {
			logrus.Errorf("Failed to unmarshal slice to object: %v", err)
			continue
		}

		ret := target.MaybeString()
		if slc, ok := ret.([]string); ok {
			for _, v := range slc {
				result = append(result, v)
			}
		} else {
			result = append(result, ret)
		}
	}

	if len(result) == 0 {
		delete(data, p.Field)
	} else {
		data[p.Field] = result
	}
}

func (p ObjectsToSlice) ToInternal(data map[string]interface{}) error {
	if data == nil {
		return nil
	}

	d, ok := data[p.Field]
	if !ok {
		return nil
	}

	if str, ok := d.(string); ok {
		d = []interface{}{str}
	}

	slc, ok := d.([]interface{})
	if !ok {
		return nil
	}

	var newSlc []interface{}

	for _, obj := range slc {
		n, err := convert.ToNumber(obj)
		if err == nil && n > 0 {
			obj = convert.ToString(n)
		}
		newObj, err := p.ToObject(obj)
		if err != nil {
			return err
		}

		if mapSlice, isMapSlice := newObj.([]map[string]interface{}); isMapSlice {
			for _, v := range mapSlice {
				newSlc = append(newSlc, v)
			}
		} else {
			if _, isMap := newObj.(map[string]interface{}); !isMap {
				newObj, err = convert.EncodeToMap(newObj)
			}

			newSlc = append(newSlc, newObj)
		}
	}

	data[p.Field] = newSlc
	return nil
}

func (p ObjectsToSlice) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	return ValidateField(p.Field, schema)
}
