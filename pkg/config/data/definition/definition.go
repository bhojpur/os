package definition

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

	"github.com/bhojpur/os/pkg/config/data/convert"
)

func IsMapType(fieldType string) bool {
	return strings.HasPrefix(fieldType, "map[") && strings.HasSuffix(fieldType, "]")
}

func IsArrayType(fieldType string) bool {
	return strings.HasPrefix(fieldType, "array[") && strings.HasSuffix(fieldType, "]")
}

func IsReferenceType(fieldType string) bool {
	return strings.HasPrefix(fieldType, "reference[") && strings.HasSuffix(fieldType, "]")
}

func SubType(fieldType string) string {
	i := strings.Index(fieldType, "[")
	if i <= 0 || i >= len(fieldType)-1 {
		return fieldType
	}

	return fieldType[i+1 : len(fieldType)-1]
}

func GetType(data map[string]interface{}) string {
	return GetShortTypeFromFull(GetFullType(data))
}

func GetShortTypeFromFull(fullType string) string {
	parts := strings.Split(fullType, "/")
	return parts[len(parts)-1]
}

func GetFullType(data map[string]interface{}) string {
	return convert.ToString(data["type"])
}
