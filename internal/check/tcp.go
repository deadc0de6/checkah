// Copyright (c) 2021 deadc0de6

package check

import (
	"fmt"

	"github.com/deadc0de6/checkah/internal/transport"
)

// Port the port struct
type Port struct {
	port    string
	command string
	options map[string]string
}

func (c *Port) returnCheck(value string, err error) *Result {
	limit := "check port is open"
	return &Result{
		Name:        c.GetName(),
		Description: c.GetDescription(),
		Value:       value,
		Limit:       limit,
		Error:       err,
	}
}

// Run executes the check
func (c *Port) Run(t transport.Transport) *Result {
	_, _, err := t.Execute(c.command)
	if err != nil {
		return c.returnCheck("", fmt.Errorf("TCP port \"%s\" is NOT open", c.port))
	}
	return c.returnCheck(fmt.Sprintf("TCP port \"%s\" is open", c.port), nil)
}

// GetName returns the check name
func (c *Port) GetName() string {
	return "tcp"
}

// GetDescription get description
func (c *Port) GetDescription() string {
	return fmt.Sprintf("TCP port \"%s\"", c.port)
}

// GetOptions returns the options
func (c *Port) GetOptions() map[string]string {
	return c.options
}

// NewCheckPort creates a disk check instance
func NewCheckPort(options map[string]string) (*Port, error) {
	port, ok := options["port"]
	if !ok {
		return nil, fmt.Errorf("TCP port value mandatory")
	}

	c := Port{
		command: fmt.Sprintf("ss -tulpn | grep ':%s'", port),
		port:    port,
		options: options,
	}

	return &c, nil
}
