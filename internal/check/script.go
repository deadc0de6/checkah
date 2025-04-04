// Copyright (c) 2021 deadc0de6

package check

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/caarlos0/log"
	"github.com/deadc0de6/checkah/internal/transport"
)

// Script the Script struct
type Script struct {
	path    string
	retCode int
	options map[string]string
}

const (
	rmScript     = "rm -f %s"
	pathOnRemote = "/tmp/checkah.check"
)

func (c *Script) returnCheck(value string, err error) *Result {
	limit := ""
	return &Result{
		Name:        c.GetName(),
		Description: c.GetDescription(),
		Value:       value,
		Limit:       limit,
		Error:       err,
	}
}

// Run executes the check
func (c *Script) Run(t transport.Transport) *Result {
	remotePath := pathOnRemote
	remoteDir := path.Dir(remotePath)

	err := t.Mkdir(remoteDir)
	if err != nil {
		return c.returnCheck("", fmt.Errorf("scp create \"%s\" failed: %v", remoteDir, err))
	}

	// copy the file over
	err = t.Copy(c.path, remotePath, "755")
	if err != nil {
		return c.returnCheck("", fmt.Errorf("scp \"%s\" to \"%s\" failed: %v", c.path, remotePath, err))
	}
	cmd := fmt.Sprintf(rmScript, remotePath)
	defer func() {
		_, _, err := t.Execute(cmd)
		if err != nil {
			log.Errorf("%v", err)
		}
	}()

	// execute script
	sout, serr, err := t.Execute(remotePath)
	if err != nil {
		// return stderr and error
		return c.returnCheck(serr, fmt.Errorf("remote script \"%s\" failed: %v", remotePath, err))
	}

	sout = strings.TrimSuffix(sout, "\n")
	return c.returnCheck(fmt.Sprintf("custom script \"%s\" was successful: %s", c.path, sout), nil)
}

// GetName returns the check name
func (c *Script) GetName() string {
	return fmt.Sprintf("script \"%s\"", c.path)
}

// GetDescription get description
func (c *Script) GetDescription() string {
	return fmt.Sprintf("custom script \"%s\"", c.path)
}

// GetOptions returns the options
func (c *Script) GetOptions() map[string]string {
	return c.options
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// NewCheckScript creates a script check instance
func NewCheckScript(options map[string]string) (*Script, error) {
	path, ok := options["path"]
	if !ok {
		return nil, fmt.Errorf("\"path\" value mandatory")
	}

	if !fileExists(path) {
		return nil, fmt.Errorf("%s does not exist", path)
	}

	c := Script{
		path:    path,
		retCode: -1,
		options: options,
	}

	return &c, nil
}
