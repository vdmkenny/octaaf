FROM golang:1.10

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR $GOPATH