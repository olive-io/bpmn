module github.com/olive-io/bpmn/v2

go 1.22

replace github.com/olive-io/bpmn/schema => ./schema

require (
	github.com/bits-and-blooms/bitset v1.22.0
	github.com/bytedance/sonic v1.13.3
	github.com/expr-lang/expr v1.17.5
	github.com/hashicorp/go-multierror v1.1.1
	github.com/muyo/sno v1.2.1
	github.com/olive-io/bpmn/schema v1.4.2
	github.com/qri-io/iso8601 v0.1.0
	github.com/stretchr/testify v1.10.0
	golang.org/x/sys v0.30.0
)

require (
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	golang.org/x/arch v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
