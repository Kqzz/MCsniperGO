package claimer

import (
	"fmt"

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
