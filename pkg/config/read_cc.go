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
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bhojpur/os/pkg/config/data/convert"
	"github.com/ghodss/yaml"
)

const (
	hostname = "/run/config/local_hostname"
	ssh      = "/run/config/ssh/authorized_keys"
	userdata = "/run/config/userdata"
)

func readCloudConfig() (map[string]interface{}, error) {
	var keys []string
	result := map[string]interface{}{}

	hostname, err := ioutil.ReadFile(hostname)
	if err == nil {
		result["hostname"] = strings.TrimSpace(string(hostname))
	}

	keyData, err := ioutil.ReadFile(ssh)
	if err != nil {
		// ignore error
		return result, nil
	}

	for _, line := range strings.Split(string(keyData), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			keys = append(keys, line)
		}
	}

	if len(keys) > 0 {
		result["ssh_authorized_keys"] = keys
	}

	return result, nil
}

func readUserData() (map[string]interface{}, error) {
	result := map[string]interface{}{}

	data, err := ioutil.ReadFile(userdata)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	cc := CloudConfig{}
	script := false
	if bytes.Contains(data, []byte{0}) {
		script = true
		cc.WriteFiles = []File{
			{
				Content:  base64.StdEncoding.EncodeToString(data),
				Encoding: "b64",
			},
		}
	} else if strings.HasPrefix(string(data), "#!") {
		script = true
		cc.WriteFiles = []File{
			{
				Content: string(data),
			},
		}
	}

	if script {
		cc.WriteFiles[0].Owner = "root"
		cc.WriteFiles[0].RawFilePermissions = "0700"
		cc.WriteFiles[0].Path = "/run/bos/userdata"
		cc.Runcmd = []string{"source /run/bos/userdata"}

		return convert.EncodeToMap(cc)
	}
	return result, yaml.Unmarshal(data, &result)
}
