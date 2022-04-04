// Copyright (c) 2021 deadc0de6

package check

import (
	"fmt"

	"github.com/deadc0de6/checkah/internal/transport"
)

// Process the process struct
type Process struct {
	invert  bool
	command string
	pattern string
	options map[string]string
}

func (c *Process) returnCheck(value string, err error) *Result {
	limit := "alert if process not is running"
	if c.invert {
		limit = "alert if process is running"
	}
	return &Result{
		Name:        c.GetName(),
		Description: c.GetDescription(),
		Value:       value,
		Limit:       limit,
		Error:       err,
	}
}

// Run executes the check
func (c *Process) Run(t transport.Transport) *Result {
	var isRunning bool
	_, _, err := t.Execute(c.command)
	if err != nil {
		isRunning = false
	} else {
		isRunning = true
	}

	if isRunning {
		if c.invert {
			// alert if is running
			return c.returnCheck("", fmt.Errorf("process \"%s\" is running", c.pattern))
		}
		return c.returnCheck(fmt.Sprintf("process \"%s\" is running", c.pattern), nil)
	}

	if c.invert {
		return c.returnCheck(fmt.Sprintf("process \"%s\" is not running", c.pattern), nil)
	}
	// alert if is not running
	return c.returnCheck("", fmt.Errorf("process \"%s\" is not running", c.pattern))
}

// GetName returns the check name
func (c *Process) GetName() string {
	return "process"
}

// GetDescription get description
func (c *Process) GetDescription() string {
	return fmt.Sprintf("process \"%s\"", c.pattern)
}

// GetOptions returns the options
func (c *Process) GetOptions() map[string]string {
	return c.options
}

// NewCheckProcess creates a disk check instance
func NewCheckProcess(options map[string]string) (*Process, error) {
	pattern, ok := options["pattern"]
	if !ok {
		return nil, fmt.Errorf("pattern value mandatory")
	}

	invert := false
	inv, ok := options["invert"]
	if ok {
		if inv == "yes" {
			invert = true
		}
	}

	c := Process{
		command: fmt.Sprintf("pgrep -f %s", pattern),
		invert:  invert,
		pattern: pattern,
		options: options,
	}

	return &c, nil
}
