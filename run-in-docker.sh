#!/bin/sh

CURRENT=$(cd "$(dirname "$0")" && pwd)
docker run --rm -it \
    -e GO111MODULE=on \
    -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=0 \
    -v "$CURRENT/.mod":/go/pkg/mod \
    -v "$CURRENT":/go/src/github.com/shogo82148/mackerel-cloudwatch-forwarder \
    -w /go/src/github.com/shogo82148/mackerel-cloudwatch-forwarder golang:1.16.0 "$@"
