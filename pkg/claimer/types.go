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
	Running   bool
	Claimer   *Claimer
}

type ClaimResponse struct {
	StatusCode int
	FailType   mc.FailType
	Error      error
}

type Claimer struct {
	Dialers               []fasthttp.DialFunc
	Accounts              []*mc.MCaccount
	AuthenticatedAccounts []*mc.MCaccount
	queue                 []*Claim
	running               []*Claim
	killChan              chan bool
	workChan              chan ClaimWork
	respchan              chan ClaimResponse
}
