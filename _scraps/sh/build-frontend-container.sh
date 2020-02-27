#!/bin/bash

# Если под Windows, добавляем команду sudo
if [[ "$OSTYPE" == "msys" ]]; then alias sudo=""; fi


# гасим бд
docker-compose down

# удаляем файлы бд
sudo rm -rf pgdata 


# build a docker image 
docker build -t rgru/auth-proxy:latest -f Dockerfile-frontend . 

# push the docker image 
docker login
docker push rgru/auth-proxy:latest


# копируем docker-compose-frontend.yml и 
mkdir ../auth
cp docker-compose-frontend.yml ../auth/docker-compose.yml
