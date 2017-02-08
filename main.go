package main

import (
	"os"

	"github.com/iamthemuffinman/cli"
	log "github.com/iamthemuffinman/logsip"

	_ "github.com/iamthemuffinman/overseer/pkg/chef"
	_ "github.com/iamthemuffinman/overseer/pkg/infoblox"
	_ "github.com/iamthemuffinman/overseer/pkg/knife"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	args := os.Args[1:]

	cli := &cli.CLI{
		Args:       args,
		Commands:   Commands,
		HelpFunc:   helpMain,
		HelpWriter: os.Stdout,
	}

	exitCode, err := cli.Run()
	if err != nil {
		log.Errorf("Error executing CLI: %s", err.Error())
	}

	return exitCode
}
