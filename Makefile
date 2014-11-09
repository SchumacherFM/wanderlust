GOFMT=gofmt -s
GOFILES=$(wildcard *.go **/*.go)
WLHOMEDIR=${GOPATH}/src/github.com/SchumacherFM/wanderlust

VERSION := $(shell cat VERSION)
GITSHA1 := $(shell git rev-parse --short HEAD)
GOARCH := amd64
GOFLAGS := -ldflags "-X main.Version $(VERSION) -X main.GitSHA $(GITSHA1)"
PREFIX := wanderlust
DOCKER_IMAGE := wanderlust-golang
DOCKER_CONTAINER := wanderlust-build
DOCKER_SRC_PATH := /go/src/github.com/SchumacherFM/wanderlust


default: dockerbuild
	@true # stop from matching "%" later


# Build binaries in Docker container. The `|| true` hack is a temporary fix for
# https://github.com/dotcloud/docker/issues/3986
dockerbuild: clean
	cd github.com/SchumacherFM/go.gzrice/gzrice && go build
	mv github.com/SchumacherFM/go.gzrice/gzrice/gzrice ${WLHOMEDIR}/
	# build with gzip support before sending the build context to the Docker daemon
	./gzrice --import-path ./picnic embed-go
	docker build -t "$(DOCKER_IMAGE)" .
	docker run --name "$(DOCKER_CONTAINER)" "$(DOCKER_IMAGE)" 
	docker cp "$(DOCKER_CONTAINER)":"$(DOCKER_SRC_PATH)"/$(PREFIX)-$(VERSION)-darwin-$(GOARCH) . || true
	docker cp "$(DOCKER_CONTAINER)":"$(DOCKER_SRC_PATH)"/$(PREFIX)-$(VERSION)-linux-$(GOARCH) . || true
	docker cp "$(DOCKER_CONTAINER)":"$(DOCKER_SRC_PATH)"/$(PREFIX)-$(VERSION)-windows-$(GOARCH).exe . || true
	docker rm "$(DOCKER_CONTAINER)"


# Remove built binaries and Docker container. Silent errors if container not found.
clean:
	rm -f $(PREFIX)*
	docker rm "$(DOCKER_CONTAINER)" 2>/dev/null || true
	find . -name "*.rice-box.go" -delete
	rm -f ./gzrice


all: darwin linux windows
	@true # stop "all" from matching "%" later


# Native Go build per OS/ARCH combo.
%:
	GOOS=$@ GOARCH=$(GOARCH) go build $(GOFLAGS) -o $(PREFIX)-$(VERSION)-$@-$(GOARCH)$(if $(filter windows, $@),.exe)


# This binary will be installed at $GOBIN or $GOPATH/bin. Requires proper
# $GOPATH setup AND the location of the source directory in $GOPATH.
goinstall:
	go install $(GOFLAGS)

format:
	${GOFMT} -w ${GOFILES}

test:
	go test -v --bench=. ./helpers/
	go test -v --bench=. ./picnic/middleware/
	go test -v --bench=. ./picnic/
	go test -v --bench=. ./provisionerApi/
	go test -v --bench=. ./rucksack/
	go test -v --bench=. ./provisioners/sitemap/
	go test -v --bench=. ./provisioners/textarea/

testAll:
	go test -v ./...

vet:
	go vet ./...

run:
	go run main.go run --rf=./build/test.db --ppd=./build/
