package repository

import "path"

type DirectoryRepository struct {
	Path string `json:"path"`
}

func (r *DirectoryRepository) Vcs() *string {
	return nil
}

func (r *DirectoryRepository) Clone(_ string) error {
	return nil
}

func (r *DirectoryRepository) Name() (string, error) {
	return path.Base(r.Path), nil
}

func (r *DirectoryRepository) Directory() (string, error) {
	return r.Path, nil
}
