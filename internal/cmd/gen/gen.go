package gen

import (
	"embed"

	"github.com/jiajia556/god/internal/cmd/addcontroller"
	"github.com/jiajia556/god/internal/service"
	"github.com/spf13/cobra"
)

// templateFS holds embedded template files for code generation
var templateFS embed.FS

func InitTemplateFS(tmplFS embed.FS) {
	templateFS = tmplFS
}

var GenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate Go code",
	Long:  `Generate Go code for controllers, models, middleware etc.`,
}

var ctrl = &cobra.Command{
	Use:     "ctrl [controller-route] [actions...]",
	Short:   "Create a new controller with optional actions",
	Long:    "Generates a new controller file with specified route and optional initial actions",
	Example: "  god addc user\n  god addc product list create update",
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
