package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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
			log("while creating accounts.txt, %s", 4, err.Error())
			return
		} else {
			log("created accounts.txt, please restart the sniper once accounts are added!", 0)
		}
	}

	if !fileExists("config.toml") {
		defaultConfig()
	}

	config, err := getConfig()

	if err != nil {
		log("error while getting config, %v", 4, err)
		return
	}

	accStrs, err := readLines("accounts.txt")
	if err != nil {
		log(err.Error(), 4)
		return
	}

	accounts = loadAccSlice(accStrs)

	if len(accounts) < 1 {
		log("Please put one account in the accounts.txt file!", 4)
		return
	}

	normCount, prenameCount := countAccounts(accounts)

	if normCount > 1 {
		log("using more than one normal account is useless", 3)
	}

	if prenameCount > 10 {
		log("using more than 10 prename accounts is useless", 3)
	}

	if targetName == "" {
		targetName = userInput("target username")
	}

	if offset == -10000 {
		var offsetStr string
		var offsetErr error

		for offsetStr == "" || offsetErr != nil {
			offsetStr = userInput("offset")
			offset, offsetErr = strconv.ParseFloat(offsetStr, 64)
			if offsetErr != nil {
				log("%v is not a valid number", 2, offsetStr)
			}
		}
	}

	droptime, err := getDroptime(targetName, config.Sniper.TimingSystemPreference)
	if err != nil {
		log(err.Error(), 4)
		return
	}

	log("Sniping %v at %v\n", 0, targetName, droptime.Format("2006/01/02 15:04:05"))

	time.Sleep(time.Until(droptime.Add(-time.Minute * time.Duration(config.Accounts.StartAuth)))) // sleep until 8 hours before droptime

	var authedAccounts []*mcgo.MCaccount

	for _, acc := range accounts {
		if authAccountErr := authAccount(acc); authAccountErr != nil {
			log("failed to authenticate %v: %v", 2, accID(acc), authAccountErr)
		} else {
			log("successfully authenticated %v", 1, accID(acc))
		}

		canSnipe, canSnipeErr := accReadyToSnipe(acc)
		if canSnipeErr != nil {
			log("failed to verify that %v can snipe: %v", 2, accID(acc), canSnipeErr)
		}

		if canSnipe {
			log("verified that %v can snipe", 1, accID(acc))
			authedAccounts = append(authedAccounts, acc)
		} else {
			log("%v not ready to snipe", 2, accID(acc))
		}
	}

	if len(authedAccounts) == 0 {
		log("no accounts successfully authenticated!", 4)
		return
	}

	changeTime := droptime.Add(time.Millisecond * time.Duration(0-offset))

	var wg sync.WaitGroup
	var resps []mcgo.NameChangeReturn

	for time.Now().Before(changeTime.Add(-time.Second * 20)) {
		color.Printf("sniping in <fg=blue>%s</>       \r", time.Until(droptime).Round(time.Second))
		time.Sleep(time.Second * 1)
	}

	fmt.Println("\nstarting in 20s...")
	var totalReqCount int // keep track of the total requests for the spread

	// snipe
	for _, acc := range authedAccounts {
		reqCount := config.Sniper.SnipeRequests
		if acc.Type == mcgo.MsPr {
			reqCount = config.Sniper.PrenameRequests
		}
		for i := 0; i < reqCount; i++ {
			totalReqCount++
			wg.Add(1)
			prename := false
			if acc.Type == mcgo.MsPr {
				prename = true
			}
			spread := float64(totalReqCount) * config.Sniper.Spread
			go func() {
				defer wg.Done()
				resp, err := acc.ChangeName(targetName, changeTime.Add(time.Millisecond*time.Duration(spread)), prename)

				if err != nil {
					log("encountered err on nc for %v: %v", 2, acc.Email, err.Error())
				} else {
					resps = append(resps, resp)
				}
			}()
		}
		time.Sleep(time.Millisecond * 1)
	}

	wg.Wait()

	logsSlice := []string{
		"accounts",
	}

	for _, acc := range accounts {
		logsSlice = append(logsSlice, formatAccount(acc))
	}

	logsSlice = append(logsSlice, "logs")

	for _, resp := range resps {
		log("sent @ %v", 0, fmtTimestamp(resp.SendTime))
		logsSlice = append(logsSlice, fmt.Sprintf("sent @ %v", fmtTimestamp(resp.SendTime)))
	}

	for _, resp := range resps {
		log("[%v] received @ %v | est process @ %v", 0, prettyStatus(resp.StatusCode), fmtTimestamp(resp.ReceiveTime), fmtTimestamp(estimatedProcess(resp.SendTime, resp.ReceiveTime)))
		logsSlice = append(logsSlice, fmt.Sprintf("[%v] received @ %v", resp.StatusCode, fmtTimestamp(resp.ReceiveTime)))
		if resp.StatusCode < 300 {
			log("sniped %v onto %v", 1, resp.Username, resp.Account.Email)
			log("if you like this sniper please consider donating @ <fg=green;op=underscore>https://mcsniperpy.com/donate</>", 0)
			if config.Sniper.AutoClaimNamemc {
				claimUrl, err := resp.Account.ClaimNamemc()
				if err != nil {
					log("failed to claim namemc: %v", 2, err)
				} else {
					log("namemc claim url: <fg=blue;op=underline>%v</>", 0, claimUrl)
				}
			}

			if config.Announce.McsnipergoAnnounceCode != "" {
				err := announceSnipe(targetName, config.Announce.McsnipergoAnnounceCode, &resp.Account)

				if err != nil {
					log("failed to announce snipe: %v", 2, err)
				} else {
					log("announced snipe!", 1)
				}
			}
		}
	}

	fmt.Print("\n")

	if !fileExists("logs") {
		err = os.Mkdir("logs", 0755)
		if err != nil {
			log("Failed to create logs folder: %v", 4, err)
		}
	}

	logFile, err := os.OpenFile(fmt.Sprintf("logs/%v.txt", targetName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log("Failed to create log file: %v", 4, err)
	}

	defer logFile.Close()

	logFile.WriteString(strings.Join(logsSlice, "\n"))

}

func pingCommand() {
	log("Coming soon™", 0)
}
