
FROM    golang:1.9-alpine

RUN     apk add -U curl git bash

ARG     DEP_TAG=v0.3.2
RUN     go get -d github.com/golang/dep/cmd/dep && \
        cd /go/src/github.com/golang/dep && \
        git checkout -q "$DEP_TAG" && \
        go build -v -o /usr/bin/dep ./cmd/dep && \
        rm -rf /go/src/* /go/pkg/* /go/bin/*

WORKDIR /go/src/github.com/seriousben/goldendir
ENV CGO_ENABLED=0
