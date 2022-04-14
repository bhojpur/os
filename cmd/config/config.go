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
	"encoding/json"
	"fmt"
	"os"

	"github.com/bhojpur/os/pkg/cc"
	"github.com/bhojpur/os/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	initrd       = false
	bootPhase    = false
	installPhase = false
	dump         = false
	dumpJSON     = false
)

// Command `config`
func Command() cli.Command {
	return cli.Command{
		Name:      "config",
		Usage:     "configure Bhojpur OS",
		ShortName: "cfg",
		// Aliases: []string{
		// 	"ccapply",
		// },
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "initrd",
				Destination: &initrd,
				Usage:       "Run initrd stage",
			},
			cli.BoolFlag{
				Name:        "boot",
				Destination: &bootPhase,
				Usage:       "Run boot stage",
			},
			cli.BoolFlag{
				Name:        "install",
				Destination: &installPhase,
				Usage:       "Run install stage",
			},
			cli.BoolFlag{
				Name:        "dump",
				Destination: &dump,
				Usage:       "Print current configuration",
			},
			cli.BoolFlag{
				Name:        "dump-json",
				Destination: &dumpJSON,
				Usage:       "Print current configuration in json",
			},
		},
		Before: func(c *cli.Context) error {
			if os.Getuid() != 0 {
				return fmt.Errorf("must be run as root")
			}
			return nil
		},
		Action: func(*cli.Context) {
			if err := Main(); err != nil {
				logrus.Error(err)
			}
		},
	}
}

// Main `config`
func Main() error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	if initrd {
		return cc.InitApply(&cfg)
	} else if bootPhase {
		return cc.BootApply(&cfg)
	} else if installPhase {
		return cc.InstallApply(&cfg)
	} else if dump {
		return config.Write(cfg, os.Stdout)
	} else if dumpJSON {
		return json.NewEncoder(os.Stdout).Encode(&cfg)
	}

	return cc.RunApply(&cfg)
}
