package claimer

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // TODO REMOVE
	"time"

	"github.com/Kqzz/MCsniperGO/pkg/log"
	"github.com/Kqzz/MCsniperGO/pkg/mc"
)

// Expected API: type Claim, Claim.Start() and it claims the username. Claimer pkg stores accounts, queue, and proxies

func (claimer *Claimer) Start(claim *Claim) {
	claim.Running = true
	claimer.running[claim.Username] = claim

}
func (claimer *Claimer) Stop(claim *Claim) {
	fmt.Println("stopping")
	claimer.running[claim.Username].Running = false
	claimer.running[claim.Username] = nil

}

func (claimer *Claimer) queueClaimsWithinRange(claims map[string]*Claim) {
	now := time.Now()
	for _, claim := range claims {
		if claim.Running {
			continue
		}

		if len(claimer.AuthenticatedAccounts) == 0 {
			log.Log("err", "no authenticated accounts")
			continue
		}

		if (claim.DropRange.Start.Before(now) && claim.DropRange.End.After(now)) || claim.DropRange.Start.IsZero() { // The username is currently dropping
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

func (claimer *Claimer) sender(ms bool, gp bool) {

	sleepTime := 15000

	if len(claimer.AuthenticatedAccounts) > 0 {
		sleepTime = 150000 / len(claimer.AuthenticatedAccounts)

		if ms {
			sleepTime = 10000 / len(claimer.AuthenticatedAccounts)
		}
	}

	sleepDuration := time.Millisecond * time.Duration(sleepTime)

	loopCount := 2
	if ms {
		loopCount = 3
	}

	for {
		select {
		case _, ok := <-claimer.killChan:
			fmt.Println("killing sender")
			if !ok {
				return
			}
		default:
			for _, claim := range claimer.running {
				for _, account := range claimer.AuthenticatedAccounts {
					if ms && account.Type != mc.Ms { // skip non ms accounts for ms
						continue
					}

					if gp && account.Type == mc.Ms { // skip ms accounts for non ms accounts
						continue
					}

					for i := 0; i < loopCount; i++ {
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

func (claimer *Claimer) SendManager() {

	go claimer.sender(true, false)
	go claimer.sender(false, true)
}

func (claimer *Claimer) Setup() {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}() // TODO rm

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
	claimer.respchan = make(chan ClaimResponse)

	fmt.Println(claimer.killChan, claimer.workChan, claimer.respchan)

	claimer.SendManager()
	go claimer.QueueManager()
	go claimer.ResponseManager()
	go claimer.AuthenticationWorker()

	claimer.queue = make(map[string]*Claim)
	claimer.running = make(map[string]*Claim)

	for _, dial := range claimer.Dialers {
		go claimer.Worker(dial)
	}
}

func (claimer *Claimer) Queue(username string, dropRange mc.DropRange) {
	// claimer.queue = append(claimer.queue, &Claim{Username: username, DropRange: dropRange})

	if claimer.queue[username] != nil {
		fmt.Println("failed to queue username")
		return
	}
	claimer.queue[username] = &Claim{Username: username, DropRange: dropRange, Claimer: claimer}
}
