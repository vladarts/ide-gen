package main

import (
	"fmt"
	"github.com/xxxbobrxxx/idea-project-manager/pkg/config"
	"github.com/xxxbobrxxx/idea-project-manager/pkg/idea"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sigs.k8s.io/yaml"
)

func main() {
	file := "/Users/vpiskunov/xxxbobrxxx@gmail.com/dev/bashrc/idea_projects/vp-v2.yaml"
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var conf config.Config

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	var modules []idea.Module

	var projectImlSubdir = "iml"

	for _, item := range conf.Repositories {
		r, err := item.Repository()
		if err != nil {
			panic(err)
		}

		n, err := r.Name()
		if err != nil {
			panic(err)
		}

		d, err := r.Directory()
		if err != nil {
			panic(err)
		}

		dir, _ := filepath.Rel(home, path.Join(home, "dev", d))

		projectImlPath := path.Join(
			projectImlSubdir,
			fmt.Sprintf("%s.iml", n))
		module := idea.Module{
			Directory: dir,
			Vcs:       r.Vcs(),
			ImlPath:   projectImlPath,
		}
		modules = append(modules, module)

		fmt.Println(idea.GenIml(module))
	}

	fmt.Println(idea.GenModules(modules))
	fmt.Println(idea.GenVcs(modules))
}
