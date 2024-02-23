package claimer

import (
	"github.com/Kqzz/MCsniperGO/pkg/mc"
	"github.com/valyala/fasthttp"
)

type ClaimWork struct {
	Claim   *Claim
	Account *mc.MCaccount
}

type Claim struct {
	Username  string
	DropRange mc.DropRange
	Claimer   Claimer
}

type Claimer struct {
	Dialers               []*fasthttp.DialFunc
	Accounts              []*mc.MCaccount
	AuthenticatedAccounts []*mc.MCaccount
	Queue                 []*Claim
	killChan              chan bool
	workChan              chan ClaimWork
}
