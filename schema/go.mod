module github.com/olive-io/bpmn/schema

go 1.18

require (
	github.com/json-iterator/go v1.1.12
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Bad imports are sometimes causing attempts to pull that code.
// This makes the error more explicit.
replace (
	github.com/olive-io/bpmn => ./FORBIDDEN_DEPENDENCY
	github.com/olive-io/bpmn/schema => ./FORBIDDEN_DEPENDENCY
)
