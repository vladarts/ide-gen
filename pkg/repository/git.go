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

type GitRepository struct {
	Url         string             `json:"url"`
	FastForward *bool              `json:"fastForward"`
	Remotes     *map[string]string `json:"remotes"`
}

func (r *GitRepository) Vcs() *string {
	value := "git"
	return &value
}

func (r *GitRepository) Clone(root string) error {
	d, err := r.Directory()
	if err != nil {
		return err
	}
	repoPath := path.Join(root, d)

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		err := os.MkdirAll(repoPath, os.ModePerm)
		if err != nil {
			return err
		}

		_, err = ExecCmd("git", []string{
			"clone", r.Url, repoPath,
		}, nil)
		if err != nil {
			return err
		}
	} else {
		curOrigin, err := ExecCmd("git", []string{
			"remote", "get-url", "origin",
		}, &repoPath)
		if err != nil {
			return err
		}

		if curOrigin != r.Url {
			return fmt.Errorf(
				"current origin %s does not match with config value %s",
				curOrigin, r.Url)
		}
	}

	if r.FastForward != nil && *r.FastForward {
		_, err := ExecCmd("git", []string{
			"config", "pull.ff", "only",
		}, &repoPath)
		if err != nil {
			return err
		}
	}

	//: TODO: add multiple upstreams

	return nil
}

func (r *GitRepository) Name() (string, error) {
	parsed, err := giturls.Parse(r.Url)
	if err != nil {
		return "", err
	}

	escapedPath := parsed.EscapedPath()
	if strings.HasSuffix(escapedPath, ".git") {
		escapedPath = escapedPath[:len(escapedPath)-4]
	}

	elements := strings.Split(escapedPath, "/")

	return strings.Join(elements, "."), nil
}

func (r *GitRepository) Directory() (string, error) {
	parsed, err := giturls.Parse(r.Url)
	if err != nil {
		return "", err
	}

	escapedPath := parsed.EscapedPath()
	if strings.HasSuffix(escapedPath, ".git") {
		escapedPath = escapedPath[:len(escapedPath)-4]
	}

	elements := []string{
		parsed.Hostname(),
	}
	elements = append(elements, strings.Split(escapedPath, "/")...)

	return path.Join(elements...), nil
}
