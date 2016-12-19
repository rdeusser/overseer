// Copyright © 2016 Robert Deusser <robert.deusser@linuxadmins.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"

	"github.com/iamthemuffinman/cli"
	log "github.com/iamthemuffinman/logsip"
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
