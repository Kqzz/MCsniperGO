package claimer

import (
	"math"
	"strings"
	"time"

	"github.com/Kqzz/MCsniperGO/pkg/mc"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"

	"github.com/Kqzz/MCsniperGO/log"
)

var workerCount = 100

type Claim struct {
	Username  string
	Running   bool
	DropRange mc.DropRange
	Accounts  []*mc.MCaccount
	Proxies   []string
}

func (c *Claim) Start() {
	c.Running = true
	go c.runClaim()
}

func (c *Claim) Stop() {
	c.Running = false
}

type ClaimAttempt struct {
	Claim   *Claim
	Name    string
	Bearer  string
	AccType mc.AccType
	AccNum  int
	Proxy   string
}

func requestGenerator(
	workChan chan ClaimAttempt,
	killChan chan bool,
	bearers []string,
	name string,
	accType mc.AccType,
	endTime time.Time,
	proxies []string,
	delay int,
) {
	noEnd := endTime.IsZero()
	if len(bearers) == 0 {
		return
	}

	sleepTime := delay

	if endTime.IsZero() {
		// catches infinite range and gives it the max delay, more information on the ratelimit below
		nMax := int(math.Min(float64(len(bearers)), float64(len(proxies))))
		day := int((time.Hour * 24).Milliseconds())
		sleepTime = int(day / 40 / nMax)
	} else if delay == -1 {
		// prevents overrequesting in cases where there are large discrepencies in account/proxy counts
		nMax := int(math.Min(float64(len(bearers)), float64(len(proxies))))
		/*
			ratelimits are set up as follows - 3 requests / 30s (2/30s for giftcard accounts), 40 requests / 24h.
			we compute both. this is technically not maximally performant - you need around 9.25 octillion accounts
			(which is a number 200 quadrillion times larger than every possible ipv4 address) for it to become relevant
			assuming a constant request stream
		*/
		day := int((time.Hour * 24).Milliseconds())
		// if under ratelimit periods for our drop range, we should use the drop range instead of the ratelimit period
		shortInterval := int(math.Min(30000, float64(time.Until(endTime).Milliseconds())))
		longInterval := int(math.Min(float64(day), float64(time.Until(endTime).Milliseconds())))
		var deltaShort int
		if accType == mc.Ms {
			deltaShort = shortInterval / 3 / nMax
		} else {
			deltaShort = shortInterval / 2 / nMax
		}
		deltaLong := longInterval / 40 / nMax
		// take the higher of the two
		sleepTime = int(math.Max(float64(deltaShort), float64(deltaLong)))
	}
	loopCount := 2
	if accType == mc.Ms {
		loopCount = 3
	}
	i := 0
	prox := 0
	for noEnd || time.Now().Before(endTime) {
		for y := 0; y < loopCount; y++ { // run n times / bearer
			if i >= len(bearers) {
				i = 0
			}

			if prox >= len(proxies) {
				prox = 0
			}

			workChan <- ClaimAttempt{
				Name:    name,
				Bearer:  bearers[i],
				AccType: accType,
				Proxy:   proxies[prox],
				AccNum:  i + 1,
			}
			time.Sleep(time.Millisecond * time.Duration(sleepTime))
			prox++
		}
		i++
	}

}

func claimName(claim ClaimAttempt, client *fasthttp.Client) {
	acc := mc.MCaccount{
		Bearer: claim.Bearer,
		Type:   claim.AccType,
	}

	status := 0
	var err error = nil
	var fail mc.FailType = mc.DUPLICATE

	if strings.HasPrefix(claim.Proxy, "socks5://") {
		client.Dial = fasthttpproxy.FasthttpSocksDialer(claim.Proxy)
	} else if claim.Proxy != "" {
		client.Dial = fasthttpproxy.FasthttpHTTPDialer(claim.Proxy)
	}

	before := time.Now()
	if claim.AccType == mc.Ms {
		status, fail, err = acc.ChangeUsername(claim.Name, client)
	} else {
		status, fail, err = acc.CreateProfile(claim.Name, client)
	}
	after := time.Now()

	if err != nil {
		log.Log("err", "%v #%d", err, claim.AccNum)
		return
	}

	Stats.Total++

	log.Log("info", "[%v] %v %vms %v %v #%d | %s", claim.Name, after.Format("15:04:05.999"), after.Sub(before).Milliseconds(), log.PrettyStatus(status), acc.Type, claim.AccNum, string(fail))
	if status == 200 {
		log.Log("success", "Claimed %v on %v acc, %v", claim.Name, acc.Type, acc.Bearer[len(acc.Bearer)/2:])
		log.Log("success", "Join https://discord.gg/2BZseKW for more!")
		Stats.Success++
		claim.Claim.Running = false
	}

	switch fail {
	case mc.DUPLICATE:
		Stats.Duplicate++
	case mc.NOT_ALLOWED:
		Stats.NotAllowed++
	case mc.TOO_MANY_REQUESTS:
		Stats.TooManyRequests++
	}

}

func worker(claimChan chan ClaimAttempt, killChan chan bool) {
	client := &fasthttp.Client{
		Dial: fasthttp.Dial,
	}
	for {
		select {
		case claim := <-claimChan:
			claimName(claim, client)
		case <-killChan:
			return
		}
	}
}

func (s *Claim) runClaim() {
	workChan := make(chan ClaimAttempt)
	killChan := make(chan bool)
	s.Running = true

	go func() {

		doChecks := true
		_, statusCode, err := mc.UsernameToUuid(s.Username)

		if err != nil {
			log.Log("err", "failed to get uuid of %v for availability checking: %v", s.Username, err)
		}

		if statusCode != 404 {
			doChecks = false
		}

		for i := 0; true; i++ {
			if i%30 == 0 && doChecks {
				i = 0
				_, statusCode, err = mc.UsernameToUuid(s.Username)

				if err != nil {
					log.Log("err", "failed to get uuid of %v for availability checking: %v", s.Username, err)
				}

				if statusCode == 200 {
					log.Log("err", "username %v is taken now", s.Username)
					s.Running = false
					close(killChan)
					return
				}
			}

			if !s.Running {
				log.Log("info", "Stopped claim of %v", s.Username)
				close(killChan)
				return
			}
			time.Sleep(time.Second * 2)
		}
	}()

	gcs := []string{}
	mss := []string{}

	for _, acc := range s.Accounts {
		if acc.Type == mc.Ms {
			mss = append(mss, acc.Bearer)
		} else {
			gcs = append(gcs, acc.Bearer)
		}
	}

	for i := 0; i < workerCount; i++ {
		go worker(workChan, killChan)
	}

	log.Log("info", "using %v accounts", len(s.Accounts))
	log.Log("info", "using %v proxies", len(s.Proxies))

	if len(s.Proxies) == 0 {
		s.Proxies = []string{""}
	}

	time.Sleep(time.Until(s.DropRange.Start))

	go requestGenerator(workChan, killChan, gcs, s.Username, mc.MsPr, s.DropRange.End, s.Proxies, -1)
	go requestGenerator(workChan, killChan, mss, s.Username, mc.Ms, s.DropRange.End, s.Proxies, -1)

	if s.DropRange.End.IsZero() {
		select {}
	}

	for time.Now().Before(s.DropRange.End) {
		time.Sleep(10 * time.Second)
	}
	s.Running = false
	_, ok := (<-killChan)
	if ok {
		close(killChan)
	}

}
