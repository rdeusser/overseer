package main

import (
	"os"
	"os/signal"

	"github.com/iamthemuffinman/cli"
	"github.com/iamthemuffinman/overseer/command"
)

var Commands map[string]cli.CommandFactory
var PlumbingCommands map[string]struct{}
var Ui cli.Ui

func init() {
	Ui = &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	PlumbingCommands = map[string]struct{}{
		"provision": {}, // inlcudes all subcommands
	}

	Commands = map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return &command.InitCommand{
				Ui: Ui,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Ui:       Ui,
				Revision: GitCommit,
				Version:  Version,
			}, nil
		},

		"provision": func() (cli.Command, error) {
			return &command.ProvisionCommand{
				Ui: Ui,
			}, nil
		},

		"provision virtual": func() (cli.Command, error) {
			return &command.ProvisionVirtualCommand{
				Ui:         Ui,
				ShutdownCh: makeShutdownCh(),
			}, nil
		},

		"provision physical": func() (cli.Command, error) {
			return &command.ProvisionPhysicalCommand{
				Ui:         Ui,
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
