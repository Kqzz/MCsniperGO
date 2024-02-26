package claimer

import (
	"fmt"
	"time"

	"github.com/Kqzz/MCsniperGO/pkg/log"
	"github.com/Kqzz/MCsniperGO/pkg/mc"
)

// Expected API: type Claim, Claim.Start() and it claims the username. Claimer pkg stores accounts, queue, and proxies

func (claimer *Claimer) Start(claim *Claim) {
	claim.Running = true
	claimer.running = append(claimer.running, claim)

}
func (claimer *Claimer) Stop(claim *Claim) {
	for i, runningClaim := range claimer.running {
		if runningClaim.Username == claim.Username {
			claimer.running = append(claimer.running[:i], claimer.running[i+1:]...)
			claim.Running = false
		}
	}

}

func (claimer *Claimer) queueClaimsWithinRange(claims []*Claim) {
	now := time.Now()
	for _, claim := range claims {
		fmt.Println(claim)

		if len(claimer.AuthenticatedAccounts) == 0 {
			log.Log("err", "no authenticated accounts")
			time.Sleep(time.Second * 20)
			continue
		}

		if (claim.DropRange.Start.Before(now) && claim.DropRange.End.After(now) && !claim.Running) || claim.DropRange.Start.IsZero() { // The username is currently dropping
			claimer.Start(claim)
		} else if claim.DropRange.End.Before(now) && claim.Running { // The username has finished dropping
			claimer.Stop(claim)
		} // TODO: stop usernames if username is claimed by other user, will involve creating checker thread
	}

}

func (claimer *Claimer) QueueManager() {
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

func (claimer *Claimer) ResponseManager() {
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

func (claimer *Claimer) sender(accountType mc.AccType) {

	sleepTime := 15000

	if len(claimer.AuthenticatedAccounts) > 0 {
		sleepTime = 150000 / len(claimer.AuthenticatedAccounts)

		if accountType == mc.Ms {
			sleepTime = 10000 / len(claimer.AuthenticatedAccounts)
		}
	}

	sleepDuration := time.Millisecond * time.Duration(sleepTime)

	loopCount := 2
	if accountType == mc.Ms {
		loopCount = 3
	}

	for {
		select {
		case _, ok := <-claimer.killChan:
			if !ok {
				return
			}
		default:
			for _, claim := range claimer.running {
				for _, account := range claimer.AuthenticatedAccounts {
					if account.Type != accountType {
						continue
					}

					for i := 0; i < loopCount; i++ {
						claimer.workChan <- ClaimWork{Claim: claim, Account: account}
						time.Sleep(sleepDuration)
					}
				}
			}
		}
	}
}

func (claimer *Claimer) SendManager() {

	go claimer.sender(mc.Ms)
	go claimer.sender(mc.MsGp)
}

func (claimer *Claimer) Setup() {

	if claimer.killChan != nil {
		close(claimer.killChan)
	}
	if claimer.workChan != nil {
		close(claimer.workChan)
	}

	if claimer.respchan != nil {
		close(claimer.respchan)
	}

	claimer.killChan = make(chan bool)
	claimer.workChan = make(chan ClaimWork)
	claimer.respchan = make(chan ClaimResponse, 1000)

	go claimer.QueueManager()
	go claimer.SendManager()
	go claimer.ResponseManager()
	go claimer.AuthenticationWorker()

	for _, dial := range claimer.Dialers {
		go claimer.StartWorker(dial)
	}
}

func (claimer *Claimer) Queue(username string, dropRange mc.DropRange) {
	fmt.Println("queueing username")
	claimer.queue = append(claimer.queue, &Claim{Username: username, DropRange: dropRange})
}
