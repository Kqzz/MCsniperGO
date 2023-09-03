package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kqzz/MCsniperGO/claimer"
	"github.com/Kqzz/MCsniperGO/log"
	"github.com/Kqzz/MCsniperGO/parser"
)

const help = `usage:
    mcsnipergo [options]
options:
    --username, -u <str>    username to snipe
`

func init() {
	flag.Usage = func() {
		fmt.Print(help)
	}
}

func isFlagPassed(names ...string) bool {
	found := false
	for _, name := range names {
		flag.Visit(func(f *flag.Flag) {
			if f.Name == name {
				found = true
			}
		})
	}
	return found
}

func main() {

	var startUsername string
	flag.StringVar(&startUsername, "username", "", "username(s) to snipe")
	flag.StringVar(&startUsername, "u", "", "username(s) to snipe")

	flag.Parse()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Print("\r")
		log.Log("err", "ctrl-c pressed, exiting...      ")
		os.Exit(0)
	}()

	for {

		log.Log("", log.GetHeader())

		var username string

		if !isFlagPassed("u", "username") {
			username = log.Input("target username(s)")
		} else {
			username = startUsername
		}

		dropRange := log.GetDropRange()

		proxies, err := parser.ReadLines("proxies.txt")

		if err != nil {
			log.Log("err", "failed to load proxies: %v", err)
		}

		err = nil

		accounts, err := getAccounts("gc.txt", "gp.txt", "ms.txt")

		err = claimer.ClaimWithinRange(username, dropRange, accounts, proxies)

		if err != nil {
			log.Log("err", "fatal: %v", err)
		}

		log.Input("snipe completed, press enter to continue")
	}

}
