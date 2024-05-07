package claimer

import (
	"fmt"

	"github.com/Kqzz/MCsniperGO/pkg/log"
	"github.com/Kqzz/MCsniperGO/pkg/mc"
	"github.com/valyala/fasthttp"
)

func (claimer *Claimer) worker(dial fasthttp.DialFunc) {
	client := &fasthttp.Client{Dial: dial}

	for {
		select {
		case claim := <-claimer.workChan:
			claim.Claim.SendRequest(claim.Account, client, claimer.respchan)
		case <-claimer.killChan:
			fmt.Println("killing worker")
			return
		}
	}
}

func (claim *Claim) SendRequest(account *mc.MCaccount, client *fasthttp.Client, respChan chan ClaimResponse) {

	if claim == nil {
		log.Log("warn", "[debug] claim == nil")
		return
	}

	var statusCode int
	var failType mc.FailType
	var err error
	switch account.Type {
	case mc.MsGc:
		statusCode, failType, err = account.CreateProfile(claim.Username, client)
	case mc.Ms:
		statusCode, failType, err = account.ChangeUsername(claim.Username, client)
	}

	respChan <- ClaimResponse{statusCode, failType, err}
}
