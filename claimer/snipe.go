package claimer

import (
	"errors"
	"fmt"
	"time"

	"github.com/Kqzz/MCsniperGO/log"
	"github.com/Kqzz/MCsniperGO/mc"
	"github.com/gookit/color"
)

const (
	authOffset = time.Hour * 8
	spread     = 0
)

func ClaimWithinRange(username string, dropRange mc.DropRange, accounts []*mc.MCaccount, proxies []string) error {

	fmt.Print("\n")
	log.Log("info", "sniping %s at %s", username, dropRange.Start.Format("02 Jan 06 15:04 MST"))

	for {
		if time.Until(dropRange.Start) > authOffset {
			color.Printf("\r[<fg=blue>*</>] authing in %v    ", time.Until(dropRange.Start.Add(-time.Hour*8)).Round(time.Second))
			time.Sleep(time.Second * 1)
		} else {
			color.Printf("\r[<fg=blue>*</>] starting auth...\n\n")
			break
		}
	}

	usableAccounts := []*mc.MCaccount{}

	for _, account := range accounts {
		authErr := account.MicrosoftAuthenticate("")
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
				log.Log("info", "checking license for %v, bearer: %v", account.Email, log.LastQuarter(account.Bearer))
				log.Log("err", "failed to license %v: %v", account.Email, licenseErr)
				continue
			}
			usableAccounts = append(usableAccounts, account)
		}

		if account.Type == mc.Ms {
			_, checkErr := account.NameChangeInfo()
			log.Log("info", "checking name change for %v, bearer: %v", account.Email, log.LastQuarter(account.Bearer))
			if checkErr != nil {
				log.Log("err", "failed to confirm name change for %v: %v", account.Email, checkErr)
				continue
			}
			usableAccounts = append(usableAccounts, account)
			continue
		}

		if account.Type == mc.MsPr {
			_, checkErr := account.HasGcApplied()

			log.Log("info", "checking gift code claim for %v, bearer: %v", account.Email, log.LastQuarter(account.Bearer))
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
		if time.Until(dropRange.Start) > time.Second*20 {
			color.Printf("\r[<fg=blue>*</>] sniping in %v    ", time.Until(dropRange.Start).Round(time.Second))
			time.Sleep(time.Second * 1)
		} else {
			color.Printf("\r[<fg=blue>*</>] starting snipe...\n")
			break
		}
	}

	snipe := &Claim{
		Username:  username,
		Accounts:  usableAccounts,
		DropRange: dropRange,
		Running:   true,
		Proxies:   proxies,
	}

	snipe.runClaim()

	return nil
}
