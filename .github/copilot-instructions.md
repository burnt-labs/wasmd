# Copilot Instructions for Wasmd

This document provides guidance for AI coding agents working with the wasmd repository. Wasmd is the first implementation of a Cosmos zone with WebAssembly (Wasm) smart contracts enabled.

## Project Overview

- **Language**: Go 1.25.3+ (primary), Protocol Buffers, Shell scripts
- **Framework**: Cosmos SDK v0.53.6, CosmWasm v3.0.3
- **Key Dependencies**: CometBFT v0.38.21, IBC-Go v10.5.0
- **Purpose**: Blockchain application enabling CosmWasm smart contracts in the Cosmos ecosystem

## Repository Structure

```
wasmd/
├── app/                  # Cosmos SDK application setup and configuration
├── cmd/wasmd/           # Main CLI binary entry point
├── x/wasm/              # Core wasm module (main focus)
│   ├── keeper/          # Business logic and state management
│   ├── types/           # Protobuf-generated types and message definitions
│   ├── client/cli/      # CLI command implementations
│   ├── simulation/      # Simulation testing for governance
│   ├── migrations/      # Version migration handlers
│   ├── exported/        # Public interface definitions
│   └── ioutils/         # File I/O utilities
├── proto/               # Protocol buffer definitions (.proto files)
├── tests/               # Test suites
│   ├── e2e/            # End-to-end tests (full chain scenarios)
│   ├── integration/    # Integration tests (module interactions)
│   ├── system/         # System-level tests (CLI, node behavior)
│   └── wasmibctesting/ # IBC testing utilities
├── benchmarks/          # Performance benchmarks
├── scripts/             # Build and utility scripts
├── docker/              # Docker setup and scripts
└── docs/                # Documentation files
```

## Development Workflow

### Essential Commands

**Build and Install:**
```bash
make install              # Install wasmd binary to $GOPATH/bin
make build               # Build wasmd binary to build/wasmd
```

**Testing:**
```bash
make test                # Run unit tests
make test-race           # Run tests with race detector
make test-cover          # Generate coverage reports
make test-system         # Run system tests
make test-all            # Run all test suites
```

**Code Quality:**
```bash
make format              # Format code (MUST run before every commit)
make lint                # Run golangci-lint (v2.1.6)
```

**Protobuf:**
```bash
make proto-all           # Generate, format, and lint proto files
make proto-gen           # Generate code from proto files (Docker-based)
make proto-format        # Format proto files with buf
make proto-lint          # Lint proto files with buf
```

### Pre-Commit Requirements

**CRITICAL**: Always run `make format` before committing. This runs:
- `gofumpt` v0.4.0 - Go code formatting
- `gci` v0.11.2 - Import sorting with custom order
- `misspell` v0.3.4 - Spelling checks

The CI will fail if code is not properly formatted.

### CI Pipeline

The repository uses GitHub Actions (`.github/workflows/checks.yml`):
1. **setup-dependencies** - Cache Go modules and build
2. **tidy-go** - Verify `go mod tidy` was run
3. **lint** - Run golangci-lint
4. **test-cover** - Run tests with coverage (4 parallel shards)
5. **test-system** - Run system tests
6. **benchmark** - Run performance benchmarks
7. **simulations** - Run deterministic simulations
8. **docker-image** - Build and verify Docker image

Tests run with:
- Race detector enabled (`GORACE=halt_on_error=1`)
- 8-minute timeout per package
- Coverage tracking uploaded to Codecov
- Build tags: `ledger test_ledger_mock`

## Coding Guidelines

### Error Handling

**ALWAYS wrap errors with context** - never return bare `err`:

```go
// ❌ BAD - No context
if err := k.bank.TransferCoins(ctx, caller, contractAddress, coins); err != nil {
    return nil, err
}

// ✅ GOOD - Wrapped with context
if err := k.bank.TransferCoins(ctx, caller, contractAddress, coins); err != nil {
    return nil, errorsmod.Wrap(err, "lock contract coins")
}
```

Use `errorsmod.Register()` for custom errors (see `x/wasm/types/errors.go`).

### Import Order

Imports MUST follow this order (enforced by `gci`):
1. Standard library packages
2. External packages
3. `cosmossdk.io` packages
4. `github.com/cosmos/cosmos-sdk` packages
5. `github.com/CosmWasm/wasmd` packages

