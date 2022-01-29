package idea

import (
	"github.com/xxxbobrxxx/idea-project-manager/pkg/repository"
)

type Project struct {
	Root string

	Modules []Module
}

func (p *Project) AddRepository(repository repository.Repository) {
	p.Modules = append(p.Modules, Module{})
}
