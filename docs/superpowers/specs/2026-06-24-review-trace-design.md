# Bot-PR Review — Pseudo-Version Compare + Traceability — Design

## Problem

For `kairos-io/go-ukify#31` (a renovate **digest bump** of `github.com/foxboron/go-uefi` from pseudo-version `v0.0.0-…fab4fdf2f2f3` to `v0.0.0-…d29549a44f29`) the review never checked the upstream change, and there was no way to see *why*. Two root issues:

1. **Pseudo-versions break the compare.** `CompareDiff` is called with `"v"+version`, so for a pseudo-version it asks GitHub to compare refs `v0.0.0-20241017…` — which don't exist → 404 → silent degrade (no upstream diff). The trailing 12 hex chars of a pseudo-version *are* the commit SHA, and GitHub compare works with SHAs.
2. **No traceability.** A review records only the verdict/reasoning/summary, so "didn't check upstream" is invisible — no record of which bumps were found, which resolved, what compare ref was used, or why a fetch was skipped.

## Goals

- For a pseudo-version bump, use the embedded **commit SHA** as the compare ref so the upstream diff is actually fetched.
- Record and surface a **trace of what the agent did** (per bump: resolved repo, compare ref, fetched/failed/skipped + reason; plus context size) in the PR comment (collapsible) and the dashboard.

## Design

### Pseudo-version compare ref (`internal/review`)
`func compareRef(version string) string` (pure):
- Pseudo-version (`…-<14-digit-timestamp>-<≥12 hex>` suffix) → return the embedded commit SHA (`regexp` `-\d{14}-([0-9a-f]{12,})$`).
- Otherwise → `"v" + version` (the release tag).
`Run` uses `compareRef(b.From)` / `compareRef(b.To)` for `CompareDiff` instead of `"v"+b.From` / `"v"+b.To`.

### Trace (`state.PRReview.Trace []string`)
`Run` builds an ordered, human-readable trace during context assembly:
- No bumps: `"no go.mod dependency bumps parsed from the PR diff"`.
- Per bump:
  - unresolvable module: `"<mod> <from>→<to>: module not resolvable to a GitHub repo (skipped)"`.
  - compare failure/empty: `"<mod> <from>→<to>: compare <baseRef>...<headRef> failed: <err> (no upstream diff)"`.
  - success: `"<mod> <from>→<to>: compare <baseRef>...<headRef> ✓ <n> bytes"`.
- Final: `"context: <n> bytes"`.
`Trace` is set on the `PRReview` and persisted to `reviews.json`.

### Surfacing
- **PR comment** (`comment()`): after the verdict/reasoning/changes lines and before the marker, append a collapsible block:
  ```
  <details><summary>review trace</summary>

  ```
  <trace line 1>
  <trace line 2>
  ```

  </details>
  ```
  Hidden by default on the PR, expandable to see exactly what the agent did. Still upserted (one comment, edited in place).
- **Dashboard** "🔎 Bot-PR reviews": under each review row, render the trace lines as indented sub-items (markdown) / a small block (HTML, escaped). Omitted when empty; deterministic.

## Out of scope

- Resolving non-GitHub / vanity modules (still degrade, but now the trace says so).
- Monorepo subdir tags (`sub/vX.Y.Z`) — the trace will show the compare failure; a future improvement can try subdir tag forms.
- Including the full upstream diff in the comment (only the trace + the model's `changesSummary`).

## Testing

- `compareRef`: `0.33.0`→`v0.33.0`; `0.0.0-20241017190036-fab4fdf2f2f3`→`fab4fdf2f2f3`; `1.2.3-0.20240101000000-abcdef123456`→`abcdef123456`.
- `Run` (fakes): a pseudo-version bump calls `CompareDiff` with the SHA refs (not `v…`); trace records the success line with bytes; an unresolvable module / compare error produces the right trace line and still assesses; no-bumps produces the "no bumps" trace line; `Trace` is on the `PRReview`.
- comment: contains the `<details>review trace</details>` block with the trace lines; still upserted.
- render: trace sub-lines shown under a review; deterministic; no raw 64-hex finding id (the commit SHAs in the trace are intended/expected, not the finding id).
- Manual: re-run on go-ukify#31 → trace shows `compare fab4fdf2f2f3...d29549a44f29 ✓ <n> bytes` and the verdict reflects the real upstream change.
