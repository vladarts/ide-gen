package config

import (
	"fmt"
	"github.com/xxxbobrxxx/idea-project-manager/pkg/repository"
)

type Config struct {
	Repositories []RepositoryConfig `json:"repositories"`
}

type RepositoryConfig struct {
	Name      *string                               `json:"name"`
	Git       *repository.GitRepositoryConfig       `json:"git"`
	Directory *repository.DirectoryRepositoryConfig `json:"directory"`
}

func (c *RepositoryConfig) NewFromConfig() (repository.Repository, error) {
	if c.Git != nil {
		return &repository.GitRepository{Config: *c.Git}, nil
	} else if c.Directory != nil {
		return &repository.DirectoryRepository{Config: *c.Directory}, nil
	}

	return nil, fmt.Errorf("can not recognize repository type")
}
