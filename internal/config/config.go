// Copyright (c) 2021 deadc0de6

package config

import (
	"checkah/internal/alert"
	"fmt"
	"github.com/spf13/viper"
)

// PrintSettings prints the settings
func PrintSettings(cfg *Config) {
	fmt.Println("settings:")
	fmt.Printf("  hosts parallel: %v\n", cfg.Settings.HostsParallel)
	fmt.Printf("  checks parallel: %v\n", cfg.Settings.ChecksParallel)

	a, err := alert.GetAlert(cfg.Settings.GlobalAlert.Type, cfg.Settings.GlobalAlert.Options)
	if err == nil {
		fmt.Println("global-alert:")
		fmt.Printf(" description: %s\n", a.GetDescription())
		for k, v := range a.GetOptions() {
			fmt.Printf("    - %s=%s\n", k, v)
		}
	}
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

// MergeConfigs merge two configs
func MergeConfigs(left *Config, right *Config) (*Config, error) {
	n := &Config{}
	// merge settings
	// last takes precedence
	n.Settings = right.Settings
	// merge hosts
	n.Hosts = append(n.Hosts, left.Hosts...)
	n.Hosts = append(n.Hosts, right.Hosts...)
	// merge profiles
	n.Profiles = append(n.Profiles, left.Profiles...)
	n.Profiles = append(n.Profiles, right.Profiles...)
	return n, nil
}

// ReadCfg reads config
func ReadCfg(path string) (*Config, error) {
	cfg := viper.New()

	cfg.SetConfigFile(path)
	err := cfg.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := getEmptyConfig()

	err = cfg.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
