module github.com/vikpe/serverstat

go 1.18

require (
	github.com/ssoroka/slice v0.0.0-20220402005549-78f0cea3df8b
	github.com/stretchr/testify v1.7.1
	github.com/vikpe/udpclient v0.1.3
	github.com/vikpe/udphelper v0.1.3
	golang.org/x/exp v0.0.0-20220518171630-0b5c67f07fdf
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v0.1.3
	v0.1.2
	v0.1.1
	v0.1.0
)
