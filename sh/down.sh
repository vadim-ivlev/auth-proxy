#!/bin/bash



echo 'гасим бд'
docker-compose down

echo 'удаляем db postgress'
sudo rm -rf pgdata 

echo 'удаляем db SQLite'
rm auth.db

