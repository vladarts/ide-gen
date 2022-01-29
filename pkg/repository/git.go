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

type GitRepositoryConfig struct {
	Url         string             `json:"url"`
	FastForward *bool              `json:"fastForward"`
	Remotes     *map[string]string `json:"remotes"`
}

type GitRepository struct {
	Config GitRepositoryConfig

	name      string
	directory string
}

func (r *GitRepository) setName() error {
	parsed, err := giturls.Parse(r.Config.Url)
	if err != nil {
		return err
	}

	escapedPath := parsed.EscapedPath()
	if strings.HasSuffix(escapedPath, ".git") {
		escapedPath = escapedPath[:len(escapedPath)-4]
	}

	elements := strings.Split(escapedPath, "/")

	r.name = strings.Join(elements, ".")

	return nil
}

func (r *GitRepository) Init(flags RepositoryFlags) error {
	err := r.setName()
	if err != nil {
		return err
	}

	err = r.setDirectory(flags.VscSourcesRoot)
	if err != nil {
		return err
	}

	return nil
}

func (r *GitRepository) setDirectory(vcsRoot string) error {
	parsed, err := giturls.Parse(r.Config.Url)
	if err != nil {
		return err
	}

	escapedPath := parsed.EscapedPath()
	if strings.HasSuffix(escapedPath, ".git") {
		escapedPath = escapedPath[:len(escapedPath)-4]
	}

	elements := []string{
		vcsRoot,
		parsed.Hostname(),
	}
	elements = append(elements, strings.Split(escapedPath, "/")...)

	r.directory = path.Join(elements...)

	return nil
}

func (r *GitRepository) Vcs() *string {
	value := "Git"
	return &value
}

func (r *GitRepository) Clone() (string, error) {
	repoPath := r.Directory()

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		err := os.MkdirAll(repoPath, os.ModePerm)
		if err != nil {
			return "", err
		}

		_, err = ExecCmd("git", []string{
			"clone", r.Config.Url, repoPath,
		}, nil)
		if err != nil {
			return "", err
		}
	} else {
		curOrigin, err := ExecCmd("git", []string{
			"remote", "get-url", "origin",
		}, &repoPath)
		if err != nil {
			return "", err
		}

		if curOrigin != r.Config.Url {
			return "", fmt.Errorf(
				"current origin %s does not match with config value %s",
				curOrigin, r.Config.Url)
		}
	}

	if r.Config.FastForward != nil && *r.Config.FastForward {
		_, err := ExecCmd("git", []string{
			"config", "pull.ff", "only",
		}, &repoPath)
		if err != nil {
			return "", err
		}
	}

	//: TODO: add multiple upstreams

	return repoPath, nil
}

func (r *GitRepository) Name() string {
	return r.name
}

func (r *GitRepository) Directory() string {
	return r.directory
}
