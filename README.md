# CHECKAH

[![Tests Status](https://github.com/deadc0de6/checkah/workflows/tests/badge.svg)](https://github.com/deadc0de6/checkah/actions)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)

[![Donate](https://img.shields.io/badge/donate-KoFi-blue.svg)](https://ko-fi.com/deadc0de6)

[checkah](https://github.com/deadc0de6/checkah) is an agentless SSH system monitoring and alerting tool.

Features:

* agentless
* check over SSH (password, keyfile, agent)
* config file based (yaml, json)
* multiple alerts (webhooks, email, script, file, ...)
* multiple checks (disk, memory, loadavg, process, opened ports, ...)

You need golang 1.16

Quick start:
```bash
git clone https://github.com/deadc0de6/checkah
cd checkah
make
./bin/checkah example --format=yaml --local > /tmp/local.yaml
./bin/checkah check /tmp/local.yaml
```

Or pick a binary from the [latest release](https://github.com/deadc0de6/checkah/releases).

Or use the Dockerfile (by changing `localhost.yaml` to the config you want to use):
```bash
docker build -t checkah .
docker run -i checkah
```

# Build

```bash
## create a binary for your current host
make
./bin/checkah --help

## create all architecture binaries
make build-all
ls ./bin/
```

# Config

A few config examples are available under the [configs directory](/configs).
Config file can be written in yaml or json.

Config examples can be generated using the `example` command directly:
```bash
## generate a generic example config in json
bin/checkah example --format=json

## generate a generic example config in yaml
bin/checkah example --format=yaml

## generate a localhost example config in json
bin/checkah example --format=json --local

## generate a localhost example config in yaml
bin/checkah example --format=yaml --local
```

A config file is made of three main blocks:

* **settings**
* **hosts**
* **profiles**

*Note* that the config can be split across multiple files

## settings block

Global settings

* **hosts-parallel**: check hosts in parallel (optional, default `false`)
* **checks-parallel**: run checks in parallel (optional, default `true`)
* **global-alert**: an alert to trigger if any of the check fails (optional, see below for available alerts)
  * *type*: the alert type
  * *options* the alert options

## hosts block

A list of hosts to monitor

* **name**: arbitrary name to identify this host
* **host**: the host ip/domain
* **port**: the SSH port (optional, default 22)
* **user**: the SSH user (optional, default to the env variable `USER`)
* **password**: the SSH password (optional)
* **keyfile**: the SSH keyfile path (optional, default `~/.ssh/id_rsa`)
* **timeout**: SSH connection timeout in seconds (optional, default "3")
* **insecure**: disable known host checking if set to true (default `false`)
* **profiles**: a list of profile to apply to this host
* **disable**: a boolean indicating if the host is disabled (optional, default `false`)

if the *host* value is either `127.0.0.1` or `localhost`, SSH is disabled
and checks are run against localhost.

## profiles block

A list of profiles for monitoring hosts

* **name**: arbitrary name to identify this profile
* **extend**: name of another profile to include in this one (optional)
* **checks**: a list of checks (see below for the available checks)
  * *type*: the check type
  * *options*: the check options
  * *disable*: a boolean indicating if this check is disabled (optional, default `false`)
* **alerts**: a list of alerts (see below for the available alerts)
  * *type*: the alert type
  * *options* the alert options
  * *disable*: a boolean indicating if this alert is disabled (optional, default `false`)

The following checks are available:

* **disk**: check disk space used
  * *mount*: mount point (optional, default to `/`)
  * *limit*: if disk use percent crosses this value, an alert is triggered
* **loadavg**: check load average
  * *limit_1min*: if load average over 1 min crosses this value, an alert is triggered
  * *limit_5min*: if load average over 5 min crosses this value, an alert is triggered
  * *limit_15min*: if load average over 15 min crosses this value, an alert is triggered
* **memory**: check memory usage
  * *limit_mem*: if memory use percent crosses this value, an alert is triggered
  * *limit_swap*: if swap use percent crosses this value, an alert is triggered
  * *limit_total*: if total memory use percent crosses this value, an alert is triggered
* **process**: check if a process is running
  * *pattern*: pattern to match process name
  * *invert*: if value "yes", alert if process is present instead of absent (optional)
* **script**: run a custom check script on remote
  * *path*: the local path to the script
* **tcp**: check a specific TCP port is opened
  * *port*: TCP port to check
* **command**: check the return code of a command run on remote
  * *command*: the command

The following alerts are available:

* **file**: append to file
  * *path*: file path
* **script**: call a script with the alert string as sole argument
  * *path*: script path
* **webhook**: call a webhook on new alert
  * *url*: webhook url
  * *header<num>*: an header key (must start at `0`, optional)
  * *value<num>*: the corresponding value to *header<num>* (optional)
* **command**: execute a command on new alert
  * *command*: command string to run
* **email**: send an email on new alert
  * *host*: SMTP server address
  * *port*: SMTP server port
  * *mailfrom*: from email address
  * *mailto*: to email address
  * *user*: plain auth username (optional)
  * *password*: plain auth password (optional)

# Thank you

If you like checkah, [buy me a coffee](https://ko-fi.com/deadc0de6).

# License

This project is licensed under the terms of the GPLv3 license.

