package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kqzz/MCsniperGO/log"
)

const help = `usage:
    mcsnipergo [options]
options:
    --username, -u <str>    username to snipe
    --offset, -o <float>    offset to use
    --autosnipe             auto snipe 3chars
    --queue, -q             run snipes from queue.txt 
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

	var startOffset float64

	flag.Float64Var(&startOffset, "offset", 0, "offset to use")
	flag.Float64Var(&startOffset, "o", 0, "offset to use")

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

		err := snipe(username)

		if err != nil {
			log.Log("err", "fatal: %v", err)
		}

		log.Input("snipe completed, press enter to continue")
	}

}
