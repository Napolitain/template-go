# template-go

Minimal Go template with hooks via [prek](https://prek.j178.dev).

```
cmd/template-go/   main package (one dir per binary)
internal/          private packages
```

```sh
prek install       # install pre-commit + pre-push hooks
go run ./cmd/template-go
go test ./...
prek run -a                          # pre-commit hooks
prek run -a --hook-stage pre-push    # pre-push hooks
```

- pre-commit: `go mod tidy`, `go fix`, `gofumpt -extra -w`, `go vet`
- pre-push: `golangci-lint run --fix`, `gofumpt -extra -w`, `go vet`, `deadcode -test`

`gofumpt` and `deadcode` are pinned in `go.mod` as `tool` dependencies (`go tool ...`); `golangci-lint` must be on `PATH`.

Coverage:

```sh
go test -race -covermode=atomic -coverprofile=coverage.out ./...
go tool cover -func=coverage.out     # per-function summary
go tool cover -html=coverage.out     # open HTML report
```
