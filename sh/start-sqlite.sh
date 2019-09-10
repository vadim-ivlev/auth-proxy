#!/bin/bash


echo 'поднимаем бд'
docker-compose up -d

echo 'спим 2 секунды'
sleep 2

echo 'запускаем приложение с SQLite'
go run main.go -serve 4000 -env=dev -sqlite

