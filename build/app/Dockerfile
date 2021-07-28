FROM golang:1.16.2

ENV GO111MODULE=off
RUN go get -u -d g.phga.de/phga/gott/...
WORKDIR /go/src/g.phga.de/phga/gott/cmd/gott
RUN go build
CMD ["./gott"]