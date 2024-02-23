package claimer

import (
	"time"
)

// Expected API: type Claim, Claim.Start() and it claims the username. Claimer pkg stores accounts, queue, and proxies

func queueClaimsWithinRange(claims []*Claim) {
	now := time.Now()
	for _, claim := range claims {
		if claim.DropRange.Start.Before(now) && claim.DropRange.End.After(now) { // The username is currently dropping
			// TODO
		}
	}

}

func (claimer *Claimer) QueueManager() {
	for {
		time.Sleep(time.Second * 1)
		queueClaimsWithinRange(claimer.Queue)
	}
}

func (claimer *Claimer) Setup() {
	go claimer.QueueManager()
	go claimer.AuthenticationWorker()
}
