package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/iamthemuffinman/cli"
	log "github.com/iamthemuffinman/logsip"
)

func helpMain(commands map[string]cli.CommandFactory) string {
	mainCommands := make(map[string]cli.CommandFactory)

	maxKeyLen := 0
	for key, f := range commands {
		if len(key) > maxKeyLen {
			maxKeyLen = len(key)
		}

		mainCommands[key] = f
	}

	helpText := fmt.Sprintf(`
Usage: overseer [--version] [--help] <command> [args]

overseer is an automation tool for provisioning of physical and virtual servers.
It will grow to encompass much more.

Commands:
%s
`, listCommands(mainCommands, maxKeyLen))

	return strings.TrimSpace(helpText)
}

func listCommands(commands map[string]cli.CommandFactory, maxKeyLen int) string {
	var buf bytes.Buffer

	// Get the list of keys so they can be sorted and also get the maximum
	// key length so they can be aligned.
	keys := make([]string, 0, len(commands))
	for key := range commands {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		commandFunc, ok := commands[key]
		if !ok {
			// Shouldn't happen since we just built the list, but you never
			// know with these bits flyin' around like they do...
			panic("command not found: " + key)
		}

		command, err := commandFunc()
		if err != nil {
			log.Errorf("Command '%s', failed to load: %s", key, err)
			continue
		}

		key = fmt.Sprintf("%s%s", key, strings.Repeat(" ", maxKeyLen-len(key)))
		buf.WriteString(fmt.Sprintf("    %s    %s\n", key, command.Synopsis()))
	}

	return buf.String()
}
