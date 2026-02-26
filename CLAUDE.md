# Skyclerk CLI - Development Guide

## Project Overview

A Go CLI tool for interacting with the Skyclerk bookkeeping API. Built with Cobra, following the same architecture as [cloudmanic/massive](https://github.com/cloudmanic/massive).

## Project Structure

```
skyclerk-cli/
├── main.go                  # Entry point → cmd.Execute()
├── Makefile                 # Build, test, cross-compile, lint
├── Formula/skyclerk.rb      # Homebrew formula
├── .github/workflows/
│   └── release.yml          # Auto-release on push to main
├── cmd/                     # Cobra commands (UI layer)
│   ├── root.go              # Root command, global flags
│   ├── helpers.go           # newClient(), printJSON()
│   ├── login.go             # login / logout
│   ├── config.go            # config show / init
│   ├── accounts.go          # accounts list / show / use
│   ├── ledger.go            # ledger CRUD + summary
│   ├── categories.go        # categories CRUD
│   ├── labels.go            # labels CRUD
│   ├── contacts.go          # contacts CRUD
│   ├── files.go             # files upload
│   ├── activities.go        # activities list
│   ├── me.go                # user profile
│   ├── users.go             # user/invite management
│   ├── reports.go           # P&L and financial reports
│   └── version.go           # version display
├── internal/
│   ├── api/                 # API client layer
│   │   ├── client.go        # HTTP client (get/post/put/delete/upload)
│   │   ├── types.go         # All request/response structs
│   │   ├── auth.go          # Login/Logout/GetAuthUser
│   │   ├── accounts.go      # Account endpoints
│   │   ├── ledger.go        # Ledger endpoints
│   │   ├── categories.go    # Category endpoints
│   │   ├── labels.go        # Label endpoints
│   │   ├── contacts.go      # Contact endpoints
│   │   ├── files.go         # File upload endpoint
│   │   ├── activities.go    # Activity endpoint
│   │   ├── reports.go       # Report endpoints
│   │   ├── me.go            # User profile endpoints
│   │   ├── users.go         # User management endpoints
│   │   └── api_test.go      # All API tests (mock HTTP servers)
│   └── config/
│       ├── config.go        # Config read/write (~/.config/skyclerk/)
│       └── config_test.go   # Config tests
```

## Key Architecture Decisions

- **One test file per source file** in each package. All API tests are in `api_test.go`, all config tests in `config_test.go`.
- **API layer returns raw types**, command layer handles display (table vs JSON).
- **Config stored at** `~/.config/skyclerk/config.json` with 0600 permissions.
- **Bearer token auth** via `Authorization` header on all authenticated requests.
- **Account ID in URL path**: All Skyclerk API routes follow `/api/v3/:account/{endpoint}`.

## Skyclerk API Notes

- **Base URL**: `https://app.skyclerk.com`
- **Auth**: `POST /oauth/token` with `username`, `password`, `grant_type=password`, and a registered `client_id`.
- **Category types**: The API returns `"income"` and `"expense"` as strings (not numeric codes).
- **Ledger summary route**: `/ledger-summary` (not `/ledger/summary` which conflicts with `/ledger/:id`).
- **Ledger P&L route**: `/ledger-pl-summary` (not `/ledger/pl`).
- **Ledger list**: Does NOT support the `order` query parameter. Using it causes an empty response.
- **Ledger summary response**: Returns `years`, `categories`, and `labels` with counts (not income/expense/profit totals).
- **Dates**: Full ISO timestamps like `"2026-02-01T08:00:00Z"`, not just `YYYY-MM-DD`.

## Build & Release

- Version injected via ldflags: `-X github.com/cloudmanic/skyclerk-cli/cmd.Version=...`
- GitHub Actions auto-releases on push to main with auto-incrementing patch versions.
- Homebrew formula at `Formula/skyclerk.rb` downloads pre-built binaries from GitHub releases.
- Install: `brew tap cloudmanic/skyclerk-cli https://github.com/cloudmanic/skyclerk-cli && brew install skyclerk`

## Commands

```
make build       # Build for current platform
make test        # Run all tests (verbose)
make test-short  # Run tests (quiet)
make coverage    # Generate coverage report
make lint        # Run fmt + vet
make cross-build # Build for all platforms
make clean       # Remove artifacts
```
