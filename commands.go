package main

import (
	"os"
	"os/signal"

	"github.com/iamthemuffinman/cli"
	"github.com/iamthemuffinman/overseer/cmd"
)

var Commands map[string]cli.CommandFactory
var PlumbingCommands map[string]struct{}
var UI cli.Ui

func init() {
	UI = &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	PlumbingCommands = map[string]struct{}{
		"provision": {}, // includes all subcommands
	}

	Commands = map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return &cmd.InitCommand{
				UI: UI,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &cmd.VersionCommand{
				UI:       UI,
				Revision: GitCommit,
				Version:  Version,
			}, nil
		},

		"provision": func() (cli.Command, error) {
			return &cmd.ProvisionCommand{
				UI: UI,
			}, nil
		},

		"provision virtual": func() (cli.Command, error) {
			return &cmd.ProvisionVirtualCommand{
				UI:         UI,
				ShutdownCh: makeShutdownCh(),
			}, nil
		},

		"provision physical": func() (cli.Command, error) {
			return &cmd.ProvisionPhysicalCommand{
				UI:         UI,
				ShutdownCh: makeShutdownCh(),
			}, nil
		},
	}
}

func makeShutdownCh() <-chan struct{} {
	resultCh := make(chan struct{})

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		for {
			<-signalCh
			resultCh <- struct{}{}

		}
	}()

	return resultCh
}
