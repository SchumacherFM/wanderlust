GOFMT=gofmt -s
GOFILES=$(wildcard *.go **/*.go)
PRERELEASE=stable
WLHOMEDIR=${GOPATH}/src/github.com/SchumacherFM/wanderlust
default: build

format:
	${GOFMT} -w ${GOFILES}

run:
	go run main.go run

build:
	mkdir -p build
	cd github.com/SchumacherFM/go.gzrice/gzrice && go build
	mv github.com/SchumacherFM/go.gzrice/gzrice/gzrice ${WLHOMEDIR}/
	#cd ${WLHOMEDIR}
	cd github.com/GeertJohan/go.rice/rice && go build
	mv github.com/GeertJohan/go.rice/rice/rice ${WLHOMEDIR}/
	#cd ${WLHOMEDIR}
	# build with gzip support
	./gzrice --import-path ./picnic embed-go
	# no gzip support
	./rice --import-path ./github.com/HouzuoGuo/tiedot/webcp/ embed-go
	#goxc -c=.goxc.json -pr="$(PRERELEASE)" -d ./build
	go build -a -v

clean:
	find . -name "*.rice-box.go" -delete
	rm -f ./gzrice ./rice ./wanderlust

.PHONY: build
