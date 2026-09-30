# template-go

Minimal Go template with hooks via [prek](https://prek.j178.dev).

```sh
prek install       # install pre-commit + pre-push hooks
go run .
go test ./...
prek run -a                          # pre-commit hooks
prek run -a --hook-stage pre-push    # pre-push hooks
```

- pre-commit: `go mod tidy`, `go fix`, `gofumpt -extra -w`, `go vet`
- pre-push: `golangci-lint run --fix`, `gofumpt -extra -w`, `go vet`, `deadcode -test`

`gofumpt` and `deadcode` are pinned in `go.mod` as `tool` dependencies (`go tool ...`); `golangci-lint` must be on `PATH`.
