#!/bin/bash


# гасим бд
docker-compose down

# удаляем файлы бд
sudo rm -rf pgdata 

# компилируем. линкуем статически под линукс
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a  .

# build a docker image 
docker build -t rgru/auth-proxy:bare -f Dockerfile-frontend-bare . 

# push the docker image 
docker login
docker push rgru/auth-proxy:bare


# копируем docker-compose-frontend.yml и 
mkdir ../auth-bare
cp docker-compose-frontend-bare.yml ../auth-bare/docker-compose.yml

