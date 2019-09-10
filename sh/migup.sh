#!/bin/bash

echo 'Создает объекты базы данных'
migrate -source=file://migrations/ -database postgres://root:root@localhost:5432/auth?sslmode=disable up
