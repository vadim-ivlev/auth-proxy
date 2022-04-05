#!/bin/bash

echo 'запускаем приложение'
# admin_url=http://localhost:5000/?url=http://localhost:4400 \
admin_url=http://localhost:4400/admin/?url=http://localhost:4400 \
go run main.go -port 4400 -config=./configs/app.env.dev -pgconfig=./configs/db.env.dev


