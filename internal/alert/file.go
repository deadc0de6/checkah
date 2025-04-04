// Copyright (c) 2021 deadc0de6

package alert

import (
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// File alert file struct
type File struct {
	path     string
	truncate bool
	options  map[string]string
}

// Notify notifies
func (a *File) Notify(content string) error {
	f, err := os.OpenFile(a.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// now
	t := time.Now()
	now := t.Format("2006-01-02 15:04:05")

	// append
	line := fmt.Sprintf("[%s] %s\n", now, content)
	_, err = f.WriteString(line)
	if err != nil {
		return err
	}
	return nil
}

// GetOptions returns this alert options
func (a *File) GetOptions() map[string]string {
	return a.options
}

// GetDescription returns a description for this alert
func (a *File) GetDescription() string {
	return fmt.Sprintf("alert to file %s", a.path)
}

// NewAlertFile creates a new file alert instance
func NewAlertFile(options map[string]string) (*File, error) {
	path, ok := options["path"]
	if !ok {
		return nil, fmt.Errorf("\"path\" option required")
	}

	truncate := false
	trunc, ok := options["truncate"]
	if ok {
		log.Debugf("truncate value: %s", trunc)
		truncate = trunc == "1" || strings.ToLower(trunc) == "true"
	}

	if truncate {
		// truncate file
		err := os.Truncate(path, 0)
		if err != nil {
			log.Errorf("%v", err)
		}
		log.Debugf("truncate %s", path)
	}

	a := &File{
		path:     path,
		truncate: truncate,
		options:  options,
	}
	return a, nil
}
