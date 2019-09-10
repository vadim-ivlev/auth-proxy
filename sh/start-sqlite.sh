#!/bin/bash


echo 'запускаем приложение с SQLite'
go run main.go -serve 4000 -env=dev -sqlite