Example:
```go
import (
    // Standard
    "context"
    "fmt"

    // External
    "github.com/pkg/errors"

    // Cosmos SDK
    "cosmossdk.io/log"

    "github.com/cosmos/cosmos-sdk/types"

    // Wasmd
    "github.com/CosmWasm/wasmd/x/wasm/keeper"
)
```

### Code Structure Best Practices

1. **Package Organization**: Limited responsibility per package, different concerns in different packages
2. **Interface Design**: Depend on abstractions, not concretions
3. **Minimal Public API**: Easier to expose later than to hide
4. **Avoid Global State**: No global variables or configurators
5. **Thread Safety**: Clearly document non-thread-safe code
6. **Determinism**: All code must be deterministic (critical for blockchain)

### Security Considerations

Be vigilant about:
- **Gas usage** - All operations must consume appropriate gas
- **Transaction verification** - Signature validation
- **Malleability attacks** - Ensure deterministic serialization
- **Code determinism** - No randomness, timestamps must be from block header
- **Input validation** - Validate at system boundaries only

### Testing Patterns

**Test Co-location**: Place `*_test.go` files next to implementation files.

**Common Test Utilities**:
- `app/test_helpers.go` - SetupOptions struct for test chain setup
- `tests/integration/common_test.go` - Shared integration test helpers
- `x/wasm/keeper/testdata/` - Embedded test contract data (e.g., `hackatom.wasm`)

**Testing Framework**: Use `github.com/stretchr/testify` for assertions.

**Test Categories**:
1. **Unit tests** - Co-located with source (`*_test.go`)
2. **Integration tests** - `tests/integration/` (module interactions)
3. **E2E tests** - `tests/e2e/` (full chain scenarios)
4. **System tests** - `tests/system/` (CLI and node behavior)
5. **Simulation tests** - `app/sim_test.go` and `x/wasm/simulation/`

## Protobuf Development

### Workflow

1. Edit `.proto` files in `proto/cosmwasm/wasm/v1/`
2. Run `make proto-all` to generate, format, and lint
3. Generated files appear in `x/wasm/types/` with `.pb.go` suffix

### Proto Files Organization

- `query.proto` - Query service definitions
- `tx.proto` - Transaction message definitions
- `types.proto` - Common types and models
- `genesis.proto` - Genesis state
- `authz.proto` - Authorization definitions
- `ibc.proto` - IBC-related messages
- `proposal_legacy.proto` - Legacy governance proposals

### Important Notes

- Protobuf generation is **Docker-based** (image: `ghcr.io/cosmos/proto-builder:0.14.0`)
- Breaking change detection runs against `main` branch
- Never manually edit `.pb.go` files - they're generated
- Use `buf` tool for linting (config in `proto/buf.yaml` and `proto/buf.lock`)

## Module Development (x/wasm)

### Key Components

**Keeper** (`x/wasm/keeper/keeper.go`):
- Central business logic (~61KB)
- Contract execution and state management
- Gas metering and limits
- Query and message handling

**Types** (`x/wasm/types/`):
- Message definitions and validation
- Codec registration (amino and protobuf)
- Error definitions
- Event helpers

**CLI** (`x/wasm/client/cli/`):
- `tx.go` - Transaction commands
- `query.go` - Query commands
- `gov_tx.go` - Governance proposal commands

### Event Emission

Emit events for observability (see `x/wasm/types/events.go`):
```go
emitEvent(ctx, sdk.NewEvent(
    types.EventTypeExecute,
    sdk.NewAttribute(types.AttributeKeyContractAddr, contractAddr.String()),
    sdk.NewAttribute(types.AttributeKeySender, sender.String()),
))
```

Events are tagged with:
- `module: wasm`
- `_contract_address` - Contract that emitted event
- `action` - Operation type (store-code, instantiate, execute)

### Configuration

Module configuration (`config/app.toml`):
```toml
[wasm]
query_gas_limit = 300000        # Max gas for smart queries
memory_cache_size = 300         # Wasm module cache size in MiB
```

CLI flags override config:
```bash
--wasm.memory_cache_size uint32
--wasm.query_gas_limit uint
```

## Git and PR Workflow

### Branching

