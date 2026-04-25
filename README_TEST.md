# Testing Guide

This file contains only test-related commands for the project.

## Run all tests

```powershell
go test ./...
```

## Run tests for one package

```powershell
go test ./cmd/server
go test ./internal/auth
go test ./internal/database
go test ./internal/handlers
go test ./internal/models
go test ./internal/repository
go test ./internal/routes
go test ./tests/e2e
```

## Run E2E test against a running server

The `tests/e2e` package contains a health endpoint test.

- Without `E2E_BASE_URL`, the test is skipped.
- With `E2E_BASE_URL`, it performs a real HTTP request to `/`.

Example:

```powershell
$env:E2E_BASE_URL = "http://localhost:8080"
go test ./tests/e2e -v
```

## Run tests with coverage summary

```powershell
go test ./... -coverprofile=coverage.out
go tool cover "-func=coverage.out"
```

## Generate HTML coverage report

```powershell
go test ./... -coverprofile=coverage.out
go tool cover "-html=coverage.out" -o coverage.html
```

Open `coverage.html` in your browser to inspect line-by-line coverage.

## Optional shortcuts (Makefile)

If you use `make` (Linux/macOS, or Windows with Make installed), these targets are available:

```powershell
make test
make test-cover
make test-cover-html
```

## PowerShell troubleshooting

If commands fail because of typos or shell parsing, copy/paste exactly:

```powershell
go test ./... "-coverprofile=.\coverage.out"
Test-Path .\coverage.out
go tool cover "-func=.\coverage.out"
go tool cover "-html=.\coverage.out" -o .\coverage.html
```
