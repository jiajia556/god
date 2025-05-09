package addcontroller

import (
	"fmt"
	"github.com/jiajia556/god/internal/cmd/addaction"
	"github.com/jiajia556/god/internal/service"
	"github.com/jiajia556/god/internal/template"
	"os"
	"path/filepath"
)

func AddController(controllerTmpl, root, controllerRoute string, actions []string) {
	var err error

	if root == "" {
		root, err = service.GetDefaultApiRoot()
		if err != nil {
			service.OutputFatal(err)
		}
	}

	if controllerRoute == "" {
		controllerRoute, err = service.InputStr("please enter controller route:")
		if err != nil {
			service.OutputFatal(err)
		}
	}
	if controllerRoute == "" {
		service.OutputFatal("Error: controller route is empty")
	}

	path, name, err := service.GetFileByRoute(controllerRoute)
	if err != nil {
		service.OutputFatal(err)
	}
	err = service.ValidateControllerName(name)
	if err != nil {
		service.OutputFatal(err)
	}
	controllerFilePath := filepath.Join(root, path)

	if service.FileExists(controllerFilePath) {
		service.OutputFatal("Error: controller already exists")
	}

	err = os.MkdirAll(filepath.Dir(controllerFilePath), os.ModePerm)
	if err != nil {
		fmt.Println("Error creating dir:", err)
		return
	}

	controllerStructName := service.CapitalizeFirstLetter(name) + "Controller"

	err = template.CreateFile(controllerTmpl,
		template.ControllerStructNameData{controllerStructName},
		controllerFilePath,
	)
	if err != nil {
		service.OutputFatal(err)
	}

	if len(actions) > 0 {
		addaction.WriteActions(controllerFilePath, controllerStructName, actions)
	}
}
