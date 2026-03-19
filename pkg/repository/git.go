package repository

import (
	"fmt"
	giturls "github.com/whilp/git-urls"
	"os"
	"os/exec"
	"path"
	"strings"
)

func ExecCmd(name string, args []string, dir *string) (string, error) {
	cmd := exec.Command(name, args...)
	if dir != nil {
		cmd.Dir = *dir
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

type GitSourcesRootCommander struct {
	Config GitSourcesRootConfig
}

func (r *GitSourcesRootCommander) Clone(targetPath string) error {
	_, err := ExecCmd("git", []string{
		"clone", r.Config.Url, targetPath,
	}, nil)
	if err != nil {
		// Remove a partially-created directory so the next run retries cleanly.
		_ = os.Remove(targetPath)
		return err
	}

	if r.Config.FastForward != nil && *r.Config.FastForward {
		_, err := ExecCmd("git", []string{
			"config", "pull.ff", "only",
		}, &targetPath)
		if err != nil {
			return err
		}
	}

	//: TODO: add multiple remotes

	return nil
}

func (r *GitSourcesRootCommander) Exists(p string) (bool, error) {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false, nil
	}

	//: Check if path is a git repository
	if _, err := os.Stat(path.Join(p, ".git")); os.IsNotExist(err) {
		return false, nil
	}

	curOrigin, err := ExecCmd("git", []string{
		"remote", "get-url", "origin",
	}, &p)
	if err != nil {
		return false, fmt.Errorf("failed to get origin for %s: %w", p, err)
	}

	if curOrigin != r.Config.Url {
		return false, fmt.Errorf(
			"current origin %s does not match with config value %s",
			curOrigin, r.Config.Url)
	}

	return true, nil
}

type GitSourcesRootConfig struct {
	Url         string             `json:"url"`
	FastForward *bool              `json:"fastForward,omitempty"`
	Remotes     *map[string]string `json:"remotes,omitempty"`
}

func (c *GitSourcesRootConfig) Directory(root string) (string, error) {
	parsed, err := giturls.Parse(c.Url)
	if err != nil {
		return "", err
	}

	escapedPath := parsed.EscapedPath()
	if strings.HasSuffix(escapedPath, ".git") {
		escapedPath = escapedPath[:len(escapedPath)-4]
	}

	elements := []string{
		root,
		parsed.Hostname(),
	}
	elements = append(elements, strings.Split(escapedPath, "/")...)

	return path.Join(elements...), nil
}

func (c *GitSourcesRootConfig) Name() (string, error) {
	parsed, err := giturls.Parse(c.Url)
	if err != nil {
		return "", err
	}

	escapedPath := parsed.EscapedPath()
	if strings.HasSuffix(escapedPath, ".git") {
		escapedPath = escapedPath[:len(escapedPath)-4]
	}

	elements := strings.FieldsFunc(escapedPath, func(r rune) bool { return r == '/' })

	return strings.Join(elements, "."), nil
}

func (c *GitSourcesRootConfig) Commander() SourcesRootCommander {
	return &GitSourcesRootCommander{Config: *c}
}

func (c *GitSourcesRootConfig) VcsType() *string {
	value := vcsTypeGit
	return &value
}
