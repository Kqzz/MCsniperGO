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

	normCount, prenameCount := countAccounts(accounts)

	if normCount > 1 {
		logWarn("using more than one normal account is useless")
	}

	if prenameCount > 10 {
		logWarn("using more than 10 prename accounts is useless")
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
	var authedAccounts []*mcgo.MCaccount
	for _, acc := range accounts {
		var authErr error
		if acc.Bearer != "" {
			logSuccess(fmt.Sprintf("successfully authenticated %v thru manual bearer", acc.Email))
			if acc.Type != mcgo.MsPr {
				loadAccErr := acc.LoadAccountInfo()
				if loadAccErr != nil {
					logErr("failed to load account info! invalid bearer, most likely.")
				} else {
					authedAccounts = append(authedAccounts, acc)
				}
			} else {
				logWarn("There are no guarentees that this bearer is correct, as it was manually inputted.")
			}
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
				authedAccounts = append(authedAccounts, acc)
			}
			time.Sleep(time.Duration(config.Accounts.AuthDelay) * time.Second)
		}

		if acc.Type != mcgo.MsPr {
			ncStatus, ncStatusErr := acc.NameChangeInfo()
			if ncStatusErr != nil {
				logWarn(fmt.Sprintf("failed to get name change info: %v", ncStatusErr))
			}
			if !ncStatus.Namechangeallowed {
				logErr(fmt.Sprintf("%v is not allowed to name change!", acc.Username))
				if acc.Email != "" {
					accounts = removeAccount(accounts, findAccByEmail(accounts, acc))
					// this might not work idk
				} else {
					logWarn("cannot remove acc due to manual bearer mode")
				}
			} else {
				logSuccess(fmt.Sprintf("verified that %v can name change", acc.Email))
			}
		}

		logInfo(fmt.Sprintf("Acc Type: %v | Bearer: %v", prettyAccType(acc.Type), censor(acc.Bearer, 260)))
	}

	fmt.Print("\n")

	if len(authedAccounts) == 0 {
		logErr("0 accounts authenticated successfully - stopping snipe.")
		return
	}

	changeTime := droptime.Add(time.Millisecond * time.Duration(0-offset))

	var wg sync.WaitGroup
	var resps []mcgo.NameChangeReturn

	for time.Now().Before(changeTime.Add(-time.Second * 20)) {
		color.Printf("sniping in <fg=blue>%s</>       \r", time.Until(droptime).Round(time.Second))
		time.Sleep(time.Second * 1)
	}

	fmt.Println("\nstarting...")
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
					logErr(fmt.Sprintf("encountered err on nc for %v: %v", acc.Email, err.Error()))
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
		logInfo(fmt.Sprintf("sent @ %v", fmtTimestamp(resp.SendTime)))
		logsSlice = append(logsSlice, fmt.Sprintf("sent @ %v", fmtTimestamp(resp.SendTime)))
	}

	for _, resp := range resps {
		logInfo(fmt.Sprintf("[%v] received @ %v", resp.StatusCode, fmtTimestamp(resp.ReceiveTime)))
		logsSlice = append(logsSlice, fmt.Sprintf("[%v] received @ %v", resp.StatusCode, fmtTimestamp(resp.ReceiveTime)))
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

			if config.Announce.McsnipergoAnnounceCode != "" {
				err := announceSnipe(targetName, config.Announce.McsnipergoAnnounceCode, &resp.Account)

				if err != nil {
					logErr(fmt.Sprintf("failed to announce snipe: %v", err))
				} else {
					logSuccess("announced snipe!")
				}
			}
		}
	}

	fmt.Print("\n")

	if !fileExists("logs") {
		err = os.Mkdir("logs", 0755)
		if err != nil {
			logFatal(fmt.Sprintf("Failed to create logs folder: %v", err))
		}
	}

	logFile, err := os.OpenFile(fmt.Sprintf("logs/%v.txt", targetName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logFatal(fmt.Sprintf("Failed to create log file: %v", err))
	}

	defer logFile.Close()

	logFile.WriteString(strings.Join(logsSlice, "\n"))

}

func pingCommand() {
	logInfo("Coming soon™")
}
