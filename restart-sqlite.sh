#!/bin/bash

# гасим бд
docker-compose down

# удаляем файлы бд
rm auth.db

# поднимаем бд
docker-compose up -d
sleep 2

# запускаем приложение
go run main.go -serve 4000 -env=dev -sqlite




