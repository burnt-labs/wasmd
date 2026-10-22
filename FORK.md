# The Burnt wasmd fork

This repository is Burnt's fork of [CosmWasm/wasmd](https://github.com/CosmWasm/wasmd),
and XION builds against it through a `replace` directive. This file records
how the fork tracks upstream, the patches we carry and why, and which commit
every shipped xion tag corresponds to. `AGENTS.md` has the rules and
`scripts/xion-sync` does the work.

## Model: one branch per upstream release line

- **Line branches.** `xion/v0.61.x`, `xion/v0.70.x`, and so on. Each one
  starts at an upstream tag and carries our patch stack as discrete commits,
  one per concern, each saying why it exists.
- **Same line, new upstream patch release** (v0.70.5 → v0.70.6): merge the
  upstream tag into the line branch with a signed merge commit, fix the
  stack if it needs it, tag `v0.70.6-xion.1`. The upstream tag stays an
  ancestor, so the next merge only brings in what is new.
- **New upstream line** (v0.70 → v0.71): start `xion/v0.71.x` at the upstream
  tag and cherry-pick the stack from the previous line, then check it with
  `git range-diff`. Upstream's release branches diverge from each other, so a
  cross-line merge conflicts in many files; a cherry-pick does not have to.
- **`xion/main`** is the latest line's head and only moves fast-forward,
  except when a new line starts (see `AGENTS.md`). Older lines stay as their
  own branches, so a fix for an older XION release needs no history surgery.
- **Tags** `vX.Y.Z-xion.N` are cut only on a line branch's head and never
  moved or deleted. `scripts/xion-sync --tag` picks the next free `N`,
  including numbers the Go proxy has already seen.

## Patch stack on `xion/v0.70.x` (XION v32)

`git log --oneline --no-merges v0.70.5..origin/xion/v0.70.x`, oldest first.
None of these are upstream PRs yet; upstreaming one shrinks the stack.

