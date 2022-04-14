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
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/bhojpur/os/pkg/config"
	"github.com/bhojpur/os/pkg/mode"
	"github.com/bhojpur/os/pkg/questions"
	"github.com/bhojpur/os/pkg/util"
)

func Ask(cfg *config.CloudConfig) (bool, error) {
	if ok, err := isInstall(cfg); err != nil {
		return false, err
	} else if ok {
		return true, AskInstall(cfg)
	}

	return false, AskServerAgent(cfg)
}

func isInstall(cfg *config.CloudConfig) (bool, error) {
	mode, err := mode.Get()
	if err != nil {
		return false, err
	}

	if mode == "install" {
		return true, nil
	} else if mode == "live-server" {
		return false, nil
	} else if mode == "live-agent" {
		return false, nil
	}

	i, err := questions.PromptFormattedOptions("Choose operation", 0,
		"Install to disk",
		"Configure server or agent")
	if err != nil {
		return false, err
	}

	return i == 0, nil
}

func AskInstall(cfg *config.CloudConfig) error {
	if cfg.BhojpurOS.Install.Silent {
		return nil
	}

	if err := AskInstallDevice(cfg); err != nil {
		return err
	}

	if err := AskConfigURL(cfg); err != nil {
		return err
	}

	if cfg.BhojpurOS.Install.ConfigURL == "" {
		if err := AskGithub(cfg); err != nil {
			return err
		}

		if err := AskPassword(cfg); err != nil {
			return err
		}

		if err := AskWifi(cfg); err != nil {
			return err
		}

		if err := AskServerAgent(cfg); err != nil {
			return err
		}
	}

	return nil
}

func AskInstallDevice(cfg *config.CloudConfig) error {
	if cfg.BhojpurOS.Install.Device != "" {
		return nil
	}

	output, err := exec.Command("/bin/sh", "-c", "lsblk -r -o NAME,TYPE | grep -w disk | awk '{print $1}'").CombinedOutput()
	if err != nil {
		return err
	}
	fields := strings.Fields(string(output))
	i, err := questions.PromptFormattedOptions("Installation target. Device will be formatted", -1, fields...)
	if err != nil {
		return err
	}

	cfg.BhojpurOS.Install.Device = "/dev/" + fields[i]
	return nil
}

func AskToken(cfg *config.CloudConfig, server bool) error {
	var (
		token string
		err   error
	)

	if cfg.BhojpurOS.Token != "" {
		return nil
	}

	msg := "Token or cluster secret"
	if server {
		msg += " (optional)"
	}
	if server {
		token, err = questions.PromptOptional(msg+": ", "")
	} else {
		token, err = questions.Prompt(msg+": ", "")
	}
	cfg.BhojpurOS.Token = token

	return err
}

func isServer(cfg *config.CloudConfig) (bool, error) {
	mode, err := mode.Get()
	if err != nil {
		return false, err
	}
	if mode == "live-server" {
		return true, nil
	} else if mode == "live-agent" || (cfg.BhojpurOS.ServerURL != "" && cfg.BhojpurOS.Token != "") {
		return false, nil
	}

	opts := []string{"server", "agent"}
	i, err := questions.PromptFormattedOptions("Run as server or agent?", 0, opts...)
	if err != nil {
		return false, err
	}

	return i == 0, nil
}

func AskServerAgent(cfg *config.CloudConfig) error {
	if cfg.BhojpurOS.ServerURL != "" {
		return nil
	}

	server, err := isServer(cfg)
	if err != nil {
		return err
	}

	if server {
		return AskToken(cfg, true)
	}

	url, err := questions.Prompt("URL of server: ", "")
	if err != nil {
		return err
	}
	cfg.BhojpurOS.ServerURL = url

	return AskToken(cfg, false)
}

func AskPassword(cfg *config.CloudConfig) error {
	if len(cfg.SSHAuthorizedKeys) > 0 || cfg.BhojpurOS.Password != "" {
		return nil
	}

	var (
		ok   = false
		err  error
		pass string
	)

	for !ok {
		pass, ok, err = util.PromptPassword()
		if err != nil {
			return err
		}
	}

	if os.Getuid() != 0 {
		cfg.BhojpurOS.Password = pass
		return nil
	}

	oldShadow, err := ioutil.ReadFile("/etc/shadow")
	if err != nil {
		return err
	}
	defer func() {
		ioutil.WriteFile("/etc/shadow", oldShadow, 0640)
	}()

	cmd := exec.Command("chpasswd")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("bhojpur:%s", pass))
	errBuffer := &bytes.Buffer{}
	cmd.Stdout = os.Stdout
	cmd.Stderr = errBuffer

	if err := cmd.Run(); err != nil {
		os.Stderr.Write(errBuffer.Bytes())
		return err
	}

	f, err := os.Open("/etc/shadow")
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ":")
		if len(fields) > 1 && fields[0] == "bhojpur" {
			cfg.BhojpurOS.Password = fields[1]
			return nil
		}
	}

	return scanner.Err()
}

func AskWifi(cfg *config.CloudConfig) error {
	if len(cfg.BhojpurOS.Wifi) > 0 {
		return nil
	}

	ok, err := questions.PromptBool("Configure WiFi?", false)
	if !ok || err != nil {
		return err
	}

	for {
		name, err := questions.Prompt("WiFi Name: ", "")
		if err != nil {
			return err
		}

		pass, err := questions.Prompt("WiFi Passphrase: ", "")
		if err != nil {
			return err
		}

		cfg.BhojpurOS.Wifi = append(cfg.BhojpurOS.Wifi, config.Wifi{
			Name:       name,
			Passphrase: pass,
		})

		ok, err := questions.PromptBool("Configure another WiFi network?", false)
		if !ok || err != nil {
			return err
		}
	}
}

func AskGithub(cfg *config.CloudConfig) error {
	if len(cfg.SSHAuthorizedKeys) > 0 || cfg.BhojpurOS.Password != "" {
		return nil
	}

	ok, err := questions.PromptBool("Authorize GitHub users to SSH?", false)
	if !ok || err != nil {
		return err
	}

	str, err := questions.Prompt("Comma separated list of GitHub users to authorize: ", "")
	if err != nil {
		return err
	}

	for _, s := range strings.Split(str, ",") {
		cfg.SSHAuthorizedKeys = append(cfg.SSHAuthorizedKeys, "github:"+strings.TrimSpace(s))
	}

	return nil
}

func AskConfigURL(cfg *config.CloudConfig) error {
	if cfg.BhojpurOS.Install.ConfigURL != "" {
		return nil
	}

	ok, err := questions.PromptBool("Config system with cloud-init file?", false)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	str, err := questions.Prompt("cloud-init file location (file path or http URL): ", "")
	if err != nil {
		return err
	}

	cfg.BhojpurOS.Install.ConfigURL = str
	return nil
}
