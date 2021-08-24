// Copyright (c) 2021 deadc0de6

package alert

import (
	"fmt"
	"os/exec"
)

// Command alert file struct
type Command struct {
	command string
	options map[string]string
}

// Notify notifies
func (a *Command) Notify(content string) error {
	cmd := exec.Command(a.command, content)
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

	a := &Command{
		command: command,
		options: options,
	}
	return a, nil
}
