package knife

import (
	"github.com/iamthemuffinman/overseer/pkg/buildspec"
)

type Knife struct {
	Hostname    string
	Environment string
	RunList     []string
}

func New(bspec *buildspec.Spec) *Knife {
	return &Knife{
		Hostname:    "",
		Environment: bspec.Chef.Environment,
		RunList:     bspec.Chef.RunList,
	}
}
