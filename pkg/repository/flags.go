package repository

import (
	"github.com/spf13/pflag"
	"os"
	"path"
)

type SourcesRootFlags struct {
	VscSourcesRoot string
}

func (f *SourcesRootFlags) AddFlags(flags *pflag.FlagSet) {
	var vscReposRootDefault string
	home, err := os.UserHomeDir()
	if err != nil {
		vscReposRootDefault = ""
	} else {
		vscReposRootDefault = path.Join(home, "dev")
	}
	flags.StringVarP(&f.VscSourcesRoot, "vcs-sources-root", "s",
		vscReposRootDefault, "")
}
