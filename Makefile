GOFMT=gofmt -s
GOFILES=$(wildcard *.go **/*.go)
PRERELEASE=stable

default: build

format:
	${GOFMT} -w ${GOFILES}

run:
	go run main.go run

build:
	mkdir -p build
	# @todo
	#cd github.com/GeertJohan/go.rice/rice && go build
	#cd ../../../picnic && rice embed-go
	goxc -c=.goxc.json -pr="$(PRERELEASE)" -d ./build

.PHONY: build
