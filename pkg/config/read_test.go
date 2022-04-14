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

import "testing"

func TestDataSource(t *testing.T) {
	cc, err := readersToObject(func() (map[string]interface{}, error) {
		return map[string]interface{}{
			"bos": map[string]interface{}{
				"datasource": "foo",
			},
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(cc.BhojpurOS.DataSources) != 1 {
		t.Fatal("no datasources")
	}
	if cc.BhojpurOS.DataSources[0] != "foo" {
		t.Fatalf("%s != foo", cc.BhojpurOS.DataSources[0])
	}
}

func TestAuthorizedKeys(t *testing.T) {
	c1 := map[string]interface{}{
		"ssh_authorized_keys": []string{
			"one...",
		},
	}
	c2 := map[string]interface{}{
		"ssh_authorized_keys": []string{
			"two...",
		},
	}
	cc, err := readersToObject(
		func() (map[string]interface{}, error) {
			return c1, nil
		},
		func() (map[string]interface{}, error) {
			return c2, nil
		},
	)
	if len(cc.SSHAuthorizedKeys) != 1 {
		t.Fatal(err, "got %d keys, expected 2", len(cc.SSHAuthorizedKeys))
	}
	if err != nil {
		t.Fatal(err)
	}
}
