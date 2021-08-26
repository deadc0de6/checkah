// Copyright (c) 2021 deadc0de6

package config

// Config config file content
type Config struct {
	Settings Settings  `mapstruture:"settings" json:"settings"`
	Hosts    []Host    `mapstructure:"hosts" json:"hosts"`
	Profiles []Profile `mapstructure:"profiles" json:"profiles"`
}

// Settings the settings
type Settings struct {
	HostsParallel  bool  `mapstructure:"hosts-parallel" json:"hosts-parallel"`
	ChecksParallel bool  `mapstructure:"checks-parallel" json:"checks-parallel"`
	GlobalAlert    Alert `mapstructure:"global-alert" json:"global-alert"`
}

// Host host block content
type Host struct {
	Name              string   `mapstructure:"name" json:"name"`
	Host              string   `mapstructure:"host" json:"host"`
	Port              string   `mapstructure:"port" json:"port"`
	User              string   `mapstructure:"user" json:"user"`
	Password          string   `mapstructure:"password" json:"password"`
	Keyfile           string   `mapstructure:"keyfile" json:"keyfile"`
	ProfileNames      []string `mapstructure:"profiles" json:"profiles"`
	KnownHostInsecure bool     `mapstructure:"insecure" json:"insecure"`
	Disable           bool     `mapstructure:"disable" json:"disable"`
	Timeout           string   `mapstructure:"timeout" json:"timeout"`
}

// Profile profile block content
type Profile struct {
	Name   string  `mapstructure:"name" json:"name"`
	Checks []Check `mapstructure:"checks" json:"checks"`
	Alerts []Alert `mapstructure:"alerts" json:"alerts"`
	Extend string  `mapstructure:"extend" json:"extend"`
}

// Check profile check block content
type Check struct {
	Type    string            `mapstructure:"type" json:"type"`
	Options map[string]string `mapstructure:"options" json:"options"`
	Disable bool              `mapstructure:"disable" json:"disable"`
}

// Alert profile alert block content
type Alert struct {
	Type    string            `mapstructure:"type" json:"type"`
	Options map[string]string `mapstructure:"options" json:"options"`
	Disable bool              `mapstructure:"disable" json:"disable"`
}
