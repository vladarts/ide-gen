package repository

import (
	"path"
)

type DirectoryRepositoryConfig struct {
	Path string `json:"path"`
}

type DirectoryRepository struct {
	Config DirectoryRepositoryConfig
}

func (r *DirectoryRepository) Init(_ RepositoryFlags) error {
	return nil
}

func (r *DirectoryRepository) Vcs() *string {
	return nil
}

func (r *DirectoryRepository) Clone() (string, error) {
	return r.Directory(), nil
}

func (r *DirectoryRepository) Name() string {
	return path.Base(r.Config.Path)
}

func (r *DirectoryRepository) Directory() string {
	return r.Config.Path
}
