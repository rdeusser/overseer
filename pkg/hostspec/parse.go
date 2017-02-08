package hostspec

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Spec struct {
	Hosts []string
	MACs  []string
}

func ParseFile(path string) (*Spec, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var spec Spec

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// We don't want to detect other spaces and exit the program
		// so I'm gonna go ahead and take care of all extraenous space
		// on both ends of the line.
		trimmedLine := strings.TrimSpace(scanner.Text())
		line := strings.Split(trimmedLine, " ")

		// We only have one field so it's probably just the hostname
		if len(line) == 1 {
			spec.Hosts = append(spec.Hosts, line[0])
		} else if len(line) == 2 {
			// If it's two fields we're handling a physical machine with a MAC address
			spec.Hosts = append(spec.Hosts, line[0])
			spec.MACs = append(spec.MACs, line[1])
		} else {
			// If there's more than two fields it must be an error
			return nil, fmt.Errorf("hostspec contains more than a server name and a MAC address: got %d lines", len(line))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// A little hacky, but this is an edge case where you forgot to add a
	// MAC address to a host while adding a MAC address to another host.
	// We don't want any mismatches here.
	if len(spec.Hosts) != len(spec.MACs) && len(spec.MACs) != 0 {
		return nil, fmt.Errorf("hostspec has an uneven number of hosts and MACs")
	}

	return &spec, nil
}
