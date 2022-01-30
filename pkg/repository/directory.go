package repository

import (
	"path"
)

type RawSourcesRootCommander struct{}

func (r *RawSourcesRootCommander) Clone(_ string) error {
	return nil
}

type RawSourcesRootConfig struct {
	Path string `json:"path"`
}

func (c *RawSourcesRootConfig) Directory(_ string) (string, error) {
	return c.Path, nil
}

func (c *RawSourcesRootConfig) Name() (string, error) {
	return path.Base(c.Path), nil
}

func (c *RawSourcesRootConfig) Commander() SourcesRootCommander {
	return &RawSourcesRootCommander{}
}

func (c *RawSourcesRootConfig) VcsType() *string {
	return nil
}
