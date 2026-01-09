# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`check-self` is a Go CLI diagnostic tool in the Senzing Garage (experimental) project suite. It verifies the environment where Senzing tools run, performing comprehensive checks on configurations, databases, licenses, and system state.

**Status**: Work-in-progress with semantic versions `0.n.x` - not production-ready.

## Build Commands

```bash
make build                 # Build for current OS/architecture
make build-all             # Build for all 6 platforms (darwin/linux/windows × amd64/arm64)
make clean                 # Remove build artifacts and caches
```

Binaries output to `target/<os>-<arch>/check-self`.

## Test Commands

```bash
make clean setup test      # Run tests (setup creates test database)
make check-coverage        # Run tests and verify coverage thresholds
make coverage              # Generate coverage report and open in browser

# Run a single test
go test -v -run TestFunctionName ./checkself/

# Run tests with trace logging
SENZING_LOG_LEVEL=TRACE go test ./...
```

Coverage targets: 80%+ overall, minimum 30% per file, 70% for `cmd/` package.

## Lint Commands

```bash
make lint                  # Run all linters (golangci-lint, govulncheck, cspell)
make fix                   # Auto-fix ~20 linter issues (gofumpt, wsl, nakedret, etc.)
make golangci-lint         # Run golangci-lint only
make govulncheck           # Vulnerability scanning
make cspell                # Spell check (includes dotfiles)
```

Linter config: `.github/linters/.golangci.yaml` - enables 100+ linters with strict settings.

## Development Setup

```bash
make dependencies-for-development  # Install dev tools (golangci-lint, gotestfmt, govulncheck, gofumpt, cspell)
make dependencies                  # Update and tidy Go dependencies
```

**Prerequisites**: Senzing C library installed at `/opt/senzing/er/lib` with SDK headers at `/opt/senzing/er/sdk/c`.

## Architecture

### Package Structure

- `main.go` - Entry point, calls `cmd.Execute()`
- `cmd/` - Cobra CLI setup with Viper configuration management
  - `root.go` - Root command definition, context variables, PreRun hooks
  - `context_<os>.go` - Platform-specific context handling
- `checkself/` - Core health-check logic
  - `checkself.go` - `BasicCheckSelf` struct and `CheckSelf()` orchestrator
  - `check*.go` - Individual health check implementations

### Core Pattern: Chain of Responsibility

The `CheckSelf()` method in `checkself/checkself.go` orchestrates checks:

1. `getTestFunctions()` returns ordered list of check functions
2. Each check follows signature: `(ctx, reportChecks, reportInfo, reportErrors) → (reportChecks, reportInfo, reportErrors, error)`
3. Checks execute sequentially; first error stops the chain
4. Reports accumulate: Information → Checks Performed → Errors → Result

### Check Functions (in execution order)

1. `Prolog` - Output header/separator
2. `ListEnvironmentVariables` - Introspect environment
3. `ListStructVariables` - Dump configuration values
4. `CheckConfigPath` - Validate config directory
5. `CheckResourcePath` - Validate resource directory
6. `CheckSupportPath` - Validate support directory
7. `CheckDatabaseURL` - Validate database connection string
8. `CheckSettings` - Parse and validate Senzing settings JSON
9. `CheckDatabaseSchema` - Verify database schema integrity
10. `CheckSenzingConfiguration` - Validate Senzing config (disabled)
11. `CheckLicense` - Validate license (disabled)

### Key Types

```go
// Main configuration struct - populated via Viper from env/flags/files
type BasicCheckSelf struct {
    ConfigPath, DatabaseURL, ResourcePath, SupportPath, Settings string
    // ... plus other config fields
}

// Factory pattern for Senzing SDK access
szfactorycreator.CreateCoreAbstractFactory() → senzing.SzAbstractFactory
szAbstractFactory.CreateConfigManager() → senzing.SzConfigManager
szAbstractFactory.CreateProduct() → senzing.SzProduct
```

### Database Support

Supports: SQLite3, MySQL, PostgreSQL, Oracle (godror), MSSQL, DB2. Database URL can come from `DatabaseURL` field or extracted from `Settings` JSON.

## Docker

```bash
make docker-build          # Build Docker image
make docker-run            # Run container
make docker-test           # Integration tests via docker-compose
make docker-build-package  # Create RPM/DEB packages
```

## Code Style Notes

- Line length limit: 120 characters
- Uses `gofumpt` for formatting (stricter than `gofmt`)
- WSL linter enforces whitespace conventions
- `exhaustruct` excludes `cobra.Command`, `BasicCheckSelf`, `ProductLicenseResponse`
- Function max complexity: 14 (cyclop), max lines: 65 (funlen)
- Use `wraperror.Errorf()` for error wrapping
