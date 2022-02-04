#!/bin/bash

# Компиляция приложения с публичной частью. Т.е может зарегистрироваться любой пользователь.

echo "Кросскомпиляция на компьютере разработчика auth-proxy:public"
env GOOS=linux GOARCH=amd64 go build -tags=jsoniter .


echo "build a docker image auth-proxy:public"
docker build -t registry.rgwork.ru:5050/masterback/auth-proxy/auth-proxy:public -f Dockerfile-public . 

echo "push the docker image" 
# docker login registry.rgwork.ru:5050
# docker push registry.rgwork.ru:5050/masterback/auth-proxy/auth-proxy:public

