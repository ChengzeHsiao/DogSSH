# LazySSH 项目概述

LazySSH 是一个基于终端的交互式 SSH 管理工具，灵感来源于 lazydocker 和 k9s。它提供了一个干净、键盘驱动的用户界面，用于快速浏览、连接、管理和在本地机器与定义在 `~/.ssh/config` 文件中的任何服务器之间传输文件。

## 技术栈与架构

- **语言**: Go (1.24.6)
- **UI 框架**: [tview](https://github.com/rivo/tview) 和 [tcell](https://github.com/gdamore/tcell)
- **命令行解析**: [cobra](https://github.com/spf13/cobra)
- **日志**: [zap](https://github.com/uber-go/zap)
- **SSH 配置解析**: [kevinburke/ssh_config](https://github.com/kevinburke/ssh_config) (定制版本)
- **架构模式**: 六边形架构 (Hexagonal Architecture)，分为核心域、端口和服务、适配器。
  - **核心域 (`internal/core/domain`)**: 定义了 `Server` 实体。
  - **端口 (`internal/core/ports`)**: 定义了 `ServerRepository` 和 `ServerService` 接口。
  - **服务 (`internal/core/services`)**: 实现了 `ServerService`，包含业务逻辑，如服务器验证、列表排序、SSH 连接、Ping 检查等。
  - **适配器**:
    - **UI (`internal/adapters/ui`)**: 实现了基于 `tview` 的终端用户界面，包括服务器列表、搜索栏、详情视图、表单等组件。
    - **数据 (`internal/adapters/data/ssh_config_file`)**: 实现了 `ServerRepository`，负责读取和写入用户的 `~/.ssh/config` 文件，并管理元数据（如访问次数、固定状态）。

## 构建、运行与测试

项目使用 `make` 进行构建和管理，相关命令定义在 `makefile` 中。

### 依赖管理
```bash
make deps       # 下载依赖
make tidy       # 整理依赖
```

### 代码质量
```bash
make fmt        # 格式化代码
make vet        # 运行 go vet
make lint       # 运行 golangci-lint
make check      # 运行 staticcheck
make quality    # 运行所有代码质量检查 (fmt, vet, lint)
```

### 构建
```bash
make build      # 构建二进制文件到 ./bin/lazyssh
make build-all  # 为所有平台构建二进制文件
```

### 运行
```bash
make run        # 直接从源码运行
make install    # 构建并安装二进制文件到 $GOBIN
```

### 测试
```bash
make test           # 运行单元测试
make test-verbose   # 运行单元测试并输出详细信息
make test-short     # 运行单元测试 (短模式)
make coverage       # 运行测试并生成覆盖率报告
make benchmark      # 运行基准测试
```

### 其他
```bash
make clean      # 清理构建产物和缓存
make help       # 显示帮助信息
```

## 开发约定

- **依赖注入**: 使用构造函数注入依赖，遵循接口隔离原则。
- **错误处理**: 使用 `zap` 进行结构化日志记录，关键错误会通过日志输出。
- **配置管理**: 所有配置通过 `~/.ssh/config` 文件进行管理，确保与系统原生 SSH 客户端兼容。
- **安全性**: 不存储、传输或修改私钥、密码等敏感信息，所有 SSH 连接通过系统原生 `ssh` 命令执行。
- **文件操作**: 对 `~/.ssh/config` 文件的修改是安全的，使用原子写入和备份机制，确保配置文件不会因程序异常而损坏。

## 目录结构

```
.
├── cmd                 # 程序入口
├── docs                # 文档和截图
├── internal            # 内部代码
│   ├── adapters        # 适配器层
│   │   ├── data        # 数据适配器 (SSH 配置文件)
│   │   └── ui          # UI 适配器 (TUI)
│   ├── core            # 核心业务逻辑
│   │   ├── domain      # 领域模型
│   │   ├── ports       # 端口接口
│   │   └── services    # 服务实现
│   └── logger          # 日志模块
├── go.mod              # Go 模块定义
├── go.sum              # Go 模块校验和
├── makefile            # 构建脚本
├── README.md           # 项目说明
└── IFLOW.md            # 此文件
```