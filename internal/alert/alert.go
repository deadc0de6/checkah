// Copyright (c) 2021 deadc0de6

package alert

import (
	"fmt"
)

// Alert the alert interface
type Alert interface {
	Notify(string) error
	GetDescription() string
	GetOptions() map[string]string
}

// GetAlert returns an alert instance
func GetAlert(name string, options map[string]string) (Alert, error) {
	switch name {
	case "file":
		return NewAlertFile(options)
	case "script":
		return NewAlertScript(options)
	case "webhook":
		return NewAlertWebhook(options)
	case "command":
		return NewAlertCommand(options)
	case "email":
		return NewAlertEmail(options)
	}
	return nil, fmt.Errorf("no such alert: %s", name)
}
