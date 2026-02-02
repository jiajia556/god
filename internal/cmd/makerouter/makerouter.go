// Package makerouter provides functionality for automatic route generation
// based on controller annotations and directory structure analysis.
package makerouter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jiajia556/god/internal/service"
	"github.com/jiajia556/god/internal/template"
)

const (
	generatedFileName    = "router.go"    // Output filename for generated router
	controllerSuffix     = "Controller"   // Suffix for controller type names
	controllerDirName    = "controller"   // Standard directory name for controllers
	httpMethodAnnotation = "@http_method" // Annotation prefix for HTTP methods
	middlewareAnnotation = "@middleware"  // Annotation prefix for middlewares
)

// routeGenerator maintains state during route generation process
type routeGenerator struct {
	imports           []string          // Import paths for controller packages
	initRegistrations []string          // Controller registration statements
	pkgAliases        map[string]string // Package import aliases
	httpMethods       map[string]string // HTTP method mappings
	middlewares       map[string]string // Middleware configurations
	projectName       string            // Current project module name
	projectRoot       string            // Current project root directory
}

// MakeRouter initiates the route generation process
func MakeRouter(routerTemplate string, rootPath string) {
	if rootPath == "" {
		var err error
		if rootPath, err = service.GetDefaultApiRoot(); err != nil {
			service.OutputFatal("Failed to get API root:", err)
		}
	}

	rg := &routeGenerator{
		pkgAliases:  make(map[string]string),
		httpMethods: make(map[string]string),
		middlewares: make(map[string]string),
	}

	var err error
	if rg.projectName, err = service.GetProjectName(); err != nil {
		service.OutputFatal("Failed to get project name:", err)
	}

	// Get project root (where gopackage.json or go.mod was discovered)
	if rg.projectRoot, err = service.GetProjectRoot(); err != nil {
		service.OutputFatal("Failed to get project root:", err)
	}

	tmplData, err := rg.generateTemplateData(rootPath)
	if err != nil {
		service.OutputFatal("Template data generation failed:", err)
	}

	outputPath := filepath.Join(rootPath, generatedFileName)
	err = template.CreateFile(routerTemplate, tmplData, outputPath)
	if err != nil {
		service.OutputFatal(err)
	}
}

// generateTemplateData collects and prepares data for template generation
func (rg *routeGenerator) generateTemplateData(root string) (template.RouterTmplData, error) {
	if err := rg.analyzeProjectStructure(root); err != nil {
		return template.RouterTmplData{}, fmt.Errorf("project analysis failed: %w", err)
	}

	return template.RouterTmplData{
		ApiRootDirName:        filepath.Base(root),
		HTTPMethodTags:        rg.formatHTTPMethods(),
		MiddlewareTags:        rg.formatMiddlewares(),
		RegisterControllers:   strings.Join(rg.initRegistrations, ""),
		MiddlewareImportPath:  rg.middlewareImport(),
		ControllersImportPath: strings.Join(rg.imports, "\n\t"),
	}, nil
}

// analyzeProjectStructure walks through project directories to find controllers
func (rg *routeGenerator) analyzeProjectStructure(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Printf("Directory access error: %v", err)
			return nil
		}

		if d.IsDir() && d.Name() == controllerDirName {
			if err := rg.processControllerPackage(path); err != nil {
				return fmt.Errorf("controller processing failed: %w", err)
			}
			return filepath.SkipDir
		}
		return nil
	})
}

// processControllerPackage processes all Go files in a controller directory
func (rg *routeGenerator) processControllerPackage(dirPath string) error {
	return filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".go") {
			return err
		}
		return rg.analyzeControllerFile(path)
	})
}

// analyzeControllerFile parses a single Go file for controller definitions
func (rg *routeGenerator) analyzeControllerFile(filePath string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("file parsing failed: %w", err)
	}

	pkgPath := constructImportPath(rg.projectName, rg.projectRoot, filePath)
	alias, exists := rg.pkgAliases[pkgPath]

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || !strings.HasSuffix(typeSpec.Name.Name, controllerSuffix) {
				continue
			}

			controllerName := typeSpec.Name.Name
			if !exists {
				alias = fmt.Sprintf("controller%d", len(rg.imports))
				// ensure alias uniqueness in rare case
				for _, a := range rg.imports {
					if strings.Contains(a, "\""+pkgPath+"\"") && rg.pkgAliases[pkgPath] == alias {
						alias = fmt.Sprintf("controller%d", len(rg.imports)+1)
					}
				}
				rg.pkgAliases[pkgPath] = alias
				rg.imports = append(rg.imports, fmt.Sprintf("\t%s \"%s\"", alias, pkgPath))
			}

			fullTypeName := fmt.Sprintf("%s.%s", alias, controllerName)
			rg.initRegistrations = append(rg.initRegistrations,
				fmt.Sprintf("\n\tRegisterController(%s{})", fullTypeName))

			rg.extractAnnotations(node, controllerName, pkgPath+"."+controllerName)
		}
	}
	return nil
}

