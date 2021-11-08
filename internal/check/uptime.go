// Copyright (c) 2021 deadc0de6

package check

import (
	"checkah/internal/transport"
	"fmt"
	"strconv"
	"strings"
)

// Uptime the uptime struct
type Uptime struct {
	command   string
	limitDays int
	options   map[string]string
}

func (c *Uptime) returnCheck(value string, err error) *Result {
	limits := fmt.Sprintf("%d days", c.limitDays)
	return &Result{
		Name:        c.GetName(),
		Description: c.GetDescription(),
		Value:       value,
		Limit:       limits,
		Error:       err,
	}
}

// Run executes the check
func (c *Uptime) Run(t transport.Transport) *Result {
	stdout, _, err := t.Execute(c.command)
	if err != nil {
		return c.returnCheck("", err)
	}

	lines := strings.Split(stdout, " ")
	if len(lines) < 1 {
		return c.returnCheck("", fmt.Errorf("getting uptime failed"))
	}
	seconds := lines[0]

	val, err := strconv.ParseFloat(seconds, 64)
	if err != nil {
		return c.returnCheck("", err)
	}

	nbDays := int(int(val) / 60 / 60 / 24)
	if nbDays > c.limitDays {
		return c.returnCheck("", fmt.Errorf("uptime above %d days: %d days", c.limitDays, nbDays))
	}

	return c.returnCheck(fmt.Sprintf("%d days", nbDays), nil)
}

// GetName returns the check name
func (c *Uptime) GetName() string {
	return "uptime"
}

// GetDescription get description
func (c *Uptime) GetDescription() string {
	return "uptime"
}

// GetOptions returns the options
func (c *Uptime) GetOptions() map[string]string {
	return c.options
}

// NewCheckUptime creates a disk check instance
func NewCheckUptime(options map[string]string) (*Uptime, error) {
	days := -1
	v, ok := options["days"]
	if ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		days = i
	}

	c := Uptime{
		command:   "cat /proc/uptime",
		limitDays: days,
		options:   options,
	}

	return &c, nil
}
