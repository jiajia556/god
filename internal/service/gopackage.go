package service

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type GoPackage struct {
	inited         bool
	ProjectName    string `json:"project_name"`
	DefaultAppRoot string `json:"default_app_root"`
	DefaultApiRoot string `json:"default_api_root"`
	DefaultGOOS    string `json:"default_goos"`
	DefaultGOARCH  string `json:"default_goarch"`
}

var (
	goPackage   GoPackage
	projectRoot string // the directory where gopackage.json or go.mod was found
	mu          sync.Mutex
)

// initGoPackage locates and loads gopackage.json or falls back to go.mod module.
// Behavior:
// 1. If env GOD_PROJECT_ROOT is set, try that directory first.
// 2. Otherwise, walk up from cwd searching for gopackage.json; if not found, use go.mod to derive module.
// On success, sets goPackage and projectRoot.
func initGoPackage() error {
	mu.Lock()
	defer mu.Unlock()

	if goPackage.inited {
		return nil
	}

	// 1) If explicit env provided, check it first
	//if root := os.Getenv("GOD_PROJECT_ROOT"); root != "" {
	//	pkgPath := filepath.Join(root, "gopackage.json")
	//	if err := loadFromFileIfExists(pkgPath); err == nil {
	//		projectRoot = filepath.Clean(root)
	//		goPackage.inited = true
	//		return nil
	//	}
	//	// try go.mod in that dir as fallback
	//	modPath := filepath.Join(root, "go.mod")
	//	if exists(modPath) {
	//		if err := loadFromGoMod(modPath); err == nil {
	//			projectRoot = filepath.Clean(root)
	//			goPackage.inited = true
	//			return nil
	//		}
	//	}
	//}

	// 2) Walk up from cwd to root looking for gopackage.json or go.mod
	startDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get working directory: %w", err)
	}

	var triedPaths []string
	dir := startDir
	for {
		tryPkg := filepath.Join(dir, "gopackage.json")
		triedPaths = append(triedPaths, tryPkg)
		if err := loadFromFileIfExists(tryPkg); err == nil {
			projectRoot = filepath.Clean(dir)
			goPackage.inited = true
			return nil
		}

		// If gopackage.json not found, check go.mod as fallback
		tryMod := filepath.Join(dir, "go.mod")
		triedPaths = append(triedPaths, tryMod)
		if exists(tryMod) {
			if err := loadFromGoMod(tryMod); err == nil {
				projectRoot = filepath.Clean(dir)
				goPackage.inited = true
				return nil
			}
			// continue walking up if parsing fails
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached filesystem root
		}
		dir = parent
	}

	return fmt.Errorf("could not find gopackage.json nor parse go.mod module; attempted: %s", strings.Join(triedPaths, "; "))
}

// loadFromFileIfExists tries to read and unmarshal the given path if it exists.
// On success, it fills goPackage (but does not set projectRoot - caller must set it).
func loadFromFileIfExists(path string) error {
	if !exists(path) {
		return fmt.Errorf("not found: %s", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &goPackage); err != nil {
		return fmt.Errorf("unmarshal %s: %w", path, err)
	}
	ensureDefaults(&goPackage)
	return nil
}

// loadFromGoMod parses module name from go.mod and fills goPackage.ProjectName.
// Caller should set projectRoot on success.
func loadFromGoMod(modPath string) error {
	f, err := os.Open(modPath)
	if err != nil {
		return fmt.Errorf("open go.mod %s: %w", modPath, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				goPackage.ProjectName = parts[1]
				ensureDefaults(&goPackage)
				return nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan go.mod %s: %w", modPath, err)
	}
	return errors.New("module directive not found in go.mod")
}

func ensureDefaults(gp *GoPackage) {
	if gp.DefaultAppRoot == "" {
		gp.DefaultAppRoot = "app"
	}
	if gp.DefaultApiRoot == "" {
		// default relative api root under app
		gp.DefaultApiRoot = filepath.ToSlash(filepath.Join(gp.DefaultAppRoot, "api", "home"))
	}
	if gp.DefaultGOOS == "" {
		gp.DefaultGOOS = "linux"
	}
	if gp.DefaultGOARCH == "" {
		gp.DefaultGOARCH = "amd64"
	}
}

// exists reports whether the named file exists (and is not a directory).
func exists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// resolvePath makes cfgPath absolute based on projectRoot if cfgPath is not already absolute.
// If projectRoot is unknown, resolves relative to current working directory.
func resolvePath(cfgPath string) (string, error) {
	if filepath.IsAbs(cfgPath) {
		return filepath.Clean(cfgPath), nil
	}
	// ensure gopackage initialized so projectRoot may be set
	if !goPackage.inited {
		if err := initGoPackage(); err != nil {
			// fallback: join with cwd
			cwd, _ := os.Getwd()
			return filepath.Clean(filepath.Join(cwd, cfgPath)), nil
		}
	}
	base := projectRoot
	if base == "" {
		// fallback to cwd
		cwd, _ := os.Getwd()
		base = cwd
	}
	return filepath.Clean(filepath.Join(base, cfgPath)), nil
}

// GetDefaultAppRoot returns the absolute path to the app root based on discovered project root.
func GetDefaultAppRoot() (string, error) {
	if !goPackage.inited {
		if err := initGoPackage(); err != nil {
			return "", err
		}
	}
	return resolvePath(goPackage.DefaultAppRoot)
}

// GetDefaultApiRoot returns the absolute path to the api root based on discovered project root.
func GetDefaultApiRoot() (string, error) {
	if !goPackage.inited {
		if err := initGoPackage(); err != nil {
			return "", err
		}
	}
	return resolvePath(goPackage.DefaultApiRoot)
}

func GetDefaultGOOS() (string, error) {
	if !goPackage.inited {
		if err := initGoPackage(); err != nil {
			return "", err
		}
	}
	return goPackage.DefaultGOOS, nil
}

func GetDefaultGOARCH() (string, error) {
	if !goPackage.inited {
		if err := initGoPackage(); err != nil {
			return "", err
		}
	}
	return goPackage.DefaultGOARCH, nil
}

func GetProjectName() (string, error) {
	if !goPackage.inited {
		if err := initGoPackage(); err != nil {
			return "", err
		}
	}
	if goPackage.ProjectName == "" {
		return "", errors.New("project name is empty in gopackage.json or go.mod")
	}
	return goPackage.ProjectName, nil
}

// GetProjectRoot returns the absolute path of the discovered project root (where gopackage.json or go.mod was found).
// If not yet initialized, it will attempt initialization.
func GetProjectRoot() (string, error) {
	if !goPackage.inited {
		if err := initGoPackage(); err != nil {
			return "", err
		}
	}
	if projectRoot == "" {
		// fallback to current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return filepath.Clean(cwd), nil
	}
	return projectRoot, nil
}
