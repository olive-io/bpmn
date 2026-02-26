# AGENTS.md

Guidance for coding agents working in `github.com/olive-io/bpmn/v2`.

## Repo Snapshot

- Primary language: Go.
- Main module: `github.com/olive-io/bpmn/v2`.
- Nested module: `schema` (`github.com/olive-io/bpmn/schema`).
- CI Go version: 1.22.
- Primary CI workflow: `.github/workflows/main.yml`.
- Lint in CI: `golint -set_exit_status ./..`.
- Tests in CI: root module tests plus schema module tests.

## Rules Files

- No `.cursorrules` file found.
- No `.cursor/rules/` directory found.
- No `.github/copilot-instructions.md` file found.
- If any of these are added later, treat them as higher-priority local agent instructions.

## Setup And Tooling

- Use Go 1.22.
- Confirm module mode is active (`go env GOMOD`).
- Install `golint` when linting: `go install golang.org/x/lint/golint@latest`.
- `go vet` is part of normal validation.
- `goimports` is used during schema generation flow.
- `changelog` CLI is used by the Makefile `build` target.
- `saxon-he` is referenced by `schema` `go:generate` directives.

## Build, Lint, Test Commands

Run from repo root unless noted.

- Full test pass (root module): `go test -v ./...`
- Schema module tests from root: `go test -v ./schema`
- Full vet pass: `go vet ./...`
- Lint (same style as CI): `golint -set_exit_status ./..`
- Makefile test target (runs vet first): `make test`
- Coverage/bench target from Makefile: `make test-coverage`
- Vendor dependencies: `make vendor`

### Single Test Commands (important)

- Single test function in one package:
  - `go test -v ./pkg/event -run '^TestForwardEvent$'`
- Single test by partial name:
  - `go test -v ./... -run 'TestUserTask'`
- Single test with count override (avoid cache):
  - `go test -v ./pkg/data -run '^TestContainer$' -count=1`
- Single subtest name:
  - `go test -v ./pkg/expression/expr -run 'TestExpr/invalid_input'`
- Single schema-module test (run inside schema module):
  - `cd schema && go test -v ./... -run '^TestParseSample$'`

Preferred approach for targeted debugging:

- Narrow package first, then add `-run` regex.
- Use `-count=1` for flaky or timing-sensitive tests.
- Keep `-v` enabled for trace-heavy tests in this repository.

## Generation Commands

- Schema generation entrypoint: `make generate`
- Make target behavior: enters `schema`, runs `go generate ./...`, then `goimports -w` on generated files.
- Generated files include `schema/schema_generated.go` and `schema/schema_di_generated.go`.
- Do not hand-edit generated schema files.

## Code Style: High-Level

- Follow existing repository patterns over generic Go style debates.
- Keep changes minimal and scoped.
- Preserve public API behavior unless task explicitly requires API changes.
- Prefer composable option-based configuration rather than adding many constructor variants.

## Formatting And Imports

- Always run `gofmt` on changed Go files.
- Use `goimports` when import blocks change significantly.
- Keep import groups in standard Go layout:
  - standard library
  - blank line
  - third-party/internal imports
- Avoid unused imports; CI and tooling should remain clean.

## Naming Conventions Used Here

- Exported names use `CamelCase`; unexported names use `camelCase`.
- Interface names frequently use `I` prefix (for example `ITracer`, `IEvent`, `IGenerator`).
- In this codebase, keep that `I*` interface style when extending existing packages.
- Constructors commonly follow:
  - `MakeX(...)` returns a value
  - `NewX(...)` returns a pointer
- Option functions typically follow `WithX(...)` naming.
- Many schema-generated struct members use `*Field` suffix; preserve that style inside `schema` package.

## Types And Data Modeling

- Use concrete types where possible; use interfaces at boundaries.
- This repository already uses `any` in maps and payload plumbing; match local usage.
- Preserve pointer-plus-presence patterns (`value, present := ...`) used across schema APIs.
- Avoid changing method signatures from value receivers to pointer receivers (or reverse) unless needed.

## Error Handling

- Return errors; do not introduce panics in normal control flow.
- Wrap errors with context using `%w` when propagating.
- Prefer structured local error types where package already defines them (`pkg/errors`).
- Keep error messages concise and action-oriented.
- For tests and examples, existing code sometimes uses `log.Fatalf`; for new production code, prefer error returns.

## Concurrency And Context

- Respect existing `context.Context` plumbing.
- Do not create background goroutines without clear lifecycle handling.
- When touching channel-based tracing/event flows, preserve unsubscribe/cleanup behavior.
- Keep lock usage consistent with current `sync.RWMutex` patterns.

## Testing Conventions

- Test files use `*_test.go` and mostly `TestXxx` naming.
- `testify/assert` is used widely; keep assertions consistent with local file style.
- Some tests use external package style (`package bpmn_test`) for public API validation.
- Fixture BPMN files live under `testdata` and are often loaded via embedded FS helpers.
- Prefer deterministic tests; avoid sleeping when possible.

## File-Specific Guardrails

- Treat files marked generated as read-only unless regenerating.
- Preserve Apache license header blocks on edited files that already include them.
- Keep build tags intact (for example `pkg/clock/host_generic.go`).
- Avoid broad refactors across root and `schema` modules in one change unless explicitly requested.

## Validation Checklist Before Finishing

- Run `go test -v ./...`.
- Run `go test -v ./schema`.
- Run `go vet ./...`.
- Run `golint -set_exit_status ./..` if lint-sensitive changes were made.
- If schema or generated code changed, run `make generate` and re-run tests.

## Commit/PR Hygiene For Agents

- Keep commits focused and atomic.
- Mention affected packages and behavior changes in commit messages.
- Note any intentionally skipped validation commands in PR/handoff text.
- Include exact single-test command(s) used when fixing a specific failing test.

## Practical Tips For Future Agents

- Root `go test ./...` does not replace schema module testing in CI; run both.
- If a test touches parsing/model elements, check both runtime and `schema` package impacts.
- Prefer incremental edits and fast package-level test runs before full-suite execution.
