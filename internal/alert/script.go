// Copyright (c) 2021 deadc0de6

package alert

import (
	"fmt"
	"os"
	"os/exec"
)

// Script alert file struct
type Script struct {
	command string
	args    []string
	options map[string]string
}

// Notify notifies
func (a *Script) Notify(content string) error {
	args := append(a.args, fmt.Sprintf("'%s'", content))
	cmd := exec.Command(a.command, args...)
	err := cmd.Run()
	return err
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetOptions returns this alert options
func (a *Script) GetOptions() map[string]string {
	return a.options
}

// GetDescription returns a description for this alert
func (a *Script) GetDescription() string {
	return fmt.Sprintf("alert to script %s", a.command)
}

// NewAlertScript creates a new script alert instance
func NewAlertScript(options map[string]string) (*Script, error) {
	command, ok := options["path"]
	if !ok {
		return nil, fmt.Errorf("\"path\" option required")
	}

	if len(command) < 1 {
		return nil, fmt.Errorf("\"path\" option required")
	}

	fields := splitArgs(command)
	var args []string
	if len(fields) > 1 {
		args = fields[1:]
	}

	if !fileExists(fields[0]) {
		return nil, fmt.Errorf("script \"%s\" does not exist", fields[0])
	}

	a := &Script{
		command: fields[0],
		args:    args,
		options: options,
	}
	return a, nil
}
