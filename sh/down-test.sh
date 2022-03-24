#!/bin/bash

echo 'гасим докер для front-end разработчиков'
docker-compose -f docker-compose-test.yml down -v --remove-orphans 
