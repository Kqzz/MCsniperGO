package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gookit/color"
	"github.com/kqzz/mcgo"
)

func snipeCommand() {

	color.Printf(genHeader())
	if !fileExists("accounts.txt") {
		_, err := os.Create("accounts.txt")
		if err != nil {
			logFatal(fmt.Sprintf("while creating accounts.txt, %s", err.Error()))
			return
		} else {
			logInfo("created accounts.txt, please restart the sniper once accounts are added!")
		}
	}

	if !fileExists("config.toml") {
		defaultConfig()
	}

	_, err := getConfig()

	if err != nil {
		logFatal(fmt.Sprintf("error while getting config, %v", err))
		return
	}

	accStrs, err := readLines("accounts.txt")
	if err != nil {
		logFatal(err.Error())
		return
	}

	accounts = loadAccSlice(accStrs)

	if len(accounts) < 1 {
		logFatal("Please put one account in the accounts.txt file!")
		return
	}

	if len(accounts) > 1 {
		logWarn("Using more than 1 account is not recommended")
	}

	targetName := userInput("target username")

	var offsetStr string
	var offset float64

	var offsetErr error

	for offsetStr == "" || offsetErr != nil {
		offsetStr = userInput("offset")
		offset, offsetErr = strconv.ParseFloat(offsetStr, 64)
		if offsetErr != nil {
			logErr(fmt.Sprintf("%v is not a valid number", offsetStr))
		}
	}

	droptime, err := getDroptime(targetName, "ckm")
	if err != nil {
		logFatal(err.Error())
	}

	logInfo(fmt.Sprintf("Sniping %v at %v\n", targetName, droptime.Format("2006/01/02 15:04:05")))

	time.Sleep(time.Until(droptime.Add(-time.Hour * 8))) // sleep until 8 hours before droptime

	for _, acc := range accounts {
		var authErr error
		if acc.Bearer != "" {
			logSuccess(fmt.Sprintf("successfully authenticated %v thru manual bearer", acc.Email))
			logWarn("There are no guarentees that this bearer is correct, as it was manually inputted.")
		} else {
			if acc.Type == mcgo.Mj {
				authErr = acc.MojangAuthenticate()
			} else {
				authErr = acc.MicrosoftAuthenticate()
			}
			if authErr != nil {
				logErr(fmt.Sprintf("Failed to authenticate %v, err: \"%v\"", acc.Email, authErr.Error()))
			} else {
				logSuccess(fmt.Sprintf("successfully authenticated %v", acc.Email))
			}
		}
	}

	fmt.Print("\n")

	changeTime := droptime.Add(time.Millisecond * time.Duration(0-offset))

	var wg sync.WaitGroup

	var resps []mcgo.NameChangeReturn

	for time.Now().Before(changeTime.Add(-time.Second * 40)) {
		color.Printf("sniping in <fg=blue>%vs</>       \r", time.Until(droptime).Round(time.Second).Seconds())
		time.Sleep(time.Second * 1)
	}

	fmt.Println("\nstarting...")

	for _, acc := range accounts {
		reqCount := 2
		if acc.Type == mcgo.MsPr {
			reqCount = 6
		}
		for i := 0; i < reqCount; i++ {
			wg.Add(1)
			prename := false
			if acc.Type == mcgo.MsPr {
				prename = true
			}
			go func() {
				defer wg.Done()
				resp, err := acc.ChangeName(targetName, changeTime, prename)
				if err != nil {
					logErr(fmt.Sprintf("encountered err on nc for %v: %v", acc.Email, err.Error()))
				} else {
					resps = append(resps, resp)
				}
			}()
		}
	}

	wg.Wait()

	for _, resp := range resps {
		logInfo(fmt.Sprintf("sent @ %v", resp.SendTime))
	}

	for _, resp := range resps {
		logInfo(fmt.Sprintf("[%v] recv @ %v", resp.StatusCode, resp.ReceiveTime))
	}

	fmt.Print("\n")

}

func pingCommand() {
	logInfo("PING PING PING")
}
