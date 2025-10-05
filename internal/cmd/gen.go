package cmd

import (
	"github.com/jiajia556/god/internal/cmd/addaction"
	"github.com/jiajia556/god/internal/cmd/addcontroller"
	"github.com/jiajia556/god/internal/cmd/addmiddleware"
	"github.com/jiajia556/god/internal/cmd/makemodel"
	"github.com/jiajia556/god/internal/service"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate Go code",
	Long:  `Generate Go code for controllers, models, middleware etc.`,
}

var ctrlCmd = &cobra.Command{
	Use:     "ctrl [controller-route] [actions...]",
	Short:   "Create a new controller with optional actions",
	Long:    "Generates a new controller file with specified route and optional initial actions",
	Example: "  god gen ctrl user\n  god gen ctrl product list create update",
	Args:    cobra.MinimumNArgs(1), // Requires at least 1 argument
	Run: func(cmd *cobra.Command, args []string) {
		// Read controller template from embedded files
		content, err := templateFS.ReadFile("templates/basic/app/api/home/controller.tmpl")
		if err != nil {
			service.OutputFatal(err)
		}

		// Extract actions from arguments
		var actions []string
		if len(args) > 1 {
			actions = args[1:]
		}

		// Get API root path from flag
		apiRoot, _ := cmd.Flags().GetString("api-root")
		addcontroller.AddController(string(content), apiRoot, args[0], actions)
	},
}

// actionCmd handles adding actions to existing controllers
var actionCmd = &cobra.Command{
	Use:     "act [controller-route] [actions...]",
	Short:   "Add actions to an existing controller",
	Long:    "Adds one or more action methods to a specified controller",
	Example: "  god gen act user getInfo\n  god gen act product search filter",
	Args:    cobra.MinimumNArgs(2), // Requires at least controller route and one action
	Run: func(cmd *cobra.Command, args []string) {
		apiRoot, _ := cmd.Flags().GetString("api-root")
		addaction.AddAction(apiRoot, args[0], args[1:])
	},
}

// middlewareCmd handles middleware creation
var middlewareCmd = &cobra.Command{
	Use:     "mdw [middleware-name...]",
	Short:   "Create new middleware components",
	Long:    "Generates middleware files with specified names",
	Example: "  god gen mdw auth\n  god gen mdw logging cache",
	Args:    cobra.MinimumNArgs(1), // Requires at least 1 middleware name
	Run: func(cmd *cobra.Command, args []string) {
		// Read middleware template from embedded files
		content, err := templateFS.ReadFile("templates/basic/lib/middleware/middleware.tmpl")
		if err != nil {
			service.OutputFatal(err)
		}
		addmiddleware.AddMiddleware(string(content), args)
	},
}

// modelCmd handles model file generation
var modelCmd = &cobra.Command{
	Use:     "model",
	Short:   "Generate database model files",
	Long:    "Generate Go model files from SQL schema definitions.\nCreates record and list type files based on SQL CREATE TABLE statements.",
	Example: "  god gen model --sql-path schema.sql\n  god gen model -s ./database/schema.sql",
	Run: func(cmd *cobra.Command, args []string) {
		recordContent, err := templateFS.ReadFile("templates/basic/model/record.go.tmpl")
		if err != nil {
			service.OutputFatal(err)
		}
		listContent, err := templateFS.ReadFile("templates/basic/model/list.go.tmpl")
		if err != nil {
			service.OutputFatal(err)
		}
		sqlPath, _ := cmd.Flags().GetString("sql-path")
		makemodel.MakeModel(sqlPath, string(recordContent), string(listContent))
	},
}
