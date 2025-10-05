# God - Go 开发加速工具

一个用于加速 Go Web 应用程序开发的 CLI 工具，支持代码生成和项目脚手架。

## 安装

安装 `god`，运行以下命令：

```bash
go install github.com/jiajia556/god@latest
```

## 命令

### `god init [项目名称]`

初始化一个新项目，包含指定名称和基本结构。

**示例：**
```bash
god init myproject
god init example.com/myapp
```

### `god gen`

生成控制器、模型、中间件等 Go 代码。

#### 子命令：

1. **`god gen ctrl [控制器路由] [动作...]`**
   - 创建一个新的控制器，可选添加动作。
   
   **示例：**
   ```bash
   god gen ctrl user
   god gen ctrl product list create update
   ```

2. **`god gen act [控制器路由] [动作...]`**
   - 向现有控制器添加动作。
   
   **示例：**
   ```bash
   god gen act user getInfo
   god gen act product search filter
   ```

3. **`god gen mdw [中间件名称...]`**
   - 创建新的中间件组件。
   
   **示例：**
   ```bash
   god gen mdw auth
   god gen mdw logging cache
   ```

4. **`god gen model`**
   - 从 SQL 模式定义生成数据库模型文件。
   
   **示例：**
   ```bash
   god gen model --sql-path schema.sql
   god gen model -s ./database/schema.sql
   ```

### `god mkrt`

基于现有控制器生成 API 路由配置。

**示例：**
```bash
god mkrt --root api
```

### `god build [应用名称]`

构建应用程序组件，支持可选版本控制。

**示例：**
```bash
god build api user-service
god build admin-console --version v1.2.0
god build payment-service --app-root services --api-root api/v1
```

## 标志

- **`--api-root` (`-a`)**
  - API 根路径（例如 `api/v1`）。

- **`--sql-path` (`-s`)**
  - 包含表定义的 SQL 文件路径。

- **`--app-root` (`-r`)**
  - 应用根路径（例如 `app`）。

- **`--version` (`-v`)**
  - 应用版本（例如 `v1.0.0`）。

- **`--goos` (`-o`)**
  - GOOS（例如 `linux`）。

- **`--goarch` (`-g`)**
  - GOARCH（例如 `amd64`）。

## 许可证

本项目采用 MIT 许可证。
