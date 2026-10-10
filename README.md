# template-go

Minimal Go template with hooks via [prek](https://prek.j178.dev).

```
cmd/template-go/   main package (one dir per binary, keep it thin)
internal/          private packages
tools/check/       quality gates run by the hooks (coverage, deadcode, mutation)
```

```sh
prek install       # install pre-commit + pre-push hooks
go run ./cmd/template-go
go test ./...
prek run -a                          # pre-commit hooks
prek run -a --hook-stage pre-push    # pre-push hooks
```

- pre-commit: `go mod tidy`, `go fix`, `gofumpt -extra -w`, `go vet`
- pre-push: `golangci-lint run --fix`, `gofumpt -extra -w`, `go vet`, `deadcode -test`, `go test -race` with coverage gate

`gofumpt`, `deadcode` and `gremlins` are pinned in `go.mod` as `tool` dependencies (`go tool ...`); `golangci-lint` must be on `PATH`.

## Optional Nix shell

Installing the tools directly on NixOS (or another OS) remains supported. The
optional shell supports x86_64/aarch64 Linux and Apple Silicon macOS. With Nix
flakes enabled:

```sh
nix develop --command "$SHELL"       # keep your shell, aliases and prompt
# Or: nix develop                   # use Nix's default Bash shell
```

Then run the normal development commands in this README. Each template has its own
`flake.nix` and `flake.lock`, using `github:NixOS/nixpkgs/nixpkgs-unstable`.
To refresh to the latest unstable packages, run `nix flake update nixpkgs`,
then leave and re-enter the shell. Commit the updated lock file with your project.

The shell supplies `go_latest`, golangci-lint, and a C compiler for race tests.
The tools declared in `go.mod` remain managed by Go.

## Complexity

`golangci-lint` enforces a maximum cyclomatic complexity of 10 (`cyclop`) and
cognitive complexity of 15 (`gocognit`) per function, including tests and tools.
Both run in the existing pre-push hook and CI. Thresholds live in `.golangci.yml`.
Refactor functions that exceed them; these checks do not have automatic fixes.

## Coverage

The pre-push hook (also run in CI) runs `go run ./tools/check coverage -min 80`: `go test -race` plus a gate on statement coverage of `./internal/...`. Change `-min` in `prek.toml`.

```sh
go tool cover -func=coverage.out     # per-function summary (after a push or the hook)
go tool cover -html=coverage.out     # open HTML report
```

## Property-based testing

[rapid](https://github.com/flyingmutant/rapid) generates inputs and shrinks failures to a minimal case; see `TestHelloProperties` in `internal/greet/greet_test.go`. It runs with the normal `go test`.

```go
rapid.Check(t, func(t *rapid.T) {
	xs := rapid.SliceOf(rapid.Int()).Draw(t, "xs")
	// assert a property that holds for every xs
})
```

## Mutation testing

[gremlins](https://gremlins.dev), on demand and weekly in CI (`.github/workflows/mutation.yml`); fails if any mutant survives:

```sh
prek run --hook-stage manual gremlins --all-files    # report in gremlins.json
```

Caveats (gremlins v0.6.0):

- `--threshold-efficacy` is ignored (exit code is always 0), so `tools/check mutation` checks `mutants_lived` in `gremlins.json`.
- Mutants that don't compile (e.g. `+` → `-` on strings) are reported as KILLED, which inflates efficacy. Look at the LIVED lines, not the percentage.
