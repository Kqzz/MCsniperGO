package main

import (
	"os"

	"github.com/gookit/color"
	"github.com/kqzz/mcgo"
	"github.com/urfave/cli/v2"
)

var accounts []*mcgo.MCaccount

func main() {
	app := &cli.App{
		Name:  "MCsniperGO",
		Usage: "mcsnipergo",
		Action: func(c *cli.Context) error {
			color.Printf(genHeader())
			err := snipeCommand("", -10000)
			if err != nil {
				log("fatal", err.Error())
			}

			userInput("press enter to exit")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "snipe",
				Aliases: []string{"s", "run", "start"},
				Usage:   "start MCsniperGO",
				Action: func(c *cli.Context) error {
					color.Printf(genHeader())
					err := snipeCommand(c.String("username"), c.Float64("offset"))
					if err != nil {
						log("fatal", err.Error())
					}
					userInput("press enter to exit")
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "username",
						Aliases: []string{"u", "user", "name"},
						Usage:   "username to snipe",
						Value:   "",
					},
					&cli.Float64Flag{
						Name:    "offset",
						Aliases: []string{"o", "delay", "time-offset"},
						Usage:   "snipe x ms early",
						Value:   -10000,
					},
				},
			},
			{
				Name: 		"queue",
				Aliases: 	[]string{"q", "list", "snipelist", "snipequeue"},
				Usage:		"queue list of names to auto snipe",
				Action: 	func(c *cli.Context) error {
					color.Printf(genHeader())
					snipeList()
					return nil
				},
			},
			{
				Name:    "autosnipe",
				Aliases: []string{"as", "auto"},
				Usage:   "Auto-snipe 3 character names",
				Action: func(c *cli.Context) error {
					color.Printf(genHeader())
					err := autoSnipeCommand(c.Float64("offset"))
					if err != nil {
						log("fatal", err.Error())
					}
					userInput("press enter to exit")
					return nil
				},
				Flags: []cli.Flag{
					&cli.Float64Flag{
						Name:    "offset",
						Aliases: []string{"o", "delay", "time-offset"},
						Usage:   "snipe x ms early",
						Value:   -1000,
					},
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
		log("fatal", err.Error())
	}
}