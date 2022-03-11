#!/bin/bash

echo 'порождаем auth_proxy_network'
docker network create auth_proxy_network

echo 'поднимаем бд'
docker-compose -f docker-compose-dev.yml up -d

# echo 'спим 2 секунды'
# sleep 2


echo 'запускаем приложение'
go run main.go -port 4400 -config=./configs/app.env.dev -pgconfig=./configs/db.env.dev


