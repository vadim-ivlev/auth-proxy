#!/bin/bash



echo 'гасим бд'
docker-compose down -v

echo 'удаляем db postgress'
sudo rm -rf pgdata 

echo 'удаляем db SQLite'
rm auth.db