// extractAnnotations parses controller method annotations
func (rg *routeGenerator) extractAnnotations(node *ast.File, typeName, pkgPrefix string) {
	ast.Inspect(node, func(n ast.Node) bool {
		fnDecl, ok := n.(*ast.FuncDecl)
		if !ok || fnDecl.Recv == nil || len(fnDecl.Recv.List) == 0 {
			return true
		}

		recvType := extractReceiverType(fnDecl.Recv.List[0].Type)
		if recvType != typeName {
			return true
		}

		annotationKey := fmt.Sprintf("%s.%s", pkgPrefix, fnDecl.Name.Name)
		rg.processMethodAnnotations(fnDecl, annotationKey)
		return true
	})
}

// Helper functions below maintain the same logic with improved readability
func constructImportPath(projectName, projectRoot, filePath string) string {
	// Normalize to absolute slash-separated paths
	absFilePath, _ := filepath.Abs(filePath)
	absFilePath = filepath.ToSlash(absFilePath)

	absProjectRoot := projectRoot
	if absProjectRoot == "" {
		if p, err := filepath.Abs("."); err == nil {
			absProjectRoot = p
		} else {
			absProjectRoot = ""
		}
	}
	absProjectRoot = filepath.ToSlash(absProjectRoot)

	// Directory containing the file
	dir := filepath.ToSlash(filepath.Dir(absFilePath))

	// Compute relative path from project root to the file's directory
	rel := dir
	if absProjectRoot != "" {
		if r, err := filepath.Rel(absProjectRoot, dir); err == nil {
			rel = filepath.ToSlash(r)
		} else {
			// fallback: if Dir contains projectRoot as prefix, trim prefix
			if strings.HasPrefix(dir, absProjectRoot+"/") {
				rel = strings.TrimPrefix(dir, absProjectRoot+"/")
			}
		}
	}

	// Clean and trim
	rel = strings.Trim(rel, "/")

	// If relative path is empty, import is module root
	if rel == "" {
		return projectName
	}

	// Make sure rel uses slashes and no leading/trailing slash
	rel = strings.ReplaceAll(rel, "\\", "/")

	// Compose module import path
	return strings.TrimRight(projectName, "/") + "/" + rel
}

func extractReceiverType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	}
	return ""
}

func (rg *routeGenerator) processMethodAnnotations(fnDecl *ast.FuncDecl, key string) {
	if fnDecl.Doc == nil {
		rg.httpMethods[key] = "POST"
		return
	}
	for _, comment := range fnDecl.Doc.List {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		switch {
		case strings.HasPrefix(text, httpMethodAnnotation):
			method := strings.ToUpper(strings.TrimSpace(strings.TrimPrefix(text, httpMethodAnnotation)))
			if method != "" {
				rg.httpMethods[key] = method
			}
		case strings.HasPrefix(text, middlewareAnnotation):
			middlewares := strings.TrimSpace(strings.TrimPrefix(text, middlewareAnnotation))
			if middlewares != "" {
				rg.middlewares[key] = middlewares
			}
		}
	}
}

func (rg *routeGenerator) middlewareImport() string {
	if len(rg.middlewares) > 0 {
		return fmt.Sprintf("\t\"%s/lib/middleware\"", rg.projectName)
	}
	return ""
}

func (rg *routeGenerator) formatHTTPMethods() string {
	var builder strings.Builder
	for k, v := range rg.httpMethods {
		builder.WriteString(fmt.Sprintf("\t\t\"%s\": \"%s\",\n", k, v))
	}
	return builder.String()
}

func (rg *routeGenerator) formatMiddlewares() string {
	var builder strings.Builder
	for k, v := range rg.middlewares {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}

		components := strings.Split(v, " ")

		for i := 0; i < len(components); i++ {
			components[i] = "middleware." + strings.TrimSpace(components[i])
		}

		formatted := "{" + strings.Join(components, ", ") + "}"
		builder.WriteString(fmt.Sprintf("\t\t\"%s\": %s,\n", k, formatted))
	}
	return builder.String()
}
