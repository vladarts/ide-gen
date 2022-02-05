package main

import (
	"github.com/spf13/cobra"
	"github.com/xxxbobrxxx/ide-gen/pkg/repository"
)

type GitCloneCommand struct {
	repository.SourcesRootFlags

	cmd *cobra.Command
}

func NewGitCloneCommand() *GitCloneCommand {
	command := &GitCloneCommand{}

	cmd := &cobra.Command{
		Use:          "git-clone",
		Short:        "Clone a single git repository following the application rules",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE:         command.Execute,
	}
	command.cmd = cmd

	command.SourcesRootFlags.AddFlags(cmd.PersistentFlags())

	return command
}

func (command *GitCloneCommand) Register() *cobra.Command {
	return command.cmd
}

func (command *GitCloneCommand) Execute(_ *cobra.Command, args []string) (err error) {
	url := args[0]

	config := repository.GitSourcesRootConfig{Url: url}
	d, err := config.Directory(command.VscSourcesRoot)
	if err != nil {
		return err
	}

	commander := repository.GitSourcesRootCommander{Config: config}

	exists, err := commander.Exists(d)
	if err != nil {
		return err
	}

	if exists {
		logger.Infof(
			"Skip clone '%s' to '%s'", url, d)
	} else {
		logger.Infof(
			"Clone '%s' to '%s'", url, d)
		err = commander.Clone(d)
		if err != nil {
			return err
		}
	}

	return nil
}
