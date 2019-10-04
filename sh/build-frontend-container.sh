#!/bin/bash


# гасим бд
docker-compose down

# удаляем файлы бд
sudo rm -rf pgdata 

# компилируем. линкуем статически под линукс
# env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a  .

# ключи нужны для компиляции sqlite3 при CGO_ENABLED=1
env CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' .

# build a docker image 
docker build -t rgru/auth-proxy:latest -f Dockerfile-frontend . 

# push the docker image 
docker login
docker push rgru/auth-proxy:latest


# копируем docker-compose-frontend.yml и 
mkdir ../auth-latest
cp docker-compose-frontend.yml ../auth-latest/docker-compose.yml
