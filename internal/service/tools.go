package service

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetFileByRoute(route string) (filePath, fileName string, err error) {
	// Validate route format constraints
	if strings.HasPrefix(route, "/") || strings.HasSuffix(route, "/") {
		err = fmt.Errorf("route must not start or end with '/'")
		return "", "", err
	}

	// Locate the last directory separator
	lastSlashPos := strings.LastIndex(route, "/")
	if lastSlashPos == -1 {
		// Handle simple case with no subdirectories
		filePath = fmt.Sprintf("controller/%s.go", route)
		return filePath, route, nil
	}

	// Split route into directory and component name
	directory := route[:lastSlashPos]
	component := route[lastSlashPos+1:]

	// Construct controller file path
	filePath = fmt.Sprintf("%s/controller/%s.go", directory, component)
	return filePath, component, nil
}

func ValidateControllerName(s string) error {
	if strings.Contains(s, " ") {
		return errors.New("controller name can not contain spaces")
	}
	if strings.Contains(s, "_") {
		return errors.New("controller name can not contain _")
	}
	if strings.Contains(s, "-") {
		return errors.New("controller name can not contain -")
	}
	return nil
}

func FileExists(filename string) bool {
	// 获取文件信息
	// Get file information
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s // Return directly if string is empty
	}

	// Capitalize first letter, keep the rest unchanged
	return strings.ToUpper(string(s[0])) + s[1:]
}

func InArray[T comparable](array []T, value T) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

var GoEnv = []string{}
var CmdDir = ""

func RunCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	if CmdDir != "" {
		cmd.Dir = CmdDir
	}
	cmd.Stdout = nil
	cmd.Stderr = nil
	if len(GoEnv) > 0 {
		cmd.Env = append(os.Environ(), GoEnv...)
	}

	if err := cmd.Run(); err != nil {
		OutputFatal(fmt.Sprintf("Command %s failed: %v", name, err))
	}
}
