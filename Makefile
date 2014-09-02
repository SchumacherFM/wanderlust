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
	cd gzrice && go build
	cd ..
	cd ./github.com/GeertJohan/go.rice/rice && go build
	cd ../../../
	# build with gzip support
	./gzrice/gzrice --import-path ./picnic/ embed-go
	# no gzip support
	./github.com/GeertJohan/go.rice/rice/rice --import-path ./github.com/HouzuoGuo/tiedot/webcp/ embed-go
	#goxc -c=.goxc.json -pr="$(PRERELEASE)" -d ./build
	go build -a -v

clean:
	find . -name "*.rice-box.go" -delete
	rm -f ./gzrice/gzrice ./github.com/GeertJohan/go.rice/rice/rice

.PHONY: build
