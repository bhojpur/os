package module

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
	"fmt"
	"os"
	"strings"

	"github.com/bhojpur/os/pkg/config"
	"github.com/paultag/go-modprobe"
	"github.com/sirupsen/logrus"
)

const (
	procModulesFile = "/proc/modules"
)

func LoadModules(cfg *config.CloudConfig) error {
	loaded := map[string]bool{}
	f, err := os.Open(procModulesFile)
	if err != nil {
		return err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		loaded[strings.SplitN(sc.Text(), " ", 2)[0]] = true
	}
	modules := cfg.BhojpurOS.Modules
	for _, m := range modules {
		if loaded[m] {
			continue
		}
		params := strings.SplitN(m, " ", -1)
		logrus.Debugf("module %s with parameters [%s] is loading", m, params)
		if err := modprobe.Load(params[0], strings.Join(params[1:], " ")); err != nil {
			return fmt.Errorf("could not load module %s with parameters [%s], err %v", m, params, err)
		}
		logrus.Debugf("module %s is loaded", m)
	}
	return sc.Err()
}
