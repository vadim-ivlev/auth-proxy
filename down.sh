#!/bin/bash



# гасим бд
docker-compose down

# удаляем файлы бд
sudo rm -rf pgdata 
