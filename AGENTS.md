# AGENTS.md

This repository is Burnt's fork of [CosmWasm/wasmd](https://github.com/CosmWasm/wasmd).
`CLAUDE.md` covers building, testing and the code layout. This file sets the
rules for our release tags.

## Xion release tags

These rules cover only our own tags, `vX.Y.Z-xion.N`. Upstream tags
(`vX.Y.Z`) and the `main` branch are upstream's: don't move, re-tag or push
to them.

1. **Cut tags only on `xion/main`.** `xion/main` is the default branch. A
   xion tag must point at a commit in `xion/main`'s history. Never tag a
   commit that sits only on a work branch, a PR branch or a deleted branch.
   Before pushing a tag, run
   `git merge-base --is-ancestor <commit> origin/xion/main`; it must succeed.
2. **A tag is the upstream tag plus our cherry-picks.** `vX.Y.Z-xion.N` is
   upstream `vX.Y.Z`, plus only the commits that release needs, one commit per
   change. Run `git log --oneline vX.Y.Z..vX.Y.Z-xion.N`. It must list only
   those commits: no other upstream `main` commits, no squashed bundles, and
   no merge commits.
3. **Never move, delete or re-create a pushed tag.** The Go module proxy and
   sum database cache a version the first time anyone fetches it. A moved tag
   breaks every `go.sum` that already records it. If a tag is wrong, fix it
   with the next `-xion.N`.

## Moving to a new upstream tag

Rewriting `xion/main` needs a force-push, so build the new line beside it:

1. Branch from the upstream tag:
   `git fetch https://github.com/CosmWasm/wasmd '+refs/tags/vX.Y.Z:refs/tags/vX.Y.Z'`,
   then `git switch -c work/<name> vX.Y.Z`.
2. Cherry-pick our commits from the previous xion tag, adapting each one
   where the new base needs it. Keep them as separate signed commits.
3. Open a PR into `xion/main` and merge it with **Create a merge commit**.
   Squash or rebase merges rewrite the line, so its commits would no longer
   sit on top of the upstream tag.
4. After the merge, tag the tip of the merged line, not the merge commit.
   That commit satisfies both rules: it is in `xion/main`'s history, and
   `git log vX.Y.Z..<tag>` lists only our commits. Then check rule 1 and rule
   2 before pushing the tag.

A fix on an existing line works the same way: a branch from the current
tag's commit, the extra commits, a merge-commit PR into `xion/main`, and the
next `-xion.N` on the branch tip.

A line for an older upstream tag than the one `xion/main` carries, such as
a patch release on the previous minor line, is merged with
`git merge -s ours`. The merge brings the line into `xion/main`'s history
without changing `xion/main`'s files. Then tag that line's tip as usual.
