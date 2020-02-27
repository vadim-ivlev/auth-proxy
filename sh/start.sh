#!/bin/bash

echo 'поднимаем бд'
docker-compose up -d

echo 'спим 5 секунд'
sleep 5

echo 'запускаем приложение'
go run main.go -serve 4400 -env=dev

