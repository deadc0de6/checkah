// Copyright (c) 2021 deadc0de6

package check

import (
	"checkah/internal/transport"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Memory the memory struct
type Memory struct {
	command       string
	limitUseMem   int
	limitUseSwap  int
	limitUseTotal int
	options       map[string]string
}

func (c *Memory) returnCheck(value string, err error) *Result {
	limit := fmt.Sprintf("used mem:%d%% swap:%d%% total:%d%%", c.limitUseMem, c.limitUseSwap, c.limitUseTotal)
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

// Run executes the check
func (c *Memory) Run(t transport.Transport) *Result {
	stdout, _, err := t.Execute(c.command)
	if err != nil {
		return c.returnCheck("", err)
	}

	// mem
	memVal, err := getPercent(stdout, 1)
	if err != nil {
		return c.returnCheck("", err)
	}
	if c.limitUseMem > -1 && memVal > c.limitUseMem {
		return c.returnCheck("", fmt.Errorf("memory used is above %d%%: %d%%", c.limitUseMem, memVal))
	}

	// swap
	swapVal, err := getPercent(stdout, 2)
	if err != nil {
		return c.returnCheck("", err)
	}
	if c.limitUseSwap > -1 && swapVal > c.limitUseSwap {
		return c.returnCheck("", fmt.Errorf("swap usage is above %d%%: %d%%", c.limitUseSwap, swapVal))
	}

	// total
	totalVal, err := getPercent(stdout, 3)
	if err != nil {
		return c.returnCheck("", err)
	}
	if c.limitUseTotal > -1 && totalVal > c.limitUseTotal {
		return c.returnCheck("", fmt.Errorf("total memory usage is above %d%%: %d%%", c.limitUseTotal, totalVal))
	}

	return c.returnCheck(fmt.Sprintf("used mem:%d%% swap:%d%% total:%d%%", memVal, swapVal, totalVal), nil)
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

	limitSwap := -1
	v, ok = options["limit_swap"]
	if ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		limitSwap = i
	}

	limitTotal := -1
	v, ok = options["limit_total"]
	if ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		limitTotal = i
	}

	c := Memory{
		command:       "free -t",
		limitUseMem:   limitMem,
		limitUseSwap:  limitSwap,
		limitUseTotal: limitTotal,
		options:       options,
	}

	return &c, nil
}
