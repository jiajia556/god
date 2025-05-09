package addaction

import (
	"fmt"
	"github.com/jiajia556/god/internal/service"
	"os"
	"path/filepath"
	"strings"
)

func AddAction(root, controllerRoute string, actions []string) {
	var err error
	if root == "" {
		root, err = service.GetDefaultApiRoot()
		if err != nil {
			service.OutputFatal(err)
		}
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
	if !service.FileExists(controllerFilePath) {
		service.OutputFatal("Error: controller is not exists")
	}
	controllerStructName := service.CapitalizeFirstLetter(name) + "Controller"
	WriteActions(controllerFilePath, controllerStructName, actions)
}

func WriteActions(controllerFilePath, controllerStructName string, actions []string) {
	actionList, err := makeActions(actions)
	if err != nil {
		service.OutputFatal(err)
	}
	file, err := os.OpenFile(controllerFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		service.OutputFatal(err)
	}
	for _, v := range actionList {
		methodStr := fmt.Sprintf(service.CONTROLLER_ACTION_TMPL,
			v.HTTPMethod,
			controllerStructName,
			v.Name,
		)
		_, err = file.WriteString(methodStr)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
}

type method struct {
	Name       string
	HTTPMethod string
}

func makeActions(actions []string) (res []method, err error) {
	length := len(actions)
	if length == 0 {
		return
	}

	res = make([]method, length)
	for k, mtd := range actions {
		mtdDetail := strings.Split(mtd, ":")
		for i, v := range mtdDetail {
			if i == 0 {
				res[k].Name = service.CapitalizeFirstLetter(v)
			} else {
				switch strings.ToLower(v) {
				case "post":
					res[k].HTTPMethod = "POST"
				case "get":
					res[k].HTTPMethod = "GET"
				default:
					err = fmt.Errorf("invalid method: %s", v)
					return
				}
			}
			if res[k].HTTPMethod == "" {
				res[k].HTTPMethod = "POST"
			}
		}
	}
	return
}
