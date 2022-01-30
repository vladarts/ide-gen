# Development Workspace Manager

`ide-gen` is a tool for development workspace prepare automation by automatic
VCS repositories discovery and clone and project generation for supported IDEs.

## Installation

For macOS users:

```
$ brew tap xxxbobrxxx/xxxbobrxxx
$ brew install ide-gen
```

You can also download the binary from the
[Releases](https://github.com/xxxbobrxxx/ide-gen/releases) page.

### Json Schema

Use [json schema](resources/config.schema.json) for autocompletion in your
favorite IDE.

## Quickstart

Define a sample project configuration in the, for example,
`${HOME}/main.yaml` file:

```yaml
#: Local non-VCS directories
directory:
  - path: /your/local/dir

#: Git repositories
git:
  -
    #: Git repository URL, required
    url: git@github.com:xxxbobrxxx/ide-gen.git

    #: Apply `git config pull.ff only` during clone.
    #: Optional, default: false
    #:
    # fastForward: true

    #: Additional origins to add to repo to.
    #: Optional, default: {}
    #:
    # remotes: {}
    #   another-origin: git@github.com:xxxbobrxxx/ide-gen-fork.git

#: Gitlab Autodiscovery configs
gitlab:
  -
    #: Gitlab API URL
    #: Optional, default: https://gitlab.com/api/v4
    #:
    # url:

    #: token or tokenEnvVar must be defined to access Gitlab API
    #:
    #: Value of Gitlab API token
    #:
    # token: XXXXXXXXXXXXX

    #: Name of environment variable containing Gitlab API token
    #:
    tokenEnvVar: GITLAB_TOKEN

    #: Use HTTPS Url instead of SSH
    #: Optional, default - false
    #:
    # httpsUrl: false

    #: Apply `git config pull.ff only` during clone.
    #: Optional, default: false
    #:
    # fastForward: true

    #: include/exclude regexp patterns for repositories paths
    #:
    #: Example: for repository `https://gitlab.com/group1/group2/repo`
    #: path is `group1/group2/repo`
    #:
    #: Optional, default - []
    #:
    include:
      - ^group-to-include/.*
    exclude:
      - ^group-to-include/subgroup-to-exclude/.*
```

Execute a command to discover and clone repositories:

```bash
$ ide-gen gen -c ${HOME}/main.yaml
```

It will:

- Keep `/your/local/dir` directory untouched
- Autodiscover repositories from a Gitlab server and add to `git` repo list.
- Clone `git@github.com:xxxbobrxxx/ide-gen.git` repository to the
  `${HOME}/dev/github.com/xxxbobrxxx/ide-gen`. If already cloned - skip.

To enable IntelliJ IDEA project generation run the command with the
`-i/--idea-project-root` param:

```bash
$ ide-gen gen -c ${HOME}/main.yaml -i ${HOME}/dev/idea_projects/main
```

It will:

- Create a `${HOME}/dev/idea_projects/main` directory if it does not exist.
  It will be a root of the IntelliJ IDEA project.
- Inside the IDEA project dir:
  - Just create `.idea/iml/<project_name>.iml` for each repository/directory.
    Overwrite is forbidden to avoid manual settings loss.
  - Create/Overwrite `.idea/modules.xml` containing all repositories
  - Create/Overwrite `.idea/vcs.xml` containing proper VCS mappings

### Repository clone path

Directories paths to clone VCS repositories to are generated automatically
following the rule

```
$VCS_ROOT/vcs-hostname/group-1[/group-2][/group-...]/repo-name
```

By default `$VCS_ROOT` is `${HOME}/dev` and it can be overridden by the
`-s/--vcs-sources-root` param.

### IntelliJ IDEA modules structure

Intellij IDEA supports adding of several modules to the one project.
The idea of the tools is to keep modules structure as they are present in
the VCS. In comparison to Github, Gitlab supports deeper subgroup level so
that modules structure may be useful.

To achieve it the tool automatically parses VCS URLs to define modules names.
Example:

- `git@github.com:xxxbobrxxx/ide-gen.git` -> `xxxbobrxxx.ide-gen`
- `https://github.com/xxxbobrxxx/ide-gen.git` -> `xxxbobrxxx.ide-gen`

For local directories it uses basename:

- `/your/local/dir` -> `dir`

Example: for the following repos list:

- `/your/local/dir`
- `git@github.com:group1/repo1.git`
- `git@github.com:group1/repo2.git`
- `git@github.com:group2/repo1.git`
- `git@github.com:group2/repo2.git`

The following modules' tree will be generated:

```
├── group1
│   ├── repo1
│   └── repo2
├── group2
│   ├── repo1
│   └── repo2
└── dir
```
