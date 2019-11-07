package main

import (
	"log"
	"os"

	"github.com/catsby/vaultstats/commands"
	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("vaultstats", "0.0.1")
	c.Args = os.Args[1:]

	ui := &cli.ColoredUi{
		OutputColor: cli.UiColorNone,
		InfoColor:   cli.UiColorNone,
		ErrorColor:  cli.UiColorRed,
		WarnColor:   cli.UiColorYellow,

		Ui: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}

	c.Commands = map[string]cli.CommandFactory{
		"stats": func() (cli.Command, error) {
			return &commands.StatsCommand{
				UI: ui,
			}, nil
		},
		"csv": func() (cli.Command, error) {
			return &commands.CSVCommand{
				UI: ui,
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