- **Main branch**: `main` (trunk-based development)
- **Branch naming**: `{issue#}-branch-name` for core developers
- **Never force push** to `main` (except reverting broken commits)

### Pull Request Process

1. **Start with Draft PR** - Get early validation
2. **Ensure clean state**:
   - Run `make format` before committing
   - Run `make lint test` before marking ready for review
   - Merge latest `main`: `git merge origin/main`
3. **Add changelog entry** - Update `CHANGELOG.md` in `Unreleased` section (top of changes)
4. **Single issue per PR** - One PR addresses one issue
5. **PR title**: Start with uppercase letter
6. **Ready for Review** - Change from Draft when complete

### Commit Process

- GitHub squashes commits and rebases on merge to `main`
- Commit messages should be descriptive
- Reference issue numbers in commits

## Common Patterns and Conventions

### Naming Conventions

- Package names match directory names
- Interfaces end with appropriate suffix (`Keeper`, `Handler`)
- Test files: `<module>_test.go`
- Private functions: `camelCase`
- Public functions: `PascalCase`

### Cosmos SDK Patterns

- Use SDK context: `sdk.Context` for all state operations
- Logger access: `ctx.Logger()`
- Store access through keeper methods
- Module registration via `AppModule` interface

### Contract Interaction

- Contract instantiation: `MsgInstantiateContract`
- Contract execution: `MsgExecuteContract`
- Contract migration: `MsgMigrateContract`
- Code storage: `MsgStoreCode`

### Gas Metering

All contract operations consume gas:
- Store operations: Gas per byte
- Contract execution: Gas per instruction
- Queries: Limited by `query_gas_limit`

## Build Configuration

### Build Tags

- `netgo` - Always enabled
- `ledger` - Hardware wallet support (requires GCC)
- `gcc` - CLevelDB support
- `muslc` - Static compilation for Alpine Linux

### Linker Flags

Set via `-ldflags`:
- Version info: `version.Version`, `version.Commit`
- App name: `version.AppName=wasmd`
- Bech32 prefix: `app.Bech32Prefix=wasm` (customizable)
- Build tags embedded in binary

### Cross-Compilation

- Linux (glibc): Default target, CentOS 7 compatible (glibc 2.12+)
- Linux (muslc): Alpine Linux, requires `-tags muslc`
- Windows client: `make build-windows-client`
- macOS: Supported, M1 experimental

## Dependencies and Compatibility

### Critical Version Constraints

- **wasmvm**: Minor version bumps are consensus-breaking
- **Cosmos SDK**: Follow wasmd's specified version closely
- **Go version**: 1.25.3+ required

### Adding Dependencies

- **Minimize third-party deps** - Only add when necessary
- **Check license compatibility** - Apache 2.0 or compatible
- **Security review** - Especially for crypto or consensus code
- Run `go mod tidy` after adding deps

## Linter Configuration

### Enabled Linters (.golangci.yml)

- `copyloopvar` - Loop variable copying
- `dogsled` - Blank identifiers (max 6)
- `errcheck` - Unchecked errors
- `goconst` - Repeated strings
- `gocritic` - Go critic checks
- `gosec` - Security checks (extensive rules)
- `govet` - Go vet checks
- `ineffassign` - Ineffectual assignments
- `misspell` - Spelling mistakes
- `nakedret` - Naked returns
- `nolintlint` - Nolint directives
- `revive` - Golint replacement
- `staticcheck` - Static analysis
- `unconvert` - Unnecessary conversions
- `unused` - Unused code

### Excluded Paths

- `*.pb.go` - Generated protobuf files
- `*.pb.gw.go` - Generated gateway files
- `testutil/testdata` - Test fixtures
- `third_party/` - External code

## Documentation

### Key Documentation Files

- `README.md` - Project overview and quick start
- `CONTRIBUTING.md` - Contribution guidelines
- `CODING_GUIDELINES.md` - Code standards
- `CHANGELOG.md` - Version history
- `UPGRADING.md` - Upgrade instructions
- `SECURITY.md` - Security policy
- `x/wasm/README.md` - Wasm module documentation
- `x/wasm/Governance.md` - Governance features
- `x/wasm/IBC.md` - IBC integration

### Documentation Style

- Use clear, concise language
- Include code examples where helpful
- Keep docs up to date with code changes
- Document breaking changes prominently

## Troubleshooting Common Issues

