package config

import (
	"github.com/spf13/pflag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sigs.k8s.io/yaml"
)

type GlobalFlags struct {
	Config         string
	VscSourcesRoot string
}

func (f *GlobalFlags) AddFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&f.Config, "config", "c",
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

func (f *GlobalFlags) ReadConfig() (*Config, error) {
	yamlFile, err := ioutil.ReadFile(f.Config)
	if err != nil {
		panic(err)
	}

	var conf Config

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return &conf, nil
}
