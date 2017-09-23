#!/bin/sh
APP=k-on-server
GOOS=linux
DIR=$(cd $(dirname $0); pwd)

mkdir -p _go
U=`id -u`:`id -g`
docker run --rm -u $U -v "$DIR":/usr/src/$APP -v $DIR/_go:/go -w /usr/src/$APP -e CGO_ENABLED=0 -e GOOS=$GOOS golang:1.9 \
	sh -c "go get -d && go build" || exit 1

# build image
cd $DIR
mkdir -p docker
cp -r $APP static docker
docker rmi $APP:latest
docker build -t $APP:latest docker
