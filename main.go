package main

import (
	"embed"
	"github.com/jiajia556/god/internal/cmd"
)

//go:embed templates/basic/*
var templateFS embed.FS

func main() {
	cmd.Execute(templateFS)
}
