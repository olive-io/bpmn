NAME=bpmn
IMAGE_NAME=olive-io/$(NAME)
ROOT=github.com/olive-io/bpmn

all: build

vendor:
	go mod vendor

test-coverage:
	go test ./... -bench=. -coverage

lint:
	golint -set_exit_status ./..

vet:
	go vet ./...

test: vet
	go test -v ./...

build:
	mkdir -p _output
	changelog --last --output _output/CHANGELOG.md

generate:
	cd schema && go generate ./... && goimports -w schema_generated.go schema_generated_test.go schema_di_generated.go

clean:
	rm -fr ./_output

