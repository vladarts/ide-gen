# CLAUDE.md — ide-gen

## Project Overview

`ide-gen` is a Go CLI tool for managing local development environments. It:
- Clones Git repositories into a structured directory layout (`$VCS_ROOT/hostname/group/repo`)
- Autodiscovers repositories via GitLab API with include/exclude filtering
- Generates IntelliJ IDEA project files (modules, VCS mappings) from a YAML config

**Module:** `github.com/xxxbobrxxx/ide-gen`
**Go version:** 1.24.0

---

## Repository Structure

```
cmd/                        # CLI layer (Cobra commands)
  main.go                   # Entry point, version command
  command-generate.go       # `generate` command
  command-git-clone.go      # `git-clone` command
  command-gopath-link.go    # `gopath-link` command
  command-jsonschema.go     # `json-schema` hidden command

pkg/config/                 # YAML config loading and resolution
  config.go                 # Config struct + ProcessProjectEntry
  flags.go                  # GlobalFlags (config file path)

pkg/repository/             # Repository abstraction layer
  base.go                   # Interfaces: SourcesRootCommander, SourcesRootConfig, ProjectEntry
  git.go                    # Git implementation
  directory.go              # Local directory implementation
  flags.go                  # SourcesRootFlags (vcs-sources-root)

pkg/idea/                   # IntelliJ IDEA project generation
  project.go                # Project struct, file write logic
  xml.go                    # XML templates (GenIml, GenModules, GenVcs)

pkg/gitlab/                 # GitLab API integration
  discover.go               # Autodiscovery with filtering
```

---

## CLI Commands

### `generate` (alias: `gen`)
Clone repositories and generate an IntelliJ IDEA project.

```
ide-gen generate \
  -c ./config.yaml \
  -s ~/dev \
  -i ~/idea-project \
  [--parallel] [--parallel-concurrency N]
```

Flags:
- `-c, --config` (required): YAML config file path
- `-s, --vcs-sources-root`: Root for cloned repos (default: `$HOME/dev`)
- `-i, --idea-project-root`: IDEA project output directory
- `-p, --parallel`: Enable parallel cloning (default: false)
- `--parallel-concurrency`: Worker count (default: `runtime.NumCPU() * 2`)

### `git-clone`
Clone a single Git repository following path conventions.

```
ide-gen git-clone git@github.com:group/repo.git -s ~/dev
```

### `gopath-link`
Create a `$GOPATH/src` symlink for a Go module directory.

```
ide-gen gopath-link ./path/to/module
```

### `json-schema`
Print JSON schema for IDE config file validation/autocompletion (hidden command).

---

## Config File Format (YAML)

```yaml
git:
  - url: git@github.com:group/repo.git
    fast_forward: true
    remotes:
      upstream: git@github.com:upstream/repo.git

raw:
  - path: /path/to/local/dir

gitlab:
  - url: https://gitlab.example.com
    token: "my-token"           # or use token_env_var
    token_env_var: GITLAB_TOKEN
    token_type: private         # private | job | oauth
    include_archived: false
    https_url: false
    fast_forward: true
    include:
      - "^group/.*"
    exclude:
      - "^group/archived-.*"
```

---

## Key Patterns

**Directory layout:** `$VCS_ROOT/{hostname}/{group...}/{repo}`
Example: `~/dev/github.com/myorg/myrepo`

**Module naming (for IDEA):** Path segments joined with dots, `.git` suffix stripped
Example: `github.com:myorg/sub/repo.git` → `myorg.sub.repo`

**Idempotency:** Repos that already exist are skipped; origin URL is validated against config.

**Parallel cloning:** Channel-based semaphore + WaitGroup pattern.

**Authentication:** GitLab token via config field or environment variable; supports private/job/OAuth2 token types.

---

## Key Interfaces

```go
// pkg/repository/base.go
type SourcesRootConfig interface {
    Name() (string, error)
    Directory(sourcesRoot string) (string, error)
    Commander() SourcesRootCommander
    VcsType() *string
}

type SourcesRootCommander interface {
    Clone(directory string) error
    Exists(directory string) (bool, error)
}
```

---

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `sigs.k8s.io/yaml` | YAML config parsing |
| `github.com/whilp/git-urls` | SSH/HTTPS Git URL parsing |
| `gitlab.com/gitlab-org/api/client-go` | GitLab API |
| `golang.org/x/oauth2` | OAuth2 for GitLab auth |
| `golang.org/x/mod` | go.mod parsing |
| `github.com/alecthomas/jsonschema` | JSON schema generation |
| `github.com/sirupsen/logrus` | Structured logging |

---

## Build & Release

```bash
go build ./cmd/...       # local build
go mod tidy              # clean dependencies (run before build in goreleaser)
goreleaser release       # full release (Linux/macOS/Windows, amd64/arm64)
```

- Binary name: `ide-gen`, main package: `./cmd/`
- Version injected at build: `-X main.version={{.Version}}`
- CGO disabled (`CGO_ENABLED=0`)
- Distributed via GitHub Releases and Homebrew (`xxxbobrxxx/homebrew-xxxbobrxxx`)

---

## Error Handling Conventions

- Config read/parse errors → `panic` (fail-fast on startup)
- Template generation errors → `panic` (non-recoverable)
- Clone/exists errors → logged, returned to caller
- CLI errors → `os.Exit(2)` via Cobra
- Parallel workers → log errors but continue other workers

## Logging

`logrus` with text formatter (no timestamps), Info level by default.
