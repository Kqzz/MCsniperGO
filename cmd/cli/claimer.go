package main

import (
	"crypto/tls"
	"strings"
	"time"

	"github.com/Kqzz/MCsniperGO/mc"

	"github.com/Kqzz/MCsniperGO/log"
)

var workerCount = 100

type Snipe struct {
	Username          string
	Spread            float64
	Running           bool
	Droptime          time.Time
	DroptimeEnd       time.Time
	Accounts          []*mc.MCaccount
	ApiKey            string
	SkinChangeUrl     string
	SkinChangeVariant string
	Proxy             string
}

func (s *Snipe) Start() {
	s.Running = true
	go s.runClaim()
}

func (s *Snipe) Stop() {
	s.Running = false
}

func PingMojang() (float64, error) {
	payload := "PUT /minecraft/profile/name/test HTTP/1.1\r\nHost: api.minecraftservices.com\r\nAuthorization: Bearer BEARER\r\n"
	conn, err := tls.Dial("tcp", "api.minecraftservices.com:443", nil)
	if err != nil {
		return 0, err
	}
	var sumNanos int64
	for i := 0; i < 3; i++ {
		junk := make([]byte, 4096)
		conn.Write([]byte(payload))
		time1 := time.Now()
		conn.Write([]byte("\r\n"))
		conn.Read(junk)
		duration := time.Since(time1)
		sumNanos += duration.Nanoseconds()
	}
	conn.Close()
	sumNanos /= 3
	avgMillis := float64(sumNanos) / float64(1000000)
	return avgMillis, nil

}

type Claim struct {
	Name    string
	Bearer  string
	AccType mc.AccType
	AccNum  int
	Proxy   string
}

func generator(workChan chan Claim, killChan chan bool, bearers []string, name string, accType mc.AccType, endTime time.Time, proxies []string, delay int) {
	noEnd := endTime.IsZero()
	if len(bearers) == 0 {
		return
	}

	sleepTime := delay

	if delay == -1 {
		sleepTime = 15500 / len(bearers)
		if accType == mc.Ms {
			sleepTime = 10500 / len(bearers)
		}
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

			workChan <- Claim{
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

func claimName(claim Claim) {
	acc := mc.MCaccount{
		Bearer: claim.Bearer,
		Type:   claim.AccType,
	}

	resp, err := acc.ChangeName(claim.Name, time.Now(), acc.Type != mc.Ms, claim.Proxy)
	if err != nil {
		log.Log("err", err.Error())
		log.Log("err", "Proxy: "+claim.Proxy)
	}

	log.Log("info", "%v %vms [%v] %v %v #%d", resp.ReceiveTime.Format("15:04:05.999"), resp.ReceiveTime.Sub(resp.SendTime).Milliseconds(), claim.Name, resp.StatusCode, acc.Type, claim.AccNum)
	if resp.StatusCode == 200 {
		log.Log("success", "Claimed %v on %v acc, %v", claim.Name, acc.Type, acc.Bearer[len(acc.Bearer)/2:])
	}

}

func worker(claimChan chan Claim, killChan chan bool) {
	for {
		select {
		case claim := <-claimChan:
			claimName(claim)
		case <-killChan:
			return
		}
	}
}

func (s *Snipe) runClaim() {
	workChan := make(chan Claim)
	killChan := make(chan bool)

	go func() {
		for {
			select {
			case <-killChan:
				return
			default:
				time.Sleep(time.Second * 15)
			}
		}
	}()

	go func() {
		for {
			if !s.Running {
				log.Log("info", "Stopped claim of %v", s.Username)
				close(killChan)
				return
			}
			time.Sleep(time.Second * 5)
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
	log.Log("info", "using %v proxies", len(strings.Split(s.Proxy, ",")))

	time.Sleep(time.Until(s.Droptime))

	go generator(workChan, killChan, gcs, s.Username, mc.MsPr, s.DroptimeEnd, strings.Split(s.Proxy, ","), -1)
	go generator(workChan, killChan, mss, s.Username, mc.Ms, s.DroptimeEnd, strings.Split(s.Proxy, ","), -1)

	for time.Now().Before(s.DroptimeEnd) {
		time.Sleep(10 * time.Second)
	}
	s.Running = false
	_, ok := (<-killChan)
	if ok {
		close(killChan)
	}

}
