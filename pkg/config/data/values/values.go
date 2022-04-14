package values

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

func RemoveValue(data map[string]interface{}, keys ...string) (interface{}, bool) {
	for i, key := range keys {
		if i == len(keys)-1 {
			val, ok := data[key]
			delete(data, key)
			return val, ok
		}
		data, _ = data[key].(map[string]interface{})
	}

	return nil, false
}

func GetValueN(data map[string]interface{}, keys ...string) interface{} {
	val, _ := GetValue(data, keys...)
	return val
}

func GetValue(data map[string]interface{}, keys ...string) (interface{}, bool) {
	for i, key := range keys {
		if i == len(keys)-1 {
			val, ok := data[key]
			return val, ok
		}
		data, _ = data[key].(map[string]interface{})
	}

	return nil, false
}

func PutValue(data map[string]interface{}, val interface{}, keys ...string) {
	if data == nil {
		return
	}

	// This is so ugly
	for i, key := range keys {
		if i == len(keys)-1 {
			data[key] = val
		} else {
			newData, ok := data[key]
			if ok {
				newMap, ok := newData.(map[string]interface{})
				if ok {
					data = newMap
				} else {
					return
				}
			} else {
				newMap := map[string]interface{}{}
				data[key] = newMap
				data = newMap
			}
		}
	}
}
