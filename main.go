package main

import (
	"os"

	"github.com/kqzz/mcgo"
	"github.com/urfave/cli/v2"
)

var accounts []*mcgo.MCaccount

func main() {
	app := &cli.App{
		Name:  "MCsniperGO",
		Usage: "mcsnipergo",
		Action: func(c *cli.Context) error {
			snipeCommand()
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "snipe",
				Aliases: []string{"s", "run", "start"},
				Usage:   "start MCsniperGO",
				Action: func(c *cli.Context) error {
					snipeCommand()
					return nil
				},
			},
			{
				Name:    "ping",
				Aliases: []string{"p"},
				Usage:   "Ping Mojang servers",
				Action: func(c *cli.Context) error {
					pingCommand()
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logFatal(err.Error())
	}
}
