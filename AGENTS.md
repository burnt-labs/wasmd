# AGENTS.md

This repository is Burnt's fork of [CosmWasm/wasmd](https://github.com/CosmWasm/wasmd).
`CLAUDE.md` covers building, testing and the code layout. This file sets the
rules for `xion/main` and for our release tags.

These rules cover only our own tags, `vX.Y.Z-xion.N`. Upstream tags
(`vX.Y.Z`) and the `main` branch belong to upstream: don't move, re-tag or
push to them.

## `xion/main` stays in line with upstream

`xion/main` is the default branch. At every xion tag it is exactly an
upstream tag plus **at most one** commit of ours. That single signed commit
holds everything the fork carries: fork CI and docs, our features and fixes,
the SDK ports and the CI fixes they need. If a release needs nothing of
ours, the count is zero.

Between tags, PRs may add commits on top of `xion/main`. They are folded
into the one fork commit when the next tag is cut.

The fork commit's message lists each part it carries and where that part
came from (PR or original commit), so nothing is lost in the squash.

## Cutting a xion tag

1. Start a branch from the upstream tag:
   `git fetch https://github.com/CosmWasm/wasmd '+refs/tags/vX.Y.Z:refs/tags/vX.Y.Z'`,
   then `git switch -c work/<name> vX.Y.Z`.
2. Build the one fork commit: the previous fork commit plus everything
   merged into `xion/main` since, adapted to the new base. Build, test and
   get CI green on the branch.
3. Check the shape: `git rev-list --count vX.Y.Z..HEAD` must print `0` or
   `1`.
4. Move `xion/main` to that commit. This replaces the branch, so it needs a
   force-push, and the force-push needs the owner's explicit go for that
   release. Push a backup first:
   `git push origin <old head>:refs/heads/backup/xion-main-<short sha>`, then
   `git push --force-with-lease=xion/main:<old head> origin HEAD:xion/main`.
5. Tag the new `xion/main` head, signed and annotated:
   `git tag -s -m vX.Y.Z-xion.N vX.Y.Z-xion.N`. Tag only a commit that is
   `xion/main`'s head at that moment, never a commit that sits on a work
   branch alone.
6. Check that proxy.golang.org resolves the new version before anything
   pins it.

Once `xion/main` moves on, the previous tag's commit stays reachable only
through its tag, which is another reason never to delete one.

## Never move, delete or re-create a pushed tag

The Go module proxy and sum database cache a version the first time anyone
fetches it. A moved tag breaks every `go.sum` that already records it. If a
tag is wrong, fix it with the next `-xion.N`.
