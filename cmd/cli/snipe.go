package main

import (
	"fmt"
	"time"

	droptimePkg "github.com/Kqzz/MCsniperGO/droptime"
	"github.com/Kqzz/MCsniperGO/log"
	"github.com/gookit/color"
)

func snipe(username string, offset float64) error {
	droptime, err := droptimePkg.GetDroptime(username)

	if err != nil {
		return err
	}

	fmt.Print("\n")
	log.Log("info", "sniping %s at %s\n\n", username, droptime.Format("02 Jan 06 15:04 MST"))

	for {
		color.Printf("\r[<fg=blue>info</>] sniping in %v        ", time.Until(droptime).Round(time.Second))
		time.Sleep(time.Second * 1)
	}

	return nil
}
