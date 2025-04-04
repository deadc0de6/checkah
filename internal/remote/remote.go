// Copyright (c) 2021 deadc0de6

package remote

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/deadc0de6/checkah/internal/alert"
	"github.com/deadc0de6/checkah/internal/check"
	"github.com/deadc0de6/checkah/internal/config"
	"github.com/deadc0de6/checkah/internal/output"
	"github.com/deadc0de6/checkah/internal/transport"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

const (
	maxJobs = 4
)

var (
	hostLocalhost = []string{"127.0.0.1", "localhost"}
)

// HostResult host result struct
type HostResult struct {
	NbCheckTotal int
	NbCheckError int
}

// Remote a remote host to check
type Remote struct {
	Name              string
	Host              string
	Port              string
	User              string
	Password          string
	Keyfile           string
	Checks            []check.Check
	Alerts            []alert.Alert
	Timeout           int
	KnownHostInsecure bool
	Disable           bool
}

type profileStruct struct {
	checks []check.Check
	alerts []alert.Alert
}

// ToRemote convert a config to a list of remote struct
func ToRemote(cfg *config.Config) ([]*Remote, error) {
	// create profile map
	profiles := make(map[string]*profileStruct)
	isReachable, _ := check.GetCheck("reachable", nil)

	for _, profile := range cfg.Profiles {
		p := profileStruct{}

		// add the checks
		for _, ch := range profile.Checks {
			if ch.Disable {
				continue
			}
			// is valid
			if len(ch.Type) < 1 {
				return nil, fmt.Errorf("check type cannot be empty")
			}
			if len(ch.Options) < 1 {
				return nil, fmt.Errorf("check \"%s\" options cannot be empty", ch.Type)
			}

			// construct check
			checker, err := check.GetCheck(ch.Type, ch.Options)
			if err != nil {
				return nil, fmt.Errorf("check %s: %v", ch.Type, err)
			}
			p.checks = append(p.checks, checker)
		}
		// add the alerts
		for _, al := range profile.Alerts {
			if al.Disable {
				continue
			}
			a, err := alert.GetAlert(al.Type, al.Options)
			if err != nil {
				return nil, fmt.Errorf("alert %s: %v", al.Type, err)
			}
			p.alerts = append(p.alerts, a)
		}
		profiles[profile.Name] = &p
	}

	// process the profile inclusion
	for _, profile := range cfg.Profiles {
		if len(profile.Extend) > 0 {
			name := profile.Name
			p, ok := profiles[name]
			if !ok {
				return nil, fmt.Errorf("unknown profile named \"%s\"", name)
			}

			// loop
			for _, other := range profile.Extend {
				// add other profile checks and alerts
				o, ok := profiles[other]
				if !ok {
					return nil, fmt.Errorf("unknown profile named \"%s\" in extend", name)
				}
				p.checks = append(p.checks, o.checks...)
				p.alerts = append(p.alerts, o.alerts...)
			}
		}
	}

	// create the remotes
	var remotes []*Remote
	for _, host := range cfg.Hosts {
		var thisChecks []check.Check
		var thisAlerts []alert.Alert

		if host.Disable {
			continue
		}

		// add the isReachable check
		thisChecks = append(thisChecks, isReachable)

		for _, proName := range host.ProfileNames {
			p, ok := profiles[proName]
			if !ok {
				return nil, fmt.Errorf("no such profile: %s", proName)
			}
			thisChecks = append(thisChecks, p.checks...)
			thisAlerts = append(thisAlerts, p.alerts...)
		}

		user := host.User
		if len(user) < 1 {
			user = os.Getenv("USER")
		}

		port := host.Port
		if len(host.Port) < 1 {
			port = "22"
		}

		timeout := host.Timeout
		if len(host.Timeout) < 1 {
			timeout = "3"
		}
		timeoutVal, err := strconv.Atoi(timeout)
		if err != nil {
			return nil, err
		}

		r := &Remote{
			Name:              host.Name,
			Host:              host.Host,
			Port:              port,
			User:              user,
			Password:          host.Password,
			Keyfile:           host.Keyfile,
			Checks:            thisChecks,
			Alerts:            thisAlerts,
			Timeout:           timeoutVal,
			KnownHostInsecure: host.KnownHostInsecure,
		}
		remotes = append(remotes, r)
	}

	return remotes, nil
}

