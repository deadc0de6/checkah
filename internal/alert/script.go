// Copyright (c) 2021 deadc0de6

package alert

import (
	"fmt"
	"os"
	"os/exec"
)

// Script alert file struct
type Script struct {
	path    string
	options map[string]string
}

// Notify notifies
func (a *Script) Notify(content string) error {
	cmd := exec.Command(a.path, content)
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
	return fmt.Sprintf("alert to script %s", a.path)
}

// NewAlertScript creates a new script alert instance
func NewAlertScript(options map[string]string) (*Script, error) {
	path, ok := options["path"]
	if !ok {
		return nil, fmt.Errorf("\"path\" option required")
	}

	if !fileExists(path) {
		return nil, fmt.Errorf("%s does not exist", path)
	}

	a := &Script{
		path:    path,
		options: options,
	}
	return a, nil
}
