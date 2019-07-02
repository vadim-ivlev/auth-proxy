#!/bin/bash

# гасим бд
docker-compose down

# удаляем файлы бд
sudo rm -rf pgdata 

# поднимаем бд
docker-compose up -d
sleep 2

# запускаем приложение
go run main.go -serve 4000 -env=dev

# ./migup.sh


