package repository

const (
	vcsTypeGit = "Git"
)

type SourcesRootCommander interface {
	Clone(string) error
}

type SourcesRootConfig interface {
	Name() (string, error)
	Directory(string) (string, error)
	Commander() SourcesRootCommander
	VcsType() *string
}

type ProjectEntry struct {
	Name      string
	Directory string
	VcsType   *string

	Commander SourcesRootCommander
}
