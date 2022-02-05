package main

import (
	"fmt"
	"github.com/alecthomas/jsonschema"
	"github.com/spf13/cobra"
	"github.com/xxxbobrxxx/ide-gen/pkg/config"
)

type JsonSchemaCommand struct {
	cmd *cobra.Command
}

func NewJsonSchemaCommandCommand() *JsonSchemaCommand {
	command := &JsonSchemaCommand{}

	cmd := &cobra.Command{
		Use:          "json-schema",
		Aliases:      []string{"gen"},
		Short:        "Generate and output json schema for the config file",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		RunE:         command.Execute,
		Hidden:       true,
	}
	command.cmd = cmd

	return command
}

func (command *JsonSchemaCommand) Register() *cobra.Command {
	return command.cmd
}

func (command *JsonSchemaCommand) Execute(_ *cobra.Command, _ []string) (err error) {
	reflector := jsonschema.Reflector{
		ExpandedStruct: true,
	}
	schema := reflector.Reflect(&config.Config{})
	b, _ := schema.MarshalJSON()

	fmt.Println(string(b))
	return nil
}
