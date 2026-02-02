# god

**god** is a command-line scaffolding and code generation tool for Go projects.
It helps you quickly bootstrap project structures, generate models from SQL, automatically build routes from controller annotations, and supports build, formatting, and common maintenance tasks.

Repository: [jiajia556/god](https://github.com/jiajia556/god)

---

## Table of Contents

* Overview
* Features
* Installation
* Quick Start
* Commands & Options
* Templates & Output
* Controller Annotations (Auto Routing)
* gopackage.json (Project Metadata)
* Examples
* Development & Testing
* FAQ & Troubleshooting
* Contributing
* License

---

## Overview

**god** is designed to simplify and accelerate Go project initialization and common code-generation workflows. Its core capabilities include:

* Initializing project structures from embedded templates (`god init`)
* Generating controllers, actions, middleware, and models (`god gen`)
* Generating Go structs (models) from SQL `CREATE TABLE`
* Automatically generating router bindings by parsing controller annotations
* Running post-generation tasks such as `goimports` and `go mod tidy`

The built-in project template includes common dependencies such as **Gin**, **GORM**, **Viper**, and **Zap**, making it easy to start an API service quickly.

---

## Features

* Project initialization (`god init`)
* Code generation (`god gen ctrl|act|mdw|model`)
* Automatic route generation (`god mkrt`)
* Build & cross-compilation (`god build`)
* SQL → Model generation
* Embedded and customizable templates (`templates/basic`)

---

## Installation

Build from source:

```bash
git clone https://github.com/jiajia556/god.git
cd god
go build -o god .
# Optional: move the binary into your PATH
mv god /usr/local/bin/
```

It is recommended to use the same Go version as the template (Go **1.24**), but this can be adjusted if needed.

---

## Quick Start

Initialize a project:

```bash
god init github.com/yourname/myapp
cd myapp
```

Generate a controller:

```bash
god gen ctrl user list create update
```

Generate models from SQL:

```bash
god gen model --sql-path ./schema.sql
```

Generate routes from controller annotations:

```bash
god mkrt --root app/api/home
```

Build the service:

```bash
god build api user-service --version v1.0.0 --goos linux --goarch amd64
```

---

## Commands & Options

Common options:

* `--api-root, -a`：API root path (e.g. `api/v1` or `app/api/home`)
* `--sql-path, -s`：SQL file path
* `--app-root, -r`：Application root path (e.g. `app`)
* `--version, -v`：Version string (e.g. `v1.0.0`)
* `--goos, -o`：Target GOOS (e.g. `linux`)
* `--goarch, -g`：Target GOARCH (e.g. `amd64`)

Run `god <command> --help` to see detailed usage and examples.

---

## Templates & Output

Templates are located in `templates/basic`, including:

* `go.mod.tmpl` – module definition and dependencies
* `gopackage.json.tmpl` – project metadata
* `config/config.go.tmpl` – configuration loader (Viper + YAML/JSON)
* `app/api/home/router.go.tmpl` – auto-generated router template
* `app/api/home/main.go.tmpl` – API service entry point
* Controller, model, and middleware templates

During initialization, these templates are rendered and written as real files into the target project directory.

---

## Controller Annotations (Auto Routing)

`makerouter` uses `go/ast` to parse controller source code and extract annotations to generate route bindings.

Conventions:

* Controller type names must end with `Controller` (e.g. `UserController`)
* Controller methods (with receivers) are treated as actions
* Annotations are placed in comments directly above methods

Supported annotations:

* `@http_method <METHOD>` – HTTP method (default: `POST`)
* `@middleware <name1 name2 ...>` – middleware names (space-separated).
  Middleware must be implemented and exported in `lib/middleware`.

Example:

```go
package controller

type UserController struct{}

// @http_method GET
// @middleware auth logging
func (c UserController) List(ctx Context) {
    // ...
}

// @http_method POST
func (c UserController) Create(ctx Context) {
    // ...
}
```

The generated router will:

* Instantiate each controller
* Bind routes based on annotations
* Attach declared middleware in order

---

## gopackage.json (Project Metadata)

The `gopackage.json.tmpl` template generates a file like this:

```json
{
  "project_name": "github.com/yourname/myapp",
  "default_app_root": "app",
  "default_api_root": "app/api/home",
  "default_goos": "linux",
  "default_goarch": "amd64"
}
```

The CLI reads defaults from `./gopackage.json` in the current working directory.

> ⚠️ Commands should be run from the project root, or the tool must be extended to support explicit project paths.

---

## Examples

Initialize and run a demo project:

```bash
god init github.com/yourname/todo
cd todo
god gen ctrl todo list create update delete
god mkrt --root app/api/home
go run app/api/home/main.go --config ./config.yaml
```

Generate models from SQL:

```bash
god gen model --sql-path ./schema.sql
```

Build with cross-compilation:

```bash
god build api todo --version v0.1.0 --goos linux --goarch amd64
```

---

## Development & Testing

It is recommended to add test coverage for:

* SQL → struct parser (various `CREATE TABLE` edge cases)
* AST annotation extractor (route generation, method matching)
* Template rendering and file generation logic

Suggested CI checks:

```bash
go vet
go test ./...
gofmt
# or golangci-lint
```

---

## FAQ & Troubleshooting

**gopackage.json not found**
Make sure you run commands from the project root, or add support for specifying a project path.

**goimports not found**
The tool attempts to auto-install `goimports`, but this requires network access and proper Go environment configuration.

**Routes not generated or methods ignored**
Check the following:

* Controller files are under the API root directory
* Controller type names end with `Controller`
* Method annotations are correctly formatted and placed above the method

---

## Contributing

Contributions are welcome!

Typical workflow:

1. Fork the repository and create a feature branch
   `git checkout -b feat/mychange`
2. Implement changes and add tests
3. Submit a PR describing the motivation and impact

Please follow Go coding conventions and add unit tests for parsing and generation logic.

---

## License

This project is licensed under the **MIT License**.
See the LICENSE file for details.

---