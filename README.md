# God - Go Development Accelerator Tool

A CLI tool to accelerate Go web application development with code generation and project scaffolding.

## Installation

To install `god`, run the following command:

```bash
go install github.com/jiajia556/god@latest
```

## Commands

### `god init [project-name]`

Initialize a new project with the specified name and basic structure.

**Example:**
```bash
god init myproject
god init example.com/myapp
```

### `god gen`

Generate Go code for controllers, models, middleware, etc.

#### Subcommands:

1. **`god gen ctrl [controller-route] [actions...]`**
   - Create a new controller with optional actions.
   
   **Example:**
   ```bash
   god gen ctrl user
   god gen ctrl product list create update
   ```

2. **`god gen act [controller-route] [actions...]`**
   - Add actions to an existing controller.
   
   **Example:**
   ```bash
   god gen act user getInfo
   god gen act product search filter
   ```

3. **`god gen mdw [middleware-name...]`**
   - Create new middleware components.
   
   **Example:**
   ```bash
   god gen mdw auth
   god gen mdw logging cache
   ```

4. **`god gen model`**
   - Generate database model files from SQL schema definitions.
   
   **Example:**
   ```bash
   god gen model --sql-path schema.sql
   god gen model -s ./database/schema.sql
   ```

### `god mkrt`

Generate API router configuration based on existing controllers.

**Example:**
```bash
god mkrt --root api
```

### `god build [app-name]`

Build application components with optional versioning.

**Example:**
```bash
god build api user-service
god build admin-console --version v1.2.0
god build payment-service --app-root services --api-root api/v1
```

## Flags

- **`--api-root` (`-a`)**
  - API root path (e.g., `api/v1`).

- **`--sql-path` (`-s`)**
  - Path to SQL file containing table definitions.

- **`--app-root` (`-r`)**
  - App root path (e.g., `app`).

- **`--version` (`-v`)**
  - App version (e.g., `v1.0.0`).

- **`--goos` (`-o`)**
  - GOOS (e.g., `linux`).

- **`--goarch` (`-g`)**
  - GOARCH (e.g., `amd64`).

## License

This project is licensed under the MIT License.
