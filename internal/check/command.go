// Copyright (c) 2021 deadc0de6

package check

import (
	"fmt"

	"github.com/deadc0de6/checkah/internal/transport"
)

// Command the struct
type Command struct {
	command string
	options map[string]string
	name    string
}

func (c *Command) returnCheck(value string, err error) *Result {
	return &Result{
		Name:        c.GetName(),
		Description: c.GetDescription(),
		Value:       value,
		Limit:       "command failed",
		Error:       err,
	}
}

// Run executes the check
func (c *Command) Run(t transport.Transport) *Result {
	_, _, err := t.Execute(c.command)
	ret := "fail"
	if err == nil {
		ret = "success"
	}
	return c.returnCheck(ret, err)
}

// GetName returns the check name
func (c *Command) GetName() string {
	return "command"
}

// GetDescription get description
func (c *Command) GetDescription() string {
	name := c.name
	if len(c.name) < 1 {
		name = c.command
	}
	return fmt.Sprintf("command \"%s\"", name)
}

// GetOptions returns the options
func (c *Command) GetOptions() map[string]string {
	return c.options
}

// NewCheckCommand creates a disk check instance
func NewCheckCommand(options map[string]string) (*Command, error) {
	cmd, ok := options["command"]
	if !ok {
		return nil, fmt.Errorf("command value mandatory")
	}

	name := options["name"]

	c := Command{
		command: cmd,
		options: options,
		name:    name,
	}

	return &c, nil
}
