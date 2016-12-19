package main

import (
	"os"
	"os/signal"

	hashicli "github.com/iamthemuffinman/cli"
	"github.com/iamthemuffinman/overseer/cli"
)

var Commands map[string]hashicli.CommandFactory
var PlumbingCommands map[string]struct{}
var Ui hashicli.Ui

func init() {
	Ui = &hashicli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	PlumbingCommands = map[string]struct{}{
		"provision": {}, // inlcudes all subcommands
	}

	Commands = map[string]hashicli.CommandFactory{
		"init": func() (hashicli.Command, error) {
			return &cli.InitCommand{
				Ui: Ui,
			}, nil
		},

		"version": func() (hashicli.Command, error) {
			return &cli.VersionCommand{
				Ui:       Ui,
				Revision: GitCommit,
				Version:  Version,
			}, nil
		},

		"provision": func() (hashicli.Command, error) {
			return &cli.ProvisionCommand{
				Ui: Ui,
			}, nil
		},

		"provision virtual": func() (hashicli.Command, error) {
			return &cli.ProvisionVirtualCommand{
				Ui:         Ui,
				ShutdownCh: makeShutdownCh(),
			}, nil
		},

		"provision physical": func() (hashicli.Command, error) {
			return &cli.ProvisionPhysicalCommand{
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
