module github.com/wallester/testify

// This should match the minimum supported version that is tested in
// .github/workflows/main.yml
go 1.17

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/stretchr/objx v0.5.2
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/stretchr/testify => github.com/wallester/testify v0.0.0-20240704140627-c236cee6c133
