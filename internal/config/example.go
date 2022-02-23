package config

// PrintExampleConfig prints config example
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
		{
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
		{
			Type:    "disk",
			Disable: false,
			Options: map[string]string{
				"mount": "/",
				"limit": "80",
			},
		},
		{
			Type:    "loadavg",
			Disable: false,
			Options: map[string]string{
				"load_15min": "1",
			},
		},
		{
			Type:    "process",
			Disable: false,
			Options: map[string]string{
				"pattern": "sshd",
			},
		},
		{
			Type:    "memory",
			Disable: false,
			Options: map[string]string{
				"limit_mem": "90",
			},
		},
		{
			Type:    "tcp",
			Disable: false,
			Options: map[string]string{
				"port": "22",
			},
		},
		{
			Type:    "uptime",
			Disable: false,
			Options: map[string]string{
				"days": "180",
			},
		},
	}

	// create alerts
	profile1Alerts := []Alert{
		{
			Type:    "file",
			Disable: false,
			Options: map[string]string{
				"path": "/tmp/alerts.txt",
			},
		},
		{
			Type:    "command",
			Disable: false,
			Options: map[string]string{
				"command": "notify-send -u critical",
			},
		},
	}

	// create the profiles block
	profiles := []Profile{
		{
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
		{
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
		{
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
		{
			Type:    "disk",
			Disable: false,
			Options: map[string]string{
				"mount": "/",
				"limit": "80",
			},
		},
		{
			Type:    "disk",
			Disable: false,
			Options: map[string]string{
				"mount": "/boot",
				"limit": "80",
			},
		},
		{
			Type:    "loadavg",
			Disable: false,
			Options: map[string]string{
				"load_15min": "1",
			},
		},
		{
			Type:    "process",
			Disable: false,
			Options: map[string]string{
				"pattern": "sshd",
			},
		},
		{
			Type:    "process",
			Disable: false,
			Options: map[string]string{
				"pattern": "firefox",
				"invert":  "yes",
			},
		},
		{
			Type:    "memory",
			Disable: false,
			Options: map[string]string{
				"limit_mem": "90",
			},
		},
		{
			Type:    "tcp",
			Disable: false,
			Options: map[string]string{
				"port": "22",
			},
		},
		{
			Type:    "uptime",
			Disable: false,
			Options: map[string]string{
				"days": "180",
			},
		},
	}

	// create profile1 alerts
	profile1Alerts := []Alert{
		{
			Type:    "file",
			Disable: false,
			Options: map[string]string{
				"path": "/tmp/alerts.txt",
			},
		},
		{
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
		{
			Type:    "command",
			Disable: false,
			Options: map[string]string{
				"command": "notify-send -u critical",
			},
		},
		{
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
		{
			Type:    "tcp",
			Disable: false,
			Options: map[string]string{
				"port": "443",
			},
		},
	}

	// create the profiles block
	profiles := []Profile{
		{
			Name:   "profile1",
			Checks: profile1Checks,
			Alerts: profile1Alerts,
		},
		{
			Name:   "profile2",
			Checks: profile2Checks,
			Extend: []string{"profile1"},
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
