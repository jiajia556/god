package makemodel

import (
	"fmt"
	"github.com/jiajia556/god/internal/service"
	"github.com/jiajia556/god/internal/template"
	"os"
	"path/filepath"
	"strings"
)

// MakeModel generates model files from SQL CREATE TABLE statements
// Parameters:
//   - sqlFilePath:  Path to SQL file containing table definitions
//   - recordTmpl:   Content of template for record generation
//   - listTmpl:     Content of template for list type generation
func MakeModel(sqlFilePath, recordTmpl, listTmpl string) {
	defer runPostGenerationTasks()

	if sqlFilePath == "" {
		service.OutputFatal("SQL file path is required")
	}

	sqls, err := service.ExtractCreateTables(sqlFilePath)
	if err != nil {
		service.OutputFatal("Error extracting SQL statements: ", err.Error())
	}

	for _, sql := range sqls {
		GenerateModelFromSQL(sql, recordTmpl, listTmpl)
	}
}

// GenerateModelFromSQL creates model files for a single SQL CREATE TABLE statement
func GenerateModelFromSQL(sql, recordTmpl, listTmpl string) {
	// Generate model structure from SQL
	structText, structName, err := service.GenerateModelStruct(sql)
	if err != nil {
		service.OutputFatal(fmt.Sprintf("Error generating model struct: %v", err))
		return
	}

	// Prepare model package name
	modelPkg := strings.ToLower(structName)

	// Generate record file
	generateModelFile(modelPkg, structName, structText, recordTmpl, "record.go")

	// Generate list file
	generateModelFile(modelPkg, structName, structText, listTmpl, "list.go")
}

// runPostGenerationTasks executes post-processing commands
func runPostGenerationTasks() {
	service.RunCommand("goimports", "-w", ".")
	service.RunCommand("go", "mod", "tidy")
}

// generateModelFile handles file creation logic for model components
func generateModelFile(modelPkg, structName, structText, templatePath, fileName string) {
	// Set up file paths
	filePath := filepath.Join("model", modelPkg, fileName)

	// Skip if file already exists
	if service.FileExists(filePath) {
		return
	}

	// Prepare template data
	projectName, err := service.GetProjectName()
	if err != nil {
		service.OutputFatal(fmt.Sprintf("Error getting project name: %v", err))
	}

	data := template.ModelData{
		ModelPkg:        modelPkg,
		ProjectName:     projectName,
		ModelStruct:     structText,
		ModelStructName: structName,
	}

	// Create directory structure
	dir := filepath.Dir(filePath)
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		service.OutputFatal(fmt.Sprintf("Error creating directory %s: %v", dir, err))
		return
	}

	// Generate file from template
	if err = template.CreateFile(templatePath, data, filePath); err != nil {
		service.OutputFatal(fmt.Sprintf("Error creating %s: %v", fileName, err))
	}
}
