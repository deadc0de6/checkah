// Copyright (c) 2021 deadc0de6

package alert

import (
	"fmt"
	"os/exec"
	"strings"
)

// Command alert file struct
type Command struct {
	command string
	args    []string
	options map[string]string
}

// Notify notifies
func (a *Command) Notify(content string) error {
	args := append(a.args, fmt.Sprintf("'%s'", content))
	cmd := exec.Command(a.command, args...)
	err := cmd.Run()
	return err
}

// GetOptions returns this alert options
func (a *Command) GetOptions() map[string]string {
	return a.options
}

// GetDescription returns a description for this alert
func (a *Command) GetDescription() string {
	return fmt.Sprintf("alert to command \"%s\"", a.command)
}

// NewAlertCommand creates a new script alert instance
func NewAlertCommand(options map[string]string) (*Command, error) {
	command, ok := options["command"]
	if !ok {
		return nil, fmt.Errorf("\"command\" option required")
	}

	if len(command) < 1 {
		return nil, fmt.Errorf("\"command\" option required")
	}

	fields := strings.Split(command, " ")
	var args []string
	if len(fields) > 1 {
		args = fields[1:]
	}

	a := &Command{
		command: fields[0],
		args:    args,
		options: options,
	}
	return a, nil
}
