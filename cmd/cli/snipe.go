package main

import (
	"errors"
	"fmt"
	"sync"
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

func snipe(username string, offset float64) error {

	accountsLines, err := readLines("accounts.txt")

	if err != nil {
		return err
	}

	accounts, errs := parseAccounts(accountsLines)

	if err != nil {
		log.Log("err", "%v", errs)
		return errors.New("failed to parse accounts")
	}

	droptime, err := droptimePkg.GetDroptime(username)

	if err != nil {
		return err
	}

	fmt.Print("\n")
	log.Log("info", "sniping %s at %s", username, droptime.Format("02 Jan 06 15:04 MST"))

	for {
		if time.Until(droptime) > authOffset {
			color.Printf("\r[<fg=blue>*</>] authing in %v    ", time.Until(droptime).Round(time.Second))
			time.Sleep(time.Second * 1)
		} else {
			color.Printf("\r[<fg=blue>*</>] starting auth...\n\n")
			break
		}
	}

	usableAccounts := []*mc.MCaccount{}

	for _, account := range accounts {
		authErr := account.MicrosoftAuthenticate()
		if authErr != nil {
			log.Log("err", "failed to authenticate %v: %v", account.Email, authErr)
			continue
		} else {
			log.Log("success", "authenticated %s", account.Email)
		}

		account.Type = mc.Ms

		ncInfo, checkErr := account.NameChangeInfo()

		if checkErr != nil {
			log.Log("err", "failed to check %v's acc: %v", account.Email, checkErr)
			continue
		}

		if ncInfo.Namechangeallowed {
			account.Type = mc.Ms
		} else {
			isGc, checkErr := account.HasGcApplied()

			if checkErr != nil {
				log.Log("err", "failed to check %v's type: %v", account.Email, checkErr)
				continue
			}

			if isGc {
				account.Type = mc.MsPr
			}
		}

		usableAccounts = append(usableAccounts, account)

	}

	if len(usableAccounts) == 0 {
		return errors.New("no accounts successfully authenticated")
	} else {
		log.Log("success", "authenticated %d account(s)\n", len(usableAccounts))
	}

	changeTime := droptime.Add(time.Millisecond * time.Duration(0-offset))

	for {
		if time.Until(changeTime) > time.Second*20 {
			color.Printf("\r[<fg=blue>*</>] sniping in %v    ", time.Until(droptime).Round(time.Second))
			time.Sleep(time.Second * 1)
		} else {
			color.Printf("\r[<fg=blue>*</>] starting snipe...\n")
			break
		}
	}

	var wg sync.WaitGroup
	var resps []mc.NameChangeReturn
	var sentReqs int

	for _, account := range usableAccounts {
		reqCount := 3
		if account.Type == mc.MsPr {
			reqCount = 6
		}

		for i := 0; i < reqCount; i++ {
			wg.Add(1)
			go func(acc *mc.MCaccount) {
				defer wg.Done()
				resp, err := acc.ChangeName(
					username,
					changeTime.Add(
						time.Millisecond*time.Duration(
							float64(sentReqs)*spread,
						),
					),
					acc.Type == mc.MsPr,
				)

				if err != nil {
					log.Log("err", "encountered err on nc for %v: %v", acc.Email, err)
				} else {
					resps = append(resps, resp)
				}

				sentReqs++
			}(account)

			time.Sleep(time.Millisecond * 2)
		}
	}

	wg.Wait()

	for _, r := range resps {
		log.Log(
			"info",
			"[%s] sent @ %s | recv @ %s | %s",
			log.PrettyStatus(r.StatusCode),
			log.FmtTimestamp(r.SendTime),
			log.FmtTimestamp(r.ReceiveTime),
			log.PrettyTimestampStatus(r.ReceiveTime, droptime, r.StatusCode),
		)

		if r.StatusCode < 300 && r.StatusCode > 199 {
			log.Log(
				"success",
				"sniped %s onto %s",
				username,
				r.Account.Email,
			)
		}
	}

	return nil
}
