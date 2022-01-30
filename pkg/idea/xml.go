package idea

import (
	"bytes"
	"strings"
	"text/template"
)

const (
	tplImlPython = `
<?xml version="1.0" encoding="UTF-8"?>
<module type="WEB_MODULE" version="4">
  <component name="FacetManager">
    <facet type="Python" name="Python">
    </facet>
  </component>
  <component name="NewModuleRootManager" inherit-compiler-output="true">
    <exclude-output />
    <content url="file://{{ .Module.Directory }}" />
    <orderEntry type="inheritedJdk" />
    <orderEntry type="sourceFolder" forTests="false" />
  </component>
</module>`

	tplModules = `
<?xml version="1.0" encoding="UTF-8"?>
<project version="4">
  <component name="ProjectModuleManager">
    <modules>
    {{- range $module := .Modules }}
      <module fileurl="file://{{ $module.ImlPath }}" filepath="{{ $module.ImlPath }}" />
    {{- end }}
    </modules>
  </component>
</project>
`

	tplVcs = `
<?xml version="1.0" encoding="UTF-8"?>
<project version="4">
  <component name="VcsDirectoryMappings">
  {{- range $module := .Modules }}
  {{- if $module.Vcs }}
    <mapping directory="{{ $module.Directory }}" vcs="{{ $module.Vcs }}" />
  {{- end }}
  {{- end }}
  </component>
</project>
`
)

func genTemplate(tpl string, ctx interface{}) string {
	tmpl, err := template.New("test").Parse(tpl)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, ctx)
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(buf.String())
}

type GenModuleContext struct {
	Module Module
}

func GenIml(module Module) string {
	ctx := GenModuleContext{
		Module: module,
	}

	return genTemplate(tplImlPython, ctx)
}

type GenModulesContext struct {
	Modules []Module
}

func GenModules(modules []Module) string {
	ctx := GenModulesContext{
		Modules: modules,
	}

	return genTemplate(tplModules, ctx)
}

func GenVcs(modules []Module) string {
	ctx := GenModulesContext{
		Modules: modules,
	}

	return genTemplate(tplVcs, ctx)
}
