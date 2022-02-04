#!/bin/bash

echo 'гасим докер для front-end разработчиков'
docker-compose -f docker-compose-front.yml down -v --remove-orphans 
