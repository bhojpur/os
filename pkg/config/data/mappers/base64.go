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
	"encoding/base64"
	"strings"

	mapper "github.com/bhojpur/os/pkg/config/data"
	"github.com/bhojpur/os/pkg/config/data/convert"
	"github.com/bhojpur/os/pkg/config/data/values"
	"github.com/sirupsen/logrus"
)

type Base64 struct {
	Field            string
	IgnoreDefinition bool
	Separator        string
}

func (m Base64) FromInternal(data map[string]interface{}) {
	if v, ok := values.RemoveValue(data, strings.Split(m.Field, m.getSep())...); ok {
		str := convert.ToString(v)
		if str == "" {
			return
		}

		newData, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			logrus.Errorf("failed to base64 decode data")
		}

		values.PutValue(data, string(newData), strings.Split(m.Field, m.getSep())...)
	}
}

func (m Base64) ToInternal(data map[string]interface{}) error {
	if v, ok := values.RemoveValue(data, strings.Split(m.Field, m.getSep())...); ok {
		str := convert.ToString(v)
		if str == "" {
			return nil
		}

		newData := base64.StdEncoding.EncodeToString([]byte(str))
		values.PutValue(data, newData, strings.Split(m.Field, m.getSep())...)
	}

	return nil
}

func (m Base64) ModifySchema(s *mapper.Schema, schemas *mapper.Schemas) error {
	if !m.IgnoreDefinition {
		if err := ValidateField(m.Field, s); err != nil {
			return err
		}
	}

	return nil
}

func (m Base64) getSep() string {
	if m.Separator == "" {
		return "/"
	}
	return m.Separator
}
