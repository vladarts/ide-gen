package repository

type Repository interface {
	Init(vcsRoot string) error
	Clone() (string, error)
	Name() string
	Directory() string
	Vcs() *string
}
