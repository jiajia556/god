package template

import (
	"os"
	"text/template"
)

type RouterTmplData struct {
	MiddlewareImportPath  string
	ControllersImportPath string
	ApiRootDirName        string
	HTTPMethodTags        string
	MiddlewareTags        string
	RegisterControllers   string
}

type OnlyProjectNameData struct {
	ProjectName string
}

type ControllerStructNameData struct {
	ControllerStructName string
}

type MiddlewareNameData struct {
	MiddlewareName string
}

type ModelData struct {
	ModelPkg        string
	ProjectName     string
	ModelStruct     string
	ModelStructName string
}

func CreateFile(routerTmpl string, data any, path string) error {
	tmpl := template.Must(template.New("").Parse(routerTmpl))
	f, _ := os.Create(path)
	defer f.Close()

	err := tmpl.Execute(f, data)
	if err != nil {
		return err
	}
	return nil
}
