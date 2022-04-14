package upgrade

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
	"os/exec"
	"path/filepath"

	"github.com/bhojpur/os/pkg/system"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/sys/unix"
)

var (
	upgradeBOS, upgradeDCP              bool
	upgradeKernel, upgradeRootFS        bool
	doRemount, doSync, doReboot         bool
	sourceDir, destinationDir, lockFile string
)

// Command is the `upgrade` sub-command, it performs upgrades to Bhojpur OS.
func Command() cli.Command {
	return cli.Command{
		Name:  "upgrade",
		Usage: "perform upgrades",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "bos",
				EnvVar:      "BOS_UPGRADE_BOS",
				Destination: &upgradeBOS,
				Hidden:      true,
			},
			cli.BoolFlag{
				Name:        "dcp",
				EnvVar:      "BOS_UPGRADE_DCP",
				Destination: &upgradeDCP,
				Hidden:      true,
			},
			cli.BoolFlag{
				Name:        "kernel",
				Usage:       "upgrade the kernel",
				EnvVar:      "BOS_UPGRADE_KERNEL",
				Destination: &upgradeKernel,
			},
			cli.BoolFlag{
				Name:        "rootfs",
				Usage:       "upgrade bos+dcp",
				EnvVar:      "BOS_UPGRADE_ROOTFS",
				Destination: &upgradeRootFS,
			},
			cli.BoolFlag{
				Name:        "remount",
				Usage:       "pre-upgrade remount?",
				EnvVar:      "BOS_UPGRADE_REMOUNT",
				Destination: &doRemount,
			},
			cli.BoolFlag{
				Name:        "sync",
				Usage:       "post-upgrade sync?",
				EnvVar:      "BOS_UPGRADE_SYNC",
				Destination: &doSync,
			},
			cli.BoolFlag{
				Name:        "reboot",
				Usage:       "post-upgrade reboot?",
				EnvVar:      "BOS_UPGRADE_REBOOT",
				Destination: &doReboot,
			},
			cli.StringFlag{
				Name:        "source",
				EnvVar:      "BOS_UPGRADE_SOURCE",
				Value:       system.RootPath(),
				Required:    true,
				Destination: &sourceDir,
			},
			cli.StringFlag{
				Name:        "destination",
				EnvVar:      "BOS_UPGRADE_DESTINATION",
				Value:       system.RootPath(),
				Required:    true,
				Destination: &destinationDir,
			},
			cli.StringFlag{
				Name:        "lock-file",
				EnvVar:      "BOS_UPGRADE_LOCK_FILE",
				Value:       system.StatePath("upgrade.lock"),
				Hidden:      true,
				Destination: &lockFile,
			},
		},
		Before: func(c *cli.Context) error {
			if destinationDir == sourceDir {
				cli.ShowSubcommandHelp(c)
				logrus.Errorf("the `destination` cannot be the `source`: %s", destinationDir)
				os.Exit(1)
			}
			if upgradeRootFS {
				upgradeDCP = true
				upgradeBOS = true
			}
			if !upgradeBOS && !upgradeDCP && !upgradeKernel {
				cli.ShowSubcommandHelp(c)
				logrus.Error("must specify components to upgrade, e.g. `rootfs`, `kernel`")
				os.Exit(1)
			}
			return nil
		},
		Action: Run,
	}
}

// Run the `upgrade` sub-command
func Run(_ *cli.Context) {
	if err := validateSystemRoot(sourceDir); err != nil {
		logrus.Fatal(err)
	}
	if err := validateSystemRoot(destinationDir); err != nil {
		logrus.Fatal(err)
	}

	// establish the lock
	lf, err := os.OpenFile(lockFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	defer lf.Close()
	if err = unix.Flock(int(lf.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		logrus.Fatal(err)
	}
	defer unix.Flock(int(lf.Fd()), unix.LOCK_UN)

	var atLeastOneComponentCopied bool

	if upgradeBOS {
		if copied, err := system.CopyComponent(sourceDir, destinationDir, doRemount, "bos"); err != nil {
			logrus.Error(err)
		} else if copied {
			atLeastOneComponentCopied = true
			doRemount = false
		}
	}
	if upgradeDCP {
		if copied, err := system.CopyComponent(sourceDir, destinationDir, doRemount, "dcp"); err != nil {
			logrus.Error(err)
		} else if copied {
			atLeastOneComponentCopied = true
			doRemount = false
		}
	}
	if upgradeKernel {
		if copied, err := system.CopyComponent(sourceDir, destinationDir, doRemount, "kernel"); err != nil {
			logrus.Error(err)
		} else if copied {
			atLeastOneComponentCopied = true
			doRemount = false
		}
	}

	if atLeastOneComponentCopied && doSync {
		unix.Sync()
	}

	if atLeastOneComponentCopied && doReboot {
		// nsenter -m -u -i -n -p -t 1 -- reboot
		if _, err := exec.LookPath("nsenter"); err != nil {
			logrus.Warn(err)
			if destinationDir != system.RootPath() {
				root := filepath.Clean(filepath.Join(destinationDir, "..", ".."))
				logrus.Debugf("attempting chroot: %v", root)
				if err := unix.Chroot(root); err != nil {
					logrus.Fatal(err)
				}
				if err := os.Chdir("/"); err != nil {
					logrus.Fatal(err)
				}
			}
		}
		cmd := exec.Command("nsenter", "-m", "-u", "-i", "-n", "-p", "-t", "1", "reboot")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			logrus.Fatal(err)
		}
	}
}

func validateSystemRoot(root string) error {
	info, err := os.Stat(root)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("stat %s: not a directory", root)
	}
	return nil
}
