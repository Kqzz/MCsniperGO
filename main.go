package main

import (
	"os"

	"github.com/kqzz/mcgo"
	"github.com/urfave/cli/v2"
)

var accounts []*mcgo.MCaccount
var presets = make(map[int]string)

func main() {
	presets[0] = "<fg=white>[</><fg=cyan;op=bold>info</><fg=white>]</>"     //Info
	presets[1] = "<fg=white>[</><fg=green;op=bold>success</><fg=white>]</>" //Succcess
	presets[2] = "<fg=white>[</><fg=red;op=bold>err</><fg=white>]</>"       //Error
	presets[3] = "<fg=white>[</><fg=yellow;op=bold>warn</><fg=white>]</>"   //Warn
	presets[4] = "<fg=white>[</><fg=red;op=bold>fatal err</><fg=white>]</>" //Fatal
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
		log(err.Error(), 2)
	}
}
