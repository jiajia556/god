# god

god 是一个面向 Go 项目的命令行脚手架与代码生成工具。它可以快速生成项目骨架、根据 SQL 生成 model、根据 controller 注释自动生成路由、并支持构建（build）与常见代码格式化/整理任务。

仓库: [jiajia556/god](https://github.com/jiajia556/god)

---

## 目录

- 概述
- 功能
- 安装
- 快速开始
- 命令与参数
- 模板与输出
- 控制器注释（自动路由）
- gopackage.json（项目元信息）
- 示例
- 开发与测试
- 常见问题与排查
- 贡献
- 许可证

---

## 概述

god 用于简化和加速 Go 项目的初始化与常见代码生成需求。主要能力：

- 基于内嵌模板初始化项目结构（`god init`）。
- 生成 controller、action、middleware、model 等代码模板（`god gen`）。
- 从 SQL（CREATE TABLE）生成 Go struct（模型）。
- 通过解析 controller 源码注释（@http_method、@middleware）生成 router 绑定代码。
- 自动执行 goimports、go mod tidy 等后处理命令。

项目模板包含常用依赖：Gin、GORM、Viper、Zap 等，便于快速启动一个 API 服务。

---

## 功能

- 初始化项目（`god init`）
- 代码生成（`god gen ctrl|act|mdw|model`）
- 路由自动生成（`god mkrt`）
- 构建组件（`god build`）
- SQL -> Model（`god gen model`）
- 嵌入模板（`templates/basic`），可定制并生成样例代码

---

## 安装

从源码构建：

```bash
git clone https://github.com/jiajia556/god.git
cd god
go build -o god .
# 可选：将可执行文件移动到 PATH 下
mv god /usr/local/bin/
```

建议使用与模板一致的 Go 版本（模板中为 go 1.24），但可根据需要调整。

---

## 快速开始

初始化项目：

```bash
god init github.com/yourname/myapp
cd myapp
```

生成控制器：

```bash
god gen ctrl user list create update
```

根据 SQL 生成 model：

```bash
god gen model --sql-path ./schema.sql
```

根据控制器注释生成路由：

```bash
god mkrt --root app/api/home
```

构建服务：

```bash
god build api user-service --version v1.0.0 --goos linux --goarch amd64
```

---

## 命令与参数

常见参数：

- `--api-root, -a`：API 根路径（例如 `api/v1` 或 `app/api/home`）
- `--sql-path, -s`：SQL 文件路径
- `--app-root, -r`：应用根路径（例如 `app`）
- `--version, -v`：版本号（例如 `v1.0.0`）
- `--goos, -o`：GOOS（例如 `linux`）
- `--goarch, -g`：GOARCH（例如 `amd64`）

运行 `god <command> --help` 查看子命令帮助与示例。

---

## 模板与输出

模板位于仓库的 `templates/basic`，包括：

- `go.mod.tmpl`：模块及依赖
- `gopackage.json.tmpl`：项目元数据（project_name、default_app_root 等）
- `config/config.go.tmpl`：配置解析（viper + yaml/json）
- `app/api/home/router.go.tmpl`：自动生成的 router 文件模板
- `app/api/home/main.go.tmpl`：API 服务入口模板
- 以及 controller、model、middleware 等模板

初始化项目会把这些模板渲染为真实文件写入目标目录。

---

## 控制器注释（自动路由）

makerouter 使用 `go/ast` 解析控制器源码并读取注释来生成路由。约定：

- 控制器类型名以 `Controller` 结尾，例如 `UserController`。
- 控制器方法（带接收者）为 action。
- 在方法上方使用注释指定 HTTP 方法与中间件。

支持注释格式（放在方法前）：

- `@http_method <METHOD>`：指定 HTTP 方法，默认为 `POST`。
- `@middleware <name1 name2 ...>`：指定中间件名称（空间分隔），中间件需在 `lib/middleware` 中实现并导出。

示例：

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

生成的 router 会为每个 controller 类型注册实例并根据注释设置方法映射和中间件。

---

## gopackage.json（项目元信息）

模板 `gopackage.json.tmpl` 生成如下结构示例：

```json
{
  "project_name": "github.com/yourname/myapp",
  "default_app_root": "app",
  "default_api_root": "app/api/home",
  "default_goos": "linux",
  "default_goarch": "amd64"
}
```

CLI 从当前工作目录的 `./gopackage.json` 读取默认值。请在项目根目录运行相关命令，或扩展工具以支持指定项目根路径。

---

## 示例

初始化并运行示例：

```bash
god init github.com/yourname/todo
cd todo
god gen ctrl todo list create update delete
god mkrt --root app/api/home
go run app/api/home/main.go --config ./config.yaml
```

从 SQL 生成模型：

```bash
god gen model --sql-path ./schema.sql
```

构建并交叉编译：

```bash
god build api todo --version v0.1.0 --goos linux --goarch amd64
```

---

## 开发与测试建议

建议为以下部分增加测试用例：

- SQL -> struct 解析器（各种 CREATE TABLE 边界情况）
- AST 注释提取器（makerouter 的注释解析、方法匹配）
- 模板渲染与文件生成逻辑

在 CI 中运行：`go vet`, `go test ./...`, `gofmt` 或 `golangci-lint`。

---

## 常见问题与排查

- 找不到 `gopackage.json`：请在项目根目录运行命令，或指定项目路径。
- 丢失 `goimports`：工具会尝试自动安装 `goimports`，但需要网络与 `go` 环境配置允许 `go install`。
- 路由未生成或方法未被识别：检查控制器文件是否位于 api root 下且类型名以 `Controller` 结尾；检查方法上方注释格式是否正确。

---

## 贡献

欢迎贡献！基本流程：

1. Fork 并创建分支：`git checkout -b feat/mychange`
2. 增加/修改代码并补充测试
3. 提交 PR 并描述变更理由与影响

请遵循 Go 语言风格与为解析/生成逻辑添加单元测试。

---

## 许可证

本项目采用 MIT 许可证，详见 LICENSE 文件。

---