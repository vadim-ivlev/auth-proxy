#!/bin/bash

echo "Кросскомпиляция на компьютере разработчика auth-proxy:dev"
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags=jsoniter .


echo "build a docker image auth-proxy:dev"

image=registry.rgwork.ru:5050/masterback/auth-proxy/auth-proxy:dev

docker build -t $image -f Dockerfile-front . 
echo "push the docker image" 
docker login registry.rgwork.ru:5050
docker push $image


