package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/hcl/parser"
	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/iamthemuffinman/cli"
	log "github.com/iamthemuffinman/logsip"
	"github.com/mitchellh/go-homedir"
)

type InitCommand struct {
	Ui         cli.Ui
	ShutdownCh <-chan struct{}
}

const overseerConfigTemplate = `
hammer {
    username = "admin"
    password = "datpass"
}

knife {
    username = "admin"
    password = "datpass"
}

infoblox {
    username = "admin"
    password = "datpass"
}
`

func (c *InitCommand) Run(args []string) int {
	home, err := homedir.Dir()
	if err != nil {
		log.Errorf("couldn't get the current user's homedir: %s", err)
		return 1
	}

	if err := os.MkdirAll(fmt.Sprintf("%s/.overseer", home), 0644); err != nil {
		log.Errorf("couldn't create .overseer in your home directory: %s", err)
		return 1
	}

	overseerConfig, err := os.Create(fmt.Sprintf("%s/.overseer/overseer.conf", home))
	if err != nil {
		log.Errorf("wasn't able to create the file: %s", err)
		return 1
	}

	ast, err := parser.Parse([]byte(overseerConfigTemplate))
	if err != nil {
		log.Errorf("error parsing config: %s", err)
		return 1
	}

	if err := printer.Fprint(overseerConfig, ast); err != nil {
		log.Errorf("error writing config: %s", err)
		return 1
	}

	return 0
}

func (c *InitCommand) Help() string {
	return c.helpInit()
}

func (c *InitCommand) Synopsis() string {
	return "Initialize Overseer"
}

func (c *InitCommand) helpInit() string {
	helpText := `
Usage: overseer init

  Initialize overseer in your home directory "~/.overseer/overseer.conf". This config
  file will contain your usernames and passwords to various parts of the
  infrastructure. Make sure you keep it safe! There are sensible defaults in
  place. Normally, I'd consider it bad practice to keep passwords in a file
  somewhere on the filesystem, but A: we do it with our knife.rb file anyway and
  B: we don't have something like Vault setup where we could "lease" the secrets
  for the life of the overseer process and then revoke them at the end. Inputting
  the passwords manually is a hassle, error-prone, and just barbaric. This way
  is best for the time being.
`

	return strings.TrimSpace(helpText)
}