| Commit | Concern | Why it exists | Upstream status |
| -- | -- | -- | -- |
| `a386ce5e` | ci/docs: fork workflows, Dockerfile, CODEOWNERS, CLAUDE.md, SECURITY.md | The fork runs its own CI and review routing and documents its security contact. | fork-only |
| `5b382e4c` | feat(genesis): contract genesis export/import | XION's genesis tooling imports codes, contracts and sequences from outside the keeper, keeping each history entry's position. | not proposed |
| `e1732403` | feat(wasm): instantiate with module-owned address hash | Lets a trusted module (XION's abstract accounts) choose a contract address through a separate capability. | not proposed |
| `9cc71786` | fix(wasm): reject contract results with neither ok nor error set | Comments and tests around the nil-ok guards; upstream carries the guards themselves since v0.61.15 / v0.70.x. | guard upstream, tests fork-only |
| `d931befa` | build: adapt the xion mods to the v0.70 line | Go 1.25.14 and grpc 1.82.1 security pins, the ibc-go/v11 import in the nil-ok tests, golangci-lint v2.11.4. Replaces the v0.61 dependency-alignment commit. | fork-only |
| `82f9fe30` | feat: port Xion wasmd fork to Cosmos SDK 0.55 | XION v32 moves to Cosmos SDK v0.55.0, CometBFT v0.40.0 and Go 1.26; removes legacy param subspaces and the protocol pool. | upstream has no SDK 0.55 release yet |
| `c0763843` | ci: push Docker images only from CosmWasm/wasmd | Only upstream holds the Docker Hub credentials, so the login and push steps are gated on the repository and every fork tag run no longer fails. | fork-only |
| `9284ec52` | build: build the Docker image on Go 1.26 | go.mod requires Go 1.26.5; the golang:1.25 builder could not build it. | fork-only |
| `f74a1cbb` | ci: run the determinism simulation outside the -race coverage loop | On SDK 0.55 the determinism simulation runs its full 500 blocks and timed out under -race; it now runs without -race in its own job. | fork-only |
| `618482dd` | fix(testnet): give each single-host node its own pprof port | CometBFT v0.40 aborts start-up when pprof cannot bind, so single-host testnets need a port per node. | could be proposed upstream |
| `1f4d9ae0` | test(system): give each node its own pprof port | Same fix for the system tests, including configs written by legacy binaries in the upgrade test. | could be proposed upstream |
| `88c5db29` | docs: give the fork release its own changelog section | Fork releases are listed separately from upstream's. | fork-only |
| `25284f70` | docs(security): pause the bug bounty program (#26) | Current security policy. | fork-only |
| `6ce91b14` | docs(changelog): v0.70.5-xion.4 section | Names the release XION v32 pins. | fork-only |

The same concerns on the v0.61 line, compared with
`git range-diff v0.61.15..v0.61.15-xion.2 v0.70.5..xion/v0.70.x`: ci/docs and
genesis changed (adapted to v0.70), module-owned address and the nil-ok
guard are identical, the v0.61 dependency-alignment commit is replaced by
`d931befa`, and the rest is new for v0.70 and SDK 0.55.

## Shipped xion tags and their line commits

Every wasmd version a XION release pins has a commit on its line whose tree
is identical to the tag (`git diff --stat <tag> <line commit>` is empty).
The tags themselves stay where they are, because XION releases and the Go
sum database pin them.

| xion tag | Tag commit | XION releases | Line commit | Upstream base |
| -- | -- | -- | -- | -- |
| v0.61.6-xion.1 | `5ab90e25` | v22.0.0 | `ee59da00` on `xion/v0.61.x` | v0.61.6 |
| v0.61.6-xion.3 | `940caaf9` | v23.0.0 | `b9096098` on `xion/v0.61.x` | v0.61.6 |
| v0.61.6-xion.4 | `d6fa4d67` | v21.1.0, v24.0.0-v26.1.0 | `b9096098` on `xion/v0.61.x` (same tree as xion.3) | v0.61.6 |
| v0.61.8-xion.1 | `ece56521` | v28.0.1 | `47aeb9cf` on `xion/v0.61.x` | v0.61.8 |
| v0.61.8-xion.2 | `a3464e53` | v28.1.0 | `e0c84c3c` on `xion/v0.61.x` | v0.61.8 |
| v0.61.10-xion.1 | `747fa521` | v29.0.0-v30.0.0 | `18a7a689` on `xion/v0.61.x` | v0.61.10 |
| v0.61.15-rc-xion.1 | `a186f1a7` | v31.0.0 | `5d1d15bb` on `xion/v0.61.x` | v0.61.14 (the tag predates upstream v0.61.15) |
| v0.61.15-xion.2 | `74aa8875` | v31.0.1, v31.0.2 | `401f73e1` on `xion/v0.61.x` (head) | v0.61.15 |
| v0.70.5-xion.4 | `80464e7b` | v32 | `6ce91b14` on `xion/v0.70.x` (head) | v0.70.5 |

Notes:

- `v0.61.15-rc-xion.1` is no longer a tag on this repository; the Go proxy
  still serves it at `a186f1a7`.
- XION v6.x and earlier pin a pseudo-version,
  `v0.45.1-0.20240312170141-3dfd94af17a8`: commit `3dfd94af` on the
  `fix/no-state-migration` branch, from before xion tags existed. It is not
  on a line branch.
- `v0.70.5-xion.2` (`25284f70`) is on `xion/v0.70.x`; its tree differs from
  v0.70.5-xion.4 only in three CHANGELOG lines.
- Interim tags that no XION release pins and that are not on a line branch:
  `v0.61.15-xion.3`, `v0.61.16-xion.1`, `v0.70.0-xion.1` to
  `v0.70.3-xion.1`, `v0.70.4-xion.3`, `v0.70.5-xion.3`, plus the older
  `v0.61.5-*`, `v0.61.6-xion.2`, `v0.61.8-xion.3`, `v0.61.9-*`,
  `v0.61.10-xion.2` and `v0.61.14-*`. They stay until twiceburnt decides
  otherwise.
- Names the Go proxy has seen can never be reused, even after a tag is
  deleted here: `v0.70.4-xion.1`, `v0.70.4-xion.2`, `v0.70.5-xion.1` and
  `v0.61.15-xion.1` are such names.
