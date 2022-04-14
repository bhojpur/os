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

type Copy struct {
	From, To string
}

func (c Copy) FromInternal(data map[string]interface{}) {
	if data == nil {
		return
	}
	v, ok := data[c.From]
	if ok {
		data[c.To] = v
	}
}

func (c Copy) ToInternal(data map[string]interface{}) error {
	if data == nil {
		return nil
	}
	t, tok := data[c.To]
	_, fok := data[c.From]
	if tok && !fok {
		data[c.From] = t
	}

	return nil
}

func (c Copy) ModifySchema(s *types.Schema, schemas *types.Schemas) error {
	f, ok := s.ResourceFields[c.From]
	if !ok {
		return fmt.Errorf("field %s missing on schema %s", c.From, s.ID)
	}

	s.ResourceFields[c.To] = f
	return nil
}
