package addmiddleware

import (
	"github.com/jiajia556/god/internal/service"
	"github.com/jiajia556/god/internal/template"
)

func AddMiddleware(middlewareTmpl string, middlewares []string) {
	for _, middleware := range middlewares {
		middlewareName := service.CapitalizeFirstLetter(middleware)

		filePath := "lib/middleware/" + middlewareName + ".go"
		err := template.CreateFile(middlewareTmpl, template.MiddlewareNameData{middlewareName}, filePath)
		if err != nil {
			service.OutputFatal(err)
		}
	}
}
