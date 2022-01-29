package config

import (
	"fmt"
	"github.com/xxxbobrxxx/idea-project-manager/pkg/repository"
)

type Config struct {
	RepositoryConfigs []RepositoryConfig `json:"repositories"`
}

type RepositoryConfig struct {
	Name      *string                         `json:"name"`
	Git       *repository.GitRepository       `json:"git"`
	Directory *repository.DirectoryRepository `json:"directory"`
}

func (c *RepositoryConfig) NewFromConfig() (repository.Repository, error) {
	if c.Git != nil {
		return c.Git, nil
	} else if c.Directory != nil {
		return c.Directory, nil
	}

	return nil, fmt.Errorf("can not recognize repository type")
}
