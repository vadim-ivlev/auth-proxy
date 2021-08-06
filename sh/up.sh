#!/bin/bash

echo 'поднимаем бд'
docker-compose -f docker-compose-dev.yml up -d

echo 'спим 5 секунд'
sleep 5

echo 'запускаем приложение'
go run main.go -port 4400 -env=dev

