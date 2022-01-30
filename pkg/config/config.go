package config

import (
	"github.com/xxxbobrxxx/ide-gen/pkg/gitlab"
	"github.com/xxxbobrxxx/ide-gen/pkg/repository"
)

type Config struct {
	GitlabEntries []*gitlab.DiscoveryConfig          `json:"gitlab"`
	GitEntries    []*repository.GitSourcesRootConfig `json:"git"`
	RawEntries    []*repository.RawSourcesRootConfig `json:"directory"`
}

func (c *Config) GetProjectEntries(flags repository.SourcesRootFlags) ([]repository.ProjectEntry, error) {
	var projectEntries []repository.ProjectEntry

	var sourceRootConfigs []repository.SourcesRootConfig

	for _, entryConfig := range c.GitEntries {
		sourceRootConfigs = append(sourceRootConfigs, entryConfig)
	}
	for _, entryConfig := range c.RawEntries {
		sourceRootConfigs = append(sourceRootConfigs, entryConfig)
	}

	for _, gl := range c.GitlabEntries {
		err := gl.Init()
		if err != nil {
			return nil, err
		}

		projects, err := gl.Discover()
		if err != nil {
			return nil, err
		}
		for _, entryConfig := range projects {
			sourceRootConfigs = append(sourceRootConfigs, entryConfig)
		}
	}

	//: Read projects from config
	for _, entryConfig := range sourceRootConfigs {
		var (
			name, directory string
			err             error
		)

		name, err = entryConfig.Name()
		if err != nil {
			return nil, err
		}

		directory, err = entryConfig.Directory(flags.VscSourcesRoot)
		if err != nil {
			return nil, err
		}

		projectEntries = append(projectEntries, repository.ProjectEntry{
			Name:      name,
			Directory: directory,
			VcsType:   entryConfig.VcsType(),
			Commander: entryConfig.Commander(),
		})
	}

	return projectEntries, nil
}
