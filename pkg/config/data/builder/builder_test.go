package builder

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
	"testing"

	mapper "github.com/bhojpur/os/pkg/config/data"
	"github.com/stretchr/testify/assert"
)

func TestEmptyStringWithDefault(t *testing.T) {
	schema := &mapper.Schema{
		ResourceFields: map[string]mapper.Field{
			"foo": {
				Default: "foo",
				Type:    "string",
				Create:  true,
			},
		},
	}
	schemas := mapper.NewSchemas()
	schemas.AddSchema(*schema)

	builder := NewBuilder(schemas)

	// Test if no field we set to "foo"
	result, err := builder.Construct(schema, nil, Create)
	if err != nil {
		t.Fatal(err)
	}
	value, ok := result["foo"]
	assert.True(t, ok)
	assert.Equal(t, "foo", value)

	// Test if field is "" we set to "foo"
	result, err = builder.Construct(schema, map[string]interface{}{
		"foo": "",
	}, Create)
	if err != nil {
		t.Fatal(err)
	}
	value, ok = result["foo"]
	assert.True(t, ok)
	assert.Equal(t, "foo", value)
}
