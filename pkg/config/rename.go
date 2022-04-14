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
	"strings"

	mapper "github.com/bhojpur/os/pkg/config/data"
	"github.com/bhojpur/os/pkg/config/data/convert"
	"github.com/bhojpur/os/pkg/config/data/mappers"
)

type FuzzyNames struct {
	mappers.DefaultMapper
	names map[string]string
}

func (f *FuzzyNames) ToInternal(data map[string]interface{}) error {
	for k, v := range data {
		if newK, ok := f.names[k]; ok && newK != k {
			data[newK] = v
		}
	}
	return nil
}

func (f *FuzzyNames) addName(name, toName string) {
	f.names[strings.ToLower(name)] = toName
	f.names[convert.ToYAMLKey(name)] = toName
	f.names[strings.ToLower(convert.ToYAMLKey(name))] = toName
}

func (f *FuzzyNames) ModifySchema(schema *mapper.Schema, schemas *mapper.Schemas) error {
	f.names = map[string]string{}

	for name := range schema.ResourceFields {
		if strings.HasSuffix(name, "s") && len(name) > 1 {
			f.addName(name[:len(name)-1], name)
		}
		if strings.HasSuffix(name, "es") && len(name) > 2 {
			f.addName(name[:len(name)-2], name)
		}
		f.addName(name, name)
	}

	f.names["pass"] = "passphrase"
	f.names["password"] = "passphrase"

	return nil
}
