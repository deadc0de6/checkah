// Copyright (c) 2026 deadc0de6

package check

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/deadc0de6/checkah/internal/transport"
)

// Zfs the zfs struct
type Zfs struct {
	command  string
	poolName string
	limit    int
	options  map[string]string
}

func (c *Zfs) getInt(str string) (int, error) {
	if len(str) < 1 {
		return -1, fmt.Errorf("empty percent value")
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		return -1, err
	}
	return val, nil
}

// Run executes the check
func (c *Zfs) Run(t transport.Transport) *Result {
	stdout, _, err := t.Execute(c.command)
	if err != nil {
		return c.returnCheck("", err)
	}

	// NAME     USED  AVAIL  MOUNTPOINT
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		fields := strings.Split(line, " ")

		if fields[0] != c.poolName {
			// bad pool
			continue
		}
		if len(fields) != 4 {
			// bad number of fields
			continue
		}

		used, err := c.getInt(fields[1])
		if err != nil {
			return c.returnCheck("", err)
		}
		avail, err := c.getInt(fields[2])
		if err != nil {
			return c.returnCheck("", err)
		}
		total := used + avail
		percent := used * 100 / total
		percentStr := fmt.Sprintf("%d%%", percent)

		if percent > c.limit {
			err := fmt.Errorf("zfs pool \"%s\" used above %d%%: %s", c.poolName, c.limit, percentStr)
			return c.returnCheck(percentStr, err)
		}

		return c.returnCheck(percentStr, nil)
	}

	// mount point not found
	return c.returnCheck("", fmt.Errorf("zfs pool \"%s\" not found", c.poolName))
}

func (c *Zfs) returnCheck(value string, err error) *Result {
	return &Result{
		Name:        c.GetName(),
		Description: c.GetDescription(),
		Value:       value,
		Limit:       fmt.Sprintf("%d", c.limit),
		Error:       err,
	}
}

// GetName returns the check name
func (c *Zfs) GetName() string {
	return "zfs"
}

// GetDescription returns description
func (c *Zfs) GetDescription() string {
	return fmt.Sprintf("zfs \"%s\" used", c.poolName)
}

// GetOptions returns the options
func (c *Zfs) GetOptions() map[string]string {
	return c.options
}

// NewCheckZfs creates a zfs check instance
func NewCheckZfs(options map[string]string) (*Zfs, error) {
	mount := "tank"
	v, ok := options["pool"]
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

	command := fmt.Sprintf("zfs list -p -H -o name,used,available,mountpoint %s", mount)
	c := Zfs{
		command:  command,
		poolName: mount,
		limit:    limit,
		options:  options,
	}

	return &c, nil
}
