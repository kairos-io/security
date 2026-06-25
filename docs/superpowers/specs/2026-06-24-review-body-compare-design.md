# Bot-PR Review — Upstream Diffs from PR-Body Compare Links — Design

## Problem

`kairos-io/AuroraBoot#566` bumps `lucide-react 0.468.0→0.577.0` — an **npm** dep (no `go.mod`). Our review only parses Go module bumps, so it fetched no upstream source diff and the trace misleadingly said "no go.mod bumps parsed". Non-Go bumps (npm, Docker, Actions) are equally important to check.

But the renovate/dependabot **PR body already embeds `github.com/<owner>/<repo>/compare/<base>...<head>` links** for the bumped dependency (renovate having already resolved the correct repo *and* tags, including monorepos like `lucide-icons/lucide`). We feed the body to the model (so its changelog is already seen), but we don't fetch those compares.

## Goal

Fetch upstream **source diffs for any ecosystem** by parsing the compare links out of the PR body and `CompareDiff`-ing them — without building per-ecosystem module→repo→tag resolvers. Keep the Go `go.mod` path as a complement; record everything in the trace.

## Design

### Pure helpers (`internal/review/deps.go`)
- `type CompareRef struct{ Repo, Base, Head, Label string }`.
- `parseCompareURLs(body string) []CompareRef` — regex over the body for `github.com/<owner>/<repo>/compare/<base>...<head>` (matches `redirect.github.com` too, since it contains `github.com`). Base/head taken **verbatim** (renovate already used the repo's real tag form — `0.576.0`, `v1.2.3`, etc.). Returns each match with `Label = "<repo> <base>..<head> (PR body)"`.
  - Regex: `github\.com/([\w.\-]+)/([\w.\-]+)/compare/([\w.\-+/@]+?)\.\.\.([\w.\-+/@]+)` — non-greedy base stops at the literal `...`, head runs to the first non-token char (`)`, `>`, space, `"`). Versions never contain three consecutive dots, so the split is unambiguous.
- `compareTargets(diff []byte, body string) []CompareRef` — unify both sources, deduped + capped:
  - From `parseBumps(diff)`: for each Go bump whose `moduleRepo` resolves → `CompareRef{repo, compareRef(from), compareRef(to), "<module> <from>→<to>"}` (Go path keeps `compareRef`'s v-tag/SHA logic).
  - From `parseCompareURLs(body)`.
  - Dedup by `repo|base|head`; cap at `maxCompares` (5).

### `review.Run`
Replace the inline Go-bump loop with `compareTargets(diff, pr.Body)`:
- If empty → trace `"no upstream comparisons available (no go.mod bumps or compare links in the PR body)"`.
- For each target (until the context cap): `gh.CompareDiff(t.Repo, t.Base, t.Head)`; on error/empty → trace `"<label>: compare <base>...<head> failed/empty (no upstream diff)"`; else cap to `maxBumpDiff`, append `Upstream <label> (<base>...<head>):\n<diff>` to the context and trace `"<label>: compare <base>...<head> ✓ <n> bytes"`.
- The PR body (changelog) is still prepended to the context as before; the PR diff still appended. Idempotency (HeadSHA), MaxPerRun cap, comment upsert, dry-run, and degrade behavior are unchanged.

## Out of scope

- npm-registry / module-resolution code (the body's compare links replace it).
- Following non-GitHub compare hosts (only `github.com` compares are fetched; others ignored).
- De-noising a huge monorepo compare beyond the byte cap (capped; the body changelog carries the semantic summary).

## Testing

- `parseCompareURLs`: a renovate body with `…/lucide-icons/lucide/compare/0.576.0...0.577.0` (and a `redirect.github.com` variant, and a markdown-link `(…compare/a...b)`) → the right `{repo,base,head}` set, deduped; a body with no compares → empty.
- `compareTargets`: a Go-bump diff → Go CompareRef via moduleRepo/compareRef; an npm body with compare links → those refs; both present → deduped + capped at 5; verbatim base/head for body links (no `v` munging).
- `review.Run` (fakes): an npm PR (no go.mod bump) with body compare links calls `CompareDiff` with the body's repo/base/head, records the ✓ trace line, and assembles the upstream diff into the context; a PR with neither records the "no upstream comparisons available" line and still assesses; cap/idempotency/upsert/dry-run preserved.
- Manual: re-run on AuroraBoot#566 → trace shows `lucide-icons/lucide 0.576.0..0.577.0 (PR body): compare 0.576.0...0.577.0 ✓ N bytes` and the verdict reflects the real lucide change.
