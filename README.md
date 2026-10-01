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
