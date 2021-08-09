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

func snipeCommand(targetName string, offset float64) {
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

	config, err := getConfig()

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

	if targetName == "" {
		targetName = userInput("target username")
	}

	if offset == 0 {
		var offsetStr string
		var offsetErr error

		for offsetStr == "" || offsetErr != nil {
			offsetStr = userInput("offset")
			offset, offsetErr = strconv.ParseFloat(offsetStr, 64)
			if offsetErr != nil {
				logErr(fmt.Sprintf("%v is not a valid number", offsetStr))
			}
		}
	}

	droptime, err := getDroptime(targetName, config.Sniper.TimingSystemPreference)
	if err != nil {
		logFatal(err.Error())
		return
	}

	logInfo(fmt.Sprintf("Sniping %v at %v\n", targetName, droptime.Format("2006/01/02 15:04:05")))

	time.Sleep(time.Until(droptime.Add(-time.Hour * 8))) // sleep until 8 hours before droptime

	// auth
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

		logInfo(fmt.Sprintf("Acc Type: %v | Bearer: %v", acc.Type, censor(acc.Bearer, 260)))
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

	// snipe
	for _, acc := range accounts {
		reqCount := config.Sniper.SnipeRequests
		if acc.Type == mcgo.MsPr {
			reqCount = config.Sniper.PrenameRequests
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
		time.Sleep(time.Millisecond * 1)
	}

	wg.Wait()

	for _, resp := range resps {
		logInfo(fmt.Sprintf("sent @ %v", resp.SendTime))
	}

	for _, resp := range resps {
		logInfo(fmt.Sprintf("[%v] recv @ %v", resp.StatusCode, resp.ReceiveTime))
		if resp.StatusCode < 300 {
			logSuccess(fmt.Sprintf("sniped %v onto %v", resp.Username, resp.Account.Email))
			logInfo("if you like this sniper please consider donating @ <fg=green;op=underscore>https://mcsniperpy.com/donate</>")
			if config.Sniper.AutoClaimNamemc {
				claimUrl, err := resp.Account.ClaimNamemc()
				if err != nil {
					logErr(fmt.Sprintf("failed to claim namemc: %v", err))
				} else {
					logInfo(fmt.Sprintf("namemc claim url: <fg=blue;op=underline>%v</>", claimUrl))
				}
			}
		}
	}

	fmt.Print("\n")

}

func pingCommand() {
	logInfo("Coming soon™")
}
