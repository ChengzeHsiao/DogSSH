<div align="center">
  <img src="./docs/logo.png" alt="dogssh logo" width="600" height="600"/>
</div>

---

# DogSSH

[![License](https://img.shields.io/github/license/chengzehsiao/dogssh)](LICENSE)
[![GitHub release](https://img.shields.io/github/v/release/chengzehsiao/dogssh)](https://github.com/chengzehsiao/dogssh/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/chengzehsiao/dogssh)](https://goreportcard.com/report/github.com/chengzehsiao/dogssh)

DogSSH 是一个基于终端的交互式 SSH 管理工具，灵感来源于 lazydocker 和 k9s 等工具，专为直接从终端管理服务器集群而设计。

使用 DogSSH，您可以快速导航、连接、管理本地计算机与 `~/.ssh/config` 文件中定义的任何服务器之间的文件传输。无需再记住 IP 地址或运行冗长的 scp 命令，只需一个干净、键盘驱动的用户界面。

---

## ✨ 核心功能

### 服务器管理
- 📜 从您的 `~/.ssh/config` 文件中读取并以可滚动列表的形式显示服务器。
- ➕ 通过 UI 添加新服务器，指定别名、主机/IP、用户名、端口和身份文件。
- ✏ 直接从 UI 编辑现有的服务器条目。
- 🗑 安全地删除服务器条目。
- 📌 固定/取消固定服务器，将收藏夹置顶。
- 🏓 Ping 服务器以检查状态。

### 快速服务器导航
- 🔍 按别名、IP 或标签进行模糊搜索。
- 🖥 一键 SSH 连接到所选服务器（Enter 键）。
- 🏷 为服务器添加标签（例如，prod、dev、test）以便快速筛选。
- ↕️ 按别名或上次 SSH 时间排序（切换 + 反向）。

### 安全性与配置安全
- 🔐 **无新增安全风险**：DogSSH 只是现有 `~/.ssh/config` 文件的 UI/TUI 包装器。所有 SSH 连接均使用系统原生的 ssh 二进制文件。
- 🛡️ **非破坏性编辑**：对 `~/.ssh/config` 的更改是最低限度的，并保留现有的注释、间距和顺序。
- 📦 **自动备份**：在进行任何更改之前，会创建一次性原始备份和滚动时间戳备份。

---

## 📦 安装指南

### 选项 1: Homebrew (macOS)

```bash
brew install chengzehsiao/homebrew-tap/dogssh
```

### 选项 2: 从 Releases 下载二进制文件

从 [GitHub Releases](https://github.com/chengzehsiao/dogssh/releases) 下载。您可以使用以下代码段自动获取适用于您的操作系统/架构（支持 Darwin/Linux 和 amd64/arm64）的最新版本：

```bash
# 检测最新版本
LATEST_TAG=$(curl -fsSL https://api.github.com/repos/chengzehsiao/dogssh/releases/latest | jq -r .tag_name)
# 下载适用于您系统的正确二进制文件
curl -LJO "https://github.com/chengzehsiao/dogssh/releases/download/${LATEST_TAG}/dogssh_$(uname)_$(uname -m).tar.gz"
# 解压二进制文件
tar -xzf dogssh_$(uname)_$(uname -m).tar.gz
# 移动到 /usr/local/bin 或 PATH 中的其他目录
sudo mv dogssh /usr/local/bin/
# 享受吧！
dogssh
```

### 选项 3: 从源代码构建

```bash
# 克隆仓库
git clone https://github.com/chengzehsiao/dogssh.git
cd dogssh

# 为 macOS 构建
make build
./bin/dogssh

# 或直接运行
make run
```

---

## 🚀 快速开始

1. 确保您的服务器已在 `~/.ssh/config` 中定义。
2. 从终端运行 `dogssh`。
3. 使用直观的键盘驱动 UI 管理和连接到您的服务器。

---

## ⌨️ 快捷键

| 按键  | 操作                     |
| ----- | ------------------------ |
| /     | 切换搜索栏               |
| ↑↓/jk | 导航服务器               |
| Enter | SSH 连接到所选服务器     |
| c     | 将 SSH 命令复制到剪贴板  |
| g     | Ping 所选服务器          |
| r     | 刷新后台数据             |
| a     | 添加服务器               |
| e     | 编辑服务器               |
| t     | 编辑标签                 |
| d     | 删除服务器               |
| p     | 固定/取消固定服务器      |
| s     | 切换排序字段             |
| S     | 反向排序                 |
| q     | 退出                     |

提示：列表顶部的提示栏显示了最有用的快捷方式。

## 🤝 贡献

欢迎贡献！

- 如果您发现错误或有功能请求，请 [提交 issue](https://github.com/chengzehsiao/dogssh/issues)。
- 如果您想贡献，请 fork 仓库并提交 pull request ❤️。

我们很高兴看到社区让 DogSSH 变得更好 🚀

---

## ⭐ 支持

如果您觉得 DogSSH 有用，请考虑给仓库点个 **star** ⭐️ 并加入 [stargazers](https://github.com/chengzehsiao/dogssh/stargazers)。

---

## 🙏 致谢

- 使用 [tview](https://github.com/rivo/tview) 和 [tcell](https://github.com/gdamore/tcell) 构建。
- 灵感来源于 [k9s](https://github.com/derailed/k9s) 和 [lazydocker](https://github.com/jesseduffield/lazydocker)。