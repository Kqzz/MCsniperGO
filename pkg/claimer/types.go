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
	Statistics            Statistics

	queue    map[string]*Claim
	running  map[string]*Claim
	killChan chan bool
	workChan chan ClaimWork
	respchan chan ClaimResponse
}

type Statistics struct {
	RequestsPerSecond int
	Requests          int // total count
	FailedRequests    int
}
