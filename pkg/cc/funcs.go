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
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/bhojpur/os/pkg/command"
	"github.com/bhojpur/os/pkg/config"
	"github.com/bhojpur/os/pkg/hostname"
	"github.com/bhojpur/os/pkg/mode"
	"github.com/bhojpur/os/pkg/module"
	"github.com/bhojpur/os/pkg/ssh"
	"github.com/bhojpur/os/pkg/sysctl"
	"github.com/bhojpur/os/pkg/version"
	"github.com/bhojpur/os/pkg/writefile"
	"github.com/sirupsen/logrus"
)

func ApplyModules(cfg *config.CloudConfig) error {
	return module.LoadModules(cfg)
}

func ApplySysctls(cfg *config.CloudConfig) error {
	return sysctl.ConfigureSysctl(cfg)
}

func ApplyHostname(cfg *config.CloudConfig) error {
	return hostname.SetHostname(cfg)
}

func ApplyPassword(cfg *config.CloudConfig) error {
	return command.SetPassword(cfg.BhojpurOS.Password)
}

func ApplyRuncmd(cfg *config.CloudConfig) error {
	return command.ExecuteCommand(cfg.Runcmd)
}

func ApplyBootcmd(cfg *config.CloudConfig) error {
	return command.ExecuteCommand(cfg.Bootcmd)
}

func ApplyInitcmd(cfg *config.CloudConfig) error {
	return command.ExecuteCommand(cfg.Initcmd)
}

func ApplyWriteFiles(cfg *config.CloudConfig) error {
	writefile.WriteFiles(cfg)
	return nil
}

func ApplySSHKeys(cfg *config.CloudConfig) error {
	return ssh.SetAuthorizedKeys(cfg, false)
}

func ApplySSHKeysWithNet(cfg *config.CloudConfig) error {
	return ssh.SetAuthorizedKeys(cfg, true)
}

func ApplyDCPWithRestart(cfg *config.CloudConfig) error {
	return ApplyDCP(cfg, true, false)
}

func ApplyDCPInstall(cfg *config.CloudConfig) error {
	return ApplyDCP(cfg, true, true)
}

func ApplyDCPNoRestart(cfg *config.CloudConfig) error {
	return ApplyDCP(cfg, false, false)
}

