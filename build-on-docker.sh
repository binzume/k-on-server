#!/bin/sh
APP=k-on-server
GOOS=linux
DIR=$(cd $(dirname $0); pwd)

mkdir -p _go
docker run --rm -v "$DIR":/usr/src/$APP -v $DIR/_go:/go -w /usr/src/$APP -e CGO_ENABLED=0 -e GOOS=$GOOS golang:1.8 bash -c "go get; go build"

# build image
cd $DIR
mkdir -p docker
cp -r $APP static docker
docker rmi $APP:latest
docker build -t $APP:latest docker

