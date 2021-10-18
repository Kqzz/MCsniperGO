package main

import (
	"os"

	"github.com/kqzz/mcgo"
	"github.com/urfave/cli/v2"
)

var accounts []*mcgo.MCaccount
var logType = make(map[string]string)

func main() {
	logType["success"] = "<fg=white>[</><fg=green;op=bold>success</><fg=white>]</>"
	logType["info"] = "<fg=white>[</><fg=cyan;op=bold>info</><fg=white>]</>"
	logType["warn"] = "<fg=white>[</><fg=yellow;op=bold>warn</><fg=white>]</>"
	logType["error"] = "<fg=white>[</><fg=red;op=bold>err</><fg=white>]</>"
	logType["fatal"] = "<fg=white>[</><fg=red;op=bold>fatal err</><fg=white>]</>"
	app := &cli.App{
		Name:  "MCsniperGO",
		Usage: "mcsnipergo",
		Action: func(c *cli.Context) error {
			snipeCommand("", -10000)
			userInput("press enter to exit")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "snipe",
				Aliases: []string{"s", "run", "start"},
				Usage:   "start MCsniperGO",
				Action: func(c *cli.Context) error {
					snipeCommand(c.String("username"), c.Float64("offset"))
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
		log(err.Error(), "fatal")
	}
}
