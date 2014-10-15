# Dockerfile to cross compile boot2docker-cli

FROM golang:1.3-cross

ADD . /go/src/github.com/SchumacherFM/wanderlust
WORKDIR /go/src/github.com/SchumacherFM/wanderlust

# Download (but not install) dependencies
RUN go get -d -v ./...

CMD ["make", "all"]
