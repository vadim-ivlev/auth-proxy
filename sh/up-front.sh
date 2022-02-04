#!/bin/bash

echo 'поднимаем докер для front-end разработчиков'
docker-compose -f docker-compose-front.yml up -d && docker logs -f auth-proxy-front 
