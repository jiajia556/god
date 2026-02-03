package service

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GetFileByRoute(route string) (filePath, fileName string, err error)
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
	out, err := RunCommandOutput(name, args...)
	if err != nil {
		// 包含命令输出便于排查
		if out != "" {
			OutputFatal(fmt.Sprintf("Command %s failed: %v\nOutput:\n%s", name, err, out))
		} else {
			OutputFatal(fmt.Sprintf("Command %s failed: %v", name, err))
		}
	}
}

func RunCommandOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	if CmdDir != "" {
		cmd.Dir = CmdDir
	}
	if len(GoEnv) > 0 {
		cmd.Env = append(os.Environ(), GoEnv...)
	}

	verbose := os.Getenv("GOD_VERBOSE") == "1"
	if verbose {
		// 在交互式/调试场景下直接把输出流到控制台
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", err
		}
		return "", nil
	}

	// 默认模式：捕获并返回 CombinedOutput，便于日志/错误中包含详细信息
	outputBytes, err := cmd.CombinedOutput()
	return string(outputBytes), err
}
