package convert

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

const (
	ArrayKey = "{ARRAY}"
	MapKey   = "{MAP}"
)

type TransformerFunc func(input interface{}) interface{}

func Transform(data map[string]interface{}, path []string, transformer TransformerFunc) {
	if len(path) == 0 || len(data) == 0 {
		return
	}

	key := path[0]
	path = path[1:]
	value := data[key]

	if value == nil {
		return
	}

	if len(path) == 0 {
		data[key] = transformer(value)
		return
	}

	// You can't end a path with ARRAY/MAP.  Not supported right now
	if len(path) > 1 {
		switch path[0] {
		case ArrayKey:
			for _, valueMap := range ToMapSlice(value) {
				Transform(valueMap, path[1:], transformer)
			}
			return
		case MapKey:
			for _, valueMap := range ToMapInterface(value) {
				Transform(ToMapInterface(valueMap), path[1:], transformer)
			}
			return
		}
	}

	Transform(ToMapInterface(value), path, transformer)
}
