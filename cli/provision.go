package cli

import (
	"strings"

	"github.com/iamthemuffinman/cli"
)

type ProvisionCommand struct {
	Ui cli.Ui
}

func (c *ProvisionCommand) Run(args []string) int {
	return cli.RunResultHelp
}

func (c *ProvisionCommand) Help() string {
	return c.helpProvision()
}

func (c *ProvisionCommand) Synopsis() string {
	return "Provision compute infrastructure"
}

func (c *ProvisionCommand) helpProvision() string {
	helpText := `
Usage: overseer provision [SUBCOMMANDS] [OPTIONS] [HOSTS]

  Provision infrastructure on virtual or physical servers.
`
	return strings.TrimSpace(helpText)
}
