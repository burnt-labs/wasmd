# AGENTS.md

This repository is Burnt's fork of [CosmWasm/wasmd](https://github.com/CosmWasm/wasmd).
`CLAUDE.md` covers building, testing and the code layout. This file sets the
process for `xion/main` and our release tags.

These rules cover only our own tags, `vX.Y.Z-xion.N`. Upstream tags
(`vX.Y.Z`) and the `main` branch belong to upstream: don't move, re-tag or
push to them.

## Process for every xion tag

1. Branch from `xion/main`.
2. Merge the upstream tag
   (`git fetch https://github.com/CosmWasm/wasmd '+refs/tags/vX.Y.Z:refs/tags/vX.Y.Z'`,
   then `git merge vX.Y.Z`).
3. Fix conflicts, add features, and modify as needed.
4. Squash it to a single commit in the PR to `xion/main`.
5. Create our xion tag for that version, `vX.Y.Z-xion.N`, on the commit
   the PR landed. Make it signed and annotated
   (`git tag -s -m vX.Y.Z-xion.N vX.Y.Z-xion.N <commit>`), and check that
   proxy.golang.org resolves it before anything pins it.

## Tag rules

- Cut tags only on `xion/main`. A tagged commit must be in `xion/main`'s
  history; `git merge-base --is-ancestor <commit> origin/xion/main` must
  succeed. Never tag a commit that sits only on a work branch.
- Never move, delete or re-create a pushed tag. The Go module proxy and sum
  database cache a version the first time anyone fetches it, and a moved
  tag breaks every `go.sum` that already records it. If a tag is wrong, fix
  it with the next `-xion.N`.

## Note: squash merges and the upstream ancestry

A squash merge in step 4 leaves the upstream tag out of `xion/main`'s git
history. The next step 2 then finds an old merge base and replays upstream
commits that are already in the tree, so expect conflicts in the files we
change. A merge-commit PR keeps the upstream tag as an ancestor, so the
next upstream merge only brings in what is new.
