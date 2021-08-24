// Copyright (c) 2021 deadc0de6

package check

import (
	"checkah/internal/transport"
	"fmt"
)

// Reachable the port struct
type Reachable struct {
	command string
}

func (c *Reachable) returnCheck(value string, err error) *Result {
	limit := "is reachable"
	return &Result{
		Name:        c.GetName(),
		Description: c.GetDescription(),
		Value:       value,
		Limit:       limit,
		Error:       err,
	}
}

// Run executes the check
func (c *Reachable) Run(t transport.Transport) *Result {
	_, _, err := t.Execute(c.command)
	if err != nil {
		return c.returnCheck("", fmt.Errorf("host is NOT reachable: %v", err))
	}
	return c.returnCheck("ok", nil)
}

// GetName returns the check name
func (c *Reachable) GetName() string {
	return "reachable"
}

// GetDescription get description
func (c *Reachable) GetDescription() string {
	return "host is reachable"
}

// GetOptions returns the options
func (c *Reachable) GetOptions() map[string]string {
	return nil
}

// NewCheckReachable creates a disk check instance
func NewCheckReachable(map[string]string) (*Reachable, error) {
	c := Reachable{
		command: "hostname",
	}

	return &c, nil
}
