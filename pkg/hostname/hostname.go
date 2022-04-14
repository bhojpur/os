package hostname

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
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/bhojpur/os/pkg/config"
)

func SetHostname(c *config.CloudConfig) error {
	hostname := c.Hostname
	if hostname == "" {
		return nil
	}
	if err := syscall.Sethostname([]byte(hostname)); err != nil {
		return err
	}
	return syncHostname()
}

func syncHostname() error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	if hostname == "" {
		return nil
	}

	if err := ioutil.WriteFile("/etc/hostname", []byte(hostname+"\n"), 0644); err != nil {
		return err
	}

	hosts, err := os.Open("/etc/hosts")
	defer hosts.Close()
	if err != nil {
		return err
	}
	lines := bufio.NewScanner(hosts)
	content := ""
	for lines.Scan() {
		line := strings.TrimSpace(lines.Text())
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == "127.0.1.1" {
			content += "127.0.1.1 " + hostname + "\n"
			continue
		}
		content += line + "\n"
	}
	return ioutil.WriteFile("/etc/hosts", []byte(content), 0600)
}
