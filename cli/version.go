package cli

import (
	"fmt"

	"github.com/iamthemuffinman/cli"
)

type VersionCommand struct {
	Ui       cli.Ui
	Revision string
	Version  string
}

func (c *VersionCommand) Run(args []string) int {
	c.Ui.Output(fmt.Sprintf("%s (%s)", c.Version, c.Revision))
	return 0
}

func (c *VersionCommand) Help() string {
	return ""
}

func (c *VersionCommand) Synopsis() string {
	return "Prints the Overseer version"
}
