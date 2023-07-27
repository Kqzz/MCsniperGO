package droptime

import (
	"strconv"
	"strings"
	"time"

	"github.com/Kqzz/MCsniperGO/log"
)

type droptimeFunction func(string) (time.Time, error)

type coolkidmachoRespStruct struct {
	Unix int64 `json:"UNIX"`
}

type starShoppingResponseStruct struct {
	Unix int64 `json:"unix"`
}

type threeChar struct {
	Name string `json:"name"`
}

func GetDroptime(username string) (time.Time, time.Time) {
	for {
		rawDroptimes := log.Input("droptime range (start-end)")

		rawDroptimesSplit := strings.Split(rawDroptimes, "-")

		if len(rawDroptimesSplit) != 2 {
			log.Log("err", "invalid droptime range")
			continue
		}

		startDroptimeNum, err := strconv.Atoi(rawDroptimesSplit[0])
		if err != nil {
			log.Log("err", "invalid droptime start")
			continue
		}
		endDroptimeNum, err := strconv.Atoi(rawDroptimesSplit[1])
		if err != nil {
			log.Log("err", "invalid droptime end")
			continue
		}
		startDroptime := time.Unix(int64(startDroptimeNum), 0)
		endDroptime := time.Unix(int64(endDroptimeNum), 0)

		return startDroptime, endDroptime
	}

}
