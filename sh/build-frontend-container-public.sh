#!/bin/bash

# Компиляция приложения с публичной частью. Т.е может зарегистрироваться любой пользователь.

echo "Кросскомпиляция на компьютере разработчика auth-proxy:devpublic"
env GOOS=linux GOARCH=amd64 go build -tags=jsoniter .


echo "build a docker image auth-proxy:devpublic"
# docker build -t rgru/auth-proxy:latest -f Dockerfile-frontend . 
docker build -t registry.rgwork.ru:5050/masterback/auth-proxy/auth-proxy:devpublic --build-arg NODE_ENV=public -f Dockerfile-frontend . 

echo "push the docker image" 
# docker login
# docker push rgru/auth-proxy:latest
docker login registry.rgwork.ru:5050
docker push registry.rgwork.ru:5050/masterback/auth-proxy/auth-proxy:devpublic


# Если под Windows, добавляем команду sudo
# if [[ "$OSTYPE" == "msys" ]]; then alias sudo=""; fi

# компилируем. линкуем статически под линукс
# env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a  .

# ключи нужны для компиляции sqlite3 при CGO_ENABLED=1
# echo "Кросскомпиляция на компьютере разработчика"
# env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -ldflags '-linkmode external -extldflags "-static"' .

# echo "Кросскомпиляция в докере. Сделано чтобы компилировать под windows."
# # docker run --rm -it -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e CGO_ENABLED=1 -e GOOS=linux golang:1.15.6 go build -a -ldflags '-linkmode external -extldflags "-static"'
# # CGO_ENABLED=0  сделан после обновления докера до version 20.10.6, build 370c289. Высыпалось исключение
# # docker run --rm -it -v "$PWD":/usr/src/auth-proxy -w /usr/src/auth-proxy -e CGO_ENABLED=0 -e GOOS=linux -e GOARCH=amd64 golang:1.15.6 go build -a 
# docker run --rm -it -v "$PWD":/usr/src/auth-proxy -w /usr/src/auth-proxy -e CGO_ENABLED=0 -e GOOS=linux -e GOARCH=amd64 golang:1.16.4 go build -a 


# echo "build a docker image"
# docker build -t rgru/auth-proxy:latest -f Dockerfile-frontend . 

# echo "push the docker image" 
# docker login
# docker push rgru/auth-proxy:latest


# echo "копируем docker-compose-frontend.yml в ../auth-proxy-front" 
# mkdir ../auth-proxy-front
# cp docker-compose-frontend.yml ../auth-proxy-front/docker-compose.yml
