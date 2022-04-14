package util

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
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"

	"github.com/sirupsen/logrus"
)

func DecodeBase64Content(content string) ([]byte, error) {
	output, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, fmt.Errorf("unable to decode base64: %q", err)
	}
	return output, nil
}

func DecodeGzipContent(content string) ([]byte, error) {
	byteContent := []byte(content)
	return DecompressGzip(byteContent)
}

func DecompressGzip(content []byte) ([]byte, error) {
	gzr, err := gzip.NewReader(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("unable to decode gzip: %q", err)
	}
	defer func() {
		if err := gzr.Close(); err != nil {
			logrus.Errorf("unable to close gzip reader: %q", err)
		}
	}()
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(gzr); err != nil {
		return nil, fmt.Errorf("unable to read gzip: %q", err)
	}
	return buf.Bytes(), nil
}

func DecodeContent(content string, encoding string) ([]byte, error) {
	switch encoding {
	case "":
		return []byte(content), nil
	case "b64", "base64":
		return DecodeBase64Content(content)
	case "gz", "gzip":
		return DecodeGzipContent(content)
	case "gz+base64", "gzip+base64", "gz+b64", "gzip+b64":
		gz, err := DecodeBase64Content(content)
		if err != nil {
			return nil, err
		}
		return DecodeGzipContent(string(gz))
	}
	return nil, fmt.Errorf("unsupported encoding %q", encoding)
}
