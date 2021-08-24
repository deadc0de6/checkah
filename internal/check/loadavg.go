// Copyright (c) 2021 deadc0de6

package check

import (
	"checkah/internal/transport"
	"fmt"
	"strconv"
	"strings"
)

// Loadavg the loadavg struct
type Loadavg struct {
	command         string
	limitOneMin     float64
	limitFiveMin    float64
	limitFifteenMin float64
	options         map[string]string
}

func (c *Loadavg) returnCheck(value string, err error) *Result {
	limits := fmt.Sprintf("%f %f %f", c.limitOneMin, c.limitFiveMin, c.limitFifteenMin)
	return &Result{
		Name:        c.GetName(),
		Description: c.GetDescription(),
		Value:       value,
		Limit:       limits,
		Error:       err,
	}
}

// Run executes the check
func (c *Loadavg) Run(t transport.Transport) *Result {
	stdout, _, err := t.Execute(c.command)
	if err != nil {
		return c.returnCheck("", err)
	}

	lines := strings.Split(stdout, "\n")
	if len(lines) < 1 {
		return c.returnCheck("", fmt.Errorf("getting loadavg failed"))
	}
	line := lines[0]
	fields := strings.Split(line, " ")
	if len(fields) < 3 {
		return c.returnCheck("", fmt.Errorf("getting loadavg failed"))
	}

	val1min, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return c.returnCheck("", err)
	}
	if c.limitOneMin > 0 && val1min > c.limitOneMin {
		return c.returnCheck("", fmt.Errorf("1 min load average above %.2f: %.2f", c.limitOneMin, val1min))
	}

	val5min, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return c.returnCheck("", err)
	}
	if c.limitFiveMin > 0 && val5min > c.limitFiveMin {
		return c.returnCheck("", fmt.Errorf("5 min load average above %.2f: %.2f", c.limitFiveMin, val5min))
	}

	val15min, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return c.returnCheck("", err)
	}
	if c.limitFifteenMin > 0 && val15min > c.limitFifteenMin {
		return c.returnCheck("", fmt.Errorf("15 min load average above %.2f: %.2f", c.limitFifteenMin, val15min))
	}

	return c.returnCheck(fmt.Sprintf("%.2f %.2f %.2f", val1min, val5min, val15min), nil)
}

// GetName returns the check name
func (c *Loadavg) GetName() string {
	return "loadavg"
}

// GetDescription get description
func (c *Loadavg) GetDescription() string {
	return "load average"
}

// GetOptions returns the options
func (c *Loadavg) GetOptions() map[string]string {
	return c.options
}

// NewCheckLoadAvg creates a disk check instance
func NewCheckLoadAvg(options map[string]string) (*Loadavg, error) {
	limitOne := -1.
	v, ok := options["load_1min"]
	if ok {
		i, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		limitOne = i
	}

	limitFive := -1.
	v, ok = options["load_5min"]
	if ok {
		i, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		limitFive = i
	}

	limitFifteen := -1.
	v, ok = options["load_15min"]
	if ok {
		i, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		limitFifteen = i
	}

	c := Loadavg{
		command:         "cat /proc/loadavg",
		limitOneMin:     limitOne,
		limitFiveMin:    limitFive,
		limitFifteenMin: limitFifteen,
		options:         options,
	}

	return &c, nil
}
