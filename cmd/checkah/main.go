// Copyright (c) 2021 deadc0de6

package main

import (
	"checkah/internal/alert"
	"checkah/internal/config"
	"checkah/internal/remote"
	"flag"
	"fmt"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

var (
	version        = "0.1.2"
	cmdShowString  = "show"
	cmdCheckString = "check"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  %s show <config-path>...\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s check <config-path>...\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func cmdShow(configs []string) int {
	cfg, remotes, err := parseConfigs(configs)
	if err != nil {
		log.Fatal(err)
	}

	config.PrintSettings(cfg)
	remote.PrintRemotes(remotes)
	return 0
}

// returns total hosts, total checks, hosts error, nbError
func cmdCheck(configs []string) (int, int, int, int) {
	cfg, remotes, err := parseConfigs(configs)
	if err != nil {
		log.Fatal(err)
	}

	config.PrintSettings(cfg)
	hostsParallel := cfg.Settings.HostsParallel
	checksParallel := cfg.Settings.ChecksParallel
	globalAlert, _ := alert.GetAlert(cfg.Settings.GlobalAlert.Type, cfg.Settings.GlobalAlert.Options)

	var wg sync.WaitGroup
	ch := make(chan *remote.HostResult, len(remotes))

	// check all hosts
	errCnt := 0
	hostErrCnt := 0
	checksCnt := 0
	for _, r := range remotes {
		wg.Add(1)
		go remote.CheckRemote(r, checksParallel, ch, &wg)
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
		globalAlert.Notify(line)
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

	return c, remotes, nil
}

func main() {
	debugArg := flag.Bool("debug", false, "debug mode")
	helpArg := flag.Bool("help", false, "Show usage")
	versArg := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *helpArg {
		usage()
	}

	if *debugArg {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}

	if *versArg {
		fmt.Printf("%s v%s\n", os.Args[0], version)
		os.Exit(0)
	}

	rest := flag.Args()
	if len(rest) < 1 {
		usage()
	}
	cmd := rest[0]
	ret := 1
	switch cmd {
	case cmdShowString:
		ret = cmdShow(rest[1:])
	case cmdCheckString:
		total, totalChecks, hostErr, errCnt := cmdCheck(rest[1:])
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
	default:
		usage()
		ret = 1
	}

	if ret != 0 {
		os.Exit(1)
	}
	os.Exit(0)
}
