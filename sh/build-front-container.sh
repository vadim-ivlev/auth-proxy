#!/bin/bash

echo "Кросскомпиляция на компьютере разработчика auth-proxy:front"
env GOOS=linux GOARCH=amd64 go build -tags=jsoniter .


echo "build a docker image auth-proxy:front"
docker build -t registry.rgwork.ru:5050/masterback/auth-proxy/auth-proxy:front  -f Dockerfile-front . 

echo "push the docker image" 
# docker login registry.rgwork.ru:5050
# docker push registry.rgwork.ru:5050/masterback/auth-proxy/auth-proxy:front

