package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type GopathLinkCommand struct {
	cmd *cobra.Command
}

func NewGopathLinkCommand() *GopathLinkCommand {
	command := &GopathLinkCommand{}

	cmd := &cobra.Command{
		Use: "gopath-link",
		Short: "Create a symlink under the ${GOPATH}/src directory for a " +
			"golang module if a valid go.mod file and the GOPATH environment " +
			"variable are defined",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE:         command.Execute,
	}
	command.cmd = cmd

	return command
}

func (command *GopathLinkCommand) Register() *cobra.Command {
	return command.cmd
}

func (command *GopathLinkCommand) Execute(_ *cobra.Command, args []string) (err error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return fmt.Errorf("GOPATH is not defined")
	}

	root, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}

	modPath := path.Join(root, "go.mod")
	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		return fmt.Errorf("go.mod does not exist")
	}
	modBody, err := ioutil.ReadFile(modPath)
	if err != nil {
		return err
	}

	f, err := modfile.Parse(modPath, modBody, func(_, vers string) (string, error) {
		return vers, nil
	})
	if err != nil {
		return err
	}

	name := (*f.Module).Mod.Path
	symlinkPath := path.Join(gopath, "src", name)

	logger.Infof(
		"Creating a symling for '%s' to '%s'", root, symlinkPath)
	err = os.MkdirAll(path.Dir(symlinkPath), os.ModePerm)
	if err != nil {
		return err
	}
	return os.Symlink(root, symlinkPath)
}
