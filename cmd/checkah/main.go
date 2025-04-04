// Copyright (c) 2021 deadc0de6

package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/deadc0de6/checkah/internal/alert"
	"github.com/deadc0de6/checkah/internal/config"
	"github.com/deadc0de6/checkah/internal/output"
	"github.com/deadc0de6/checkah/internal/remote"

	"github.com/docopt/docopt-go"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

// Switches the command line options
type Switches struct {
	// actions
	Print   bool `docopt:"print"`
	Check   bool `docopt:"check"`
	Example bool `docopt:"example"`
	// args
	Paths []string `docopt:"<path>"`
	// options
	Local   bool   `docopt:"-l,--local"`
	Format  string `docopt:"-f,--format"`
	Verbose bool   `docopt:"-v,--verbose"`
	Version bool   `docopt:"--version"`
	Help    bool   `docopt:"-h,--help"`
}

var (
	version = "0.2.4"
	name    = "checkah"
	usage   = `checkah.

Usage:
	checkah check [-v] <path>...
	checkah print [-v] [--format=<format>] <path>...
	checkah example [-lv] [--format=<format>]
	checkah -h | --help
	checkah --version

Options:
  -l --local              Generate localhost config example.
  -f --format=<format>    Output format [default: yaml].
  -v --verbose            Debug logs.
  -h --help               Show this screen.
  --version               Show version.`
)

func printUsage() {
	fmt.Println(usage)
	os.Exit(1)
}

func cmdExample(format string, local bool) int {
	err := config.PrintExampleConfig(format, local)
	if err != nil {
		log.Fatal(err)
	}
	return 0
}

func cmdPrint(configs []string, format string) int {
	cfg, _, err := parseConfigs(configs)
	if err != nil {
		log.Fatal(err)
	}

	err = config.PrintConfig(cfg, format)
	if err != nil {
		log.Fatal(err)
	}

	return 0
}

// returns total hosts, total checks, hosts error, nbError
func cmdCheck(configs []string) (int, int, int, int) {
	cfg, remotes, err := parseConfigs(configs)
	if err != nil {
		log.Fatal(err)
	}

	hostsParallel := cfg.Settings.HostsParallel
	checksParallel := cfg.Settings.ChecksParallel
	globalAlert, _ := alert.GetAlert(cfg.Settings.GlobalAlert.Type, cfg.Settings.GlobalAlert.Options)

	log.Debugf("hosts parallel: %t", hostsParallel)
	log.Debugf("checks parallel: %t", checksParallel)
	log.Debugf("global alert: %v", globalAlert)

	var wg sync.WaitGroup
	ch := make(chan *remote.HostResult, len(remotes))

	// create the output
	out, _ := output.GetOutput("stdout", nil)

	// check all hosts
	errCnt := 0
	hostErrCnt := 0
	checksCnt := 0
	for _, r := range remotes {
		wg.Add(1)
		log.Debugf("launching checks on %s", r.Name)
		go remote.CheckRemote(r, checksParallel, ch, &wg, out)
		if !hostsParallel {
			wg.Wait()
		}
	}

	wg.Wait()
	close(ch)

	// process results
	for res := range ch {
		if res.NbCheckError > 0 {
			hostErrCnt++
		}
		errCnt += res.NbCheckError
		checksCnt += res.NbCheckTotal
	}

	if globalAlert != nil && errCnt > 0 {
		line := fmt.Sprintf("check failed: %d/%d host(s) failed (check error: %d)", hostErrCnt, len(remotes), errCnt)
		err := globalAlert.Notify(line)
		if err != nil {
			log.Errorf("%v", err)
		}
	}
	return len(remotes), checksCnt, hostErrCnt, errCnt
}

func parseConfigs(paths []string) (*config.Config, []*remote.Remote, error) {
	c := &config.Config{}
	for _, path := range paths {
		cfg, err := config.ReadCfg(path)
		if err != nil {
			return nil, nil, err
		}
		c, err = config.MergeConfigs(c, cfg)
		if err != nil {
			return nil, nil, err
		}
	}
	remotes, err := remote.ToRemote(c)
	if err != nil {
		return nil, nil, err
	}

	if log.GetLevel() == log.DebugLevel {
		remote.PrintRemotes(remotes)
	}

	return c, remotes, nil
}

func main() {
	// parse cli switches
	args, err := docopt.ParseArgs(usage, nil, version)
	if err != nil {
		log.Fatal(err)
	}

	var opts Switches
	err = args.Bind(&opts)
	if err != nil {
		log.Fatal(err)
	}

	if opts.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	if opts.Help {
		printUsage()
	}

	fmt.Printf("%s v%s\n", name, version)
	if opts.Version {
		os.Exit(0)
	}

	ret := 1
	if opts.Print {
		paths := opts.Paths
		if len(paths) < 1 {
			printUsage()
		}
		ret = cmdPrint(paths, opts.Format)
	} else if opts.Example {
		ret = cmdExample(opts.Format, opts.Local)
	} else if opts.Check {
		paths := opts.Paths
		if len(paths) < 1 {
			printUsage()
		}
		total, totalChecks, hostErr, errCnt := cmdCheck(paths)
		ret = errCnt
		errStr := fmt.Sprintf("%d", errCnt)
		if errCnt > 0 {
			red := color.New(color.FgRed).SprintFunc()
			errStr = red(errCnt)
		}
		hostErrStr := fmt.Sprintf("%d", hostErr)
		if hostErr > 0 {
			red := color.New(color.FgRed).SprintFunc()
			hostErrStr = red(hostErr)
		}
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("\nChecked %d hosts (%d checks):\n", total, totalChecks)
		fmt.Printf("%s success, %s failed (%s total errors)\n", green(total-hostErr), hostErrStr, errStr)
	}

	if ret != 0 {
		os.Exit(1)
	}
	os.Exit(0)
}
