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
	"github.com/docker/go-units"
)

var abbrs = []string{"", "k", "m", "g", "t", "p"}

type Bytes struct {
	Field string
}

func (d Bytes) FromInternal(data map[string]interface{}) {
	v, ok := data[d.Field]
	if !ok {
		return
	}

	n, err := convert.ToNumber(v)
	if err != nil {
		return
	}

	data[d.Field] = units.CustomSize("%.4g%s", float64(n), 1024.0, abbrs)
}

func (d Bytes) ToInternal(data map[string]interface{}) error {
	v, ok := data[d.Field]
	if !ok {
		return nil
	}

	if str, ok := v.(string); ok {
		sec, err := units.RAMInBytes(str)
		if err != nil {
			return err
		}
		data[d.Field] = sec
	}

	return nil
}

func (d Bytes) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	return ValidateField(d.Field, schema)
}
