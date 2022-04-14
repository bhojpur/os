package writefile

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
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/bhojpur/os/pkg/config"
	"github.com/bhojpur/os/pkg/util"
	"github.com/sirupsen/logrus"
)

func WriteFiles(cfg *config.CloudConfig) {
	for i, f := range cfg.WriteFiles {
		c, err := util.DecodeContent(f.Content, f.Encoding)
		if err != nil {
			logrus.Errorf("failed to decode content from write_files item [%d]: %v", i, err)
			continue
		}
		f.Content = string(c)
		f.Encoding = ""
		p, err := WriteFile(&f, "/")
		if err != nil {
			logrus.WithFields(logrus.Fields{"err": err, "path": p}).Errorln("failed to write file")
			continue
		}
		logrus.Infof("wrote file %s to filesystem", p)
	}
}

func WriteFile(f *config.File, root string) (string, error) {
	if f.Encoding != "" {
		return "", fmt.Errorf("unable to write file with encoding %s", f.Encoding)
	}
	p := path.Join(root, f.Path)
	d := path.Dir(p)
	logrus.Infof("writing file to %q", d)
	if err := util.EnsureDirectoryExists(d); err != nil {
		return "", err
	}
	perm, err := f.Permissions()
	if err != nil {
		return "", err
	}
	var tmp *os.File
	// create a temporary file in the same directory to ensure it's on the same filesystem
	if tmp, err = ioutil.TempFile(d, "wfs-temp"); err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(tmp.Name(), []byte(f.Content), perm); err != nil {
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	// ensure the permissions are as requested (since WriteFile can be affected by sticky bit)
	if err := os.Chmod(tmp.Name(), perm); err != nil {
		return "", err
	}
	if f.Owner != "" {
		// we shell out since we don't have a way to look up unix groups natively
		cmd := exec.Command("chown", f.Owner, tmp.Name())
		if err := cmd.Run(); err != nil {
			return "", err
		}
	}
	if err := os.Rename(tmp.Name(), p); err != nil {
		return "", err
	}
	return p, nil
}
