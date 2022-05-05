# serverstat [![Go Reference](https://pkg.go.dev/badge/github.com/vikpe/serverstat.svg)](https://pkg.go.dev/github.com/vikpe/serverstat) [![Test](https://github.com/vikpe/serverstat/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/vikpe/serverstat/actions/workflows/test.yml) [![codecov](https://codecov.io/gh/vikpe/serverstat/branch/main/graph/badge.svg?token=nW6fiGr7hJ)](https://codecov.io/gh/vikpe/serverstat)

> Fetch info from QuakeWorld servers

## Development

### Tools

* **gow** (run go command on file change): `go install github.com/mitranim/gow@latest`

### Tests

Run tests on file change and write coverage data.

```shell
gow -c test ./... --cover -coverprofile coverage.out
```

## See also

* [masterstat](https://github.com/vikpe/masterstat)
* [masterstat-cli](https://github.com/vikpe/masterstat-cli)
