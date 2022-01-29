package repository

type Repository interface {
	Init(RepositoryFlags) error
	Clone() (string, error)
	Name() string
	Directory() string
	Vcs() *string
}
