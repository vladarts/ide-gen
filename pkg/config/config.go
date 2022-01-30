package config

import (
	"fmt"
	"github.com/xxxbobrxxx/idea-project-manager/pkg/repository"
)

type Config struct {
	ProjectEntries []RepositoryConfig `json:"repositories"`
}

type RepositoryConfig struct {
	Git       *repository.GitSourcesRootConfig `json:"git"`
	Directory *repository.RawSourcesRootConfig `json:"directory"`
}

func (c *Config) GetProjectEntries(flags repository.SourcesRootFlags) ([]repository.ProjectEntry, error) {
	var projectEntries []repository.ProjectEntry

	//: Read projects from config
	for _, projectEntryConfig := range c.ProjectEntries {
		var (
			sourcesRootConfig repository.SourcesRootConfig
			name, directory   string
			err               error
		)

		if projectEntryConfig.Git != nil {
			sourcesRootConfig = projectEntryConfig.Git
		} else if projectEntryConfig.Directory != nil {
			sourcesRootConfig = projectEntryConfig.Directory
		} else {
			return nil, fmt.Errorf("can not determine repository type")
		}

		name, err = sourcesRootConfig.Name()
		if err != nil {
			return nil, err
		}

		directory, err = sourcesRootConfig.Directory(flags.VscSourcesRoot)
		if err != nil {
			return nil, err
		}

		projectEntries = append(projectEntries, repository.ProjectEntry{
			Name:      name,
			Directory: directory,
			VcsType:   sourcesRootConfig.VcsType(),
			Commander: sourcesRootConfig.Commander(),
		})
	}

	return projectEntries, nil
}
