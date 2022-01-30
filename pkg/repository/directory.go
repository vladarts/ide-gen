package repository

import (
	"fmt"
	"os"
	"path"
)

type RawSourcesRootCommander struct{}

func (r *RawSourcesRootCommander) Clone(_ string) error {
	return fmt.Errorf("directory resource root can not be cloned")
}

func (r *RawSourcesRootCommander) Exists(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
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
