// Copyright (c) 2026 deadc0de6

package check

import (
	"fmt"
	"strings"

	"github.com/deadc0de6/checkah/internal/transport"
)

// systemctl is-enabled <service>
// 	- enabled	   Service is enabled at boot
// 	- disabled	 Service is disabled at boot
// 	- static	   Service unit cannot be enabled directly (used by another unit)
// 	- masked	   Service is completely disabled / symlinked to /dev/null
// 	- indirect	 Unit is enabled indirectly by another unit
// 	- generated	 Service was generated dynamically
// 	- bad	       Something is wrong (rare)

// systemctl is-active <service>
//  - active	      Service is running
//  - inactive	    Service is not running
//  - failed	      Service terminated with an error
//  - activating	  Service is in the process of starting
//  - deactivating	Service is in the process of stopping
//  - reloading	    Service is reloading configuration
//  - unknown	      Service unit not found
//  - maintenance	  Systemd is in maintenance mode (rare)

// Systemd the systemd struct
type Systemd struct {
	serviceName    string
	enabledCommand string
	serviceEnabled string
	runningCommand string
	serviceRunning string
	options        map[string]string
}

func (c *Systemd) returnCheck(value string, err error) *Result {
	limit := fmt.Sprintf("alert if service \"%s\" is not enabled=%s not running=%s", c.serviceName, c.serviceEnabled, c.serviceRunning)
	return &Result{
		Name:        c.GetName(),
		Description: c.GetDescription(),
		Value:       value,
		Limit:       limit,
		Error:       err,
	}
}

// Run executes the check
func (c *Systemd) Run(t transport.Transport) *Result {
	// check service enabled
	enabled, _, err := t.Execute(c.enabledCommand)
	if err != nil {
		err2 := fmt.Errorf("no such service \"%s\": %v", c.serviceName, err)
		return c.returnCheck("", err2)
	}
	enabled = strings.TrimSpace(enabled)
	if enabled != c.serviceEnabled {
		err := fmt.Errorf("service \"%s\" state is not \"%s\" but \"%s\"", c.serviceName, c.serviceEnabled, enabled)
		return c.returnCheck("", err)
	}

	// check service running
	running, _, err := t.Execute(c.runningCommand)
	if err != nil {
		err2 := fmt.Errorf("no such service \"%s\": %v", c.serviceName, err)
		return c.returnCheck("", err2)
	}
	running = strings.TrimSpace(running)
	if running != c.serviceRunning {
		err := fmt.Errorf("service \"%s\" running state is not \"%s\" but \"%s\"", c.serviceName, c.serviceRunning, running)
		return c.returnCheck("", err)
	}

	return c.returnCheck("ok", nil)
}

// GetName returns the check name
func (c *Systemd) GetName() string {
	return "systemd"
}

// GetDescription returns description
func (c *Systemd) GetDescription() string {
	return fmt.Sprintf("%s enabled=%s running=%s", c.serviceName, c.serviceEnabled, c.serviceRunning)
}

// GetOptions returns the options
func (c *Systemd) GetOptions() map[string]string {
	return c.options
}

// NewCheckSystemd creates a systemd check instance
func NewCheckSystemd(options map[string]string) (*Systemd, error) {
	serviceName, ok := options["service"]
	if !ok {
		return nil, fmt.Errorf("service value mandatory")
	}

	// service state
	state := "enabled"
	v, ok := options["state"]
	if ok {
		state = v
	}
	enabledCommand := fmt.Sprintf("systemctl is-enabled %s", serviceName)

	// service is running
	running := "active"
	v, ok = options["running"]
	if ok {
		running = v
	}
	activeCommand := fmt.Sprintf("systemctl is-active %s", serviceName)

	c := Systemd{
		serviceName:    serviceName,
		enabledCommand: enabledCommand,
		serviceEnabled: state,
		runningCommand: activeCommand,
		serviceRunning: running,
		options:        options,
	}

	return &c, nil
}
