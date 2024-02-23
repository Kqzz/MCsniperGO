package claimer

import (
	"time"

	"github.com/Kqzz/MCsniperGO/pkg/log"
	"github.com/Kqzz/MCsniperGO/pkg/mc"
)

func (claimer *Claimer) AuthenticationWorker() {
	// TODO: proxies
	for {
		for _, account := range claimer.Accounts {
			err := account.MicrosoftAuthenticate("")
			if err != nil {
				log.Log("err", "failed to authenticate %v: %v", account.Email, err)
				continue
			}

			if account.Type == mc.MsPr || account.Type == mc.MsGp {
				err = account.License()
				if err != nil {
					log.Log("err", "failed to license %v: %v", account.Email, err)
					continue
				}
			}

		}
		time.Sleep(time.Hour * 24)
	}
}
