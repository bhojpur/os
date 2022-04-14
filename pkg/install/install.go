package install

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

	"github.com/bhojpur/os/pkg/config"
	"github.com/bhojpur/os/pkg/questions"
	"github.com/ghodss/yaml"
)

func Run() error {
	fmt.Println("\nRunning Bhojpur OS configuration")

	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	isInstall, err := Ask(&cfg)
	if err != nil {
		return err
	}

	if isInstall {
		return runInstall(cfg)
	}

	bytes, err := config.ToBytes(cfg)
	if err != nil {
		return err
	}

	f, err := os.Create(config.SystemConfig)
	if err != nil {
		f, err = os.Create(config.LocalConfig)
		if err != nil {
			return err
		}
	}
	defer f.Close()

	if _, err := f.Write(bytes); err != nil {
		return err
	}

	f.Close()
	return runCCApply()
}

func runCCApply() error {
	cmd := exec.Command(os.Args[0], "config", "--install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func runInstall(cfg config.CloudConfig) error {
	var (
		err      error
		tempFile *os.File
	)

	installBytes, err := config.PrintInstall(cfg)
	if err != nil {
		return err
	}

	if !cfg.BhojpurOS.Install.Silent {
		val, err := questions.PromptBool("\nConfiguration\n"+"-------------\n\n"+
			string(installBytes)+
			"\nYour disk will be formatted and Bhojpur OS will be installed with the above configuration.\nContinue?", false)
		if err != nil || !val {
			return err
		}
	}

	if cfg.BhojpurOS.Install.ConfigURL == "" {
		tempFile, err = ioutil.TempFile("/tmp", "bos.XXXXXXXX")
		if err != nil {
			return err
		}
		defer tempFile.Close()

		cfg.BhojpurOS.Install.ConfigURL = tempFile.Name()
	}

	ev, err := config.ToEnv(cfg)
	if err != nil {
		return err
	}

	if tempFile != nil {
		cfg.BhojpurOS.Install = nil
		bytes, err := yaml.Marshal(&cfg)
		if err != nil {
			return err
		}
		if _, err := tempFile.Write(bytes); err != nil {
			return err
		}
		if err := tempFile.Close(); err != nil {
			return err
		}
		defer os.Remove(tempFile.Name())
	}

	cmd := exec.Command("/usr/libexec/os/install")
	cmd.Env = append(os.Environ(), ev...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
