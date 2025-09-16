# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

LazySSH is a terminal-based, interactive SSH manager inspired by tools like lazydocker and k9s. It provides a clean, keyboard-driven UI for managing SSH connections defined in `~/.ssh/config`.

**Tech Stack:**
- Go 1.24.6
- UI: tview + tcell for terminal interface
- CLI: cobra for command parsing
- Logging: zap for structured logging
- SSH config parsing: Custom fork of kevinburke/ssh_config

## Development Commands

### Build and Run
```bash
make build          # Build binary to ./bin/lazyssh
make run            # Run directly from source
make install        # Build and install to $GOBIN
make build-all      # Cross-compile for all platforms
```

### Code Quality and Testing
```bash
make quality        # Run all quality checks (fmt, vet, lint)
make fmt            # Format code using gofumpt
make lint           # Run golangci-lint
make check          # Run staticcheck analyzer

make test           # Run unit tests with race detection
make test-verbose   # Run tests with verbose output
make coverage       # Generate coverage report (coverage.html)
make benchmark      # Run benchmarks
```

### Dependency Management
```bash
make deps           # Download and verify dependencies
make tidy           # Tidy dependencies
make tools          # Install development tools locally
```

### Development Tools
```bash
make run-race       # Run with race detector
make security       # Run security checks
make clean          # Clean build artifacts and caches
```

## Architecture

The project follows **Hexagonal Architecture** with clear separation of concerns:

### Core Domain (`internal/core/`)
- **`domain/server.go`**: Server entity with fields like Alias, Host, User, Port, Tags, SSH metadata
- **`ports/`**: Defines interfaces for repositories and services
- **`services/server_service.go`**: Business logic including validation, SSH operations, server management

### Adapters (`internal/adapters/`)
- **`ui/`**: Terminal UI using tview - handles TUI components, key bindings, and user interactions
- **`data/ssh_config_file/`**: SSH config file operations - reads/writes `~/.ssh/config`, manages metadata and passwords

### Key Components
- **TUI (`internal/adapters/ui/tui.go`)**: Main application orchestrator with layout management
- **ServerList, SearchBar, ServerDetails**: UI components for different views
- **SSH Config Repository**: Handles safe, atomic config file operations with backups
- **Metadata Manager**: Stores additional server data (pins, SSH counts) in `~/.lazyssh/metadata.json`
- **Password Manager**: Encrypted password storage in `~/.lazyssh/passwords.json`

## Development Guidelines

### File Structure Patterns
- All Go files include Apache 2.0 license headers
- Use structured logging with zap throughout
- Repository pattern for data access
- Dependency injection via constructors

### SSH Config Safety
- Non-destructive edits preserving comments and formatting
- Atomic writes via temporary files
- Automatic backups: `config.original.backup` (one-time) and timestamped rolling backups
- System SSH compatibility - uses native `ssh` command for connections

### Code Quality Standards
- golangci-lint with comprehensive rule set (.golangci.yml)
- Static analysis with staticcheck
- Race detection enabled in tests
- Code formatting with gofumpt

### Testing Strategy
- Unit tests with race detection: `make test`
- Coverage reporting: `make coverage`
- Benchmarks available: `make benchmark`
- Currently no test files present - tests would go alongside source files as `*_test.go`

## Key Interfaces

**ServerService**: Core business operations (List, Add, Update, Delete, SSH, Ping)
**ServerRepository**: Data persistence layer for SSH config and metadata
**FileSystem**: Abstraction for file operations (enables testing with mocks)

## Security Notes

- No storage/transmission of private keys or sensitive SSH credentials
- Password storage is encrypted and optional
- All SSH connections use system's native SSH binary
- Config file permissions preserved during operations