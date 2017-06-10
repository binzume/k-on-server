#!/bin/sh
APP=k-on-server
GOOS=linux
docker run --rm -v "$PWD":/usr/src/$APP -v $PWD/build_src:/go/src -w /usr/src/$APP -e CGO_ENABLED=0 -e GOOS=$GOOS golang:1.8 bash -c "go get; go build"

