# AGENTS.md

This repository is Burnt's fork of [CosmWasm/wasmd](https://github.com/CosmWasm/wasmd).
`CLAUDE.md` covers building, testing and the code layout. `FORK.md` lists the
patches we carry, why, and which line commit each shipped xion tag matches.
This file sets the rules for our branches and tags.

These rules cover only our own branches and tags (`xion/*`,
`vX.Y.Z-xion.N`). Upstream tags (`vX.Y.Z`) and the `main` branch belong to
upstream: don't move, re-tag or push to them.

## Branches

- One line branch per upstream release line: `xion/v0.61.x`, `xion/v0.70.x`,
  and so on. A line starts at an upstream tag and carries our patch stack as
  discrete commits, one per concern. Each commit message says why the patch
  exists. Never squash the stack into one commit, and never rebase a line.
- `xion/main` (the default branch) is the head of the latest line. It moves
  only fast-forward. The one exception is starting a new line, which is not
  a descendant of the old one: push `backup/xion-main-<short sha>` first,
  then `git push --force-with-lease=xion/main:<old sha>`, and only on
  twiceburnt's go.
- Older lines stay as their own branches. Fix an older XION release on its
  line; don't rewrite history to do it.

## Bringing in an upstream tag: use the script

Run `scripts/xion-sync <upstream-tag>` from a checkout of this repository.
Try `--dry-run` first; it builds and tests in a throwaway worktree and
pushes nothing.

- **Same line** (for example v0.70.5 → v0.70.6): the script branches from
  `xion/v0.70.x`, merges the upstream tag with a signed merge commit, builds,
  runs the x/wasm tests and pushes `sync/<tag>`. Fix conflicts or the stack
  in further commits on that branch, never by rewriting it.
- **New line** (for example v0.71.0): the script starts from the upstream tag,
  cherry-picks the previous line's stack (`--from-line`, default
  `origin/xion/main`) and prints `git range-diff` so you can see how each
  patch changed. Update `FORK.md` for patches that changed or were dropped.
- **Land the PR with a merge commit.** Never squash or rebase a sync PR: a
  squash drops the upstream tag from the line's history and the next merge
  replays upstream's commits and conflicts. The PR goes to the line branch,
  or to `xion/main` when the line is the current one; then fast-forward the
  line branch to `xion/main`.
- Update `FORK.md` in the same PR: the stack table and the tag mapping.

## Tags

- Cut a tag only on a line branch's head, with
  `scripts/xion-sync --tag <upstream-tag>`. It checks that the upstream tag
  is in the line, builds and tests, picks the next free `-xion.N`, signs and
  annotates the tag, pushes it and waits for proxy.golang.org to resolve it.
  Nothing may pin a tag before the proxy resolves it.
- Never tag a commit that sits only on a work branch, `main`, or an
  integration branch.
- Never move, delete or re-create a pushed tag. The Go module proxy and sum
  database cache a version the first time anyone fetches it, a moved tag
  breaks every `go.sum` that records it, and a deleted name can never be
  reused. If a tag is wrong, fix it with the next `-xion.N`.
- Tags and merges to `xion/main` are twiceburnt's call unless he says
  otherwise.
