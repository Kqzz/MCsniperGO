package claimer

import (
	"fmt"

	"github.com/Kqzz/MCsniperGO/pkg/mc"
	"github.com/valyala/fasthttp"
)

func (claimer *Claimer) worker(dial fasthttp.DialFunc) {
	client := &fasthttp.Client{Dial: dial}

	for {
		select {
		case claim := <-claimer.workChan:
			claim.Claim.SendRequest(claim.Account, client)
		case <-claimer.killChan:
			fmt.Println("killing worker")
			return
		}
	}
}

func (claim *Claim) SendRequest(account *mc.MCaccount, client *fasthttp.Client) {
	var statusCode int
	var failType mc.FailType
	var err error
	switch account.Type {
	case mc.MsGc:
		statusCode, failType, err = account.CreateProfile(claim.Username, client)
	case mc.Ms:
		statusCode, failType, err = account.ChangeUsername(claim.Username, client)
	}

	claim.Claimer.respchan <- ClaimResponse{statusCode, failType, err}
}
