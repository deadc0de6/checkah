// Copyright (c) 2021 deadc0de6

package check

import (
	"checkah/internal/transport"
	"fmt"
)

// Result check result struct
type Result struct {
	Name        string
	Description string
	Value       string
	Limit       string
	Error       error
}

// Check the check interface
type Check interface {
	GetName() string
	GetDescription() string
	Run(transport.Transport) *Result
	GetOptions() map[string]string
}

// GetCheck returns a check instance
func GetCheck(name string, options map[string]string) (Check, error) {
	switch name {
	case "reachable":
		return NewCheckReachable(options)
	case "disk":
		return NewCheckDisk(options)
	case "loadavg":
		return NewCheckLoadAvg(options)
	case "process":
		return NewCheckProcess(options)
	case "memory":
		return NewCheckMemory(options)
	case "tcp":
		return NewCheckPort(options)
	case "script":
		return NewCheckScript(options)
	case "command":
		return NewCheckCommand(options)
	case "uptime":
		return NewCheckUptime(options)
	}
	return nil, fmt.Errorf("no such check: %s", name)
}
