package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Kqzz/MCsniperGO/log"
)

func main() {

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Print("\r")
		log.Log("err", "ctrl-c pressed, exiting...")
		os.Exit(0)
	}()

	for {

		log.Log("", log.GetHeader())

		username := log.Input("target username")

		var offset float64

	offsetLoop:
		for {
			o := log.Input("offset")
			var err error
			offset, err = strconv.ParseFloat(o, 64)
			if err == nil {
				break offsetLoop
			}
		}

		err := snipe(username, offset)

		if err != nil {
			log.Log("err", "fatal: %v", err)
		}

		log.Input("snipe completed, press enter to continue")
	}

}
