# God - Go Development Accelerator Tool

God (Go Development Accelerator Tool) is a command-line tool designed to accelerate Go web application development. It provides a set of commands to generate code and scaffold projects, making it easier and faster to build web applications with Go.

## Installation

```bash
go install github.com/jiajia556/god@latest
```

Or clone the repository and build from source:

```bash
git clone https://github.com/jiajia556/god.git
cd god
go build -o god
```

## Usage

### Initialize a new project

```bash
god init myproject
```

This command creates a new project with the basic structure and necessary files.

### Add a controller

```bash
god addc user/user
```

This command adds a new controller with the path "user/user" to your project.

### Add an action to a controller

```bash
god adda user/user login
```

This command adds a new action named "login" to the "user/user" controller.

### Add middleware

```bash
god addm auth
```

This command adds a new middleware named "auth" to your project.

### Generate database model files

```bash
god mkmd --sql-path ./mydb.sql
```

This command extracts table creation statements from the mydb.sql file and creates corresponding models.

### Generate API router configuration

```bash
god mkrt
```

This command generates router configuration for your API endpoints.

### Build the application

```bash
god build api home
```

This command builds your application.

## Project Structure

When you initialize a new project with God, it creates the following structure:

```
myproject/
├── app/
│   └── api/
│       └── home/
│           ├── main.go
│           └── router.go
├── bin/
├── config/
│   └── config.go
├── lib/
│   ├── db/
│   │   └── mysql/
│   │       └── mysql.go
│   ├── middleware/
│   ├── mylog/
│   │   └── mylog.go
│   ├── mytime/
│   │   └── mytime.go
│   └── output/
│       ├── output.go
│       └── outputmsg.go
├── model/
├── go.mod
└── gopackage.json
```

## Features

- **Automatic Route Registration**: Routes are automatically registered based on controller methods.
- **Middleware Support**: Easy integration of middleware into your routes.
- **HTTP Method Annotations**: Define HTTP methods for your routes using annotations.
- **Unified Output Format**: Standardized output format for API responses.
- **Template-Based Code Generation**: Generate code from templates to ensure consistency.


## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
