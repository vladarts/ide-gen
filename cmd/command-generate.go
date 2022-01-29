package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xxxbobrxxx/idea-project-manager/pkg/config"
)

type GenerateCommand struct {
	config.GlobalFlags

	cmd *cobra.Command
}

func NewGenerateCommand() *GenerateCommand {
	command := &GenerateCommand{}

	cmd := &cobra.Command{
		Use:          "generate",
		Aliases:      []string{"gen"},
		Short:        "Generate an IDEA project",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		RunE:         command.Execute,
	}
	command.cmd = cmd

	command.GlobalFlags.AddFlags(cmd.PersistentFlags())
	_ = command.cmd.MarkPersistentFlagRequired("config")
	_ = command.cmd.MarkPersistentFlagRequired("idea-project-root")

	return command
}

func (command *GenerateCommand) Register() *cobra.Command {
	return command.cmd
}

func (command *GenerateCommand) Execute(_ *cobra.Command, _ []string) (err error) {
	fmt.Println(command.GlobalFlags.Config)
	fmt.Println(command.GlobalFlags.IdeaProjectRoot)
	fmt.Println(command.GlobalFlags.VscSourcesRoot)

	return nil
}
