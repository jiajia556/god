package initproject

import (
	"embed"
	"github.com/jiajia556/god/internal/service"
	"github.com/jiajia556/god/internal/template"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func InitProject(name string, tmplFS embed.FS) {
	defer func() {
		service.CmdDir = "./" + name
		service.RunCommand("go", "mod", "tidy")
		_, err := exec.LookPath("goimports")
		if err != nil {
			service.RunCommand("go", "install", "golang.org/x/tools/cmd/goimports@latest")
		}
	}()
	fs.WalkDir(tmplFS, "templates/basic", func(originalPath string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		path := name + strings.TrimPrefix(originalPath, "templates/basic")
		dirPath := filepath.Dir(path)
		os.MkdirAll(dirPath, 0755)
		tmplName := filepath.Base(originalPath)
		if strings.HasSuffix(tmplName, "router.go.tmpl") ||
			strings.HasSuffix(tmplName, "controller.tmpl") ||
			strings.HasSuffix(tmplName, "middleware.tmpl") ||
			strings.HasSuffix(tmplName, "record.go.tmpl") ||
			strings.HasSuffix(tmplName, "list.go.tmpl") {
			return nil
		}

		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}
		targetPath := path[:len(path)-5]

		contentByte, err := fs.ReadFile(tmplFS, originalPath)
		if err != nil {
			service.OutputFatal(err)
		}
		content := string(contentByte)

		fileName := filepath.Base(targetPath)

		onlyProjectNameTmpls := []string{"mysql.go", "go.mod", "main.go", "gopackage.json"}
		if service.InArray(onlyProjectNameTmpls, fileName) {
			data := template.OnlyProjectNameData{ProjectName: name}
			err = template.CreateFile(content, data, targetPath)
			if err != nil {
				service.OutputFatal(err)
			}
		} else {
			f, _ := os.Create(targetPath)
			_, err = f.WriteString(content)
			if err != nil {
				service.OutputFatal(err)
			}
			err = f.Close()
			if err != nil {
				service.OutputFatal(err)
			}
		}
		return nil
	})
}
