#!/bin/bash

# поднимаем бд
docker-compose up -d
sleep 2

# запускаем приложение
go run main.go -serve 4000 -env=dev -sqlite

