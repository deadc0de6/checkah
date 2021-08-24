// Copyright (c) 2021 deadc0de6

package config

// Config config file content
type Config struct {
	Settings Settings  `mapstruture:"settings"`
	Hosts    []Host    `mapstructure:"hosts"`
	Profiles []Profile `mapstructure:"profiles"`
}

// Settings the settings
type Settings struct {
	HostsParallel  bool  `mapstructure:"hosts-parallel"`
	ChecksParallel bool  `mapstructure:"checks-parallel"`
	GlobalAlert    Alert `mapstructure:"global-alert"`
}

// Host host block content
type Host struct {
	Name              string   `mapstructure:"name"`
	Host              string   `mapstructure:"host"`
	Port              string   `mapstructure:"port"`
	User              string   `mapstructure:"user"`
	Password          string   `mapstructure:"password"`
	Keyfile           string   `mapstructure:"keyfile"`
	ProfileNames      []string `mapstructure:"profiles"`
	KnownHostInsecure bool     `mapstructure:"insecure"`
	Disable           bool     `mapstructure:"disable"`
	Timeout           string   `mapstructure:"timeout"`
}

// Profile profile block content
type Profile struct {
	Name   string  `mapstructure:"name"`
	Checks []Check `mapstructure:"checks"`
	Alerts []Alert `mapstructure:"alerts"`
	Extend string  `mapstructure:"extend"`
}

// Check profile check block content
type Check struct {
	Type    string            `mapstructure:"type"`
	Options map[string]string `mapstructure:"options"`
	Disable bool              `mapstructure:"disable"`
}

// Alert profile alert block content
type Alert struct {
	Type    string            `mapstructure:"type"`
	Options map[string]string `mapstructure:"options"`
	Disable bool              `mapstructure:"disable"`
}
