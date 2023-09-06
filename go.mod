module github.com/stretchr/testify

// This should match the minimum supported version that is tested in
// .github/workflows/main.yml
go 1.17

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/stretchr/objx v0.5.0
	golang.org/x/term v0.12.0
	gopkg.in/yaml.v3 v3.0.1
)

require golang.org/x/sys v0.12.0 // indirect
