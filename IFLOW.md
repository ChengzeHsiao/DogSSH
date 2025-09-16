# DogSSH 项目概述

DogSSH 是一个基于终端的交互式 SSH 管理工具，灵感来源于 lazydocker 和 k9s。它提供了一个干净、键盘驱动的用户界面 (TUI)，用于快速导航、连接、管理和在本地机器与 `~/.ssh/config` 文件中定义的服务器之间传输文件。

## 技术栈

- **语言**: Go (1.24.6)
- **框架/库**:
  - `github.com/gdamore/tcell/v2`：用于构建终端用户界面。
  - `github.com/rivo/tview`：基于 tcell 构建的终端 UI 组件库。
  - `github.com/spf13/cobra`：用于构建命令行应用。
  - `go.uber.org/zap`：用于日志记录。
  - `github.com/kevinburke/ssh_config` (已替换为 `github.com/adembc/ssh_config`)：用于解析和操作 SSH 配置文件。
  - `github.com/atotto/clipboard`：用于复制 SSH 命令到剪贴板。

## 项目架构

项目遵循简洁架构 (Clean Architecture) 原则，分为以下主要部分：

1.  **`cmd/`**: 包含应用程序的入口点 `main.go`。它负责初始化依赖项（如日志记录器、SSH 配置仓库、服务器服务和 TUI）并启动 Cobra 命令。
2.  **`internal/core/`**: 核心业务逻辑。
    *   `domain/`: 定义核心实体，如 `Server` 结构体。
    *   `ports/`: 定义核心服务和数据仓库的接口，实现依赖倒置。
    *   `services/`: 实现核心业务逻辑，如服务器列表、添加、编辑、删除、SSH 连接、Ping 等。`ServerService` 是主要的服务。
3.  **`internal/adapters/`**: 适配器层，负责与外部系统交互。
    *   `data/ssh_config_file/`: 实现 `ports.ServerRepository` 接口，用于读取和写入用户的 `~/.ssh/config` 文件。它还管理元数据（如最后连接时间、固定状态）和加密存储的密码。
    *   `ui/`: 实现基于 `tview` 的终端用户界面。`tui.go` 是主界面控制器，其他文件如 `server_list.go`, `server_form.go`, `search_bar.go` 等定义了具体的 UI 组件。

## 构建和运行

项目使用 `make` 作为构建工具。关键命令包括：

- `make build`: 构建项目二进制文件到 `./bin/dogssh`。
- `make run`: 直接从源代码运行应用。
- `make test`: 运行单元测试。
- `make fmt`: 格式化 Go 代码。
- `make vet`: 运行 `go vet`。
- `make lint`: 运行 `golangci-lint` 检查代码质量。
- `make quality`: 执行格式化、检查和 linting。

## 开发约定

- **依赖管理**: 使用 Go modules (`go.mod`, `go.sum`)。
- **代码质量**: 通过 `golangci-lint`, `gofumpt`, `staticcheck` 等工具保证代码质量。
- **测试**: 鼓励编写单元测试，使用 `make test` 运行。
- **日志**: 使用 `zap` 进行结构化日志记录。
- **配置管理**: 通过解析和操作 `~/.ssh/config` 文件进行服务器管理，确保与现有 SSH 设置兼容。
- **安全性**: 
  - 不存储、传输或修改私钥、密码等敏感信息。
  - 使用系统原生的 `ssh` 客户端进行连接。
  - 对写入 `~/.ssh/config` 的更改采用非破坏性方式，并创建备份。
  - 密码（如果使用）会加密后存储在单独的文件中。

## 安装

- **Homebrew (macOS)**: `brew install Adembc/homebrew-tap/dogssh`
- **从 Releases 下载**: 从 [GitHub Releases](https://github.com/ChengzeHsiao/dogssh/releases) 下载适用于您系统的二进制文件。
- **从源码构建**: 克隆仓库后运行 `make build` 或 `make run`。