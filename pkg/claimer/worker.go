package claimer

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func (claimer *Claimer) StartWorker(dial fasthttp.DialFunc) {
	client := &fasthttp.Client{Dial: dial}

	for {
		select {
		case claim := <-claimer.workChan:
			fmt.Println(claim, client)
			claim.Claim.SendRequest(claim.Account, client)
			// TODO RUN CLAIM
		case <-claimer.killChan:
			return
		}
	}

}
