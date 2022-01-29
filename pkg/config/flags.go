package config

import (
	"github.com/spf13/pflag"
	"os"
	"path"
)

type GlobalFlags struct {
	Config          string
	IdeaProjectRoot string
	VscSourcesRoot  string
}

func (f *GlobalFlags) AddFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&f.Config, "config", "c",
		"", "")

	flags.StringVarP(&f.IdeaProjectRoot, "idea-project-root", "i",
		"", "")

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
