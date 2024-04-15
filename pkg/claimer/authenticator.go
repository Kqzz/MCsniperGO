package claimer

import (
	"slices"
	"time"

	"github.com/Kqzz/MCsniperGO/pkg/log"
	"github.com/Kqzz/MCsniperGO/pkg/mc"
)

const authPause = 60 * 30 // check every 30 minutes

func (claimer *Claimer) authenticationWorker() {
	// TODO: proxies
	for {
		select {
		case _, ok := <-claimer.killChan:
			if !ok {
				return
			}
		default:
			for _, account := range claimer.Accounts {

				if time.Until(account.AuthenticationExpiration()) > time.Hour*2 {
					if account.BearerAccount {
						if !slices.Contains(claimer.AuthenticatedAccounts, account) {
							claimer.AuthenticatedAccounts = append(claimer.AuthenticatedAccounts, account)
						}
					}
					log.Log("info", "%v authenticated, expires in %v", account.Email, time.Until(account.AuthenticationExpiration()))
					continue
				}

				err := account.MicrosoftAuthenticate("")
				if err != nil {
					log.Log("err", "failed to authenticate %v: %v", account.Email, err)
					continue
				}

				if account.Type == mc.MsGc {
					err = account.License()
					if err != nil {
						log.Log("err", "failed to license %v: %v", account.Email, err)
						continue
					}
				}

				// Remove the account from the authenticated accounts list if it's already there

				for i, authenticatedAccount := range claimer.AuthenticatedAccounts {
					if authenticatedAccount.Email == account.Email {
						claimer.AuthenticatedAccounts = append(claimer.AuthenticatedAccounts[:i], claimer.AuthenticatedAccounts[i+1:]...)
					}
				}

				claimer.AuthenticatedAccounts = append(claimer.AuthenticatedAccounts, account)

				log.Log("success", "authenticated %v", account.Email)

			}
			time.Sleep(time.Second * authPause)
		}
	}
}
