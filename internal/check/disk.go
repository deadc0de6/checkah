// Copyright (c) 2021 deadc0de6

package check

import (
	"checkah/internal/transport"
	"fmt"
	"strconv"
	"strings"
)

// Disk the disk struct
type Disk struct {
	command    string
	mountPoint string
	limit      int
	options    map[string]string
}

func (c *Disk) getValue(percent string) (int, error) {
	if len(percent) < 1 {
		return -1, fmt.Errorf("empty percent value")
	}
	i := percent[:len(percent)-1]
	val, err := strconv.Atoi(i)
	if err != nil {
		return -1, err
	}
	return val, nil
}

// Run executes the check
func (c *Disk) Run(t transport.Transport) *Result {
	stdout, _, err := t.Execute(c.command)
	if err != nil {
		return c.returnCheck("", err)
	}

	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		fields := strings.Split(line, " ")
		if fields[len(fields)-1] == c.mountPoint {
			// find first percent value
			var value = ""
			for _, f := range fields {
				if strings.Contains(f, "%") {
					value = f
					break
				}
			}
			v, err := c.getValue(value)
			if err != nil {
				return c.returnCheck("", err)
			}

			// check with limit
			if v > c.limit {
				err := fmt.Errorf("disk used of mount point \"%s\" above %d%%: %s", c.mountPoint, c.limit, value)
				return c.returnCheck(value, err)
			}
			return c.returnCheck(value, nil)
		}
	}

	// mount point not found
	return c.returnCheck("", fmt.Errorf("mount point %s not found", c.mountPoint))
}

func (c *Disk) returnCheck(value string, err error) *Result {
	return &Result{
		Name:        c.GetName(),
		Description: c.GetDescription(),
		Value:       value,
		Limit:       fmt.Sprintf("%d", c.limit),
		Error:       err,
	}
}

// GetName returns the check name
func (c *Disk) GetName() string {
	return "disk"
}

// GetDescription returns description
func (c *Disk) GetDescription() string {
	return fmt.Sprintf("disk \"%s\" used", c.mountPoint)
}

// GetOptions returns the options
func (c *Disk) GetOptions() map[string]string {
	return c.options
}

// NewCheckDisk creates a disk check instance
func NewCheckDisk(options map[string]string) (*Disk, error) {
	mount := "/"
	v, ok := options["mount"]
	if ok {
		mount = v
	}

	limit := 90
	v, ok = options["limit"]
	if ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		limit = i
	}

	c := Disk{
		command:    "df",
		mountPoint: mount,
		limit:      limit,
		options:    options,
	}

	return &c, nil
}
