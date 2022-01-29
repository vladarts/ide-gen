package config

import (
	"fmt"
	"github.com/xxxbobrxxx/idea-project-manager/pkg/repository"
)

type Config struct {
	Repositories []RepositoryItem `json:"repositories"`
}

type RepositoryItem struct {
	Name      *string                         `json:"name"`
	Git       *repository.GitRepository       `json:"git"`
	Directory *repository.DirectoryRepository `json:"directory"`
}

func (item RepositoryItem) Repository() (repository.Repository, error) {
	if item.Git != nil {
		return item.Git, nil
	} else if item.Directory != nil {
		return item.Directory, nil
	}

	return nil, fmt.Errorf("can not determine repository type")
}
