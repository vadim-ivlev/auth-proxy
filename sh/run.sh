#!/bin/bash

echo 'запускаем приложение'
go run main.go -port 4400 -config=./configs/app.env -pgconfig=./configs/db.env


