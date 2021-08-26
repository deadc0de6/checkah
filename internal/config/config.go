// Copyright (c) 2021 deadc0de6

package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

// print a config file to stdout
func PrintConfig(cfg *Config, format string) error {
	// serialize config
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	// read struct to viper
	v := viper.New()
	v.SetConfigType(format)
	err = v.ReadConfig(bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	// viper doesn't allow to write to string
	file, err := ioutil.TempFile("", "checkah")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	// write config
	err = v.WriteConfigAs(file.Name())
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return err
	}

	fmt.Print(string(content))
	return nil
}

func getEmptyConfig() Config {
	// set defaults
	content := Config{
		Settings: Settings{
			HostsParallel:  false,
			ChecksParallel: true,
		},
		Hosts:    []Host{},
		Profiles: []Profile{},
	}
	return content
}

func duplicates(left []string, right []string) error {
	for _, l := range left {
		for _, r := range right {
			if l == r {
				return fmt.Errorf("duplicate entry name: %s", l)
			}
		}
	}
	return nil
}

// MergeConfigs merge two configs
func MergeConfigs(left *Config, right *Config) (*Config, error) {
	n := &Config{}

	// merge settings
	// last takes precedence
	n.Settings = right.Settings

	// merge hosts
	var leftNames []string
	var rightNames []string
	for _, h := range left.Hosts {
		leftNames = append(leftNames, h.Name)
	}
	for _, h := range right.Hosts {
		rightNames = append(rightNames, h.Name)
	}
	err := duplicates(leftNames, rightNames)
	if err != nil {
		return nil, err
	}
	n.Hosts = append(n.Hosts, left.Hosts...)
	n.Hosts = append(n.Hosts, right.Hosts...)

	// merge profiles
	leftNames = []string{}
	rightNames = []string{}
	for _, h := range left.Profiles {
		leftNames = append(leftNames, h.Name)
	}
	for _, h := range right.Profiles {
		rightNames = append(rightNames, h.Name)
	}
	err = duplicates(leftNames, rightNames)
	if err != nil {
		return nil, err
	}
	n.Profiles = append(n.Profiles, left.Profiles...)
	n.Profiles = append(n.Profiles, right.Profiles...)

	return n, nil
}

// ReadCfg reads config
func ReadCfg(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := getEmptyConfig()
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