func ApplyDCP(cfg *config.CloudConfig, restart, install bool) error {
	mode, err := mode.Get()
	if err != nil {
		return err
	}
	if mode == "install" {
		return nil
	}

	dcpExists := false
	dcpLocalExists := false
	if _, err := os.Stat("/sbin/dcp"); err == nil {
		dcpExists = true
	}
	if _, err := os.Stat("/usr/local/bin/dcp"); err == nil {
		dcpLocalExists = true
	}

	args := cfg.BhojpurOS.DcpArgs
	vars := []string{
		"INSTALL_DCP_NAME=service",
	}

	if !dcpExists && !restart {
		return nil
	}

	if dcpExists {
		vars = append(vars, "INSTALL_DCP_SKIP_DOWNLOAD=true")
		vars = append(vars, "INSTALL_DCP_BIN_DIR=/sbin")
		vars = append(vars, "INSTALL_DCP_BIN_DIR_READ_ONLY=true")
	} else if dcpLocalExists {
		vars = append(vars, "INSTALL_DCP_SKIP_DOWNLOAD=true")
	} else if !install {
		return nil
	}

	if !restart {
		vars = append(vars, "INSTALL_DCP_SKIP_START=true")
	}

	if cfg.BhojpurOS.ServerURL == "" {
		if len(args) == 0 {
			args = append(args, "server")
		}
	} else {
		vars = append(vars, fmt.Sprintf("DCP_URL=%s", cfg.BhojpurOS.ServerURL))
		if len(args) == 0 {
			args = append(args, "agent")
		}
	}

	if strings.HasPrefix(cfg.BhojpurOS.Token, "K10") {
		vars = append(vars, fmt.Sprintf("DCP_TOKEN=%s", cfg.BhojpurOS.Token))
	} else if cfg.BhojpurOS.Token != "" {
		vars = append(vars, fmt.Sprintf("DCP_CLUSTER_SECRET=%s", cfg.BhojpurOS.Token))
	}

	var labels []string
	for k, v := range cfg.BhojpurOS.Labels {
		labels = append(labels, fmt.Sprintf("%s=%s", k, v))
	}
	if mode != "" {
		labels = append(labels, fmt.Sprintf("os.bhojpur.net/mode=%s", mode))
	}
	labels = append(labels, fmt.Sprintf("os.bhojpur.net/version=%s", version.Version))
	sort.Strings(labels)

	for _, l := range labels {
		args = append(args, "--node-label", l)
	}

	for _, taint := range cfg.BhojpurOS.Taints {
		args = append(args, "--kubelet-arg", "register-with-taints="+taint)
	}

	cmd := exec.Command("/usr/libexec/os/dcp-install.sh", args...)
	cmd.Env = append(os.Environ(), vars...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	logrus.Debugf("Running %s %v %v", cmd.Path, cmd.Args, vars)

	return cmd.Run()
}

func ApplyInstall(cfg *config.CloudConfig) error {
	mode, err := mode.Get()
	if err != nil {
		return err
	}
	if mode != "install" {
		return nil
	}

	cmd := exec.Command("opsutl", "install")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func ApplyDNS(cfg *config.CloudConfig) error {
	buf := &bytes.Buffer{}
	buf.WriteString("[General]\n")
	buf.WriteString("NetworkInterfaceBlacklist=veth\n")
	buf.WriteString("PreferredTechnologies=ethernet,wifi\n")
	if len(cfg.BhojpurOS.DNSNameservers) > 0 {
		dns := strings.Join(cfg.BhojpurOS.DNSNameservers, ",")
		buf.WriteString("FallbackNameservers=")
		buf.WriteString(dns)
		buf.WriteString("\n")
	} else {
		buf.WriteString("FallbackNameservers=8.8.8.8\n")
	}

	if len(cfg.BhojpurOS.NTPServers) > 0 {
		ntp := strings.Join(cfg.BhojpurOS.NTPServers, ",")
		buf.WriteString("FallbackTimeservers=")
		buf.WriteString(ntp)
		buf.WriteString("\n")
	}

	err := ioutil.WriteFile("/etc/connman/main.conf", buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write /etc/connman/main.conf: %v", err)
	}

	return nil
}

func ApplyWifi(cfg *config.CloudConfig) error {
	if len(cfg.BhojpurOS.Wifi) == 0 {
		return nil
	}

	buf := &bytes.Buffer{}

	buf.WriteString("[WiFi]\n")
	buf.WriteString("Enable=true\n")
	buf.WriteString("Tethering=false\n")

	if buf.Len() > 0 {
		if err := os.MkdirAll("/var/lib/connman", 0755); err != nil {
			return fmt.Errorf("failed to mkdir /var/lib/connman: %v", err)
		}
		if err := ioutil.WriteFile("/var/lib/connman/settings", buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write to /var/lib/connman/settings: %v", err)
		}
	}

	buf = &bytes.Buffer{}

	buf.WriteString("[global]\n")
	buf.WriteString("Name=cloud-config\n")
	buf.WriteString("Description=Services defined in the cloud-config\n")

	for i, w := range cfg.BhojpurOS.Wifi {
		name := fmt.Sprintf("wifi%d", i)
		buf.WriteString("[service_")
		buf.WriteString(name)
		buf.WriteString("]\n")
		buf.WriteString("Type=wifi\n")
		buf.WriteString("Passphrase=")
		buf.WriteString(w.Passphrase)
		buf.WriteString("\n")
		buf.WriteString("Name=")
		buf.WriteString(w.Name)
		buf.WriteString("\n")
	}

	if buf.Len() > 0 {
		return ioutil.WriteFile("/var/lib/connman/cloud-config.config", buf.Bytes(), 0644)
	}

	return nil
}

func ApplyDataSource(cfg *config.CloudConfig) error {
	if len(cfg.BhojpurOS.DataSources) == 0 {
		return nil
	}

	args := strings.Join(cfg.BhojpurOS.DataSources, " ")
	buf := &bytes.Buffer{}

	buf.WriteString("command_args=\"")
	buf.WriteString(args)
	buf.WriteString("\"\n")

	if err := ioutil.WriteFile("/etc/conf.d/cloud-config", buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write to /etc/conf.d/cloud-config: %v", err)
	}

	return nil
}

func ApplyEnvironment(cfg *config.CloudConfig) error {
	if len(cfg.BhojpurOS.Environment) == 0 {
		return nil
	}
	env := make(map[string]string, len(cfg.BhojpurOS.Environment))
	if buf, err := ioutil.ReadFile("/etc/environment"); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(buf))
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "#") {
				continue
			}
			line = strings.TrimPrefix(line, "export")
			line = strings.TrimSpace(line)
			if len(line) > 1 {
				parts := strings.SplitN(line, "=", 2)
				key := parts[0]
				val := ""
				if len(parts) > 1 {
					if val, err = strconv.Unquote(parts[1]); err != nil {
						val = parts[1]
					}
				}
				env[key] = val
			}
		}
	}
	for key, val := range cfg.BhojpurOS.Environment {
		env[key] = val
	}
	buf := &bytes.Buffer{}
	for key, val := range env {
		buf.WriteString(key)
		buf.WriteString("=")
		buf.WriteString(strconv.Quote(val))
		buf.WriteString("\n")
	}
	if err := ioutil.WriteFile("/etc/environment", buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write to /etc/environment: %v", err)
	}

	return nil
}
