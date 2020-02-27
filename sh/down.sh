#!/bin/bash

# Если под Windows, добавляем команду sudo
if [[ "$OSTYPE" == "msys" ]]; then alias sudo=""; fi


echo 'гасим бд'
docker-compose down -v

echo 'удаляем db postgress'
sudo rm -rf pgdata 

echo 'удаляем db SQLite'
rm auth.db

