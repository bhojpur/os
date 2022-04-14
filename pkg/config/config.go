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
	"fmt"
	"os"
	"strconv"
)

type BhojpurOS struct {
	DataSources    []string          `json:"dataSources,omitempty"`
	Modules        []string          `json:"modules,omitempty"`
	Sysctls        map[string]string `json:"sysctls,omitempty"`
	NTPServers     []string          `json:"ntpServers,omitempty"`
	DNSNameservers []string          `json:"dnsNameservers,omitempty"`
	Wifi           []Wifi            `json:"wifi,omitempty"`
	Password       string            `json:"password,omitempty"`
	ServerURL      string            `json:"serverUrl,omitempty"`
	Token          string            `json:"token,omitempty"`
	Labels         map[string]string `json:"labels,omitempty"`
	DcpArgs        []string          `json:"dcpArgs,omitempty"`
	Environment    map[string]string `json:"environment,omitempty"`
	Taints         []string          `json:"taints,omitempty"`
	Install        *Install          `json:"install,omitempty"`
}

type Wifi struct {
	Name       string `json:"name,omitempty"`
	Passphrase string `json:"passphrase,omitempty"`
}

type Install struct {
	ForceEFI  bool   `json:"forceEfi,omitempty"`
	Device    string `json:"device,omitempty"`
	ConfigURL string `json:"configUrl,omitempty"`
	Silent    bool   `json:"silent,omitempty"`
	ISOURL    string `json:"isoUrl,omitempty"`
	PowerOff  bool   `json:"powerOff,omitempty"`
	NoFormat  bool   `json:"noFormat,omitempty"`
	Debug     bool   `json:"debug,omitempty"`
	TTY       string `json:"tty,omitempty"`
}

type CloudConfig struct {
	SSHAuthorizedKeys []string  `json:"sshAuthorizedKeys,omitempty"`
	WriteFiles        []File    `json:"writeFiles,omitempty"`
	Hostname          string    `json:"hostname,omitempty"`
	BhojpurOS         BhojpurOS `json:"bos,omitempty"`
	Runcmd            []string  `json:"runCmd,omitempty"`
	Bootcmd           []string  `json:"bootCmd,omitempty"`
	Initcmd           []string  `json:"initCmd,omitempty"`
}

type File struct {
	Encoding           string `json:"encoding"`
	Content            string `json:"content"`
	Owner              string `json:"owner"`
	Path               string `json:"path"`
	RawFilePermissions string `json:"permissions"`
}

func (f *File) Permissions() (os.FileMode, error) {
	if f.RawFilePermissions == "" {
		return os.FileMode(0644), nil
	}
	// parse string representation of file mode as integer
	perm, err := strconv.ParseInt(f.RawFilePermissions, 8, 32)
	if err != nil {
		return 0, fmt.Errorf("unable to parse file permissions %q as integer", f.RawFilePermissions)
	}
	return os.FileMode(perm), nil
}
