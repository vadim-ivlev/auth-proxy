#!/bin/bash

echo 'Создает объекты базы данных'
migrate -source=file://migrations/ -database sqlite3://./auth.db up
