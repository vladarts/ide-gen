package repository

type Repository interface {
	Clone(root string) error
	Name() (string, error)
	Directory() (string, error)
	Vcs() *string
}
