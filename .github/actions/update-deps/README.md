# update-deps action

Have `nib` (driven by a self-hosted LocalAI) bump a repository's dependencies to
their latest versions, verify the build, and open a CI-triggering pull request.

## Usage

Copy [`examples/caller-workflow.yml`](examples/caller-workflow.yml) to
`.github/workflows/update-deps.yml` in the target repo. Minimal call:

    - uses: kairos-io/security/.github/actions/update-deps@main
      with:
        language: go
        token: ${{ steps.app-token.outputs.token }}

## Inputs

See [`action.yml`](action.yml). Key ones: `token` (required), `language`
(only `go` today), `model` (default `gemma-4-e2b-it`), `branch`, `base`,
`dry-run`.

## Token: use a GitHub App (not the built-in token)

The built-in `GITHUB_TOKEN` opens a PR but its checks never run (GitHub
suppresses workflow runs triggered by that token). Use a **GitHub App** token so
CI triggers, with no personal PAT:

1. Create an org GitHub App (`kairos-deps-bot`) with **Contents: write** and
   **Pull requests: write**.
2. Generate a private key; note the App ID.
3. Install the App on the target repos.
4. Store `DEPS_BOT_APP_ID` and `DEPS_BOT_APP_KEY` as org secrets.
5. Mint a token per run with `actions/create-github-app-token@v2` and pass it as
   `token` (see the example workflow).

## Behavior

- nib is the primary engine; if LocalAI can't load, a deterministic
  `go get -u ./... && go mod tidy` fallback runs so a PR still opens.
- Verify gate is `go build ./... && go vet ./...`. The repo's own CI runs tests
  on the PR.
- No PR is opened when there is no dependency change; the action fails (no PR)
  when the build can't be made to pass.
- An already-open PR on the same branch is force-updated instead of duplicated.
- The branch is pushed using the provided `token`, so updates to a reused PR
  trigger CI (git itself authenticates with the token, not just `gh`).
- LocalAI's binary and model weights are downloaded under the runner temp dir,
  never into the checkout, so they are never committed.
