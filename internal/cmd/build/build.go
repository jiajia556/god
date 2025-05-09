// Package build provides functionality for building Go applications
package build

import (
	"github.com/jiajia556/god/internal/cmd/makerouter"
	"github.com/jiajia556/god/internal/service"
	"path/filepath"
)

// Build compiles the application
// Parameters:
//   - routerTmpl: Path to router template file
//   - app:        Application name
//   - appRoot:    Root directory of the application
//   - apiRoot:    API root directory
//   - version:    Version number for the build
//   - goos:       Target operating system
//   - goarch:     Target architecture
//   - isApi:      Flag indicating if building an API application
func Build(routerTmpl, app, appRoot, apiRoot, version, goos, goarch string, isApi bool) {
	// Set default application root if not specified
	if appRoot == "" {
		var err error
		appRoot, err = service.GetDefaultAppRoot()
		if err != nil {
			service.OutputFatal(err)
		}
	}

	// Set default API root if not specified
	if apiRoot == "" {
		var err error
		apiRoot, err = service.GetDefaultApiRoot()
		if err != nil {
			service.OutputFatal(err)
		}
	}

	// Determine build path based on application type
	buildPath := filepath.Join(appRoot, app)
	if isApi {
		buildPath = filepath.Join(filepath.Dir(apiRoot), app)
		makerouter.MakeRouter(routerTmpl, apiRoot)
	}
	buildPath = "./" + buildPath

	// Set default target OS if not specified
	if goos == "" {
		var err error
		goos, err = service.GetDefaultGOOS()
		if err != nil {
			service.OutputFatal(err)
		}
	}

	// Set default target architecture if not specified
	if goarch == "" {
		var err error
		goarch, err = service.GetDefaultGOARCH()
		if err != nil {
			service.OutputFatal(err)
		}
	}

	// Construct output filename with version and extension
	outName := filepath.Join("bin", app)
	if version != "" {
		outName += "-v" + version
	}
	if goos == "windows" {
		outName += ".exe"
	}

	// Set build environment variables
	service.GoEnv = []string{
		"GOOS=" + goos,
		"GOARCH=" + goarch,
	}
	defer func() {
		service.GoEnv = []string{}
	}()

	// Execute the build command
	service.RunCommand("go", "build", "-o", outName, buildPath)
}
