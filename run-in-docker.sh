#!/bin/sh

CURRENT=$(cd "$(dirname "$0")" && pwd)
docker volume create cfn-mackerel-macro-cache > /dev/null 2>&1
docker run --rm -it \
    -e GO111MODULE=on \
    -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=0 \
    -v cfn-mackerel-macro-cache:/go/pkg/mod \
    -v "$CURRENT":/go/src/github.com/shogo82148/cfn-mackerel-macro \
    -w /go/src/github.com/shogo82148/cfn-mackerel-macro golang:1.16.4 "$@"
