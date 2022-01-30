# Development Workspace Manager

`ide-gen` is a tool to simplify local workspace management by defining a
strict directories structure for VCS repositories required to be cloned and
automatic workspaces generation for supported IDE projects.

## Quickstart

Define a sample project configuration in the, for example,
`${HOME}/main.yaml` file:

```yaml
#: Local non-VCS directories
directory:
  - path: /your/local/dir

#: Git repository
git:
  -
    #: Git repository URL, required
    url: git@github.com:xxxbobrxxx/ide-gen.git

    #: Apply `git config pull.ff only` during clone.
    #: Optional, default: false
    # fastForward: true

    #: Additional origins to add to repo to.
    #: Optional, default: {}
    # remotes: {}
    #   another-origin: git@github.com:xxxbobrxxx/ide-gen-fork.git

#: Gitlab Autodiscovery configs
gitlab:
  -
    #: Gitlab API URL
    #: Optional, default: https://gitlab.com/api/v4
    # url:

    #: token or tokenEnvVar must be defined to access Gitlab API
    # token: XXXXXXXXXXXXX
    tokenEnvVar: GITLAB_TOKEN

    #: Use HTTPS Url instead of SSH
    #: Optional, default - false
    # httpsUrl: false

    #: Apply `git config pull.ff only` during clone.
    #: Optional, default: false
    # fastForward: true

    #: include/exclude patterns for repositories paths
    #:
    #: Example: for repository `https://gitlab.com/group1/group2/repo`
    #: path is `group1/group2/repo`
    #:
    #: Optional, default - []
    include:
      - ^group-to-include/.*
    exclude:
      - ^group-to-include/subgroup-to-exclude/.*
```

Execute a command to clone repositories only:

```bash
$ ide-gen gen -c ${HOME}/main.yaml
```

It will:

- Keep `/your/local/dir` directory untouched
- Clone `git@github.com:xxxbobrxxx/ide-gen.git` repository to the
  `${HOME}/dev/github.com/xxxbobrxxx/ide-gen`. If already cloned - skip.

To enable IntelliJ IDEA project generation command with params:

```bash
$ ide-gen gen -c ${HOME}/main.yaml -i ${HOME}/dev/idea_projects/main
```

It will:

- Create a `${HOME}/dev/idea_projects/main` directory if it does not exist
- Inside the IDEA project dir:
  - Create **only** `.idea/iml/<project_name>.iml` for each repository.
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

Program automatically parses VCS URLs to define modules names.

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

The following modules tree will be generated:

```
├── group1
│   ├── repo1
│   └── repo2
├── group2
│   ├── repo1
│   └── repo2
└── dir
```
