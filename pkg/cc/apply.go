package cc

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
	"github.com/bhojpur/os/pkg/config"
	"github.com/urfave/cli"
)

type applier func(cfg *config.CloudConfig) error

func runApplies(cfg *config.CloudConfig, appliers ...applier) error {
	var errors []error

	for _, a := range appliers {
		err := a(cfg)
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return cli.NewMultiError(errors...)
	}

	return nil
}

func RunApply(cfg *config.CloudConfig) error {
	return runApplies(cfg,
		ApplyModules,
		ApplySysctls,
		ApplyHostname,
		ApplyDNS,
		ApplyWifi,
		ApplyPassword,
		ApplySSHKeysWithNet,
		ApplyWriteFiles,
		ApplyEnvironment,
		ApplyRuncmd,
		ApplyInstall,
		ApplyDCPInstall,
	)
}

func InstallApply(cfg *config.CloudConfig) error {
	return runApplies(cfg,
		ApplyDCPWithRestart,
	)
}

func BootApply(cfg *config.CloudConfig) error {
	return runApplies(cfg,
		ApplyDataSource,
		ApplyModules,
		ApplySysctls,
		ApplyHostname,
		ApplyDNS,
		ApplyWifi,
		ApplyPassword,
		ApplySSHKeys,
		ApplyDCPNoRestart,
		ApplyWriteFiles,
		ApplyEnvironment,
		ApplyBootcmd,
	)
}

func InitApply(cfg *config.CloudConfig) error {
	return runApplies(cfg,
		ApplyModules,
		ApplySysctls,
		ApplyHostname,
		ApplyWriteFiles,
		ApplyEnvironment,
		ApplyInitcmd,
	)
}
