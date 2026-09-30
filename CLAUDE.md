# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a fork of [CosmWasm/wasmd](https://github.com/CosmWasm/wasmd), a Cosmos SDK blockchain with CosmWasm smart contract support. The fork is maintained by Burnt (Xion) with custom modifications aligned to the Xion chain. The primary module is `x/wasm`, which integrates WebAssembly smart contract execution via [wasmvm](https://github.com/CosmWasm/wasmvm).

**Go 1.24+ required.** Module path: `github.com/CosmWasm/wasmd`

## Build & Test Commands

```bash
# Build binary to ./build/wasmd
make build

# Install to $GOPATH/bin
make install

# Run all unit tests
make test

# Run a single test
go test -mod=readonly -tags='ledger test_ledger_mock' -run TestFunctionName ./x/wasm/keeper/...

# Run tests with race detection
make test-race

# Run tests with coverage
make test-cover

# Run system/integration tests (requires install first)
make test-system

# Lint (requires golangci-lint)
make lint

# Format code (installs gofumpt, misspell, gci)
make format

# Protobuf generation (requires Docker)
make proto-gen
```

## Architecture

### Core Module: `x/wasm`

The entire wasm module lives under `x/wasm/` following Cosmos SDK module conventions:

- **`keeper/`** — State machine logic. `Keeper` manages contract lifecycle (upload, instantiate, execute, migrate, sudo). `ContractKeeper` wraps Keeper with permission checks. `handler_plugin.go` and `handler_plugin_encoders.go` handle dispatching CosmWasm messages to SDK messages.
- **`types/`** — All protobuf-generated types, message definitions, params, keys, errors, and keeper interfaces (`ViewKeeper`, `ContractOpsKeeper`). Protobuf sources are in `proto/cosmwasm/wasm/v1/`.
- **`keeper/query_plugins.go`** — Custom querier plugins that bridge CosmWasm queries to SDK queries (bank, staking, distribution, IBC, etc.).
- **`keeper/msg_dispatcher.go`** — Dispatches sub-messages from contract execution, handling reply logic.
- **`keeper/ibc.go` / `ibc2.go`** — IBC integration for wasm contracts (ICS-20, ICS-27, IBC v2).
- **`migrations/`** — State migration handlers for module version upgrades (v1, v2, v3).
- **`module.go`** — AppModule registration, ABCI hooks, CLI registration.
- **`client/cli/`** — CLI commands for transactions and queries.

### App Wiring: `app/`

- **`app.go`** — Full app construction with all Cosmos SDK and IBC modules wired together. The wasm keeper is configured here with its dependencies (bank, staking, IBC, etc.).
- **`ante.go`** — Custom AnteHandler chain including wasm-specific decorators.
- **`upgrades.go` / `upgrades/`** — Chain upgrade handlers.

### Testing

- **`x/wasm/keeper/test_common.go`** — Shared test setup for keeper tests; creates a full app context with all keepers.
- **`x/wasm/keeper/wasmtesting/`** — Mock implementations for testing (mock keepers, messengers, query handlers).
- **`x/wasm/keeper/testdata/`** — Pre-compiled `.wasm` contract binaries used in tests.
- **`tests/e2e/`** — End-to-end tests.
- **`tests/system/`** — System-level integration tests (run via `make test-system`).

## Key Conventions

- **Error wrapping**: Use `errorsmod.Wrap(err, "context")` (from `cosmossdk.io/errors`), not `fmt.Errorf`. This predates Go's `errors.Is` pattern.
- **Import ordering** (enforced by `gci`): standard library, then `cosmossdk.io`, then `github.com/cosmos/cosmos-sdk`, then `github.com/CosmWasm/wasmd`, with blank line separators.
- **Determinism**: All state-machine code must be deterministic. No floating point, no maps with non-deterministic iteration, no system time.
- **Protobuf types**: Generated `.pb.go` files live alongside hand-written code in `x/wasm/types/`. Regenerate with `make proto-gen` (Docker required).
- **Build tags**: Tests use `-tags='ledger test_ledger_mock'`. The `VERSION` env var is set during test runs.
