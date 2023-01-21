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

	// normCount, prenameCount := countAccounts(accounts)

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
		authAccountErr := authAccount(acc, droptime)
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
				msg := fmt.Sprintf("%v cannot change name: %v", accID(acc), canSnipeErr)
				if acc.Type == mcgo.MsPr {
					msg = fmt.Sprintf("%v cannot create profile: %v", accID(acc), canSnipeErr)
				}
				log("error", msg)
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
		for i := 0; i < func() int {
			if acc.Type == mcgo.MsPr {
				return config.Sniper.PrenameRequests
			} else {
				return config.Sniper.SnipeRequests
			}
		}(); i++ {
			wg.Add(1)
			go func(acc *mcgo.MCaccount) {
				defer wg.Done()
				resp, err := acc.ChangeName(targetName, changeTime.Add(time.Millisecond*time.Duration(float64(totalReqCount)*config.Sniper.Spread)), acc.Type == mcgo.MsPr)

				if err != nil {
					log("error", "encountered err on nc for %v: %v", acc.Email, err.Error())
				} else {
					resps = append(resps, resp)
				}

				totalReqCount++

			}(acc)
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

	logsSlice = append(logsSlice, fmt.Sprintf("offset = %v\nlogs", offset))

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

	var (
		start = time.Now()
		end   time.Time
	)

	for _, resp := range resps {
		if resp.SendTime.Before(start) {
			start = resp.SendTime
		} else if resp.SendTime.After(end) {
			end = resp.SendTime
		}
	}

	fmt.Printf("%.3f req/ms\n", float64(len(resps))/(float64(end.UnixMicro()-start.UnixMicro())/1000))

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

func pingCommand() {
	ping, err := pingMojang()
	if err != nil {
		log("fatal", "failed to ping mojang: %v", err)
		return
	}
	log("info", "ping: %vms", ping)
}

func turbo(username string) {
	if !fileExists("accounts.txt") {
		_, err := os.Create("accounts.txt")
		if err == nil {
			log("info", "created accounts.txt, please restart the sniper once accounts are added!")
		}
	}

	if !fileExists("config.toml") {
		defaultConfig()
	}

	for {

		accStrs, _ := readLines("accounts.txt")
		accounts = loadAccSlice(accStrs)
		var authedAccounts []*mcgo.MCaccount
		config, _ := getConfig()
		var wg sync.WaitGroup
		var resps []mcgo.NameChangeReturn

		// auth + checking if ready to snipe
		for _, acc := range accounts {
			authAccountErr := authAccount(acc, time.Time{})
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
					msg := fmt.Sprintf("%v cannot change name: %v", accID(acc), canSnipeErr)
					if acc.Type == mcgo.MsPr {
						msg = fmt.Sprintf("%v cannot create profile: %v", accID(acc), canSnipeErr)
					}
					log("error", msg)
				}
			}
			time.Sleep(time.Duration(config.Accounts.AuthDelay) * time.Second)
		}

		for _, acc := range authedAccounts {
			for i := 0; i < func() int {
				if acc.Type == mcgo.MsPr {
					return config.Sniper.PrenameRequests
				} else {
					return config.Sniper.SnipeRequests
				}
			}(); i++ {
				wg.Add(1)
				go func(acc *mcgo.MCaccount) {
					defer wg.Done()
					resp, err := acc.ChangeName(username, time.Now(), acc.Type == mcgo.MsPr)

					if err != nil {
						log("error", "encountered err on nc for %v: %v", acc.Email, err.Error())
					} else {
						resps = append(resps, resp)
					}
				}(acc)
			}
			time.Sleep(time.Millisecond * 1)
		}

		wg.Wait()

		for _, resp := range resps {
			log(
				"info", "[%v] sent @ %v | recv @ %v | %v",
				prettyStatus(resp.StatusCode),
				fmtTimestamp(resp.SendTime),
				fmtTimestamp(resp.ReceiveTime),
				accID(&resp.Account),
			)

			if resp.StatusCode < 300 {

				log("success", "sniped %v onto %v", resp.Username, resp.Account.Email)
				log("info", "if you like this sniper please consider donating @ <fg=green;op=underscore>https://mcsniperpy.com/donate</>")

				if config.Announce.McsnipergoAnnounceCode != "" {
					err := announceSnipe(username, config.Announce.McsnipergoAnnounceCode, &resp.Account)

					if err != nil {
						log("error", "failed to announce snipe: %v", err)
					} else {
						log("success", "announced snipe!")
					}
				}

				if config.Announce.WebhookURL != "" {
					err := customServerAnnounce(username)
					if err != nil {
						log("error", "failed to announce snipe to your webhook: %v", err)
					} else {
						log("succes", "announced your snipe!")
					}
				}
			}
		}
		time.Sleep(time.Minute)
	}
}
