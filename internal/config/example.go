package config

func PrintExampleConfig(format string, local bool) error {
	// create the example config
	var cfg *Config
	if local {
		cfg = createLocalhostConfig()
	} else {
		cfg = createExampleConfig()
	}

	return PrintConfig(cfg, format)
}

func createLocalhostConfig() *Config {
	settings := Settings{
		HostsParallel:  false,
		ChecksParallel: true,
	}

	// create the config hosts block
	hosts := []Host{
		Host{
			Name:              "localhost",
			Host:              "127.0.0.1",
			Port:              "22",
			User:              "",
			Password:          "",
			Keyfile:           "",
			KnownHostInsecure: false,
			Disable:           false,
			Timeout:           "5",
			ProfileNames: []string{
				"default",
			},
		},
	}

	// create profile1 checks
	profile1Checks := []Check{
		Check{
			Type:    "disk",
			Disable: false,
			Options: map[string]string{
				"mount": "/",
				"limit": "80",
			},
		},
		Check{
			Type:    "loadavg",
			Disable: false,
			Options: map[string]string{
				"load_15min": "1",
			},
		},
		Check{
			Type:    "process",
			Disable: false,
			Options: map[string]string{
				"pattern": "sshd",
			},
		},
		Check{
			Type:    "memory",
			Disable: false,
			Options: map[string]string{
				"limit_mem":   "90",
				"limit_swap":  "90",
				"limit_total": "90",
			},
		},
		Check{
			Type:    "tcp",
			Disable: false,
			Options: map[string]string{
				"port": "22",
			},
		},
	}

	// create alerts
	profile1Alerts := []Alert{
		Alert{
			Type:    "file",
			Disable: false,
			Options: map[string]string{
				"path": "/tmp/alerts.txt",
			},
		},
		Alert{
			Type:    "command",
			Disable: false,
			Options: map[string]string{
				"command": "notify-send",
			},
		},
	}

	// create the profiles block
	profiles := []Profile{
		Profile{
			Name:   "default",
			Checks: profile1Checks,
			Alerts: profile1Alerts,
		},
	}

	// create the config
	config := Config{
		Settings: settings,
		Hosts:    hosts,
		Profiles: profiles,
	}

	return &config

}

func createExampleConfig() *Config {
	// create the config settings block
	settings := Settings{
		HostsParallel:  false,
		ChecksParallel: true,
		GlobalAlert: Alert{
			Type: "file",
			Options: map[string]string{
				"path": "/tmp/global-alerts.txt",
			},
			Disable: false,
		},
	}

	// create the config hosts block
	hosts := []Host{
		Host{
			Name:              "local",
			Host:              "127.0.0.1",
			Port:              "22",
			User:              "",
			Password:          "",
			Keyfile:           "",
			KnownHostInsecure: false,
			Disable:           false,
			Timeout:           "5",
			ProfileNames: []string{
				"profile1",
			},
		},
		Host{
			Name:              "remote",
			Host:              "10.0.0.1",
			Port:              "22",
			User:              "",
			Password:          "",
			Keyfile:           "",
			KnownHostInsecure: false,
			Disable:           false,
			Timeout:           "5",
			ProfileNames: []string{
				"profile2",
			},
		},
	}

	// create profile1 checks
	profile1Checks := []Check{
		Check{
			Type:    "disk",
			Disable: false,
			Options: map[string]string{
				"mount": "/",
				"limit": "80",
			},
		},
		Check{
			Type:    "disk",
			Disable: false,
			Options: map[string]string{
				"mount": "/boot",
				"limit": "80",
			},
		},
		Check{
			Type:    "loadavg",
			Disable: false,
			Options: map[string]string{
				"load_15min": "1",
			},
		},
		Check{
			Type:    "process",
			Disable: false,
			Options: map[string]string{
				"pattern": "sshd",
			},
		},
		Check{
			Type:    "process",
			Disable: false,
			Options: map[string]string{
				"pattern": "firefox",
				"invert":  "yes",
			},
		},
		Check{
			Type:    "memory",
			Disable: false,
			Options: map[string]string{
				"limit_mem":   "90",
				"limit_swap":  "90",
				"limit_total": "90",
			},
		},
		Check{
			Type:    "tcp",
			Disable: false,
			Options: map[string]string{
				"port": "22",
			},
		},
	}

	// create profile1 alerts
	profile1Alerts := []Alert{
		Alert{
			Type:    "file",
			Disable: false,
			Options: map[string]string{
				"path": "/tmp/alerts.txt",
			},
		},
		Alert{
			Type:    "webhook",
			Disable: false,
			Options: map[string]string{
				"url":     "http://127.0.0.1",
				"header0": "h1",
				"value0":  "val1",
				"header1": "h2",
				"value1":  "val2",
			},
		},
		Alert{
			Type:    "command",
			Disable: false,
			Options: map[string]string{
				"command": "notify-send",
			},
		},
		Alert{
			Type:    "email",
			Disable: false,
			Options: map[string]string{
				"host":     "mail.example.com",
				"port":     "25",
				"mailfrom": "foo@example.com",
				"mailto":   "bar@example.com",
				"user":     "username",
				"password": "password",
			},
		},
	}

	// create profile2 checks
	profile2Checks := []Check{
		Check{
			Type:    "tcp",
			Disable: false,
			Options: map[string]string{
				"port": "443",
			},
		},
	}

	// create the profiles block
	profiles := []Profile{
		Profile{
			Name:   "profile1",
			Checks: profile1Checks,
			Alerts: profile1Alerts,
		},
		Profile{
			Name:   "profile2",
			Checks: profile2Checks,
			Extend: "profile1",
		},
	}

	// create the config
	config := Config{
		Settings: settings,
		Hosts:    hosts,
		Profiles: profiles,
	}

	return &config
}
