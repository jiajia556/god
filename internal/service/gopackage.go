package service

import (
	"encoding/json"
	"os"
)

type GoPackage struct {
	inited         bool
	ProjectName    string `json:"project_name"`
	DefaultAppRoot string `json:"default_app_root"`
	DefaultApiRoot string `json:"default_api_root"`
	DefaultGOOS    string `json:"default_goos"`
	DefaultGOARCH  string `json:"default_goarch"`
}

var goPackage GoPackage

func initGoPackage() error {
	dat, err := os.ReadFile("./gopackage.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(dat, &goPackage)
	if err != nil {
		return err
	}
	goPackage.inited = true
	return nil
}

func GetDefaultAppRoot() (string, error) {
	if !goPackage.inited {
		err := initGoPackage()
		if err != nil {
			return "", err
		}
	}
	return goPackage.DefaultAppRoot, nil
}

func GetDefaultApiRoot() (string, error) {
	if !goPackage.inited {
		err := initGoPackage()
		if err != nil {
			return "", err
		}
	}
	return goPackage.DefaultApiRoot, nil
}

func GetDefaultGOOS() (string, error) {
	if !goPackage.inited {
		err := initGoPackage()
		if err != nil {
			return "", err
		}
	}
	return goPackage.DefaultGOOS, nil
}

func GetDefaultGOARCH() (string, error) {
	if !goPackage.inited {
		err := initGoPackage()
		if err != nil {
			return "", err
		}
	}
	return goPackage.DefaultGOARCH, nil
}

func GetProjectName() (string, error) {
	if !goPackage.inited {
		err := initGoPackage()
		if err != nil {
			return "", err
		}
	}
	return goPackage.ProjectName, nil
}
