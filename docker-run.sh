#!/bin/sh
IMAGE=binzume/k-on-server
APP=$IMAGE
PORT=14201

APP=${APP#*/}
APP=${APP%:*}
DIR=$(cd $(dirname $0); pwd)

mkdir -p $DIR/data

docker rm -f $APP 2>/dev/null
docker run -d -u `id -u` --restart="always" --name $APP -v $DIR/data:/data -p $PORT:8080 $IMAGE

