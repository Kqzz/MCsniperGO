package legacy

import "time"

var queue []*Claim

func AddClaim(c *Claim) {
	queue = append(queue, c)
}

func RemoveClaim(c *Claim) {
	for i, claim := range queue {
		if claim.Username == c.Username {
			queue = append(queue[:i], queue[i+1:]...)
		}
	}
}

func GetQueue() []*Claim {
	return queue
}

func QueueManager() {
	for {
		now := time.Now()
		if len(queue) == 0 {
			continue
		}

		for _, claim := range queue {
			if claim.DropRange.Start.Before(now) && claim.DropRange.End.After(now) && !claim.Running {
				claim.Start()
			}
			if claim.DropRange.End.Before(now) {
				claim.Stop()
			}
		}

		time.Sleep(time.Second * 1)
	}
}
