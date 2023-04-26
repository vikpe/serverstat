module github.com/vikpe/serverstat

go 1.20

require (
	github.com/go-playground/validator/v10 v10.12.0
	github.com/goccy/go-json v0.10.2
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629
	github.com/ssoroka/slice v0.0.0-20220402005549-78f0cea3df8b
	github.com/stretchr/testify v1.8.2
	github.com/valyala/fastjson v1.6.4
	github.com/vikpe/qw-hub-api v0.2.2
	github.com/vikpe/udpclient v1.0.0
	github.com/vikpe/udphelper v1.0.1
	github.com/vikpe/wildcard v0.1.0
	golang.org/x/exp v0.0.0-20230425010034-47ecfdc1ba53
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/leodido/go-urn v1.2.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v0.1.3
	v0.1.2
	v0.1.1
	v0.1.0
)
