#!/bin/bash

echo 'порождаем auth_proxy_network'
docker network create auth_proxy_network

echo 'поднимаем бд'
docker-compose -f docker-compose-dev.yml up -d

# echo 'спим 2 секунды'
# sleep 2


echo 'запускаем приложение'
# admin_url=http://localhost:5000/?url=http://localhost:4400 \
admin_url=http://localhost:4400/admin/?url=http://localhost:4400 \
go run main.go -port 4400 -config=./configs/app.env.dev -pgconfig=./configs/db.env.dev