### Build Issues

**Problem**: `wasmvm` version mismatch
**Solution**: Check `go list -m github.com/CosmWasm/wasmvm/v3` and verify against `go.mod`

**Problem**: Protobuf generation fails
**Solution**: Ensure Docker is running and accessible

**Problem**: Ledger support build fails
**Solution**: Install GCC or set `LEDGER_ENABLED=false`

### Test Issues

**Problem**: Tests fail with "X server" errors on headless Linux
**Solution**: See issue #31 for DBus/keyring workarounds

**Problem**: Race detector failures
**Solution**: These indicate real concurrency issues - must fix, not ignore

**Problem**: System tests fail
**Solution**: Check `tests/system/testnet` artifacts (auto-uploaded on failure)

### Lint Issues

**Problem**: Import order incorrect
**Solution**: Run `make format` - `gci` will fix automatically

**Problem**: Golangci-lint timeout
**Solution**: Increase timeout: `make lint ARGS="--timeout=10m0s"`

**Problem**: Generated files causing lint errors
**Solution**: Regenerate with `make proto-gen` and run `make format`

## Performance Considerations

### Optimization Guidelines

- **Avoid unnecessary allocations** - Reuse objects when safe
- **Minimize store reads/writes** - Cache expensive lookups
- **Gas efficiency** - Every operation costs gas
- **Batch operations** - Group state changes when possible

### Benchmarking

Run benchmarks before/after performance changes:
```bash
cd x/wasm/keeper && go test -bench .
cd benchmarks && go test -bench .
```

Benchmarks compare:
- Contract execution performance
- State operation costs
- Gas calculation accuracy

## Reference: Cosmos SDK Integration

### Module Interface

Implement `AppModule` interface:
- `RegisterServices` - Register gRPC/query services
- `RegisterStoreDecoder` - Simulation support
- `RegisterMigrations` - Version migrations
- `InitGenesis` / `ExportGenesis` - Genesis handling

### Cosmos SDK Context Usage

```go
// Logging
logger := ctx.Logger().With("module", "x/wasm")

// Store access
store := ctx.KVStore(k.storeKey)

// Block info
height := ctx.BlockHeight()
time := ctx.BlockTime()

// Gas metering
ctx.GasMeter().ConsumeGas(amount, "description")
```

### Events and Telemetry

Emit events for indexing:
```go
ctx.EventManager().EmitEvent(
    sdk.NewEvent(
        eventType,
        sdk.NewAttribute(key, value),
    ),
)
```

## Quick Reference Commands

```bash
# Development cycle
make install                    # Build and install
make test                       # Run tests
make lint                       # Check code quality
make format                     # Format code (pre-commit)

# Full verification (run before PR)
make lint test test-race        # Main CI checks
git merge origin/main           # Ensure up to date

# Protobuf workflow
make proto-all                  # Full proto pipeline

# Docker development
docker build -t cosmwasm/wasmd:latest .
docker run --rm -it cosmwasm/wasmd:latest /usr/bin/wasmd

# Version checking
wasmd version                   # Binary version
wasmd query wasm libwasmvm-version  # Runtime wasmvm version
go list -m github.com/CosmWasm/wasmvm/v3  # Go module version
```

## Additional Resources

- **CosmWasm Documentation**: https://docs.cosmwasm.com
- **Cosmos SDK Documentation**: https://docs.cosmos.network
- **Protobuf Buf Registry**: https://buf.build/cosmwasm/wasmd
- **Issue Tracker**: https://github.com/CosmWasm/wasmd/issues
- **Good First Issues**: Label `good first issue`

## Notes for AI Agents

1. **Always read files before modifying** - Understand context first
2. **Follow existing patterns** - Consistency is critical in blockchain code
3. **Security first** - This is consensus-critical code
4. **Test thoroughly** - Include unit, integration, and system tests
5. **Document breaking changes** - Update CHANGELOG.md
6. **Never skip `make format`** - CI will fail
7. **Wrap all errors** - Context is essential for debugging
8. **Be deterministic** - No randomness, no timestamps except from block header
9. **Mind the gas** - All operations must consume appropriate gas
10. **Check CI** - Ensure all checks pass before marking PR ready

When in doubt, reference the existing codebase for patterns and conventions. The wasmd codebase is well-established and follows consistent patterns throughout.
