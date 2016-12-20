package main

import (
	"os"
	"os/signal"

	"github.com/iamthemuffinman/cli"
	"github.com/iamthemuffinman/overseer/cmd"
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
			return &cmd.InitCommand{
				Ui: Ui,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &cmd.VersionCommand{
				Ui:       Ui,
				Revision: GitCommit,
				Version:  Version,
			}, nil
		},

		"provision": func() (cli.Command, error) {
			return &cmd.ProvisionCommand{
				Ui: Ui,
			}, nil
		},

		"provision virtual": func() (cli.Command, error) {
			return &cmd.ProvisionVirtualCommand{
				Ui:         Ui,
				ShutdownCh: makeShutdownCh(),
			}, nil
		},

		"provision physical": func() (cli.Command, error) {
			return &cmd.ProvisionPhysicalCommand{
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
