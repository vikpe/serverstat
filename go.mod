module github.com/vikpe/serverstat

go 1.18

require (
	github.com/go-playground/validator/v10 v10.11.0
	github.com/goccy/go-json v0.9.10
	github.com/ssoroka/slice v0.0.0-20220402005549-78f0cea3df8b
	github.com/stretchr/testify v1.7.2
	github.com/valyala/fastjson v1.6.3
	github.com/vikpe/udpclient v0.1.3
	github.com/vikpe/udphelper v0.1.3
	golang.org/x/exp v0.0.0-20220713135740-79cabaa25d75
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v0.1.3
	v0.1.2
	v0.1.1
	v0.1.0
)
