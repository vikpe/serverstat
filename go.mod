module github.com/vikpe/serverstat

go 1.23

toolchain go1.23.3

require (
	github.com/go-playground/validator/v10 v10.23.0
	github.com/goccy/go-json v0.10.3
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629
	github.com/samber/lo v1.47.0
	github.com/stretchr/testify v1.9.0
	github.com/valyala/fastjson v1.6.4
	github.com/vikpe/qw-hub-api v0.8.0
	github.com/vikpe/udpclient v1.0.0
	github.com/vikpe/udphelper v1.0.1
	github.com/vikpe/wildcard v0.1.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.6 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	golang.org/x/crypto v0.29.0 // indirect
	golang.org/x/net v0.31.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v0.1.3
	v0.1.2
	v0.1.1
	v0.1.0
)
