/*
Copyright 2024 The OpenYurt Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

var (
	keepalivedTmpl = "/etc/keepalived/keepalived.tmpl"
	keepalivedPid  = "/var/run/keepalived.pid"
)

type Keepalived struct {
	Iface          string
	Priority       int
	VIPs           []string
	KeepalivedTmpl *template.Template
	Vrid           int
	cmd            *exec.Cmd
	started        bool
}

func (k *Keepalived) Start() error {
	args := []string{"--dont-fork", "--log-console", "--log-detail", "--release-vips"}

	k.cmd = exec.Command("keepalived", args...)
	k.cmd.Stdout = os.Stdout
	k.cmd.Stderr = os.Stderr

	k.cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	k.started = true

	if err := k.cmd.Run(); err != nil {
		return fmt.Errorf("start keepalived error, %v", err)
	}
	return nil
}

func (k *Keepalived) LoadTemplates() error {
	tmpl, err := template.ParseFiles(keepalivedTmpl)
	if err != nil {
		return err
	}
	k.KeepalivedTmpl = tmpl
	return nil
}

func (k *Keepalived) ReloadKeepalived() error {
	for !k.IsRunning() {
		time.Sleep(time.Second)
	}

	err := syscall.Kill(k.cmd.Process.Pid, syscall.SIGHUP)
	if err != nil {
		return fmt.Errorf("reload keepalived error, %v", err)
	}

	return nil
}

func (k *Keepalived) IsRunning() bool {
	if !k.started {
		return false
	}

	if _, err := os.Stat(keepalivedPid); os.IsNotExist(err) {
		return false
	}

	return true
}

func (k *Keepalived) WriteCfg(path string, vips map[string]struct{}) error {
	vip := make([]string, 0, len(vips))
	for ip := range vips {
		vip = append(vip, ip)
	}
	k.VIPs = vip
	conf := make(map[string]interface{})
	conf["iface"] = k.Iface
	conf["vips"] = k.VIPs
	conf["priority"] = k.Priority
	conf["vrid"] = k.Vrid
	conf["iface"] = k.Iface

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()
	err = k.KeepalivedTmpl.Execute(w, conf)
	if err != nil {
		return fmt.Errorf("unexpected error creating keepalived.conf: %v", err)
	}

	return nil
}
