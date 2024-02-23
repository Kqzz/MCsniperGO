package claimer

import (
	"fmt"
	"time"

	"github.com/Kqzz/MCsniperGO/pkg/mc"
)

// Expected API: type Claim, Claim.Start() and it claims the username. Claimer pkg stores accounts, queue, and proxies

func queueClaimsWithinRange(claims []*Claim) {
	now := time.Now()
	for _, claim := range claims {
		fmt.Println(claim)
		if (claim.DropRange.Start.Before(now) && claim.DropRange.End.After(now) && !claim.Running) || claim.DropRange.Start.IsZero() { // The username is currently dropping
			// TODO
			fmt.Println("Claiming", claim.Username)
		}
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
			queueClaimsWithinRange(claimer.queue)
		}
	}
}

func (claimer *Claimer) Setup() {

	if claimer.killChan != nil {
		close(claimer.killChan)
	}
	if claimer.workChan != nil {
		close(claimer.workChan)
	}

	claimer.killChan = make(chan bool)
	claimer.workChan = make(chan ClaimWork)

	go claimer.QueueManager()
	go claimer.AuthenticationWorker()
}

func (claimer *Claimer) Queue(username string, dropRange mc.DropRange) {
	fmt.Println("queueing username")
	claimer.queue = append(claimer.queue, &Claim{Username: username, DropRange: dropRange})
}
