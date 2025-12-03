# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

check-self is a Senzing CLI tool that reviews the environment in which it runs and returns a diagnostic report. It validates Senzing configuration paths, database connections, settings, and license information. Part of the senzing-tools suite.

## Build Commands

```bash
# Build for current platform
make build

# Build for all platforms (darwin/linux/windows, amd64/arm64)
make build-all

# Run directly
make run
# or: go run main.go

# Build Docker image
make docker-build
```

## Testing

```bash
# Setup test environment (copies SQLite test databases to /tmp/sqlite)
make setup

# Run tests
make test

# Run tests with coverage
make coverage

# Check coverage thresholds
make check-coverage
```

Running a single test:
```bash
go test -v -run TestFunctionName ./checkself/...
```

## Linting

```bash
# Run all linters (golangci-lint, govulncheck, cspell)
make lint

# Individual linters
make golangci-lint
make govulncheck
make cspell

# Auto-fix lint issues
make fix
```

Linter config: `.github/linters/.golangci.yaml`

## Architecture

- `main.go` - Entry point, calls `cmd.Execute()`
- `cmd/` - Cobra CLI command definitions using spf13/cobra and viper for configuration
  - `root.go` - Main command with all CLI flags mapped via go-cmdhelping
  - Platform-specific context: `context_linux.go`, `context_darwin.go`, `context_windows.go`
- `checkself/` - Core checker implementation
  - `checkself.go` - `BasicCheckSelf` struct and main `CheckSelf()` method that orchestrates all checks
  - Individual check files: `checkdatabaseurl.go`, `checksettings.go`, `checkdatabaseschema.go`, `checkconfigpath.go`, `checkresourcepath.go`, `checksupportpath.go`, `checklicense.go`, `checksenzingconfiguration.go`

## Key Dependencies

- `github.com/senzing-garage/sz-sdk-go` - Senzing SDK Go bindings
- `github.com/senzing-garage/go-cmdhelping` - CLI helper utilities and option definitions
- `github.com/senzing-garage/go-databasing` - Database connector abstractions
- `github.com/senzing-garage/go-helpers` - Settings parsing and error wrapping
- `github.com/spf13/cobra` and `github.com/spf13/viper` - CLI framework

## Environment Requirements

Senzing C library must be installed:
- `/opt/senzing/er/lib` - Shared objects
- `/opt/senzing/er/sdk/c` - SDK header files
- `/etc/opt/senzing` - Configuration

Set `LD_LIBRARY_PATH=/opt/senzing/er/lib` when running.

Default test database: `sqlite3://na:na@nowhere/tmp/sqlite/G2C.db`
