package main

import (
	"github.com/spf13/cobra"
	"github.com/xxxbobrxxx/ide-gen/pkg/config"
	"github.com/xxxbobrxxx/ide-gen/pkg/idea"
	"github.com/xxxbobrxxx/ide-gen/pkg/repository"
	"runtime"
	"sync"
)

type GenerateCommand struct {
	config.GlobalFlags
	repository.SourcesRootFlags
	idea.Project

	cmd *cobra.Command
}

func NewGenerateCommand() *GenerateCommand {
	command := &GenerateCommand{}

	cmd := &cobra.Command{
		Use:          "generate",
		Aliases:      []string{"gen"},
		Short:        "Clone repositories and generate IDE project",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		RunE:         command.Execute,
	}
	command.cmd = cmd

	command.Project.AddFlags(cmd.PersistentFlags())
	command.GlobalFlags.AddFlags(cmd.PersistentFlags())
	command.SourcesRootFlags.AddFlags(cmd.PersistentFlags())

	_ = command.cmd.MarkPersistentFlagRequired("config")

	return command
}

func (command *GenerateCommand) Register() *cobra.Command {
	return command.cmd
}

func (command *GenerateCommand) Execute(_ *cobra.Command, _ []string) (err error) {
	c, err := command.ReadConfig()
	if err != nil {
		return err
	}

	projectEntries, err := c.GetProjectEntries(command.SourcesRootFlags)
	if err != nil {
		return err
	}

	// Determine the number of workers based on CPU count
	numWorkers := runtime.NumCPU() * 2
	// Create a buffered channel to limit concurrent goroutines
	sem := make(chan struct{}, numWorkers)

	var wg sync.WaitGroup
	for _, projectEntry := range projectEntries {
		wg.Add(1)
		sem <- struct{}{} // Acquire a slot
		go func(entry repository.ProjectEntry) {
			defer wg.Done()
			defer func() { <-sem }() // Release the slot

			exists, err := entry.Commander.Exists(entry.Directory)
			if err != nil {
				logger.Errorf("Error checking existence: %v", err)
				return
			}
			if exists {
				logger.Infof("Skip clone project '%s' to '%s'", entry.Name, entry.Directory)
			} else {
				logger.Infof("Clone project '%s' to '%s'", entry.Name, entry.Directory)
				if err := entry.Commander.Clone(entry.Directory); err != nil {
					logger.Errorf("Error cloning project: %v", err)
				}
			}
		}(projectEntry)
	}
	wg.Wait() // Wait for all goroutines to finish

	// Idea project
	if command.Project.Root != "" {
		project := command.Project
		for _, projectEntry := range projectEntries {
			project.AddEntry(projectEntry)
		}

		logger.Infof("Writing idea project %s", project.Root)
		if err := project.Write(); err != nil {
			return err
		}
	}

	return nil
}
