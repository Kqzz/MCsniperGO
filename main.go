package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gookit/color"
	"github.com/kqzz/mcgo"
)

var accounts []mcgo.MCaccount

func init() {
	color.Printf(genHeader())
	accStrs, err := readLines("accounts.txt")
	if err != nil {
		logFatal(err.Error())
	}

	accounts = loadAccSlice(accStrs)
}

func main() {
	if len(accounts) < 1 {
		logFatal("Please put one account in the accounts.txt file!")
	}

	targetName := userInput("target username")
	offsetStr := userInput("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		logFatal(fmt.Sprintf("%v is not a valid integer", offsetStr))
	}

	fmt.Printf("Starting snipe for %v using offset %v\n", targetName, offset)

	droptime, err := coolkidmachoDroptime(targetName)
	if err != nil {
		logFatal(err.Error())
	}

	logInfo(fmt.Sprintf("Sniping %v at %v", targetName, droptime.Format("2006/01/02 15:04:05")))

	time.Sleep(time.Until(droptime.Add(-time.Hour * 8))) // sleep until 8 hours before droptime

	for _, acc := range accounts {
		err := acc.MojangAuthenticate()
		if err != nil {
			logErr(err.Error())
		}
	}

}
