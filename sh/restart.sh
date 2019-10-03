#!/bin/bash

echo 'гасим бд'
docker-compose down

echo 'удаляем файлы бд'
sudo rm -rf pgdata 

echo 'поднимаем бд'
docker-compose up -d
sleep 2

echo 'запускаем приложение'
go run main.go -serve 4400 -env=dev




