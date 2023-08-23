package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	droptimePkg "github.com/Kqzz/MCsniperGO/droptime"
	"github.com/Kqzz/MCsniperGO/log"
	"github.com/Kqzz/MCsniperGO/mc"
	"github.com/gookit/color"
)

const (
	authOffset = time.Hour * 8
	spread     = 0
)

func snipe(username string) error {

	proxies, err := readLines("proxies.txt")

	if err != nil {
		log.Log("err", "failed to load proxies: %v", err)
	}

	err = nil

	giftCodeLines, _ := readLines("gc.txt")
	gamepassLines, _ := readLines("gp.txt")
	microsoftLines, _ := readLines("ms.txt")

	gcs, parseErrors := parseAccounts(giftCodeLines, mc.MsPr)

	for _, er := range parseErrors {
		if er == nil {
			continue
		}
		log.Log("err", "%v", err)
	}
	microsofts, msParseErrors := parseAccounts(microsoftLines, mc.Ms)

	for _, er := range msParseErrors {
		if er == nil {
			continue
		}
		log.Log("err", "%v", err)
	}

	gamepasses, gpParseErrors := parseAccounts(gamepassLines, mc.MsGp)

	for _, er := range gpParseErrors {
		if er == nil {
			continue
		}
		log.Log("err", "%v", err)
	}

	startDroptime, endDroptime := droptimePkg.GetDroptime(username)

	if err != nil {
		log.Log("err", "%v", err)
		return errors.New("failed to parse accounts")
	}

	fmt.Print("\n")
	log.Log("info", "sniping %s at %s", username, startDroptime.Format("02 Jan 06 15:04 MST"))

	for {
		if time.Until(startDroptime) > authOffset {
			color.Printf("\r[<fg=blue>*</>] authing in %v    ", time.Until(startDroptime.Add(-time.Hour*8)).Round(time.Second))
			time.Sleep(time.Second * 1)
		} else {
			color.Printf("\r[<fg=blue>*</>] starting auth...\n\n")
			break
		}
	}

	accounts := append(gcs, microsofts...)
	accounts = append(accounts, gamepasses...)

	if len(accounts) == 0 {
		return errors.New("no accounts loaded")
	}

	usableAccounts := []*mc.MCaccount{}

	for _, account := range accounts {
		authErr := account.MicrosoftAuthenticate()
		if authErr != nil {
			log.Log("err", "failed to authenticate %v: %v", account.Email, authErr)
			time.Sleep(time.Second * 21)
			continue
		} else {
			log.Log("success", "authenticated %s", account.Email)
		}

		time.Sleep(time.Millisecond * 500)
		if account.Type == mc.MsGp {
			licenseErr := account.License()
			if licenseErr != nil {
				log.Log("err", "failed to license %v: %v", account.Email, licenseErr)
				continue
			}
			usableAccounts = append(usableAccounts, account)
		}

		if account.Type == mc.Ms {
			_, checkErr := account.NameChangeInfo()
			if checkErr != nil {
				log.Log("err", "failed to confirm name change for %v: %v", account.Email, checkErr)
				continue
			}
			usableAccounts = append(usableAccounts, account)
			continue
		}

		if account.Type == mc.MsPr {
			_, checkErr := account.HasGcApplied()

			if checkErr != nil {
				log.Log("err", "failed to confirm gift code claim for %v: %v", account.Email, checkErr)
				continue
			}

			usableAccounts = append(usableAccounts, account)
		}
		time.Sleep(time.Second * 21)

	}

	if len(usableAccounts) == 0 {
		return errors.New("no accounts successfully authenticated")
	} else {
		log.Log("success", "authenticated %d account(s)\n", len(usableAccounts))
	}

	for {
		if time.Until(startDroptime) > time.Second*20 {
			color.Printf("\r[<fg=blue>*</>] sniping in %v    ", time.Until(startDroptime).Round(time.Second))
			time.Sleep(time.Second * 1)
		} else {
			color.Printf("\r[<fg=blue>*</>] starting snipe...\n")
			break
		}
	}

	snipe := &Snipe{
		Username:    username,
		Accounts:    usableAccounts,
		Droptime:    startDroptime,
		DroptimeEnd: endDroptime,
		Running:     true,
		Proxy:       strings.Join(proxies, ","),
	}

	snipe.runClaim()

	return nil
}
