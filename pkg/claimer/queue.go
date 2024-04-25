package claimer

import (
	"fmt"
	_ "net/http/pprof" // TODO REMOVE
	"time"

	"github.com/Kqzz/MCsniperGO/pkg/log"
	"github.com/Kqzz/MCsniperGO/pkg/mc"
)

func (claimer *Claimer) start(claim *Claim) {
	claim.Running = true
	claimer.running[claim.Username] = claim

}
func (claimer *Claimer) stop(claim *Claim) {
	log.Log("info", "stopping %s", claim.Username)

	_, e := claimer.running[claim.Username]
	if !e {
		return
	}

	claimer.running[claim.Username].Running = false
	claimer.running[claim.Username] = nil

}

func (claimer *Claimer) queueClaimsWithinRange(claims map[string]*Claim) {
	now := time.Now()
	for _, claim := range claims {

		if claim.Running && !claim.DropRange.End.IsZero() && claim.DropRange.End.Before(now) {
			claimer.stop(claim)
			continue
		}

		if claim.Running {
			continue
		}

		if len(claimer.AuthenticatedAccounts) == 0 {
			log.Log("err", "no authenticated accounts")
			continue
		}

		if (claim.DropRange.Start.Before(now) && claim.DropRange.End.After(now)) || claim.DropRange.Start.IsZero() { // The username is currently dropping
			claimer.start(claim)
		}
		// TODO: stop usernames if username is claimed by other user, will involve creating checker thread
	}

}

func (claimer *Claimer) queueManager() {
	for {
		select {
		case _, ok := <-claimer.killChan:
			if !ok {
				return
			}
		default:
			time.Sleep(time.Second * 1)
			claimer.queueClaimsWithinRange(claimer.queue)
		}
	}
}

func (claimer *Claimer) responseManager() {
	for {
		select {
		case _, ok := <-claimer.killChan:
			if !ok {
				return
			}
		case response, ok := <-claimer.respchan:
			if !ok {
				return
			}
			log.Log("info", "response: %v", response)
		}
	}
}

func determineSleep(accType mc.AccType, accountCount int) time.Duration {
	sleepTime := 15000

	if accountCount > 0 {
		sleepTime = 15000 / accountCount

		if accType == mc.Ms {
			sleepTime = 10000 / accountCount
		}
	}

	sleepDuration := time.Millisecond * time.Duration(sleepTime)
	return sleepDuration
}

func (claimer *Claimer) sender(accType mc.AccType) {

	loopCount := 2 // # of requests per account per loop
	if accType == mc.Ms {
		loopCount = 3
	}

	var accounts []*mc.MCaccount
	var sleepDuration time.Duration

	go func() {
		time.Sleep(time.Microsecond * 200)
		for {
			accounts = filter(claimer.AuthenticatedAccounts, func(acc *mc.MCaccount) bool { return acc.Type == accType })
			sleepDuration = determineSleep(accType, len(accounts)) // time between each request send
			time.Sleep(time.Second * 15)
		}
	}()

	for {
		select {
		case _, ok := <-claimer.killChan:
			fmt.Println("killing sender")
			if !ok {
				return
			}
		default:
			for _, claim := range claimer.running {
				for _, account := range accounts {
					for i := 0; i < loopCount; i++ {
						if claim == nil {
							continue
						}
						claimer.workChan <- ClaimWork{Claim: claim, Account: account}
						if i != loopCount-1 {
							time.Sleep(sleepDuration)
						}
					}
				}
			}
			time.Sleep(sleepDuration)
		}
	}
}
