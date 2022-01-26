package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"bufio"

	"github.com/gookit/color"
	"github.com/kqzz/mcgo"
)

func snipeCommand(targetName string, offset float64) error {
	if !fileExists("accounts.txt") {
		_, err := os.Create("accounts.txt")
		if err != nil {
			return fmt.Errorf("while creating accounts.txt, %s", err)
		} else {
			log("info", "created accounts.txt, please restart the sniper once accounts are added!")
		}
	}

	if !fileExists("config.toml") {
		defaultConfig()
	}

	config, err := getConfig()

	if err != nil {
		return fmt.Errorf("error while getting config, %v", err)
	}

	accStrs, err := readLines("accounts.txt")
	if err != nil {
		return err
	}

	accounts = loadAccSlice(accStrs)

	if len(accounts) < 1 {
		return fmt.Errorf("please put one account in the accounts.txt file")
	}

	normCount, prenameCount := countAccounts(accounts)

	if normCount > 1 {
		log("warn", "using more than one normal account is useless")
	}

	if prenameCount > 10 {
		log("warn", "using more than 10 prename accounts is useless")
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
				log("error", "%v is not a valid number", offsetStr)
			}
		}
	}

	droptime, err := getDroptime(targetName, config.Sniper.TimingSystemPreference)
	if err != nil {
		return err
	}

	log("info", "Sniping %v at %v\n", targetName, droptime.Format("2006/01/02 15:04:05"))

	time.Sleep(time.Until(droptime.Add(-time.Minute * time.Duration(config.Accounts.StartAuth)))) // sleep until StartAuth minutes before drop

	var authedAccounts []*mcgo.MCaccount

	// auth + checking if ready to snipe
	for _, acc := range accounts {
		authAccountErr := authAccount(acc)
		if authAccountErr != nil {
			log("error", "failed to authenticate %v: %v", accID(acc), authAccountErr)
		} else {
			log("success", "successfully authenticated %v", accID(acc))
		}

		if authAccountErr == nil {
			canSnipe, canSnipeErr := accReadyToSnipe(acc)
			if canSnipeErr != nil {
				log("error", "failed to verify that %v can snipe: %v", accID(acc), canSnipeErr)
			}

			if canSnipe {
				log("success", "verified that %v can snipe", accID(acc))
				authedAccounts = append(authedAccounts, acc)
			} else {
				log("error", "%v not ready to snipe", accID(acc))
			}
		}
		time.Sleep(time.Duration(config.Accounts.AuthDelay) * time.Second)
	}

	if len(authedAccounts) == 0 {
		return fmt.Errorf("no accounts successfully authenticated")
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
					log("error", "encountered err on nc for %v: %v", acc.Email, err.Error())
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

	logsSlice = append(logsSlice, fmt.Sprintf("offset = %v", offset))

	logsSlice = append(logsSlice, "logs")


	for _, resp := range resps {
		log(
			"info", "[%v] sent @ %v | recv @ %v | %v",
			prettyStatus(resp.StatusCode),
			fmtTimestamp(resp.SendTime),
			fmtTimestamp(resp.ReceiveTime),
			accID(&resp.Account),
		)

		logsSlice = append(logsSlice, fmt.Sprintf(
			"[%v] sent @ %v | recv @ %v | %v",
			resp.StatusCode,
			fmtTimestamp(resp.SendTime),
			fmtTimestamp(resp.ReceiveTime),
			accID(&resp.Account),
		))

		if resp.StatusCode < 300 {

			log("success", "sniped %v onto %v", resp.Username, resp.Account.Email)
			log("info", "if you like this sniper please consider donating @ <fg=green;op=underscore>https://mcsniperpy.com/donate</>")

			if config.Announce.McsnipergoAnnounceCode != "" {
				err := announceSnipe(targetName, config.Announce.McsnipergoAnnounceCode, &resp.Account)

				if err != nil {
					log("error", "failed to announce snipe: %v", err)
				} else {
					log("success", "announced snipe!")
				}
			}

			if config.Announce.WebhookURL != "" {
				err := customServerAnnounce(targetName)
					if err != nil {
						log("error", "failed to announce snipe to your webhook: %v", err)
					} else {
						log("succes", "announced your snipe!")
					}
			}
		}
	}

	fmt.Print("\n")

	if !fileExists("logs") {
		err = os.Mkdir("logs", 0755)
		if err != nil {
			return fmt.Errorf("failed to create logs folder: %v", err)
		}
	}

	logFile, err := os.OpenFile(fmt.Sprintf("logs/%v.txt", targetName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create log file: %v", err)
	}

	defer logFile.Close()

	logFile.WriteString(strings.Join(logsSlice, "\n"))

	return nil
}
func autoSnipeCommand(offset float64) error {
	for {
		nameSlice, err := getNext3c()
		if err != nil {
			return err
		}
		for _, i := range nameSlice {

			if offset == -10000 {
				var offsetStr string
				var offsetErr error

				for offsetStr == "" || offsetErr != nil {
					offsetStr = userInput("offset")
					offset, offsetErr = strconv.ParseFloat(offsetStr, 64)
					if offsetErr != nil {
						log("error", "%v is not a valid number", offsetStr)
					}
				}
			}
			_, err = starShoppingDroptime(i.Name)
			if err == nil {
				err = snipeCommand(i.Name, offset)
				if err != nil {
					return err
				}
			}
		}
	}
}
func pingCommand() {
	ping, err := pingMojang()
	if err != nil {
		log("fatal", "failed to ping mojang: %v", err)
		return
	}
	log("info", "ping: %vms", ping)
}
func snipeList() {
	offsetStr := userInput("offset")
		offset, err := strconv.ParseFloat(offsetStr, 64)
			if err !=  nil {
				log("error", "invalid offset number: %v", err)
			}

	if !fileExists("names.txt") {
		_, err := os.Create("names.txt")
		if err != nil {
			log("fatal", "failed to create names.txt: %v", err)
		} else {
			log("info", "created names.txt")
		}
	}

	file, err := os.Open("names.txt")
		if err != nil {
		log("error", "failed to open 'names.txt': %v", err)
		}

		var names []string
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		names = append(names, fileScanner.Text())
	}
	if len(names) > 0 {
		for _, name := range names {
		snipeCommand(name, offset)
		}
		log("info", "finished sniping each name on your list")
		userInput("press enter to exit")
	} else {
		log("error", "Your names.txt file is empty! One name per line.")
		userInput("press enter to exit")
	}
}