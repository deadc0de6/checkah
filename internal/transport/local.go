// Copyright (c) 2021 deadc0de6

package transport

import (
	"bytes"
	"os/exec"
)

// Local the localhost fake object
type Local struct{}

// Execute executes a command through SSH
// returns stdout, stderr, error
func (l *Local) Execute(cmd string) (string, string, error) {
	var stdout, stderr bytes.Buffer

	c := exec.Command("bash", "-c", cmd)
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()
	if err != nil {
		return "", "", err
	}

	return stdout.String(), stderr.String(), nil
}

// Copy fakes copy
func (l *Local) Copy(string, string, string) error {
	return nil
}

// Mkdir fakes mkdir
func (l *Local) Mkdir(string) error {
	return nil
}

// Close fakes closing connection
func (l *Local) Close() {}

// NewLocal creates a new local object
func NewLocal() (*Local, error) {
	l := &Local{}
	return l, nil
}
