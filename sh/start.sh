#!/bin/bash

echo 'поднимаем бд'
docker-compose up -d

echo 'спим 2 секунды'
sleep 2

echo 'запускаем приложение'
go run main.go -serve 4400 -env=dev

