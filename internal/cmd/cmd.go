// Package cmd provides command-line interface functionality for the god tool
package cmd

import (
	"embed"

	"github.com/jiajia556/god/internal/cmd/build"
	"github.com/jiajia556/god/internal/cmd/initproject"
	"github.com/jiajia556/god/internal/cmd/makerouter"
	"github.com/jiajia556/god/internal/service"
	"github.com/spf13/cobra"
)

// templateFS holds embedded template files for code generation
var templateFS embed.FS

// rootCmd is the base command for the CLI tool
var rootCmd = &cobra.Command{
	Use:   "god",
	Short: "God - Go Development Accelerator Tool",
	Long:  `A CLI tool to accelerate Go web application development with code generation and project scaffolding.`,
}

// initCmd handles project initialization
var initCmd = &cobra.Command{
	Use:     "init [project-name]",
	Short:   "Create a new project",
	Long:    "Initialize a new project with the specified name and basic structure",
	Example: "  god init myproject\n  god init example.com/myapp",
	Args:    cobra.ExactArgs(1), // Requires exactly 1 argument
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		initproject.InitProject(projectName, templateFS)
	},
}

// makeRouterCmd handles router file generation
var makeRouterCmd = &cobra.Command{
	Use:     "mkrt",
	Short:   "Generate API router configuration",
	Long:    "Creates or updates the main router file based on existing controllers",
	Example: "  god mkrt --root api",
	Run: func(cmd *cobra.Command, args []string) {
		// Read router template from embedded files
		content, err := templateFS.ReadFile("templates/basic/app/api/home/router.go.tmpl")
		if err != nil {
			service.OutputFatal(err)
		}

		// Get API root path from flag
		apiRoot, _ := cmd.Flags().GetString("api-root")
		makerouter.MakeRouter(string(content), apiRoot)
	},
}

// buildCmd handles app building
var buildCmd = &cobra.Command{
	Use:     "build [app-name]",
	Short:   "Build application components",
	Long:    "Build application components with optional versioning.\nFor API applications, use 'build api [app-name]'.\nFor regular applications, use 'build [app-name]'.",
	Example: "  god build api user-service\n  god build admin-console --version v1.2.0\n  god build payment-service --app-root services --api-root api/v1",
	Args:    cobra.RangeArgs(1, 2), // Accepts 1 or 2 arguments
	Run: func(cmd *cobra.Command, args []string) {
		content, err := templateFS.ReadFile("templates/basic/app/api/home/router.go.tmpl")
		if err != nil {
			service.OutputFatal(err)
		}
		app := args[0]
		isApi := false
		appRoot, _ := cmd.Flags().GetString("app-root")
		apiRoot, _ := cmd.Flags().GetString("api-root")
		version, _ := cmd.Flags().GetString("version")
		goos, _ := cmd.Flags().GetString("goos")
		goarch, _ := cmd.Flags().GetString("goarch")
		if args[0] == "api" {
			app = args[1]
			isApi = true
		}

		build.Build(string(content), app, appRoot, apiRoot, version, goos, goarch, isApi)
	},
}

// Execute initializes and runs the CLI application
func Execute(tmplFS embed.FS) {
	templateFS = tmplFS

	rootCmd.AddCommand(genCmd)
	genCmd.AddCommand(ctrlCmd)
	genCmd.AddCommand(actionCmd)
	genCmd.AddCommand(middlewareCmd)
	genCmd.AddCommand(modelCmd)

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(makeRouterCmd)
	rootCmd.AddCommand(buildCmd)

	// Configure persistent flags for relevant commands
	for _, cmd := range []*cobra.Command{ctrlCmd, actionCmd, makeRouterCmd, buildCmd} {
		cmd.Flags().StringP("api-root", "a", "", "API root path (e.g., 'api/v1')")
	}
	modelCmd.Flags().StringP("sql-path", "s", "", "Path to SQL file containing table definitions")
	buildCmd.Flags().StringP("app-root", "r", "", "App root path (e.g., 'app')")
	buildCmd.Flags().StringP("version", "v", "", "App version (e.g., 'v1.0.0')")
	buildCmd.Flags().StringP("goos", "o", "", "GOOS (e.g., 'linux')")
	buildCmd.Flags().StringP("goarch", "g", "", "GOARCH (e.g., 'amd64')")

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		service.OutputFatal(err)
	}
}
