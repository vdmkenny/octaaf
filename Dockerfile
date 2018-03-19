FROM golang:alpine

RUN apk add --no-cache git

RUN go get -u github.com/golang/dep/cmd/dep

RUN adduser -D -s /bin/ash -G wheel  octaaf

USER octaaf

WORKDIR $GOPATH