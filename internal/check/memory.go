// Copyright (c) 2021 deadc0de6

package check

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/deadc0de6/checkah/internal/transport"
)

// Memory the memory struct
type Memory struct {
	command     string
	limitUseMem int
	checker     func(string) (int, error)
	options     map[string]string
}

func (c *Memory) returnCheck(value string, err error) *Result {
	limit := fmt.Sprintf("%d%%", c.limitUseMem)
	return &Result{
		Name:        c.GetName(),
		Description: c.GetDescription(),
		Value:       value,
		Limit:       limit,
		Error:       err,
	}
}

func getField(output string, lineNb int, fieldNb int) (int, error) {
	lines := strings.Split(output, "\n")
	if len(lines) < lineNb {
		return 0, fmt.Errorf("cannot parse memory use")
	}
	line := lines[lineNb]
	// remove spaces
	r := regexp.MustCompile(`\s+`)
	line = r.ReplaceAllString(line, " ")
	fields := strings.Split(line, " ")
	if len(fields) < fieldNb {
		return 0, fmt.Errorf("cannot parse memory use")
	}
	val := fields[fieldNb]
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("cannot parse memory use")
	}
	return i, nil
}

func getPercent(output string, lineOffset int) (int, error) {
	total, err := getField(output, lineOffset, 1)
	if err != nil {
		return 0, err
	}
	use, err := getField(output, lineOffset, 2)
	if err != nil {
		return 0, err
	}

	val := 0
	if total > 0 {
		val = use * 100 / total
	}

	return val, nil
}

func memoryFromMemoryPressure(stdout string) (int, error) {
	var val string
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "System-wide memory free percentage:") {
			fields := strings.Split(line, ": ")
			val = fields[1][:len(fields[1])-1]
			break
		}
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("cannot parse memory use")
	}
	return 100 - i, nil
}

func memoryFromFree(stdout string) (int, error) {
	// mem
	memVal, err := getPercent(stdout, 1)
	if err != nil {
		return 0, err
	}
	return memVal, nil
}

// Run executes the check
func (c *Memory) Run(t transport.Transport) *Result {
	stdout, _, err := t.Execute(c.command)
	if err != nil {
		return c.returnCheck("", err)
	}

	val, err := c.checker(stdout)
	if err != nil {
		return c.returnCheck("", err)
	}

	if val > c.limitUseMem {
		return c.returnCheck("", fmt.Errorf("memory used is above %d%%: %d%%", c.limitUseMem, val))
	}

	return c.returnCheck(fmt.Sprintf("%d%%", val), nil)
}

// GetName returns the check name
func (c *Memory) GetName() string {
	return "memory"
}

// GetDescription get description
func (c *Memory) GetDescription() string {
	return "memory used"
}

// GetOptions returns the options
func (c *Memory) GetOptions() map[string]string {
	return c.options
}

// NewCheckMemory creates a disk check instance
func NewCheckMemory(options map[string]string) (*Memory, error) {
	limitMem := -1
	v, ok := options["limit_mem"]
	if ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		limitMem = i
	}

	cmd := "free -t"
	checker := memoryFromFree
	if !cmdExist("free") {
		cmd = "memory_pressure"
		checker = memoryFromMemoryPressure
	}

	c := Memory{
		command:     cmd,
		limitUseMem: limitMem,
		checker:     checker,
		options:     options,
	}

	return &c, nil
}
