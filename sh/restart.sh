#!/bin/bash

# Если под Windows, добавляем команду sudo
if [[ "$OSTYPE" == "msys" ]]; then alias sudo=""; fi


echo 'гасим бд'
docker-compose down

echo 'удаляем файлы бд'
sudo rm -rf pgdata 

echo 'поднимаем бд'
docker-compose up -d
sleep 5

echo 'запускаем приложение'
go run main.go -port 4400 -env=dev




