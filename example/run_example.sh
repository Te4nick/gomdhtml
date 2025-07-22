#!/bin/bash
name="gomdhtml-nginx"

go run ../cmd

if [ "$(docker ps -aq -f name=$name)" ]; then
    docker rm -f $name # cleanup
fi

docker build -t $name .
docker run -d -p 80:80 --name $name $name