// PrintRemote prints a remote
func PrintRemote(remote *Remote) {
	fmt.Printf("Remote \"%s\" (%s):\n", remote.Name, remote.Host)

	// checks
	fmt.Printf("  Checks:\n")
	for _, check := range remote.Checks {
		fmt.Printf("    %s: %s\n", check.GetName(), check.GetDescription())
		for k, v := range check.GetOptions() {
			fmt.Printf("      - %s=%s\n", k, v)
		}
	}

	// alerts
	fmt.Printf("  Alerts:\n")
	for _, alert := range remote.Alerts {
		fmt.Printf("    description: %s\n", alert.GetDescription())
		for k, v := range alert.GetOptions() {
			fmt.Printf("      - %s=%s\n", k, v)
		}
	}
}

// PrintRemotes print remotes
func PrintRemotes(remotes []*Remote) {
	for _, remote := range remotes {
		PrintRemote(remote)
	}
}

func notify(host string, content string, alerts []alert.Alert) {
	for _, a := range alerts {
		log.Debugf("notify with %s", a.GetDescription())
		line := fmt.Sprintf("ALERT \"%s\" - %s", host, content)
		err := a.Notify(line)
		if err != nil {
			c := fmt.Sprintf("notification error for \"%s\": ", a.GetDescription())
			col := color.New(color.FgRed)
			fmt.Print(c + col.Sprintln(err.Error()))
		}
	}
}

func isLocalhost(host string) bool {
	for _, n := range hostLocalhost {
		if n == host {
			return true
		}
	}
	return false
}

// CheckRemote runs the check against a remote
func CheckRemote(remote *Remote, parallel bool, resChan chan *HostResult, doneFunc *sync.WaitGroup, out output.Output) {
	// create the transport
	var trans transport.Transport
	var err error
	var nbChecks int

	outputKey := fmt.Sprintf("%s (%s:%s)", remote.Name, remote.Host, remote.Port)

	log.Debugf("connecting to %s...", remote.Name)
	if isLocalhost(remote.Host) {
		trans, err = transport.NewLocal()
	} else {
		var keyfiles []string
		if len(remote.Keyfile) > 0 {
			keyfiles = append(keyfiles, remote.Keyfile)
		}
		trans, err = transport.NewSSH(remote.Host, remote.Port, remote.User, remote.Password, keyfiles, remote.Timeout, remote.KnownHostInsecure)
	}

	if err != nil {
		out.StackErr(outputKey, "not reachable", err.Error())
		notify(remote.Name, fmt.Sprintf("host %s is not reachable: %v", remote.Name, err), remote.Alerts)
		out.Flush(outputKey)
		resChan <- &HostResult{
			NbCheckTotal: nbChecks,
			NbCheckError: 1,
		}
		doneFunc.Done()
		return
	}

	// defer closing the sessions
	defer trans.Close()

	// create the result channel
	ch := make(chan *check.Result, len(remote.Checks))
	// create the jobs channel
	maxJob := 1
	if parallel {
		maxJob = maxJobs
	}
	jobs := make(chan check.Check, maxJob)
	// create the end of process channel
	endChan := make(chan int, 1)

	// checker worker
	// reads checks from jobs channel
	// and push results to result channel
	go func() {
		for check := range jobs {
			log.Debugf("running check %s", check.GetDescription())
			res := check.Run(trans)
			ch <- res
		}
		close(ch)
	}()

	// process results worker
	// handles the results and construct output
	go func() {
		errCnt := 0
		for res := range ch {
			if res.Error != nil {
				// alert notification
				errStr := fmt.Sprintf("%s: %s", res.Description, res.Error)
				notify(remote.Name, errStr, remote.Alerts)
				// output
				out.StackErr(outputKey, res.Description, res.Error.Error())
				errCnt++
			} else {
				out.StackOk(outputKey, res.Description, res.Value)
			}
		}
		endChan <- errCnt
	}()

	// send the jobs
	for _, c := range remote.Checks {
		nbChecks++
		jobs <- c
	}
	close(jobs)

	// wait for result processing to end
	errCnt := <-endChan

	// print output
	out.Flush(outputKey)
	resChan <- &HostResult{
		NbCheckTotal: nbChecks,
		NbCheckError: errCnt,
	}
	doneFunc.Done()
}
