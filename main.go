package main

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
	"os"
	"path/filepath"

	"github.com/bhojpur/os/cmd/app"
	"github.com/bhojpur/os/pkg/enterchroot"
	"github.com/bhojpur/os/pkg/transferroot"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/reexec"
	"github.com/sirupsen/logrus"
)

func main() {
	reexec.Register("/init", initrd)      // mode=live
	reexec.Register("/sbin/init", initrd) // mode=local
	reexec.Register("enter-root", enterchroot.Enter)

	if !reexec.Init() {
		app := app.New()
		args := []string{app.Name}
		path := filepath.Base(os.Args[0])
		if path != app.Name && app.Command(path) != nil {
			args = append(args, path)
		}
		args = append(args, os.Args[1:]...)
		// this will bomb if the application has any non-defaulted, required flags
		err := app.Run(args)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func initrd() {
	enterchroot.DebugCmdline = "bos.debug"
	transferroot.Relocate()
	if err := mount.Mount("", "/", "none", "rw,remount"); err != nil {
		logrus.Errorf("failed to remount root as rw: %v", err)
	}
	if err := enterchroot.Mount("./bos/data", os.Args, os.Stdout, os.Stderr); err != nil {
		logrus.Fatalf("failed to enter root: %v", err)
	}
}
