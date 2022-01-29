package idea

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/xxxbobrxxx/idea-project-manager/pkg/repository"
	"os"
	"path"
)

const (
	projectImlSubdir = "iml"
	ideaSubdir       = ".idea"
	modulesFileName  = "modules.xml"
	vcsFileName      = "vcs.xml"
)

type Module struct {
	Directory string
	Vcs       *string
	ImlPath   string
}

type Project struct {
	Root string

	Modules []Module
}

func (p *Project) AddFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&p.Root, "idea-project-root", "i",
		"", "")
}

func (p *Project) ImlDir() string {
	return path.Join(p.Root, ideaSubdir, projectImlSubdir)
}

func (p *Project) ModulesPath() string {
	return path.Join(p.Root, ideaSubdir, modulesFileName)
}

func (p *Project) VcsPath() string {
	return path.Join(p.Root, ideaSubdir, vcsFileName)
}

func (p *Project) AddRepository(r repository.Repository) {
	directory := r.Directory()
	name := r.Name()

	projectImlPath := path.Join(
		p.ImlDir(),
		fmt.Sprintf("%s.iml", name))
	module := Module{
		Directory: directory,
		Vcs:       r.Vcs(),
		ImlPath:   projectImlPath,
	}
	p.Modules = append(p.Modules, module)
}

func (p *Project) Write() error {
	if _, err := os.Stat(p.ImlDir()); os.IsNotExist(err) {
		err := os.MkdirAll(p.ImlDir(), os.ModePerm)
		if err != nil {
			return err
		}
	}

	err := os.WriteFile(p.ModulesPath(), []byte(GenModules(p.Modules)), 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile(p.VcsPath(), []byte(GenVcs(p.Modules)), 0644)
	if err != nil {
		return err
	}

	for _, module := range p.Modules {
		err := os.WriteFile(module.ImlPath, []byte(GenIml(module)), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
