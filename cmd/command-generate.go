package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

	Parallel            bool
	ParallelConcurrency int

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
	command.AddFlags(cmd.PersistentFlags())

	_ = command.cmd.MarkPersistentFlagRequired("config")

	return command
}

func (command *GenerateCommand) Register() *cobra.Command {
	return command.cmd
}

func (command *GenerateCommand) AddFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&command.Parallel, "parallel", "p",
		false, "Use parallel download")
	flags.IntVar(&command.ParallelConcurrency, "parallel-concurrency", runtime.NumCPU()*2,
		"Parallel download concurrency, default is runtime.NumCPU() * 2")
}

func (command *GenerateCommand) ProcessProjectEntry(projectEntry repository.ProjectEntry) (err error) {
	exists, err := projectEntry.Commander.Exists(projectEntry.Directory)
	if err != nil {
		return
	}

	if exists {
		logger.Infof(
			"Skip clone project '%s' to '%s'", projectEntry.Name, projectEntry.Directory)
	} else {
		logger.Infof(
			"Clone project '%s' to '%s'", projectEntry.Name, projectEntry.Directory)
		return projectEntry.Commander.Clone(projectEntry.Directory)
	}
	return
}

func (command *GenerateCommand) Execute(_ *cobra.Command, _ []string) (err error) {
	c, err := command.ReadConfig()
	if err != nil {
		return err
	}

	//: Read entries from config and flags
	projectEntries, err := c.GetProjectEntries(command.SourcesRootFlags)
	if err != nil {
		return err
	}

	//: Clone repos
	if command.Parallel {
		sem := make(chan struct{}, command.ParallelConcurrency)
		var wg sync.WaitGroup

		for _, projectEntry := range projectEntries {
			wg.Add(1)
			sem <- struct{}{}
			go func(projectEntry repository.ProjectEntry) {
				defer wg.Done()
				defer func() { <-sem }()

				err := command.ProcessProjectEntry(projectEntry)
				if err != nil {
					logger.WithError(err).Errorf("error cloning the repo: %v", projectEntry.Name)
				}
			}(projectEntry)
		}
		wg.Wait()
	} else {
		for _, projectEntry := range projectEntries {
			if err := command.ProcessProjectEntry(projectEntry); err != nil {
				logger.WithError(err).Errorf("error cloning the repo: %v", projectEntry.Name)
			}
		}
	}

	//: Idea project
	if command.Project.Root != "" {
		project := command.Project
		for _, projectEntry := range projectEntries {
			project.AddEntry(projectEntry)
		}

		logger.Infof("Writing idea project %s", project.Root)
		err = project.Write()
		if err != nil {
			return err
		}
	}

	return nil
}
