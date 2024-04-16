package claimer

import (
	"fmt"

	"github.com/Kqzz/MCsniperGO/pkg/mc"
)

// contains the external API for the claimer, all backend logic will be handled elsewhere
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
	claimer.respchan = make(chan ClaimResponse)

	// senders are responsible for putting username requests onto the queue and timing everything
	go claimer.sender(mc.Ms)
	go claimer.sender(mc.MsGc)

	go claimer.queueManager()         //
	go claimer.responseManager()      // reads responses and prints (TODO callback function)
	go claimer.authenticationWorker() // authenticates accounts and keeps the accounts slice up to date

	claimer.queue = make(map[string]*Claim)
	claimer.running = make(map[string]*Claim)

	for _, dial := range claimer.Dialers {
		go claimer.worker(dial)
	}
}
func (claimer *Claimer) Queue(username string, dropRange mc.DropRange) error {
	// claimer.queue = append(claimer.queue, &Claim{Username: username, DropRange: dropRange})

	if claimer.queue[username] != nil {
		return fmt.Errorf("%s is already in queue", username)
	}
	claimer.queue[username] = &Claim{Username: username, DropRange: dropRange, Claimer: claimer}
	return nil
}

func (claimer *Claimer) Dequeue(username string) {
	// TODO
}

func (claimer *Claimer) Shutdown() {
	// TODO
}
