package data

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
	"encoding/json"
	"io"
	"regexp"

	"github.com/ghodss/yaml"
)

var (
	commenter = regexp.MustCompile("(?m)^( *)zzz#\\((.*)\\)\\((.*)\\)([a-z]+.*):(.*)")
)

func JSONEncoder(writer io.Writer, v interface{}) error {
	return json.NewEncoder(writer).Encode(v)
}

func YAMLEncoder(writer io.Writer, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	buf, err := yaml.JSONToYAML(data)
	if err != nil {
		return err
	}
	//buf = commenter.ReplaceAll(buf, []byte("${1}# ${2}type: ${3}\n${1}# ${4}:${5}"))
	buf = commenter.ReplaceAll(buf, []byte("${1}# ${4}:${5}"))
	_, err = writer.Write(buf)
	return err
}
